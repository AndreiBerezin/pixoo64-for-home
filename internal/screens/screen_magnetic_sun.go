package screens

import (
	"image/color"

	"github.com/AndreiBerezin/pixoo64/internal/collector/types"
	"github.com/AndreiBerezin/pixoo64/internal/screens/image/fonts"
)

func (s *Screens) DrawMagneticSun(magneticData *types.MagneticData, yandexData *types.YandexData) error {
	if magneticData == nil || yandexData == nil {
		return nil
	}

	startY := 45

	offset := 2
	for _, day := range magneticData.Days {
		s.image.DrawString(day.Day, offset, startY, color.RGBA{255, 255, 255, 255}, fonts.FontMicro5Normal)
		offset += 9
		for _, hour := range day.Hours {
			//hour.Level = rand.Float32() * 10
			col := color.RGBA{100, 255, 100, 255}
			if hour.Level >= 4 && hour.Level < 5 {
				col = color.RGBA{255, 255, 100, 255}
			} else if hour.Level >= 5 {
				col = color.RGBA{255, 100, 100, 255}
			}
			level := max(1, int(hour.Level))

			s.image.DrawRect(offset, startY-level, 1, level, col)
			offset += 1
		}
		offset += 4
	}

	s.image.DrawString(yandexData.Sun.SunriseTime, 2, startY+9, color.RGBA{255, 255, 255, 255}, fonts.FontMicro5Normal)
	s.image.DrawPNGFromFile("static/images/sunrise.png", 21, startY+1, 10)
	s.image.DrawString(yandexData.Sun.SunsetTime, 34, startY+9, color.RGBA{255, 255, 255, 255}, fonts.FontMicro5Normal)

	return nil
}
