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

	// RegionIDs maps every walkable tile to a connected-component ID so that
	// reachability checks can be answered with a map lookup instead of an
	// A* search. Two tiles share an ID iff a worker can walk between them
	// (with stair pairs bridging adjacent Z-levels). Unwalkable / unknown
	// tiles have no entry, which RegionAt reports as 0.
	RegionIDs map[Position]int
	// RegionDirty signals that the cached RegionIDs are out of date and
	// must be rebuilt before the next consumer reads them. Set by
	// MarkRegionDirty whenever a tile's walkability changes.
	RegionDirty bool
}

func NewGameMapSingleton() GameMapSingleton {
	gm := GameMapSingleton{
		Grids:       make(map[int]*paths.Grid),
		Tiles:       make(map[Position]*Tile),
		Stockpiles:  make(map[Position]enums.ItemTypeEnum),
		RegionIDs:   make(map[Position]int),
		RegionDirty: true,
	}

	return gm
}
