package renderer

import (
	"math"

	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw4-fractal-flame/internal/domain"
)

type SymmetryApplierImpl struct{}

func NewSymmetryApplier() *SymmetryApplierImpl {
	return &SymmetryApplierImpl{}
}

func (s *SymmetryApplierImpl) Apply(x, y float64, level int) []domain.Point {
	if level <= 1 {
		return []domain.Point{{X: x, Y: y}}
	}

	points := make([]domain.Point, level)
	angleStep := 2.0 * math.Pi / float64(level)

	for i := 0; i < level; i++ {
		angle := float64(i) * angleStep
		cosA := math.Cos(angle)
		sinA := math.Sin(angle)

		points[i] = domain.Point{
			X: x*cosA - y*sinA,
			Y: x*sinA + y*cosA,
		}
	}

	return points
}
