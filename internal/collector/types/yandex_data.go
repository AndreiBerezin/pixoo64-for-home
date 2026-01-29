package types

type YandexData struct {
	CurrentWeather YandexCurrentWeather
	ByDays         []YandexDayWeather
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
	Morning YandexDayItem
	Day     YandexDayItem
	Evening YandexDayItem
	Night   YandexDayItem
}

type YandexDayItem struct {
	Icon        string
	Temperature int
}

type YandexSun struct {
	SunriseTime string
	SunsetTime  string
}

type YandexMoon struct {
	Icon        string
	NewMoonDate string
}
