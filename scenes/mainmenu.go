package scenes

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sedyh/mizu/pkg/engine"
	"github.com/tomknightdev/dwarven-fortresses/assets"
	"github.com/tomknightdev/dwarven-fortresses/components"
	"github.com/tomknightdev/dwarven-fortresses/entities"
	"github.com/tomknightdev/dwarven-fortresses/enums"
	"github.com/tomknightdev/dwarven-fortresses/gui"
	"github.com/tomknightdev/dwarven-fortresses/systems"
	"github.com/yohamta/furex/v2"
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

	//furex.Debug = true

	ui := func(screen *ebiten.Image) {
		width, height := screen.Size()
		view := &furex.View{
			Width:      width,
			Height:     height,
			Direction:  furex.Column,
			Justify:    furex.JustifySpaceBetween,
			AlignItems: furex.AlignItemCenter,
		}

		view.AddChild(
			&furex.View{
				Width:     500,
				Height:    100,
				MarginTop: 50,
				Handler: &gui.Label{
					Text: "Dwarven Fortress",
					Font: assets.MainFont,
				},
			},
			(&furex.View{
				Width:      500,
				Height:     500,
				Direction:  furex.Column,
				Justify:    furex.JustifyCenter,
				AlignItems: furex.AlignItemCenter,
			}).AddChild(
				&furex.View{
					Width:  200,
					Height: 200,
					Handler: &gui.Image{
						Image: assets.OpaqueImages[enums.TileTypeNewGame],
						Scale: 10.0,
					},
				},
				&furex.View{
					Width:  500,
					Height: 100,
					Handler: &gui.Label{
						Text: "Play",
						Font: assets.MainFont,
					},
				},
			),
			&furex.View{
				Width:  500,
				Height: 100,
				Handler: &gui.Label{
					Text: "By Tom Knight and Leigh Lawley",
					Font: assets.MainFont,
				},
			},
		)

		view.Update()
		view.Draw(screen)
	}
	w.AddEntities(
		// 	&entities.Gui{
		// 		Gui:    components.NewGui(0, 0, enums.GuiPositionCenter, 10.0, enums.GuiActionNewGame),
		// 		Sprite: components.NewSprite(assets.OpaqueImages[enums.TileTypeNewGame]),
		// 	},
		&entities.GuiText{
			Gui: components.Gui{
				UIUpdate: ui,
			},
			Text: components.Text{Content: "Hi"},
		},
		// 	&entities.GuiText{
		// 		Gui:  components.NewGui(0, 120, enums.GuiPositionCenter, 10.0, enums.GuiActionNewGame),
		// 		Text: components.NewText("Play"),
		// 	},
		// 	&entities.GuiText{
		// 		Gui:  components.NewGui(0, 20, enums.GuiPositionBottom, 10.0, enums.GuiActionNewGame),
		// 		Text: components.NewText("By Tom Knight & Leigh Lawley"),
		// 	},
		// 	&entities.GuiText{
		// 		Gui:  components.NewGui(0, 0, enums.GuiPositionRight, 10.0, enums.GuiActionNewGame),
		// 		Text: components.NewText(">"),
		// 	},
		// 	&entities.GuiText{
		// 		Gui:  components.NewGui(0, 0, enums.GuiPositionLeft, 10.0, enums.GuiActionNewGame),
		// 		Text: components.NewText("<"),
	)
}

func NewGameScene(w engine.World) {
	w.ChangeScene(NewGame())
}
