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

	w.AddEntities(&entities.Gui{
		Gui:    components.NewGui(0, 0, enums.GuiPositionRelative, 10.0, enums.GuiActionNewGame),
		Sprite: components.NewSprite(assets.OpaqueImages[enums.TileTypeNewGame]),
	})
}

func NewGameScene(w engine.World) {
	w.ChangeScene(NewGame())
}
