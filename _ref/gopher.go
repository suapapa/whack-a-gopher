package main

import (
	"context"
	"log"
	"math/rand"
	"sync"
	"time"
)

// EyePosition is position of eye of a gopher
type EyePosition int

const (
	// EyeX means dead eye
	EyeX EyePosition = iota
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
	eye EyePosition
	sync.RWMutex
}

// Eye returns currnet shape of eye
func (g *Gopher) Eye() EyePosition {
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
		default:
		}
		time.Sleep(time.Duration(r.Intn(1000)) * time.Millisecond)
		g.Lock()
		g.eye = EyePosition(r.Intn(3))
		g.Unlock()
	}
}
