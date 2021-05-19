/*
@Time : 5/5/21 8:28 PM
@Author : jideam
@File : config
@Software: GoLand
*/
package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path"
)

var (
	Cnf Config
)

const (
	DotByte = byte(46)
)

func Initialize() error {
	if len(os.Args) < 2 {
		return errors.New("usage: project <json file>")
	}

	content, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		return err
	}

	return json.Unmarshal(content, &Cnf)
}

type Atlas struct {
	Width         string `json:"width"`
	Height        string `json:"height"`
	Size          int    `json:"size"`
	Quality       int    `json:"quality"`
	PixelFormat   string `json:"pixelFormat"`
	PowerOfTwo    bool   `json:"powerOfTwo"`
	TextureFormat string `json:"textureFormat"`
}

type Data struct {
	Format  string `json:"format"`
	Compact bool   `json:"compact"`
}

type Sprite struct {
	Width     string `json:"width"`
	Height    string `json:"height"`
	Size      int    `json:"size"`
	Rotation  bool   `json:"rotation"`
	Extrude   int    `json:"extrude"`
	Padding   int    `json:"padding"`
	CropAlpha bool   `json:"cropAlpha"`
}

type Config struct {
	InputDir    string   `json:"inputDir"`
	OutputDir   string   `json:"outputDir"`
	ResDir      string   `json:"resDir"`
	Force       bool     `json:"force"`
	IncludeList []string `json:"includeList"`
	ExcludeList []string `json:"excludeList"`
	ExtrudeList []string `json:"extrudeList"`
	Atlas       Atlas    `json:"atlas"`
	Data        Data     `json:"data"`
	Sprite      Sprite   `json:"sprite"`
}

func (c Config) String() string {
	bs, _ := json.Marshal(c)
	return string(bs)
}

func RecPath() string {
	return path.Join(Cnf.OutputDir, ".rec")
}

func IsCopyRes() bool {
	return Cnf.ResDir != ""
}

func IsExclude(absPath string) bool {
	for _, v := range Cnf.ExcludeList {
		if v == absPath {
			return true
		}
	}

	return false
}

func IsInclude(absPath string) bool {
	for _, v := range Cnf.IncludeList {
		if v == absPath {
			return true
		}
	}
	return false
}

func IsExtrude(absPath string) bool {
	for _, v := range Cnf.ExtrudeList {
		if v == absPath {
			return true
		}
	}
	return false
}
