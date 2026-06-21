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
}

type Calculation struct {
	data		*models.Data
	attacker	*CalcData
	defender	*CalcData
}

func NewCalc(data *models.Data, atkInput *models.UserInput, defInput *models.UserInput) (*Calculation, error) {
	// validate inputs
	attacker, err := validateInput(data, atkInput)
	if err != nil {
		return nil, fmt.Errorf("invalid attack input: %v", err)
	}

	defender, err := validateInput(data, defInput)
	if err != nil {
		return nil, fmt.Errorf("invalid defender input: %v", err)
	}
	
	return &Calculation{
		data: data,
		attacker: attacker,
		defender: defender,
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
	
	atkCategory := c.attacker.usedMove.Category
	switch (MoveCategory(atkCategory)) {
		case Physical:
			attackStat = c.attacker.effectiveStat(StatAtk)
			DefenseStat = c.defender.effectiveStat(StatDef)
			fmt.Println("[DEBUG]: Effectice Atk stat:", attackStat)
			fmt.Println("[DEBUG]: Effective Def stat:", DefenseStat)
		case Special:
			attackStat = c.attacker.effectiveStat(StatSpa)
			DefenseStat = c.defender.effectiveStat(StatSpd)
			fmt.Println("[DEBUG]: Effectice sp Atk stat:", attackStat)
			fmt.Println("[DEBUG]: Effective sp Def stat:", DefenseStat)
		default:
			return 0, fmt.Errorf("invalid attack move category")
	}
	
	dmg := ((((2 * float64(c.attacker.level)) / 5) + 2) * float64(c.attacker.usedMove.BasePower) * (float64(attackStat) / float64(DefenseStat))) / 50
	fmt.Println("[DEBUG]: Base damage:", dmg)
	return float64(dmg), nil
}

func (c *Calculation) calcConditionalDamage(baseDmg float64) (float64, error) {
	burnMod := 1.0
	screenMod := 1.0
	targetsMod := 1.0
	weatherMod := 1.0
	flashFireMod := 1.0
	
	// TODO: implement conditions
	dmg := (baseDmg * burnMod * screenMod * targetsMod * weatherMod * flashFireMod) + 2
	fmt.Println("[DEBUG]: Conditional damage:", dmg)
	return dmg, nil
}

func (c *Calculation) calcModificationDamage(condDmg float64) (int, int, error) {
	// TODO: implement modifications
	critMod := 1.0
	itemMod := 1.0
	firstMod := 1.0
	solidRockFilterMod := 1.0
	expertBeltMod := 1.0
	tintedLensMod := 1.0
	berryMod := 1.0
	stabMod := 1.0
	typeMod := 1.0

	stabMod, err := c.attacker.stabModifier()
	if err != nil {
		return 0, 0, err
	}

	typeMod = c.typeDamageModifier()

	dmgMin := condDmg * critMod * itemMod * firstMod * 0.85 * stabMod * typeMod * solidRockFilterMod * expertBeltMod * tintedLensMod * berryMod
	dmgMax:= condDmg * critMod * itemMod * firstMod * 1.0 * stabMod * typeMod * solidRockFilterMod * expertBeltMod * tintedLensMod * berryMod

	fmt.Println("[DEBUG]: Min damage:", dmgMin)
	fmt.Println("[DEBUG]: Max damage:", dmgMax)

	return int(math.Floor(dmgMin)), int(math.Floor(dmgMax)), nil
}
