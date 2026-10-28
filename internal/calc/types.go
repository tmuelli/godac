package calc

import (
	"fmt"
	"slices"

	"github.com/tmuelli/godac/internal/models"
)

type TypeEffectivness struct {
	typeMod1    float64
	typeMod2    float64
	activeValue float64
}

// type precedence following smogons precedence rules
var gen4TypePrecedence = []string{
	"normal", "fire", "water", "electric", "grass", "ice", "fighting", "poison", "ground",
	"flying", "psychic", "bug", "rock", "ghost", "dragon", "dark", "steel",
}

func calcTypeModifier(moveType, defType *models.Type, canHitGhostType, isGravity bool) float64 {
	if canHitGhostType && defType.ID == "ghost" && (moveType.ID == "normal" || moveType.ID == "fighting") {
		return 1.0
	}

	if isGravity && defType.ID == "flying" && moveType.ID == "ground" {
		return 1.0
	}

	mod := moveType.Effectiveness[defType.ID]
	// if !ok {
	// 	// FUTURE FEATURES: add to logging
	// 	// fmt.Printf("Error trying to look up type %s in type map while trying to calc type modifier\n", defType.Name)
	// }

	return mod
}

func (c *Calculation) setTypeModifier() error {
	// get used move type
	atkType, err := getTypeFromID(c.data, c.attacker.usedMove.Type)
	if err != nil {
		c.typeEff = TypeEffectivness{
			typeMod1:    0.0,
			typeMod2:    0.0,
			activeValue: 0.0,
		}
		return fmt.Errorf("Error during type damage modifier calc: %v\n", err)
	}

	canHitGhostType := c.field.Right.Foresight || c.attacker.hasAbility("scrappy")

	// check type precedence rules and swap if specie types doesnt match precedence
	defTypes := slices.Clone(c.defender.specieTypes)
	if (len(defTypes) > 1) && (defTypes[0].ID != defTypes[1].ID) &&
		(slices.Index(gen4TypePrecedence, defTypes[0].ID) > slices.Index(gen4TypePrecedence, defTypes[1].ID)) {
		defTypes[0], defTypes[1] = defTypes[1], defTypes[0]
	}

	defHasSecondType := len(defTypes) > 1
	typeMod1 := 1.0
	typeMod2 := 1.0

	typeMod1 = calcTypeModifier(atkType, defTypes[0], canHitGhostType, c.field.IsGravity)

	if defHasSecondType {
		typeMod2 = calcTypeModifier(atkType, defTypes[1], canHitGhostType, c.field.IsGravity)
	}

	typeMod := typeMod1 * typeMod2

	// handle special case => ability 'Klutz'
	if typeMod == 0.0 && atkType.ID == "ground" &&
		c.defender.hasItem("ironball") && !c.defender.hasAbility("klutz") {
		if typeMod1 == 0.0 {
			typeMod1 = 1.0
		} else if defHasSecondType && typeMod2 == 0.0 {
			typeMod2 = 1.0
		}
	}

	c.typeEff = TypeEffectivness{
		typeMod1: typeMod1,
		typeMod2: typeMod2,
	}

	c.typeEff.activeValue = c.typeEff.typeMod1 * c.typeEff.typeMod2

	return nil
}

func (d *CalcData) stabModifier() (float64, error) {
	if d.hasType(d.usedMove.Type) {
		if d.hasAbility("adaptability") {
			return 2.0, nil
		} else {
			return 1.5, nil
		}
	}

	return 1.0, nil
}

func getTypeFromID(data *models.Data, id string) (*models.Type, error) {
	if data == nil {
		return nil, fmt.Errorf("invalid parameters: getTypeFromID")
	}

	var result *models.Type
	if t, ok := data.Types[id]; ok {
		result = &t
	}

	if result == nil {
		return nil, fmt.Errorf("invalid input: Unknown type %s", id)
	}

	return result, nil
}

func getTypesFromSpecie(data *models.Data, s *models.Specie) ([]*models.Type, error) {
	if data == nil || s == nil {
		return nil, fmt.Errorf("invalid parameters: getTypeFromSpecie")
	}

	var specieTypes []*models.Type
	for _, st := range s.Types {
		t, ok := data.Types[st]
		if !ok {
			return nil, fmt.Errorf("invalid type: %v", st)
		}

		specieTypes = append(specieTypes, &t)
	}

	return specieTypes, nil
}

func getTypesFromIDs(data *models.Data, ids []string) ([]*models.Type, error) {
	// check for valid type numbers
	if len(ids) < 1 || len(ids) > 2 {
		return nil, fmt.Errorf("invalid input: A specie can only have one or two types, got %d", len(ids))
	}

	// check if both types are the same
	if len(ids) == 2 && ids[0] == ids[1] {
		return nil, fmt.Errorf("invalid input: Types must differ, got %s twice", ids[0])
	}

	var types []*models.Type
	for _, id := range ids {
		t, err := getTypeFromID(data, id)
		if err != nil {
			return nil, err
		}

		types = append(types, t)
	}

	return types, nil
}

func (d *CalcData) hasType(ids ...string) bool {
	for _, t := range d.specieTypes {
		if slices.Contains(ids, t.ID) {
			return true
		}
	}

	return false
}
