package types

type PressureData struct {
	Days []PressureDay
}

type PressureDay struct {
	Day   string
	Hours []PressureHour
}

type PressureHour struct {
	Hour     int
	Pressure float32
}
