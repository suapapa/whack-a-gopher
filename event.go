// Copyright 2013, Homin Lee. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/0xe2-0x9a-0x9b/Go-SDL/sdl"
	"os"
)

type Point struct {
	X, Y uint
}

func eventLoop(mouseC chan Point) {
	var p Point

EVENT_LOOP:
	_event := <-sdl.Events
	switch e := _event.(type) {
	case sdl.MouseButtonEvent:
		if e.Type == sdl.MOUSEBUTTONDOWN {
			p.X, p.Y = uint(e.X), uint(e.Y)
			mouseC <- p
		}
	case sdl.KeyboardEvent:
		if e.State == 0 {
			break
		}

		keysym := e.Keysym.Sym
		if keysym == sdl.K_q {
			os.Exit(0)
		}
	case sdl.QuitEvent:
		os.Exit(0)
	}

	goto EVENT_LOOP
}
