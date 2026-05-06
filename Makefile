.PHONY: build wasm

build:
	go build -o DwarvenFortresses .

wasm:
	GOOS=js GOARCH=wasm go build -o web/game.wasm .
	cp "$$(go env GOROOT)/lib/wasm/wasm_exec.js" web/wasm_exec.js
