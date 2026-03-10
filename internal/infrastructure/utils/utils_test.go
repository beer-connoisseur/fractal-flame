package utils

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw4-fractal-flame/internal/domain"
)

type MockTransformation struct {
	weight float64
}

func (m *MockTransformation) Weight() float64 {
	return m.weight
}

func (m *MockTransformation) Apply(x, y float64) (float64, float64) {
	return x * m.weight, y * m.weight
}

func TestCalculateTotalWeight(t *testing.T) {
	tests := []struct {
		name            string
		transformations []domain.Transformation
		want            float64
	}{
		{
			name:            "empty slice",
			transformations: []domain.Transformation{},
			want:            0.0,
		},
		{
			name: "single transformation",
			transformations: []domain.Transformation{
				&MockTransformation{weight: 2.5},
			},
			want: 2.5,
		},
		{
			name: "multiple transformations",
			transformations: []domain.Transformation{
				&MockTransformation{weight: 1.0},
				&MockTransformation{weight: 2.0},
				&MockTransformation{weight: 3.0},
			},
			want: 6.0,
		},
		{
			name: "with zero weight",
			transformations: []domain.Transformation{
				&MockTransformation{weight: 1.5},
				&MockTransformation{weight: 0.0},
				&MockTransformation{weight: 2.5},
			},
			want: 4.0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateTotalWeight(tt.transformations)
			assert.InDelta(t, tt.want, got, 0.001)
		})
	}
}

func TestSelectTransformation(t *testing.T) {
	transformation1 := &MockTransformation{weight: 1.0}
	transformation2 := &MockTransformation{weight: 2.0}
	transformation3 := &MockTransformation{weight: 3.0}

	tests := []struct {
		name            string
		r               *rand.Rand
		transformations []domain.Transformation
		totalWeight     float64
		want            domain.Transformation
		wantErr         bool
	}{
		{
			name:            "single transformation",
			r:               rand.New(rand.NewSource(1)),
			transformations: []domain.Transformation{transformation1},
			totalWeight:     1.0,
			want:            transformation1,
		},
		{
			name:            "choice falls in first bucket",
			r:               rand.New(rand.NewSource(9)),
			transformations: []domain.Transformation{transformation1, transformation2, transformation3},
			totalWeight:     6.0,
			want:            transformation1,
		},
		{
			name:            "choice falls in second bucket",
			r:               rand.New(rand.NewSource(2)),
			transformations: []domain.Transformation{transformation1, transformation2, transformation3},
			totalWeight:     6.0,
			want:            transformation2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chosenTransformation := SelectTransformation(tt.r, tt.transformations, tt.totalWeight)

			require.NotNil(t, chosenTransformation)
			assert.Equal(t, tt.want, chosenTransformation,
				"SelectTransformation() = %v, want %v", chosenTransformation, tt.want)
		})
	}
}
