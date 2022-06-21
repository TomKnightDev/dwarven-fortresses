package systems

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sedyh/mizu/pkg/engine"
	"github.com/tomknightdev/dwarven-fortresses/assets"
	"github.com/tomknightdev/dwarven-fortresses/components"
	"github.com/tomknightdev/dwarven-fortresses/enums"
	"github.com/tomknightdev/dwarven-fortresses/helpers"
)

type Job struct {
}

func NewJob() *Job {
	return &Job{}
}

func (j *Job) Update(w engine.World) {
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
	ents := w.View(components.Job{})
	var t *components.Job

	ents.Each(func(e engine.Entity) {
		e.Get(&t)

		for _, task := range t.Tasks {
			helpers.DrawImage(w, task.Position, assets.TransImages[enums.TileTypeCursor])
		}
	})

}
