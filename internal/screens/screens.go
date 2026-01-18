package screens

import (
	imagePkg "image"

	"github.com/AndreiBerezin/pixoo64/internal/screens/image"
)

type Screens struct {
	image *image.Image
}

func New(width int, height int) *Screens {
	return &Screens{
		image: image.NewImage(width, height),
	}
}

func (s *Screens) Reset() {
	s.image.Reset()
}

func (s *Screens) Image() *imagePkg.RGBA {
	return s.image.Image()
}
