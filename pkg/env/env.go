package env

import (
	"os"
	"slices"
)

func IsDebug() bool {
	return slices.Contains([]string{"development", "dev"}, os.Getenv("ENV"))
}

func Lang() string {
	lang := os.Getenv("APP_LANG")
	if slices.Contains([]string{"ru", "en"}, lang) {
		return lang
	}
	return "en"
}
