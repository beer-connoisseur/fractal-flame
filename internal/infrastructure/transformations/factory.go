package transformations

import (
	"fmt"

	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw4-fractal-flame/internal/domain"
)

type TransformationFactoryImpl struct {
	registry map[string]func(weight float64) domain.Transformation
}

func NewTransformationFactory() *TransformationFactoryImpl {
	factory := &TransformationFactoryImpl{
		registry: make(map[string]func(weight float64) domain.Transformation),
	}

	factory.Register("swirl", func(weight float64) domain.Transformation {
		return NewSwirlTransformation(weight)
	})

	factory.Register("horseshoe", func(weight float64) domain.Transformation {
		return NewHorseshoeTransformation(weight)
	})

	factory.Register("polar", func(weight float64) domain.Transformation {
		return NewPolarTransformation(weight)
	})

	factory.Register("spherical", func(weight float64) domain.Transformation {
		return NewSphericalTransformation(weight)
	})

	factory.Register("heart", func(weight float64) domain.Transformation {
		return NewHeartTransformation(weight)
	})

	return factory
}

func (f *TransformationFactoryImpl) Register(
	name string,
	constructor func(weight float64) domain.Transformation,
) {
	f.registry[name] = constructor
}

func (f *TransformationFactoryImpl) Create(name string, weight float64) (domain.Transformation, error) {
	constructor, exists := f.registry[name]
	if !exists {
		return nil, fmt.Errorf("transformation %s not found", name)
	}

	return constructor(weight), nil
}
