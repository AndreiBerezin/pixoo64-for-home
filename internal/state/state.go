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

	drawInterval  = 1 * time.Minute
	errorInterval = 5 * time.Minute
)

type drawFn func(*types.CollectedData) error

type rotation struct {
	timerTopScreens []drawFn
	timerTopIdx     int

	defaultBottomScreens []drawFn
	defaultBottomIdx     int
}

type State struct {
	device       *pixoo64.Pixoo64
	collector    *collector.Collector
	screens      *screens.Screens
	timerManager *timer.Manager
	rotation     rotation
}

func New(collector *collector.Collector, timerManager *timer.Manager) *State {
	sc := screens.New(deviceWidth, deviceHeight)
	return &State{
		device:       pixoo64.New(deviceWidth, deviceHeight),
		collector:    collector,
		screens:      sc,
		timerManager: timerManager,
		rotation: rotation{
			timerTopScreens: []drawFn{
				sc.DrawTopCurrentWeather,
				sc.DrawTopExtraWeater,
			},
			defaultBottomScreens: []drawFn{
				sc.DrawBottomExtraWeater,
				sc.DrawBottomMagneticPressure,
				sc.DrawBottomSunMoon,
			},
		},
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

	if activeTimer := s.timerManager.ActiveTimer(); activeTimer != nil {
		if err = s.drawTimerState(data, activeTimer); err != nil {
			return fmt.Errorf("failed to draw timer state: %w", err)
		}
	} else {
		if err = s.drawDefaultState(data); err != nil {
			return fmt.Errorf("failed to draw default state: %w", err)
		}
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

func (s *State) drawTimerState(data *types.CollectedData, activeTimer *timer.ActiveTimer) error {
	if err := s.rotation.timerTopScreens[s.rotation.timerTopIdx](data); err != nil {
		return fmt.Errorf("failed to draw top screen: %w", err)
	}
	s.rotation.timerTopIdx = (s.rotation.timerTopIdx + 1) % len(s.rotation.timerTopScreens)

	if err := s.screens.DrawBottomTimer(activeTimer.From, activeTimer.To); err != nil {
		return fmt.Errorf("failed to draw bottom timer screen: %w", err)
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
	s.rotation.defaultBottomIdx = (s.rotation.defaultBottomIdx + 1) % len(s.rotation.defaultBottomScreens)

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
