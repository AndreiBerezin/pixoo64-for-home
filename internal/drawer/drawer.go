package drawer

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"net/http"
	"os"

	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
	"golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

type Drawer struct {
	img *image.RGBA
}

func NewDrawer() *Drawer {
	img := image.NewRGBA(image.Rect(0, 0, 64, 64))
	for y := 0; y < 64; y++ {
		for x := 0; x < 64; x++ {
			img.Set(x, y, color.RGBA{0, 0, 0, 255})
		}
	}

	return &Drawer{
		img: img,
	}
}

func (d *Drawer) DrawString(text string, x int, y int, color color.Color, fontFace FontFace) {
	fontDrawer := &font.Drawer{
		Dst:  d.img,
		Src:  image.NewUniform(color),
		Face: Fonts[fontFace],
		Dot:  fixed.P(x, y),
	}
	fontDrawer.DrawString(text)
}

func (d *Drawer) DrawPNGFromFile(filename string, x int, y int, targetSize int) error {
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

	draw.Draw(d.Image(), d.Image().Bounds(), targetImage, image.Point{X: x * -1, Y: y * -1}, draw.Over)

	return nil
}

func (d *Drawer) DrawSVGFromURL(url string, x int, y int, targetSize int) error {
	response, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to get SVG from URL: %w", err)
	}
	defer response.Body.Close()

	svg, err := oksvg.ReadIconStream(response.Body)
	if err != nil {
		return fmt.Errorf("failed to parse SVG from URL: %w", err)
	}

	svg.SetTarget(0, 0, float64(targetSize), float64(targetSize))

	targetImage := image.NewRGBA(image.Rect(0, 0, targetSize, targetSize))

	scanner := rasterx.NewScannerGV(targetSize, targetSize, targetImage, targetImage.Bounds())
	raster := rasterx.NewDasher(targetSize, targetSize, scanner)
	svg.Draw(raster, 1.0)

	draw.Draw(d.Image(), d.Image().Bounds(), targetImage, image.Point{X: x * -1, Y: y * -1}, draw.Over)

	return nil
}

func (d *Drawer) DrawRect(x int, y int, width int, height int, color color.Color) {
	for i := x; i < x+width; i++ {
		for j := y; j < y+height; j++ {
			d.Image().Set(i, j, color)
		}
	}
}

func (d *Drawer) Image() *image.RGBA {
	return d.img
}
