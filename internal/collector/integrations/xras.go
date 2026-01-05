package integrations

import (
	"log"

	"github.com/AndreiBerezin/pixoo64/internal/collector/types"
)

func GetMagneticData() (*types.MagneticData, error) {
	log.Print("Getting magnetic data...")
	return &types.MagneticData{}, nil
}
