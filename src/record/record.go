/*
@Time : 5/5/21 8:47 PM
@Author : jideam
@File : record
@Software: GoLand
*/
package record

import (
	"atlasGen/src/config"
	"atlasGen/src/utils"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type Item struct {
	Typ      string
	Crc      uint32
	Name     string
	Pictures []string
}

var (
	lastMap = map[string]*Item{}
	records = make([]string, 0)
)

const (
	TypDirectory = "D"
	TypPicture   = "P"
	TypOther     = "R"
)

func AddDirectory(path string) {
	records = append(records, fmt.Sprintf("%s %s", TypDirectory, path))
}

func AddPicture(crc uint32, path string) {
	records = append(records, fmt.Sprintf("%s %08X %s", TypPicture, crc, path))
}

func AddOther(crc uint32, path string) {
	records = append(records, fmt.Sprintf("%s %08X %s", TypOther, crc, path))
}

func Write() {
	err := ioutil.WriteFile(config.RecPath(), []byte(strings.Join(records, "\n")), os.ModePerm)
	if err != nil {
		fmt.Println("failed to write record file:", err)
	}
}

func ReadLast() {
	if config.Cnf.Force {
		fmt.Println("Force Publication.")
		return
	}

	content, err := ioutil.ReadFile(config.RecPath())
	if err != nil {
		return
	}
	lastDir := ""
	for _, line := range strings.Split(string(content), "\n") {
		cells := strings.Split(line, " ")
		if len(cells) > 1 {
			switch cells[0] {
			case TypDirectory:
				name := strings.TrimSpace(cells[1])
				lastMap[name] = &Item{
					Typ:      TypDirectory,
					Crc:      0,
					Name:     name,
					Pictures: make([]string, 0),
				}
				lastDir = name
			case TypPicture:
				name := strings.TrimSpace(cells[2])
				lastMap[name] = &Item{
					Typ:  TypPicture,
					Crc:  utils.Hex2Uint(cells[1]),
					Name: name,
				}
				if lastDir != "" {
					if _, ok := lastMap[lastDir]; ok {
						lastMap[lastDir].Pictures = append(lastMap[lastDir].Pictures, name)
					}
				}
			case TypOther:
				name := strings.TrimSpace(cells[2])
				lastMap[name] = &Item{
					Typ:  TypOther,
					Crc:  utils.Hex2Uint(cells[1]),
					Name: name,
				}
			}
		}
	}

}

func IsModified(rPath string, crc uint32) bool {
	if v, ok := lastMap[rPath]; ok {
		if v.Crc == crc {
			return false
		}
	}
	return true
}

func Get(rPath string) *Item {
	if v, ok := lastMap[rPath]; ok {
		return v
	}
	return nil
}
