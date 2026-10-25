package calc

import (
	"fmt"
	"math"

	"github.com/tmuelli/godac/internal/models"
)

type CalcData struct {
	specie      *models.Specie
	specieTypes []*models.Type
	nature      *models.Nature
	ability     *models.Ability
	item        *models.Item
	usedMove    *models.Move
	ivs         *models.IvStats
	evs         *models.EvStats
	level       int
	health      int
	dataType    CalcDataType
	status      StatusCondition
	boosts      *models.Boosts
	gender      *Gender
	metrHits    int
}

type Calculation struct {
	data     *models.Data
	attacker *CalcData
	defender *CalcData
	field    *models.Field
	typeEff  TypeEffectivness
	Result   models.BattleResult
}

func NewCalc(data *models.Data, atkInput *models.UserInput, defInput *models.UserInput, fieldInput *models.UserFieldInput) (*Calculation, error) {
	// validate inputs
	attacker, err := validateInput(data, atkInput)
	if err != nil {
		return nil, fmt.Errorf("invalid attack input: %v", err)
	}

	defender, err := validateInput(data, defInput)
	if err != nil {
		return nil, fmt.Errorf("invalid defender input: %v", err)
	}

	field, err := validateFieldInput(fieldInput)
	if err != nil {
		return nil, fmt.Errorf("invalid field input: %v", err)
	}

	// TODO: decide on precedence of level input =>
	// if attacker / defender level is set or field level is set or both
	// determine on defaults and "last set" option
	// FOR NOW: only UserInput level is used

	return &Calculation{
		data:     data,
		attacker: attacker,
		defender: defender,
		field:    field,
		Result: models.BattleResult{
			Range: models.BattleResultRange{
				Min:     0,
				MinPerc: 0,
				Max:     0,
				MaxPerc: 0,
			},
			DmgMatrix: [][]int{},
		},
	}, nil
}

func (c *Calculation) needCalc() bool {
	if (c.attacker.usedMove.Category == string(Status) && c.attacker.usedMove.ID != "naturepower") ||
		(c.field.Right.Protect && !c.attacker.usedMove.BreaksProtect) {
		return false
	}

	return true
}

func (c *Calculation) Calc() error {
	// setup calculation
	err := c.setup()
	if err != nil {
		return err
	}

	if !c.needCalc() {
		return nil
	}

	// handle pain split
	if c.attacker.usedMove.ID == "painsplit" {
		dmg := c.attacker.usedMove.HandlePainSplit(c.attacker.health, c.defender.health)
		c.Result.Range = models.BattleResultRange{
			Min:     dmg,
			Max:     dmg,
			MinPerc: 0,
			MaxPerc: 0,
		}
		c.Result.DmgMatrix = [][]int{{dmg}}
		return nil
	}

	// check damage prevention from ability
	err = c.setTypeModifier()
	if err != nil {
		return err
	}

	if c.typeEff.activeValue == 0.0 {
		return nil
	}

	if c.defender.ability.IsPreventingDamage(c.attacker.usedMove, c.defender.item, c.typeEff.activeValue, c.field.IsGravity) {
		return nil
	}

	// handle fixed damage
	if dmg := c.handleFixedDamage(); dmg != 0 {
		c.Result.Range = models.BattleResultRange{
			Min:     dmg,
			Max:     dmg,
			MinPerc: 0,
			MaxPerc: 0,
		}
		c.Result.DmgMatrix = [][]int{{dmg}}
		return nil
	}

	// calculate base damage
	baseDmg, err := c.calcBaseDamage(DEFAULT_HIT_COUNT)
	if err != nil {
		return err
	}

	// calculate conditional damage
	condDmg, err := c.calcConditionalDamage(baseDmg)
	if err != nil {
		return err
	}

	// calcute modification damage
	dmgResult, err := c.calcModificationDamage(condDmg)
	if err != nil {
		return err
	}

	c.Result.DmgMatrix = append(c.Result.DmgMatrix, dmgResult)
	dmgMin := dmgResult[0]
	dmgMax := dmgResult[15]

	// TODO: implement final dot damage (e.g. stealth rock, spikes, leech seed) from field
	// in smogon damage calc this is not inside the gen4 damage calc but rather calculated afterwards
	// TODO: also implement chip damage from sand / hail here (from field weather)

	c.Result.Range.Min = int(dmgMin)
	c.Result.Range.Max = int(dmgMax)

	minPerc, maxPerc := c.percentRange()
	c.Result.Range.MinPerc = minPerc
	c.Result.Range.MaxPerc = maxPerc

	if c.attacker.usedMove.Hits > DEFAULT_HIT_COUNT {
		usedItems := models.MultihitItemUsage{
			HasAtkUsedItem: false,
			HasDefUsedItem: false,
		}

		for hit := 1; hit < c.attacker.usedMove.Hits; hit++ {
			// check multihit item usage
			c.checkMultihitBoost(&usedItems)

			// calculate base damage
			baseDmg, err := c.calcBaseDamage(hit + 1)
			if err != nil {
				return err
			}

			// calculate conditional damage
			condDmg, err := c.calcConditionalDamage(baseDmg)
			if err != nil {
				return err
			}

			// calcute modification damage
			dmgResult, err := c.calcModificationDamage(condDmg)
			if err != nil {
				return err
			}

			c.Result.DmgMatrix = append(c.Result.DmgMatrix, dmgResult)
		}
	}

	return nil
}

