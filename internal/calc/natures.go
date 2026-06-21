package calc

import (
	"github.com/tmuelli/godac/internal/models"
)

func getNatureModifier(stat StatKind, nature *models.Nature) float64 {
	switch(string(stat)) {
	case nature.Plus:
		return 1.1
	case nature.Minus:
		return 0.9
	default:
		return 1.0
	}
}