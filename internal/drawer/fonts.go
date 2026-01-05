package drawer

import (
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
	Fonts = map[FontFace]font.Face{
		FontMicro5Normal: newFont("static/fonts/micro5.ttf", 11),
		FontMicro5Big:    newFont("static/fonts/micro5.ttf", 22),
		FontTiny5Normal:  newFont("static/fonts/tiny5.ttf", 8),
		FontTiny5Big:     newFont("static/fonts/tiny5.ttf", 16),
	}
}

func newFont(filename string, size int) font.Face {
	fontBytes, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal("failed to load font: %w", err)
		return nil
	}

	ttfFont, err := opentype.Parse(fontBytes)
	if err != nil {
		log.Fatal("failed to parse font: %w", err)
		return nil
	}

	face, err := opentype.NewFace(ttfFont, &opentype.FaceOptions{
		Size:    float64(size),
		DPI:     72,
		Hinting: font.HintingNone,
	})
	if err != nil {
		log.Fatal("failed to create face: %w", err)
		return nil
	}

	return face
}
