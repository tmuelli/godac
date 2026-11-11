package calc_test

import (
	"fmt"
	"log"

	"github.com/tmuelli/godac/calc"
	"github.com/tmuelli/godac/models"
)

func ExampleNewCalc() {
	data, err := models.NewData()
	if err != nil {
		log.Fatal(err)
	}

	poochyena := data.Species["Poochyena"]
	bite := data.Moves["Bite"]

	attacker := data.DefaultUserInput(&poochyena)
	attacker.AttackMove = bite.Name
	attacker.Moves = []string{bite.Name}
	attacker.InputSide = string(models.Attacker)

	defender := data.DefaultUserInput(&poochyena)
	defender.InputSide = string(models.Defender)

	field := models.UserFieldInput{
		Type:      string(models.Singles),
		Level:     int(models.Level100),
		Terrain:   string(models.Clear),
		IsGravity: false,
	}

	c, err := calc.NewCalc(data, &attacker, &defender, &field)
	if err != nil {
		log.Fatal(err)
	}
	if err := c.Calc(); err != nil {
		log.Fatal(err)
	}

	fmt.Println(c.Result.Range.Min, c.Result.Range.Max)
	// Output: 45 53
}
