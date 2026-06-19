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
	// TODO: implement
	fmt.Println("type damage modifier not yet implemented")

	for _, t := range dComb.t {
		if aComb.m.Type == t.Name {
			return 0.5
		}
	}

	return 1.0
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