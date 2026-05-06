package helpers

import (
	"sort"

	"github.com/sedyh/mizu/pkg/engine"
	"github.com/tomknightdev/dwarven-fortresses/components"
	"github.com/tomknightdev/dwarven-fortresses/entities"
	"github.com/tomknightdev/dwarven-fortresses/enums"
)

// GetJob looks for an unclaimed, unblocked job that the worker at pos can
// reach, and returns the job entity together with a path to it.
//
// The search uses the cached region index to discard jobs that are
// definitely unreachable without running A* on them, then sorts the
// surviving candidates by distance and only A*-searches them in
// closest-first order. The dwarf therefore picks up the nearest reachable
// job rather than whichever one happens to come first in the ECS view.
func GetJob(w engine.World, pos components.Position) (engine.Entity, []components.Path) {
	gms := GetGameMapSingleton(w)
	EnsureRegions(gms)

	workerRegion := RegionAt(gms, pos)

	type candidate struct {
		ent  engine.Entity
		job  *components.Job
		dist int
	}

	jobs := w.View(components.Job{}).Filter()
	var job *components.Job
	var candidates []candidate

	for _, j := range jobs {
		j.Get(&job)

		if job.ClaimedByID > 0 || job.Tasks[0].Blocked {
			continue
		}

		target := job.Tasks[0].Position
		if Matches(pos, target) {
			// Standing on it — return the job with no movement required.
			return j, nil
		}

		// Region filter: skip jobs the worker provably can't reach. This
		// turns "run A* on every job to find out it's hopeless" into a
		// cheap map lookup. workerRegion == 0 means the worker itself is
		// on an unwalkable tile (shouldn't happen) — fall back to the
		// pre-region behaviour by accepting every candidate.
		if workerRegion != 0 {
			taskType := job.Tasks[0].TaskTypeEnum
			if taskType == enums.TaskTypePickUp || taskType == enums.TaskTypeDrop {
				// Worker walks onto the target tile, so the target itself
				// must be in the same region.
				if RegionAt(gms, target) != workerRegion {
					continue
				}
			} else {
				// Worker only needs to reach a neighbour (chop, mine,
				// build, add-to-stockpile).
				reachable := false
				for _, r := range AdjacentRegions(gms, target) {
					if r == workerRegion {
						reachable = true
						break
					}
				}
				if !reachable {
					continue
				}
			}
		}

		candidates = append(candidates, candidate{
			ent:  j,
			job:  job,
			dist: travelHeuristic(pos, target),
		})
	}

	sort.Slice(candidates, func(i, k int) bool {
		return candidates[i].dist < candidates[k].dist
	})

	for _, c := range candidates {
		target := c.job.Tasks[0].Position
		var paths []components.Path
		if c.job.Tasks[0].TaskTypeEnum == enums.TaskTypePickUp || c.job.Tasks[0].TaskTypeEnum == enums.TaskTypeDrop {
			paths = GetPath(w, pos, target)
		} else {
			paths = GetPathToAdjacent(w, pos, target)
		}

		if len(paths) == 0 {
			// Region said reachable but A* failed (rare — happens if the
			// region index was stale for this tick, or if there's no
			// stair path despite same-region membership). Mark the task
			// blocked; it will un-block automatically on the next
			// walkability change via MarkRegionDirty.
			c.job.Tasks[0].Blocked = true
			continue
		}

		return c.ent, paths
	}

	return nil, nil
}

// travelHeuristic estimates how far the worker has to walk to reach a job.
// It's a Chebyshev distance on XY (matches the 8-direction grid) plus a
// per-Z penalty, since changing levels requires routing through a stair.
// Only used to order candidates — not a cost the pathfinder consumes.
func travelHeuristic(a, b components.Position) int {
	dx := a.X - b.X
	if dx < 0 {
		dx = -dx
	}
	dy := a.Y - b.Y
	if dy < 0 {
		dy = -dy
	}
	dz := a.Z - b.Z
	if dz < 0 {
		dz = -dz
	}
	cheb := dx
	if dy > cheb {
		cheb = dy
	}
	return cheb + 6*dz
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
