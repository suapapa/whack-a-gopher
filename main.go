package main

import (
	"github.com/0xe2-0x9a-0x9b/Go-SDL/sdl"
	"log"
	"time"
)

var (
	gopherBody     = sdl.Load("res/gopher_body_normal.png")
	gopherHulkBody = sdl.Load("res/gopher_body_hulk.png")
	gopherEye      = sdl.Load("res/gopher_doteye.png")
	gopherEyeX     = sdl.Load("res/gopher_xeye.png")
)

func main() {
	if sdl.Init(sdl.INIT_EVERYTHING) != 0 {
		log.Fatal(sdl.GetError())
	}

	bg := sdl.SetVideoMode(640, 480, 32, sdl.RESIZABLE)
	if bg == nil {
		log.Fatal(sdl.GetError())
	}

	gopherRect := &sdl.Rect{20, 20, 200, 200}

	tker := time.NewTicker(time.Second / 4) // 4fps
	var i int16
	for {
		select {
		case <-tker.C:
			gopherAnimation(bg, gopherRect, i)
			bg.Flip()
		}
		i += 1
	}
}

func gopherAnimation(bg *sdl.Surface, pos *sdl.Rect, i int16) {
	bg.FillRect(pos, 0)
	bg.Blit(pos, gopherBody, nil)

	var eye *sdl.Surface
	eyeRect := *pos
	switch {
	case i%3 == 1:
		eye = gopherEye
		eyeRect.X += 12
		eyeRect.X += 4
	case i%3 == 2:
		eye = gopherEye
		eyeRect.X -= 12
		eyeRect.X += 4
	default:
		eye = gopherEyeX
	}
	bg.Blit(&eyeRect, eye, nil)
}
