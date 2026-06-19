package calc

import (
	"fmt"
	"math"
)

// TODO: use embedded types instead of untyped structs
func getEffectiveStat(kind string, data CalcData) int {
	base := 0.0
	iv := 0.0
	ev := 0.0

	// get nature modifier
	nature := getNatureModifier(kind, data.nature.Name)

	switch (kind) {
		case "hp":
			fmt.Println("Not implemented yet")
			return 0
		case "atk":
			base = float64(data.specie.BaseStats.Atk)
			iv = float64(data.ivs.Atk)
			ev = float64(data.evs.Atk)
		case "def":
			base = float64(data.specie.BaseStats.Def)
			iv = float64(data.ivs.Def)
			ev = float64(data.evs.Def)
		case "spa":
			base = float64(data.specie.BaseStats.Spa)
			iv = float64(data.ivs.Spa)
			ev = float64(data.evs.Spa)
		case "spd":
			base = float64(data.specie.BaseStats.Spd)
			iv = float64(data.ivs.Spd)
			ev = float64(data.evs.Spd)
		case "spe":
			fmt.Println("Not implemented yet")
			return 0
	}

	fmt.Println("Base stat:", base)
	fmt.Println("IV stat:", iv)
	fmt.Println("EV stat:", ev)

	// calculate effective stat
	effStat := math.Abs((math.Abs((((2 * base) + iv + math.Abs(ev / 4)) * float64(data.level)) / 100) + 5) * nature)
	return int(math.Floor(effStat))
}