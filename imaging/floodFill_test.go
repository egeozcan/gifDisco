package imaging

import (
	"image"
	"image/color"
	"math"
	"testing"

	"github.com/egeozcan/gifDisco/colorutils"
)

func TestSmoothFloodFill(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	for y := 2; y < 8; y++ {
		for x := 2; x < 8; x++ {
			img.Set(x, y, color.White)
		}
	}

	result := SmoothFloodFill(img, 5, 5, color.RGBA{255, 0, 0, 255}, 100).(*image.RGBA)

	// Verify filled area
	for y := 2; y < 8; y++ {
		for x := 2; x < 8; x++ {
			r, g, b, _ := result.At(x, y).RGBA()
			if r>>8 != 255 || g>>8 != 0 || b>>8 != 0 {
				t.Errorf("Pixel at (%d,%d) not filled correctly", x, y)
			}
		}
	}

	// Verify boundaries
	boundaries := []image.Point{{1, 1}, {8, 8}, {2, 1}, {7, 8}}
	for _, pt := range boundaries {
		r, g, b, _ := result.At(pt.X, pt.Y).RGBA()
		if r != 0 || g != 0 || b != 0 {
			t.Errorf("Boundary pixel at (%d,%d) was modified", pt.X, pt.Y)
		}
	}
}

func TestColorBlending(t *testing.T) {
	start := color.RGBA{100, 100, 100, 255}
	end := color.RGBA{200, 200, 200, 255}

	blended := blendColors(start, end, 0.5).(color.RGBA)

	expected := color.RGBA{150, 150, 150, 255}
	if blended != expected {
		t.Errorf("Expected %v, got %v", expected, blended)
	}

	// Corrected distance calculation
	dist := colorutils.ColorDistanceSquared(start, end)
	expectedDist := uint64(math.Pow(100*257, 2) * 3) // (200-100) in 8-bit = 100*257 in 16-bit
	if dist != expectedDist {
		t.Errorf("Expected distance %d, got %d", expectedDist, dist)
	}
}
