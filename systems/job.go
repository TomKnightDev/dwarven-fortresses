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
	// Create a haul job for the first unclaimed, haulable item not yet in a stockpile.
	items := w.View(components.Item{}, components.Position{}).Filter()
	var i *components.Item
	var p *components.Position
	for _, e := range items {
		e.Get(&i, &p)
		if !i.Claimed && i.Haulable && !i.InStockpile {
			spPoses := helpers.StockpileLocations(w, i.ItemType, true)
			if len(spPoses) > 0 {
				hj := entities.Job{
					Job: components.NewJob(components.NewTask(*p, enums.TaskTypePickUp, 1), components.NewTask(spPoses[0], enums.TaskTypeAddToStockpile, 1)),
				}
				hj.Job.EntityId = e.ID()
				w.AddEntities(&hj)
				i.Claimed = true
				break
			}
		}
	}

	is := helpers.GetInputSingleton(w)

	if len(is.LeftClickedTiles) > 0 {
		switch is.InputMode {
		case enums.InputModeChop:
			helpers.CreateJobs(w, enums.TaskTypeChop, 5, is.LeftClickedTiles...)

		case enums.InputModeMine:
			helpers.CreateJobs(w, enums.TaskTypeMine, 10, is.LeftClickedTiles...)

		case enums.InputModeBuild:
			helpers.CreateJobs(w, enums.TaskTypeBuild, 15, components.NewPosition(is.LeftClickedTiles[0].X, is.LeftClickedTiles[0].Y, is.LeftClickedTiles[0].Z))
		}
	} else if len(is.RightClickedTiles) > 0 {
		helpers.CancelJobs(w, is.RightClickedTiles...)
	}

}

func (j *Job) Draw(w engine.World, screen *ebiten.Image) {
	camera, found := w.View(components.Zoom{}, components.Position{}).Get()
	if !found {
		return
	}
	var camPos *components.Position
	camera.Get(&camPos)

	rs := helpers.GetRenderSingleton(w)
	cursorImg := assets.TransImages[enums.TileTypeCursor]
	tileSize := float64(assets.TileSize)

	var op ebiten.DrawImageOptions
	var t *components.Job
	w.View(components.Job{}).Each(func(e engine.Entity) {
		e.Get(&t)
		for _, task := range t.Tasks {
			if task.Position.Z != camPos.Z {
				continue
			}
			op.GeoM.Reset()
			op.GeoM.Translate(float64(task.Position.X)*tileSize, float64(task.Position.Y)*tileSize)
			rs.OffScreen.DrawImage(cursorImg, &op)
		}
	})
}
