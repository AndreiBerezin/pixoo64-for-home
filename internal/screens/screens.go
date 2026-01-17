package screens

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"strings"
	"time"

	"github.com/AndreiBerezin/pixoo64/internal/screens/helpers"
	"github.com/AndreiBerezin/pixoo64/internal/screens/helpers/cache"
	"github.com/AndreiBerezin/pixoo64/pkg/http_client"
)

type Screens struct {
	client *http_client.Client
	img    *image.RGBA
	cache  *cache.Cache
}

func New(width int, height int) *Screens {
	return &Screens{
		client: http_client.New(),
		img:    image.NewRGBA(image.Rect(0, 0, width, height)),
		cache: cache.NewCache("cache", 1*time.Hour, func(url string) string {
			parts := strings.Split(url, "/")
			return parts[len(parts)-1]
		}),
	}
}

func (s *Screens) Reset() {
	helpers.DrawRect(s.img, 0, 0, s.img.Rect.Size().X, s.img.Rect.Size().Y, color.RGBA{0, 0, 0, 255})
}

func (s *Screens) Image() *image.RGBA {
	return s.img
}

func (s *Screens) drawSVGFromURL(url string, x int, y int, targetSize int) error {
	data, err := s.getCachedSvg(url)
	if err != nil {
		return fmt.Errorf("failed to get cached SVG: %w", err)
	}

	helpers.DrawSVGFromBytes(s.img, data, x, y, targetSize)
	return nil
}

func (s *Screens) getCachedSvg(url string) ([]byte, error) {
	data, err := s.cache.Get(url)
	if err == nil {
		return data, nil
	}
	if !errors.Is(err, cache.ErrCacheExpired) {
		return nil, fmt.Errorf("failed to get SVG from cache: %w", err)
	}

	var response []byte
	err = s.client.Get(url, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to get SVG from URL: %w", err)
	}
	err = s.cache.Set(url, response)
	if err != nil {
		return nil, fmt.Errorf("failed to set SVG in cache: %w", err)
	}

	return response, nil
}
