package main

import (
	"github.com/0xe2-0x9a-0x9b/Go-SDL/sdl"
)

type Point struct {
	X, Y int
}

func runMouseListener(outC chan Point) {
	var p Point
	for {
		_event := <-sdl.Events
		switch e := _event.(type) {
		case sdl.MouseButtonEvent:
			if e.Type == sdl.MOUSEBUTTONDOWN {
				p.X, p.Y = int(e.X), int(e.Y)
				outC <- p
			}
		}
	}
}
