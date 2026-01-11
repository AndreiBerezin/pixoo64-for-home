package pixoo64

import (
	"fmt"
	"os"

	"github.com/AndreiBerezin/pixoo64/pkg/http_client"
)

const (
	DeviceWidth  = 64
	DeviceHeight = 64
)

type Pixoo64 struct {
	client *http_client.Client
	addr   string
}

func NewPixoo64() *Pixoo64 {
	return &Pixoo64{
		client: http_client.New(),
		addr:   fmt.Sprintf("http://%s/post", os.Getenv("PIXOO_ADDRESS")),
	}
}

func (p *Pixoo64) callApi(bodyObject any) error {
	return p.client.Post(p.addr, bodyObject, nil)
}

func (p *Pixoo64) callApiWithResponse(bodyObject any, target any) error {
	return p.client.Post(p.addr, bodyObject, target)
}
