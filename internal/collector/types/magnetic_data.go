package types

type MagneticData struct {
	Days []MagneticDay
}

type MagneticDay struct {
	Hours []MagneticHour
}

type MagneticHour struct {
	Hour  int
	Level float32
}
