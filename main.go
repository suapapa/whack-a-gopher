package main

import (
	"log"
)

const (
	GAME_W   = 600
	GAME_H   = 600
	GOPHER_W = 200
	GOPHER_H = 200
)

func main() {
	if err := initGraphic(GAME_W, GAME_H); err != nil {
		log.Fatal("Failed to init graphic:", err)
	}

	gophers := makeGophers((GAME_W/GOPHER_W)*(GAME_H/GOPHER_H),
		GAME_W/GOPHER_W, GAME_H/GOPHER_H)
	go runGophers(gophers)

	mouseC := make(chan Point, 20)
	go runMouseListener(mouseC)
	for {
		p := <-mouseC
		hammerIdx := p.X/GOPHER_W + (p.Y / GOPHER_H * (GAME_W/GOPHER_W))
		log.Println("Hammer to", hammerIdx)
		gophers[hammerIdx].headC <- true
	}
}
