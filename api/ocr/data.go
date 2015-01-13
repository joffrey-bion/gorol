package ocr

import (
	"encoding/base64"
	"image"
	"image/png"
	"strings"
)

const (
	IMG_0 string = ""
	IMG_1 string = ""
	IMG_2 string = ""
	IMG_3 string = ""
	IMG_4 string = ""
	IMG_5 string = ""
	IMG_6 string = ""
	IMG_7 string = ""
	IMG_8 string = ""
	IMG_9 string = ""
	IMG_DOT string = ""
)

var (
	DIGITS []*image.Image
	DOT *image.Image
)

func init() {
	DIGITS[0] = loadPNG(IMG_0)
	DIGITS[1] = loadPNG(IMG_1)
	DIGITS[2] = loadPNG(IMG_2)
	DIGITS[3] = loadPNG(IMG_3)
	DIGITS[4] = loadPNG(IMG_4)
	DIGITS[5] = loadPNG(IMG_5)
	DIGITS[6] = loadPNG(IMG_6)
	DIGITS[7] = loadPNG(IMG_7)
	DIGITS[8] = loadPNG(IMG_8)
	DIGITS[9] = loadPNG(IMG_9)
	DOT = loadPNG(IMG_DOT)
}

func loadPNG(b64 string) *image.Image {
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(b64))
	img, err := png.Decode(reader)
	if (err != nil) {
		panic("error loading embedded reference images")
	}
	return &img
}
