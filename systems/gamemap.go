package systems

import (
	"math/rand"
	"time"

	"github.com/OpenSauce/paths"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sedyh/mizu/pkg/engine"
	"github.com/tomknightdev/dwarven-fortresses/assets"
	"github.com/tomknightdev/dwarven-fortresses/components"
	"github.com/tomknightdev/dwarven-fortresses/entities"
	"github.com/tomknightdev/dwarven-fortresses/enums"
)

type GameMap struct {
}

func NewGameMap() *GameMap {
	return &GameMap{}
}

func (gm *GameMap) Update(w engine.World) {
	gms, found := w.View(components.GameMapSingleton{}).Get()
	if !found {
		panic("game map singleton not found")
	}

	var gmComp *components.GameMapSingleton
	gms.Get(&gmComp)

	if !gmComp.WorldGenerated {
		generateWorld(w, gmComp)
		gmComp.WorldGenerated = true
	}
}

func (gm *GameMap) Draw(w engine.World, screen *ebiten.Image) {
	gms, found := w.View(components.GameMapSingleton{}).Get()
	if !found {
		panic("game map singleton not found")
	}

	var gmComp *components.GameMapSingleton
	gms.Get(&gmComp)

	// Camera
	camera, found := w.View(components.Zoom{}, components.Position{}).Get()
	if !found {
		return
	}
	var zoom *components.Zoom
	var camPos *components.Position
	camera.Get(&zoom, &camPos)

	// Entities with position and sprite components
	ents := w.View(components.Position{}, components.Sprite{}, components.TileMap{})
	ents.Each(func(e engine.Entity) {
		var pos *components.Position
		var spr *components.Sprite
		e.Get(&pos, &spr)
		op := &ebiten.DrawImageOptions{}

		if camPos.Z > 5 {
			if pos.Z < 5 {
				return
			}

			op.ColorM.Scale(1, 1, 1, 0.5)

		} else if pos.Z != camPos.Z {
			return
		}

		op.GeoM.Translate(float64(pos.X*assets.TileSize), float64(pos.Y*assets.TileSize))
		gmComp.OffScreen.DrawImage(spr.Image, op)
	})

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(zoom.Value, zoom.Value)

	ww, wh := ebiten.WindowSize()
	op.GeoM.Translate(-float64(camPos.X-(ww/2)), -float64(camPos.Y-(wh/2)))
	// op.Filter = ebiten.FilterNearest
	screen.DrawImage(gmComp.OffScreen, op)
	gmComp.OffScreen.Clear()
}

// func UpdateTile(w engine.World, fromTileType, newTileType enums.TileTypeEnum, tileByTypeIndex int, gmComp *components.GameMapSingleton) {
// 	tile := gmComp.TilesByType[fromTileType][tileByTypeIndex]
// 	tileMap := w.View(components.TileMap{}, components.Sprite{}, components.Position{}).Filter()
// 	rand.Seed(time.Now().UnixNano())

// 	for _, tm := range tileMap {
// 		var tmPos *components.Position
// 		var tmSprite *components.Sprite

// 		tm.Get(&tmPos, &tmSprite)
// 		if tmPos.Z == tile.Z {
// 			op := &ebiten.DrawImageOptions{}
// 			op.GeoM.Translate(float64(tile.X*assets.CellSize), float64(tile.Y*assets.CellSize))

// 			switch newTileType {
// 			case enums.TileTypeGrass:
// 				r := rand.Intn(3)
// 				tmSprite.Image.DrawImage(assets.Images[fmt.Sprintf("grass%d", r)], op)
// 			case enums.TileTypeRockFloor:
// 				tmSprite.Image.DrawImage(assets.Images["rockfloor"], op)
// 				cell := gmComp.Grids[tile.Z].Get(tile.X, tile.Y)
// 				cell.Walkable = true
// 				updateAdjacentTiles(w, gmComp, tile, enums.TileTypeRockFloor)
// 			case enums.TileTypeRock:
// 				tmSprite.Image.DrawImage(assets.Images["rock"], op)
// 			}
// 			// Update maps
// 			if fromTileType != newTileType {
// 				gmComp.TilesByType[fromTileType] = append(gmComp.TilesByType[fromTileType][:tileByTypeIndex], gmComp.TilesByType[fromTileType][tileByTypeIndex+1:]...)
// 				gmComp.TilesByType[newTileType] = append(gmComp.TilesByType[newTileType], tile)
// 			}
// 			break
// 		}
// 	}
// }

