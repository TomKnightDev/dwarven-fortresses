package scenes

import (
	"github.com/sedyh/mizu/pkg/engine"
	"github.com/tomknightdev/dwarven-fortresses/assets"
	"github.com/tomknightdev/dwarven-fortresses/components"
	"github.com/tomknightdev/dwarven-fortresses/entities"
	"github.com/tomknightdev/dwarven-fortresses/enums"
	"github.com/tomknightdev/dwarven-fortresses/systems"
)

type MainMenu struct {
}

func NewMainMenu() *MainMenu {
	return &MainMenu{}
}

func (mm *MainMenu) Setup(w engine.World) {
	w.AddComponents(
		components.InputSingleton{},
		components.Mouse{},
		components.Gui{},
		components.Text{},
		components.Sprite{},
		components.Position{},
		components.RenderSingleton{},
		components.GameMapSingleton{},
		components.NatureSingleton{},
	)

	w.AddSystems(
		systems.NewInput(),
		systems.NewGui(NewGameScene),
	)

	// Admin entity
	w.AddEntities(&entities.Admin{
		InputSingleton:   components.NewInputSingleton(),
		RenderSingleton:  components.NewRenderSingleton(),
		GameMapSingleton: components.NewGameMapSingleton(),
		NatureSingleton:  components.NewNatureSingleton(),
	})

	w.AddEntities(
		&entities.Gui{
			Gui:    components.NewGui(0, 0, enums.GuiPositionCenter, 10.0, enums.GuiActionNewGame),
			Sprite: components.NewSprite(assets.OpaqueImages[enums.TileTypeNewGame]),
		},
		&entities.GuiText{
			Gui:  components.NewGui(0, 100, enums.GuiPositionTop, 10.0, enums.GuiActionNewGame),
			Text: components.NewText("Dwarven Fortress"),
		},
		&entities.GuiText{
			Gui:  components.NewGui(0, 120, enums.GuiPositionCenter, 10.0, enums.GuiActionNewGame),
			Text: components.NewText("Play"),
		},
		&entities.GuiText{
			Gui:  components.NewGui(0, 20, enums.GuiPositionBottom, 10.0, enums.GuiActionNewGame),
			Text: components.NewText("By Tom Knight & Leigh Lawley"),
		},
		&entities.GuiText{
			Gui:  components.NewGui(0, 0, enums.GuiPositionRight, 10.0, enums.GuiActionNewGame),
			Text: components.NewText(">"),
		},
		&entities.GuiText{
			Gui:  components.NewGui(0, 0, enums.GuiPositionLeft, 10.0, enums.GuiActionNewGame),
			Text: components.NewText("<"),
		})
}

func NewGameScene(w engine.World) {
	w.ChangeScene(NewGame())
}
