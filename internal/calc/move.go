package calc

import (
	"math"

	"github.com/tmuelli/godac/internal/models"
)

func (c *Calculation) getMoveBasePower() int {
	basePower := c.attacker.usedMove.BasePower
	atkHp := c.attacker.effectiveHpStat()
	defHp := c.defender.effectiveHpStat()

	switch c.attacker.usedMove.ID {
	case "brine":
		if float64(c.defender.health) <= float64(defHp / 2) {
			basePower *= 2
		}
	case "eruption", "waterspout":
		basePower = int(math.Max(1, math.Floor(float64(basePower * c.attacker.health) / float64(atkHp))))
	case "facade":
		if c.attacker.status == BadlyPoisoned ||
		   c.attacker.status == Burned ||
		   c.attacker.status == Paralyzed ||
		   c.attacker.status == Poisoned {
			basePower *= 2
		}
	case "flail", "reversal":
		perc := math.Floor(float64(64 * c.attacker.health) / float64(atkHp))
		if perc <= 1 {
			basePower = 200
		} else if perc <= 5 {
			basePower = 150
		} else if perc <= 12 {
			basePower = 100
		} else if perc <= 21 {
			basePower = 80
		} else if perc <= 42 {
			basePower = 40
		} else {
			basePower = 20
		}
	case "fling":
		if c.attacker.item != nil {
			basePower = c.attacker.item.FlingPower
		}
	case "grassknot", "lowkick":
		weight := c.defender.specie.Weightkg
		if weight >= 200 {
			basePower = 120
		} else if weight >= 100 {
			basePower = 100
		} else if weight >= 50 {
			basePower = 80
		} else if weight >= 25 {
			basePower = 60
		} else if weight >= 10 {
			basePower = 40
		} else {
			basePower = 20
		}
	case "gyroball":
		atkSpeed := c.attacker.effectiveStat(StatSpe)
		defSpeed := c.defender.effectiveStat(StatSpe)

		basePower = int(math.Min(150.0, math.Floor(float64((25 * defSpeed) / atkSpeed))))
	case "payback":
		atkSpeed := c.attacker.effectiveStat(StatSpe)
		defSpeed := c.defender.effectiveStat(StatSpe)

		if defSpeed >= atkSpeed {
			basePower *= 2
		}
	case "punishment":
		basePower = int(math.Min(200.0, float64(60 + (20 * c.defender.sumStatBoosts()))))
	case "pursuit":
		if c.field.Right.SwitchingOut {
			basePower *= 2
		}
	case "wakeupslap":
		if c.defender.status == Asleep {
			basePower *= 2
		}
	case "naturepower":
		// TODO: handle naturepower special base power case
		// damage-calc repo also doesnt implement further terrain checks
		basePower = 80
	case "crushgrip", "wringout":
		basePower = int(math.Floor(float64(c.defender.health * 120) / float64(defHp)) + 1)
	case "triplekick":
		// TODO: implement hits field on move
	case "weatherball":
		if c.field.Terrain != models.Clear {
			basePower *= 2
		}
	}

	return basePower
}