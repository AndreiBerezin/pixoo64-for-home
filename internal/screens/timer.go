package screens

import (
	"image/color"
	"time"

	"github.com/AndreiBerezin/pixoo64/internal/screens/image/fonts"
)

func (s *Screens) DrawBottomTimer(from time.Time, to time.Time) error {
	return s.drawTimer(45, from, to)
}

func (s *Screens) drawTimer(startY int, from time.Time, to time.Time) error {
	allDuration := to.Sub(from)
	diffDuration := time.Until(to)

	progressColor := color.RGBA{255, 255, 255, 255}
	if diffDuration.Minutes() < 10 {
		progressColor = color.RGBA{255, 100, 100, 255}
	}

	diffTime := time.Unix(int64(diffDuration.Seconds())+60, 0).UTC()
	s.image.DrawString(diffTime.Format("15:04"), 2, startY+3, progressColor, fonts.FontMicro5Big)
	s.image.DrawRect(2, startY+5, 60, 12, color.RGBA{50, 50, 50, 255})

	progress := min(1.0, max(0.0, 1-diffDuration.Seconds()/allDuration.Seconds()))
	s.image.DrawRect(2, startY+5, int(progress*60), 12, progressColor)

	return nil
}
