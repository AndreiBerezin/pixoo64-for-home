package types

type YandexData struct {
	CurrentWeather YandexCurrentWeather
	DayWeather     YandexDayWeather
	Sun            YandexSun
	Moon           YandexMoon
}

type YandexCurrentWeather struct {
	Temperature          int
	FeelsLikeTemperature int
	Icon                 string
	WindSpeed            int
	WindDirection        string
}

type YandexDayWeather struct {
	Items []YandexDayItem
}

type YandexDayItem struct {
	Name        string
	Icon        string
	Temperature int
}

type YandexSun struct {
	SunriseTime string
	SunsetTime  string
}

type YandexMoon struct {
	MoonPhase string
	MoonDay   int
}
