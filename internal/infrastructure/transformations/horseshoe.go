package transformations

import "math"

type HorseshoeTransformation struct {
	*BaseTransformation
}

func NewHorseshoeTransformation(weight float64) *HorseshoeTransformation {
	return &HorseshoeTransformation{
		BaseTransformation: NewBaseTransformation("horseshoe", weight),
	}
}

func (t *HorseshoeTransformation) Apply(x, y float64) (float64, float64) {
	r := math.Sqrt(x*x + y*y)
	if r == 0 {
		return 0, 0
	}
	newX := (x - y) * (x + y) / r
	newY := (2 * x * y) / r

	return newX, newY
}
