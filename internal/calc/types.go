package calc

import (
	"fmt"

	"github.com/tmuelli/godac/internal/models"
)

func (c *Calculation) typeDamageModifier() float64 {
	mod := 1.0

	for _, defType := range c.defender.specieTypes {
		// TODO: check if type is in effectivness matrix => map: val, ok := map
		mod *= defType.Effectiveness[c.attacker.usedMove.Type]
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