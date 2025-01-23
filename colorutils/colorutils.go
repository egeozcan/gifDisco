package colorutils

import (
	"image/color"
	"math/rand"
)

func ColorDistanceSquared(c1, c2 color.Color) uint64 {
	r1, g1, b1, _ := c1.RGBA()
	r2, g2, b2, _ := c2.RGBA()
	dr := int64(r2) - int64(r1)
	dg := int64(g2) - int64(g1)
	db := int64(b2) - int64(b1)
	return uint64(dr*dr + dg*dg + db*db)
}

func RandomPastelColor() color.RGBA {
	return color.RGBA{
		R: uint8(rand.Intn(128) + 128),
		G: uint8(rand.Intn(128) + 128),
		B: uint8(rand.Intn(128) + 128),
		A: 255,
	}
}

func RandomDarkColor() color.RGBA {
	return color.RGBA{
		R: uint8(rand.Intn(128)),
		G: uint8(rand.Intn(128)),
		B: uint8(rand.Intn(128)),
		A: 255,
	}
}
