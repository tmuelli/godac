package calc

import "github.com/tmuelli/godac/internal/models"

func (c *Calculation) weatherModifier() float64 {
	if c.attacker.usedMove.ID == "solarbeam" {
		if c.field.Terrain == models.Clear || c.field.Terrain == models.Sun {
			return 1.0
		}

		return 0.5
	}

	switch c.field.Terrain {
	case models.Clear:
	case models.Sun:
		switch c.attacker.usedMove.Type {
		case "Fire":
			return 1.5
		case "Water":
			return 0.5
		}

	case models.Rain:
		switch c.attacker.usedMove.Type {
		case "Water":
			return 1.5
		case "Fire":
			return 0.5
		}

	case models.Sand:
		// chip damage is calculated after real damage calculation
		return 1.0
	case models.Hail:
		// chip damage is calculated after real damage calculation
		return 1.0
	}

	return 1.0
}