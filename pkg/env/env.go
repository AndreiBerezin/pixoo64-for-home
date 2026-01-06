package env

import (
	"os"
	"slices"
)

func IsDebug() bool {
	return slices.Contains([]string{"development", "dev"}, os.Getenv("ENV"))
}
