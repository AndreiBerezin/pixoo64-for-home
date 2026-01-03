package frame

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"path"

	"golang.org/x/image/draw"
)

type Frame struct {
	rgba  *image.RGBA
	speed int
}

func NewFrame(filename string, speed int) (*Frame, error) {
	path := path.Join("images", filename)
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open image: %w", err)
	}
	defer file.Close()

	rawImg, err := png.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("failed to decode png: %w", err)
	}

	// todo: унести в константы
	img := image.NewRGBA(image.Rect(0, 0, 64, 64))
	draw.NearestNeighbor.Scale(img, img.Bounds(), rawImg, rawImg.Bounds(), draw.Over, nil)

	return &Frame{
		rgba:  img,
		speed: speed,
	}, nil
}

func NewFrameImage(img *image.RGBA, speed int) (*Frame, error) {
	return &Frame{
		rgba:  img,
		speed: speed,
	}, nil
}

func (f *Frame) ToBytes() []byte {
	count := 0
	raw := make([]byte, f.rgba.Rect.Size().X*f.rgba.Rect.Size().Y*3)
	for y := 0; y < f.rgba.Rect.Size().Y; y++ {
		for x := 0; x < f.rgba.Rect.Size().X; x++ {
			raw[count*3] = f.rgba.RGBAAt(x, y).R
			raw[count*3+1] = f.rgba.RGBAAt(x, y).G
			raw[count*3+2] = f.rgba.RGBAAt(x, y).B
			count++
		}
	}
	return raw
}

func (f *Frame) Speed() int {
	return f.speed
}
