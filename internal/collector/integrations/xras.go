package integrations

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/AndreiBerezin/pixoo64/internal/collector/types"
	"github.com/AndreiBerezin/pixoo64/pkg/env"
	"github.com/AndreiBerezin/pixoo64/pkg/log"
)

type Xras struct{}

func NewXras() *Xras {
	return &Xras{}
}

func (x *Xras) Data() (*types.MagneticData, error) {
	log.Info("Getting magnetic data...")

	var response *xrasResponse
	var err error
	if env.IsDebug() {
		response, err = x.mockResponse()
	} else {
		response, err = x.callApi()
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get xras response: %w", err)
	}

	slices.SortFunc(response.Data, func(a, b xrasData) int {
		return strings.Compare(a.Time, b.Time)
	})

	var days []types.MagneticDay
	for _, day := range response.Data {
		dayTime, err := time.Parse("2006-01-02", day.Time)
		if err != nil {
			return nil, fmt.Errorf("failed to parse day time: %w", err)
		}
		days = append(days, types.MagneticDay{
			Day: dayTime.Format("02"),
			Hours: []types.MagneticHour{
				{
					Hour:  0,
					Level: day.H00.Float32(),
				},
				{
					Hour:  3,
					Level: day.H03.Float32(),
				},
				{
					Hour:  6,
					Level: day.H06.Float32(),
				},
				{
					Hour:  9,
					Level: day.H09.Float32(),
				},
				{
					Hour:  12,
					Level: day.H12.Float32(),
				},
				{
					Hour:  15,
					Level: day.H15.Float32(),
				},
				{
					Hour:  18,
					Level: day.H18.Float32(),
				},
				{
					Hour:  21,
					Level: day.H21.Float32(),
				},
			},
		})
	}

	return &types.MagneticData{
		Days: days,
	}, nil
}

func (x *Xras) mockResponse() (*xrasResponse, error) {
	file, err := os.Open("mocks/xras.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var response *xrasResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (x *Xras) callApi() (*xrasResponse, error) {
	req, err := http.NewRequest("GET", "https://xras.ru/txt/kpf_RAL5.json", nil)
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

	var response *xrasResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

type xrasResponse struct {
	Data []xrasData `json:"data"`
}

type xrasData struct {
	Time string `json:"time"`
	H00  level  `json:"h00"`
	H03  level  `json:"h03"`
	H06  level  `json:"h06"`
	H09  level  `json:"h09"`
	H12  level  `json:"h12"`
	H15  level  `json:"h15"`
	H18  level  `json:"h18"`
	H21  level  `json:"h21"`
}

type level string

func (l level) Float32() float32 {
	value, err := strconv.ParseFloat(string(l), 32)
	if err != nil {
		return 0
	}
	return float32(value)
}
