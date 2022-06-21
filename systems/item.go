package systems

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sedyh/mizu/pkg/engine"
	"github.com/tomknightdev/dwarven-fortresses/components"
	"github.com/tomknightdev/dwarven-fortresses/helpers"
)

type Item struct{}

func NewItem() *Item {
	return &Item{}
}

func (i *Item) Update(w engine.World) {

}

func (i *Item) Draw(w engine.World, screen *ebiten.Image) {
	ents := w.View(components.Item{}, components.Position{}, components.Sprite{})

	// var p *components.Position
	// var s *components.Sprite

	helpers.DrawImages(w, ents)

	// ents.Each(func(e engine.Entity) {
	// 	e.Get(&p, &s)

	// 	if !s.Drawn {
	// 		return
	// 	}

	// 	helpers.DrawImage(w, screen, *p, s.Image)
	// })
}
