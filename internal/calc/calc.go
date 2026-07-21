package calc

import (
	"fmt"
	"math"

	"github.com/tmuelli/godac/internal/models"
)

type CalcData struct {
	specie		*models.Specie
	specieTypes	[]*models.Type
	nature		*models.Nature
	ability		*models.Ability
	item		*models.Item
	usedMove	*models.Move
	ivs			*models.IvStats
	evs			*models.EvStats
	level		int
	health		int
	dataType	CalcDataType
	status		StatusCondition
}

type Calculation struct {
	data		*models.Data
	attacker	*CalcData
	defender	*CalcData
	field		*models.Field
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

	attacker.dataType = Attacker
	defender.dataType = Defender

	// TODO: decide on precedence of level input =>
	// if attacker / defender level is set or field level is set or both
	// determine on defaults and "last set" option
	// FOR NOW: only UserInput level is used

	return &Calculation{
		data: data,
		attacker: attacker,
		defender: defender,
		field: field,
	}, nil
}

func (c *Calculation) Calc() (models.BattleResult, error) {
	// calculate base damage
	baseDmg, err := c.calcBaseDamage()
	if err != nil {
		return models.BattleResult{}, err
	}

	// calculate conditional damage
	condDmg, err := c.calcConditionalDamage(baseDmg)
	if err != nil {
		return models.BattleResult{}, err
	}

	// calcute modification damage
	minDmg, maxDmg, err := c.calcModificationDamage(condDmg)
	if err != nil {
		return models.BattleResult{}, err
	}

	// TODO: implement final dot damage (e.g. stealth rock, spikes, leech seed) from field
	// in smogon damage calc this is not inside the gen4 damage calc but rather calculated afterwards
	// TODO: also implement chip damage from sand / hail here (from field weather)

	// TODO: implement percent calculation
	calcResult := models.BattleResult{
		Range: models.BattleResultRange{
			Min: minDmg,
			MinPerc: 0,
			Max: maxDmg,
			MaxPerc: 0,
		},
		Possibilities: make([]int, 0),
	}

	// populate possibilities array for result and set min and max values
	// TODO: implement possibility calculation correct
	for i := 0; i <= (calcResult.Range.Max - calcResult.Range.Min); i++ {
		calcResult.Possibilities = append(calcResult.Possibilities, calcResult.Range.Min + i)
	}

	return calcResult, nil
}

func (c *Calculation) calcBaseDamage() (float64, error) {
	// get attack and Defense stat based on category
	var attackStat int
	var DefenseStat int
	moveItemModifier := 1.0

	// check protect
	if c.field.Right.Protect && !c.attacker.usedMove.BreaksProtect {
		return 0.0, nil
	}

	// check powertrick
	if c.field.Left.PowerTrick {
		c.attacker.specie.BaseStats.Atk, c.attacker.specie.BaseStats.Def = c.attacker.specie.BaseStats.Def, c.attacker.specie.BaseStats.Atk
	}

	if c.field.Right.PowerTrick {
		c.defender.specie.BaseStats.Atk, c.defender.specie.BaseStats.Def = c.defender.specie.BaseStats.Def, c.defender.specie.BaseStats.Atk
	}

	atkCategory := c.attacker.usedMove.Category
	switch (MoveCategory(atkCategory)) {
		case Physical:
			attackStat = c.attacker.effectiveStat(StatAtk)
			DefenseStat = c.defender.effectiveStat(StatDef)
			fmt.Println("[DEBUG]: Effectice Atk stat:", attackStat)
			fmt.Println("[DEBUG]: Effective Def stat:", DefenseStat)

			kind := StatAtk
			mod, err := c.attacker.itemModifier(&kind, models.MoveEffect, nil)
			if err != nil {
				return 0, err
			}
			moveItemModifier = mod
		case Special:
			attackStat = c.attacker.effectiveStat(StatSpa)
			DefenseStat = c.defender.effectiveStat(StatSpd)
			fmt.Println("[DEBUG]: Effectice sp Atk stat:", attackStat)
			fmt.Println("[DEBUG]: Effective sp Def stat:", DefenseStat)

			kind := StatSpa
			mod, err := c.attacker.itemModifier(&kind, models.MoveEffect, nil)
			if err != nil {
				return 0, err
			}
			moveItemModifier = mod
		default:
			return 0, fmt.Errorf("invalid attack move category")
	}
	
	// TODO: add new receiver method on move to calc base power depending on different moves
	// for example Gyroball => depends on effectiveStat (Spe)
	// TODO: implement helping hand => only relevant in field type doubles => 50% move power boost
	moveBP := c.attacker.usedMove.BasePower
	movePower := float64(moveBP) * moveItemModifier

	// check field for switch out
	if c.field.Right.SwitchingOut && c.attacker.usedMove.ID == "pursuit" {
		moveBP *= 2
	}

	// TODO: change calculation to math.Floor every part
	dmg := ((((2 * float64(c.attacker.level)) / 5) + 2) * movePower * (float64(attackStat) / float64(DefenseStat))) / 50
	fmt.Println("[DEBUG]: Base damage:", dmg)
	return float64(dmg), nil
}

