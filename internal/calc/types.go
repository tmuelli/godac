package calc

import (
	"fmt"

	"github.com/tmuelli/godac/internal/models"
)

type MoveTypeComb struct {
	m *models.Move
	t []*models.Type
}

func getTypeDamageModifier(aComb MoveTypeComb, dComb MoveTypeComb) float64 {
	mod := 1.0

	for _, defType := range dComb.t {
		mod *= defType.Effectiveness[aComb.m.Type]
	}

	return mod
}

func getStabModifier(comb MoveTypeComb) (float64, error) {
	return 0.0, nil
}

func getTypesFromSpecie(data *models.Data, s *models.Specie) ([]*models.Type, error) {
	if data == nil || s == nil {
		return nil, fmt.Errorf("Invalid parameters: getTypeFromSpecie")
	}
	
	var specieTypes []*models.Type
	for _, st := range s.Types {
		for _, dt := range data.Types {
			if st == dt.Name {
				specieTypes = append(specieTypes, &dt)
			}
		}
	}

	return specieTypes, nil
}