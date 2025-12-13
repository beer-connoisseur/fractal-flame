package transformations

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHorseshoeTransformation_Apply(t *testing.T) {
	tests := []struct {
		name  string
		x     float64
		y     float64
		wantX float64
		wantY float64
	}{
		{
			name:  "origin",
			x:     0.0,
			y:     0.0,
			wantX: 0.0,
			wantY: 0.0,
		},
		{
			name:  "right point",
			x:     1.0,
			y:     0.0,
			wantX: 1.0,
			wantY: 0.0,
		},
		{
			name:  "top point",
			x:     0.0,
			y:     1.0,
			wantX: -1.0,
			wantY: 0.0,
		},
		{
			name:  "left point",
			x:     -1.0,
			y:     0.0,
			wantX: 1.0,
			wantY: 0.0,
		},
		{
			name:  "bottom point",
			x:     0.0,
			y:     -1.0,
			wantX: -1.0,
			wantY: 0.0,
		},
		{
			name:  "45 degrees",
			x:     math.Sqrt2 / 2,
			y:     math.Sqrt2 / 2,
			wantX: 0.0,
			wantY: 1.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transformation := NewHorseshoeTransformation(1.0)
			gotX, gotY := transformation.Apply(tt.x, tt.y)

			assert.InDelta(t, tt.wantX, gotX, 0.001)
			assert.InDelta(t, tt.wantY, gotY, 0.001)
		})
	}
}

func TestNewHorseshoeTransformation(t *testing.T) {
	wantName := "horseshoe"
	tests := []struct {
		name   string
		weight float64
	}{
		{
			name:   "positive weight",
			weight: 1.5,
		},
		{
			name:   "zero weight",
			weight: 0.0,
		},
		{
			name:   "negative weight",
			weight: -0.5,
		},
		{
			name:   "fractional weight",
			weight: 0.333333,
		},
		{
			name:   "large weight",
			weight: 1000.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewHorseshoeTransformation(tt.weight)

			require.NotNil(t, got)
			require.NotNil(t, got.BaseTransformation)

			require.Equal(t, wantName, got.name)

			require.Equal(t, tt.weight, got.Weight())
		})
	}
}
