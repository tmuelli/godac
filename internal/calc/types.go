package calc

import (
	"fmt"
)

func getTypeDamageModifier(atkType string, defType string) float64 {
	// TODO: implement
	fmt.Println("type damage modifier not yet implemented")

	if atkType == defType {
		return 0.5
	}

	return 1.0
}