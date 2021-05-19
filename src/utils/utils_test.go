/*
@Time : 5/6/21 12:39 PM
@Author : jideam
@File : utils_test
@Software: GoLand
*/
package utils

import (
	"fmt"
	"testing"
)

func TestHex2Uint(t *testing.T) {
	s := fmt.Sprintf("%08X", uint32(1000))
	fmt.Println("s:", s)
	i := Hex2Uint(s)
	fmt.Println("i:", i)
}
