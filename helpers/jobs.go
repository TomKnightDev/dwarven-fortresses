package helpers

import (
	"github.com/sedyh/mizu/pkg/engine"
	"github.com/tomknightdev/dwarven-fortresses/components"
	"github.com/tomknightdev/dwarven-fortresses/entities"
	"github.com/tomknightdev/dwarven-fortresses/enums"
)

func GetJob(w engine.World, pos components.Position) (engine.Entity, []components.Path) {
	jobs := w.View(components.Job{}).Filter()
	var job *components.Job

	for _, j := range jobs {
		j.Get(&job)

		if job.ClaimedByID > 0 || job.Tasks[0].Blocked {
			continue
		}

		if Matches(pos, job.Tasks[0].Position) {
			return j, nil
		}

		var paths []components.Path

		if job.Tasks[0].TaskTypeEnum == enums.TaskTypePickUp || job.Tasks[0].TaskTypeEnum == enums.TaskTypeDrop {
			paths = GetPath(w, pos, job.Tasks[0].Position)
		} else {
			paths = GetPathToAdjacent(w, pos, job.Tasks[0].Position)
		}

		if len(paths) == 0 {
			job.Tasks[0].Blocked = true
			continue
		}

		return j, paths
	}

	// Create jobs for haulable items not in a stockpile
	// items := w.View(components.Item{}, components.Position{}).Filter()
	// var i *components.Item
	// var p *components.Position

	// for _, e := range items {
	// 	e.Get(&i, &p)
	// 	if !i.Claimed && i.Haulable && !i.InStockpile {
	// 		spPoses := StockpileLocations(w, i.ItemType, true)
	// 		if len(spPoses) > 0 {
	// 			j := entities.Job{
	// 				Job: components.NewJob(components.NewTask(*p, enums.TaskTypePickUp, 1), components.NewTask(spPoses[0], enums.TaskTypeAddToStockpile, 1)),
	// 			}
	// 			j.Job.EntityId = e.ID()
	// 			w.AddEntities(&j)

	// 			i.Claimed = true
	// 			break
	// 		}
	// 	}
	// }

	return nil, nil
}

func CreateJobs(w engine.World, jobType enums.TaskTypeEnum, jobCost int, positions ...components.Position) {
	for _, pos := range positions {
		if e := FindJobAtTile(w, pos, jobType); e != nil {
			continue
		}

		switch jobType {
		case enums.TaskTypeChop:
			if !TileHasResource(w, pos, enums.ResourceTypeTree) {
				continue
			}
		case enums.TaskTypeMine:
			if !IsTileOfType(w, pos, enums.TileTypeRock) {
				continue
			}
		case enums.TaskTypeBuild:

		}

		w.AddEntities(&entities.Job{
			Job: components.NewJob(components.NewTask(pos, jobType, jobCost)),
		})
	}
}

func CancelJobs(w engine.World, positions ...components.Position) {
	var entitiesToRemove []engine.Entity

	for _, p := range positions {
		if e := FindJobAtTile(w, p, enums.TaskTypeNone); e != nil {
			entitiesToRemove = append(entitiesToRemove, e)
		}
	}

	for _, job := range entitiesToRemove {
		w.RemoveEntity(job)
	}
}

func FindJobAtTile(w engine.World, pos components.Position, jobType enums.TaskTypeEnum) engine.Entity {
	jobs := w.View(components.Job{}).Filter()
	var tasks *components.Job

	for _, e := range jobs {
		e.Get(&tasks)

		for _, t := range tasks.Tasks {
			if Matches(t.Position, pos) && (t.TaskTypeEnum == jobType || jobType == enums.TaskTypeNone) {
				return e
			}
		}
	}

	return nil
}

func RescheduleJob(w engine.World, job engine.Entity) {
	var tasks *components.Job
	job.Get(&tasks)

	w.AddEntities(&entities.Job{
		Job: *tasks,
	})

	w.RemoveEntity(job)
}
