// Copyright 2013, Homin Lee. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/banthar/Go-SDL/sdl"
)

var (
	pkgPath       = "github.com/suapapa/whac-a-gopher"
	resCandidates = []string{
		"res",
		filepath.Join(os.Getenv("GOPATH"), "src", pkgPath, "res"),
		filepath.Join(os.Getenv("GOROOT"), "src", "pkg", pkgPath, "res"),
	}
)

func loadImage(filename string) *sdl.Surface {
	// search res directory in following order
	// - ./res/
	// - $(GOPATH)/src/github.com/suapapa/whac-a-gopher/res/
	// - $(GOROOT)/src/pkg/github.com/suapapa/whac-a-gopher/res/
	var imagePath string
	for _, d := range resCandidates {
		tryPath := filepath.Join(d, filename)
		if _, err := os.Stat(tryPath); err == nil {
			imagePath = tryPath
			break
		}
	}

	if imagePath == "" {
		log.Fatalf("failed to find image, %s\n", imagePath)
	}

	if s := sdl.Load(imagePath); s != nil {
		log.Println(imagePath, "loaded")
		return s
	}

	log.Fatalf("failed to load image, %s\n", imagePath)
	return nil
}
