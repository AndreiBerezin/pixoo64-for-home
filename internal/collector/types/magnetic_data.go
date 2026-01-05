package types

type MagneticData struct {
	ByHours []MagneticHour
	ByDays  []MagneticDay
}

type MagneticDay struct {
	Day   int
	Level float32
}

type MagneticHour struct {
	Hour  int
	Level float32
}
