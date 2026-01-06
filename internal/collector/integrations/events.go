package integrations

import (
	"github.com/AndreiBerezin/pixoo64/internal/collector/types"
	"github.com/AndreiBerezin/pixoo64/pkg/log"
)

type Events struct{}

func NewEvents() *Events {
	return &Events{}
}

func (e *Events) Data() (*types.EventsData, error) {
	log.Info("Getting events data...")
	return &types.EventsData{}, nil
}

func getLahtaHollData() (*types.EventsData, error) {
	return &types.EventsData{}, nil
}

func getGazpromArenaData() (*types.EventsData, error) {
	return &types.EventsData{}, nil
}
