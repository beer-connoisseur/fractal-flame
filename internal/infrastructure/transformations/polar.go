package transformations

import "math"

type PolarTransformation struct {
	*BaseTransformation
}

func NewPolarTransformation(weight float64) *PolarTransformation {
	return &PolarTransformation{
		BaseTransformation: NewBaseTransformation("polar", weight),
	}
}

func (t *PolarTransformation) Apply(x, y float64) (float64, float64) {
	newX := math.Atan2(y, x) / math.Pi
	newY := math.Sqrt(x*x+y*y) - 1

	return newX, newY
}
