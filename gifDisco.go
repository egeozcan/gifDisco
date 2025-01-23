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

	"github.com/egeozcan/gifDisco/colorutils"
	"github.com/egeozcan/gifDisco/imaging"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

const (
	totalFrames = 30
	frameDelay  = 8
)

func main() {
	rand.Seed(time.Now().UnixNano())

	// Load input image
	img := loadInputImage()
	bounds := img.Bounds()

	// Initialize GIF components
	frames := make([]*image.Paletted, totalFrames)
	delays := make([]int, totalFrames)
	for i := range delays {
		delays[i] = frameDelay
	}

	// Process frames concurrently
	var wg sync.WaitGroup
	wg.Add(totalFrames)

	for i := 0; i < totalFrames; i++ {
		go func(frameNum int) {
			defer wg.Done()
			processFrame(frameNum, img, bounds, frames)
		}(i)
	}

	wg.Wait()
	addLightCones(frames, bounds)
	saveGIF(frames, delays)
}

func loadInputImage() image.Image {
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		img, _, err := image.Decode(os.Stdin)
		if err != nil {
			panic("Error decoding piped input: " + err.Error())
		}
		return img
	}

	file, err := os.Open("input.png")
	if err != nil {
		panic("No input provided: " + err.Error())
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		panic("Error decoding input.png: " + err.Error())
	}
	return img
}

func processFrame(frameNum int, orig image.Image, bounds image.Rectangle, frames []*image.Paletted) {
	// Create base image
	img := image.NewRGBA(bounds)
	draw.Draw(img, bounds, orig, bounds.Min, draw.Src)

	// Apply 60 flood fills
	for j := 0; j < 60; j++ {
		var c color.RGBA
		if frameNum%2 == 0 {
			c = colorutils.RandomPastelColor()
		} else {
			c = colorutils.RandomDarkColor()
		}
		x := rand.Intn(bounds.Dx())
		y := rand.Intn(bounds.Dy())
		filled := imaging.SmoothFloodFill(img, x, y, c, 800)
		draw.Draw(img, bounds, filled, image.Point{}, draw.Over)
	}

	// Invert colors for first half of frames
	if frameNum < totalFrames/2 {
		invertImage(img)
	}

	// Convert to paletted image
	paletted := image.NewPaletted(bounds, palette.Plan9)
	draw.Draw(paletted, bounds, img, bounds.Min, draw.Src)
	frames[frameNum] = paletted
	fmt.Printf("Processed frame %d\n", frameNum)
}

func invertImage(img *image.RGBA) {
	pix := img.Pix
	for i := 0; i < len(pix); i += 4 {
		pix[i] = 255 - pix[i]
		pix[i+1] = 255 - pix[i+1]
		pix[i+2] = 255 - pix[i+2]
	}
}

func addLightCones(frames []*image.Paletted, bounds image.Rectangle) {
	height := bounds.Dy()
	for i := 0; i < len(frames); i += 10 {
		targetY := height
		if height > 90 {
			targetY = height - 10
		}

		startX := rand.Intn(bounds.Dx() / 2)
		if (i/10)%2 == 0 {
			startX += bounds.Dx() / 2
		}
		endX := rand.Intn(bounds.Dx())

		for y := 0; y < 10 && i+y < len(frames); y++ {
			if frames[i+y] != nil {
				applyLightCone(frames[i+y], startX, 0, endX, targetY)
			}
		}
	}
}

func applyLightCone(img *image.Paletted, x1, y1, x2, y2 int) {
	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, img.Bounds().Min, draw.Src)
	addLightCone(rgba, x1, y1, x2, y2, 255)
	draw.Draw(img, img.Bounds(), rgba, rgba.Bounds().Min, draw.Src)
}

func addLightCone(img *image.RGBA, x1, y1, x2, y2 int, maxAlpha uint8) {
	dx := float64(x2 - x1)
	dy := float64(y2 - y1)
	distanceSquared := dx*dx + dy*dy
	if distanceSquared == 0 {
		return
	}

	distance := math.Sqrt(distanceSquared)
	dirX, dirY := dx/distance, dy/distance
	cosTheta := math.Cos(math.Pi / 8)

	bounds := img.Bounds()
	stride := img.Stride
	pix := img.Pix

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			ldx := float64(x - x1)
			ldy := float64(y - y1)
			ldSquared := ldx*ldx + ldy*ldy

			if ldSquared > distanceSquared {
				continue
			}

			ld := math.Sqrt(ldSquared)
			if ld == 0 {
				offset := (y-bounds.Min.Y)*stride + (x-bounds.Min.X)*4
				pix[offset] = 255 - pix[offset]
				pix[offset+1] = 255 - pix[offset+1]
				pix[offset+2] = 255 - pix[offset+2]
				continue
			}

			dot := (ldx*dirX + ldy*dirY) / ld
			if dot >= cosTheta {
				alpha := uint8((1 - (ld / distance)) * float64(maxAlpha))
				if alpha > 0 {
					offset := (y-bounds.Min.Y)*stride + (x-bounds.Min.X)*4
					r := pix[offset]
					g := pix[offset+1]
					b := pix[offset+2]

					pix[offset] = uint8(float64(r) + (float64(255-r)*float64(alpha))/255)
					pix[offset+1] = uint8(float64(g) + (float64(255-g)*float64(alpha))/255)
					pix[offset+2] = uint8(float64(b) + (float64(255-b)*float64(alpha))/255)
				}
			}
		}
	}
}

func saveGIF(frames []*image.Paletted, delays []int) {
	f, err := os.Create(time.Now().Format("20060102150405") + "_disco.gif")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if err := gif.EncodeAll(f, &gif.GIF{
		Image:     frames,
		Delay:     delays,
		LoopCount: 0,
	}); err != nil {
		panic(err)
	}
}
