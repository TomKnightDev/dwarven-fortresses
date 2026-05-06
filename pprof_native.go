//go:build !js

package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
)

func startPprof() {
	go func() {
		fmt.Println(http.ListenAndServe("localhost:6060", nil))
	}()
}
