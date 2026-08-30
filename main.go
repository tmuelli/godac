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

	testAttackerData, err := os.ReadFile("test_attacker_v2.json")
	if err != nil {
		log.Fatal("Error reading attacker test data", err)
	}

	testDefenderData, err := os.ReadFile("test_defender_v2.json")
	if err != nil {
		log.Fatal("Error reading defender test data", err)
	}

	testFieldData, err := os.ReadFile("test_field.json")
	if err != nil {
		log.Fatal("Error reading field test data", err)
	}

	var attacker models.UserInput
	var defender models.UserInput
	var field models.UserFieldInput

	if err := json.Unmarshal(testAttackerData, &attacker); err != nil {
		log.Fatal("Error while decoding attacker test data to json")
	}

	if err := json.Unmarshal(testDefenderData, &defender); err != nil {
		log.Fatal("Error while decoding defender test data to json")
	}

	if err := json.Unmarshal(testFieldData, &field); err != nil {
		log.Fatal("Error while decoding field test data to json")
	}

	// create new calculation
	calculation, err := calc.NewCalc(data, &attacker, &defender, &field)
	if err != nil {
		log.Fatal("Error while initializing calculation data:", err)
	}

	fmt.Println("Start calculation...")
	err = calculation.Calc()
	if err != nil {
		log.Fatal("Error while damage calculation: ", err)
	}

	fmt.Println("Damage result successfully calculated")
	fmt.Println("Min dmg:", calculation.Result.Range.Min)
	fmt.Println("Max dmg:", calculation.Result.Range.Max)
	fmt.Println("Possible range:")
	for _, v := range calculation.Result.DmgMatrix[0] {
		fmt.Printf("%v ", v)
	}
	fmt.Printf("\n")
}