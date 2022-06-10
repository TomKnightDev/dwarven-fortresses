package systems

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/sedyh/mizu/pkg/engine"
	"github.com/tomknightdev/dwarven-fortresses/components"
	"github.com/tomknightdev/dwarven-fortresses/enums"
)

type Debug struct {
}

func NewDebug() *Debug {
	return &Debug{}
}

func (d *Debug) Update(w engine.World) {

}

func (d *Debug) Draw(w engine.World, screen *ebiten.Image) {
	// Debug information
	msg := fmt.Sprintf("TPS: %0.2f FPS: %0.2f\n",
		ebiten.CurrentTPS(), ebiten.CurrentFPS())

	msg += fmt.Sprintf("JOB COUNT: %d\n", len(w.View(components.Job{}).Filter()))

	// msg += showStockpileInventory(w)

	ebitenutil.DebugPrint(screen, msg)

}

func showStockpileInventory(w engine.World) string {
	ents := w.View(components.Designation{}, components.Position{}, components.Inventory{}).Filter()
	var d *components.Designation
	var i *components.Inventory

	items := make(map[enums.ItemTypeEnum]int)
	for _, e := range ents {
		e.Get(&d, &i)

		items[d.ItemType] += len(i.Items)
	}

	var msg string

	for t, c := range items {
		msg += fmt.Sprintf("%d: %d\n", t, c)
	}

	return msg
}
