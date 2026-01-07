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

func (s *MagneticScreen) DrawStatic(magneticData *types.MagneticData, yandexData *types.YandexData) error {
	if magneticData == nil || yandexData == nil {
		return nil
	}

	startY := 45

	offset := 2
	for _, day := range magneticData.Days {
		s.drawer.DrawString(day.Day, offset, startY, color.RGBA{255, 255, 255, 255}, drawer.FontMicro5Normal)
		offset += 9
		for _, hour := range day.Hours {
			//hour.Level = rand.Float32() * 10
			level := max(1, int(hour.Level))
			col := color.RGBA{100, 255, 100, 255}
			if hour.Level >= 3 && hour.Level < 5 {
				col = color.RGBA{255, 255, 100, 255}
			} else if hour.Level >= 5 {
				col = color.RGBA{255, 100, 100, 255}
			}
			s.drawer.DrawRect(offset, startY-level, 1, level, col)
			offset += 1
		}
		offset += 4
	}

	s.drawer.DrawString(yandexData.Sun.SunriseTime, 2, startY+9, color.RGBA{255, 255, 255, 255}, drawer.FontMicro5Normal)
	s.drawer.DrawPNGFromFile("static/images/sunrise.png", 20, startY+1, 10)
	s.drawer.DrawString(yandexData.Sun.SunsetTime, 33, startY+9, color.RGBA{255, 255, 255, 255}, drawer.FontMicro5Normal)

	return nil
}
