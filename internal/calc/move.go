package calc

import (
	"fmt"
	"math"

	"github.com/tmuelli/godac/internal/models"
)

func (c *Calculation) moveBasePower() int {
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
		// TODO: use getFinalSpeed => to be implemented
		// maybe move effectiveStat(StatSpe) into effectiveSpeedStat()
		// inside call effectiveStat and additionally check tailwind and stuff
		atkSpeed := c.attacker.effectiveStat(StatSpe, nil, false, nil)
		defSpeed := c.defender.effectiveStat(StatSpe, nil, false, nil)

		basePower = int(math.Min(150.0, math.Floor(float64((25 * defSpeed) / atkSpeed))))
	case "payback":
		// TODO: use getFinalSpeed => to be implemented
		atkSpeed := c.attacker.effectiveStat(StatSpe, nil, false, nil)
		defSpeed := c.defender.effectiveStat(StatSpe, nil, false, nil)

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

func (c *Calculation) modifiedMovePower(basePower float64) (float64, error) {
	// TODO: maybe find better approach to handle abiltiy modifier of move base power
	var kind StatKind

	if MoveCategory(c.attacker.usedMove.Category) == Special {
		kind = StatSpa
	} else if MoveCategory(c.attacker.usedMove.Category) == Physical {
		kind = StatAtk
	} else {
		return 0.0, fmt.Errorf("invalid move category %v", c.attacker.usedMove.Category)
	}

	// helping hand modifier
	if c.field.Left.HelpingHand {
		basePower = math.Floor(basePower * 1.5)
	}

	// abilitiy technician modifier for move base power
	if c.attacker.hasAbility("technician") && c.attacker.usedMove.BasePower <= 60 {
		basePower = math.Floor(basePower * 1.5)
	}

	// item modifier for move base power
	itemMod, err := c.attacker.itemModifier(&kind, models.MoveEffect, nil)
	if err != nil {
		return 0, err
	}

	basePower = math.Floor(basePower * itemMod)

	// ability reckless and iron fist
	if (c.attacker.hasAbility("reckless") && (len(c.attacker.usedMove.Recoil) > 1 || c.attacker.usedMove.HasCrashDamage)) ||
	   (c.attacker.hasAbility("ironfist") && c.attacker.usedMove.Flags.Punch != 0) {
		basePower = math.Floor(basePower * 1.2)
	}

	// ability overgrow, blaze, torrent, swarm
	atkHp := c.attacker.effectiveHpStat()
	if float64(c.attacker.health) <= float64(atkHp / 3) {
		if (c.attacker.hasAbility("overgrow") && c.attacker.usedMove.Type == "grass") ||
		   (c.attacker.hasAbility("blaze") && c.attacker.usedMove.Type == "fire") ||
		   (c.attacker.hasAbility("torrent") && c.attacker.usedMove.Type == "water") ||
		   (c.attacker.hasAbility("swarm") && c.attacker.usedMove.Type == "bug") {
			basePower = math.Floor(basePower * 1.5)
		}
	}

	// ability heatproff and thick fat
	if (c.defender.hasAbility("heatproof") && c.attacker.usedMove.Type == "fire") ||
	   (c.defender.hasAbility("thickfat") && (c.attacker.usedMove.Type == "fire" || c.attacker.usedMove.Type == "ice")) {
		basePower = math.Floor(basePower * 0.5)
	}

	// ability dry skin
	if c.defender.hasAbility("dryskin") && c.attacker.usedMove.Type == "fire" {
		basePower = math.Floor(basePower * 1.25)
	}

	// ability rivalry
	if c.attacker.gender != nil && c.defender.gender != nil {
		atkGender := *c.attacker.gender
		defGender := *c.defender.gender
		if c.attacker.hasAbility("rivalry") && atkGender != Neutral && defGender != Neutral {
			if atkGender == defGender {
				basePower = math.Floor(basePower * 1.25)
			} else {
				basePower = math.Floor(basePower * 0.75)
			}
		}
	}

	return basePower, nil
}