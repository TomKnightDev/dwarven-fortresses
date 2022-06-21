package components

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/tomknightdev/dwarven-fortresses/enums"
)

type Tile struct {
	Position
	enums.TileTypeEnum
	*ebiten.Image
	Resources []Resource
	Buildings []Building
	Items     []Item
}

func NewTile(pos Position, tt enums.TileTypeEnum) *Tile {
	return &Tile{
		Position:     pos,
		TileTypeEnum: tt,
	}
}