// func updateAdjacentTiles(w engine.World, gmComp *components.GameMapSingleton, tile components.Position, centreTileType enums.TileTypeEnum) {
// 	for x := -1; x <= 1; x++ {
// 		for y := -1; y <= 1; y++ {
// 			if x == 0 && y == 0 {
// 				continue
// 			}

// 			currentTile := components.NewPosition(tile.X+x, tile.Y+y, tile.Z)

// 			if currentTile.X < 0 || currentTile.Y < 0 {
// 				continue
// 			}

// 			cell := gmComp.Grids[currentTile.Z].Get(currentTile.X, currentTile.Y)
// 			if cell.Walkable {
// 				continue
// 			}

// 			if centreTileType == enums.TileTypeRockFloor {
// 				index, err := helpers.GetTileByTypeIndexFromPos(currentTile, gmComp.TilesByType[enums.TileTypeRock])
// 				if err != nil {
// 					log.Printf("failed to find index for %v at %v\n", enums.TileTypeRock, currentTile)
// 				}
// 				UpdateTile(w, enums.TileTypeRock, enums.TileTypeRock, index, gmComp)
// 			}
// 		}
// 	}
// }

func generateWorld(w engine.World, gms *components.GameMapSingleton) {
	// Setup world tiles
	for z := 1; z <= assets.WorldLevels; z++ {
		g := paths.NewGrid(assets.WorldWidth, assets.WorldHeight, assets.TileSize, assets.TileSize)
		for x := 0; x < assets.WorldWidth; x++ {
			for y := 0; y < assets.WorldHeight; y++ {
				c := g.Get(x, y)
				t := struct {
					components.Position
					components.TileType
					components.Sprite
				}{
					Position: components.NewPosition(x, y, z),
				}

				if z == 5 {
					t.TileType = components.NewTileType(enums.TileTypeGrass0)
					t.Image = assets.OpaqueImages[enums.TileTypeGrass0]
				} else if z < 5 {
					t.TileType = components.NewTileType(enums.TileTypeRockWallHz)
					// t.Image = assets.Images["rock"]
					c.Walkable = false
				} else {
					t.TileType = components.NewTileType(enums.TileTypeEmpty)
					c.Walkable = false
				}

				gms.TilesByZ[z] = append(gms.TilesByZ[z], t)
				gms.TilesByType[t.TileTypeEnum] = append(gms.TilesByType[t.TileTypeEnum], t.Position)
			}
		}
		gms.Grids[z] = g
	}

	// Setup resource tiles
	rand.Seed(time.Now().UnixNano())

	for _, tile := range gms.TilesByType[enums.TileTypeGrass0] {
		if rand.Intn(100) < 5 {
			g := gms.Grids[tile.Z]
			c := g.Get(tile.X, tile.Y)
			c.Walkable = false

			t := struct {
				components.Position
				components.Sprite
			}{
				Position: components.NewPosition(tile.X, tile.Y, tile.Z),
				Sprite:   components.NewSprite(assets.OpaqueImages[enums.TileTypeTree0]),
			}

			gms.ResourcesByZ[tile.Z] = append(gms.ResourcesByZ[tile.Z], t)
		}
	}

	for z := 0; z < assets.WorldLevels; z++ {
		// Tiles
		tmImage := ebiten.NewImage(assets.WorldWidth*assets.TileSize, assets.WorldHeight*assets.TileSize)
		for _, t := range gms.TilesByZ[z] {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(t.X*assets.TileSize), float64(t.Y*assets.TileSize))

			if t.Image != nil {
				tmImage.DrawImage(t.Image, op)
			}
		}
		w.AddEntities(&entities.TileMap{
			Sprite:   components.NewSprite(tmImage),
			Position: components.NewPosition(0, 0, z),
			TileMap:  components.NewTileMap(),
		})

		// Resources
		for _, r := range gms.ResourcesByZ[z] {
			w.AddEntities(&entities.Tree{
				Sprite:    r.Sprite,
				Position:  r.Position,
				Resource:  components.NewResource(),
				Choppable: components.NewChoppable(),
				Drops:     components.NewDrops(enums.DropTypeLog, 3),
				Nature:    components.NewNature(),
			})
		}
	}
}
