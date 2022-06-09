package systems

import (
	"log"

	"github.com/OpenSauce/paths"
	"github.com/sedyh/mizu/pkg/engine"
	"github.com/tomknightdev/dwarven-fortresses/assets"
	"github.com/tomknightdev/dwarven-fortresses/components"
	"github.com/tomknightdev/dwarven-fortresses/entities"
	"github.com/tomknightdev/dwarven-fortresses/enums"
	"github.com/tomknightdev/dwarven-fortresses/helpers"
)

type Pathfinder struct {
}

func NewPathfinder() *Pathfinder {
	return &Pathfinder{}
}

func (p *Pathfinder) Update(w engine.World) {
	gms, found := w.View(components.GameMapSingleton{}).Get()
	if !found {
		panic("game map singleton not found")
	}
	var gmComp *components.GameMapSingleton
	gms.Get(&gmComp)

	// pathCalcsForTick := 0

	view := w.View(components.Move{}, components.Position{}, components.Inventory{})
	view.Each((func(e engine.Entity) {
		var move *components.Move
		var pos *components.Position
		var inv *components.Inventory
		e.Get(&move, &pos, &inv)

		if move.CurrentEnergy < move.TotalEnergy+inv.Weight {
			move.CurrentEnergy++
			return
		}

		if move.GettingRoute {
			return
		}

		if (move.Adjacent && helpers.IsAdjacent(*move, *pos)) || helpers.Matches(components.NewPosition(move.X, move.Y, move.Z), *pos) {
			if len(move.CurrentPaths) > 1 {
				move.CurrentPaths = move.CurrentPaths[1:]
				pos.Z = move.CurrentPaths[0].Level

			} else {
				move.CurrentPaths = nil
				move.Arrived = true
			}
		} else {
			move.Arrived = false

			if move.CurrentPaths == nil {
				// If there is a lot of entities, having too many looking for paths in one game tick will really hit TPS.  Limit to 10
				// if pathCalcsForTick == 9 {
				// 	return
				// }
				// pathCalcsForTick++
				move.GettingRoute = true

				var paths []components.Path

				if move.Adjacent {
					adjacents := helpers.GetAdjacents(gmComp.Grids[move.Z], components.NewPosition(move.X, move.Y, move.Z), true)

					for _, a := range adjacents {
						paths = p.GetPath(*pos, a, gmComp.Grids, gmComp.TilesByType)
						if len(paths) > 0 {
							if len(paths[0].Cells) == 0 {
								log.Println("path with no cells found")
								continue
							}

							break
						}
					}

					// Path to job not found
					if len(paths) == 0 {
						var wk *components.Worker
						e.Get(&wk)
						wk.HasJob = false

						job, found := w.GetEntity(wk.JobID)
						if found {
							var jobComp *components.Job
							job.Get(&jobComp)
							if jobComp == nil {
								log.Println("failed to find job component")
							} else {
								jobComp.ClaimedByID = 0

								w.AddEntities(&entities.Job{
									Job: *jobComp,
								})

								w.RemoveEntity(job)
							}
						}
					}

				} else {
					paths = p.GetPath(*pos, components.NewPosition(move.X, move.Y, move.Z), gmComp.Grids, gmComp.TilesByType)
				}

				if len(paths) == 0 {
					// move.Arrived = true
					move.GettingRoute = false

					return
				}
				move.CurrentPaths = paths
			}

			if move.CurrentPaths[0].AtEnd() {
				if len(move.CurrentPaths) > 1 {
					move.CurrentPaths = move.CurrentPaths[1:]
					pos.Z = move.CurrentPaths[0].Level

				} else {
					move.CurrentPaths = nil
					move.Arrived = true
				}
			} else {
				c := move.CurrentPaths[0].Next()

				if c == nil {
					panic("no next cell in path")
				}

				pos.X = c.X
				pos.Y = c.Y
				move.CurrentPaths[0].Advance()

				for _, c := range move.CurrentPaths[0].Cells[move.CurrentPaths[0].CurrentIndex:] {
					if !c.Walkable {
						move.CurrentPaths = nil
						break
					}
				}
			}

			move.CurrentEnergy = 0
			move.GettingRoute = false
		}
	}))
}

func (p Pathfinder) GetPath(startPos components.Position, endPos components.Position, grids map[int]*paths.Grid, tileByType map[enums.TileSpriteEnum][]components.Position) []components.Path {
	if startPos.Z == endPos.Z {
		path := grids[endPos.Z].GetPath(float64(startPos.X*assets.TileSize), float64(startPos.Y*assets.TileSize), float64(endPos.X*assets.TileSize), float64(endPos.Y*assets.TileSize), true, false)

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

	travTiles[false] = append(travTiles[false], tileByType[enums.TileSpriteStairUp]...)
	travTiles[true] = append(travTiles[true], tileByType[enums.TileSpriteStairDown]...)

	direction := true
	if endPos.Z > startPos.Z {
		direction = false
	}

	var paths []components.Path
	var reached bool
	// Find route to each stair from current location, checking if each can get to final destination
	for _, tt := range travTiles[direction] {
		path := grids[startPos.Z].GetPath(float64(startPos.X*assets.TileSize), float64(startPos.Y*assets.TileSize), float64(tt.X*assets.TileSize), float64(tt.Y*assets.TileSize), true, false)

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

		paths, reached = p.traverseLevel(paths, travTiles[direction], direction, components.NewPosition(tt.X, tt.Y, tt.Z+zChange), endPos, grids)
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

func (p Pathfinder) traverseLevel(paths []components.Path, travTiles []components.Position, direction bool, startPos, finalDest components.Position, grids map[int]*paths.Grid) ([]components.Path, bool) {
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

				paths, reachedFinalDest := p.traverseLevel(paths, travTiles, direction, components.NewPosition(tt.X, tt.Y, tt.Z+zChange), finalDest, grids)
				if reachedFinalDest {
					return paths, reachedFinalDest
				}
			}
		}
	}

	return nil, false
}
