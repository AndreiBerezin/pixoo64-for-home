package state

import (
	"sync/atomic"
	"time"

	"github.com/AndreiBerezin/pixoo64/internal/collector"
	"github.com/AndreiBerezin/pixoo64/internal/collector/types"
	"github.com/AndreiBerezin/pixoo64/internal/pixoo64"
	"github.com/AndreiBerezin/pixoo64/internal/screens"
	"github.com/AndreiBerezin/pixoo64/internal/timer"
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
	customTopScreens     []drawFn
	customTopIdx         int
	defaultBottomScreens []drawFn
	defaultBottomIdx     int
}

func (r *rotation) nextCustomTop() {
	r.customTopIdx = (r.customTopIdx + 1) % len(r.customTopScreens)
}

func (r *rotation) nextDefaultBottom() {
	r.defaultBottomIdx = (r.defaultBottomIdx + 1) % len(r.defaultBottomScreens)
}

type State struct {
	device       *pixoo64.Pixoo64
	collector    *collector.Collector
	screens      *screens.Screens
	timerManager *timer.Manager
	rotation     rotation

	onAir          atomic.Bool
	onAirStartTime time.Time
}

func (s *State) SetOnAir(on bool) {
	s.onAir.Store(on)
	if on {
		s.onAirStartTime = time.Now()
	}
}

func New(collector *collector.Collector, timerManager *timer.Manager) *State {
	sc := screens.New(deviceWidth, deviceHeight)
	return &State{
		device:       pixoo64.New(deviceWidth, deviceHeight),
		collector:    collector,
		screens:      sc,
		timerManager: timerManager,
		rotation: rotation{
			customTopScreens: []drawFn{
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
