/*
@Time : 5/19/21 12:19 PM
@Author : jideam
@File : box
@Software: GoLand
*/
package box

import (
	"atlasGen/src/packer/maxrect"
	"image"
	"image/draw"
	"sort"
)

type Box struct {
	Name    string
	Padding int
	Size    image.Point
	Place   *image.Rectangle
	Data    image.Image

	SrcSize image.Point
	Rotated bool
}

func Sort(boxes []*Box) {
	sort.Slice(boxes, func(i, j int) bool {
		a, b := boxes[i], boxes[j]
		if a.Size.X > b.Size.X {
			return true
		} else if a.Size.X < b.Size.X {
			return false
		}

		if a.Size.Y > b.Size.Y {
			return true
		} else if a.Size.Y < b.Size.Y {
			return false
		}

		return a.Name < b.Name
	})
}

func SizesFrom(boxes []*Box) []image.Point {
	sizes := make([]image.Point, 0, len(boxes))
	for _, b := range boxes {
		sizes = append(sizes, b.Size.Add(image.Point{X: b.Padding * 2, Y: b.Padding * 2}))
	}
	return sizes
}

func Place(maxCtxSize image.Point, rule maxrect.Rule, boxes []*Box) (ctxSize image.Point, ok bool) {
	sizes := SizesFrom(boxes)

	var rects []image.Rectangle
	ctxSize, rects, ok = minimizeFit(maxCtxSize, rule, sizes)
	for i := range rects {
		iBox := boxes[i]
		*iBox.Place = rects[i].Inset(iBox.Padding)
	}
	return
}

func minimizeFit(maxCtxSize image.Point, rule maxrect.Rule, sizes []image.Point) (ctxSize image.Point, rects []image.Rectangle, ok bool) {
	try := func(size image.Point) ([]image.Rectangle, bool) {
		ctx := maxrect.New(size)
		ctx.SetRule(rule)
		return ctx.Adds(sizes...)
	}

	ctxSize = maxCtxSize
	rects, ok = try(ctxSize)
	if !ok {
		return
	}

	shrunk, shrinkX, shrinkY := true, true, true
	for shrunk {
		shrunk = false
		if shrinkX {
			trySize := image.Point{X: ctxSize.X - 32, Y: ctxSize.Y}
			tryRects, tryOk := try(trySize)
			if tryOk {
				ctxSize = trySize
				rects = tryRects
				shrunk = true
			} else {
				shrinkX = false
			}
		}

		if shrinkY {
			trySize := image.Point{X: ctxSize.X, Y: ctxSize.Y - 32}
			tryRects, tryOk := try(trySize)
			if tryOk {
				ctxSize = trySize
				rects = tryRects
				shrunk = true
			} else {
				shrinkY = false
			}
		}
	}
	return
}

func (b *Box) Draw(dst draw.Image) {
	draw.Draw(dst, *b.Place, b.Data, image.Point{}, draw.Src)

	ex := (b.Size.X - b.SrcSize.X) / 2
	if ex > 0 {
		*b.Place = b.Place.Inset(ex)
	}
}

func (b *Box) Extrude(val int) {
	b.Size = b.Size.Add(image.Point{X: val * 2, Y: val * 2})

	dst := image.NewRGBA(image.Rectangle{Min: image.Point{}, Max: b.Size})
	rect := image.Rectangle{
		Min: image.Point{X: val, Y: val},
		Max: b.SrcSize.Add(image.Point{X: val, Y: val}),
	}
	draw.Draw(dst, rect, b.Data, image.Point{}, draw.Src)

	if b.Size.Y > val+1 {
		for i := val; i < b.Size.X-val; i++ {

			for j := 0; j < val; j++ {
				dst.Set(i, j, dst.At(i, val))
				dst.Set(i, b.Size.Y-1-j, dst.At(i, b.Size.Y-1-val))
			}

		}
	}

	if b.Size.X > val+1 {
		for i := val; i < b.Size.Y-val; i++ {
			for j := 0; j < val; j++ {
				dst.Set(j, i, dst.At(val, i))
				dst.Set(b.Size.X-1-j, i, dst.At(b.Size.X-1-val, i))
			}
		}
	}

	b.Data = dst
}
