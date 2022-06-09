package systems

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sedyh/mizu/pkg/engine"
	"github.com/tomknightdev/dwarven-fortresses/components"
	"github.com/tomknightdev/dwarven-fortresses/helpers"
)

type Nature struct {
}

func NewNature() *Nature {
	return &Nature{}
}

func (n *Nature) Update(w engine.World) {
	// ne, found := w.View(components.NatureSingleton{}).Get()
	// if !found {
	// 	panic("unable to find entity with nature component")
	// }
	// var nc *components.NatureSingleton
	// ne.Get(&nc)

	// if nc.CurrentGrowTimer < nc.GrowTimer {
	// 	nc.CurrentGrowTimer++
	// 	return
	// }
	// nc.CurrentGrowTimer = 0

	// gms, found := w.View(components.GameMapSingleton{}).Get()
	// if !found {
	// 	panic("game map singleton not found")
	// }

	// var gmComp *components.GameMapSingleton
	// gms.Get(&gmComp)

	// // Pick a random tile, if dirt, make grass
	// tiles := gmComp.TilesByType[enums.TileSpriteGrass0]
	// rand.Seed(time.Now().UnixNano())
	// r := rand.Intn(len(tiles))

	// helpers.UpdateTile(w, enums.TileSpriteGrass0, enums.TileSpriteGrass0, r, gmComp)
}

func (n *Nature) Draw(w engine.World, screen *ebiten.Image) {
	ne, found := w.View(components.NatureSingleton{}).Get()
	if !found {
		panic("unable to find entity with nature component")
	}
	var nc *components.NatureSingleton
	ne.Get(&nc)

	ents := w.View(components.Nature{}, components.Sprite{}, components.Position{}).Filter()
	helpers.DrawImages(w, screen, nc.OffScreen, ents)
}
