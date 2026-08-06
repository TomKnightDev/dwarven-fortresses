package systems

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sedyh/mizu/pkg/engine"
	"github.com/tomknightdev/dwarven-fortresses/components"
	"github.com/tomknightdev/dwarven-fortresses/enums"
	"github.com/tomknightdev/dwarven-fortresses/helpers"
)

type Actor struct {
}

func NewActor() *Actor {
	return &Actor{}
}

func (a *Actor) Update(w engine.World) {
	actors := w.View(components.Worker{}, components.Move{}, components.Position{})
	actors.Each(func(e engine.Entity) {
		var worker *components.Worker
		var move *components.Move
		var pos *components.Position

		e.Get(&worker, &move, &pos)

		if worker.GettingJob {
			return
		}

		if !worker.HasJob {
			worker.GettingJob = true
			// go func() {
			jobEnt, paths := helpers.GetJob(w, *pos)

			if jobEnt == nil {
				worker.GettingJob = false
				return
			}

			var job *components.Job
			jobEnt.Get(&job)

			move.CurrentPaths = paths
			worker.HasJob = true
			worker.JobID = jobEnt.ID()

			currentTask := job.Tasks[0]

			move.Adjacent = true
			if currentTask.TaskTypeEnum == enums.TaskTypePickUp || currentTask.TaskTypeEnum == enums.TaskTypeDrop {
				move.Adjacent = false
			}

			move.X = currentTask.Position.X
			move.Y = currentTask.Position.Y
			move.Z = currentTask.Position.Z
			job.ClaimedByID = e.ID()
			worker.GettingJob = false
			// }()

			// return
		} else if move.Arrived {
			if move.CurrentEnergy < move.TotalEnergy {
				move.CurrentEnergy++
				return
			}
			jobEnt, found := w.GetEntity(worker.JobID)
			if !found {
				log.Println("unable to find job entity", worker.JobID)
				worker.HasJob = false
				worker.JobID = 0
				return
			}
			var job *components.Job
			jobEnt.Get(&job)
			if job == nil {
				log.Println("unable to find job component")
				worker.HasJob = false
				worker.JobID = 0
				return
			}
			currentTask := job.Tasks[0]

			if currentTask.ActionsComplete < currentTask.RequiredActions {
				currentTask.ActionsComplete++
				move.CurrentEnergy = 0
				return
			}

			currentTask.CompleteTask()
			if len(job.Tasks) > 1 {
				task := job.Tasks[1]
				move.Adjacent = true
				if task.TaskTypeEnum == enums.TaskTypePickUp || task.TaskTypeEnum == enums.TaskTypeDrop || task.TaskTypeEnum == enums.TaskTypeAddToStockpile {
					move.Adjacent = false
				}

				move.X = task.Position.X
				move.Y = task.Position.Y
				move.Z = task.Position.Z
				job.ClaimedByID = e.ID()
				move.Arrived = false
				return
			}

			worker.HasJob = false
			worker.JobID = 0
		} else {
			// Check for job cancellation
			_, found := w.GetEntity(worker.JobID)
			if !found {
				worker.HasJob = false
				worker.JobID = 0
				move.X = pos.X
				move.Y = pos.Y
				move.Z = pos.Z

				move.Arrived = true
			}

		}
	})
}

func (a *Actor) Draw(w engine.World, screen *ebiten.Image) {
	ents := w.View(components.Move{}, components.Position{}, components.Sprite{}, components.Inventory{})
	helpers.DrawMovingImages(w, ents)
}
