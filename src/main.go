/*
@Time : 5/5/21 8:20 PM
@Author : jideam
@File : main
@Software: GoLand
*/
package main

import (
	"atlasGen/src/config"
	"atlasGen/src/packer"
	"atlasGen/src/record"
	"atlasGen/src/utils"
	"fmt"
	"io/ioutil"
	"path"
)

// 1. 解析json文件
// 2. 遍历inputDir，找出需要打包的目录
// 3. 打包目录
// 4. 更新记录
func main() {
	if err := config.Initialize(); err != nil {
		fmt.Println(err)
		return
	}

	record.ReadLast()

	rootDir, err := ioutil.ReadDir(config.Cnf.InputDir)
	if err != nil {
		return
	}

	record.AddDirectory(".")
	for _, v := range rootDir {
		if v.Name()[0] == config.DotByte {
			continue
		}

		if v.IsDir() {
			// 根目录下文件夹 稍后将会打包
			packer.AddPath(v.Name())
		} else {

			absPath := path.Join(config.Cnf.InputDir, v.Name())
			record.AddOther(utils.Crc32File(absPath), v.Name())
			utils.Copy2ResDir(v.Name())
			fmt.Println("INROOT ", v.Name())
		}
	}

	packer.Execute()

	// 写记录
	record.Write()

	// 完成
	fmt.Println("finish")
}
