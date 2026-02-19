package timer

import (
	"time"

	"github.com/robfig/cron/v3"
)

type Timer struct {
	schedule          cron.Schedule
	notifyDurationMin int
}

type ActiveTimer struct {
	From time.Time
	To   time.Time
}

func (a *ActiveTimer) IsBoundary() bool {
	now := time.Now()
	atBoundary := func(t time.Time) bool {
		return now.Hour() == t.Hour() && now.Minute() == t.Minute()
	}
	return atBoundary(a.From) || atBoundary(a.To)
}

func (t *Timer) Active(now time.Time) (*ActiveTimer, bool) {
	nowMinute := now.Truncate(time.Minute)
	notifyDuration := time.Duration(t.notifyDurationMin) * time.Minute
	candidate := t.schedule.Next(nowMinute.Add(-1 * time.Nanosecond))

	if candidate.Sub(nowMinute) <= notifyDuration {
		return &ActiveTimer{
			From: candidate.Add(-notifyDuration),
			To:   candidate,
		}, true
	}

	return nil, false
}
