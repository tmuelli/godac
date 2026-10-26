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
		case "fire":
			return 1.5
		case "water":
			return 0.5
		}

	case models.Rain:
		switch c.attacker.usedMove.Type {
		case "water":
			return 1.5
		case "fire":
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

func getTypeStringFromWeather(terrain models.FieldTerrain) string {
	switch terrain {
	case models.Sun:
		return "fire"
	case models.Rain:
		return "water"
	case models.Hail:
		return "ice"
	case models.Sand:
		return "rock"
	default:
		return "normal"
	}
}
