package ocr

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"net/http"
	"strconv"
)

func GetImage(url string) (image.Image, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	img, err := png.Decode(resp.Body)
	if err != nil {
		return nil, err
	}
	return img, err
}

func width(img *image.Image) int {
	return (*img).Bounds().Max.X - (*img).Bounds().Min.X
}

func height(img *image.Image) int {
	return (*img).Bounds().Max.Y - (*img).Bounds().Min.Y
}

func ReadValue(img *image.Image) (int, error) {
	if width(img) != 70 {
		panic("image width is not 70")
	}
	if height(img) != 8 {
		panic("image height is not 8")
	}
	digits := getDigitsImages(img)
	var buffer bytes.Buffer
	for _, digit := range digits {
		buffer.WriteString(getDigit(digit))
	}
	return strconv.Atoi(buffer.String())
}

func getEmptyColumns(img *image.Image) map[int]bool {
	emptyCols := map[int]bool{}
col_loop:
	for i := (*img).Bounds().Min.X; i < (*img).Bounds().Max.X; i++ {
		for j := (*img).Bounds().Min.Y; j < (*img).Bounds().Max.Y; j++ {
			_, _, _, alpha := (*img).At(i, j).RGBA()
			if alpha > 0 {
				continue col_loop
			}
		}
		emptyCols[i] = true
	}
	return emptyCols
}

func getDigitsBounds(img *image.Image) [][]int {
	emptyCols := getEmptyColumns(img)
	digitsBounds := [][]int{}
	start := -1
	end := -1
	for i := (*img).Bounds().Min.X; i < (*img).Bounds().Max.X; i++ {
		if emptyCols[i] {
			if start > 0 && end > 0 {
				digitsBounds = append(digitsBounds, []int{start, end})
			}
			start = -1
			end = -1
			continue
		}
		if start == -1 {
			start = i
		}
		end = i
	}
	return digitsBounds
}

func getDigitsImages(img *image.Image) []*image.Image {
	digitsBounds := getDigitsBounds(img)
	digitsImages := []*image.Image{}
	for _, bounds := range digitsBounds {
		digitImg := (*img).(*image.RGBA).SubImage(image.Rect(bounds[0], bounds[1], 0, (*img).Bounds().Max.Y))
		digitsImages = append(digitsImages, &digitImg)
	}
	return digitsImages
}

func areSimilarPixels(recoPixel, refPixel color.RGBA) bool {
	if recoPixel.A == refPixel.A {
		return true
	}
	return false
}

func areSimilarImages(recoDigit, refDigit *image.Image) bool {
	// check dimensions
	if width(recoDigit) != width(refDigit) {
		return false
	}
	if height(recoDigit) != height(refDigit) {
		return false
	}
	// check pixels
	for i := (*recoDigit).Bounds().Min.X; i < (*recoDigit).Bounds().Max.X; i++ {
		for j := 0; j < (*recoDigit).Bounds().Max.Y; j++ {
			_, _, _, recoPixelAlpha := (*recoDigit).At(i, j).RGBA()
			_, _, _, refPixelAlpha := (*recoDigit).At(i, j).RGBA()
			if recoPixelAlpha != refPixelAlpha {
				return false
			}
		}
	}
	return true
}

func getDigit(digitImg *image.Image) string {
	for i := 0; i < len(DIGITS); i++ {
		if areSimilarImages(digitImg, DIGITS[i]) {
			return strconv.Itoa(i)
		}
	}
	if areSimilarImages(digitImg, DOT) {
		return ""
	}
	return ""
}
