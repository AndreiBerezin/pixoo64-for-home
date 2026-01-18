package collector

import (
	"fmt"
	"sync"
	"time"

	"github.com/AndreiBerezin/pixoo64/internal/collector/types"
)

const (
	collectInterval = 5 * time.Minute
)

type Collector struct {
	sync.RWMutex
	collectedData *types.CollectedData
	integrations  []*Integration
}

func New() *Collector {
	collector := &Collector{
		collectedData: &types.CollectedData{},
		integrations: []*Integration{
			NewYandexIntegration(1 * time.Hour),
			NewMagneticIntegration(1 * time.Hour),
			NewEventsIntegration(1 * time.Hour),
		},
	}
	collector.collect()

	return collector
}

func (c *Collector) Start() {
	go func() {
		for {
			c.collect()

			time.Sleep(collectInterval)
		}
	}()
}

func (c *Collector) collect() {
	c.Lock()
	defer c.Unlock()

	for _, integration := range c.integrations {
		integration.Collect(c.collectedData)
	}
}

func (c *Collector) CollectedData() (*types.CollectedData, error) {
	c.RLock()
	defer c.RUnlock()

	collectedData, err := c.collectedData.Clone()
	if err != nil {
		return nil, fmt.Errorf("failed to clone collected data: %w", err)
	}

	return collectedData, nil
}
