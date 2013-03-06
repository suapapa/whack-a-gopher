// Copyright 2013, Homin Lee. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/0xe2-0x9a-0x9b/Go-SDL/sdl"
	"log"
)

var (
	gopherBodyFile     = "gopher_body_normal.png"
	gopherHulkBodyFile = "gopher_body_hulk.png"
	gopherEyeFile      = "gopher_doteye.png"
	gopherEyeXFile     = "gopher_xeye.png"
)

func loadImage(fn string) *sdl.Surface {
	// TODO: fix to also find image from $GOPATH.
	if s := sdl.Load("res/" + fn); s != nil {
		return s
	}

	log.Fatalf("failed to load image, %s\n", fn)
	return nil
}
