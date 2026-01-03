package integrations

import (
	"log"

	intTypes "github.com/AndreiBerezin/pixoo64/internal/integrations/types"
)

func GetEventsData() (*intTypes.EventsData, error) {
	log.Print("Getting events data...")
	return &intTypes.EventsData{}, nil
}

func getLahtaHollData() (*intTypes.EventsData, error) {
	return &intTypes.EventsData{}, nil
}

func getGazpromArenaData() (*intTypes.EventsData, error) {
	return &intTypes.EventsData{}, nil
}
