package integrations

import (
	"encoding/json"
	"io"
	"os"

	"github.com/AndreiBerezin/pixoo64/internal/collector/types"
	"github.com/AndreiBerezin/pixoo64/pkg/log"
)

type Xras struct{}

func NewXras() *Xras {
	return &Xras{}
}

func (x *Xras) Data() (*types.MagneticData, error) {
	log.Info("Getting magnetic data...")

	/*var response *xrasResponse
	var err error
	if env.IsDebug() {
		response, err = x.mockResponse()
	} else {
		response, err = x.callApi()
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get xras response: %w", err)
	}*/

	return &types.MagneticData{}, nil
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
	//https://xras.ru/txt/kpf_RAL5.json

	return &xrasResponse{}, nil
}

type xrasResponse struct {
	Data []xrasData `json:"data"`
}

type xrasData struct {
	Time string `json:"time"`
	H00  string `json:"h00"`
	H03  string `json:"h03"`
	H06  string `json:"h06"`
	H09  string `json:"h09"`
	H12  string `json:"h12"`
	H15  string `json:"h15"`
	H18  string `json:"h18"`
	H21  string `json:"h21"`
}
