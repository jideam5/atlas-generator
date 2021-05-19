/*
@Time : 5/19/21 12:13 PM
@Author : jideam
@File : Atlas
@Software: GoLand
*/
package atlas

import (
	"encoding/json"
	"os"
)

type Atlas struct {
	Meta   Meta             `json:"meta"`
	Frames map[string]Frame `json:"frames"`
}

func (a *Atlas) WriteFile(filename string) error {
	fd, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer fd.Close()

	enc := json.NewEncoder(fd)
	return enc.Encode(a)
}

type Rect struct {
	X int `json:"x"`
	Y int `json:"y"`
	W int `json:"w"`
	H int `json:"h"`
}

type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type Size struct {
	W int `json:"w"`
	H int `json:"h"`
}

type Meta struct {
	Image  string `json:"image"`
	Prefix string `json:"prefix"`
}

type Frame struct {
	Name             string `json:"-"`
	Rotated          bool   `json:"rotated"`
	Trimmed          bool   `json:"trimmed"`
	Frame            Rect   `json:"frame"`
	SpriteSourceSize Rect   `json:"spriteSourceSize"`
	SourceSize       Size   `json:"sourceSize"`
	Pivot            Point  `json:"-"`
}
