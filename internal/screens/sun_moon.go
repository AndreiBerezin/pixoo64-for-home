package screens

import (
	"fmt"
	"image/color"

	"github.com/AndreiBerezin/pixoo64/internal/collector/types"
	"github.com/AndreiBerezin/pixoo64/internal/screens/image/fonts"
)

func (s *Screens) DrawBottomSunMoon(data *types.CollectedData) error {
	return s.drawSunMoon(35, data.YandexData)
}

func (s *Screens) drawSunMoon(startY int, data *types.YandexData) error {
	if data == nil {
		return nil
	}

	s.image.DrawString(data.Sun.SunriseTime, 2, startY+9, color.RGBA{255, 255, 255, 255}, fonts.FontMicro5Normal)
	s.image.DrawPNGFromFile("static/images/sunrise.png", 21, startY+1, 10)
	s.image.DrawString(data.Sun.SunsetTime, 34, startY+9, color.RGBA{255, 255, 255, 255}, fonts.FontMicro5Normal)

	s.image.DrawPNGFromFile(data.Moon.Icon, 2, startY+12, 10)
	s.image.DrawString(fmt.Sprintf("%d", data.Moon.MoonPhaseDay), 15, startY+20, color.RGBA{255, 255, 255, 255}, fonts.FontMicro5Normal)

	return nil
}
