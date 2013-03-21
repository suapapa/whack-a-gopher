package main

import (
	"log"
	"time"
)

var (
	updateFrameC chan bool
)

func init() {
	updateFrameC = make(chan bool)
}

func debugFpsLoop() {
	var frameCnt uint
	printFpsIntervalTkr := time.NewTicker(3 * time.Second)
DEBUG_FPS_LOOP:
	select {
	case <-updateFrameC:
		frameCnt += 1
	case <-printFpsIntervalTkr.C:
		log.Println("fps =", frameCnt/3)
		frameCnt = 0
	}
	goto DEBUG_FPS_LOOP
}
