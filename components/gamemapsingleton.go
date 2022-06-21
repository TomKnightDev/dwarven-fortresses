package components

import (
	"github.com/OpenSauce/paths"
)

type GameMapSingleton struct {
	WorldGenerated bool
	Grids          map[int]*paths.Grid
	Tiles          map[Position]*Tile
	Ups            []Position
	Downs          []Position
}

func NewGameMapSingleton() GameMapSingleton {
	gm := GameMapSingleton{
		Grids: make(map[int]*paths.Grid),
		Tiles: make(map[Position]*Tile),
	}

	return gm
}
