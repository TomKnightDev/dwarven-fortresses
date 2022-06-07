package systems

import (
	"log"

	"github.com/sedyh/mizu/pkg/engine"
	"github.com/tomknightdev/dwarven-fortresses/assets"
	"github.com/tomknightdev/dwarven-fortresses/components"
	"github.com/tomknightdev/dwarven-fortresses/entities"
	"github.com/tomknightdev/dwarven-fortresses/enums"
	"github.com/tomknightdev/dwarven-fortresses/helpers"
)

type Task struct {
}

func NewTask() *Task {
	return &Task{}
}

func (t *Task) Update(w engine.World) {
	ents := w.View(components.Job{}).Filter()

	if len(ents) == 0 {
		return
	}

	gms, found := w.View(components.GameMapSingleton{}).Get()
	if !found {
		panic("game map singleton not found")
	}

	var gmComp *components.GameMapSingleton
	gms.Get(&gmComp)

	var entitiesToRemove []engine.Entity
	var job *components.Job
	for _, e := range ents {
		e.Get(&job)

		currentTask := job.Tasks[0]

		if currentTask.Completed {
			pos := currentTask.Position
			switch currentTask.TaskTypeEnum {
			case enums.TaskTypeChop:
				ent, ok := w.GetEntity(job.EntityID)
				if !ok {
					log.Println("entity not found ", job.EntityID)
				}

				cell := gmComp.Grids[pos.Z].Get(pos.X, pos.Y)
				cell.Walkable = true

				var drop *components.Drops
				ent.Get(&drop)

				for i := 0; i < drop.DropCount; i++ {
					w.AddEntities(&entities.Item{
						Position: pos,
						Sprite:   components.NewSprite(assets.TransImages[enums.TileTypeLog0]),
						Item:     components.NewItem(true, 25, enums.ItemTypeLog),
					})
				}

				entitiesToRemove = append(entitiesToRemove, ent)
			case enums.TaskTypeBuild:
				w.AddEntities(&entities.Building{
					Position: pos,
					Sprite:   components.NewSprite(assets.OpaqueImages[enums.TileTypeStairDown]),
					TileType: components.NewTileType(enums.TileTypeStairDown),
					Building: components.NewBuilding(),
				})
				gmComp.TilesByType[enums.TileTypeStairDown] = append(gmComp.TilesByType[enums.TileTypeStairDown], pos)

				index, err := helpers.GetTileByTypeIndexFromPos(components.NewPosition(pos.X, pos.Y, pos.Z-1), gmComp.TilesByType[enums.TileTypeRockWallHz])
				if err != nil {
					log.Println(err)
					entitiesToRemove = append(entitiesToRemove, e)
					continue
				}

				helpers.UpdateTile(w, enums.TileTypeRockWallHz, enums.TileTypeRockFloor, index, gmComp)

				w.AddEntities(&entities.Building{
					Position: components.NewPosition(pos.X, pos.Y, pos.Z-1),
					Sprite:   components.NewSprite(assets.OpaqueImages[enums.TileTypeStairUp]),
					TileType: components.NewTileType(enums.TileTypeStairUp),
					Building: components.NewBuilding(),
				})
				gmComp.TilesByType[enums.TileTypeStairUp] = append(gmComp.TilesByType[enums.TileTypeStairUp], components.NewPosition(pos.X, pos.Y, pos.Z-1))

			case enums.TaskTypeMine:
				tileType := enums.TileTypeRockWallHz
				index, err := helpers.GetTileByTypeIndexFromPos(components.NewPosition(pos.X, pos.Y, pos.Z), gmComp.TilesByType[enums.TileTypeRockWallHz])
				if err != nil {
					tileType = enums.TileTypeRockWallVt
					index, err = helpers.GetTileByTypeIndexFromPos(components.NewPosition(pos.X, pos.Y, pos.Z), gmComp.TilesByType[enums.TileTypeRockWallVt])

					if err != nil {
						log.Println(err)
						entitiesToRemove = append(entitiesToRemove, e)
						continue
					}
				}

				helpers.UpdateTile(w, tileType, enums.TileTypeRockFloor, index, gmComp)

				for i := 0; i < 1; i++ {
					w.AddEntities(&entities.Item{
						Position: pos,
						Sprite:   components.NewSprite(assets.TransImages[enums.TileTypeRocks]),
						Item:     components.NewItem(true, 25, enums.ItemTypeStone),
					})
				}

			case enums.TaskTypePickUp:
				item, found := w.GetEntity(job.EntityID)
				if !found {
					log.Println("item entity not found")
				}
				var ite *components.Item
				item.Get(&ite)

				actor, found := w.GetEntity(job.ClaimedByID)
				if !found {
					log.Println("actor entity not found")
				}

				var inv *components.Inventory
				actor.Get(&inv)
				inv.Items = append(inv.Items, job.EntityID)
				inv.Weight += ite.Weight

				var itemSprite *components.Sprite
				item.Get(&itemSprite)
				itemSprite.Drawn = false

			case enums.TaskTypeAddToStockpile:
				actor, found := w.GetEntity(job.ClaimedByID)
				if !found {
					log.Println("actor entity not found")
				}

				var inv *components.Inventory
				var move *components.Move
				actor.Get(&inv, &move)

				item, found := w.GetEntity(job.EntityID)
				if !found {
					log.Println("item entity not found")
				}
				var ite *components.Item
				var itemPos *components.Position
				var itemSprite *components.Sprite
				item.Get(&ite, &itemPos, &itemSprite)
				itemPos.X = currentTask.Position.X
				itemPos.Y = currentTask.Position.Y
				itemPos.Z = currentTask.Position.Z
				itemSprite.Drawn = true

				for i, it := range inv.Items {
					if it == job.EntityID {
						inv.Items = append(inv.Items[:i], inv.Items[i+1:]...)
					}
				}
				inv.Weight -= ite.Weight

				helpers.AddItemToStockpile(w, *itemPos, job.EntityID, 1)
				// }
			}

			job.Tasks = job.Tasks[1:]
			if len(job.Tasks) == 0 {
				entitiesToRemove = append(entitiesToRemove, e)
			}
		}
	}

	for _, job := range entitiesToRemove {
		w.RemoveEntity(job)
	}
}
