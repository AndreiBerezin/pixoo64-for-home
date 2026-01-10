package screens

import (
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
)

type State struct {
	collector           *collector.Collector
	currentBottomScreen int
}

func NewState(collector *collector.Collector) *State {
	return &State{
		collector:           collector,
		currentBottomScreen: BottomScreenExtraWeather,
	}
}

func (s *State) Start() {
	go func() {
		for {
			data, err := s.collector.GetCollectedData()
			if err != nil {
				log.Fatal("failed to get collected data: ", zap.Error(err))
			}

			drawer := drawer.NewDrawer()

			weatherScreen := NewCurrentWeatherScreen(drawer)

			err = weatherScreen.DrawStatic(data.YandexData)
			if err != nil {
				log.Fatal("failed to draw weather screen: ", zap.Error(err))
			}

			switch s.currentBottomScreen {
			case BottomScreenExtraWeather:
				extraWeatherScreen, err := NewExtraWeatherScreen(drawer)
				if err != nil {
					log.Fatal("failed to create extra weather screen: ", zap.Error(err))
				}
				err = extraWeatherScreen.DrawTodayStatic(data.YandexData)
				if err != nil {
					log.Fatal("failed to draw extra weather screen: ", zap.Error(err))
				}

				s.currentBottomScreen = BottomScreenMagneticSun
			case BottomScreenMagneticSun:
				magneticSunScreen, err := NewMagneticSunScreen(drawer)
				if err != nil {
					log.Fatal("failed to create magnetic screen: ", zap.Error(err))
				}
				err = magneticSunScreen.DrawStatic(data.MagneticData, data.YandexData)
				if err != nil {
					log.Fatal("failed to draw magnetic screen: ", zap.Error(err))
				}

				s.currentBottomScreen = BottomScreenExtraWeather
			}

			if env.IsDebug() {
				devImgDraw(drawer)
			}
			pixoo64Draw(drawer)

			log.Debug("data draw finished")
			time.Sleep(1 * time.Minute)
		}
	}()
}

func devImgDraw(drawer *drawer.Drawer) {
	filename := "dev_img.png"
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal("failed to create file: ", zap.Error(err))
	}
	defer file.Close()

	png.Encode(file, drawer.Image())

	log.Debug("success draw dev image to " + filename)
}

func pixoo64Draw(drawer *drawer.Drawer) {
	client := pixoo64.NewClient(os.Getenv("PIXOO_ADDRESS"))
	var frames []frame.Frame
	frame, err := frame.NewFrameImage(drawer.Image(), 400)
	if err != nil {
		log.Fatal("failed to create frame: ", zap.Error(err))
	}
	frames = append(frames, *frame)

	pixoo64.ResetHttpGifId(client)
	err = pixoo64.SendHttpGif(client, 0, frames)
	if err != nil {
		log.Fatal("failed to send http gif: ", zap.Error(err))
	}

	/*time.Sleep(1 * time.Second)
	err = pixoo64.SendHttpText(client, 0, "hello world", image.Point{X: 4, Y: 50}, "#00ff00", 0)
	if err != nil {
		log.Fatal("failed to send http text: ", err)
	}*/

	log.Debug("success draw on pixoo64")
}
