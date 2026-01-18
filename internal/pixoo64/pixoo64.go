package pixoo64

import (
	"fmt"
	"image"
	"os"

	"github.com/AndreiBerezin/pixoo64/internal/pixoo64/frame"
	"github.com/AndreiBerezin/pixoo64/pkg/http_client"
)

type Pixoo64 struct {
	client *http_client.Client
	addr   string
	width  int
	height int
}

func New(width int, height int) *Pixoo64 {
	return &Pixoo64{
		client: http_client.New(),
		addr:   fmt.Sprintf("http://%s/post", os.Getenv("PIXOO_ADDRESS")),
		width:  width,
		height: height,
	}
}

func (p *Pixoo64) callApi(bodyObject any) error {
	return p.client.Post(p.addr, bodyObject, nil)
}

func (p *Pixoo64) callApiWithResponse(bodyObject any, target any) error {
	return p.client.Post(p.addr, bodyObject, target)
}

func (p *Pixoo64) DrawImage(image *image.RGBA) error {
	frames := []frame.Frame{frame.New(image, 400)}

	p.ResetHttpGifId()
	if err := p.SendHttpGif(0, frames); err != nil {
		return fmt.Errorf("failed to send http gif: %w", err)
	}

	return nil
}
