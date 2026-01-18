package image

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

	"github.com/AndreiBerezin/pixoo64/internal/screens/image/cache"
	"github.com/AndreiBerezin/pixoo64/internal/screens/image/fonts"
	"github.com/AndreiBerezin/pixoo64/pkg/http_client"
	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
	"golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

type Image struct {
	image  *image.RGBA
	client *http_client.Client
	cache  *cache.Cache
}

func NewImage(width int, height int) *Image {
	return &Image{
		image:  image.NewRGBA(image.Rect(0, 0, width, height)),
		client: http_client.New(),
		cache: cache.NewCache("cache", 1*time.Hour, func(url string) string {
			parts := strings.Split(url, "/")
			return parts[len(parts)-1]
		}),
	}
}

func (i *Image) Reset() {
	i.DrawRect(0, 0, i.image.Rect.Size().X, i.image.Rect.Size().Y, color.RGBA{0, 0, 0, 255})
}

func (i *Image) Image() *image.RGBA {
	return i.image
}

func (i *Image) DrawString(text string, x int, y int, color color.Color, fontFace fonts.FontFace) {
	fontDrawer := &font.Drawer{
		Dst:  i.image,
		Src:  image.NewUniform(color),
		Face: fonts.Fonts[fontFace],
		Dot:  fixed.P(x, y),
	}
	fontDrawer.DrawString(text)
}

func (i *Image) DrawPNGFromFile(filename string, x int, y int, targetSize int) error {
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

	draw.Draw(i.image, i.image.Bounds(), targetImage, image.Point{X: x * -1, Y: y * -1}, draw.Over)

	return nil
}

func (i *Image) DrawSVGFromURL(url string, x int, y int, targetSize int) error {
	data, err := i.getCachedSvg(url)
	if err != nil {
		return fmt.Errorf("failed to get cached SVG: %w", err)
	}

	svg, err := oksvg.ReadIconStream(bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("failed to parse SVG from bytes: %w", err)
	}

	svg.SetTarget(0, 0, float64(targetSize), float64(targetSize))

	targetImage := image.NewRGBA(image.Rect(0, 0, targetSize, targetSize))

	scanner := rasterx.NewScannerGV(targetSize, targetSize, targetImage, targetImage.Bounds())
	raster := rasterx.NewDasher(targetSize, targetSize, scanner)
	svg.Draw(raster, 1.0)

	draw.Draw(i.image, i.image.Bounds(), targetImage, image.Point{X: x * -1, Y: y * -1}, draw.Over)

	return nil
}

func (i *Image) getCachedSvg(url string) ([]byte, error) {
	data, err := i.cache.Get(url)
	if err == nil {
		return data, nil
	}
	if !errors.Is(err, cache.ErrCacheExpired) {
		return nil, fmt.Errorf("failed to get SVG from cache: %w", err)
	}

	var response []byte
	err = i.client.Get(url, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to get SVG from URL: %w", err)
	}
	err = i.cache.Set(url, response)
	if err != nil {
		return nil, fmt.Errorf("failed to set SVG in cache: %w", err)
	}

	return response, nil
}

func (i *Image) DrawRect(x int, y int, width int, height int, color color.Color) {
	for k := x; k < x+width; k++ {
		for l := y; l < y+height; l++ {
			i.image.Set(k, l, color)
		}
	}
}
