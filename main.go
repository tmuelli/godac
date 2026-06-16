package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/tmuelli/godac/internal/calc"
	"github.com/tmuelli/godac/internal/models"
)

func main() {
	// load data
	data, err := models.NewData()
	if err != nil {
		log.Fatal("Failed to load data ", err)
	}

	// create new calculation
	calculation := calc.NewCalc(data)

	testAttackerData, err := os.ReadFile("test_attacker.json")
	if err != nil {
		log.Fatal("Error reading attacker test data", err)
	}

	testDefenderData, err := os.ReadFile("test_defender.json")
	if err != nil {
		log.Fatal("Error reading defender test data", err)
	}

	var attacker models.BattleMon
	var defender models.BattleMon

	if err := json.Unmarshal(testAttackerData, &attacker); err != nil {
		log.Fatal("Error while decoding attacker test data to json")
	}

	if err := json.Unmarshal(testDefenderData, &defender); err != nil {
		log.Fatal("Error while decoding defender test data to json")
	}

	fmt.Println("Attacker mon:", attacker.Name)
	fmt.Println("Defender mon:", defender.Name)
	calculation.Atk = &attacker
	calculation.Def = &defender

	fmt.Println("Start calculation...")

	var result models.BattleResult
	result, err = calculation.Calc()
	if err != nil {
		log.Fatal("Error while damage calculation")
	}

	fmt.Println("Damage result successfully calculated")
	fmt.Println("Min dmg:", result.Range.Min)
	fmt.Println("Max dmg:", result.Range.Max)
	fmt.Println("Possible range:")
	for _, v := range result.Possibilities {
		fmt.Printf("%v ", v)
	}
	fmt.Printf("\n")
}