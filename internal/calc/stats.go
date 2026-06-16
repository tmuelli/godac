package calc

import (
	"fmt"
	"math"

	"github.com/tmuelli/godac/internal/models"
)

// TODO: use embedded types instead of untyped structs
func getEffectiveStat(kind string, mon models.BattleMon) int {
	base := 0.0
	iv := 0.0
	ev := 0.0

	// get nature modifier
	nature := getNatureModifier(kind, mon.Nature)

	switch (kind) {
		case "hp":
			fmt.Println("Not implemented yet")
			return 0
		case "atk":
			base = float64(mon.BaseStats.Atk)
			iv = float64(mon.Ivs.Atk)
			ev = float64(mon.Evs.Atk)
		case "def":
			base = float64(mon.BaseStats.Def)
			iv = float64(mon.Ivs.Def)
			ev = float64(mon.Evs.Def)
		case "spa":
			base = float64(mon.BaseStats.Spa)
			iv = float64(mon.Ivs.Spa)
			ev = float64(mon.Evs.Spa)
		case "spd":
			base = float64(mon.BaseStats.Spd)
			iv = float64(mon.Ivs.Spd)
			ev = float64(mon.Evs.Spd)
		case "spe":
			fmt.Println("Not implemented yet")
			return 0
	}

	// calculate effective stat
	effStat := math.Abs((math.Abs((((2 * base) + iv + math.Abs(ev / 4)) * float64(mon.Level)) / 100) + 5) * nature)
	return int(math.Floor(effStat))
}