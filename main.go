package main

import (
	"log"
)

func main() {
	if err := initGraphic(600, 600); err != nil {
		log.Fatal("Failed to init graphic:", err)
	}

	runGophers(makeGophers(9))
}
