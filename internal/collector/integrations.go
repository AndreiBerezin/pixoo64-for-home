package collector

import (
	"time"

	"github.com/AndreiBerezin/pixoo64/internal/collector/integrations"
	"github.com/AndreiBerezin/pixoo64/internal/collector/types"
	"github.com/AndreiBerezin/pixoo64/pkg/log"
	"go.uber.org/zap"
)

type Integration struct {
	ttl         time.Duration
	collectedAt time.Time
	collectFunc func(collectedData *types.CollectedData)
}

func (i *Integration) isExpired() bool {
	return time.Since(i.collectedAt) > i.ttl
}

func (i *Integration) Collect(collectedData *types.CollectedData) {
	if !i.isExpired() {
		return
	}
	i.collectedAt = time.Now()
	i.collectFunc(collectedData)
}

func NewYandexIntegration(ttl time.Duration) *Integration {
	return &Integration{
		ttl: ttl,
		collectFunc: func(collectedData *types.CollectedData) {
			integration := integrations.NewYandexWeather()
			yandexData, err := integration.Data()
			if err != nil {
				log.Error("failed to get yandex data: ", zap.Error(err))
				return
			}
			collectedData.YandexData = yandexData
		},
	}
}

func NewMagneticIntegration(ttl time.Duration) *Integration {
	return &Integration{
		ttl: ttl,
		collectFunc: func(collectedData *types.CollectedData) {
			integration := integrations.NewXras()
			magneticData, err := integration.Data()
			if err != nil {
				log.Error("failed to get magnetic data: ", zap.Error(err))
				return
			}
			collectedData.MagneticData = magneticData
		},
	}
}

func NewEventsIntegration(ttl time.Duration) *Integration {
	return &Integration{
		ttl: ttl,
		collectFunc: func(collectedData *types.CollectedData) {
			integration := integrations.NewEvents()
			eventsData, err := integration.Data()
			if err != nil {
				log.Error("failed to get events data: ", zap.Error(err))
				return
			}
			collectedData.EventsData = eventsData
		},
	}
}
