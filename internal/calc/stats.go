package calc

import (
	"math"

	"github.com/tmuelli/godac/internal/models"
	"github.com/tmuelli/godac/internal/utils"
)

func getEVItems() []string {
	return []string{
		"machobrace",
		"poweranklet",
		"powerband",
		"powerbelt",
		"powerbracer",
		"powerlens",
		"powerweight",
	}
}

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
	modStat := OF16(stat * boostTable[6+boost][0])
	modStat = math.Floor(modStat / boostTable[6+boost][1])

	return modStat
}

func getSimpleModifiedStat(stat float64, boost int) float64 {
	simpleMod := math.Min(6, math.Max(-6, float64(boost*2)))

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

	hp := utils.GetMaxHP(base, iv, ev, d.level)
	return int(hp)
}

func getBaseSpeed(base, iv, ev float64, level int, nature *models.Nature) float64 {
	// nature modifier
	n := getNatureModifier(StatSpe, nature)
	return math.Floor((math.Floor((((2*base)+iv+math.Floor(ev/4))*float64(level))/100) + 5) * n)
}

func (d *CalcData) effectiveSpeedStat(field *models.Field) int {
	speedMods := []int{}

	// first calculate base speed with base stat, iv and ev and nature
	baseSpeed := getBaseSpeed(
		float64(d.specie.BaseStats.Spe),
		float64(d.ivs.Spe),
		float64(d.evs.Spe),
		d.level,
		d.nature,
	)

	baseModifiedSpeed := getModifiedStat(baseSpeed, d.boosts.Get(string(StatSpe)))

	if (d.dataType == models.Attacker && field.Left.Tailwind) ||
		(d.dataType == models.Defender && field.Right.Tailwind) {
		speedMods = append(speedMods, SPEED_MOD_2X)
	}

	if d.hasAbility("unburden") ||
		(d.hasAbility("chlorophyll") && field.Terrain == models.Sun) ||
		(d.hasAbility("swiftswim") && field.Terrain == models.Rain) {
		speedMods = append(speedMods, SPEED_MOD_2X)
	} else if d.hasAbility("quickfeet") && d.status != NoneStatus {
		speedMods = append(speedMods, SPEED_MOD_1_5X)
	} else if d.hasAbility("slowstart") {
		speedMods = append(speedMods, SPEED_MOD_0_5X)
	}

	if !d.hasAbility("unburden") {
		if d.hasItem("choicescarf") {
			speedMods = append(speedMods, SPEED_MOD_1_5X)
		} else if d.hasItem(append([]string{"ironball"}, getEVItems()...)...) {
			speedMods = append(speedMods, SPEED_MOD_0_5X)
		} else if d.hasItem("quickpowder") && d.specie.ID == "ditto" {
			speedMods = append(speedMods, SPEED_MOD_2X)
		}
	}

	modSpeed := OF32(pokeRound((baseModifiedSpeed * combineMods(speedMods, COMBINE_MODS_MIN, COMBINE_MODS_MAX)) / FIXED_POINT_SCALE_BASE))

	if d.status == Paralyzed && !d.hasAbility("quickfeet") {
		modSpeed = int(math.Floor(float64(OF32(modSpeed*25) / 100)))
	}

	effSpeed := math.Min(float64(MAX_SPEED), float64(modSpeed))
	return int(math.Max(0, effSpeed))
}

func (d *CalcData) effectiveStat(kind StatKind, usedMove *models.Move, hasOpponentUnaware bool, field *models.Field) int {
	base := 0.0
	iv := 0.0
	ev := 0.0

	switch kind {
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
		return d.effectiveSpeedStat(field)
	}

	// get nature modifier
	nature := getNatureModifier(kind, d.nature)

	// calculate effective stat
	baseEffStat := math.Floor((math.Floor((((2*base)+iv+math.Floor(ev/4))*float64(d.level))/100) + 5) * nature)

	// check ability stats modifications
	modStat := d.modifiedStat(kind, hasOpponentUnaware, usedMove, baseEffStat, field)

	// get item modifier
	item, err := d.itemModifier(&kind, models.StatEffect, nil)
	if err != nil {
		// TODO: add to logging
		// fmt.Println("Error during calculating effective stat:", err)
		return 0
	}

	effStat := math.Floor(modStat * item)
	return int(effStat)
}

func (d *CalcData) modifiedStat(kind StatKind, hasOpponentUnaware bool, usedMove *models.Move, statValue float64, field *models.Field) float64 {
	// dont calculate for hp or speed stat
	switch kind {
	case StatAtk, StatDef, StatSpa, StatSpd:
	default:
		return statValue
	}

	statBoosts := d.boosts.Get(string(kind))

	// boosts are not accounted into modified stat if ability is unaware => ignores stat boosts
	if !hasOpponentUnaware {
		if d.hasAbility("simple") {
			statValue = getSimpleModifiedStat(statValue, statBoosts)
		} else {
			if (d.dataType == models.Attacker && (statBoosts > 0 || (!usedMove.Critical && statBoosts < 0))) ||
				(d.dataType == models.Defender && (statBoosts < 0 || (!usedMove.Critical && statBoosts > 0))) {
				statValue = getModifiedStat(statValue, statBoosts)
			}
		}
	}

	isPhysical := usedMove.Category == string(Physical)

	switch d.dataType {
	case models.Attacker:
		if isPhysical && d.hasAbility("purepower", "hugepower") {
			statValue *= 2
		} else if (field.Terrain == models.Sun) && ((isPhysical && d.hasAbility("flowergift")) || (!isPhysical && d.hasAbility("solarpower"))) {
			statValue = math.Floor(statValue * 1.5)
		} else if (isPhysical && (d.hasAbility("hustle") || (d.hasAbility("guts") && d.status != NoneStatus))) ||
			(!isPhysical && (d.hasAbility("plus", "minus"))) {
			statValue = math.Floor(statValue * 1.5)
		} else if isPhysical && d.hasAbility("slowstart") {
			statValue = math.Floor(statValue / 2)
		}

		if field.Left.FlowerGift && !d.hasAbility("flowergift") && field.Terrain == models.Sun && isPhysical {
			statValue = math.Floor(statValue * 1.5)
		}
	case models.Defender:
		if isPhysical && d.status != NoneStatus && d.hasAbility("marvelscale") {
			statValue = math.Floor(statValue * 1.5)
		} else if !isPhysical && (field.Right.FlowerGift || d.hasAbility("flowergift")) && field.Terrain == models.Sun {
			statValue = math.Floor(statValue * 1.5)
		}

		if field.Terrain == models.Sand && d.hasType("rock") && !isPhysical {
			statValue = math.Floor(statValue * 1.5)
		}

		if usedMove.ID == "explosion" || usedMove.ID == "selfdestruct" {
			statValue = math.Floor(statValue * 0.5)
		}

		if statValue < 1 {
			statValue = 1
		}
	}

	return statValue
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
