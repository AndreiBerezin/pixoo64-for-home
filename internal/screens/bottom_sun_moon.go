package screens

import (
	"fmt"
	"image/color"

	"github.com/AndreiBerezin/pixoo64/internal/collector/types"
	"github.com/AndreiBerezin/pixoo64/internal/screens/image/fonts"
)

func (s *Screens) DrawBottomSunMoon(yandexData *types.YandexData) error {
	if yandexData == nil {
		return nil
	}

	startY := 35

	s.image.DrawString(yandexData.Sun.SunriseTime, 2, startY+9, color.RGBA{255, 255, 255, 255}, fonts.FontMicro5Normal)
	s.image.DrawPNGFromFile("static/images/sunrise.png", 21, startY+1, 10)
	s.image.DrawString(yandexData.Sun.SunsetTime, 34, startY+9, color.RGBA{255, 255, 255, 255}, fonts.FontMicro5Normal)

	s.image.DrawPNGFromFile(yandexData.Moon.Icon, 2, startY+12, 10)
	s.image.DrawString(fmt.Sprintf("%d", yandexData.Moon.MoonPhaseDay), 15, startY+20, color.RGBA{255, 255, 255, 255}, fonts.FontMicro5Normal)

	return nil
}
