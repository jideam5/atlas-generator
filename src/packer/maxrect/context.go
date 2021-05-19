/*
@Time : 5/18/21 7:00 PM
@Author : jideam
@File : context
@Software: GoLand
*/
package maxrect

import "image"

const maxInt = int(^uint(0) >> 1)

type Context struct {
	Size   image.Point
	Rotate bool
	Score  func(size image.Point) (image.Rectangle, int, int)

	Used []image.Rectangle
	Free []image.Rectangle

	DebugPlace func(rect image.Rectangle)
}

func New(size image.Point) *Context {
	ctx := &Context{}
	ctx.Size = size
	ctx.Score = ctx.BottomLeft
	ctx.Free = append(ctx.Free, image.Rect(0, 0, size.X, size.Y))
	return ctx
}

func (ctx *Context) Add(size image.Point) (rect image.Rectangle, ok bool) {
	node, _, _ := ctx.Score(size)
	if node == image.ZR {
		return node, false
	}

	ctx.PlaceRect(node)
	return node, true
}

func (ctx *Context) Adds(size ...image.Point) (rects []image.Rectangle, ok bool) {
	rs := make([]image.Rectangle, len(size))

	toAdd := make([]int, len(size))

	for i := range toAdd {
		toAdd[i] = i
	}

	for len(toAdd) > 0 {
		best, bestMajor, bestMinor, bestIndex := image.Rectangle{}, maxInt, maxInt, -1

		for pi, i := range toAdd {
			node, major, minor := ctx.Score(size[i])
			if node == image.ZR {
				continue
			}

			if major < bestMajor || (major == bestMajor && minor < bestMinor) {
				best = node
				bestIndex = pi
				bestMajor = major
				bestMinor = minor
			}
		}

		if bestIndex == -1 {
			return rs, false
		}

		ctx.PlaceRect(best)
		rs[toAdd[bestIndex]] = best
		toAdd = append(toAdd[:bestIndex], toAdd[bestIndex+1:]...)
	}

	return rs, true
}

func (ctx *Context) PlaceRect(rect image.Rectangle) {
	n := len(ctx.Free)

	for i := 0; i < n; i++ {
		if ctx.splitFree(ctx.Free[i], &rect) {
			ctx.Free = removeAt(ctx.Free, i)
			i--
			n--
		}
	}

	ctx.mergeFree()
	ctx.removeSmall(3)

	if ctx.DebugPlace != nil {
		ctx.DebugPlace(rect)
	}

	ctx.Used = append(ctx.Used, rect)
}

func removeAt(rects []image.Rectangle, i int) []image.Rectangle {
	return append(rects[:i], rects[i+1:]...)
}

func (ctx *Context) splitFree(free image.Rectangle, used *image.Rectangle) bool {
	// 如果没有重叠直接返回
	if !free.Overlaps(*used) {
		return false
	}

	if used.Min.X < free.Max.X && used.Max.X > free.Min.X {

		if used.Min.Y > free.Min.Y && used.Min.Y < free.Max.Y {
			split := free
			split.Max.Y = used.Min.Y
			ctx.Free = append(ctx.Free, split)
		}

		if used.Max.Y < free.Max.Y {
			split := free
			split.Min.Y = used.Max.Y
			ctx.Free = append(ctx.Free, split)
		}

	}

	if used.Min.Y < free.Max.Y && used.Max.Y > free.Min.Y {

		if used.Min.X > free.Min.X && used.Min.X < free.Max.X {
			split := free
			split.Max.X = used.Min.X
			ctx.Free = append(ctx.Free, split)
		}

		if used.Max.X < free.Max.X {
			split := free
			split.Min.X = used.Max.X
			ctx.Free = append(ctx.Free, split)
		}
	}

	return true
}

func (ctx *Context) mergeFree() {
	for i := 0; i < len(ctx.Free); i++ {
		for k := i + 1; k < len(ctx.Free); k++ {
			if ctx.Free[i].In(ctx.Free[k]) {
				ctx.Free = removeAt(ctx.Free, i)
				i--
				break
			}

			if ctx.Free[k].In(ctx.Free[i]) {
				ctx.Free = removeAt(ctx.Free, k)
				k--
			}
		}
	}
}

func (ctx *Context) removeSmall(limit int) {
	for i := 0; i < len(ctx.Free); i++ {
		sz := ctx.Free[i].Size()
		if sz.X <= limit || sz.Y <= limit {
			ctx.Free = removeAt(ctx.Free, i)
			i--
		}
	}
}
