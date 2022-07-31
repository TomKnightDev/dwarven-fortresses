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
			if (x == 0 && y == 0) || pos.X+x < 0 || pos.Y+y < 0 || pos.X+x > assets.WorldWidth-1 || pos.Y+y > assets.WorldHeight-1 {
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

func StockpileLocations(w engine.World, itemType enums.ItemTypeEnum, assignItemType bool) []components.Position {
	gm := GetGameMapSingleton(w)

	var empty []components.Position
	for pos, it := range gm.Stockpiles {
		if it == itemType {
			return []components.Position{
				pos,
			}
		} else if it == enums.ItemTypeNone {
			empty = append(empty, pos)
		}
	}

	return empty
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

			break
		}
	}
}
