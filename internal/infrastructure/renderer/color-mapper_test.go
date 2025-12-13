package renderer

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw4-fractal-flame/internal/domain"
)

func TestColorMapper_GeneratePalette(t *testing.T) {
	tests := []struct {
		name       string
		size       int
		seed       int64
		wantLength int
		validate   func(t *testing.T, palette []domain.Color)
	}{
		{
			name:       "generate palette with size 1",
			size:       1,
			seed:       12345,
			wantLength: 1,
			validate: func(t *testing.T, palette []domain.Color) {
				require.Len(t, palette, 1)
				color := palette[0]
				require.InDelta(t, 0.8487, color.R, 0.001)
				require.InDelta(t, 0.6458, color.G, 0.001)
				require.InDelta(t, 0.7382, color.B, 0.001)
			},
		},
		{
			name:       "generate palette with size 5",
			size:       5,
			seed:       42,
			wantLength: 5,
			validate: func(t *testing.T, palette []domain.Color) {
				require.Len(t, palette, 5)

				for i, color := range palette {
					require.True(t, color.R >= 0.0 && color.R <= 1.0,
						"Color %d R component out of range: %f", i, color.R)
					require.True(t, color.G >= 0.0 && color.G <= 1.0,
						"Color %d G component out of range: %f", i, color.G)
					require.True(t, color.B >= 0.0 && color.B <= 1.0,
						"Color %d B component out of range: %f", i, color.B)
				}

				require.InDelta(t, 0.3730, palette[0].R, 0.001)
				require.InDelta(t, 0.066, palette[0].G, 0.001)
				require.InDelta(t, 0.6040, palette[0].B, 0.001)
			},
		},
		{
			name:       "generate palette with size 0",
			size:       0,
			seed:       100,
			wantLength: 0,
			validate: func(t *testing.T, palette []domain.Color) {
				require.Empty(t, palette)
			},
		},
		{
			name:       "generate palette with negative size",
			size:       -5,
			seed:       100,
			wantLength: 0,
			validate: func(t *testing.T, palette []domain.Color) {
				require.Empty(t, palette)
			},
		},
		{
			name:       "same seed produces same palette",
			size:       3,
			seed:       777,
			wantLength: 3,
			validate: func(t *testing.T, palette []domain.Color) {
				mapper := &ColorMapper{}
				palette2 := mapper.GeneratePalette(3, 777)

				require.Equal(t, len(palette), len(palette2))
				for i := range palette {
					require.InDelta(t, palette[i].R, palette2[i].R, 0.000001)
					require.InDelta(t, palette[i].G, palette2[i].G, 0.000001)
					require.InDelta(t, palette[i].B, palette2[i].B, 0.000001)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &ColorMapper{}
			palette := m.GeneratePalette(tt.size, tt.seed)

			require.Len(t, palette, tt.wantLength)

			if tt.validate != nil {
				tt.validate(t, palette)
			}
		})
	}
}

func TestNewHSVColorMapper(t *testing.T) {
	t.Run("creates new color mapper", func(t *testing.T) {
		mapper := NewHSVColorMapper()
		require.NotNil(t, mapper)
		require.IsType(t, &ColorMapper{}, mapper)
	})
}