func (c *Calculation) handleFixedDamage() int {
	switch c.attacker.usedMove.ID {
	case "seismictoss", "nightshade":
		return c.attacker.level

	case "dragonrage":
		return FIXED_DMG_DRAGON_RAGE

	case "sonicboom":
		return FIXED_DMG_SONIC_BOOM
	}

	return 0
}

func (c *Calculation) calcBaseDamage(hitCount int) (float64, error) {
	// get attack and Defense stat based on category
	var attackStat int
	var defenseStat int
	var atkStat StatKind
	var defStat StatKind

	atkCategory := c.attacker.usedMove.Category
	atkMove := c.attacker.usedMove

	switch MoveCategory(atkCategory) {
	case Physical:
		atkStat = StatAtk
		defStat = StatDef
	case Special:
		atkStat = StatSpa
		defStat = StatSpd
	default:
		return 0.0, fmt.Errorf("invalid attack move category")
	}

	attackStat = c.attacker.effectiveStat(atkStat, atkMove, c.defender.hasAbility("unaware"), c.field)
	defenseStat = c.defender.effectiveStat(defStat, atkMove, c.attacker.hasAbility("unaware"), c.field)

	moveBP := c.moveBasePower(hitCount)
	movePower, err := c.modifiedMovePower(float64(moveBP))
	if err != nil {
		return 0.0, err
	}

	dmg := math.Floor(math.Floor((math.Floor(float64(2*c.attacker.level)/5+2)*movePower*float64(attackStat))/50) / float64(defenseStat))
	return float64(dmg), nil
}

func (c *Calculation) calcConditionalDamage(baseDmg float64) (float64, error) {
	burnMod := 1.0
	screenMod := 1.0
	weatherMod := 1.0
	flashFireMod := 1.0

	// TODO: implement conditions
	targetsMod := 1.0

	if c.attacker.status == Burned && c.attacker.usedMove.Category == string(Physical) && c.attacker.ability.Id != "guts" {
		burnMod = 0.5
	}

	// TODO: implement double battles as well => reduced screen mod
	// check light screen or reflect
	if (c.attacker.usedMove.Category == string(Physical) && c.field.Right.Reflect) ||
		(c.attacker.usedMove.Category == string(Special) && c.field.Right.LightScreen) {
		screenMod = 0.5
	}

	// check weather
	weatherMod = c.weatherModifier()

	if c.attacker.hasAbility("flashfire") && c.attacker.ability.Active && c.attacker.usedMove.HasType("fire") {
		flashFireMod = 1.5
	}

	dmg := math.Floor(baseDmg*burnMod*screenMod*targetsMod*weatherMod*flashFireMod) + 2
	return dmg, nil
}

