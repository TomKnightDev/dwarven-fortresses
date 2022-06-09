package components

import "github.com/tomknightdev/dwarven-fortresses/enums"

type TileType struct {
	enums.TileTypeEnum
	enums.TileSpriteEnum
}

func NewTileType(tt enums.TileTypeEnum, ts enums.TileSpriteEnum) TileType {
	return TileType{
		TileTypeEnum:   tt,
		TileSpriteEnum: ts,
	}
}
