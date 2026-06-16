package calc

import (
	"fmt"
	"math"

	"github.com/tmuelli/godac/internal/models"
)

type calculation struct {
	atk models.BattleMon
	def models.BattleMon
}

func Calc(attacker models.BattleMon, defender models.BattleMon) (models.BattleResult, error) {
	calcResult := models.BattleResult{
		Range: models.BattleResultRange{
			Min: 0,
			MinPerc: 0,
			Max: 0,
			MaxPerc: 0,
		},
		Possibilities: make([]int, 0),
	}
	
	dmgCalc := calculation{
		atk: attacker,
		def: defender,
	}
	
	// calculate base damage
	baseDmg, err := dmgCalc.calcBaseDamage()
	if err != nil {
		return models.BattleResult{}, nil
	}

	// calculate conditional damage
	condDmg, err := dmgCalc.calcConditionalDamage(baseDmg)
	if err != nil {
		return models.BattleResult{}, nil
	}

	// calcute modification damage
	minDmg, maxDmg, err := dmgCalc.calcModificationDamage(condDmg)
	if err != nil {
		return models.BattleResult{}, nil
	}

	// populate possibilities array for result and set min and max values
	// TODO: implement percent calculation
	calcResult.Range.Min = minDmg
	calcResult.Range.Max = maxDmg
	// TODO: implement possibility calculation correct
	for i := 0; i <= (calcResult.Range.Max - calcResult.Range.Min); i++ {
		calcResult.Possibilities = append(calcResult.Possibilities, calcResult.Range.Min + i)
	}

	return calcResult, nil
}

func (c *calculation) calcBaseDamage() (float64, error) {
	// TODO: use the chosen move from the user instead of the first one
	if len(c.atk.Moves) == 0 {
		return 0, fmt.Errorf("No moves to attack with")
	}
	
	// get attack and defense stat based on category
	var attackStat int
	var defenseStat int
	
	// TODO: use move type instead
	atkMove := c.atk.Moves[0]
	switch (atkMove.Category) {
		case Physical:
			// TODO: use embedded types instead of untyped structs
			attackStat = getEffectiveStat("atk", c.atk)
			defenseStat = getEffectiveStat("def", c.def)
			fmt.Println("Effectice atk stat:", attackStat)
			fmt.Println("Effective def stat:", defenseStat)
		case Special:
			attackStat = getEffectiveStat("spa", c.atk)
			defenseStat = getEffectiveStat("spd", c.def)
			fmt.Println("Effectice sp atk stat:", attackStat)
			fmt.Println("Effective sp def stat:", defenseStat)
		default:
			return 0, fmt.Errorf("Invalid attack move category")
	}
	
	dmg := ((((2 * float64(c.atk.Level)) / 5) + 2) * float64(atkMove.BasePower) * (float64(attackStat) / float64(defenseStat))) / 50
	fmt.Println("Base damage:", dmg)
	return float64(dmg), nil
}

func (c *calculation) calcConditionalDamage(baseDmg float64) (float64, error) {
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

func (c *calculation) calcModificationDamage(condDmg float64) (int, int, error) {
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
	atkMove := c.atk.Moves[0]
	if atkMove.Type.Name == c.atk.Types[0].Name {
		stabMod = 1.5
	}

	type1Mod = getTypeDamageModifier(atkMove.Type.Name, c.def.Types[0].Name)


	dmgMin := condDmg * critMod * itemMod * firstMod * 0.85 * stabMod * type1Mod * type2Mod * solidRockFilterMod * expertBeltMod * tintedLensMod * berryMod
	dmgMax:= condDmg * critMod * itemMod * firstMod * 1.0 * stabMod * type1Mod * type2Mod * solidRockFilterMod * expertBeltMod * tintedLensMod * berryMod

	fmt.Println("Min damage:", dmgMin)
	fmt.Println("Max damage:", dmgMax)

	return int(math.Floor(dmgMin)), int(math.Floor(dmgMax)), nil
}
