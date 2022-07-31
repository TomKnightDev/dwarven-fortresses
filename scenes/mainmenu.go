package scenes

import (
	"image/color"

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
		components.Flex{},
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

	// Debug the GUI to show component boundaries
	// furex.Debug = true

	view := &furex.View{
		Width:      0,
		Height:     0,
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
				Text:  "Dwarven Fortress",
				Font:  assets.MainFont36,
				Color: color.RGBA{0xff, 0xff, 0, 0xff},
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
					Button: gui.Button{
						OnClick: func() {
							w.ChangeScene(NewGame())
						},
					},
				},
			},
			&furex.View{
				Width:  500,
				Height: 100,
				Handler: &gui.Label{
					Text:  "Play",
					Font:  assets.MainFont24,
					Color: color.White,
					Button: gui.Button{
						OnClick: func() {
							w.ChangeScene(NewGame())
						},
					},
				},
			},
		),
		&furex.View{
			Width:  500,
			Height: 100,
			Handler: &gui.Label{
				Text:  "By Tom Knight and Leigh Lawley",
				Color: color.White,
				Font:  assets.MainFont12,
			},
		},
	)

	w.AddEntities(
		&entities.GuiFlex{
			Flex: components.Flex{View: view},
		},
	)
}

func NewGameScene(w engine.World) {
	w.ChangeScene(NewGame())
}
