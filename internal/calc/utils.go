package calc

import "math"

const (
	MAX_UNSIGNED_INT_16 = 65535
	MAX_UNSIGNED_INT_32 = 4294967295
)

func combineMods(mods []int, lowerBound, upperBound int) float64 {
	mod := FIXED_POINT_SCALE_BASE

	for _, m := range mods {
		if m != FIXED_POINT_SCALE_BASE {
			mod = ((mod * m) + HALF_POINT_SCALE_BASE) >> 12
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

func OF16(v float64) float64 {
	value := int(v)

	if value > MAX_UNSIGNED_INT_16 {
		value = value % MAX_UNSIGNED_INT_16
	}

	return float64(value)
}

func OF32(v int) int {
	if v > MAX_UNSIGNED_INT_32 {
		return v % MAX_UNSIGNED_INT_32
	}

	return v
}