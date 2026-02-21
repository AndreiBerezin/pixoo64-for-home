package integrations

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"
	"time"

	"github.com/AndreiBerezin/pixoo64/internal/collector/types"
	"github.com/AndreiBerezin/pixoo64/pkg/env"
	"github.com/AndreiBerezin/pixoo64/pkg/log"
)

type OpenMeteo struct{}

func NewOpenMeteo() *OpenMeteo {
	return &OpenMeteo{}
}

func (o *OpenMeteo) Data() (*types.PressureData, error) {
	log.Info("Getting pressure data...")

	var response *openMeteoResponse
	var err error
	if env.IsDebug() {
		response, err = o.mockResponse()
	} else {
		response, err = o.callApi()
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get openmeteo response: %w", err)
	}

	collectHours := []int{0, 3, 6, 9, 12, 15, 18, 21}

	var minDate time.Time
	hoursByDay := make(map[string][]types.PressureHour)
	for i, t := range response.Hourly.Time {
		parsedTime, err := time.Parse("2006-01-02T15:04", t)
		if err != nil {
			return nil, fmt.Errorf("failed to parse time: %w", err)
		}

		date := parsedTime.Format("2006-01-02")
		hour := parsedTime.Hour()
		if !slices.Contains(collectHours, hour) {
			continue
		}

		hoursByDay[date] = append(hoursByDay[date], types.PressureHour{
			Hour:     hour,
			Pressure: float32(response.Hourly.SurfacePressure[i]) * 0.7500638, // hPa -> mmHg
		})

		if minDate.IsZero() || parsedTime.Before(minDate) {
			minDate = parsedTime
		}
	}

	var days []types.PressureDay
	for i := 0; i < 3; i++ {
		date := minDate.AddDate(0, 0, i)
		dateStr := date.Format("2006-01-02")
		if _, ok := hoursByDay[dateStr]; !ok {
			return nil, fmt.Errorf("no pressure data for day: %s", dateStr)
		}
		dayHours := hoursByDay[dateStr]

		days = append(days, types.PressureDay{
			Day:   date.Format("02"),
			Hours: dayHours,
		})
	}

	return &types.PressureData{
		Days: days,
	}, nil
}

func (o *OpenMeteo) mockResponse() (*openMeteoResponse, error) {
	file, err := os.Open("mocks/openmeteo.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var response *openMeteoResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (o *OpenMeteo) callApi() (*openMeteoResponse, error) {
	url := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%s&longitude=%s&hourly=surface_pressure&forecast_days=3&timezone=auto",
		os.Getenv("LAT"), os.Getenv("LON"),
	)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response *openMeteoResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

type openMeteoResponse struct {
	Hourly hourlyData `json:"hourly"`
}

type hourlyData struct {
	Time            []string  `json:"time"`
	SurfacePressure []float64 `json:"surface_pressure"`
}
