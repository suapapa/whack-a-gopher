package main

import (
	"errors"
	"github.com/0xe2-0x9a-0x9b/Go-SDL/sdl"
	"math/rand"
	"time"
)

var (
	gopherBody     = sdl.Load("res/gopher_body_normal.png")
	gopherHulkBody = sdl.Load("res/gopher_body_hulk.png")
	gopherEye      = sdl.Load("res/gopher_doteye.png")
	gopherEyeX     = sdl.Load("res/gopher_xeye.png")

	bg *sdl.Surface
)

func initGraphic(w, h int) error {
	if sdl.Init(sdl.INIT_EVERYTHING) != 0 {
		return errors.New(sdl.GetError())
	}

	bg = sdl.SetVideoMode(w, h, 32, sdl.RESIZABLE)
	if bg == nil {
		return errors.New(sdl.GetError())
	}

	return nil
}

type Gopher struct {
	rect, rectEyeL, rectEyeR *sdl.Rect

	lastAnimTS time.Time

	hit       bool
	dizzyTill time.Time

	readyC chan int
	syncC  chan time.Time
	headC  chan bool // the hammer will come here
}

func NewGopher(x, y int16) *Gopher {
	g := &Gopher{
		rect:     &sdl.Rect{x, y, 0, 0},
		rectEyeL: &sdl.Rect{x - 10, y, 0, 0},
		rectEyeR: &sdl.Rect{x + 18, y, 0, 0},
	}

	g.readyC = make(chan int)
	g.syncC = make(chan time.Time)
	g.headC = make(chan bool)
	g.lastAnimTS = time.Unix(0, 0)
	g.dizzyTill = time.Unix(0, 0)

	return g
}

func (g *Gopher) Run() {
	for {
		select {
		case now := <-g.syncC:
			if g.dizzyTill != time.Unix(0, 0) {
				if now.After(g.dizzyTill) {
					g.dizzyTill = time.Unix(0, 0)
				} else {
					g.readyC <- 0
					continue
				}
			}

			duration := time.Duration(250+rand.Intn(250)) * time.Millisecond
			if time.Since(g.lastAnimTS) < duration {
				g.readyC <- 0
				continue
			}
			g.lastAnimTS = now

			bg.FillRect(g.rect, 0)
			bg.Blit(g.rect, gopherBody, nil)

			if g.hit {
				g.hit = false
				bg.Blit(g.rect, gopherEyeX, nil)
				g.dizzyTill = now.Add(time.Second)
				g.readyC <- 1
				continue
			}

			if rand.Intn(2) == 0 {
				// XXX: Workaround for
				// g.rectEyeL.X is recovered to 0 after Blit
				x := g.rectEyeL.X
				bg.Blit(g.rectEyeL, gopherEye, nil)
				g.rectEyeL.X = x
			} else {
				bg.Blit(g.rectEyeR, gopherEye, nil)
			}

			g.readyC <- 1

		case <-g.headC:
			g.hit = true
		}
	}
}

func makeGophers(n4W, n4H int) []*Gopher {
	gs := make([]*Gopher, n4W*n4H)
	var i, x, y int16
	for y = 0; y < int16(n4H); y++ {
		for x = 0; x < int16(n4W); x++ {
			gs[i] = NewGopher(x*GOPHER_W, y*GOPHER_H)
			i += 1
		}
	}

	return gs
}

func runGophers(gs []*Gopher) {
	now := time.Now()
	for _, g := range gs {
		go g.Run()
		g.syncC <- now
	}

	fpsTker := time.NewTicker(time.Second / 30) // 30fps
	var dirtyCnt int
	for {
		dirtyCnt = 0
		select {
		case <-fpsTker.C:
			for _, g := range gs {
				dirtyCnt += (<-g.readyC)
			}
			if dirtyCnt > 0 {
				bg.Flip()
			}

			now = time.Now()
			for _, g := range gs {
				g.syncC <- now
			}
		}
	}
}
