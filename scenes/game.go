package scenes

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/sedyh/mizu/pkg/engine"
	"github.com/tomknightdev/dwarven-fortresses/assets"
	"github.com/tomknightdev/dwarven-fortresses/components"
	"github.com/tomknightdev/dwarven-fortresses/entities"
	"github.com/tomknightdev/dwarven-fortresses/enums"
	"github.com/tomknightdev/dwarven-fortresses/systems"
)

type Game struct {
}

func NewGame() *Game {
	return &Game{}
}

func (g *Game) Setup(w engine.World) {
	w.AddComponents(
		components.InputSingleton{},
		components.GameMapSingleton{},
		components.Position{},
		components.Sprite{},
		components.Move{},
		components.Mouse{},
		components.Zoom{},
		components.TileType{},
		components.Task{},
		components.Job{},
		components.Worker{},
		components.TileMap{},
		components.Gui{},
		components.Resource{},
		components.Choppable{},
		components.Drops{},
		components.NatureSingleton{},
		components.Nature{},
		components.Item{},
		components.Building{},
		components.Designation{},
		components.Inventory{},
	)

	w.AddSystems(
		systems.NewInput(),
		systems.NewGameMap(),
		systems.NewCamera(),
		systems.NewMouse(),
		systems.NewPathfinder(),
		systems.NewGui(),
		systems.NewNature(),
		systems.NewBuilding(),
		systems.NewDesignations(),
		systems.NewItem(),
		systems.NewJob(),
		systems.NewActor(),
		systems.NewTask(),
		systems.NewDebug(),
	)

	// Admin entity
	w.AddEntities(&entities.Admin{
		InputSingleton:   components.NewInputSingleton(),
		GameMapSingleton: components.NewGameMapSingleton(),
		NatureSingleton:  components.NewNatureSingleton(),
	})

	// Actors
	for i := 0; i < assets.StartingDwarfCount; i++ {
		w.AddEntities(&entities.Actor{
			Position:  components.NewPosition(1, 1, 5),
			Sprite:    components.NewSprite(assets.TransImages[enums.TileTypeDwarf]),
			Move:      components.NewMove(1, 1, 5),
			Worker:    components.NewWorker(),
			Inventory: components.NewInventory(),
		})
	}

	// Input
	cx, cy := ebiten.CursorPosition()
	w.AddEntities(&entities.Mouse{
		Position: components.NewPosition(cx, cy, 5),
		Sprite:   components.NewSprite(assets.Images[enums.TileTypeEmpty]),
		Mouse:    components.NewMouse(),
	})

	// Camera
	w.AddEntities(&entities.Camera{
		Zoom:     components.NewZoom(),
		Position: components.NewPosition(0, 0, 5),
	})

	setupGui(w)
	setupAudio()
}

func setupGui(w engine.World) {
	w.AddEntities(&entities.Gui{
		Gui:    components.NewGui(10, 200, 3.0, enums.GuiActionStair),
		Sprite: components.NewSprite(assets.Images[enums.TileTypeStairDown]),
	})
	w.AddEntities(&entities.Gui{
		Gui:    components.NewGui(10, 250, 3.0, enums.GuiActionChop),
		Sprite: components.NewSprite(assets.Images[enums.TileTypeTree0]),
	})
	w.AddEntities(&entities.Gui{
		Gui:    components.NewGui(10, 300, 3.0, enums.GuiActionMine),
		Sprite: components.NewSprite(assets.Images[enums.TileTypePickaxe]),
	})
	w.AddEntities(&entities.Gui{
		Gui:    components.NewGui(10, 350, 3.0, enums.GuiActionStockpile),
		Sprite: components.NewSprite(assets.Images[enums.TileTypeStockpile]),
	})
}

func setupAudio() {
	audioContext := audio.NewContext(44100)

	p, err := audioContext.NewPlayer(assets.MainAudio)
	if err != nil {
		log.Fatal(err)
	}

	// TODO Be able to set the volume in the settings
	p.SetVolume(0.2)

	p.Play()
}
