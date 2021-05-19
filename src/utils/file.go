/*
@Time : 5/6/21 11:05 AM
@Author : jideam
@File : file
@Software: GoLand
*/
package utils

import (
	"atlasGen/src/config"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

func Copy2ResDir(relationPath string) {
	if !config.IsCopyRes() {
		return
	}

	srcPath := path.Join(config.Cnf.InputDir, relationPath)
	dstPath := path.Join(config.Cnf.ResDir, relationPath)
	srcInfo, err := os.Stat(srcPath)
	if err != nil {
		return
	}

	if srcInfo.IsDir() {
		_ = CopyDir(dstPath, srcPath)
		return
	}

	resDir, _ := filepath.Split(dstPath)
	_ = os.MkdirAll(resDir, 0755)

	_ = CopyFile(dstPath, srcPath)
}

func CopyFile(dst, src string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}

func CopyDir(dst, src string) error {
	var err error
	var fds []os.FileInfo
	var srcInfo os.FileInfo

	if srcInfo, err = os.Stat(src); err != nil {
		return err
	}

	if err = os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return err
	}

	if fds, err = ioutil.ReadDir(src); err != nil {
		return err
	}

	for _, fd := range fds {
		srcPath := path.Join(src, fd.Name())
		dstPath := path.Join(dst, fd.Name())

		if fd.IsDir() {
			if err = CopyDir(dstPath, srcPath); err != nil {
				fmt.Println(err)
			}
		} else {
			if err = CopyFile(dstPath, srcPath); err != nil {
				fmt.Println(err)
			}
		}
	}
	return nil
}
