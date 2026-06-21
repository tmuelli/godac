package calc

import (
	"fmt"

	"github.com/tmuelli/godac/internal/models"
)

func getAbilityFromSpecie(data *models.Data, s *models.Specie) (*models.Ability, error) {
	if data == nil || s == nil {
		return nil, fmt.Errorf("invalid input parameters: getAbilityFromSpecie")
	}

	abilityName, ok := s.Abilities["0"]
	if !ok {
		return nil, fmt.Errorf("invalid ability %v on specie %v", abilityName, s.Name)
	}

	var ability *models.Ability

	if a, ok := data.Abilities[s.Abilities["0"]]; ok {
		ability = &a
	}

	if ability == nil {
		return nil, fmt.Errorf("invalid ability %v on specie %v", abilityName, s.Name)
	}

	return ability, nil
}