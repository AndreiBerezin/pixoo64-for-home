package screens

import (
	"fmt"
	"image/color"
	"time"

	"github.com/AndreiBerezin/pixoo64/internal/screens/image/fonts"
)

func (s *Screens) DrawBottomOnAir(startedAt time.Time) error {
	return s.drawOnAir(45, startedAt)
}

func (s *Screens) drawOnAir(startY int, startedAt time.Time) error {
	s.image.DrawRect(2, startY+3, 60, 14, color.RGBA{255, 50, 50, 255})
	s.image.DrawString("ON AIR", 12, startY+15, color.RGBA{255, 255, 255, 255}, fonts.FontMicro5Big)

	onAirDuration := int(time.Since(startedAt).Minutes())
	s.image.DrawString(fmt.Sprintf("%d min", onAirDuration), 2, startY+1, color.RGBA{255, 255, 255, 255}, fonts.FontMicro5Normal)

	return nil
}
