package types

import "time"

type EventsData struct {
	Events []Event
}

type Event struct {
	Icon string
	Time time.Time
}
