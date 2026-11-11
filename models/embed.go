package models

import (
	_ "embed"
	"encoding/json"
)

//go:embed abilities.json
var rawAbilities []byte

//go:embed items.json
var rawItems []byte

//go:embed moves.json
var rawMoves []byte

//go:embed natures.json
var rawNatures []byte

//go:embed species.json
var rawSpecies []byte

//go:embed types.json
var rawTypes []byte

func loadAbilities() ([]Ability, error) {
	var abilities []Ability
	return abilities, json.Unmarshal(rawAbilities, &abilities)
}

func loadItems() ([]Item, error) {
	var items []Item
	return items, json.Unmarshal(rawItems, &items)
}

func loadMoves() ([]Move, error) {
	var moves []Move
	return moves, json.Unmarshal(rawMoves, &moves)
}

func loadNatures() ([]Nature, error) {
	var natures []Nature
	return natures, json.Unmarshal(rawNatures, &natures)
}

func loadSpecies() ([]Specie, error) {
	var species []Specie
	return species, json.Unmarshal(rawSpecies, &species)
}

func loadTypes() ([]Type, error) {
	var types []Type
	return types, json.Unmarshal(rawTypes, &types)
}
