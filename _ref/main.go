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
	canvas := doc.Call("getElementById", "canvas")
	// bodyW := doc.Get("body").Get("clientWidth").Float()
	// bodyH := doc.Get("body").Get("clientHeight").Float()
	bodyW := js.Global().Get("innerWidth").Float() * 0.9
	bodyH := js.Global().Get("innerHeight").Float() * 0.9
	log.Println(bodyW, bodyH)
	canvas.Set("width", bodyW)
	canvas.Set("height", bodyH)
	ctx := canvas.Call("getContext", "2d")

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

	// canvas.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
	// 	js.Global().Get("window").Call("alert", "Don't click me!")
	// 	return nil
	// }))

	n := 0
	js.Global().Call("setInterval", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		drawGopher(ctx, 0, 0)
		n = 1 - n
		return nil
	}), js.ValueOf(50))

	select {}
}

func drawGopher(cctx js.Value, x, y int) {
	cctx.Call("clearRect", x, y, 200, 200)
	cctx.Call("drawImage", images["gopher_body_normal"], x, y)
	if r.Intn(2) != 0 {
		cctx.Call("drawImage", images["gopher_doteye"], x-10, y)
	} else {
		cctx.Call("drawImage", images["gopher_doteye"], x+10, y)
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
