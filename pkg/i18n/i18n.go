package i18n

import "github.com/AndreiBerezin/pixoo64/pkg/env"

var windDirections = map[string]map[string]string{
	"ru": {"n": "с", "s": "ю", "e": "в", "w": "з", "nw": "св", "ne": "сз", "sw": "юз", "se": "юв"},
	"en": {"n": "n", "s": "s", "e": "e", "w": "w", "nw": "nw", "ne": "ne", "sw": "sw", "se": "se"},
}

var dayPeriodsLabels = map[string]map[string]string{
	"ru": {"morning": "у", "day": "д", "evening": "в", "night": "н"},
	"en": {"morning": "m", "day": "d", "evening": "e", "night": "n"},
}

func WindDirection(direction string) string {
	return windDirections[env.Lang()][direction]
}

func MorningLabel() string { return dayPeriodsLabels[env.Lang()]["morning"] }
func DayLabel() string     { return dayPeriodsLabels[env.Lang()]["day"] }
func EveningLabel() string { return dayPeriodsLabels[env.Lang()]["evening"] }
func NightLabel() string   { return dayPeriodsLabels[env.Lang()]["night"] }
