// Copyright 2013, Homin Lee. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"math/rand"
	"time"
)

const (
	GAME_W   = 600
	GAME_H   = 600
	GOPHER_W = 200
	GOPHER_H = 200
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	log.Println("opts =", opts)
	if err := initGraphic(opts.scrnW, opts.scrnH, opts.fullscreen); err != nil {
		log.Fatal("Failed to init graphic:", err)
	}

	n4W, n4H := opts.scrnW/GOPHER_W, opts.scrnH/GOPHER_H
	gophers := makeGophers(n4W, n4H)
	log.Printf("%d gophers ready\n", len(gophers))
	go runGophers(gophers)

	runPoker := func(d time.Duration) {
		pokeTkr := time.NewTicker(d)
		for {
			<-pokeTkr.C
			pokeIdx := rand.Intn(len(gophers))
			log.Println("Poke ", pokeIdx)
			gophers[pokeIdx].buttC <- true
		}
	}

	go runPoker(3 * time.Second)
	go runPoker(time.Second)
	go runPoker(time.Second / 2)

	mouseC := make(chan Point, 20)
	go runMouseListener(mouseC)
	for {
		p := <-mouseC
		hammerIdx := p.X/GOPHER_W + (p.Y / GOPHER_H * n4W)
		log.Println("Hammer to", hammerIdx)
		gophers[hammerIdx].headC <- true
	}
}
