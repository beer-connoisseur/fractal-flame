package transformations

type BaseTransformation struct {
	name   string
	weight float64
}

func NewBaseTransformation(name string, weight float64) *BaseTransformation {
	return &BaseTransformation{
		name:   name,
		weight: weight,
	}
}

func (t *BaseTransformation) Weight() float64 {
	return t.weight
}
