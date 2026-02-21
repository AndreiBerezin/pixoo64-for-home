package screens

import (
	"image/color"

	"github.com/AndreiBerezin/pixoo64/internal/collector/types"
	"github.com/AndreiBerezin/pixoo64/internal/screens/image/fonts"
)

func (s *Screens) DrawBottomMagneticPressure(magneticData *types.MagneticData, pressureData *types.PressureData) error {
	if magneticData == nil || pressureData == nil {
		return nil
	}

	startY := 37

	s.image.DrawPNGFromFile("static/images/magnet.png", 2, startY, 6)

	offsetX := 13
	for _, day := range magneticData.Days {
		s.image.DrawRect(offsetX-1, startY-4, 10, 10, color.RGBA{35, 35, 35, 255})

		s.image.DrawString(day.Day, offsetX, startY+13, color.RGBA{255, 255, 255, 255}, fonts.FontMicro5Normal)

		for _, hour := range day.Hours {
			//hour.Level = rand.Float32() * 10
			col := color.RGBA{100, 255, 100, 255}
			if hour.Level >= 4 && hour.Level < 5 {
				col = color.RGBA{255, 255, 100, 255}
			} else if hour.Level >= 5 {
				col = color.RGBA{255, 100, 100, 255}
			}
			level := max(1, int(hour.Level))

			s.image.DrawRect(offsetX, startY+6-level, 1, level, col)
			offsetX += 1
		}

		offsetX += 4
	}

	s.image.DrawPNGFromFile("static/images/pressure.png", 2, startY+18, 7)

	offsetX = 13
	for _, day := range pressureData.Days {
		s.image.DrawRect(offsetX-1, startY+15, 10, 10, color.RGBA{35, 35, 35, 255})

		for _, hour := range day.Hours {
			var level int
			if hour.Pressure <= 750 {
				level = 1
			} else if hour.Pressure >= 770 {
				level = 10
			} else {
				level = int((hour.Pressure-750)/20*9) + 1
			}

			col := color.RGBA{100, 255, 100, 255}
			switch level {
			case 1, 10:
				col = color.RGBA{255, 100, 100, 255}
			case 9, 8, 3, 2:
				col = color.RGBA{255, 255, 100, 255}
			}

			s.image.DrawRect(offsetX, startY+25-level, 1, level, col)
			offsetX += 1
		}

		offsetX += 4
	}

	return nil
}
