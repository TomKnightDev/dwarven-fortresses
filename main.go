package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sedyh/mizu/pkg/engine"
	"github.com/tomknightdev/dwarven-fortresses/scenes"
)

func main() {
	go func() {
		fmt.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	ebiten.SetWindowSize(1920, 1024)
	ebiten.SetWindowTitle("DWARVEN FORTRESSES")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	g := engine.NewGame(scenes.NewMainMenu())
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
