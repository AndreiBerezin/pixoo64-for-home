package timer

import (
	"testing"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTimer_Active(t *testing.T) {
	// active at 8:20-8:40
	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	schedule, err := parser.Parse("40 8 * * 1-5")
	require.NoError(t, err)
	timer := Timer{schedule: schedule, notifyDurationMin: 20}

	monday := time.Date(2026, 2, 16, 0, 0, 0, 0, time.Local)

	tests := []struct {
		now        time.Time
		wantActive bool
	}{
		{monday.Add(8*time.Hour + 19*time.Minute + 37*time.Second), false},
		{monday.Add(8*time.Hour + 20*time.Minute), true},
		{monday.Add(8*time.Hour + 20*time.Minute + 45*time.Second), true},
		{monday.Add(8*time.Hour + 30*time.Minute + 12*time.Second), true},
		{monday.Add(8*time.Hour + 40*time.Minute + 59*time.Second), true},
		{monday.Add(8*time.Hour + 41*time.Minute + 3*time.Second), false},
		{monday.AddDate(0, 0, 6).Add(8*time.Hour + 30*time.Minute + 22*time.Second), false},
	}

	for _, tt := range tests {
		t.Run(tt.now.Format("Mon 15:04:05"), func(t *testing.T) {
			_, ok := timer.Active(tt.now)
			assert.Equal(t, tt.wantActive, ok)
		})
	}
}
