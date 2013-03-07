// Copyright 2013, Homin Lee. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

const (
	DFLT_SCRN_W   = 600
	DFLT_SCRN_H   = 600
	DFLT_FULLSCRN = false
	DFLT_FPS      = 30
)

var (
	opts  *Options
	flags []string
)

func init() {
	opts, flags = parseFlags()
}

type Options struct {
	scrnW, scrnH uint
	fullscreen   bool
	maxFps       uint
}

func setupFlags(o *Options) *flag.FlagSet {
	prgName := os.Args[0]
	fs := flag.NewFlagSet(prgName, flag.ExitOnError)

	fs.UintVar(&o.scrnW, "w", DFLT_SCRN_W, "default screen width")
	fs.UintVar(&o.scrnH, "h", DFLT_SCRN_H, "default screen height")
	fs.BoolVar(&o.fullscreen, "f", DFLT_FULLSCRN, "fullscreen")
	fs.UintVar(&o.maxFps, "fps", DFLT_FPS, "max fps")

	fs.Usage = func() {
		fmt.Printf("Usage: %s [options]\n", prgName)
		fs.PrintDefaults()
	}

	return fs
}

func verifyFlags(o *Options, fs *flag.FlagSet) {
	if o.scrnW%GOPHER_W != 0 {
		o.scrnW = GOPHER_W * (o.scrnW / GOPHER_W)
		log.Println("forced screen width to", o.scrnW)
	}

	if o.scrnH%GOPHER_H != 0 {
		o.scrnH = GOPHER_H * (o.scrnH / GOPHER_H)
		log.Println("forced screen height to", o.scrnH)
	}
}

func parseFlags() (*Options, []string) {
	var o Options
	fs := setupFlags(&o)
	fs.Parse(os.Args[1:])
	verifyFlags(&o, fs)
	return &o, fs.Args()
}
