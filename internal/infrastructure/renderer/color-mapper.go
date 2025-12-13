package renderer

import (
	"math/rand"

	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw4-fractal-flame/internal/domain"
)

type ColorMapper struct{}

func NewHSVColorMapper() *ColorMapper {
	return &ColorMapper{}
}

func (m *ColorMapper) GeneratePalette(size int, seed int64) []domain.Color {
	if size <= 0 {
		return nil
	}
	r := rand.New(rand.NewSource(seed))
	colors := make([]domain.Color, size)

	for i := 0; i < size; i++ {
		colors[i] = domain.Color{
			R: r.Float64(),
			G: r.Float64(),
			B: r.Float64(),
		}
	}

	return colors
}
