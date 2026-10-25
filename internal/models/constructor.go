package models

import "github.com/tmuelli/godac/internal/utils"

func NewData() (*Data, error) {
	abilities, err := LoadAbilities()
	if err != nil {
		return nil, err
	}

	items, err := LoadItems()
	if err != nil {
		return nil, err
	}

	moves, err := LoadMoves()
	if err != nil {
		return nil, err
	}

	natures, err := LoadNatures()
	if err != nil {
		return nil, err
	}

	species, err := LoadSpecies()
	if err != nil {
		return nil, err
	}

	types, err := LoadTypes()
	if err != nil {
		return nil, err
	}

	data := Data{
		Abilities: utils.SliceToMap(abilities, Ability.Key),
		Items:     utils.SliceToMap(items, Item.Key),
		Moves:     utils.SliceToMap(moves, Move.Key),
		Natures:   utils.SliceToMap(natures, Nature.Key),
		Species:   utils.SliceToMap(species, Specie.Key),
		Types:     utils.SliceToMap(types, Type.Key),
	}
	return &data, nil
}
