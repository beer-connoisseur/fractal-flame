package transformations

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBaseTransformation_Weight(t *testing.T) {
	tests := []struct {
		name   string
		weight float64
		want   float64
	}{
		{
			name:   "get positive weight",
			weight: 3.14,
			want:   3.14,
		},
		{
			name:   "get zero weight",
			weight: 0.0,
			want:   0.0,
		},
		{
			name:   "get negative weight",
			weight: -1.5,
			want:   -1.5,
		},
		{
			name:   "get very small weight",
			weight: 1e-10,
			want:   1e-10,
		},
		{
			name:   "get very large weight",
			weight: 1e+10,
			want:   1e+10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transformation := &BaseTransformation{
				name:   tt.name,
				weight: tt.weight,
			}
			got := transformation.Weight()

			assert.InDelta(t, tt.want, got, 0.001)
		})
	}
}

func TestNewBaseTransformation(t *testing.T) {
	tests := []struct {
		name   string
		weight float64
		want   *BaseTransformation
	}{
		{
			name:   "create with positive weight",
			weight: 1.5,
			want: &BaseTransformation{
				name:   "test",
				weight: 1.5,
			},
		},
		{
			name:   "create with zero weight",
			weight: 0.0,
			want: &BaseTransformation{
				name:   "zero_weight",
				weight: 0.0,
			},
		},
		{
			name:   "create with negative weight",
			weight: -2.0,
			want: &BaseTransformation{
				name:   "negative",
				weight: -2.0,
			},
		},
		{
			name:   "create with fractional weight",
			weight: 0.333333,
			want: &BaseTransformation{
				name:   "fractional",
				weight: 0.333333,
			},
		},
		{
			name:   "create with empty name",
			weight: 1.0,
			want: &BaseTransformation{
				name:   "",
				weight: 1.0,
			},
		},
		{
			name:   "create with special characters in name",
			weight: 2.0,
			want: &BaseTransformation{
				name:   "test-transformation_123",
				weight: 2.0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewBaseTransformation(tt.want.name, tt.weight)
			require.NotNil(t, got)
			require.Equal(t, tt.want.name, got.name)
			require.Equal(t, tt.want.weight, got.weight)
			require.NotNil(t, got, "NewBaseTransformation returned nil")
		})
	}
}
