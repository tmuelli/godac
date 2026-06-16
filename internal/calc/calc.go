package calc

import (
	"fmt"
	"math"

	"github.com/tmuelli/godac/internal/models"
)

type Calculation struct {
	data	*models.Data
	Atk		*models.BattleMon
	Def		*models.BattleMon
}

func NewCalc(data *models.Data) *Calculation {
	return &Calculation{
		data: data,
		Atk: nil,
		Def: nil,
	}
}

func (c *Calculation) Calc() (models.BattleResult, error) {
	// calculate base damage
	baseDmg, err := c.calcBaseDamage()
	if err != nil {
		return models.BattleResult{}, nil
	}

	// calculate conditional damage
	condDmg, err := c.calcConditionalDamage(baseDmg)
	if err != nil {
		return models.BattleResult{}, nil
	}

	// calcute modification damage
	minDmg, maxDmg, err := c.calcModificationDamage(condDmg)
	if err != nil {
		return models.BattleResult{}, nil
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
	// TODO: use the chosen move from the user instead of the first one
	if len(c.Atk.Moves) == 0 {
		return 0, fmt.Errorf("No moves to attack with")
	}
	
	// get attack and Defense stat based on category
	var attackStat int
	var DefenseStat int
	
	// TODO: use move type instead
	AtkMove := c.Atk.Moves[0]
	switch (AtkMove.Category) {
		case Physical:
			// TODO: use embedded types instead of untyped structs
			// TODO: use consistent base stat constants => "atk" vs "Atk"
			attackStat = getEffectiveStat("atk", *c.Atk)
			DefenseStat = getEffectiveStat("def", *c.Def)
			fmt.Println("Effectice Atk stat:", attackStat)
			fmt.Println("Effective Def stat:", DefenseStat)
		case Special:
			attackStat = getEffectiveStat("spa", *c.Atk)
			DefenseStat = getEffectiveStat("spd", *c.Def)
			fmt.Println("Effectice sp Atk stat:", attackStat)
			fmt.Println("Effective sp Def stat:", DefenseStat)
		default:
			return 0, fmt.Errorf("Invalid attack move category")
	}
	
	dmg := ((((2 * float64(c.Atk.Level)) / 5) + 2) * float64(AtkMove.BasePower) * (float64(attackStat) / float64(DefenseStat))) / 50
	fmt.Println("Base damage:", dmg)
	return float64(dmg), nil
}

func (c *Calculation) calcConditionalDamage(baseDmg float64) (float64, error) {
	burnMod := 1.0
	screenMod := 1.0
	targetsMod := 1.0
	weatherMod := 1.0
	flashFireMod := 1.0
	
	// TODO: implement conditions
	fmt.Println("Conditions not implemented yet")
	dmg := (baseDmg * burnMod * screenMod * targetsMod * weatherMod * flashFireMod) + 2
	fmt.Println("Conditional damage:", dmg)
	return dmg, nil
}

func (c *Calculation) calcModificationDamage(condDmg float64) (int, int, error) {
	critMod := 1.0
	itemMod := 1.0
	firstMod := 1.0
	solidRockFilterMod := 1.0
	expertBeltMod := 1.0
	tintedLensMod := 1.0
	berryMod := 1.0
	stabMod := 1.0
	type1Mod := 1.0
	type2Mod := 1.0

	// TODO: implement modifications
	fmt.Println("Modifications not implemented yet")

	// TODO: use selected move by user instead of first move
	AtkMove := c.Atk.Moves[0]
	if AtkMove.Type.Name == c.Atk.Types[0].Name {
		stabMod = 1.5
	}

	type1Mod = getTypeDamageModifier(AtkMove.Type.Name, c.Def.Types[0].Name)


	dmgMin := condDmg * critMod * itemMod * firstMod * 0.85 * stabMod * type1Mod * type2Mod * solidRockFilterMod * expertBeltMod * tintedLensMod * berryMod
	dmgMax:= condDmg * critMod * itemMod * firstMod * 1.0 * stabMod * type1Mod * type2Mod * solidRockFilterMod * expertBeltMod * tintedLensMod * berryMod

	fmt.Println("Min damage:", dmgMin)
	fmt.Println("Max damage:", dmgMax)

	return int(math.Floor(dmgMin)), int(math.Floor(dmgMax)), nil
}
