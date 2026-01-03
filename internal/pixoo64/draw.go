package pixoo64

import (
	"encoding/base64"
	"encoding/json"
	"image"

	"github.com/AndreiBerezin/pixoo64/internal/frame"
)

func ResetHttpGifId(client *Client) error {
	data := map[string]any{
		"Command": "Draw/ResetHttpGifId",
	}

	_, err := client.Post(data)
	if err != nil {
		return err
	}
	return err
}

type GetHttpGifIdResponse struct {
	ErrorCode int `json:"error_code"`
	PicID     int `json:"PicId"`
}

func GetHttpGifId(client *Client) (int, error) {
	data := map[string]any{
		"Command": "Draw/GetHttpGifId",
	}
	body, err := client.Post(data)
	if err != nil {
		return 0, err
	}

	result := GetHttpGifIdResponse{}
	err = json.Unmarshal([]byte(body), &result)
	if err != nil {
		return 0, err
	}

	return result.PicID, nil
}

func SendHttpGif(client *Client, httpGifID int, frames []frame.Frame) error {
	for i, f := range frames {
		data := map[string]any{
			"Command":   "Draw/SendHttpGif",
			"PicNum":    len(frames),
			"PicWidth":  DeviceWidth,
			"PicOffset": i,
			"PicID":     httpGifID,
			"PicSpeed":  f.Speed(),
			"PicData":   base64.StdEncoding.EncodeToString(f.ToBytes()),
		}
		_, err := client.Post(data)
		if err != nil {
			return err
		}
	}

	return nil
}

func SendHttpText(client *Client, textID int, text string, point image.Point, colorHex string, font int) error {
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

	_, err := client.Post(data)
	if err != nil {
		return err
	}
	return nil
}
