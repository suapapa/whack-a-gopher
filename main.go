// +build js,wasm

package main

//go:generate cp $GOROOT/misc/wasm/wasm_exec.js .

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"syscall/js"
	"time"
)

type point struct {
	x, y float64
}

const (
	GOPHER_W = 200
	GOPHER_H = 200
)

var (
	images   map[string]js.Value
	mousePos point
	n4W, n4H int
)

func main() {
	doc := js.Global().Get("document")
	cvs := doc.Call("getElementById", "canvas")
	bodyW := js.Global().Get("innerWidth").Float() * 0.9
	bodyH := js.Global().Get("innerHeight").Float() * 0.9
	cvs.Set("width", bodyW)
	cvs.Set("height", bodyH)
	cvsCtx := cvs.Call("getContext", "2d")

	files := []string{
		"/res/gopher_body_normal.png",
		"/res/gopher_doteye.png",
		"/res/gopher_xeye.png",
	}
	images = make(map[string]js.Value)
	for _, file := range files {
		key := filepath.Base(file)
		key = strings.TrimSuffix(key, ".png")
		images[key] = js.Global().Call("eval", "new Image()")
		images[key].Set("src", "data:image/png;base64,"+loadImage(file))
	}

	mouseClickCh := make(chan point, 10)

	mouseClickEvt := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		e := args[0]
		mousePos.x = e.Get("clientX").Float()
		mousePos.y = e.Get("clientY").Float()
		mouseClickCh <- mousePos
		return nil
	})
	defer mouseClickEvt.Release()
	doc.Call("addEventListener", "click", mouseClickEvt)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	gophers, gopherPos := makeGophersAndPositions(ctx, bodyW, bodyH)
	go func() {
		for {
			p := <-mouseClickCh
			hammerIdx := int(p.x)/GOPHER_W + (int(p.y) / GOPHER_H * n4W)
			gophers[hammerIdx].HeadCh <- struct{}{}
		}
	}()
	go porkGophers(ctx, gophers)

	var renderFrame js.Func
	renderFrame = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		// // Handle window resizing
		// curBodyW := js.Global().Get("innerWidth").Float() * 0.9
		// curBodyH := js.Global().Get("innerHeight").Float() * 0.9
		// if curBodyW != bodyW || curBodyH != bodyH {
		// 	bodyW, bodyH = curBodyW, curBodyH
		// 	cvs.Set("width", bodyW)
		// 	cvs.Set("height", bodyH)
		// 	cancel()
		// 	ctx, cancel = context.WithCancel(context.Background())
		// 	gophers, gopherPos = makeGophersAndPositions(ctx, bodyW, bodyH)
		// }

		drawGophers(cvsCtx, gopherPos, gophers)

		js.Global().Call("requestAnimationFrame", renderFrame)
		return nil
	})
	defer renderFrame.Release()

	js.Global().Call("requestAnimationFrame", renderFrame)

	select {}
}

func makeGophersAndPositions(ctx context.Context, cvsW, cvsH float64) ([]*Gopher, []point) {
	n4W, n4H = int(cvsW)/GOPHER_W, int(cvsH)/GOPHER_H
	gophers := make([]*Gopher, n4W*n4H)
	positions := make([]point, len(gophers))
	var i, x, y int
	for y = 0; y < n4H; y++ {
		for x = 0; x < n4W; x++ {
			gophers[i] = NewGopher()
			go gophers[i].Start(ctx)
			positions[i] = point{
				x: float64(x * GOPHER_W),
				y: float64(y * GOPHER_H),
			}
			i += 1
		}
	}
	return gophers, positions
}

func porkGophers(ctx context.Context, gophers []*Gopher) {
porkLoop:
	select {
	case <-ctx.Done():
		return
	default:
	}
	for _, g := range gophers {
		if r.Intn(100) < 5 { // this set difficulty
			g.ButtCh <- struct{}{}
		}
	}
	time.Sleep(100 * time.Millisecond)
	goto porkLoop
}

func drawGophers(cvsCtx js.Value, ps []point, gs []*Gopher) {
	// log.Println("drawGophers: len = ", len(gs))
	for i, g := range gs {
		drawGopher(cvsCtx, ps[i], g)
	}
}

func drawGopher(cvsCtx js.Value, p point, g *Gopher) {
	x, y := p.x, p.y
	cvsCtx.Call("clearRect", x, y, 200, 200)
	if g.Status() == Hide {
		return
	}
	cvsCtx.Call("drawImage", images["gopher_body_normal"], x, y)
	switch g.Eye {
	case EyeLeft:
		cvsCtx.Call("drawImage", images["gopher_doteye"], x-10, y)
	case EyeRight:
		cvsCtx.Call("drawImage", images["gopher_doteye"], x+20, y)
	default:
		cvsCtx.Call("drawImage", images["gopher_xeye"], x-5, y)
	}
}

func loadImage(path string) string {
	href := js.Global().Get("location").Get("href")
	u, err := url.Parse(href.String())
	if err != nil {
		log.Fatal(err)
	}
	u.Path = path
	u.RawQuery = fmt.Sprint(time.Now().UnixNano())

	log.Println("loading image file: " + u.String())
	resp, err := http.Get(u.String())
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return base64.StdEncoding.EncodeToString(b)
}
