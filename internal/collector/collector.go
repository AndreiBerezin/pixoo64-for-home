package collector

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/AndreiBerezin/pixoo64/internal/integrations"
	intTypes "github.com/AndreiBerezin/pixoo64/internal/integrations/types"
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
}

type CollectedData struct {
	YandexData   *intTypes.YandexData
	MagneticData *intTypes.MagneticData
	EventsData   *intTypes.EventsData
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
	}
}

func (c *Collector) Start() {
	go func() {
		for {
			errors := c.collect()
			for _, err := range errors {
				log.Print(err)
			}

			time.Sleep(5 * time.Second)
		}
	}()
}

func (c *Collector) collect() []error {
	c.Lock()
	defer c.Unlock()

	var errors []error

	for name, meta := range c.meta {
		if !meta.isExpired() {
			continue
		}

		meta.collectedAt = time.Now()
		c.meta[name] = meta

		switch name {
		case YandexDataName:
			yandexData, err := integrations.GetYandexData()
			if err != nil {
				errors = append(errors, fmt.Errorf("failed to get yandex data: %w", err))
				continue
			}
			c.collectedData.YandexData = yandexData
		case MagneticDataName:
			magneticData, err := integrations.GetMagneticData()
			if err != nil {
				errors = append(errors, fmt.Errorf("failed to get magnetic data: %w", err))
				continue
			}
			c.collectedData.MagneticData = magneticData
		case EventsDataName:
			eventsData, err := integrations.GetEventsData()
			if err != nil {
				errors = append(errors, fmt.Errorf("failed to get events data: %w", err))
				continue
			}
			c.collectedData.EventsData = eventsData
		}
	}

	return errors
}

// todo: тут нифига не лочится
func (c *Collector) GetCollectedData() CollectedData {
	c.RLock()
	defer c.RUnlock()

	return *c.collectedData
}
