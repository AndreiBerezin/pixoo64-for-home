package screens

import (
	"fmt"
	"image/color"

	"github.com/AndreiBerezin/pixoo64/internal/collector/types"
	"github.com/AndreiBerezin/pixoo64/internal/screens/helpers"
	"github.com/AndreiBerezin/pixoo64/internal/screens/helpers/fonts"
)

func (s *Screens) DrawToday(data *types.YandexData) error {
	if data == nil {
		return nil
	}

	startY := 45

	for i, item := range data.DayWeather.Items {
		temperature := fmt.Sprintf("%ḋ", item.Temperature)
		if item.Temperature > 0 {
			temperature = "+" + temperature
		} else if item.Temperature == 0 {
			temperature = " " + temperature
		}
		helpers.DrawString(s.img, temperature, 2+i*16, startY, color.RGBA{255, 255, 255, 255}, fonts.FontMicro5Normal)

		err := s.drawSVGFromURL(item.Icon, 2+i*16, startY, 12)
		if err != nil {
			return fmt.Errorf("failed to draw icon: %w", err)
		}
		helpers.DrawString(s.img, item.Name, 6+i*16, startY+16, color.RGBA{255, 255, 255, 255}, fonts.FontTiny5Normal)
	}

	return nil
}

func (s *Screens) DrawTomorrow(data *types.YandexData) error {
	if data == nil {
		return nil
	}

	startY := 45

	nearestItems := []item{
		{name: "24", icon: data.CurrentWeather.Icon, temperature: -23},
		{name: "25", icon: data.CurrentWeather.Icon, temperature: +11},
		{name: "26", icon: data.CurrentWeather.Icon, temperature: 0},
	}
	for i, item := range nearestItems {
		temperature := fmt.Sprintf("%d", item.temperature)
		if item.temperature > 0 {
			temperature = "+" + temperature
		} else if item.temperature == 0 {
			temperature = " " + temperature
		}
		helpers.DrawString(s.img, temperature, 2+i*16, startY, color.RGBA{255, 255, 255, 255}, fonts.FontMicro5Normal)

		err := s.drawSVGFromURL(item.icon, 2+i*16, startY, 12)
		if err != nil {
			return fmt.Errorf("failed to draw icon: %w", err)
		}
		helpers.DrawString(s.img, item.name, 5+i*16, startY+17, color.RGBA{255, 255, 255, 255}, fonts.FontTiny5Normal)
	}

	return nil
}

type item struct {
	name        string
	icon        string
	temperature int
}
