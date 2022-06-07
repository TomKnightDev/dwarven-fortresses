package systems

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sedyh/mizu/pkg/engine"
	"github.com/tomknightdev/dwarven-fortresses/assets"
	"github.com/tomknightdev/dwarven-fortresses/components"
	"github.com/tomknightdev/dwarven-fortresses/enums"
)

type Gui struct {
}

func NewGui() *Gui {
	return &Gui{}
}

func (g *Gui) Update(w engine.World) {
	var inputSingleton *components.InputSingleton
	is, found := w.View(components.InputSingleton{}).Get()
	if !found {
		panic("input singleton was not found")
	}
	is.Get(&inputSingleton)

	inputSingleton.InGui = false

	ents := w.View(components.Gui{}, components.Sprite{}).Filter()
	for _, e := range ents {
		var gsp *components.Sprite
		var gui *components.Gui
		e.Get(&gsp, &gui)

		if g.Within(*gui, inputSingleton.MousePosX, inputSingleton.MousePosY) {
			inputSingleton.InGui = true
			if inputSingleton.IsMouseLeftPressed {
				switch gui.Action {
				case enums.GuiActionStair:
					inputSingleton.InputMode = enums.InputModeBuild
				case enums.GuiActionChop:
					inputSingleton.InputMode = enums.InputModeChop
				case enums.GuiActionMine:
					inputSingleton.InputMode = enums.InputModeMine
				case enums.GuiActionStockpile:
					inputSingleton.InputMode = enums.InputModeStockpile
				}
			}
		}
	}
}

func (g *Gui) Draw(w engine.World, screen *ebiten.Image) {
	view := w.View(components.Gui{}, components.Sprite{})
	view.Each(func(e engine.Entity) {
		var gui *components.Gui
		var sprite *components.Sprite
		e.Get(&gui, &sprite)
		op := &ebiten.DrawImageOptions{}

		op.GeoM.Translate(float64(gui.X)/gui.Scale, float64(gui.Y)/gui.Scale)
		op.GeoM.Scale(gui.Scale, gui.Scale)
		screen.DrawImage(sprite.Image, op)
	})

}

func (g Gui) Within(gui components.Gui, x, y int) bool {
	sx, sy := g.scalePos(gui)
	if x > gui.X && x < sx && y > gui.Y && y < sy {
		return true
	}

	return false
}

func (g Gui) scalePos(gui components.Gui) (int, int) {
	x := gui.X + int(gui.Scale)*assets.TileSize
	y := gui.Y + int(gui.Scale)*assets.TileSize
	return x, y
}
