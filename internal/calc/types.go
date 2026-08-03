package calc

import (
	"fmt"

	"github.com/tmuelli/godac/internal/models"
)

func calcTypeModifier(moveType, defType *models.Type, canHitGhostType, isGravity bool) float64 {
	if canHitGhostType && defType.ID == "ghost" && (moveType.ID == "normal" || moveType.ID == "fighting") {
		return 1.0
	}

	if isGravity && defType.ID == "flying" && moveType.ID == "ground" {
		return 1.0
	}

	mod, ok := moveType.Effectiveness[defType.Name]
	if !ok {
		fmt.Printf("Error trying to look up type %s in type map while trying to calc type modifier\n", defType.Name)
	}

	return mod
}

func (c *Calculation) typeModifier() (float64, error) {
	// get used move type
	atkType, err := getTypeFromName(c.data, c.attacker.usedMove.Type)
	if err != nil {
		return 0.0, fmt.Errorf("Error during type damage modifier calc: %v\n", err)
	}

	canHitGhostType := c.field.Right.Foresight || c.attacker.ability.Id == "scrappy"
	defHasSecondType := len(c.defender.specieTypes) > 1
	typeMod1 := 1.0
	typeMod2 := 1.0

	typeMod1 = calcTypeModifier(atkType, c.defender.specieTypes[0], canHitGhostType, c.field.IsGravity)

	if defHasSecondType {
		typeMod2 = calcTypeModifier(atkType, c.defender.specieTypes[1], canHitGhostType, c.field.IsGravity)
	}

	typeMod := typeMod1 * typeMod2

	// handle special case => ability 'Klutz'
	if typeMod == 0.0 && atkType.ID == "ground" &&
	   c.defender.item != nil && c.defender.item.Id == "ironball" &&
	   c.defender.ability.Id != "klutz" {
		if typeMod1 == 0.0 {
			typeMod1 = 1.0
		} else if defHasSecondType && typeMod2 == 0.0 {
			typeMod2 = 1.0
		}

		typeMod = typeMod1 * typeMod2
	}

	return typeMod, nil
}

func (d *CalcData) stabModifier() (float64, error) {
	for _, t := range d.specieTypes {
		if d.usedMove.Type == t.Name {
			if d.ability.Id == "adaptability" {
				return 2.0, nil
			} else {
				return 1.5, nil
			}
		}
	}

	return 1.0, nil
}

func getTypeFromName(data *models.Data, name string) (*models.Type, error) {
	if data == nil {
		return nil, fmt.Errorf("invalid parameters: getTypeFromName")
	}

	var result *models.Type
	if t, ok := data.Types[name]; ok {
		result = &t
	}

	if result == nil {
		return nil, fmt.Errorf("invalid input: Unknown type %s", name)
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