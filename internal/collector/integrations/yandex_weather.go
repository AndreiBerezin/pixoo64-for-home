package integrations

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/AndreiBerezin/pixoo64/internal/collector/types"
	"github.com/AndreiBerezin/pixoo64/pkg/env"
	"github.com/AndreiBerezin/pixoo64/pkg/log"
)

type YandexWeather struct{}

func NewYandexWeather() *YandexWeather {
	return &YandexWeather{}
}

func (y *YandexWeather) Data() (*types.YandexData, error) {
	log.Info("Getting yandex data...")

	var response *yandexWeatherResponse
	var err error
	if env.IsDebug() {
		response, err = y.mockResponse()
	} else {
		response, err = y.callApi()
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get yandex weather response: %w", err)
	}

	todayData := response.Forecasts[0]

	var byDays []types.YandexDayWeather
	for _, forecast := range response.Forecasts {
		byDays = append(byDays, types.YandexDayWeather{
			Morning: types.YandexDayItem{
				Icon:        forecast.Parts["morning"].Icon.GetUrl(),
				Temperature: int(forecast.Parts["morning"].TempAverage),
			},
			Day: types.YandexDayItem{
				Icon:        forecast.Parts["day"].Icon.GetUrl(),
				Temperature: int(forecast.Parts["day"].TempAverage),
			},
			Evening: types.YandexDayItem{
				Icon:        forecast.Parts["evening"].Icon.GetUrl(),
				Temperature: int(forecast.Parts["evening"].TempAverage),
			},
			Night: types.YandexDayItem{
				Icon:        forecast.Parts["night"].Icon.GetUrl(),
				Temperature: int(forecast.Parts["night"].TempAverage),
			},
		})
	}

	return &types.YandexData{
		CurrentWeather: types.YandexCurrentWeather{
			Temperature:          int(response.Fact.Temperature),
			FeelsLikeTemperature: int(response.Fact.FeelsLikeTemperature),
			Icon:                 response.Fact.Icon.GetUrl(),
			WindSpeed:            int(response.Fact.WindSpeed),
			WindDirection:        response.Fact.WindDirection,
		},
		ByDays: byDays,
		Sun: types.YandexSun{
			SunriseTime: todayData.Sunrise,
			SunsetTime:  todayData.Sunset,
		},
		Moon: types.YandexMoon{
			Icon:         todayData.MoonCode.GetIcon(),
			MoonPhaseDay: y.getMoonPhaseDay(),
		},
	}, nil
}

func (y *YandexWeather) getMoonPhaseDay() int {
	synodicMonthDays := 29.530588853
	referenceNewMoonUTC := time.Date(2000, 1, 6, 18, 14, 0, 0, time.UTC)

	now := time.Now().UTC()
	elapsed := now.Sub(referenceNewMoonUTC)
	daysSinceReference := elapsed.Hours() / 24
	moonAge := math.Mod(daysSinceReference, synodicMonthDays)

	return int(math.Round(moonAge))
}

func (y *YandexWeather) mockResponse() (*yandexWeatherResponse, error) {
	file, err := os.Open("mocks/yandex_weather.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var response *yandexWeatherResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (y *YandexWeather) callApi() (*yandexWeatherResponse, error) {
	url := fmt.Sprintf("https://api.weather.yandex.ru/v2/forecast?lat=%s&lon=%s", os.Getenv("LAT"), os.Getenv("LON"))
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Yandex-Weather-Key", os.Getenv("YANDEX_WEATHER_KEY"))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response *yandexWeatherResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

type yandexWeatherResponse struct {
	Fact      fact       `json:"fact"`
	Forecasts []forecast `json:"forecasts"`
}

type fact struct {
	Temperature          float32 `json:"temp"`
	FeelsLikeTemperature float32 `json:"feels_like"`
	Icon                 Icon    `json:"icon"`
	WindSpeed            float32 `json:"wind_speed"`
	WindDirection        string  `json:"wind_dir"`
}

type forecast struct {
	Sunrise  string          `json:"sunrise"`
	Sunset   string          `json:"sunset"`
	MoonCode moonCode        `json:"moon_code"`
	Parts    map[string]part `json:"parts"`
}

type part struct {
	Icon        Icon    `json:"icon"`
	TempAverage float32 `json:"temp_avg"`
}

type Icon string

func (i Icon) GetUrl() string {
	return fmt.Sprintf("https://yastatic.net/weather/i/icons/funky/dark/%s.svg", string(i))
}

type moonCode int

func (m moonCode) GetIcon() string {
	switch int(m) {
	case 0: // full moon
		return "static/images/moon_0.png"
	case 1, 2, 3: // waning moon
		return "static/images/moon_1-3.png"
	case 4: // last quarter
		return "static/images/moon_4.png"
	case 5, 6, 7: // waning moon
		return "static/images/moon_5-7.png"
	case 8: // new moon
		return "static/images/moon_8.png"
	case 9, 10, 11: // waxing moon
		return "static/images/moon_9-11.png"
	case 12: // first quarter
		return "static/images/moon_12.png"
	case 13, 14, 15: // waxing moon
		return "static/images/moon_13-15.png"
	}

	return ""
}
