package helpers

import (
	"github.com/sedyh/mizu/pkg/engine"
	"github.com/tomknightdev/dwarven-fortresses/assets"
	"github.com/tomknightdev/dwarven-fortresses/components"
)

// MarkRegionDirty signals that walkability has changed and the cached
// region index must be rebuilt before the next consumer reads it. It also
// clears the Blocked flag on every open task so that previously
// unreachable jobs get re-evaluated against the new map state.
//
// Call this from any code path that flips a cell's Walkable bit.
func MarkRegionDirty(w engine.World, gms *components.GameMapSingleton) {
	gms.RegionDirty = true

	jobs := w.View(components.Job{}).Filter()
	var job *components.Job
	for _, j := range jobs {
		j.Get(&job)
		for i := range job.Tasks {
			job.Tasks[i].Blocked = false
		}
	}
}

// EnsureRegions rebuilds the region index if it is dirty. Cheap to call
// every tick — it is a no-op when clean.
func EnsureRegions(gms *components.GameMapSingleton) {
	if !gms.RegionDirty && gms.RegionIDs != nil {
		return
	}
	rebuildRegions(gms)
	gms.RegionDirty = false
}

// RegionAt returns the region ID of the given position, or 0 if the tile
// is unwalkable or out of bounds. Region IDs start at 1.
func RegionAt(gms *components.GameMapSingleton, pos components.Position) int {
	if gms.RegionIDs == nil {
		return 0
	}
	return gms.RegionIDs[pos]
}

// SameRegion reports whether two positions are part of the same connected
// region. An unwalkable tile (region 0) is never equal to anything,
// including another unwalkable tile.
func SameRegion(gms *components.GameMapSingleton, a, b components.Position) bool {
	ra := RegionAt(gms, a)
	if ra == 0 {
		return false
	}
	return ra == RegionAt(gms, b)
}

// AdjacentRegions returns the distinct region IDs of the 8 walkable
// neighbours of pos. Used when a job's target tile is itself unwalkable
// (a tree, a rock wall) and the worker only has to reach one of its
// neighbours.
func AdjacentRegions(gms *components.GameMapSingleton, pos components.Position) []int {
	var ids []int
	seen := make(map[int]struct{})
	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			if dx == 0 && dy == 0 {
				continue
			}
			n := components.NewPosition(pos.X+dx, pos.Y+dy, pos.Z)
			if n.X < 0 || n.Y < 0 || n.X >= assets.Settings.MapWidth() || n.Y >= assets.Settings.MapHeight() {
				continue
			}
			id := RegionAt(gms, n)
			if id == 0 {
				continue
			}
			if _, ok := seen[id]; ok {
				continue
			}
			seen[id] = struct{}{}
			ids = append(ids, id)
		}
	}
	return ids
}

// rebuildRegions recomputes the region index from scratch by flood-filling
// every walkable tile across all Z-levels. Stair pairs (a Down at z and an
// Up at z-1 sharing X/Y) act as bidirectional bridges between levels, so
// floors connected only by stairs end up in the same region.
func rebuildRegions(gms *components.GameMapSingleton) {
	regions := make(map[components.Position]int, len(gms.Tiles))

	// Build the cross-Z adjacency from stair pairs.
	crossZ := make(map[components.Position][]components.Position)
	addLink := func(a, b components.Position) {
		crossZ[a] = append(crossZ[a], b)
		crossZ[b] = append(crossZ[b], a)
	}
	for _, d := range gms.Downs {
		addLink(d, components.NewPosition(d.X, d.Y, d.Z-1))
	}
	for _, u := range gms.Ups {
		addLink(u, components.NewPosition(u.X, u.Y, u.Z+1))
	}

	walkable := func(p components.Position) bool {
		g, ok := gms.Grids[p.Z]
		if !ok {
			return false
		}
		if p.X < 0 || p.Y < 0 || p.X >= assets.Settings.MapWidth() || p.Y >= assets.Settings.MapHeight() {
			return false
		}
		return g.Get(p.X, p.Y).Walkable
	}

	nextID := 1
	stack := make([]components.Position, 0, 64)

	for p := range gms.Tiles {
		if _, done := regions[p]; done {
			continue
		}
		if !walkable(p) {
			continue
		}

		id := nextID
		nextID++
		regions[p] = id
		stack = append(stack[:0], p)

		for len(stack) > 0 {
			cur := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			// 8-neighbour flood on the same Z.
			for dx := -1; dx <= 1; dx++ {
				for dy := -1; dy <= 1; dy++ {
					if dx == 0 && dy == 0 {
						continue
					}
					n := components.NewPosition(cur.X+dx, cur.Y+dy, cur.Z)
					if !walkable(n) {
						continue
					}
					if _, done := regions[n]; done {
						continue
					}
					regions[n] = id
					stack = append(stack, n)
				}
			}

			// Stair links bridge to the floor above/below.
			for _, n := range crossZ[cur] {
				if !walkable(n) {
					continue
				}
				if _, done := regions[n]; done {
					continue
				}
				regions[n] = id
				stack = append(stack, n)
			}
		}
	}

	gms.RegionIDs = regions
}
