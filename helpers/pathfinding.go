package helpers

import (
	"log"

	"github.com/OpenSauce/paths"
	"github.com/sedyh/mizu/pkg/engine"
	"github.com/tomknightdev/dwarven-fortresses/assets"
	"github.com/tomknightdev/dwarven-fortresses/components"
)

func GetPathToAdjacent(w engine.World, startPos components.Position, endPos components.Position) []components.Path {
	gms := GetGameMapSingleton(w)
	adjacents := GetAdjacents(gms.Grids[endPos.Z], endPos, true)
	adjacents = append(adjacents, endPos)

	for _, a := range adjacents {
		paths := GetPath(w, startPos, a)
		if len(paths) > 0 {
			if len(paths[0].Cells) == 0 {
				log.Println("path with no cells found")
				continue
			}

			return paths
		}
	}

	return nil
}

func GetPath(w engine.World, startPos components.Position, endPos components.Position) []components.Path {
	gms := GetGameMapSingleton(w)
	if startPos.Z == endPos.Z {
		path := gms.Grids[endPos.Z].GetPath(float64(startPos.X*assets.TileSize), float64(startPos.Y*assets.TileSize), float64(endPos.X*assets.TileSize), float64(endPos.Y*assets.TileSize), true, false)

		if path != nil && len(path.Cells) > 0 {
			return []components.Path{
				{
					Path:  path,
					Level: endPos.Z,
				},
			}
		}
	}

	// Use golden path approach, assume that the dwarf can reach the end pos from any stair
	// true = down, false = up
	travTiles := make(map[bool][]components.Position)

	travTiles[false] = gms.Ups
	travTiles[true] = gms.Downs

	direction := endPos.Z <= startPos.Z

	var paths []components.Path
	var reached bool
	// Find route to each stair from current location, checking if each can get to final destination
	for _, tt := range travTiles[direction] {
		path := gms.Grids[startPos.Z].GetPath(float64(startPos.X*assets.TileSize), float64(startPos.Y*assets.TileSize), float64(tt.X*assets.TileSize), float64(tt.Y*assets.TileSize), true, false)

		if path == nil {
			continue
		}

		paths = append(paths, components.Path{
			Path:  path,
			Level: startPos.Z,
		})

		zChange := 1
		if direction {
			zChange = -1
		}

		paths, reached = traverseLevel(paths, travTiles[direction], direction, components.NewPosition(tt.X, tt.Y, tt.Z+zChange), endPos, gms.Grids)
		if reached {
			return paths
		}
	}

	if !reached {
		log.Println("unable to reach final destination")
		return nil
	}
	return paths
}

func traverseLevel(paths []components.Path, travTiles []components.Position, direction bool, startPos, finalDest components.Position, grids map[int]*paths.Grid) ([]components.Path, bool) {
	// First thing to try is if we can reach final destination
	if startPos.Z == finalDest.Z {
		path := grids[finalDest.Z].GetPath(float64(startPos.X*assets.TileSize), float64(startPos.Y*assets.TileSize), float64(finalDest.X*assets.TileSize), float64(finalDest.Y*assets.TileSize), true, false)

		if path != nil {
			paths = append(paths, components.Path{
				Path:  path,
				Level: startPos.Z,
			})
			return paths, true
		}
	}

	for _, tt := range travTiles {
		if tt.Z == startPos.Z {
			path := grids[startPos.Z].GetPath(float64(startPos.X*assets.TileSize), float64(startPos.Y*assets.TileSize), float64(tt.X*assets.TileSize), float64(tt.Y*assets.TileSize), true, false)

			if path != nil {
				paths = append(paths, components.Path{
					Path:  path,
					Level: tt.Z,
				})

				// Final destination not found, next start pos will be the opposite stair in the same location
				zChange := 1
				if direction {
					zChange = -1
				}

				paths, reachedFinalDest := traverseLevel(paths, travTiles, direction, components.NewPosition(tt.X, tt.Y, tt.Z+zChange), finalDest, grids)
				if reachedFinalDest {
					return paths, reachedFinalDest
				}
			}
		}
	}

	return nil, false
}
