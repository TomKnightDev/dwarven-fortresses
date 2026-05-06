package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sedyh/mizu/pkg/engine"
	"github.com/tomknightdev/dwarven-fortresses/scenes"
)

func main() {
	startPprof()

	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("DWARVEN FORTRESSES")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	g := engine.NewGame(scenes.NewMainMenu())
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
