// +build js,wasm

package main

//go:generate cp $GOROOT/misc/wasm/wasm_exec.js .

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"syscall/js"
	"time"
)

var (
	images map[string]js.Value
	r      *rand.Rand
)

func init() {
	// random generator
	src := rand.NewSource(time.Now().UnixNano())
	r = rand.New(src)
}

func main() {
	doc := js.Global().Get("document")
	cvs := doc.Call("getElementById", "canvas")
	// dispW := js.Global().Get("innerWidth").Float() * 0.9
	// dispH := js.Global().Get("innerHeight").Float() * 0.9
	// log.Println(dispW, dispH)
	bodyW := js.Global().Get("innerWidth").Float() * 0.9
	bodyH := js.Global().Get("innerHeight").Float() * 0.9
	cvs.Set("width", bodyW)
	cvs.Set("height", bodyH)
	ctx := cvs.Call("getContext", "2d")

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

	n := 0
	var renderFrame js.Func
	renderFrame = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		// Handle window resizing
		curBodyW := js.Global().Get("innerWidth").Float() * 0.9
		curBodyH := js.Global().Get("innerHeight").Float() * 0.9
		if curBodyW != bodyW || curBodyH != bodyH {
			bodyW, bodyH = curBodyW, curBodyH
			cvs.Set("width", bodyW)
			cvs.Set("height", bodyH)
			log.Println("cvs WxH =", bodyW, bodyH)
		}
		ctx.Call("clearRect", 0, 0, bodyW, bodyH)

		drawGopher(ctx, n*200, 0)
		n = 1 - n

		log.Println("call of renderFrame")
		js.Global().Call("requestAnimationFrame", renderFrame)
		return nil
	})
	defer renderFrame.Release()

	log.Println("first call of renderFrame")
	js.Global().Call("requestAnimationFrame", renderFrame)

	select {}
}

func drawGopher(ctx js.Value, x, y int) {
	ctx.Call("clearRect", x, y, 200, 200)
	ctx.Call("drawImage", images["gopher_body_normal"], x, y)
	if r.Intn(2) != 0 {
		ctx.Call("drawImage", images["gopher_doteye"], x-10, y)
	} else {
		ctx.Call("drawImage", images["gopher_doteye"], x+10, y)
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
