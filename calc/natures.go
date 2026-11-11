package calc

import (
	"github.com/tmuelli/godac/models"
)

func getNatureModifier(stat statKind, nature *models.Nature) float64 {
	if nature.Plus == nature.Minus {
		return 1.0
	}

	switch string(stat) {
	case nature.Plus:
		return 1.1
	case nature.Minus:
		return 0.9
	default:
		return 1.0
	}
}
