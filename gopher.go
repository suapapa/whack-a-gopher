package main

import (
	"context"
	"log"
	"math/rand"
	"sync"
	"time"
)

// EyeShape is position of eye of a gopher
type EyeShape int

// GopherStatus is status of gopher
type GopherStatus int

const (
	// EyeX means dead eye
	EyeX EyeShape = iota
	// EyeLeft means look left
	EyeLeft
	// EyeRight means look right
	EyeRight

	// Hide means gopher is in the hole
	Hide GopherStatus = iota
	// Peak means gopher peaks
	Peak
	// Dizzy means gopher is dizzed by hammer
	Dizzy
)

var (
	r *rand.Rand
)

func init() {
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
}

// Gopher reperesent a gopher in a hole
type Gopher struct {
	Eye                   EyeShape
	HeadCh, ButtCh        chan struct{}
	dizzyUntil, peakUntil time.Time
	status                GopherStatus
	sync.RWMutex          // Lock for status and
}

// NewGopher return adress of a Gopher
func NewGopher() *Gopher {
	return &Gopher{
		Eye:    EyeX,
		HeadCh: make(chan struct{}, 1),
		ButtCh: make(chan struct{}, 1),
		status: Hide,
	}
}

// Status returns currnet status of gopher
func (g *Gopher) Status() GopherStatus {
	g.RLock()
	defer g.RUnlock()
	return g.status
}

// Start makes a gopher run
func (g *Gopher) Start(ctx context.Context) {
	log.Printf("start gopher: %v", g)
	defer func() {
		log.Printf("end gopher: %v", g)
	}()
	go g.updateStatus(ctx)
	for {
		select {
		case <-ctx.Done():
			return
		case <-g.HeadCh:
			if g.Status() == Hide {
				continue
			}
			log.Println("ouch!! :%v", g)
			g.Lock()
			g.status = Dizzy
			g.Unlock()
			g.Eye = EyeX
			g.dizzyUntil = time.Now().Add(500 * time.Millisecond)
			continue
		case <-g.ButtCh:
			if g.Status() != Hide {
				continue
			}
			g.Lock()
			g.status = Peak
			g.Unlock()
			g.Eye = EyeShape(1 + r.Intn(2))
			g.peakUntil = time.Now().Add(time.Duration(r.Intn(1000))*time.Millisecond + 100*time.Millisecond)
			continue
		default:
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func (g *Gopher) updateStatus(ctx context.Context) {
loop:
	select {
	case <-ctx.Done():
		return
	default:
	}

	switch g.Status() {
	case Dizzy:
		if time.Now().After(g.dizzyUntil) {
			g.Lock()
			g.status = Hide
			g.Unlock()
		}
	case Peak:
		if time.Now().After(g.peakUntil) {
			g.Lock()
			g.status = Hide
			g.Unlock()
		} else {
			g.Eye = EyeLeft + EyeShape(r.Intn(2))
		}
	default:
	}
	time.Sleep(100 * time.Millisecond)
	goto loop
}
