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

const (
	// EyeX means dead eye
	EyeX EyeShape = iota
	// EyeLeft means look left
	EyeLeft
	// EyeRight means look right
	EyeRight
)

var (
	r *rand.Rand
)

func init() {
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
}

// Gopher reperesent a gopher in a hole
type Gopher struct {
	eye  EyeShape
	head chan struct{}
	sync.RWMutex
}

// NewGopher return adress of a Gopher
func NewGopher() *Gopher {
	return &Gopher{
		eye:  EyeX,
		head: make(chan struct{}),
	}
}

// Eye returns currnet shape of eye
func (g *Gopher) Eye() EyeShape {
	g.RLock()
	defer g.RUnlock()
	return g.eye
}

// Start makes a gopher run
func (g *Gopher) Start(ctx context.Context) {
	log.Printf("start gopher: %v", g)
	defer func() {
		log.Printf("end gopher: %v", g)
	}()
	for {
		select {
		case <-ctx.Done():
			return
		case <-g.head:
			log.Println("ouch!! :%v", g)
			g.Lock()
			g.eye = EyeX
			g.Unlock()
			time.Sleep(500 * time.Millisecond)
			continue
		default:
		}
		g.Lock()
		g.eye = EyeShape(1 + r.Intn(2))
		g.Unlock()
		time.Sleep(time.Duration(r.Intn(1000)) * time.Millisecond)
	}
}
