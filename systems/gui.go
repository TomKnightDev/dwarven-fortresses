package systems

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/sedyh/mizu/pkg/engine"
	"github.com/tomknightdev/dwarven-fortresses/assets"
	"github.com/tomknightdev/dwarven-fortresses/components"
	"github.com/tomknightdev/dwarven-fortresses/enums"
)

type Gui struct {
	NewGame func(engine.World)
}

func NewGui(newGame func(engine.World)) *Gui {
	return &Gui{
		NewGame: newGame,
	}
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
				case enums.GuiActionNewGame:
					g.NewGame(w)
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
	RenderSprites(w, screen)
	RenderText(w, screen)
}

func RenderSprites(w engine.World, screen *ebiten.Image) {
	width, height := ebiten.WindowSize()
	xCentre := width / 2
	yCentre := height / 2

	view := w.View(components.Gui{}, components.Sprite{})
	view.Each(func(e engine.Entity) {
		var gui *components.Gui
		var sprite *components.Sprite
		e.Get(&gui, &sprite)
		op := &ebiten.DrawImageOptions{}

		if gui.Position == enums.GuiPositionRelative {
			gui.X = xCentre - assets.TileSize*int(gui.Scale)/2
			gui.Y = yCentre - assets.TileSize*int(gui.Scale)/2
		}

		op.GeoM.Translate(float64(gui.X)/gui.Scale, float64(gui.Y)/gui.Scale)

		op.GeoM.Scale(gui.Scale, gui.Scale)
		screen.DrawImage(sprite.Image, op)
	})
}

func RenderText(w engine.World, screen *ebiten.Image) {
	width, height := ebiten.WindowSize()
	xCentre := width / 2
	yCentre := height / 2

	view := w.View(components.Gui{}, components.Text{})
	view.Each(func(e engine.Entity) {
		var gui *components.Gui
		var ctext *components.Text
		e.Get(&gui, &ctext)
		op := &ebiten.DrawImageOptions{}

		if gui.Position == enums.GuiPositionRelative {
			gui.X = xCentre - assets.TileSize*int(gui.Scale)/2
			gui.Y = yCentre - assets.TileSize*int(gui.Scale)/2
		}

		op.GeoM.Translate(float64(gui.X)/gui.Scale, float64(gui.Y)/gui.Scale)

		op.GeoM.Scale(gui.Scale, gui.Scale)
		text.Draw(screen, ctext.Content, mplusBigFont, gui.X, gui.Y, color.White)
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
