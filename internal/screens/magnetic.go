package screens

import (
	"image/color"

	"github.com/AndreiBerezin/pixoo64/internal/collector/types"
	"github.com/AndreiBerezin/pixoo64/internal/drawer"
)

type MagneticScreen struct {
	drawer *drawer.Drawer
}

func NewMagneticScreen(drawer *drawer.Drawer) (*MagneticScreen, error) {
	return &MagneticScreen{drawer: drawer}, nil
}

func (s *MagneticScreen) DrawStatic(data *types.MagneticData) error {
	if data == nil {
		return nil
	}

	startY := 35

	for i := 1; i < 10; i++ {
		s.drawSquare(i*4, startY, color.RGBA{255, 100, 100, 255})
	}
	for i := 1; i < 10; i++ {
		s.drawSquare(i*4, startY, color.RGBA{100, 255, 100, 255})
	}

	return nil
}

func (s *MagneticScreen) drawSquare(x int, y int, color color.Color) {
	for i := x - 1; i < x+1; i++ {
		for j := y - 1; j < y+1; j++ {
			s.drawer.Image().Set(i, j, color)
		}
	}
}
