/*
@Time : 5/19/21 12:10 PM
@Author : jideam
@File : packer
@Software: GoLand
*/
package packer

import (
	"atlasGen/src/config"
	"atlasGen/src/packer/atlas"
	"atlasGen/src/packer/box"
	"atlasGen/src/packer/maxrect"
	"atlasGen/src/record"
	"atlasGen/src/utils"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

var (
	directories = make([]string, 0)
	boxes       = make([]*box.Box, 0)
)

func AddPath(filename string) {
	directories = append(directories, filename)
}

func Execute() {

	for {
		if len(directories) == 0 {
			break
		}

		v := directories[0]
		directories = directories[1:]
		pack(v)
	}
}

func pack(dirPath string) {
	// 清空
	boxes = boxes[0:0]

	record.AddDirectory(dirPath)
	fmt.Println("DIRECTORY", dirPath)
	absPath := path.Join(config.Cnf.InputDir, dirPath)

	if config.IsExclude(absPath) {
		fmt.Println("EXCLUDE", dirPath)
		utils.Copy2ResDir(dirPath)
		return
	}

	// 是否已在record中
	dirs, err := ioutil.ReadDir(absPath)
	if err != nil {
		return
	}
	for _, v := range dirs {
		if v.Name()[0] == config.DotByte {
			return
		}

		rPath := path.Join(dirPath, v.Name())
		if v.IsDir() {
			// 根目录下文件夹 稍后将会打包
			AddPath(rPath)
		} else {
			ProcessFile(rPath)
		}
	}

	// 查看是否需要打包
	if !needPack(dirPath) {
		fmt.Println("directory has not changes")
		return
	}

	fmt.Println("directory has been changed")

	// 排序
	box.Sort(boxes)

	maxSize := image.Point{X: config.Cnf.Sprite.Size, Y: config.Cnf.Sprite.Size}

	size, ok := box.Place(maxSize, maxrect.ParseRule("automatic"), boxes)
	if !ok {
		fmt.Println("image do not fit")
		return
	}

	dst := image.NewRGBA(image.Rectangle{Min: image.Point{}, Max: size})
	for _, b := range boxes {
		b.Draw(dst)
	}

	outPath := path.Join(config.Cnf.OutputDir, dirPath)
	pPath, _ := filepath.Split(outPath)
	_ = os.MkdirAll(pPath, 0777)
	outFilename := outPath + ".png"
	atlasFilename := outPath + ".json"

	if err := writeImage(dst, outFilename); err != nil {
		fmt.Println("write dst image error:", err)
		return
	}

	atl := &atlas.Atlas{}
	atl.Frames = make(map[string]atlas.Frame)
	atl.Meta.Image = path.Base(outFilename)
	atl.Meta.Prefix = dirPath + "/"
	for _, b := range boxes {
		frame := atlas.Frame{
			Name:    b.Name,
			Rotated: b.Rotated,
			Frame: atlas.Rect{
				X: b.Place.Min.X,
				Y: b.Place.Min.Y,
				W: b.Place.Max.X - b.Place.Min.X,
				H: b.Place.Max.Y - b.Place.Min.Y,
			},
			SpriteSourceSize: atlas.Rect{
				X: 0,
				Y: 0,
				W: b.Place.Max.X - b.Place.Min.X,
				H: b.Place.Max.Y - b.Place.Min.Y,
			},
			SourceSize: atlas.Size{
				W: b.Place.Max.X - b.Place.Min.X,
				H: b.Place.Max.Y - b.Place.Min.Y,
			},
		}

		atl.Frames[frame.Name] = frame
	}
	if err := atl.WriteFile(atlasFilename); err != nil {
		fmt.Println("failed to write data json:", err)
		return
	}
}

func needPack(dirPath string) bool {
	if len(boxes) == 0 {
		return false
	}

	rItem := record.Get(dirPath)
	if rItem == nil || config.Cnf.Force {
		return true
	}

	if len(rItem.Pictures) != len(boxes) {
		return true
	}

	for _, v := range boxes {
		if !utils.StrInArray(v.Name, rItem.Pictures) {
			return true
		}
	}

	return false
}

func ProcessFile(rPath string) {
	absPath := path.Join(config.Cnf.InputDir, rPath)
	crc := utils.Crc32File(absPath)
	m := "="
	if config.IsExclude(absPath) {

		if record.IsModified(rPath, crc) {
			utils.Copy2ResDir(rPath)
			m = "*"
		}
		fmt.Println(m, "EXCLUDE", rPath)
		record.AddOther(crc, rPath)
		return
	}

	iBox, err := box.LoadImage(absPath, config.Cnf.Sprite.Padding)
	if err != nil {
		// 不是图片
		if record.IsModified(rPath, crc) {
			utils.Copy2ResDir(rPath)
			m = "*"
		}
		fmt.Println(m, "NOTIMAGE", rPath)
		record.AddOther(crc, rPath)
		return
	}

	// 如果图片的size小于spriteSize则打包， 如果大于，则查看是否include，否则直接copy
	if iBox.Size.X > config.Cnf.Sprite.Size || iBox.Size.Y > config.Cnf.Sprite.Size {

		if !config.IsInclude(absPath) {
			if record.IsModified(rPath, crc) {
				utils.Copy2ResDir(rPath)
				m = "*"
			}
			fmt.Println(m, "OVERFLOW", rPath)
			record.AddOther(crc, rPath)
			return
		}

		record.AddPicture(crc, rPath)
		if record.IsModified(rPath, crc) {
			m = "*"
		}

		fmt.Println(m, "INCLUDE", rPath)
	} else {
		record.AddPicture(crc, rPath)
		if record.IsModified(rPath, crc) {
			m = "*"
		}
		fmt.Println(m, "LOAD", rPath)
	}

	// 查看是否需要拉伸
	if config.Cnf.Sprite.Extrude > 0 && config.IsExtrude(absPath) {
		fmt.Println("EXTRUDE:", iBox.Name)
		iBox.Extrude(config.Cnf.Sprite.Extrude)
	}

	boxes = append(boxes, iBox)
}

func writeImage(image *image.RGBA, filename string) error {
	writer, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer writer.Close()

	err = png.Encode(writer, image)
	if err != nil {
		return err
	}

	return nil
}
