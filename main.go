package main

import (
	"image/png"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AndreiBerezin/pixoo64/internal/collector"
	"github.com/AndreiBerezin/pixoo64/internal/drawer"
	"github.com/AndreiBerezin/pixoo64/internal/frame"
	"github.com/AndreiBerezin/pixoo64/internal/pixoo64"
	"github.com/AndreiBerezin/pixoo64/internal/screens"
	"github.com/AndreiBerezin/pixoo64/pkg/env"
	"github.com/AndreiBerezin/pixoo64/pkg/log"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	_ = godotenv.Load()

	log.Init()
	defer log.Sync()

	pixooAddress := os.Getenv("PIXOO_ADDRESS")
	if pixooAddress == "" {
		log.Fatal("PIXOO_ADDRESS is not set, please check your environment variables")
	}

	log.Info("app started")

	collector := collector.NewCollector()
	collector.Start()

	time.Sleep(2 * time.Second)
	go func() {
		for {
			data := collector.GetCollectedData()

			drawer := drawer.NewDrawer()

			weatherScreen := screens.NewCurrentWeatherScreen(drawer)

			err := weatherScreen.DrawStatic(data.YandexData)
			if err != nil {
				log.Fatal("failed to draw weather screen: ", zap.Error(err))
			}

			/*extraWeatherScreen, err := screens.NewExtraWeatherScreen(drawer)
			if err != nil {
				log.Fatal("failed to create extra weather screen: ", zap.Error(err))
			}
			err = extraWeatherScreen.DrawTodayStatic(data.YandexData)
			if err != nil {
				log.Fatal("failed to draw extra weather screen: ", zap.Error(err))
			}*/

			magneticScreen, err := screens.NewMagneticScreen(drawer)
			if err != nil {
				log.Fatal("failed to create magnetic screen: ", zap.Error(err))
			}
			err = magneticScreen.DrawStatic(data.MagneticData, data.YandexData)
			if err != nil {
				log.Fatal("failed to draw magnetic screen: ", zap.Error(err))
			}

			if env.IsDebug() {
				devImgDraw(drawer)
			}
			pixoo64Draw(drawer)

			log.Debug("data draw finished")
			time.Sleep(1 * time.Minute)
		}
	}()

	waitShutdownSignal()
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

func waitShutdownSignal() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigChan
	log.Info("received shutdown signal: " + sig.String())
}
