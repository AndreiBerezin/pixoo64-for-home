package frame

import (
	"image"
)

type Frame struct {
	rgba  *image.RGBA
	speed int
}

func New(img *image.RGBA, speed int) Frame {
	return Frame{
		rgba:  img,
		speed: speed,
	}
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
