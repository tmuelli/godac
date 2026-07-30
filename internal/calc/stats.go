package calc

import (
	"fmt"
	"math"

	"github.com/tmuelli/godac/internal/models"
)

func (d *CalcData) effectiveHpStat() int {
	base := float64(d.specie.BaseStats.Hp)
	iv := float64(d.ivs.Hp)
	ev := float64(d.evs.Hp)

	hp := math.Floor(math.Floor(2 * base) + iv + math.Floor(math.Floor(ev / 4) * float64(d.level)) / 100) + float64(d.level) + 10
	return int(hp)
}

func (d *CalcData) effectiveStat(kind StatKind) int {
	base := 0.0
	iv := 0.0
	ev := 0.0

	// get nature modifier
	nature := getNatureModifier(kind, d.nature)

	// get item modifier
	item, err := d.itemModifier(&kind, models.StatEffect, nil)
	if err != nil {
		fmt.Println("Error during calculating effective stat:", err)
		return 0
	}

	switch (kind) {
		case StatHp:
			return d.effectiveHpStat()
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
			// TODO: implement and check for tailwind field flag
			base = float64(d.specie.BaseStats.Spe)
			iv = float64(d.ivs.Spe)
			ev = float64(d.evs.Spe)
	}

	fmt.Println("[DEBUG]: Base stat:", base)
	fmt.Println("[DEBUG]: IV stat:", iv)
	fmt.Println("[DEBUG]: EV stat:", ev)
	fmt.Println("[DEBUG]: Nature modifier:", nature)

	// TODO: implement field flower gift (from ability) combined with weather

	// calculate effective stat
	baseEffStat := math.Floor((math.Floor((((2 * base) + iv + math.Floor(ev / 4)) * float64(d.level)) / 100) + 5) * nature)
	effStat := math.Floor(baseEffStat * item)

	return int(effStat)
}