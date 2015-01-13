package ocr

import (
	"encoding/base64"
	"image"
	"image/png"
	"strings"
)

const (
	IMG_0   string = "iVBORw0KGgoAAAANSUhEUgAAAAUAAAAIBAMAAADD3ygIAAAAIVBMVEWfn59fX19/f3/f398fHx8/Pz+/v78AAACc4/4AAP//gABOeO/lAAAAAXRSTlMAQObYZgAAAAFiS0dEAIgFHUgAAAAJcEhZcwAACxMAAAsTAQCanBgAAAAHdElNRQfeCBAROhivqETJAAAAHklEQVQI12MwLTVgEGISYHBmdmAoAEIQDeKDxIEAAEaEA+mAdl79AAAAAElFTkSuQmCC"
	IMG_1   string = "iVBORw0KGgoAAAANSUhEUgAAAAUAAAAIBAMAAADD3ygIAAAAIVBMVEWfn58fHx9/f39fX18/Pz/f39+/v78AAACc4/4AAP//gAAKcjd9AAAAAXRSTlMAQObYZgAAAAFiS0dEAIgFHUgAAAAJcEhZcwAACxMAAAsTAQCanBgAAAAHdElNRQfeCBASAAySnfCUAAAAFklEQVQI12NgK2BgKAdiBjRcDhZkAABKIgR8AVHj5wAAAABJRU5ErkJggg=="
	IMG_2   string = "iVBORw0KGgoAAAANSUhEUgAAAAUAAAAIBAMAAADD3ygIAAAAIVBMVEV/f3+fn58fHx9fX18/Pz/f39+/v78AAACc4/4AAP//gAAgV5JqAAAAAXRSTlMAQObYZgAAAAFiS0dEAIgFHUgAAAAJcEhZcwAACxMAAAsTAQCanBgAAAAHdElNRQfeCBASABiIRyTpAAAAJUlEQVQI12MQKUpgMGVUYGBgBeKQBAZWNwYG4QAGhvLyAgYgAABDOgQikbL7UQAAAABJRU5ErkJggg=="
	IMG_3   string = "iVBORw0KGgoAAAANSUhEUgAAAAUAAAAIBAMAAADD3ygIAAAAIVBMVEWc4/6fn5/f399fX1+/v78fHx8/Pz9/f38AAAAAAP//gAB1Sdq0AAAAAXRSTlMAQObYZgAAAAFiS0dEAIgFHUgAAAAJcEhZcwAACxMAAAsTAQCanBgAAAAHdElNRQfeCBASARyWMdGxAAAAJUlEQVQI12Moa3VgMGIJYGBgTGBgaFVgYACykxgDGErbHBiAAABZHQUKBdN1HwAAAABJRU5ErkJggg=="
	IMG_4   string = "iVBORw0KGgoAAAANSUhEUgAAAAUAAAAIBAMAAADD3ygIAAAAIVBMVEV/f39fX1+/v78/Pz8fHx/f39+fn58AAACc4/4AAP//gAAHF60HAAAAAXRSTlMAQObYZgAAAAFiS0dEAIgFHUgAAAAJcEhZcwAACxMAAAsTAQCanBgAAAAHdElNRQfeCBASAiSVHjrsAAAAI0lEQVQI12NgSGdgYHVnYGBWZ2AQYmdgKC8vYGAA0mDMwAAALP0CY6hYOk0AAAAASUVORK5CYII="
	IMG_5   string = "iVBORw0KGgoAAAANSUhEUgAAAAUAAAAIBAMAAADD3ygIAAAAIVBMVEVfX19/f3+/v78/Pz8fHx/f39+fn58AAACc4/4AAP//gABTSap+AAAAAXRSTlMAQObYZgAAAAFiS0dEAIgFHUgAAAAJcEhZcwAACxMAAAsTAQCanBgAAAAHdElNRQfeCBASAi3swoJIAAAAJElEQVQI12MoLy9gKGBgYCgvDmBgCHQAsgoYTBkNGERAfAYGAGiPBUfqt8YDAAAAAElFTkSuQmCC"
	IMG_6   string = "iVBORw0KGgoAAAANSUhEUgAAAAUAAAAIBAMAAADD3ygIAAAAIVBMVEVfX19/f3/f398fHx+fn58/Pz+/v78AAACc4/4AAP//gACpWu9/AAAAAXRSTlMAQObYZgAAAAFiS0dEAIgFHUgAAAAJcEhZcwAACxMAAAsTAQCanBgAAAAHdElNRQfeCBAROgcioEk8AAAAJUlEQVQI12NgMWdgcFVgYDBjYGAoLXZgKGIzYAgBYtXSBKAIAwBADwQrm36bpwAAAABJRU5ErkJggg=="
	IMG_7   string = "iVBORw0KGgoAAAANSUhEUgAAAAUAAAAIBAMAAADD3ygIAAAAIVBMVEVfX19/f3/f398fHx+fn58/Pz+/v78AAACc4/4AAP//gACpWu9/AAAAAXRSTlMAQObYZgAAAAFiS0dEAIgFHUgAAAAJcEhZcwAACxMAAAsTAQCanBgAAAAHdElNRQfeCBAROTi/6zfCAAAAI0lEQVQI12MoLy9gYGALYGBgTWBgSGVgYAgBcg2AXAcGEAAASVwDTm4cYkAAAAAASUVORK5CYII="
	IMG_8   string = "iVBORw0KGgoAAAANSUhEUgAAAAUAAAAIBAMAAADD3ygIAAAAIVBMVEWc4/5fX19/f3/f398fHx+fn58/Pz+/v78AAAAAAP//gAAyGfFPAAAAAXRSTlMAQObYZgAAAAFiS0dEAIgFHUgAAAAJcEhZcwAACxMAAAsTAQCanBgAAAAHdElNRQfeCBAROhHWdPxtAAAAJklEQVQI12MoaSlgcGd3YEgyF2AoSzFgcGZ2YHBldWAoaStgAAIAb9kGJPcvpx0AAAAASUVORK5CYII="
	IMG_9   string = "iVBORw0KGgoAAAANSUhEUgAAAAUAAAAIBAMAAADD3ygIAAAAIVBMVEVfX1+fn5/f39+/v78fHx8/Pz9/f38AAACc4/4AAP//gABIo3ygAAAAAXRSTlMAQObYZgAAAAFiS0dEAIgFHUgAAAAJcEhZcwAACxMAAAsTAQCanBgAAAAHdElNRQfeCBASASXJNFm5AAAAJUlEQVQI12MwLVVgcGYMYHBmKmAQKS1gYGAOYGBQFWBgd2QAAQBPVgPdUYkOjQAAAABJRU5ErkJggg=="
	IMG_DOT string = "iVBORw0KGgoAAAANSUhEUgAAAAEAAAAIBAMAAADKNIhyAAAAIVBMVEVfX19/f3/f398fHx+fn58/Pz+/v78AAACc4/4AAP//gACpWu9/AAAAAXRSTlMAQObYZgAAAAFiS0dEAIgFHUgAAAAJcEhZcwAACxMAAAsTAQCanBgAAAAHdElNRQfeCBAROi8XFeHGAAAAD0lEQVQI12NggIMCIGQAAAOQAOFBuM1DAAAAAElFTkSuQmCC"
)

var (
	DIGITS [10]*image.Image
	DOT    *image.Image
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
	if err != nil {
		panic("error loading embedded reference images")
	}
	printAsciiImage(&img)
	return &img
}