func (c *Calculation) calcConditionalDamage(baseDmg float64) (float64, error) {
	burnMod := 1.0
	screenMod := 1.0
	weatherMod := 1.0

	// TODO: implement conditions
	targetsMod := 1.0
	flashFireMod := 1.0

	if c.attacker.status == Burned && c.attacker.usedMove.Category == string(Physical) {
		burnMod = 0.5
	}

	// check light screen or reflect
	if (c.attacker.usedMove.Category == string(Physical) && c.field.Right.Reflect) ||
	   (c.attacker.usedMove.Category == string(Special) && c.field.Right.LightScreen) {
		screenMod = 0.5
	}

	// check weather
	weatherMod = c.weatherModifier()

	dmg := (baseDmg * burnMod * screenMod * targetsMod * weatherMod * flashFireMod) + 2
	fmt.Println("[DEBUG]: Conditional damage:", dmg)
	return dmg, nil
}

func (c *Calculation) calcModificationDamage(condDmg float64) (int, int, error) {
	stabMod := 1.0
	typeMod := 1.0

	// TODO: implement modifications
	critMod := 1.0
	itemMod := 1.0
	firstMod := 1.0
	solidRockFilterMod := 1.0
	expertBeltMod := 1.0
	tintedLensMod := 1.0
	berryMod := 1.0

	stabMod, err := c.attacker.stabModifier()
	if err != nil {
		return 0, 0, err
	}

	typeMod = c.typeDamageModifier()

	// final attacker item modifier
	itemAtkMod, err := c.attacker.itemModifier(nil, models.DamageEffect, &models.ItemModifierOptions{})
	if err != nil {
		return 0, 0, err
	}

	// final defender item modifier
	itemDefMod, err := c.defender.itemModifier(nil, models.DamageEffect, &models.ItemModifierOptions{
		MoveType:	c.attacker.usedMove.Type,
		TypeMod:	typeMod,
	})
	if err != nil {
		return 0, 0, err
	}

	// final item modifier
	itemMod = itemAtkMod * itemDefMod

	dmgMin := condDmg * critMod * itemMod * firstMod * 0.85 * stabMod * typeMod * solidRockFilterMod * expertBeltMod * tintedLensMod * berryMod
	dmgMax:= condDmg * critMod * itemMod * firstMod * 1.0 * stabMod * typeMod * solidRockFilterMod * expertBeltMod * tintedLensMod * berryMod

	fmt.Println("[DEBUG]: Min damage:", dmgMin)
	fmt.Println("[DEBUG]: Max damage:", dmgMax)

	return int(math.Floor(dmgMin)), int(math.Floor(dmgMax)), nil
}
