package main

import (
	"fmt"
	"image"
	"image/color"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"math"
	"math/rand"
	"os"
	"sync"
	"time"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

func colorDistance(c1, c2 color.Color) float64 {
	r1, g1, b1, _ := c1.RGBA()
	r2, g2, b2, _ := c2.RGBA()
	return math.Sqrt(float64((r2-r1)*(r2-r1) + (g2-g1)*(g2-g1) + (b2-b1)*(b2-b1)))
}

// Linear interpolation between a and b by t
func lerp(a, b, t float64) float64 {
	return (1-t)*a + t*b
}

// Smooth flood fill algorithm
// dear GC, please don't kill my allocations, all I want is a little memory
// ok maybe not so little, but I'm not leaking it, I promise
// ok maybe I am, a little
// ok I'm totally "leaking" it, but you'll get it back eventually, which is your whole thing, right?
// just don't jump in and start collecting while I'm in the middle of this
// and don't be inconsistent as F*
func smoothFloodFill(img draw.Image, x, y int, fill color.Color, tolerance float64) {
	q := []image.Point{{x, y}}
	origColor := img.At(x, y)

	// make a visited map to avoid cycles
	visited := make(map[image.Point]bool)

	for len(q) > 0 {
		current := q[0]
		q = q[1:]

		currColor := img.At(current.X, current.Y)
		distance := colorDistance(currColor, origColor)

		if distance > tolerance {
			continue
		}

		// Interpolate between existing color and fill color based on distance
		t := distance / tolerance
		r1, g1, b1, a1 := fill.RGBA()
		r2, g2, b2, a2 := currColor.RGBA()

		r := uint8(lerp(float64(r1), float64(r2), t) / 256)
		g := uint8(lerp(float64(g1), float64(g2), t) / 256)
		b := uint8(lerp(float64(b1), float64(b2), t) / 256)
		a := uint8(lerp(float64(a1), float64(a2), t) / 256)

		img.Set(current.X, current.Y, color.RGBA{r, g, b, a})

		for _, d := range []image.Point{{-1, 0}, {1, 0}, {0, -1}, {0, 1}} {
			next := image.Point{X: current.X + d.X, Y: current.Y + d.Y}

			// Check if we've already visited this pixel
			if visited[next] {
				continue
			}

			visited[next] = true

			// Check bounds to make sure (next.X, next.Y) is inside img.
			if next.X < img.Bounds().Min.X || next.X >= img.Bounds().Max.X || next.Y < img.Bounds().Min.Y || next.Y >= img.Bounds().Max.Y {
				continue
			}

			if colorDistance(img.At(next.X, next.Y), origColor) <= tolerance {
				q = append(q, next)
			}
		}
	}
}

func randomPastelColor() color.RGBA {
	r := uint8(rand.Intn(128) + 128) // 128 to 255
	g := uint8(rand.Intn(128) + 128) // 128 to 255
	b := uint8(rand.Intn(128) + 128) // 128 to 255
	return color.RGBA{R: r, G: g, B: b, A: 255}
}

func randomDarkColor() color.RGBA {
	r := uint8(rand.Intn(128)) // 0 to 127
	g := uint8(rand.Intn(128)) // 0 to 127
	b := uint8(rand.Intn(128)) // 0 to 127
	return color.RGBA{R: r, G: g, B: b, A: 255}
}

// adds a light cone to the image. warning: has trigonometry inside
func addLightCone(img draw.Image, x1, y1, x2, y2 int, maxAlpha uint8) {
	dx := float64(x2 - x1)
	dy := float64(y2 - y1)
	distance := math.Sqrt(dx*dx + dy*dy)
	angle := math.Atan2(dy, dx)

	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			// Calculate angle and distance from (x1, y1) to (x, y)
			localDx := float64(x - x1)
			localDy := float64(y - y1)
			localDistance := math.Sqrt(localDx*localDx + localDy*localDy)
			localAngle := math.Atan2(localDy, localDx) - angle

			// Normalize the angle to between -π and π
			// that sounded way more complicated than it is, see...
			for localAngle < -math.Pi {
				localAngle += 2 * math.Pi
			}
			for localAngle > math.Pi {
				localAngle -= 2 * math.Pi
			}

			// Calculate alpha based on angle and distance
			// (close to the light source is brighter)
			// (which looks absolutely horrible in GIFs but oh well)
			alpha := uint8(0)
			if math.Abs(localAngle) < math.Pi/8 && localDistance <= distance {
				alpha = uint8((1 - (localDistance / distance)) * float64(maxAlpha))
			}

			// let there be light
			if alpha > 0 {
				r, g, b, a := img.At(x, y).RGBA()
				r, g, b = r>>8, g>>8, b>>8 // Normalize to 8-bit color
				newR := uint8(float64(r) + (float64(255-r) * float64(alpha) / 255.0))
				newG := uint8(float64(g) + (float64(255-g) * float64(alpha) / 255.0))
				newB := uint8(float64(b) + (float64(255-b) * float64(alpha) / 255.0))
				img.Set(x, y, color.RGBA{newR, newG, newB, uint8(a >> 8)})
			}
		}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	const totalFrames = 30
	const frameDelay = 8

	frames := [totalFrames]*image.Paletted{}
	delays := [totalFrames]int{}

	for i := 0; i < totalFrames; i++ {
		delays[i] = frameDelay
	}

	file, _ := os.Open("input.png")
	origImg, _, _ := image.Decode(file)

	imgWidth := origImg.Bounds().Max.X
	imgHeight := origImg.Bounds().Max.Y
	wg := sync.WaitGroup{}

	for i := 0; i < totalFrames; i++ {
		img := image.NewRGBA(origImg.Bounds())
		draw.Draw(img, img.Bounds(), origImg, image.Point{}, draw.Src)
		wg.Add(1)

		go func(i int) {
			defer func() {
				wg.Done()
				fmt.Printf("Finished frame %d\n", i)
			}()

			fmt.Printf("Starting frame %d\n", i)

			for j := 0; j < 60; j++ {
				fmt.Printf("Frame %d, circle %d\n", i, j)

				var c color.RGBA

				// alternate between pastel and dark colors
				// makes it look like a "proper" disco
				// as in all things, I'm a disco expert too :)
				if i%2 == 0 {
					c = randomPastelColor()
				} else {
					c = randomDarkColor()
				}

				x := rand.Intn(imgWidth)
				y := rand.Intn(imgHeight)
				smoothFloodFill(img, x, y, c, 800)
			}

			// inverse the image colors for the first half of the frames
			// I tried making a strobe effect by alternating every frame, but perhaps it could cause seizures :/
			if i < (totalFrames / 2) {
				for x := 0; x < imgWidth; x++ {
					for y := 0; y < imgHeight; y++ {
						r, g, b, a := img.At(x, y).RGBA()
						// invert the colors
						img.Set(x, y, color.RGBA{R: uint8(255 - r), G: uint8(255 - g), B: uint8(255 - b), A: uint8(a)})
					}
				}
			}

			// today is the day that you learn about the Plan9 palette, child
			paletted := image.NewPaletted(img.Bounds(), palette.Plan9)
			draw.Draw(paletted, paletted.Bounds(), img, image.Point{}, draw.Over)

			frames[i] = paletted
		}(i)
	}

	wg.Wait()

	for i := 0; i < totalFrames; i += 10 {
		coneTarget := imgHeight

		if imgHeight > 90 {
			coneTarget = imgHeight - 10
		}

		// add a light cone
		startX := rand.Intn(imgWidth / 2)

		if (i/10)%2 == 0 {
			startX += imgWidth / 2
		}

		endX := rand.Intn(imgWidth)

		for y := 0; y < 10 && y < totalFrames; y++ {
			addLightCone(frames[i+y], startX, 0, endX, coneTarget, 255)
		}
	}

	f, err := os.Create(time.Now().Format("20060102150405") + "_disco.gif")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	err = gif.EncodeAll(f, &gif.GIF{
		Image:     frames[:],
		Delay:     delays[:],
		LoopCount: 0,
	})

	if err != nil {
		panic(err)
	}
}
