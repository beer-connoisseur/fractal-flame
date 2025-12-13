package renderer

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw4-fractal-flame/internal/domain"
)

func TestNewSymmetryApplier(t *testing.T) {
	t.Run("creates new symmetry applier", func(t *testing.T) {
		applier := NewSymmetryApplier()
		require.NotNil(t, applier)
		require.IsType(t, &SymmetryApplierImpl{}, applier)
	})
}

func TestSymmetryApplierImpl_Apply(t *testing.T) {
	type args struct {
		x     float64
		y     float64
		level int
	}
	tests := []struct {
		name       string
		x          float64
		y          float64
		level      int
		wantPoints int
		validate   func(t *testing.T, points []domain.Point)
		wantErr    bool
	}{
		{
			name:       "level 1",
			x:          1.0,
			y:          0.0,
			level:      1,
			wantPoints: 1,
			validate: func(t *testing.T, points []domain.Point) {
				require.Len(t, points, 1)
				require.InDelta(t, 1.0, points[0].X, 0.001)
				require.InDelta(t, 0.0, points[0].Y, 0.001)
			},
		},
		{
			name:       "level 2",
			x:          1.0,
			y:          0.0,
			level:      2,
			wantPoints: 2,
			validate: func(t *testing.T, points []domain.Point) {
				require.Len(t, points, 2)
				require.InDelta(t, 1.0, points[0].X, 0.001)
				require.InDelta(t, 0.0, points[0].Y, 0.001)
				require.InDelta(t, -1.0, points[1].X, 0.001)
				require.InDelta(t, 0.0, points[1].Y, 0.001)
			},
		},
		{
			name:       "level 3",
			x:          1.0,
			y:          0.0,
			level:      3,
			wantPoints: 3,
			validate: func(t *testing.T, points []domain.Point) {
				require.Len(t, points, 3)
				angles := []float64{0.0, 120.0, 240.0}
				for i, angle := range angles {
					rad := angle * math.Pi / 180.0
					expectedX := math.Cos(rad)
					expectedY := math.Sin(rad)
					require.InDelta(t, expectedX, points[i].X, 0.001)
					require.InDelta(t, expectedY, points[i].Y, 0.001)
				}
			},
		},
		{
			name:       "level 4",
			x:          1.0,
			y:          0.0,
			level:      4,
			wantPoints: 4,
			validate: func(t *testing.T, points []domain.Point) {
				require.Len(t, points, 4)
				expectedPoints := []domain.Point{
					{X: 1.0, Y: 0.0},
					{X: 0.0, Y: 1.0},
					{X: -1.0, Y: 0.0},
					{X: 0.0, Y: -1.0},
				}
				for i, expected := range expectedPoints {
					require.InDelta(t, expected.X, points[i].X, 0.001)
					require.InDelta(t, expected.Y, points[i].Y, 0.001)
				}
			},
		},
		{
			name:       "level 0",
			x:          1.0,
			y:          2.0,
			level:      0,
			wantPoints: 1,
			validate: func(t *testing.T, points []domain.Point) {
				require.Len(t, points, 1)
				require.InDelta(t, 1.0, points[0].X, 0.001)
				require.InDelta(t, 2.0, points[0].Y, 0.001)
			},
		},
		{
			name:       "level negative",
			x:          1.0,
			y:          2.0,
			level:      -3,
			wantPoints: 1,
			validate: func(t *testing.T, points []domain.Point) {
				require.Len(t, points, 1)
				require.InDelta(t, 1.0, points[0].X, 0.001)
				require.InDelta(t, 2.0, points[0].Y, 0.001)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SymmetryApplierImpl{}
			points := s.Apply(tt.x, tt.y, tt.level)

			require.Len(t, points, tt.wantPoints)

			if tt.validate != nil {
				tt.validate(t, points)
			}
		})
	}
}
