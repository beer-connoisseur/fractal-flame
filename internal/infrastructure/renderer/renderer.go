package renderer

import (
	"image"
	"image/color"
	"image/png"
	"log/slog"
	"os"

	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw4-fractal-flame/internal/domain"
)

type PNGRenderer struct{}

func NewPNGRenderer() *PNGRenderer {
	return &PNGRenderer{}
}

func (re *PNGRenderer) Render(histogram domain.Histogram) image.Image {
	width, height := histogram.GetSize()
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			histColor := histogram.GetColor(x, y)

			r8 := uint8(histColor.R * 255.0)
			g8 := uint8(histColor.G * 255.0)
			b8 := uint8(histColor.B * 255.0)

			img.Set(x, y, color.RGBA{R: r8, G: g8, B: b8, A: 255})
		}
	}

	return img
}

func (re *PNGRenderer) Save(img image.Image, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			slog.Warn(err.Error())
		}
	}(file)

	encoder := png.Encoder{
		CompressionLevel: png.BestCompression,
	}

	if err = encoder.Encode(file, img); err != nil {
		return err
	}

	return nil
}
