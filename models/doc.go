// Package models provides the data and types used by the godac damage
// calculator for the 4th generation of pokemon fights.
//
// Contains game data (species, moves, abilities, items, natures and types), the input and result type of a godac calculation and various helpers.
// Data is embedded form Smogon's damage-calc data
//
// To get started, load the game data once with [NewData]. For the calculation to run you'll need [UserInput] and [UserFieldInput]. Pass them to calc.NewCalc
// together with the created data.
//
// Hint: Some exports only exists because they are needed in the calc package. They are engine internals and not a stable API and are marked as those.
package models
