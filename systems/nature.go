package systems

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sedyh/mizu/pkg/engine"
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
	// tiles := gmComp.TilesByType[enums.TileTypeGrass0]
	// rand.Seed(time.Now().UnixNano())
	// r := rand.Intn(len(tiles))

	// helpers.UpdateTile(w, enums.TileTypeGrass0, enums.TileTypeGrass0, r, gmComp)
}

func (n *Nature) Draw(w engine.World, screen *ebiten.Image) {
	// ents := w.View(components.Nature{}, components.Sprite{}, components.Position{})
	// helpers.DrawImages(w, ents)
}
