package state

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"time"

	"github.com/AndreiBerezin/pixoo64/internal/collector"
	"github.com/AndreiBerezin/pixoo64/internal/pixoo64"
	"github.com/AndreiBerezin/pixoo64/internal/screens"
	"github.com/AndreiBerezin/pixoo64/internal/state/frame"
	"github.com/AndreiBerezin/pixoo64/pkg/env"
	"github.com/AndreiBerezin/pixoo64/pkg/log"
	"go.uber.org/zap"
)

const (
	deviceWidth  = 64
	deviceHeight = 64

	BottomScreenExtraWeather = 0
	BottomScreenMagneticSun  = 1

	drawInterval  = 1 * time.Minute
	errorInterval = 5 * time.Minute
)

type State struct {
	device              *pixoo64.Pixoo64
	collector           *collector.Collector
	currentBottomScreen int

	screens *screens.Screens
}

func New(collector *collector.Collector) *State {
	return &State{
		device:              pixoo64.New(deviceWidth, deviceHeight),
		collector:           collector,
		screens:             screens.New(deviceWidth, deviceHeight),
		currentBottomScreen: BottomScreenExtraWeather,
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

	if err = s.screens.DrawCurrentWeather(data.YandexData); err != nil {
		return fmt.Errorf("failed to draw weather screen: %w", err)
	}

	switch s.currentBottomScreen {
	case BottomScreenExtraWeather:
		if err = s.screens.DrawToday(data.YandexData); err != nil {
			return fmt.Errorf("failed to draw extra weather screen: %w", err)
		}

		s.currentBottomScreen = BottomScreenMagneticSun
	case BottomScreenMagneticSun:
		if err = s.screens.DrawMagneticSun(data.MagneticData, data.YandexData); err != nil {
			return fmt.Errorf("failed to draw magnetic sun screen: %w", err)
		}

		s.currentBottomScreen = BottomScreenExtraWeather
	}

	if env.IsDebug() {
		if err = devImgDraw(s.screens.Image()); err != nil {
			return fmt.Errorf("failed to draw dev image: %w", err)
		}
	}
	if err = pixoo64Draw(s.screens.Image()); err != nil {
		return fmt.Errorf("failed to draw pixoo64: %w", err)
	}

	log.Debug("data draw finished")

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

func pixoo64Draw(image *image.RGBA) error {
	pixoo64 := pixoo64.New(deviceWidth, deviceHeight)

	var frames []frame.Frame
	frame, err := frame.NewFrameImage(image, 400)
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