func (c *Calculation) calcModificationDamage(condDmg float64) ([]int, error) {
	resultDmg := condDmg
	resultDmgSlice := make([]int, RESULT_DAMAGE_SLICE_SIZE)

	// modifiers which get applied before random multiplier
	critMod := 1.0
	itemMod := 1.0
	firstMod := 1.0

	// modifiers which get applied after random multiplier
	stabMod := 1.0
	typeMod1 := 1.0
	typeMod2 := 1.0
	solidRockFilterMod := 1.0
	tintedLensMod := 1.0
	expertBeltMod := 1.0
	berryMod := 1.0

	// critical modifier if ability sniper
	if c.attacker.usedMove.Critical {
		if c.attacker.hasAbility("sniper") {
			critMod = 3
		} else {
			critMod = 2
		}
	}

	resultDmg = math.Floor(resultDmg * critMod)

	// final attacker item modifier (does include expert belt mod)
	itemAtkMod, err := c.attacker.itemModifier(nil, models.DamageEffect, &models.ItemModifierOptions{
		SubStage: models.DamagePreLoop,
	})
	if err != nil {
		return resultDmgSlice, err
	}

	// final defender item modifier (does include berry mod)
	itemDefMod, err := c.defender.itemModifier(nil, models.DamageEffect, &models.ItemModifierOptions{
		MoveType: c.attacker.usedMove.Type,
		TypeMod:  c.typeEff.activeValue,
		SubStage: models.DamagePreLoop,
	})
	if err != nil {
		return resultDmgSlice, err
	}

	// final item modifier
	itemMod = itemAtkMod * itemDefMod
	resultDmg = math.Floor(resultDmg * itemMod)

	// me first mod => technically useless because the other move has to be known => for future use
	if c.attacker.usedMove.ID == "mefirst" {
		firstMod = 1.5
	}

	resultDmg = math.Floor(resultDmg * firstMod)

	// calculate loop modifiers (stab, typeMod1, typeMod2, filterMod, tintedLensMod)
	stabMod, err = c.attacker.stabModifier()
	if err != nil {
		return resultDmgSlice, err
	}

	typeMod1 = c.typeEff.typeMod1
	typeMod2 = c.typeEff.typeMod2

	// solid rock filter mod
	if (c.defender.hasAbility("filter", "solidrock")) && c.typeEff.activeValue > 1 {
		solidRockFilterMod = 0.75
	}

	// tinted lens mod
	if c.attacker.hasAbility("tintedlens") && c.typeEff.activeValue < 1 {
		tintedLensMod = 2.0
	}

	// expert belt mod
	expertBeltMod, err = c.attacker.itemModifier(nil, models.DamageEffect, &models.ItemModifierOptions{
		SubStage: models.DamageExpertBelt,
		TypeMod:  c.typeEff.activeValue,
	})

	// berry mod
	berryMod, err = c.defender.itemModifier(nil, models.DamageEffect, &models.ItemModifierOptions{
		SubStage: models.DamageBerry,
		MoveType: c.attacker.usedMove.Type,
		TypeMod:  c.typeEff.activeValue,
	})

	// add random multiplier and other modifiers inside loop to get correct damage spread with floor rounding
	for i := range resultDmgSlice {
		dmg := resultDmg

		// random multiplier
		dmg = math.Floor((dmg * float64(85+i)) / 100)

		// stab modifier
		dmg = math.Floor(dmg * stabMod)

		// type effectiveness modifier
		dmg = math.Floor(dmg * typeMod1)
		dmg = math.Floor(dmg * typeMod2)

		// filter / solid rock modifier
		dmg = math.Floor(dmg * solidRockFilterMod)

		// expert belt modifier
		dmg = math.Floor(dmg * expertBeltMod)

		// tinted lens modifier
		dmg = math.Floor(dmg * tintedLensMod)

		// berry modifier
		dmg = math.Floor(dmg * berryMod)

		resultDmgSlice[i] = int(math.Max(1, dmg))
	}

	return resultDmgSlice, nil
}
