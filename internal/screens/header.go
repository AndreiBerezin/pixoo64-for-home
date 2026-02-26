package screens

import (
	"fmt"
	"image/color"
	"time"

	"github.com/AndreiBerezin/pixoo64/internal/screens/image/fonts"
)

func (s *Screens) DrawHeader() error {
	now := time.Now()
	s.image.DrawString(fmt.Sprintf("%02d:%02d", now.Hour(), now.Minute()), 2, 7, color.RGBA{100, 255, 255, 255}, fonts.FontMicro5Normal)
	s.image.DrawString(fmt.Sprintf("%02d.%02d", now.Day(), now.Month()), 45, 7, color.RGBA{100, 255, 255, 255}, fonts.FontMicro5Normal)

	s.image.DrawRect(2, 9, 60, 1, color.RGBA{50, 50, 50, 255})

	return nil
}
