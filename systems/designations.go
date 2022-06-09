package systems

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sedyh/mizu/pkg/engine"
	"github.com/tomknightdev/dwarven-fortresses/assets"
	"github.com/tomknightdev/dwarven-fortresses/components"
	"github.com/tomknightdev/dwarven-fortresses/entities"
	"github.com/tomknightdev/dwarven-fortresses/enums"
	"github.com/tomknightdev/dwarven-fortresses/helpers"
)

type Designations struct {
}

func NewDesignations() *Designations {
	return &Designations{}
}

func (d *Designations) Update(w engine.World) {
	var inputSingleton *components.InputSingleton
	is, found := w.View(components.InputSingleton{}).Get()
	if !found {
		panic("input singleton was not found")
	}
	is.Get(&inputSingleton)

	if len(inputSingleton.LeftClickedTiles) > 0 {
		switch inputSingleton.InputMode {
		case enums.InputModeStockpile:
			for _, st := range inputSingleton.LeftClickedTiles {
				w.AddEntities(&entities.Stockpile{
					Designation: components.NewDesignation(enums.DesignationTypeStockpile),
					Position:    st,
					Sprite:      components.NewSprite(assets.TransImages[enums.TileSpriteStockpile]),
					Inventory:   components.NewInventory(),
				})
			}
		}
	} else if len(inputSingleton.RightClickedTiles) > 0 {
		switch inputSingleton.InputMode {
		case enums.InputModeStockpile:
			var entsToDelete []engine.Entity
			ents := w.View(components.Designation{}, components.Position{}, components.Inventory{})
			var d *components.Designation
			var p *components.Position
			var i *components.Inventory
			ents.Each(func(e engine.Entity) {
				e.Get(&d, &p, &i)
				if d.DesignationType != enums.DesignationTypeStockpile {
					return
				}

				for _, st := range inputSingleton.RightClickedTiles {
					if helpers.Matches(st, *p) {
						for _, item := range i.Items {
							itemEnt, found := w.GetEntity(item)
							if !found {
								log.Println("item not found")
							}
							var itemComp *components.Item
							itemEnt.Get(&itemComp)
							itemComp.InStockpile = false
						}
						entsToDelete = append(entsToDelete, e)
					}
				}
			})

			for _, e := range entsToDelete {
				w.RemoveEntity(e)
			}
		}
	}
}

func (d *Designations) Draw(w engine.World, screen *ebiten.Image) {
	ents := w.View(components.Designation{}, components.Position{}, components.Sprite{})
	var p *components.Position
	var s *components.Sprite
	ents.Each(func(e engine.Entity) {
		e.Get(&p, &s)

		helpers.DrawImage(w, screen, *p, s.Image)
	})
}
