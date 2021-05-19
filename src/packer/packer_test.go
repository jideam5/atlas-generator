/*
@Time : 5/19/21 4:32 PM
@Author : jideam
@File : packer_test
@Software: GoLand
*/
package packer

import (
	"atlasGen/src/packer/box"
	"image"
	"image/draw"
	"testing"
)

func TestProcessFile(t *testing.T) {
	iBox, _ := box.LoadImage("/Users/jideam/workspace/gitlab/project/atlasGen/tmp/dist/chat/0.png", 1)

	iBox.Extrude(2)

	rect := iBox.Data.Bounds()
	dst := image.NewRGBA(rect)
	draw.Draw(dst, rect, iBox.Data, rect.Min, draw.Src)
	writeImage(dst, "/Users/jideam/workspace/gitlab/project/atlasGen/tmp/dist/chat/0-v.png")
}
