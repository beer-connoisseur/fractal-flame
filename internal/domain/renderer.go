package domain

import (
	"image"
)

type Renderer interface {
	Render(histogram Histogram) image.Image
	Save(image image.Image, path string) error
}

type ColorMapper interface {
	GeneratePalette(size int, seed int64) []Color
}

type SymmetryApplier interface {
	Apply(x, y float64, level int) []Point
}
