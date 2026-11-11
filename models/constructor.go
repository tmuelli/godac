package models

import "github.com/tmuelli/godac/internal/utils"

func NewData() (*Data, error) {
	abilities, err := loadAbilities()
	if err != nil {
		return nil, err
	}

	items, err := loadItems()
	if err != nil {
		return nil, err
	}

	moves, err := loadMoves()
	if err != nil {
		return nil, err
	}

	natures, err := loadNatures()
	if err != nil {
		return nil, err
	}

	species, err := loadSpecies()
	if err != nil {
		return nil, err
	}

	types, err := loadTypes()
	if err != nil {
		return nil, err
	}

	data := Data{
		Abilities: utils.SliceToMap(abilities, Ability.key),
		Items:     utils.SliceToMap(items, Item.key),
		Moves:     utils.SliceToMap(moves, Move.key),
		Natures:   utils.SliceToMap(natures, Nature.key),
		Species:   utils.SliceToMap(species, Specie.key),
		Types:     utils.SliceToMap(types, Type.key),
	}
	return &data, nil
}
