package utils

import (
	"math/rand"

	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw4-fractal-flame/internal/domain"
)

func CalculateTotalWeight(transformations []domain.Transformation) float64 {
	total := 0.0
	for _, t := range transformations {
		total += t.Weight()
	}
	return total
}

func SelectTransformation(
	r *rand.Rand,
	transformations []domain.Transformation,
	totalWeight float64,
) domain.Transformation {
	if len(transformations) == 1 {
		return transformations[0]
	}

	choice := r.Float64() * totalWeight
	cumulative := 0.0

	for _, t := range transformations {
		cumulative += t.Weight()
		if choice <= cumulative {
			return t
		}
	}

	return transformations[0]
}
