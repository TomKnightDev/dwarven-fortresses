package components

import (
	"github.com/OpenSauce/paths"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/tomknightdev/dwarven-fortresses/assets"
)

type GameMapSingleton struct {
	WorldGenerated bool
	OffScreen      *ebiten.Image
	Tiles          map[Position]TileInfo
	Resources      map[Position][]TileInfo
	Grids          map[int]*paths.Grid
}

func NewGameMapSingleton() GameMapSingleton {
	gm := GameMapSingleton{
		Grids:     make(map[int]*paths.Grid),
		OffScreen: ebiten.NewImage(assets.WorldWidth*assets.TileSize, assets.WorldHeight*assets.TileSize),
		Tiles:     make(map[Position]TileInfo),
		Resources: make(map[Position][]TileInfo),
	}

	return gm
}
