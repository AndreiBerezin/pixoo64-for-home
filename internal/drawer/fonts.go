package drawer

import (
	"fmt"
	"log"
	"os"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

type FontFace string

const (
	FontMicro5Normal FontFace = "micro5_normal"
	FontMicro5Big    FontFace = "micro5_big"
	FontTiny5Normal  FontFace = "tiny5_normal"
	FontTiny5Big     FontFace = "tiny5_big"
)

var Fonts map[FontFace]font.Face

func init() {
	// 4 для pico8
	micro5Normal, err := newFont("static/fonts/micro5.ttf", 11)
	if err != nil {
		log.Fatal(err)
	}

	micro5Big, err := newFont("static/fonts/micro5.ttf", 22)
	if err != nil {
		log.Fatal(err)
	}

	tiny5Normal, err := newFont("static/fonts/tiny5.ttf", 8)
	if err != nil {
		log.Fatal(err)
	}

	tiny5Big, err := newFont("static/fonts/tiny5.ttf", 16)
	if err != nil {
		log.Fatal(err)
	}

	Fonts = map[FontFace]font.Face{
		FontMicro5Normal: micro5Normal,
		FontMicro5Big:    micro5Big,
		FontTiny5Normal:  tiny5Normal,
		FontTiny5Big:     tiny5Big,
	}
}

func newFont(filename string, size int) (font.Face, error) {
	fontBytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to load font: %w", err)
	}

	ttfFont, err := opentype.Parse(fontBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse font: %w", err)
	}

	return opentype.NewFace(ttfFont, &opentype.FaceOptions{
		Size:    float64(size),
		DPI:     72,
		Hinting: font.HintingNone,
	})
}
