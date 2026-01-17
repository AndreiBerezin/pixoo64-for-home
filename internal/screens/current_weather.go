package screens

import (
	"fmt"
	"image/color"
	"math"
	"strings"
	"time"

	"github.com/AndreiBerezin/pixoo64/internal/collector/types"
	"github.com/AndreiBerezin/pixoo64/internal/screens/helpers"
	"github.com/AndreiBerezin/pixoo64/internal/screens/helpers/fonts"
)

func (s *Screens) DrawCurrentWeather(data *types.YandexData) error {
	if data == nil {
		return nil
	}

	/*data.Weather.Temperature = -22
	data.Weather.FeelsLikeTemperature = -22
	data.Weather.WindSpeed = 22
	data.Weather.WindDirection = "sw"*/

	now := time.Now()
	helpers.DrawString(s.img, fmt.Sprintf("%02d:%02d", now.Hour(), now.Minute()), 2, 7, color.RGBA{100, 255, 255, 255}, fonts.FontMicro5Normal)
	helpers.DrawString(s.img, fmt.Sprintf("%02d.%02d", now.Day(), now.Month()), 46, 7, color.RGBA{100, 255, 255, 255}, fonts.FontMicro5Normal)

	helpers.DrawRect(s.img, 2, 9, 60, 1, color.RGBA{50, 50, 50, 255})

	err := s.drawSVGFromURL(data.CurrentWeather.Icon, 1, 11, 20)
	if err != nil {
		return fmt.Errorf("failed to draw icon: %w", err)
	}

	temperature := fmt.Sprintf("%ḋ", int(math.Abs(float64(data.CurrentWeather.Temperature)))) // тут спрятан ̇(символ градуса)
	helpers.DrawString(s.img, temperature, 29, 24, color.RGBA{255, 255, 255, 255}, fonts.FontMicro5Big)
	if data.CurrentWeather.Temperature > 0 {
		helpers.DrawString(s.img, "+", 24, 22, color.RGBA{255, 255, 255, 255}, fonts.FontMicro5Normal)
	} else if data.CurrentWeather.Temperature < 0 {
		helpers.DrawString(s.img, "-", 24, 22, color.RGBA{255, 255, 255, 255}, fonts.FontMicro5Normal)
	}

	sign := ""
	if data.CurrentWeather.FeelsLikeTemperature > 0 {
		sign = "+"
	}
	helpers.DrawString(s.img, fmt.Sprintf("%s%ḋ", sign, data.CurrentWeather.FeelsLikeTemperature), 25, 31, color.RGBA{100, 100, 255, 255}, fonts.FontMicro5Normal)

	helpers.DrawString(s.img, fmt.Sprintf("%d·%s", data.CurrentWeather.WindSpeed, s.windDirectionToRus(data.CurrentWeather.WindDirection)), 45, 31, color.RGBA{100, 255, 100, 255}, fonts.FontTiny5Normal)

	s.drawHouseWind(data.CurrentWeather.WindDirection)

	return nil
}

func (s *Screens) windDirectionToRus(direction string) string {
	rusMap := map[string]string{
		"n":  "с",
		"s":  "ю",
		"e":  "в",
		"w":  "з",
		"nw": "св",
		"ne": "сз",
		"sw": "юз",
		"se": "юз",
	}
	return rusMap[direction]
}

func (s *Screens) drawHouseWind(direction string) {
	startX := 54
	startY := 17
	arrowColor := color.RGBA{100, 255, 100, 255}

	if strings.HasPrefix(direction, "n") {
		helpers.DrawRect(s.img, startX, startY, 1, 5, arrowColor)

		s.img.Set(startX+1, startY+3, arrowColor)
		s.img.Set(startX+2, startY+2, arrowColor)
		s.img.Set(startX-2, startY+2, arrowColor)
		s.img.Set(startX-1, startY+3, arrowColor)

		helpers.DrawRect(s.img, startX-3, startY+6, 7, 1, color.RGBA{255, 255, 255, 255})
	} else if strings.HasPrefix(direction, "s") {
		helpers.DrawRect(s.img, startX-3, startY, 7, 1, color.RGBA{255, 255, 255, 255})

		helpers.DrawRect(s.img, startX, startY+2, 1, 5, arrowColor)

		s.img.Set(startX+1, startY+3, arrowColor)
		s.img.Set(startX+2, startY+4, arrowColor)
		s.img.Set(startX-1, startY+3, arrowColor)
		s.img.Set(startX-2, startY+4, arrowColor)
	}
}
