package helpers

import (
	"log"

	"github.com/OpenSauce/paths"
	"github.com/sedyh/mizu/pkg/engine"
	"github.com/tomknightdev/dwarven-fortresses/assets"
	"github.com/tomknightdev/dwarven-fortresses/components"
	"github.com/tomknightdev/dwarven-fortresses/enums"
)

func Matches(a components.Position, b components.Position) bool {
	if a.X == b.X && a.Y == b.Y && a.Z == b.Z {
		return true
	}

	return false
}

func GetAdjacents(grid *paths.Grid, pos components.Position, walkable bool) []components.Position {
	var adjacents []components.Position

	for x := -1; x < 2; x++ {
		for y := -1; y < 2; y++ {
			if (x == 0 && y == 0) || pos.X+x < 0 || pos.Y+y < 0 || pos.X+x > assets.Settings.MapWidth()-1 || pos.Y+y > assets.Settings.MapHeight()-1 {
				continue
			}

			if grid.Get(pos.X+x, pos.Y+y).Walkable == walkable {
				adjacents = append(adjacents, components.NewPosition(pos.X+x, pos.Y+y, pos.Z))
			}
		}
	}

	return adjacents
}

func IsAdjacent(dest components.Move, current components.Position) bool {
	if current.X >= dest.X-1 && current.X <= dest.X+1 && current.Y >= dest.Y-1 && current.Y <= dest.Y+1 && current.Z == dest.Z {
		return true
	}

	return false
}

const StockpileMaxItems = 10

func stockpileEffectiveCounts(w engine.World) map[components.Position]int {
	// Start with items already deposited
	counts := make(map[components.Position]int)
	ents := w.View(components.Designation{}, components.Position{}, components.Inventory{}).Filter()
	var p *components.Position
	var inv *components.Inventory
	for _, e := range ents {
		e.Get(&p, &inv)
		counts[*p] = len(inv.Items)
	}

	// Add in-flight haul jobs so we don't over-assign a tile before dwarves arrive
	jobs := w.View(components.Job{}).Filter()
	var job *components.Job
	for _, e := range jobs {
		e.Get(&job)
		for _, task := range job.Tasks {
			if task.TaskTypeEnum == enums.TaskTypeAddToStockpile && !task.Completed {
				counts[task.Position]++
				break
			}
		}
	}

	return counts
}

func StockpileLocations(w engine.World, itemType enums.ItemTypeEnum, assignItemType bool) []components.Position {
	gm := GetGameMapSingleton(w)
	counts := stockpileEffectiveCounts(w)

	// Prefer tiles already locked to this item type, then fall back to empty tiles.
	// In both cases, pick the fullest available tile so we fill one before starting another.
	bestTyped, bestTypedCount := components.Position{}, -1
	bestEmpty, bestEmptyCount := components.Position{}, -1
	foundTyped, foundEmpty := false, false

	for pos, it := range gm.Stockpiles {
		if counts[pos] >= StockpileMaxItems {
			continue
		}
		switch it {
		case itemType:
			if !foundTyped || counts[pos] > bestTypedCount {
				bestTyped, bestTypedCount = pos, counts[pos]
				foundTyped = true
			}
		case enums.ItemTypeNone:
			if !foundEmpty || counts[pos] > bestEmptyCount {
				bestEmpty, bestEmptyCount = pos, counts[pos]
				foundEmpty = true
			}
		}
	}

	if foundTyped {
		return []components.Position{bestTyped}
	}
	if foundEmpty {
		return []components.Position{bestEmpty}
	}
	return nil
}

func AddItemToStockpile(w engine.World, pos components.Position, itemID, quatity int) {
	ents := w.View(components.Designation{}, components.Position{}, components.Inventory{}).Filter()
	var p *components.Position
	var d *components.Designation
	var i *components.Inventory

	gm := GetGameMapSingleton(w)

	for _, e := range ents {
		e.Get(&p, &d, &i)

		if Matches(*p, pos) {
			i.Items = append(i.Items, itemID)

			item, found := w.GetEntity(itemID)
			if !found {
				log.Println("item not found")
			}

			var it *components.Item
			item.Get(&it)
			it.InStockpile = true
			it.Claimed = false

			gm.Stockpiles[pos] = it.ItemType

			break
		}
	}
}

func RemoveItemFromStockpile(w engine.World, pos components.Position, itemID, quantity int) {
	ents := w.View(components.Designation{}, components.Position{}, components.Inventory{}).Filter()
	var p *components.Position
	var d *components.Designation
	var i *components.Inventory

	for _, e := range ents {
		e.Get(&p, &d, &i)

		if Matches(*p, pos) {
			for index, t := range i.Items {
				if t == itemID {
					i.Items = append(i.Items[:index], i.Items[index+1:]...)
				}
			}

			item, found := w.GetEntity(itemID)
			if !found {
				log.Println("item not found")
			}

			var it *components.Item
			item.Get(&it)
			it.InStockpile = false

			if len(i.Items) == 0 {
				gm := GetGameMapSingleton(w)
				gm.Stockpiles[pos] = enums.ItemTypeNone
			}

			break
		}
	}
}
