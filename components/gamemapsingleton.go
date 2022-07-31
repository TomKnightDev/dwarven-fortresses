package components

import (
	"github.com/OpenSauce/paths"
	"github.com/tomknightdev/dwarven-fortresses/enums"
)

type GameMapSingleton struct {
	WorldGenerated bool
	Grids          map[int]*paths.Grid
	Tiles          map[Position]*Tile
	Ups            []Position
	Downs          []Position
	Stockpiles     map[Position]enums.ItemTypeEnum
}

func NewGameMapSingleton() GameMapSingleton {
	gm := GameMapSingleton{
		Grids:      make(map[int]*paths.Grid),
		Tiles:      make(map[Position]*Tile),
		Stockpiles: make(map[Position]enums.ItemTypeEnum),
	}

	return gm
}
