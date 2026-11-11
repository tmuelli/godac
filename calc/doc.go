// Package calc implements the Generation 4 damaga formula.
//
// A godac calculation takes two [models.UserInput] arguements as sides (attacker and defender) and a [models.UserFieldInput].
//
// Steps to use:
//
//  1. Load the game data with models.NewData.
//
//  2. Call [NewCalc] with data, attacker and defender input [models.UserInput] and a
//     [models.UserFieldInput]. The input data is validated and calculation
//     is prepared.
//
// 3. Call [Calculation.Calc] and read the damage results and the damage rolls from [Calculation.Result].
//
// Errors from calling [NewCalc] indicate an invalid input. Any error from [Calculation.Calc] indicates an internal error during calculation.
//
// Game mechanics follow Smogon's damage-calc.
package calc
