package main

import (
	"log"
)

func main() {
	if err := initGraphic(600, 600); err != nil {
		log.Fatal("Failed to init graphic:", err)
	}

	go runGophers(makeGophers(9))

	mouseC := make(chan Point)
	go runMouseListener(mouseC)
	for {
		p := <-mouseC
		hammerIdx := p.X/200 + (p.Y / 200 * 3)
		log.Println("Hammer to", hammerIdx)
	}
}
