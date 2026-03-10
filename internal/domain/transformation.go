package domain

type Transformation interface {
	Apply(x, y float64) (float64, float64)
	Weight() float64
}

type TransformationFactory interface {
	Create(name string, weight float64) (Transformation, error)
}
