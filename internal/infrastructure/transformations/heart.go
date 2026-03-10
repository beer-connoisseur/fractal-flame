package transformations

import "math"

type HeartTransformation struct {
	*BaseTransformation
}

func NewHeartTransformation(weight float64) *HeartTransformation {
	return &HeartTransformation{
		BaseTransformation: NewBaseTransformation("heart", weight),
	}
}

func (t *HeartTransformation) Apply(x, y float64) (float64, float64) {
	r := math.Sqrt(x*x + y*y)
	newX := r * math.Sin(r*math.Atan2(y, x))
	newY := -r * math.Cos(r*math.Atan2(y, x))

	return newX, newY
}
