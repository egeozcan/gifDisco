package colorutils

import (
	"image/color"
	"testing"
)

func TestColorDistanceSquared(t *testing.T) {
	tests := []struct {
		name     string
		c1       color.Color
		c2       color.Color
		expected uint64
	}{
		{"Black to Black", color.RGBA{0, 0, 0, 255}, color.RGBA{0, 0, 0, 255}, 0},
		{"White to White", color.RGBA{255, 255, 255, 255}, color.RGBA{255, 255, 255, 255}, 0},
		{"Black to White", color.RGBA{0, 0, 0, 255}, color.RGBA{255, 255, 255, 255}, 0xffff * 0xffff * 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ColorDistanceSquared(tt.c1, tt.c2)
			if result != tt.expected {
				t.Errorf("Expected %d, got %d", tt.expected, result)
			}
		})
	}
}

func TestRandomColorRanges(t *testing.T) {
	// Test pastel colors
	for i := 0; i < 100; i++ {
		c := RandomPastelColor()
		if c.R < 128 || c.G < 128 || c.B < 128 {
			t.Errorf("Pastel color out of range: %v", c)
		}
	}

	// Test dark colors
	for i := 0; i < 100; i++ {
		c := RandomDarkColor()
		if c.R >= 128 || c.G >= 128 || c.B >= 128 {
			t.Errorf("Dark color out of range: %v", c)
		}
	}
}
