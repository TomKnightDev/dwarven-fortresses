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

	// gms := helpers.GetGameMapSingleton(w)

	var entitiesToRemove []engine.Entity
	var job *components.Job
	for _, e := range ents {
		e.Get(&job)

		currentTask := job.Tasks[0]

		if currentTask.Completed {
			pos := currentTask.Position
			switch currentTask.TaskTypeEnum {
			case enums.TaskTypeChop:

				helpers.RemoveResourceFromTile(w, pos, enums.ResourceTypeTree, true)

				w.AddEntities(&entities.Item{
					Position: pos,
					Sprite:   components.NewSprite(assets.TransImages[enums.TileTypeLog0]),
					Item:     components.NewItem(true, 25, enums.ItemTypeLog),
				})

			case enums.TaskTypeBuild:
				helpers.AddBuildingToTile(w, pos, enums.TileTypeStairDown, true)
				helpers.AddBuildingToTile(w, components.NewPosition(pos.X, pos.Y, pos.Z-1), enums.TileTypeStairUp, true)
				helpers.MineTile(w, components.NewPosition(pos.X, pos.Y, pos.Z-1))

			case enums.TaskTypeMine:
				helpers.MineTile(w, pos)

				for i := 0; i < 1; i++ {
					w.AddEntities(&entities.Item{
						Position: pos,
						Sprite:   components.NewSprite(assets.TransImages[enums.TileTypeRocks]),
						Item:     components.NewItem(true, 25, enums.ItemTypeStone),
					})
				}

			case enums.TaskTypePickUp:
				item, found := w.GetEntity(job.EntityId)
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
				inv.Items = append(inv.Items, job.EntityId)
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

				item, found := w.GetEntity(job.EntityId)
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
					if it == job.EntityId {
						inv.Items = append(inv.Items[:i], inv.Items[i+1:]...)
					}
				}
				inv.Weight -= ite.Weight

				helpers.AddItemToStockpile(w, *itemPos, job.EntityId, 1)
				// }
			}

			// Update surronding tasks that are blocked
			updateSurroundingTasks(job.Tasks[0].Position, ents)

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

func updateSurroundingTasks(pos components.Position, ents []engine.Entity) {
	var job *components.Job
	for _, e := range ents {
		e.Get(&job)

		if len(job.Tasks) <= 0 || !job.Tasks[0].Blocked {
			continue
		}

		if helpers.IsAdjacent(components.NewMove(pos.X, pos.Y, pos.Z), job.Tasks[0].Position) {
			job.Tasks[0].Blocked = false
		}
	}
}
