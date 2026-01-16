package drawer

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"strings"
	"time"

	"github.com/AndreiBerezin/pixoo64/pkg/http_client"
	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
	"golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

type Drawer struct {
	client *http_client.Client
	img    *image.RGBA
	cache  *Cache
}

func New() *Drawer {
	drawer := &Drawer{
		client: http_client.New(),
		img:    image.NewRGBA(image.Rect(0, 0, 64, 64)),
		cache: NewCache("cache", 1*time.Hour, func(url string) string {
			parts := strings.Split(url, "/")
			return parts[len(parts)-1]
		}),
	}
	drawer.Reset()

	return drawer
}

func (d *Drawer) Reset() {
	for y := 0; y < d.img.Rect.Size().Y; y++ {
		for x := 0; x < d.img.Rect.Size().X; x++ {
			d.img.Set(x, y, color.RGBA{0, 0, 0, 255})
		}
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
	data, err := d.getCachedSvg(url)
	if err != nil {
		return fmt.Errorf("failed to get cached SVG: %w", err)
	}

	svg, err := oksvg.ReadIconStream(bytes.NewReader(data))
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

func (d *Drawer) getCachedSvg(url string) ([]byte, error) {
	data, err := d.cache.Get(url)
	if err == nil {
		return data, nil
	}
	if !errors.Is(err, ErrCacheExpired) {
		return nil, fmt.Errorf("failed to get SVG from cache: %w", err)
	}

	var response []byte
	err = d.client.Get(url, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to get SVG from URL: %w", err)
	}
	err = d.cache.Set(url, response)
	if err != nil {
		return nil, fmt.Errorf("failed to set SVG in cache: %w", err)
	}

	return response, nil
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
