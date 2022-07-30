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
	RenderGame(w, screen)
	RenderUI(w, screen)
}

func RenderGame(w engine.World, screen *ebiten.Image) {
	width, height := ebiten.WindowSize()
	xCentre := width / 2
	yCentre := height / 2

	view := w.View(components.Gui{}, components.Sprite{})
	view.Each(func(e engine.Entity) {
		var gui *components.Gui
		var sprite *components.Sprite
		e.Get(&gui, &sprite)
		op := &ebiten.DrawImageOptions{}

		if gui.Position == enums.GuiPositionCenter {
			gui.X = xCentre - assets.TileSize*int(gui.Scale)/2
			gui.Y = yCentre - assets.TileSize*int(gui.Scale)/2
		}

		op.GeoM.Translate(float64(gui.X)/gui.Scale, float64(gui.Y)/gui.Scale)

		op.GeoM.Scale(gui.Scale, gui.Scale)
		screen.DrawImage(sprite.Image, op)
	})
}

func RenderUI(w engine.World, screen *ebiten.Image) {
	view := w.View(components.Gui{}, components.Text{})
	view.Each(func(e engine.Entity) {
		var gui *components.Gui
		var ctext *components.Text
		e.Get(&gui, &ctext)

		if gui.UIUpdate != nil {
			gui.UIUpdate(screen)
			return
		}

		offsetRect := text.BoundString(assets.MainFont, ctext.Content)

		x, y := calculatePosition(gui.X, gui.Y, offsetRect.Max.X, offsetRect.Dy(), gui)

		text.Draw(screen, ctext.Content, assets.MainFont, x, y, color.White)
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

func calculatePosition(x, y, xOffset, yOffset int, gui *components.Gui) (int, int) {
	width, height := ebiten.WindowSize()
	xCentre := width / 2
	yCentre := height / 2
	switch gui.Position {
	case enums.GuiPositionCenter:
		return x + xCentre - xOffset/2, y + yCentre - yOffset/2
	case enums.GuiPositionTop:
		return x + xCentre - xOffset/2, y + yOffset
	case enums.GuiPositionBottom:
		return x + xCentre - xOffset/2, height - y
	case enums.GuiPositionLeft:
		return x + 0, y + yCentre - yOffset/2
	case enums.GuiPositionRight:
		return width - x - xOffset, y + yCentre - yOffset
	default:
		return x, y
	}
}
