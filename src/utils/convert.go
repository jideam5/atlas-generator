/*
@Time : 5/6/21 12:38 PM
@Author : jideam
@File : convert
@Software: GoLand
*/
package utils

import "strconv"

func Hex2Uint(s string) uint32 {
	i, _ := strconv.ParseUint(s, 16, 32)
	return uint32(i)
}

func StrInArray(seed string, sArr []string) bool {
	for _, v := range sArr {
		if v == seed {
			return true
		}
	}
	return false
}
