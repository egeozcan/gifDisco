package imaging

import (
	"image"
	"image/color"
	"image/draw"
	"math"

	"github.com/egeozcan/gifDisco/colorutils"
)

func SmoothFloodFill(src image.Image, x, y int, fill color.Color, tolerance float64) image.Image {
	bounds := src.Bounds()
	dst := image.NewRGBA(bounds)
	draw.Draw(dst, bounds, src, bounds.Min, draw.Src)

	origColor := src.At(x, y)
	toleranceSquared := uint64(tolerance * tolerance)
	visited := make([]bool, bounds.Dx()*bounds.Dy())
	queue := []image.Point{{X: x, Y: y}}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if !current.In(bounds) {
			continue
		}

		idx := (current.Y-bounds.Min.Y)*bounds.Dx() + (current.X - bounds.Min.X)
		if visited[idx] {
			continue
		}
		visited[idx] = true

		currColor := src.At(current.X, current.Y)
		distanceSquared := colorutils.ColorDistanceSquared(currColor, origColor)

		if distanceSquared > toleranceSquared {
			continue
		}

		blended := blendColors(currColor, fill, math.Sqrt(float64(distanceSquared))/tolerance)
		dst.Set(current.X, current.Y, blended)

		for _, d := range []image.Point{{-1, 0}, {1, 0}, {0, -1}, {0, 1}} {
			next := current.Add(d)
			if next.In(bounds) {
				nextIdx := (next.Y-bounds.Min.Y)*bounds.Dx() + (next.X - bounds.Min.X)
				if !visited[nextIdx] {
					queue = append(queue, next)
				}
			}
		}
	}

	return dst
}

func blendColors(curr, fill color.Color, t float64) color.Color {
	cr, cg, cb, ca := curr.RGBA()
	fr, fg, fb, fa := fill.RGBA()

	return color.RGBA{
		R: uint8(lerp(float64(fr), float64(cr), t) / 256),
		G: uint8(lerp(float64(fg), float64(cg), t) / 256),
		B: uint8(lerp(float64(fb), float64(cb), t) / 256),
		A: uint8(lerp(float64(fa), float64(ca), t) / 256),
	}
}

func lerp(a, b, t float64) float64 {
	return (1-t)*a + t*b
}
