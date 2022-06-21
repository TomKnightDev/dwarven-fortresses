package systems

import (
	"log"

	"github.com/sedyh/mizu/pkg/engine"
	"github.com/tomknightdev/dwarven-fortresses/components"
	"github.com/tomknightdev/dwarven-fortresses/entities"
	"github.com/tomknightdev/dwarven-fortresses/helpers"
)

type Pathfinder struct {
}

func NewPathfinder() *Pathfinder {
	return &Pathfinder{}
}

func (p *Pathfinder) Update(w engine.World) {
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

			return
		}

		move.Arrived = false

		if move.CurrentPaths == nil {
			move.GettingRoute = true

			var paths []components.Path

			if move.Adjacent {
				paths = helpers.GetPathToAdjacent(w, *pos, components.NewPosition(move.X, move.Y, move.Z))
			} else {
				paths = helpers.GetPath(w, *pos, components.NewPosition(move.X, move.Y, move.Z))
			}

			if len(paths) > 0 {
				move.CurrentPaths = paths
			} else {
				// Path to job not found
				var wk *components.Worker
				e.Get(&wk)
				wk.HasJob = false

				jobent, found := w.GetEntity(wk.JobID)
				if !found {
					log.Println("unable to find job entity", wk.JobID)
					wk.HasJob = false
					wk.JobID = 0
					return
				}
				var job *components.Job
				jobent.Get(&job)
				if job == nil {
					log.Println("unable to find job component", wk.JobID)
					wk.HasJob = false
					wk.JobID = 0
					return
				}

				job.ClaimedByID = 0
				wk.JobID = 0
				move.GettingRoute = false

				w.AddEntities(&entities.Job{
					Job: *job,
				})

				w.RemoveEntity(jobent)

				return
			}
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
	}))
}
