/*
@Time : 5/5/21 8:59 PM
@Author : jideam
@File : crc
@Software: GoLand
*/
package utils

import (
	"hash/crc32"
	"io/ioutil"
)

func Crc32File(filename string) uint32 {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return 0
	}

	return crc32.ChecksumIEEE(content)
}
