package scenes

import (
	"fmt"
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

	startGame := func() {
		w.ChangeScene(NewGame())
	}

	// furex.Debug = true

	white := color.White
	yellow := color.RGBA{0xff, 0xff, 0, 0xff}
	dim := color.RGBA{0xaa, 0xaa, 0xaa, 0xff}

	// Helper to build a [-] / [+] row for a setting.
	makeSettingRow := func(name string, getText func() string, onMinus, onPlus func()) *furex.View {
		return (&furex.View{
			Width:      480,
			Height:     56,
			Direction:  furex.Row,
			Justify:    furex.JustifySpaceBetween,
			AlignItems: furex.AlignItemCenter,
		}).AddChild(
			// Setting name
			&furex.View{
				Width:  180,
				Height: 56,
				Handler: &gui.Label{
					Text:  name,
					Font:  assets.MainFont24,
					Color: white,
				},
			},
			// Minus button
			&furex.View{
				Width:  56,
				Height: 56,
				Handler: &gui.Label{
					Text:  "-",
					Font:  assets.MainFont24,
					Color: yellow,
					Button: gui.Button{OnClick: onMinus},
				},
			},
			// Current value
			&furex.View{
				Width:  120,
				Height: 56,
				Handler: &gui.Label{
					GetText: getText,
					Font:    assets.MainFont24,
					Color:   white,
				},
			},
			// Plus button
			&furex.View{
				Width:  56,
				Height: 56,
				Handler: &gui.Label{
					Text:  "+",
					Font:  assets.MainFont24,
					Color: yellow,
					Button: gui.Button{OnClick: onPlus},
				},
			},
		)
	}

	dwarfRow := makeSettingRow(
		"Dwarves",
		func() string { return fmt.Sprintf("%d", assets.Settings.DwarfCount) },
		func() {
			if assets.Settings.DwarfCount > 1 {
				assets.Settings.DwarfCount--
			}
		},
		func() {
			if assets.Settings.DwarfCount < 20 {
				assets.Settings.DwarfCount++
			}
		},
	)

	mapRow := makeSettingRow(
		"Map Size",
		func() string { return assets.Settings.MapSizeName() },
		func() {
			if assets.Settings.MapSizeIndex > 0 {
				assets.Settings.MapSizeIndex--
			}
		},
		func() {
			if assets.Settings.MapSizeIndex < assets.Settings.MapSizeCount()-1 {
				assets.Settings.MapSizeIndex++
			}
		},
	)

	view := &furex.View{
		Width:      0,
		Height:     0,
		Direction:  furex.Column,
		Justify:    furex.JustifySpaceBetween,
		AlignItems: furex.AlignItemCenter,
	}

	view.AddChild(
		// Title
		&furex.View{
			Width:     500,
			Height:    100,
			MarginTop: 50,
			Handler: &gui.Label{
				Text:  "Dwarven Fortress",
				Font:  assets.MainFont36,
				Color: yellow,
			},
		},

		// Centre section: settings + play
		(&furex.View{
			Width:      500,
			Height:     460,
			Direction:  furex.Column,
			Justify:    furex.JustifyCenter,
			AlignItems: furex.AlignItemCenter,
		}).AddChild(
			// Settings panel
			(&furex.View{
				Width:      500,
				Height:     140,
				Direction:  furex.Column,
				Justify:    furex.JustifyCenter,
				AlignItems: furex.AlignItemCenter,
			}).AddChild(dwarfRow, mapRow),

			// Play image
			&furex.View{
				Width:  160,
				Height: 160,
				Handler: &gui.Image{
					Image: assets.OpaqueImages[enums.TileTypeNewGame],
					Scale: 8.0,
					Button: gui.Button{OnClick: startGame},
				},
			},

			// Play label
			&furex.View{
				Width:  500,
				Height: 60,
				Handler: &gui.Label{
					Text:  "Play",
					Font:  assets.MainFont24,
					Color: white,
					Button: gui.Button{OnClick: startGame},
				},
			},
		),

		// Footer
		&furex.View{
			Width:  500,
			Height: 60,
			Handler: &gui.Label{
				Text:  "By Tom Knight and Leigh Lawley",
				Color: dim,
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
