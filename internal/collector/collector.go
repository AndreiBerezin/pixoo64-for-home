package collector

import (
	"sync"
	"time"

	"github.com/AndreiBerezin/pixoo64/internal/collector/integrations"
	"github.com/AndreiBerezin/pixoo64/internal/collector/types"
	"github.com/AndreiBerezin/pixoo64/pkg/log"
	"go.uber.org/zap"
)

const (
	YandexDataName   = "yandex_data"
	MagneticDataName = "magnetic_data"
	EventsDataName   = "events_data"
)

type Collector struct {
	sync.RWMutex
	collectedData *CollectedData
	meta          map[string]metaItem

	yandexWeather *integrations.YandexWeather
	xras          *integrations.Xras
	events        *integrations.Events
}

type CollectedData struct {
	YandexData   *types.YandexData
	MagneticData *types.MagneticData
	EventsData   *types.EventsData
}

type metaItem struct {
	ttl         time.Duration
	collectedAt time.Time
}

func (c metaItem) isExpired() bool {
	return time.Since(c.collectedAt) > c.ttl
}

func NewCollector() *Collector {
	return &Collector{
		collectedData: &CollectedData{},
		meta: map[string]metaItem{
			YandexDataName: {
				ttl:         1 * time.Hour,
				collectedAt: time.Unix(0, 0),
			},
			MagneticDataName: {
				ttl:         1 * time.Hour,
				collectedAt: time.Unix(0, 0),
			},
			EventsDataName: {
				ttl:         1 * time.Hour,
				collectedAt: time.Unix(0, 0),
			},
		},
		yandexWeather: integrations.NewYandexWeather(),
		xras:          integrations.NewXras(),
		events:        integrations.NewEvents(),
	}
}

func (c *Collector) Start() {
	go func() {
		for {
			c.collect()

			time.Sleep(5 * time.Second)
		}
	}()
}

func (c *Collector) collect() {
	c.Lock()
	defer c.Unlock()

	for name, meta := range c.meta {
		if !meta.isExpired() {
			continue
		}

		meta.collectedAt = time.Now()
		c.meta[name] = meta

		switch name {
		case YandexDataName:
			yandexData, err := c.yandexWeather.Data()
			if err != nil {
				log.Error("failed to get yandex data: ", zap.Error(err))
				continue
			}
			c.collectedData.YandexData = yandexData
		case MagneticDataName:
			magneticData, err := c.xras.Data()
			if err != nil {
				log.Error("failed to get magnetic data: ", zap.Error(err))
				continue
			}
			c.collectedData.MagneticData = magneticData
		case EventsDataName:
			eventsData, err := c.events.Data()
			if err != nil {
				log.Error("failed to get events data: ", zap.Error(err))
				continue
			}
			c.collectedData.EventsData = eventsData
		}
	}
}

// todo: тут нифига не лочится
func (c *Collector) GetCollectedData() CollectedData {
	c.RLock()
	defer c.RUnlock()

	return *c.collectedData
}
