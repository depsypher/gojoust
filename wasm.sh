#!/bin/zsh

env GOOS=js GOARCH=wasm go build -o ./wasm/gojoust.wasm github.com/depsypher/gojoust
cp "$(go env GOROOT)"/misc/wasm/wasm_exec.js ./wasm/

