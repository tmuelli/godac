package models

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
		Abilities:	abilities,
		Items:		items,
		Moves:		moves,
		Natures:	natures,
		Species:	species,
		Types:		types,
	}
	return &data, nil
}