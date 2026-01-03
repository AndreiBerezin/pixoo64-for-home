package main

import (
	"image/png"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AndreiBerezin/pixoo64/internal/collector"
	"github.com/AndreiBerezin/pixoo64/internal/drawer"
	"github.com/AndreiBerezin/pixoo64/internal/frame"
	"github.com/AndreiBerezin/pixoo64/internal/pixoo64"
	"github.com/AndreiBerezin/pixoo64/internal/screens"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file", err)
	}

	collector := collector.NewCollector()
	collector.Start()

	time.Sleep(2 * time.Second)
	go func() {
		for {
			data := collector.GetCollectedData()

			drawer, err := drawer.NewDrawer()
			if err != nil {
				log.Fatal(err)
			}

			weatherScreen, err := screens.NewCurrentWeatherScreen(drawer)
			if err != nil {
				log.Fatal(err)
			}
			weatherScreen.DrawStatic(data.YandexData)
			if err != nil {
				log.Fatal(err)
			}

			extraWeatherScreen, err := screens.NewExtraWeatherScreen(drawer)
			if err != nil {
				log.Fatal(err)
			}
			extraWeatherScreen.DrawTodayStatic(data.YandexData)
			if err != nil {
				log.Fatal(err)
			}

			draw(drawer)
			pixoo64Draw(drawer)

			log.Print("draw data")
			time.Sleep(1 * time.Minute)
		}
	}()

	waitShutdownSignal()
}

func draw(drawer *drawer.Drawer) {
	filename := "test.png"
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	png.Encode(file, drawer.Image())

	log.Print("Image saved to ", filename)
}

func pixoo64Draw(drawer *drawer.Drawer) {
	client := pixoo64.NewClient(os.Getenv("PIXOO_ADDRESS"))
	var frames []frame.Frame
	frame, err := frame.NewFrameImage(drawer.Image(), 400)
	if err != nil {
		log.Fatal("failed to create frame", err)
	}
	frames = append(frames, *frame)

	pixoo64.ResetHttpGifId(client)
	err = pixoo64.SendHttpGif(client, 0, frames)
	if err != nil {
		log.Fatal("failed to send http gif", err)
	}

	/*time.Sleep(1 * time.Second)
	err = pixoo64.SendHttpText(client, 0, "hello world", image.Point{X: 4, Y: 50}, "#00ff00", 0)
	if err != nil {
		log.Fatal("failed to send http text", err)
	}*/
}

func waitShutdownSignal() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigChan
	log.Print("received shutdown signal", sig)
}
