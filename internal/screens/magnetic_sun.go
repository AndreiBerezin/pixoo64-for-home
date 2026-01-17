package screens

import (
	"image/color"

	"github.com/AndreiBerezin/pixoo64/internal/collector/types"
	"github.com/AndreiBerezin/pixoo64/internal/screens/helpers"
	"github.com/AndreiBerezin/pixoo64/internal/screens/helpers/fonts"
)

func (s *Screens) DrawMagneticSun(magneticData *types.MagneticData, yandexData *types.YandexData) error {
	if magneticData == nil || yandexData == nil {
		return nil
	}

	startY := 45

	offset := 2
	for _, day := range magneticData.Days {
		helpers.DrawString(s.img, day.Day, offset, startY, color.RGBA{255, 255, 255, 255}, fonts.FontMicro5Normal)
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

			helpers.DrawRect(s.img, offset, startY-level, 1, level, col)
			offset += 1
		}
		offset += 4
	}

	helpers.DrawString(s.img, yandexData.Sun.SunriseTime, 2, startY+9, color.RGBA{255, 255, 255, 255}, fonts.FontMicro5Normal)
	helpers.DrawPNGFromFile(s.img, "static/images/sunrise.png", 21, startY+1, 10)
	helpers.DrawString(s.img, yandexData.Sun.SunsetTime, 34, startY+9, color.RGBA{255, 255, 255, 255}, fonts.FontMicro5Normal)

	return nil
}
