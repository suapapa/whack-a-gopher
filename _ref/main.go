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
	x, y int
}

var (
	images map[string]js.Value
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
		"/data/gopher_body_normal.png",
		"/data/gopher_doteye.png",
		"/data/gopher_xeye.png",
	}
	images = make(map[string]js.Value)
	for _, file := range files {
		key := filepath.Base(file)
		key = strings.TrimSuffix(key, ".png")
		images[key] = js.Global().Call("eval", "new Image()")
		images[key].Set("src", "data:image/png;base64,"+loadImage(file))
	}

	ctx := context.TODO()
	gophers := make([]*Gopher, 3) // TODO: more gophers
	gopherPos := make([]point, len(gophers))
	for i := range gophers {
		gophers[i] = new(Gopher)
		gopherPos[i] = point{200 * i, 0}
		go gophers[i].Start(ctx)
	}

	var renderFrame js.Func
	renderFrame = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		// Handle window resizing
		curBodyW := js.Global().Get("innerWidth").Float() * 0.9
		curBodyH := js.Global().Get("innerHeight").Float() * 0.9
		if curBodyW != bodyW || curBodyH != bodyH {
			bodyW, bodyH = curBodyW, curBodyH
			cvs.Set("width", bodyW)
			cvs.Set("height", bodyH)
		}
		// cvsCtx.Call("clearRect", 0, 0, bodyW, bodyH)

		drawGophers(cvsCtx, gopherPos, gophers)

		js.Global().Call("requestAnimationFrame", renderFrame)
		return nil
	})
	defer renderFrame.Release()

	log.Println("first call of renderFrame")
	js.Global().Call("requestAnimationFrame", renderFrame)

	select {}
}

func drawGophers(cvsCtx js.Value, ps []point, gs []*Gopher) {
	for i, g := range gs {
		drawGopher(cvsCtx, ps[i], g)
	}
}

func drawGopher(cvsCtx js.Value, p point, g *Gopher) {
	x, y := p.x, p.y
	cvsCtx.Call("clearRect", x, y, 200, 200)
	cvsCtx.Call("drawImage", images["gopher_body_normal"], x, y)
	switch g.Eye() {
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
