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
	ents := w.View(components.Designation{}, components.Position{}, components.Inventory{})

	var d *components.Designation
	var p *components.Position
	var i *components.Inventory

	var spPositions []components.Position
	var firstFree int
	var itemTypePositions []components.Position
	ents.Each(func(e engine.Entity) {
		e.Get(&d, &p, &i)

		if len(i.Items) >= d.MaxItems {
			return
		}

		if d.DesignationType == enums.DesignationTypeStockpile {
			spPositions = append(spPositions, *p)

			if d.ItemType == itemType {
				itemTypePositions = append(itemTypePositions, *p)
			}

			if firstFree == 0 && d.ItemType == enums.ItemTypeNone {
				firstFree = e.ID()
			}
		}
	})

	if len(spPositions) == 0 {
		return spPositions
	}

	if len(itemTypePositions) > 0 {
		return itemTypePositions
	}

	if firstFree > 0 && assignItemType && itemType != enums.ItemTypeNone {
		e, found := w.GetEntity(firstFree)
		if !found {
			log.Println("somehow failed to find entity")
			return spPositions
		}

		e.Get(&d, &p)
		if d == nil {
			log.Println("designation component not found")
			return nil
		}

		d.ItemType = itemType
		return []components.Position{*p}
	}

	return spPositions
}

func AddItemToStockpile(w engine.World, pos components.Position, itemID, quatity int) {
	ents := w.View(components.Designation{}, components.Position{}, components.Inventory{}).Filter()
	var p *components.Position
	var d *components.Designation
	var i *components.Inventory

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
