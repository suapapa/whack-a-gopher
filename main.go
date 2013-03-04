package main

import (
	"log"
)

func main() {
	if err := initGraphic(640, 480); err != nil {
		log.Fatal("Failed to init graphic:", err)
	}

	runGophers(makeGophers(2))
}
