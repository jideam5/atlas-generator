/*
@Time : 5/18/21 6:57 PM
@Author : jideam
@File : utils
@Software: GoLand
*/
package maxrect

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
