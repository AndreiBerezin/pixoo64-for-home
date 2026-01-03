package integrations

import (
	"log"

	intTypes "github.com/AndreiBerezin/pixoo64/internal/integrations/types"
)

func GetMagneticData() (*intTypes.MagneticData, error) {
	log.Print("Getting magnetic data...")
	return &intTypes.MagneticData{}, nil
}
