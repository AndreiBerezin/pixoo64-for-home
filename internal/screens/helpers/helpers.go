package helpers

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"

	"github.com/AndreiBerezin/pixoo64/internal/screens/helpers/fonts"
	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
	"golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

func DrawString(img *image.RGBA, text string, x int, y int, color color.Color, fontFace fonts.FontFace) {
	fontDrawer := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(color),
		Face: fonts.Fonts[fontFace],
		Dot:  fixed.P(x, y),
	}
	fontDrawer.DrawString(text)
}

func DrawPNGFromFile(img *image.RGBA, filename string, x int, y int, targetSize int) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open image: %w", err)
	}
	defer file.Close()

	rawImg, err := png.Decode(file)
	if err != nil {
		return fmt.Errorf("failed to decode image: %w", err)
	}

	targetImage := image.NewRGBA(image.Rect(0, 0, targetSize, targetSize))
	draw.NearestNeighbor.Scale(targetImage, targetImage.Bounds(), rawImg, rawImg.Bounds(), draw.Over, nil)

	draw.Draw(img, img.Bounds(), targetImage, image.Point{X: x * -1, Y: y * -1}, draw.Over)

	return nil
}

func DrawSVGFromBytes(img *image.RGBA, data []byte, x int, y int, targetSize int) error {
	svg, err := oksvg.ReadIconStream(bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("failed to parse SVG from bytes: %w", err)
	}

	svg.SetTarget(0, 0, float64(targetSize), float64(targetSize))

	targetImage := image.NewRGBA(image.Rect(0, 0, targetSize, targetSize))

	scanner := rasterx.NewScannerGV(targetSize, targetSize, targetImage, targetImage.Bounds())
	raster := rasterx.NewDasher(targetSize, targetSize, scanner)
	svg.Draw(raster, 1.0)

	draw.Draw(img, img.Bounds(), targetImage, image.Point{X: x * -1, Y: y * -1}, draw.Over)

	return nil
}

func DrawRect(img *image.RGBA, x int, y int, width int, height int, color color.Color) {
	for i := x; i < x+width; i++ {
		for j := y; j < y+height; j++ {
			img.Set(i, j, color)
		}
	}
}
