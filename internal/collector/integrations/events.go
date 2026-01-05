package integrations

import (
	"log"

	"github.com/AndreiBerezin/pixoo64/internal/collector/types"
)

func GetEventsData() (*types.EventsData, error) {
	log.Print("Getting events data...")
	return &types.EventsData{}, nil
}

func getLahtaHollData() (*types.EventsData, error) {
	return &types.EventsData{}, nil
}

func getGazpromArenaData() (*types.EventsData, error) {
	return &types.EventsData{}, nil
}
