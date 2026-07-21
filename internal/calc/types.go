package calc

import (
	"fmt"

	"github.com/tmuelli/godac/internal/models"
)

func (c *Calculation) typeDamageModifier() float64 {
	mod := 1.0

	// get used move type
	atkType, err := getTypeFromName(c.data, c.attacker.usedMove.Type)
	if err != nil {
		fmt.Printf("Error during type damage modifier calc: %v\n", err)
		return 1.0
	}

	foundGhost := false
	foundFlying := false

	for _, defType := range c.defender.specieTypes {
		if defType.ID == "ghost" {
			foundGhost = true
		}

		if defType.ID == "flying" {
			foundFlying = true
		}

		// TODO: check if type is in effectivness matrix => map: val, ok := map
		mod *= atkType.Effectiveness[defType.Name]
	}

	// check foresight
	if mod == 0 && c.field.Right.Foresight {
		if (atkType.ID == "normal" || atkType.ID == "fighting") && foundGhost {
			mod = 1.0
		}
	}

	// check gravity
	if mod == 0 && c.field.Terrain == models.Gravity && foundFlying {
		if atkType.ID == "ground" {
			mod = 1.0
		}
	}

	return mod
}

func (d *CalcData) stabModifier() (float64, error) {
	for _, t := range d.specieTypes {
		if d.usedMove.Type == t.Name {
			return 1.5, nil
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