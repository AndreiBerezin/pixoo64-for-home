package screens

import (
	"fmt"
	"image/color"
	"time"

	"github.com/AndreiBerezin/pixoo64/internal/collector/types"
	"github.com/AndreiBerezin/pixoo64/internal/screens/image/fonts"
	"github.com/AndreiBerezin/pixoo64/pkg/i18n"
)

const splitterName = "|"

func (s *Screens) DrawTopExtraWeater(data *types.CollectedData) error {
	return s.drawExtraWeather(18, data.YandexData)
}

func (s *Screens) DrawBottomExtraWeater(data *types.CollectedData) error {
	return s.drawExtraWeather(45, data.YandexData)
}

func (s *Screens) drawExtraWeather(startY int, data *types.YandexData) error {
	if data == nil {
		return nil
	}

	offsetX := 2
	for _, item := range s.futureWeatherItems(data) {
		if item.name == splitterName {
			s.image.DrawRect(offsetX-2, startY-6, 1, 23, color.RGBA{100, 100, 100, 255})
			offsetX += 2
			continue
		}

		temperature := fmt.Sprintf("%ḋ", item.dataItem.Temperature)
		if item.dataItem.Temperature > 0 {
			temperature = "+" + temperature
		} else if item.dataItem.Temperature == 0 {
			temperature = " " + temperature
		}
		s.image.DrawString(temperature, offsetX, startY, color.RGBA{255, 255, 255, 255}, fonts.FontMicro5Normal)

		err := s.image.DrawSVGFromURL(item.dataItem.Icon, offsetX, startY, 12)
		if err != nil {
			return fmt.Errorf("failed to draw icon: %w", err)
		}
		s.image.DrawString(item.name, offsetX+4, startY+16, color.RGBA{255, 255, 255, 255}, fonts.FontTiny5Normal)

		offsetX += 16
	}

	return nil
}

func (s *Screens) futureWeatherItems(data *types.YandexData) []item {
	morningLabel := i18n.MorningLabel()
	dayLabel := i18n.DayLabel()
	eveningLabel := i18n.EveningLabel()
	nightLabel := i18n.NightLabel()

	currentHour := time.Now().Hour()
	if currentHour <= 11 {
		return []item{
			{name: morningLabel, dataItem: data.ByDays[0].Morning},
			{name: dayLabel, dataItem: data.ByDays[0].Day},
			{name: eveningLabel, dataItem: data.ByDays[0].Evening},
			{name: nightLabel, dataItem: data.ByDays[0].Night},
		}
	} else if currentHour <= 17 {
		return []item{
			{name: dayLabel, dataItem: data.ByDays[0].Day},
			{name: eveningLabel, dataItem: data.ByDays[0].Evening},
			{name: nightLabel, dataItem: data.ByDays[0].Night},
			{name: splitterName},
			{name: morningLabel, dataItem: data.ByDays[1].Morning},
		}
	} else if currentHour <= 21 {
		return []item{
			{name: eveningLabel, dataItem: data.ByDays[0].Evening},
			{name: nightLabel, dataItem: data.ByDays[0].Night},
			{name: splitterName},
			{name: morningLabel, dataItem: data.ByDays[1].Morning},
			{name: dayLabel, dataItem: data.ByDays[1].Day},
		}
	} else if currentHour <= 23 {
		return []item{
			{name: nightLabel, dataItem: data.ByDays[0].Night},
			{name: splitterName},
			{name: morningLabel, dataItem: data.ByDays[1].Morning},
			{name: dayLabel, dataItem: data.ByDays[1].Day},
			{name: eveningLabel, dataItem: data.ByDays[1].Evening},
		}
	} else {
		return []item{}
	}
}

type item struct {
	name     string
	dataItem types.YandexDayItem
}
