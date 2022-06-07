package components

import (
	"github.com/OpenSauce/paths"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/tomknightdev/dwarven-fortresses/assets"
	"github.com/tomknightdev/dwarven-fortresses/enums"
)

type GameMapSingleton struct {
	WorldGenerated bool
	OffScreen      *ebiten.Image
	TilesByZ       map[int][]struct {
		Position
		TileType
		Sprite
	}
	TilesByType  map[enums.TileTypeEnum][]Position
	ResourcesByZ map[int][]struct {
		Position
		Sprite
	}
	Grids map[int]*paths.Grid
}

func NewGameMapSingleton() GameMapSingleton {
	gm := GameMapSingleton{
		Grids:     make(map[int]*paths.Grid),
		OffScreen: ebiten.NewImage(assets.WorldWidth*assets.TileSize, assets.WorldHeight*assets.TileSize),
		TilesByZ: map[int][]struct {
			Position
			TileType
			Sprite
		}{},
		TilesByType: make(map[enums.TileTypeEnum][]Position),
		ResourcesByZ: make(map[int][]struct {
			Position
			Sprite
		}),
	}

	return gm
}
