package screens

import (
	"fmt"
	"image/png"
	"os"
	"time"

	"github.com/AndreiBerezin/pixoo64/internal/collector"
	"github.com/AndreiBerezin/pixoo64/internal/drawer"
	"github.com/AndreiBerezin/pixoo64/internal/frame"
	"github.com/AndreiBerezin/pixoo64/internal/pixoo64"
	"github.com/AndreiBerezin/pixoo64/pkg/env"
	"github.com/AndreiBerezin/pixoo64/pkg/log"
	"go.uber.org/zap"
)

const (
	BottomScreenExtraWeather = 0
	BottomScreenMagneticSun  = 1

	drawInterval  = 1 * time.Minute
	errorInterval = 5 * time.Minute
)

type State struct {
	collector           *collector.Collector
	currentBottomScreen int

	drawer             *drawer.Drawer
	weatherScreen      *CurrentWeatherScreen
	extraWeatherScreen *ExtraWeatherScreen
	magneticSunScreen  *MagneticSunScreen
}

func NewState(collector *collector.Collector) *State {
	draw := drawer.New()

	return &State{
		collector:           collector,
		currentBottomScreen: BottomScreenExtraWeather,
		drawer:              draw,
		weatherScreen:       NewCurrentWeatherScreen(draw),
		extraWeatherScreen:  NewExtraWeatherScreen(draw),
		magneticSunScreen:   NewMagneticSunScreen(draw),
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
	data, err := s.collector.GetCollectedData()
	if err != nil {
		return fmt.Errorf("failed to get collected data: %w", err)
	}

	s.drawer.Reset()

	if err = s.weatherScreen.DrawStatic(data.YandexData); err != nil {
		return fmt.Errorf("failed to draw weather screen: %w", err)
	}

	switch s.currentBottomScreen {
	case BottomScreenExtraWeather:
		if err = s.extraWeatherScreen.DrawTodayStatic(data.YandexData); err != nil {
			return fmt.Errorf("failed to draw extra weather screen: %w", err)
		}

		s.currentBottomScreen = BottomScreenMagneticSun
	case BottomScreenMagneticSun:
		if err = s.magneticSunScreen.DrawStatic(data.MagneticData, data.YandexData); err != nil {
			return fmt.Errorf("failed to draw magnetic sun screen: %w", err)
		}

		s.currentBottomScreen = BottomScreenExtraWeather
	}

	if env.IsDebug() {
		if err = devImgDraw(s.drawer); err != nil {
			return fmt.Errorf("failed to draw dev image: %w", err)
		}
	}
	if err = pixoo64Draw(s.drawer); err != nil {
		return fmt.Errorf("failed to draw pixoo64: %w", err)
	}

	log.Debug("data draw finished")

	return nil
}

func devImgDraw(drawer *drawer.Drawer) error {
	filename := "dev_img.png"
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	png.Encode(file, drawer.Image())

	log.Debug("success draw dev image to " + filename)

	return nil
}

func pixoo64Draw(drawer *drawer.Drawer) error {
	pixoo64 := pixoo64.NewPixoo64()

	var frames []frame.Frame
	frame, err := frame.NewFrameImage(drawer.Image(), 400)
	if err != nil {
		return fmt.Errorf("failed to create frame: %w", err)
	}
	frames = append(frames, *frame)

	pixoo64.ResetHttpGifId()
	if err = pixoo64.SendHttpGif(0, frames); err != nil {
		return fmt.Errorf("failed to send http gif: %w", err)
	}

	log.Debug("success draw on pixoo64")

	return nil
}
