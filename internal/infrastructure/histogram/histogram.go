package histogram

import (
	"errors"
	"math"

	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw4-fractal-flame/internal/domain"
)

type HistField struct {
	width   int
	height  int
	pixels  [][]domain.Pixel
	maxHits int
}

func NewHistogram(width, height int) *HistField {
	pixels := make([][]domain.Pixel, height)
	for y := range pixels {
		pixels[y] = make([]domain.Pixel, width)
	}

	return &HistField{
		width:  width,
		height: height,
		pixels: pixels,
	}
}

func (h *HistField) Add(x, y float64, color domain.Color) {
	xNorm := (domain.XMAX - float64(x)) / (domain.XMAX - domain.XMIN)
	yNorm := (domain.YMAX - float64(y)) / (domain.YMAX - domain.YMIN)
	px := h.width - int(xNorm*float64(h.width))
	py := h.height - int(yNorm*float64(h.height))

	if px < 0 || px >= h.width || py < 0 || py >= h.height {
		return
	}

	pixel := &h.pixels[py][px]

	if pixel.Counter == 0 {
		pixel.Color.R = color.R
		pixel.Color.G = color.G
		pixel.Color.B = color.B
	} else {
		pixel.Color.R = (pixel.Color.R + color.R) / 2
		pixel.Color.G = (pixel.Color.G + color.G) / 2
		pixel.Color.B = (pixel.Color.B + color.B) / 2
	}

	pixel.Counter++

	if pixel.Counter > h.maxHits {
		h.maxHits = pixel.Counter
	}

	return
}

func (h *HistField) GetColor(x, y int) domain.Color {
	if x < 0 || x >= h.width || y < 0 || y >= h.height {
		return domain.Color{}
	}

	pixel := h.pixels[y][x]
	return domain.Color{
		R: pixel.Color.R,
		G: pixel.Color.G,
		B: pixel.Color.B,
	}
}

func (h *HistField) ApplyCorrection(gamma float64) error {
	if h.maxHits == 0 {
		return errors.New("there are no pixels in the histogram")
	}

	logMaxHits := math.Log10(float64(h.maxHits))

	for y := 0; y < h.height; y++ {
		for x := 0; x < h.width; x++ {
			pixel := &h.pixels[y][x]

			if pixel.Counter > 0 {
				normal := math.Log10(float64(pixel.Counter))
				var normalized float64
				if h.maxHits == 1 {
					normalized = 1.0
				} else {
					normalized = normal / logMaxHits
				}

				gammaCorrected := math.Pow(normalized, 1.0/gamma)
				pixel.Color.R *= gammaCorrected
				pixel.Color.G *= gammaCorrected
				pixel.Color.B *= gammaCorrected
			}
		}
	}

	return nil
}

func (h *HistField) Merge(other domain.Histogram) error {
	otherHist, ok := other.(*HistField)
	if !ok {
		return errors.New("can't merge histogram with wrong histogram type")
	}

	if h.width != otherHist.width || h.height != otherHist.height {
		return errors.New("can't merge histogram with histogram of other sizes")
	}

	for y := 0; y < h.height; y++ {
		for x := 0; x < h.width; x++ {
			otherPixel := otherHist.pixels[y][x]

			if otherPixel.Counter > 0 {
				currentPixel := &h.pixels[y][x]

				totalCounter := currentPixel.Counter + otherPixel.Counter

				weight1 := float64(currentPixel.Counter) / float64(totalCounter)
				weight2 := float64(otherPixel.Counter) / float64(totalCounter)

				currentPixel.Color.R = currentPixel.Color.R*weight1 + otherPixel.Color.R*weight2
				currentPixel.Color.G = currentPixel.Color.G*weight1 + otherPixel.Color.G*weight2
				currentPixel.Color.B = currentPixel.Color.B*weight1 + otherPixel.Color.B*weight2

				currentPixel.Counter = totalCounter
				if currentPixel.Counter > h.maxHits {
					h.maxHits = currentPixel.Counter
				}
			}
		}
	}

	return nil
}

func (h *HistField) GetSize() (width, height int) {
	return h.width, h.height
}

type HistFactory struct{}

func NewHistogramFactory() *HistFactory {
	return &HistFactory{}
}

func (f *HistFactory) Create(width, height int) domain.Histogram {
	return NewHistogram(width, height)
}
