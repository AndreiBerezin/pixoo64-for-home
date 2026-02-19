package timer

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/robfig/cron/v3"
)

type Manager struct {
	timers []Timer
}

type timerConfig struct {
	At                string `json:"at"`
	NotifyDurationMin int    `json:"notify_duration_min"`
}

func NewManager() (*Manager, error) {
	raw := os.Getenv("TIMERS")
	if raw == "" {
		return &Manager{}, nil
	}

	var configs []timerConfig
	if err := json.Unmarshal([]byte(raw), &configs); err != nil {
		return nil, fmt.Errorf("failed to parse TIMERS env: %w", err)
	}

	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)

	timers := make([]Timer, 0, len(configs))
	for _, c := range configs {
		schedule, err := parser.Parse(c.At)
		if err != nil {
			return nil, fmt.Errorf("failed to parse cron expression %q: %w", c.At, err)
		}
		timers = append(timers, Timer{
			schedule:          schedule,
			notifyDurationMin: c.NotifyDurationMin,
		})
	}

	return &Manager{timers: timers}, nil
}

func (m *Manager) ActiveTimer() *ActiveTimer {
	now := time.Now()
	for i := range m.timers {
		if active, ok := m.timers[i].Active(now); ok {
			return active
		}
	}
	return nil
}
