package main

import (
	"context"
	"log"
	"math/rand"
	"sync"
	"time"
)

// EyeStatus is position of eye of a gopher
type EyeStatus int

const (
	// EyeX means dead eye
	EyeX EyeStatus = iota
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
	eye EyeStatus
	sync.RWMutex
}

// Eye returns currnet shape of eye
func (g *Gopher) Eye() EyeStatus {
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
		g.eye = EyeStatus(r.Intn(2))
		g.Unlock()
	}
}
