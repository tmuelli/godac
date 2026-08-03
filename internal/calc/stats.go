package calc

import (
	"fmt"
	"math"

	"github.com/tmuelli/godac/internal/models"
)

func getModifiedStat(stat float64, boost int) float64 {
	boostTable := [][]float64{
    	{2, 8},
    	{2, 7},
    	{2, 6},
    	{2, 5},
    	{2, 4},
    	{2, 3},
    	{2, 2},
    	{3, 2},
    	{4, 2},
    	{5, 2},
    	{6, 2},
    	{7, 2},
    	{8, 2},
	}

	// implemented OF16 and OF32 to be cartridge "correct" => following smogon damage calculation logic which implements this as well
	// technically this should never overflow because stats cant be that high => there are exceptions though e.g. hacked mons
	modStat := OF16(stat * boostTable[6 + boost][0]);
	modStat = math.Floor(stat / boostTable[6 + boost][1]);

	return modStat
}

func getSimpleModifiedStat(stat float64, boost int) float64 {
	simpleMod := math.Min(6, math.Max(-6, float64(boost * 2)))

	if simpleMod > 0 {
		return math.Floor((stat * (2 + simpleMod)) / 2)
	} else if simpleMod < 0 {
		return math.Floor((stat * 2) / (2 - simpleMod))
	} else {
		return stat
	}
}

func (d *CalcData) effectiveHpStat() int {
	base := float64(d.specie.BaseStats.Hp)
	iv := float64(d.ivs.Hp)
	ev := float64(d.evs.Hp)

	hp := math.Floor(math.Floor(2 * base) + iv + math.Floor(math.Floor(ev / 4) * float64(d.level)) / 100) + float64(d.level) + 10
	return int(hp)
}

func (d *CalcData) effectiveStat(kind StatKind, usedMove *models.Move, opponentAbility string, field *models.Field) int {
	base := 0.0
	iv := 0.0
	ev := 0.0

	// TODO: maybe put base, iv, ef stat into same form as boosts to get the stats more conveniently
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

	// get nature modifier
	nature := getNatureModifier(kind, d.nature)
	fmt.Println("[DEBUG]: Nature modifier:", nature)

	// calculate effective stat
	baseEffStat := math.Floor((math.Floor((((2 * base) + iv + math.Floor(ev / 4)) * float64(d.level)) / 100) + 5) * nature)

	// check ability stats modifications
	modStat := d.modifiedStat(kind, opponentAbility, usedMove, baseEffStat, field)

	// get item modifier
	item, err := d.itemModifier(&kind, models.StatEffect, nil)
	if err != nil {
		fmt.Println("Error during calculating effective stat:", err)
		return 0
	}

	effStat := math.Floor(modStat * item)
	return int(effStat)
}

func (d *CalcData) modifiedStat(kind StatKind, opponentAbility string, usedMove *models.Move, statValue float64, field *models.Field) float64 {
	modStat := 0.0

	// dont calculate for hp or speed stat
	switch kind {
	case StatAtk, StatDef, StatSpa, StatSpd:
	default:
		return statValue
	}

	statBoosts := d.boosts.Get(string(kind))

	// boosts are not accounted into modified stat if ability is unaware => ignores stat boosts
	if opponentAbility != "unaware" {
		if d.ability.Id == "simple" {
			statValue = getSimpleModifiedStat(statValue, statBoosts)
		} else {
			if (d.dataType == Attacker && (statBoosts > 0 || (!usedMove.Critical && statBoosts < 0))) ||
			(d.dataType == Defender && (statBoosts < 0 || (!usedMove.Critical && statBoosts > 0))) {
				statValue = getModifiedStat(statValue, statBoosts)
			}
		}
	}

	isPhysical := usedMove.Category == string(Physical)

	if d.ability.Active {
		switch d.dataType {
		case Attacker:
			if isPhysical && (d.ability.Id == "purepower" || d.ability.Id == "hugepower") {
				statValue *= 2
			} else if field.Terrain == models.Sun {
				if (isPhysical && d.ability.Id == "flowergift") ||
				(!isPhysical && d.ability.Id == "solarpower") {
					statValue = math.Floor(statValue * 1.5)
				}
			} else if (isPhysical && (d.ability.Id == "hustle" || d.ability.Id == "guts") && d.status != NoneStatus) ||
					(!isPhysical && (d.ability.Id == "plus" || d.ability.Id == "minus")) {
				statValue = math.Floor(statValue * 1.5)
			} else if isPhysical && d.ability.Id == "slowstart" {
				statValue = math.Floor(statValue / 2)
			}

			// TODO: check if this case can be combined with the defense case
			if (field.Right.FlowerGift || d.ability.Id == "flowergift") && field.Terrain == models.Sun && isPhysical {
				statValue = math.Floor(statValue * 1.5)
			}
		case Defender:
			if isPhysical && d.status != NoneStatus && d.ability.Id == "marvelscale" {
				statValue = math.Floor(statValue * 1.5)
			} else if !isPhysical && (field.Right.FlowerGift || d.ability.Id == "flowergift") && field.Terrain == models.Sun {
				// TODO: check if this case can be combined with the defense case
				statValue = math.Floor(statValue * 1.5)
			}

			if field.Terrain == models.Sand && (d.specieTypes[0].ID == "rock" || (len(d.specieTypes) > 1 && d.specieTypes[1].ID == "rock")) && !isPhysical {
				statValue = math.Floor(statValue * 1.5)
			}

			if usedMove.ID == "explosion" || usedMove.ID == "selfdestruct" {
				statValue = math.Floor(statValue * 0.5)
			}

			if statValue < 1 {
				statValue = 1
			}
		}
	}

	return modStat
}

func (d *CalcData) sumStatBoosts() int {
	sum := 0
	boosts := []string{string(StatAtk), string(StatDef), string(StatSpa), string(StatSpd), string(StatSpe)}
	
	for _, b := range boosts {
		boost := d.boosts.Get(b)
		if boost > 0 {
			sum += boost
		}
	}

	return sum
}