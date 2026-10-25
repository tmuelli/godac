package main

import (
	"fmt"
	"log"
	"os"

	"github.com/tmuelli/godac/internal/models"
	"github.com/tmuelli/godac/internal/tui"
)

func main() {
	// load data
	data, err := models.NewData()
	if err != nil {
		log.Fatal("Failed to load data ", err)
	}

	err = tui.Run(data)
	if err != nil {
		fmt.Printf("Error from tui: %v\n", err)
		os.Exit(1)
	}
}
