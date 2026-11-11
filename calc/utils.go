package calc

import "math"

const (
	maxUnsignedInt16 = 65535
	maxUnsignedInt32 = 4294967295
)

func combineMods(mods []int, lowerBound, upperBound int) float64 {
	mod := fixedPointScaleBase

	for _, m := range mods {
		if m != fixedPointScaleBase {
			mod = ((mod * m) + halfPointScaleBase) >> 12
		}
	}

	return math.Max(math.Min(float64(mod), float64(upperBound)), float64(lowerBound))
}

// Game Freak rounds down on .5
func pokeRound(number float64) int {
	var roundedValue float64

	if math.Mod(number, 1) > 0.5 {
		roundedValue = math.Ceil(number)
	} else {
		roundedValue = math.Floor(number)
	}

	return int(roundedValue)
}

func of16(v float64) float64 {
	value := int(v)

	if value > maxUnsignedInt16 {
		value = value % maxUnsignedInt16
	}

	return float64(value)
}

func of32(v int) int {
	if v > maxUnsignedInt32 {
		return v % maxUnsignedInt32
	}

	return v
}
