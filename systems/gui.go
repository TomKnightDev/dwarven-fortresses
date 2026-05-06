package systems

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/sedyh/mizu/pkg/engine"
	"github.com/tomknightdev/dwarven-fortresses/assets"
	"github.com/tomknightdev/dwarven-fortresses/components"
	"github.com/tomknightdev/dwarven-fortresses/enums"
	"github.com/tomknightdev/dwarven-fortresses/helpers"
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

	view := w.View(components.Gui{}, components.Flex{})
	view.Each(func(e engine.Entity) {
		var gui *components.Gui
		var flex *components.Flex
		e.Get(&gui, &flex)

		if flex.View != nil {
			flex.View.UpdateWithSize(ebiten.WindowSize())
			return
		}
	})
}

func (g *Gui) Draw(w engine.World, screen *ebiten.Image) {
	RenderGame(w, screen)
	RenderUI(w, screen)
	RenderResourceCounts(w, screen)
	RenderStockpileTooltip(w, screen)
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
	view := w.View(components.Gui{}, components.Flex{})
	view.Each(func(e engine.Entity) {
		var gui *components.Gui
		var flex *components.Flex
		e.Get(&gui, &flex)

		if flex.View != nil {
			flex.View.Draw(screen)
			return
		}
	})
}

func RenderResourceCounts(w engine.World, screen *ebiten.Image) {
	var logs, stones int
	view := w.View(components.Item{}, components.Position{})
	view.Each(func(e engine.Entity) {
		var item *components.Item
		e.Get(&item)
		if !item.InStockpile {
			return
		}
		switch item.ItemType {
		case enums.ItemTypeLog:
			logs++
		case enums.ItemTypeStone:
			stones++
		}
	})

	width, _ := ebiten.WindowSize()
	woodLabel := fmt.Sprintf("Wood: %d", logs)
	stoneLabel := fmt.Sprintf("Rock: %d", stones)

	woodBounds := text.BoundString(assets.MainFont12, woodLabel)
	stoneBounds := text.BoundString(assets.MainFont12, stoneLabel)

	padding := 10
	x := width - stoneBounds.Dx() - padding
	if woodBounds.Dx() > stoneBounds.Dx() {
		x = width - woodBounds.Dx() - padding
	}

	text.Draw(screen, woodLabel, assets.MainFont12, x, 20, color.White)
	text.Draw(screen, stoneLabel, assets.MainFont12, x, 40, color.White)
}

func RenderStockpileTooltip(w engine.World, screen *ebiten.Image) {
	is := helpers.GetInputSingleton(w)
	if is.InGui || is.MouseWorldPosX < 0 || is.MouseWorldPosY < 0 {
		return
	}

	mouseEnt, found := w.View(components.Mouse{}, components.Position{}).Get()
	if !found {
		return
	}
	var mComp *components.Mouse
	var mousePos *components.Position
	mouseEnt.Get(&mComp, &mousePos)

	gm := helpers.GetGameMapSingleton(w)
	if _, ok := gm.Stockpiles[*mousePos]; !ok {
		return
	}

	// Count items by type in this stockpile tile
	counts := make(map[enums.ItemTypeEnum]int)
	ents := w.View(components.Designation{}, components.Position{}, components.Inventory{}).Filter()
	var p *components.Position
	var d *components.Designation
	var inv *components.Inventory
	for _, e := range ents {
		e.Get(&d, &p, &inv)
		if !helpers.Matches(*p, *mousePos) {
			continue
		}
		for _, itemID := range inv.Items {
			itemEnt, found := w.GetEntity(itemID)
			if !found {
				continue
			}
			var item *components.Item
			itemEnt.Get(&item)
			counts[item.ItemType]++
		}
		break
	}

	var lines []string
	if len(counts) == 0 {
		lines = []string{"Empty stockpile"}
	} else {
		for itemType, count := range counts {
			var name string
			switch itemType {
			case enums.ItemTypeLog:
				name = "Wood"
			case enums.ItemTypeStone:
				name = "Rock"
			default:
				name = "Unknown"
			}
			lines = append(lines, fmt.Sprintf("%s: %d/%d", name, count, helpers.StockpileMaxItems))
		}
	}

	const lineH = 16
	const padding = 4
	maxW := 0
	for _, l := range lines {
		b := text.BoundString(assets.MainFont12, l)
		if b.Dx() > maxW {
			maxW = b.Dx()
		}
	}

	bgW := maxW + padding*2
	bgH := len(lines)*lineH + padding*2
	bg := ebiten.NewImage(bgW, bgH)
	bg.Fill(color.RGBA{0, 0, 0, 200})

	tx := is.MousePosX + 16
	ty := is.MousePosY + 8
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(tx-padding), float64(ty-lineH+padding))
	screen.DrawImage(bg, op)

	for i, line := range lines {
		text.Draw(screen, line, assets.MainFont12, tx, ty+i*lineH, color.White)
	}
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

// func calculatePosition(x, y, xOffset, yOffset int, gui *components.Gui) (int, int) {
// 	width, height := ebiten.WindowSize()
// 	xCentre := width / 2
// 	yCentre := height / 2
// 	switch gui.Position {
// 	case enums.GuiPositionCenter:
// 		return x + xCentre - xOffset/2, y + yCentre - yOffset/2
// 	case enums.GuiPositionTop:
// 		return x + xCentre - xOffset/2, y + yOffset
// 	case enums.GuiPositionBottom:
// 		return x + xCentre - xOffset/2, height - y
// 	case enums.GuiPositionLeft:
// 		return x + 0, y + yCentre - yOffset/2
// 	case enums.GuiPositionRight:
// 		return width - x - xOffset, y + yCentre - yOffset
// 	default:
// 		return x, y
// 	}
// }
