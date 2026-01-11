package screens

import (
	"fmt"
	"image/color"

	"github.com/AndreiBerezin/pixoo64/internal/collector/types"
	"github.com/AndreiBerezin/pixoo64/internal/drawer"
)

type ExtraWeatherScreen struct {
	drawer *drawer.Drawer
}

func NewExtraWeatherScreen(drawer *drawer.Drawer) *ExtraWeatherScreen {
	return &ExtraWeatherScreen{drawer: drawer}
}

func (s *ExtraWeatherScreen) DrawTodayStatic(data *types.YandexData) error {
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
		s.drawer.DrawString(temperature, 2+i*16, startY, color.RGBA{255, 255, 255, 255}, drawer.FontMicro5Normal)

		err := s.drawer.DrawSVGFromURL(item.Icon, 2+i*16, startY, 12)
		if err != nil {
			return fmt.Errorf("failed to draw icon: %w", err)
		}
		s.drawer.DrawString(item.Name, 6+i*16, startY+16, color.RGBA{255, 255, 255, 255}, drawer.FontTiny5Normal)
	}

	return nil
}

func (s *ExtraWeatherScreen) DrawTomorrowStatic(data *types.YandexData) error {
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
		s.drawer.DrawString(temperature, 2+i*16, startY, color.RGBA{255, 255, 255, 255}, drawer.FontMicro5Normal)

		err := s.drawer.DrawSVGFromURL(item.icon, 2+i*16, startY, 12)
		if err != nil {
			return fmt.Errorf("failed to draw icon: %w", err)
		}
		s.drawer.DrawString(item.name, 5+i*16, startY+17, color.RGBA{255, 255, 255, 255}, drawer.FontTiny5Normal)
	}

	return nil
}

type item struct {
	name        string
	icon        string
	temperature int
}
