package helpers

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sedyh/mizu/pkg/engine"
	"github.com/tomknightdev/dwarven-fortresses/assets"
	"github.com/tomknightdev/dwarven-fortresses/components"
	"github.com/tomknightdev/dwarven-fortresses/enums"
)

func GetTileByTypeIndexFromPos(pos components.Position, tilesByType []components.Position) (int, error) {
	for i, t := range tilesByType {
		if t.X == pos.X && t.Y == pos.Y && t.Z == pos.Z {
			return i, nil
		}
	}

	return 0, fmt.Errorf("unable to find tile at %v", pos)
}

func DrawImage(w engine.World, screen *ebiten.Image, pos components.Position, image *ebiten.Image) {
	camera, found := w.View(components.Zoom{}, components.Position{}).Get()
	if !found {
		return
	}
	var zoom *components.Zoom
	var camPos *components.Position
	camera.Get(&zoom, &camPos)

	if pos.Z != camPos.Z {
		return
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(pos.X*assets.TileSize), float64(pos.Y*assets.TileSize))

	op.GeoM.Scale(zoom.Value, zoom.Value)

	ww, wh := ebiten.WindowSize()
	op.GeoM.Translate(-float64(camPos.X-(ww/2)), -float64(camPos.Y-(wh/2)))
	// op.Filter = ebiten.FilterNearest
	screen.DrawImage(image, op)
}

func DrawImages(w engine.World, screen *ebiten.Image, offScreen *ebiten.Image, ents []engine.Entity) {
	camera, found := w.View(components.Zoom{}, components.Position{}).Get()
	if !found {
		return
	}
	var zoom *components.Zoom
	var camPos *components.Position
	camera.Get(&zoom, &camPos)

	var s *components.Sprite
	var p *components.Position
	for _, e := range ents {
		e.Get(&s, &p)

		if p.Z != camPos.Z {
			return
		}

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(p.X*assets.TileSize), float64(p.Y*assets.TileSize))
		offScreen.DrawImage(s.Image, op)
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(zoom.Value, zoom.Value)

	ww, wh := ebiten.WindowSize()
	op.GeoM.Translate(-float64(camPos.X-(ww/2)), -float64(camPos.Y-(wh/2)))
	// op.Filter = ebiten.FilterNearest
	screen.DrawImage(offScreen, op)
	offScreen.Clear()
}

func GetTilesForType(tiles map[components.Position]components.TileInfo, tt enums.TileTypeEnum) []components.TileInfo {
	var ti []components.TileInfo
	for _, t := range tiles {
		if t.TileTypeEnum == tt {
			ti = append(ti, t)
		}
	}

	return ti
}

func UpdateTile(w engine.World, fromTileType, newTileType enums.TileSpriteEnum, tileByTypeIndex int, gmComp *components.GameMapSingleton) {
	tile := gmComp.TilesByType[fromTileType][tileByTypeIndex]
	tileMap := w.View(components.TileMap{}, components.Sprite{}, components.Position{}).Filter()
	rand.Seed(time.Now().UnixNano())

	for _, tm := range tileMap {
		var tmPos *components.Position
		var tmSprite *components.Sprite

		tm.Get(&tmPos, &tmSprite)
		if tmPos.Z == tile.Z {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(tile.X*assets.TileSize), float64(tile.Y*assets.TileSize))

			switch newTileType {
			case enums.TileSpriteGrass0:
				// r := rand.Intn(3)
				tmSprite.Image.DrawImage(assets.OpaqueImages[enums.TileSpriteGrass0], op)
				cell := gmComp.Grids[tile.Z].Get(tile.X, tile.Y)
				cell.Walkable = true
			case enums.TileSpriteRockFloor:
				tmSprite.Image.DrawImage(assets.OpaqueImages[enums.TileSpriteRockFloor], op)
				cell := gmComp.Grids[tile.Z].Get(tile.X, tile.Y)
				cell.Walkable = true
				updateAdjacentTiles(w, gmComp, tile, enums.TileSpriteRockFloor)
			case enums.TileSpriteRock:
				tmSprite.Image.DrawImage(assets.OpaqueImages[enums.TileSpriteRock], op)
			}
			// Update maps
			if fromTileType != newTileType {
				gmComp.TilesByType[fromTileType] = append(gmComp.TilesByType[fromTileType][:tileByTypeIndex], gmComp.TilesByType[fromTileType][tileByTypeIndex+1:]...)
				gmComp.TilesByType[newTileType] = append(gmComp.TilesByType[newTileType], tile)
			}
			break
		}
	}
}

func updateAdjacentTiles(w engine.World, gmComp *components.GameMapSingleton, tile components.Position, centreTileType enums.TileSpriteEnum) {
	for x := -1; x <= 1; x++ {
		for y := -1; y <= 1; y++ {
			if x == 0 && y == 0 {
				continue
			}

			currentTile := components.NewPosition(tile.X+x, tile.Y+y, tile.Z)

			if currentTile.X < 0 || currentTile.Y < 0 {
				continue
			}

			cell := gmComp.Grids[currentTile.Z].Get(currentTile.X, currentTile.Y)
			if cell.Walkable {
				continue
			}

			if centreTileType == enums.TileSpriteRockFloor {
				index, err := GetTileByTypeIndexFromPos(currentTile, gmComp.TilesByType[enums.TileSpriteRock])
				if err != nil {
					log.Printf("failed to find index for %v at %v\n", enums.TileSpriteRock, currentTile)
				}
				UpdateTile(w, enums.TileSpriteRock, enums.TileSpriteRock, index, gmComp)
			}
		}
	}
}
