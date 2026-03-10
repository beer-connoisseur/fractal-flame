package transformations

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw4-fractal-flame/internal/domain"
)

type MockTransformation struct {
	name   string
	weight float64
}

func NewMockTransformation(name string, weight float64) *MockTransformation {
	return &MockTransformation{
		name:   name,
		weight: weight,
	}
}

func (m *MockTransformation) Weight() float64 {
	return m.weight
}

func (m *MockTransformation) Apply(x, y float64) (float64, float64) {
	return x * m.weight, y * m.weight
}

func TestNewTransformationFactory(t *testing.T) {
	t.Run("creates factory with default transformations", func(t *testing.T) {
		factory := NewTransformationFactory()

		require.NotNil(t, factory)
		require.NotNil(t, factory.registry)
		require.Greater(t, len(factory.registry), 0)

		expectedTransformations := []string{
			"swirl", "horseshoe", "polar", "spherical", "heart",
		}

		for _, name := range expectedTransformations {
			_, exists := factory.registry[name]
			require.True(t, exists, "Transformation %s should be registered", name)
		}
	})
}

func TestTransformationFactoryImpl_Create(t *testing.T) {
	t.Run("create existing transformation", func(t *testing.T) {
		factory := NewTransformationFactory()

		testCases := []struct {
			name   string
			weight float64
		}{
			{"swirl", 1.0},
			{"horseshoe", 0.5},
			{"polar", 2.0},
			{"spherical", 0.3},
			{"heart", 1.5},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				transformation, err := factory.Create(tc.name, tc.weight)
				require.NoError(t, err)
				require.NotNil(t, transformation)

				_, ok := transformation.(domain.Transformation)
				require.True(t, ok)

				require.Equal(t, tc.weight, transformation.Weight())

				x, y := transformation.Apply(1.0, 1.0)
				require.False(t, math.IsNaN(x))
				require.False(t, math.IsNaN(y))
			})
		}
	})

	t.Run("create non-existent transformation", func(t *testing.T) {
		factory := NewTransformationFactory()

		transformation, err := factory.Create("non_existent", 1.0)
		require.Error(t, err)
		require.Nil(t, transformation)
		require.Contains(t, err.Error(), "not found")
	})
}

func TestTransformationFactoryImpl_Register(t *testing.T) {
	t.Run("register new transformation", func(t *testing.T) {
		factory := &TransformationFactoryImpl{
			registry: make(map[string]func(weight float64) domain.Transformation),
		}

		factory.Register("custom", func(weight float64) domain.Transformation {
			return NewMockTransformation("custom", weight)
		})

		_, exists := factory.registry["custom"]
		require.True(t, exists)
	})
}
