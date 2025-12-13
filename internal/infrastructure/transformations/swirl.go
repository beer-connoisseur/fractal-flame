package transformations

import "math"

type SwirlTransformation struct {
	*BaseTransformation
}

func NewSwirlTransformation(weight float64) *SwirlTransformation {
	return &SwirlTransformation{
		BaseTransformation: NewBaseTransformation("swirl", weight),
	}
}

func (t *SwirlTransformation) Apply(x, y float64) (float64, float64) {
	r2 := x*x + y*y
	sinR2 := math.Sin(r2)
	cosR2 := math.Cos(r2)
	newX := x*sinR2 - y*cosR2
	newY := x*cosR2 + y*sinR2

	return newX, newY
}
