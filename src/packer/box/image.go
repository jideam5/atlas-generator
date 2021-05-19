/*
@Time : 5/19/21 12:26 PM
@Author : jideam
@File : image
@Software: GoLand
*/
package box

import (
	"image"
	"os"
	"path"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

func LoadImage(filename string, padding int) (*Box, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	data, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}

	return &Box{
		Name:    path.Base(filename),
		Padding: padding,
		Size:    data.Bounds().Size(),
		Place:   &image.Rectangle{},
		Data:    data,
		SrcSize: data.Bounds().Size(),
	}, nil
}
