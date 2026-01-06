package integrations

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

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

	return &types.YandexData{
		CurrentWeather: types.YandexCurrentWeather{
			Temperature:          int(response.Fact.Temperature),
			FeelsLikeTemperature: int(response.Fact.FeelsLikeTemperature),
			Icon:                 response.Fact.Icon.GetUrl(),
			WindSpeed:            int(response.Fact.WindSpeed),
			WindDirection:        response.Fact.WindDirection,
		},
		DayWeather: types.YandexDayWeather{
			Items: []types.YandexDayItem{
				{
					Name:        "у",
					Icon:        todayData.Parts["morning"].Icon.GetUrl(),
					Temperature: int(todayData.Parts["morning"].TempAverage),
				},
				{
					Name:        "д",
					Icon:        todayData.Parts["day"].Icon.GetUrl(),
					Temperature: int(todayData.Parts["day"].TempAverage),
				},
				{
					Name:        "в",
					Icon:        todayData.Parts["evening"].Icon.GetUrl(),
					Temperature: int(todayData.Parts["evening"].TempAverage),
				},
				{
					Name:        "н",
					Icon:        todayData.Parts["night"].Icon.GetUrl(),
					Temperature: int(todayData.Parts["night"].TempAverage),
				},
			},
		},
		Sun: types.YandexSun{
			SunriseTime: todayData.Sunrise,
			SunsetTime:  todayData.Sunset,
		},
		Moon: types.YandexMoon{ // todo: вытащить откуда нибудь
			MoonPhase: "",
			MoonDay:   1,
		},
	}, nil
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
	url := fmt.Sprintf("https://api.weather.yandex.ru/v2/forecast?lat=%s&lon=%s", os.Getenv("YANDEX_WEATHER_LAT"), os.Getenv("YANDEX_WEATHER_LON"))
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
	MoonCode int             `json:"moon_code"`
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
