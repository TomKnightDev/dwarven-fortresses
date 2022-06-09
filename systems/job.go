package systems

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sedyh/mizu/pkg/engine"
	"github.com/tomknightdev/dwarven-fortresses/assets"
	"github.com/tomknightdev/dwarven-fortresses/components"
	"github.com/tomknightdev/dwarven-fortresses/entities"
	"github.com/tomknightdev/dwarven-fortresses/enums"
	"github.com/tomknightdev/dwarven-fortresses/helpers"
)

type Job struct {
}

func NewJob() *Job {
	return &Job{}
}

func (j *Job) Update(w engine.World) {
	var inputSingleton *components.InputSingleton
	is, found := w.View(components.InputSingleton{}).Get()
	if !found {
		panic("input singleton was not found")
	}
	is.Get(&inputSingleton)

	if len(inputSingleton.LeftClickedTiles) > 0 {
		switch inputSingleton.InputMode {
		case enums.InputModeChop:
			resources := w.View(components.Choppable{}, components.Position{}).Filter()
			var rPos *components.Position
			for _, r := range resources {
				r.Get(&rPos)
				for _, st := range inputSingleton.LeftClickedTiles {
					if rPos.X == st.X && rPos.Y == st.Y && rPos.Z == st.Z {
						w.AddEntities(&entities.Job{
							Job: components.NewJob(r.ID(), components.NewTask(components.NewPosition(st.X, st.Y, st.Z), enums.TaskTypeChop, 10)),
						})
					}
				}

			}
		case enums.InputModeMine:
			gms, found := w.View(components.GameMapSingleton{}).Get()
			if !found {
				panic("game map singleton not found")
			}
			var gmComp *components.GameMapSingleton
			gms.Get(&gmComp)

			for _, st := range inputSingleton.LeftClickedTiles {
				t := gmComp.Tiles[components.NewPosition(st.X, st.Y, st.Z)]
				if t.TileTypeEnum != enums.TileTypeRock {
					continue
				}

				w.AddEntities(&entities.Job{
					Job: components.NewJob(0, components.NewTask(components.NewPosition(st.X, st.Y, st.Z), enums.TaskTypeMine, 10)),
				})
			}

		case enums.InputModeBuild:
			w.AddEntities(&entities.Job{
				Job: components.NewJob(0, components.NewTask(components.NewPosition(inputSingleton.LeftClickedTiles[0].X, inputSingleton.LeftClickedTiles[0].Y, inputSingleton.LeftClickedTiles[0].Z), enums.TaskTypeBuild, 10)),
			})
		}
	} else if len(inputSingleton.RightClickedTiles) > 0 {
		// Cancel ents
		ents := w.View(components.Job{}).Filter()

		if len(ents) == 0 {
			return
		}

		var entitiesToRemove []engine.Entity

		for _, e := range ents {
			var job *components.Job
			e.Get(&job)
			for _, t := range job.Tasks {
				for _, st := range inputSingleton.RightClickedTiles {
					if helpers.Matches(st, t.Position) {
						entitiesToRemove = append(entitiesToRemove, e)
					}
				}
			}
		}

		for _, job := range entitiesToRemove {
			w.RemoveEntity(job)
		}
	}

	// Create jobs for haulable items not in a stockpile
	items := w.View(components.Item{}, components.Position{})
	var i *components.Item
	var p *components.Position

	items.Each(func(e engine.Entity) {
		e.Get(&i, &p)
		if !i.Claimed && i.Haulable && !i.InStockpile {
			spPoses := helpers.StockpileLocations(w, i.ItemType, true)
			if len(spPoses) > 0 {
				w.AddEntities(&entities.Job{
					Job: components.NewJob(e.ID(), components.NewTask(*p, enums.TaskTypePickUp, 1), components.NewTask(spPoses[0], enums.TaskTypeAddToStockpile, 1)),
				})

				i.Claimed = true
			}
		}
	})
}

func (j *Job) Draw(w engine.World, screen *ebiten.Image) {
	ents := w.View(components.Job{})
	var t *components.Job

	ents.Each(func(e engine.Entity) {
		e.Get(&t)

		for _, task := range t.Tasks {
			helpers.DrawImage(w, screen, task.Position, assets.TransImages[enums.TileSpriteCursor])
		}
	})

}
