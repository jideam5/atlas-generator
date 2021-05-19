/*
@Time : 5/18/21 7:05 PM
@Author : jideam
@File : fit
@Software: GoLand
*/
package maxrect

import "image"

// (maxInt, maxInt) - worst
// (0, 0) - best
type ScoreFn func(free image.Rectangle, size image.Point) (int, int)

func FitsInto(free image.Rectangle, size image.Point) bool {
	return size.X <= free.Dx() && size.Y <= free.Dy()
}

func FindBest(frees []image.Rectangle, size image.Point, rotate bool, score ScoreFn) (image.Rectangle, int, int) {
	best := image.Rectangle{}
	bestMajor, bestMinor := maxInt, maxInt

	rotSize := image.Point{X: size.Y, Y: size.X}

	for _, free := range frees {
		if FitsInto(free, size) {
			major, minor := score(free, size)

			if major < bestMajor || (major == bestMajor && minor < bestMinor) {
				best.Min = free.Min
				best.Max = free.Min.Add(size)
				bestMajor, bestMinor = major, minor
			}
		}

		if rotate && FitsInto(free, rotSize) {
			major, minor := score(free, rotSize)
			if major < bestMajor || (major == bestMajor && minor < bestMinor) {
				best.Min = free.Min
				best.Max = free.Min.Add(rotSize)
				bestMajor, bestMinor = major, minor
			}
		}
	}

	return best, bestMajor, bestMinor
}

func CommonInterval(a0, a1, b0, b1 int) int {
	if a1 < b0 || b1 < a0 {
		return 0
	}

	return min(a1, b1) - max(a0, b0)
}
