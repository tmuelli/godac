package calc

import (
	"fmt"
	"math"
)

func (d *CalcData) effectiveStat(kind StatKind) int {
	base := 0.0
	iv := 0.0
	ev := 0.0

	// get nature modifier
	nature := getNatureModifier(kind, d.nature)

	switch (kind) {
		case StatHp:
			fmt.Println("[DEBUG]: Not implemented yet")
			return 0
		case StatAtk:
			base = float64(d.specie.BaseStats.Atk)
			iv = float64(d.ivs.Atk)
			ev = float64(d.evs.Atk)
		case StatDef:
			base = float64(d.specie.BaseStats.Def)
			iv = float64(d.ivs.Def)
			ev = float64(d.evs.Def)
		case StatSpa:
			base = float64(d.specie.BaseStats.Spa)
			iv = float64(d.ivs.Spa)
			ev = float64(d.evs.Spa)
		case StatSpd:
			base = float64(d.specie.BaseStats.Spd)
			iv = float64(d.ivs.Spd)
			ev = float64(d.evs.Spd)
		case StatSpe:
			fmt.Println("[DEBUG]: Not implemented yet")
			return 0
	}

	fmt.Println("[DEBUG]: Base stat:", base)
	fmt.Println("[DEBUG]: IV stat:", iv)
	fmt.Println("[DEBUG]: EV stat:", ev)
	fmt.Println("[DEBUG]: Nature modifier:", nature)

	// calculate effective stat
	effStat := math.Abs((math.Abs((((2 * base) + iv + math.Abs(ev / 4)) * float64(d.level)) / 100) + 5) * nature)
	return int(math.Floor(effStat))
}