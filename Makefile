ifeq ($(OS),Windows_NT)
    BROWSER = start
else
    BROWSER = xdg-open
endif

.PHONY: all clean serve

all: main.wasm serve

%.wasm: *.go
	GOOS=js GOARCH=wasm go generate
	GOOS=js GOARCH=wasm go build -o "$@" *.go
	cp main.wasm build/
	cp wasm_exec.js build/
	cp -r whack-a-gopher build/

serve:
	$(BROWSER) 'http://localhost:5000'
	serve || (go get -v github.com/mattn/serve && serve)

deploy: main.wasm
	rm -rf _deploy
	mkdir -p _deploy/build
	cp index.html _deploy/
	cp wasm_exec.js _deploy/build
	cp main.wasm _deploy/build
	cp -r whack-a-gopher/res _deploy/
	cd _deploy && git init . && git add . && \
		git commit -m "deploy" && \
		git push -f https://github.com/suapapa/whack-a-gopher master:gh-pages

clean:
	rm -f *.wasm
