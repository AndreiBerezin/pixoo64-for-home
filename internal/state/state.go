package state

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"time"

	"github.com/AndreiBerezin/pixoo64/internal/collector"
	"github.com/AndreiBerezin/pixoo64/internal/collector/types"
	"github.com/AndreiBerezin/pixoo64/internal/pixoo64"
	"github.com/AndreiBerezin/pixoo64/internal/screens"
	"github.com/AndreiBerezin/pixoo64/internal/timer"
	"github.com/AndreiBerezin/pixoo64/pkg/env"
	"github.com/AndreiBerezin/pixoo64/pkg/log"
	"go.uber.org/zap"
)

const (
	deviceWidth  = 64
	deviceHeight = 64

	BottomScreenExtraWeather     = 0
	BottomScreenMagneticPressure = 1
	BottomScreenSunMoon          = 2

	drawInterval  = 1 * time.Minute
	errorInterval = 5 * time.Minute
)

type State struct {
	device              *pixoo64.Pixoo64
	collector           *collector.Collector
	currentBottomScreen int
	screens             *screens.Screens
	timerManager        *timer.Manager
}

func New(collector *collector.Collector, timerManager *timer.Manager) *State {
	return &State{
		device:              pixoo64.New(deviceWidth, deviceHeight),
		collector:           collector,
		screens:             screens.New(deviceWidth, deviceHeight),
		currentBottomScreen: BottomScreenExtraWeather,
		timerManager:        timerManager,
	}
}

func (s *State) Start() {
	go func() {
		for {
			if err := s.draw(); err != nil {
				log.Error("failed to draw screen: ", zap.Error(err))
				time.Sleep(errorInterval)
				continue
			}

			time.Sleep(drawInterval)
		}
	}()
}

func (s *State) draw() error {
	data, err := s.collector.CollectedData()
	if err != nil {
		return fmt.Errorf("failed to get collected data: %w", err)
	}

	s.screens.Reset()

	if err = s.screens.DrawHeader(); err != nil {
		return fmt.Errorf("failed to draw header screen: %w", err)
	}
	if err = s.drawTopState(data); err != nil {
		return fmt.Errorf("failed to draw top state: %w", err)
	}
	if err := s.drawBottomState(data); err != nil {
		return fmt.Errorf("failed to draw bottom state: %w", err)
	}

	if env.IsDebug() {
		if err = devImgDraw(s.screens.Image()); err != nil {
			return fmt.Errorf("failed to draw dev image: %w", err)
		}
	}

	if err = s.device.DrawImage(s.screens.Image()); err != nil {
		return fmt.Errorf("failed to draw pixoo64: %w", err)
	}

	log.Debug("data draw finished")

	return nil
}

func (s *State) drawTopState(data *types.CollectedData) error {
	if active := s.timerManager.ActiveTimer(); active != nil {
		if err := s.screens.DrawTopTimer(active.From, active.To); err != nil {
			return fmt.Errorf("failed to draw timer screen: %w", err)
		}

		if active.IsBoundary() {
			s.device.PlayBuzzer(100, 100, 500)
		} else {
			s.device.PlayBuzzer(100, 0, 100)
		}
	} else {
		if err := s.screens.DrawTopCurrentWeather(data.YandexData); err != nil {
			return fmt.Errorf("failed to draw current weather screen: %w", err)
		}
	}

	return nil
}

func (s *State) drawBottomState(data *types.CollectedData) error {
	switch s.currentBottomScreen {
	case BottomScreenExtraWeather:
		if err := s.screens.DrawBottomExtraWeater(data.YandexData); err != nil {
			return fmt.Errorf("failed to draw extra weather screen: %w", err)
		}
		s.currentBottomScreen = BottomScreenMagneticPressure
	case BottomScreenMagneticPressure:
		if err := s.screens.DrawBottomMagneticPressure(data.MagneticData, data.PressureData); err != nil {
			return fmt.Errorf("failed to draw magnetic pressure screen: %w", err)
		}
		s.currentBottomScreen = BottomScreenSunMoon
	case BottomScreenSunMoon:
		if err := s.screens.DrawBottomSunMoon(data.YandexData); err != nil {
			return fmt.Errorf("failed to draw sun moon screen: %w", err)
		}
		s.currentBottomScreen = BottomScreenExtraWeather
	}

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
