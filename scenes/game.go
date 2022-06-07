package scenes

import (
	"github.com/hajimehoshi/ebiten/v2"
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
		Sprite:   components.NewSprite(assets.OpaqueImages[enums.TileTypeEmpty]),
		Mouse:    components.NewMouse(),
	})

	// Camera
	w.AddEntities(&entities.Camera{
		Zoom:     components.NewZoom(),
		Position: components.NewPosition(0, 0, 5),
	})

	setupGui(w)
}

func setupGui(w engine.World) {
	w.AddEntities(&entities.Gui{
		Gui:    components.NewGui(10, 200, 3.0, enums.GuiActionStair),
		Sprite: components.NewSprite(assets.OpaqueImages[enums.TileTypeStairDown]),
	})
	w.AddEntities(&entities.Gui{
		Gui:    components.NewGui(10, 250, 3.0, enums.GuiActionChop),
		Sprite: components.NewSprite(assets.OpaqueImages[enums.TileTypeTree0]),
	})
	w.AddEntities(&entities.Gui{
		Gui:    components.NewGui(10, 300, 3.0, enums.GuiActionMine),
		Sprite: components.NewSprite(assets.OpaqueImages[enums.TileTypePickaxe]),
	})
	w.AddEntities(&entities.Gui{
		Gui:    components.NewGui(10, 350, 3.0, enums.GuiActionStockpile),
		Sprite: components.NewSprite(assets.OpaqueImages[enums.TileTypeStockpile]),
	})
}
