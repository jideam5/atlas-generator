/*
@Time : 5/18/21 6:58 PM
@Author : jideam
@File : best
@Software: GoLand
*/
package maxrect

import "image"

type Rule byte

const (
	Automatic Rule = iota
	ShortSide
	LongSide
	BottomLeft
	Area
	ContactPoint
)

func ParseRule(s string) Rule {
	switch s {
	case "short-side":
		return ShortSide
	case "long-side":
		return LongSide
	case "bottom-left":
		return BottomLeft
	case "area":
		return Area
	case "contact-point":
		return ContactPoint
	default:
		return Automatic
	}
}

func (ctx *Context) SetRule(rule Rule) {
	switch rule {
	case Automatic:
		ctx.Score = ctx.ContactPoint
	case ShortSide:
		ctx.Score = ctx.ShortSide
	case LongSide:
		ctx.Score = ctx.LongSide
	case BottomLeft:
		ctx.Score = ctx.BottomLeft
	case Area:
		ctx.Score = ctx.Area
	case ContactPoint:
		ctx.Score = ctx.ContactPoint
	default:
		panic("unknown rule")
	}
}

func (ctx *Context) ShortSide(size image.Point) (image.Rectangle, int, int) {
	return FindBest(ctx.Free, size, ctx.Rotate,
		func(free image.Rectangle, size image.Point) (int, int) {
			leftX, leftY := free.Dx()-size.X, free.Dy()-size.Y
			if leftX < leftY {
				return leftX, leftY
			} else {
				return leftY, leftX
			}
		})
}

func (ctx *Context) LongSide(size image.Point) (image.Rectangle, int, int) {
	return FindBest(ctx.Free, size, ctx.Rotate,
		func(free image.Rectangle, size image.Point) (int, int) {
			leftX, leftY := free.Dx()-size.X, free.Dy()-size.Y
			if leftX > leftY {
				return leftX, leftY
			} else {
				return leftY, leftX
			}
		})
}

func (ctx *Context) BottomLeft(size image.Point) (image.Rectangle, int, int) {
	return FindBest(ctx.Free, size, ctx.Rotate,
		func(free image.Rectangle, size image.Point) (int, int) {
			return free.Min.Y + size.Y, free.Min.X
		})
}

func (ctx *Context) Area(size image.Point) (image.Rectangle, int, int) {
	area := size.X * size.Y

	return FindBest(ctx.Free, size, ctx.Rotate,
		func(free image.Rectangle, size image.Point) (int, int) {
			areaFit := free.Dx()*free.Dy() - area
			leftX := free.Dx() - size.X
			leftY := free.Dy() - size.Y
			shortFit := min(leftX, leftY)
			return areaFit, shortFit
		})
}

func (ctx *Context) ContactPoint(size image.Point) (image.Rectangle, int, int) {
	return FindBest(ctx.Free, size, ctx.Rotate,
		func(free image.Rectangle, size image.Point) (int, int) {
			contact := maxInt
			target := image.Rectangle{Min: free.Min, Max: free.Min.Add(size)}

			if target.Min.X == 0 || target.Max.X == ctx.Size.X {
				contact -= size.X
			}

			for _, used := range ctx.Used {
				if used.Min.X == target.Max.X || used.Max.X == target.Min.X {
					contact -= CommonInterval(used.Min.Y, used.Max.Y, target.Min.Y, target.Max.Y)
				}

				if used.Min.Y == target.Max.Y || used.Max.Y == target.Min.Y {
					contact -= CommonInterval(used.Min.X, used.Max.X, target.Min.X, target.Max.X)
				}
			}

			return contact, contact
		})
}
