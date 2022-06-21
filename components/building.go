package components

import "github.com/tomknightdev/dwarven-fortresses/enums"

type Building struct {
	enums.TileTypeEnum
}

func NewBuilding(tt enums.TileTypeEnum) Building {
	return Building{
		TileTypeEnum: tt,
	}
}
