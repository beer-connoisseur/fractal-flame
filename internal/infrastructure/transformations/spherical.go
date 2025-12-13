package transformations

type SphericalTransformation struct {
	*BaseTransformation
}

func NewSphericalTransformation(weight float64) *SphericalTransformation {
	return &SphericalTransformation{
		BaseTransformation: NewBaseTransformation("spherical", weight),
	}
}

func (t *SphericalTransformation) Apply(x, y float64) (float64, float64) {
	r2 := x*x + y*y
	if r2 == 0 {
		return 0, 0
	}
	newX := x / r2
	newY := y / r2

	return newX, newY
}
