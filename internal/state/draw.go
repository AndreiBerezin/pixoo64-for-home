package state

import (
	"fmt"
	"image"
	"image/png"
	"os"

	"github.com/AndreiBerezin/pixoo64/internal/collector/types"
	"github.com/AndreiBerezin/pixoo64/pkg/env"
	"github.com/AndreiBerezin/pixoo64/pkg/log"
)

func (s *State) draw() error {
	if err := s.prepareScreens(); err != nil {
		return fmt.Errorf("failed to prepare screens: %w", err)
	}

	if env.IsDebug() {
		if err := devImgDraw(s.screens.Image()); err != nil {
			return fmt.Errorf("failed to draw dev image: %w", err)
		}
	}

	if err := s.device.DrawImage(s.screens.Image()); err != nil {
		return fmt.Errorf("failed to draw pixoo64: %w", err)
	}

	log.Debug("data draw finished")

	return nil
}

func (s *State) prepareScreens() error {
	data, err := s.collector.CollectedData()
	if err != nil {
		return fmt.Errorf("failed to get collected data: %w", err)
	}

	s.screens.Reset()

	if err = s.screens.DrawHeader(); err != nil {
		return fmt.Errorf("failed to draw header screen: %w", err)
	}

	switch {
	case s.onAir.Load():
		if err = s.drawOnAirState(data); err != nil {
			return fmt.Errorf("failed to draw on air state: %w", err)
		}
	case s.timerManager.ActiveTimer() != nil:
		if err = s.drawTimerState(data); err != nil {
			return fmt.Errorf("failed to draw timer state: %w", err)
		}
	default:
		if err = s.drawDefaultState(data); err != nil {
			return fmt.Errorf("failed to draw default state: %w", err)
		}
	}

	return nil
}

func (s *State) drawOnAirState(data *types.CollectedData) error {
	if err := s.rotation.customTopScreens[s.rotation.customTopIdx](data); err != nil {
		return fmt.Errorf("failed to draw top screen: %w", err)
	}
	s.rotation.nextCustomTop()

	if err := s.screens.DrawBottomOnAir(s.onAirStartTime); err != nil {
		return fmt.Errorf("failed to draw bottom on air screen: %w", err)
	}

	return nil
}

func (s *State) drawTimerState(data *types.CollectedData) error {
	activeTimer := s.timerManager.ActiveTimer()
	if activeTimer == nil {
		return fmt.Errorf("active timer is nil")
	}

	if err := s.rotation.customTopScreens[s.rotation.customTopIdx](data); err != nil {
		return fmt.Errorf("failed to draw top screen: %w", err)
	}
	s.rotation.nextCustomTop()

	if err := s.screens.DrawBottomTimer(activeTimer.From, activeTimer.To); err != nil {
		return fmt.Errorf("failed to draw bottom timer screen: %w", err)
	}

	if activeTimer.IsBoundary() {
		s.device.PlayBuzzer(100, 100, 500)
	} else {
		s.device.PlayBuzzer(100, 0, 100)
	}

	return nil
}

func (s *State) drawDefaultState(data *types.CollectedData) error {
	if err := s.screens.DrawTopCurrentWeather(data); err != nil {
		return fmt.Errorf("failed to draw top current weather screen: %w", err)
	}

	if err := s.rotation.defaultBottomScreens[s.rotation.defaultBottomIdx](data); err != nil {
		return fmt.Errorf("failed to draw bottom screen: %w", err)
	}
	s.rotation.nextDefaultBottom()

	return nil
}

func devImgDraw(image *image.RGBA) error {
	filename := "dev_img.png"
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	png.Encode(file, image)

	log.Debug("success draw dev image to " + filename)

	return nil
}
