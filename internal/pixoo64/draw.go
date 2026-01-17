package pixoo64

import (
	"encoding/base64"
	"image"

	"github.com/AndreiBerezin/pixoo64/internal/state/frame"
)

func (p *Pixoo64) ResetHttpGifId() error {
	data := map[string]any{
		"Command": "Draw/ResetHttpGifId",
	}

	return p.callApi(data)
}

func (p *Pixoo64) GetHttpGifId() (int, error) {
	data := map[string]any{
		"Command": "Draw/GetHttpGifId",
	}
	var result struct {
		ErrorCode int `json:"error_code"`
		PicID     int `json:"PicId"`
	}
	err := p.callApiWithResponse(data, &result)
	if err != nil {
		return 0, err
	}

	return result.PicID, nil
}

func (p *Pixoo64) SendHttpGif(httpGifID int, frames []frame.Frame) error {
	for i, f := range frames {
		data := map[string]any{
			"Command":   "Draw/SendHttpGif",
			"PicNum":    len(frames),
			"PicWidth":  p.width,
			"PicOffset": i,
			"PicID":     httpGifID,
			"PicSpeed":  f.Speed(),
			"PicData":   base64.StdEncoding.EncodeToString(f.ToBytes()),
		}
		if err := p.callApi(data); err != nil {
			return err
		}
	}

	return nil
}

func (p *Pixoo64) SendHttpText(textID int, text string, point image.Point, colorHex string, font int) error {
	data := map[string]any{
		"Command":    "Draw/SendHttpText",
		"LcdId":      0,
		"TextID":     textID,
		"x":          point.X,
		"y":          point.Y,
		"dir":        0,
		"font":       font,
		"TextWidth":  64,
		"TextString": text,
		"speed":      0,
		"color":      colorHex,
		"align":      1,
	}

	return p.callApi(data)
}
