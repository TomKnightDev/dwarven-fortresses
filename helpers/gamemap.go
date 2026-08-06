package helpers

import (
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sedyh/mizu/pkg/engine"
	"github.com/tomknightdev/dwarven-fortresses/assets"
	"github.com/tomknightdev/dwarven-fortresses/components"
	"github.com/tomknightdev/dwarven-fortresses/enums"
)

var treeVariants = []enums.TileTypeEnum{
	enums.TileTypeTree0,
	enums.TileTypeTree1,
	enums.TileTypeTree2,
	enums.TileTypeTree3,
	enums.TileTypeTree4,
	enums.TileTypeTree5,
}

func getImageForResourceType(rt enums.ResourceTypeEnum) *ebiten.Image {
	switch rt {
	case enums.ResourceTypeTree:
		return assets.OpaqueImages[treeVariants[rand.Intn(len(treeVariants))]]
	default:
		log.Println("resource type not handled: ", rt)
		return nil
	}
}

func getImageForItemType(it enums.ItemTypeEnum) *ebiten.Image {
	switch it {
	case enums.ItemTypeLog:
		return assets.OpaqueImages[enums.TileTypeLog0]
	default:
		log.Println("item type not handled: ", it)
		return nil
	}
}

func getImageForTileType(tt enums.TileTypeEnum) *ebiten.Image {
	switch tt {
	case enums.TileTypeStairDown:
		return assets.OpaqueImages[enums.TileTypeStairDown]
	case enums.TileTypeStairUp:
		return assets.OpaqueImages[enums.TileTypeStairUp]
	default:
		log.Println("tile type not handled: ", tt)
		return nil
	}
}

func GetTileForZ(w engine.World, level int) (tiles []*components.Tile) {
	gms := GetGameMapSingleton(w)
	for p, t := range gms.Tiles {
		if p.Z == level {
			tiles = append(tiles, t)
		}
	}

	return tiles
}

func GetTilesForType(w engine.World, tt enums.TileTypeEnum) (tiles []*components.Tile) {
	gms := GetGameMapSingleton(w)
	for _, t := range gms.Tiles {
		if t.TileTypeEnum == tt {
			tiles = append(tiles, t)
		}
	}

	return tiles
}

func IsTileOfType(w engine.World, pos components.Position, tt enums.TileTypeEnum) bool {
	gms := GetGameMapSingleton(w)
	tile := gms.Tiles[pos]
	return tile.TileTypeEnum == tt
}

func TileHasResource(w engine.World, pos components.Position, rt enums.ResourceTypeEnum) bool {
	gms := GetGameMapSingleton(w)
	tile := gms.Tiles[pos]
	for _, r := range tile.Resources {
		if r.ResourceTypeEnum == rt {
			return true
		}
	}

	return false
}

func RemoveResourceFromTile(w engine.World, pos components.Position, rt enums.ResourceTypeEnum, walkable bool) {
	gms := GetGameMapSingleton(w)
	tile := gms.Tiles[pos]

	for i, r := range tile.Resources {
		if r.ResourceTypeEnum == rt {
			tile.Resources = append(tile.Resources[:i], tile.Resources[i+1:]...)
		}
	}

	cell := gms.Grids[pos.Z].Get(pos.X, pos.Y)
	if cell.Walkable != walkable {
		cell.Walkable = walkable
		MarkRegionDirty(w, gms)
	}

	RenderTile(w, pos)
}

func AddBuildingToTile(w engine.World, pos components.Position, tt enums.TileTypeEnum, walkable bool) {
	gms := GetGameMapSingleton(w)
	tile := gms.Tiles[pos]

	tile.Buildings = append(tile.Buildings, components.NewBuilding(tt))

	cell := gms.Grids[pos.Z].Get(pos.X, pos.Y)
	walkChanged := cell.Walkable != walkable
	cell.Walkable = walkable

	addedStair := false
	switch tt {
	case enums.TileTypeStairDown:
		gms.Downs = append(gms.Downs, pos)
		addedStair = true
	case enums.TileTypeStairUp:
		gms.Ups = append(gms.Ups, pos)
		addedStair = true
	}

	// A new stair changes inter-Z connectivity even if its own cell was
	// already walkable, so always invalidate when one is placed.
	if walkChanged || addedStair {
		MarkRegionDirty(w, gms)
	}

	RenderTile(w, pos)
}

func MineTile(w engine.World, pos components.Position) {
	gms := GetGameMapSingleton(w)
	tile := gms.Tiles[pos]

	if tile.TileTypeEnum != enums.TileTypeRock {
		return
	}

	tile.Image = assets.TransImages[enums.TileTypeRockFloor]
	tile.TileTypeEnum = enums.TileTypeRockFloor

	cell := gms.Grids[pos.Z].Get(pos.X, pos.Y)
	if !cell.Walkable {
		cell.Walkable = true
		MarkRegionDirty(w, gms)
	}

	RenderTile(w, pos)
}

func RenderTile(w engine.World, pos components.Position) {
	gms := GetGameMapSingleton(w)
	tile := gms.Tiles[pos]

	if tile.Image == nil {
		return
	}

	tmEnts := w.View(components.TileMap{}).Filter()
	var sprite *components.Sprite
	var tmpos *components.Position

	for _, tmEnt := range tmEnts {
		tmEnt.Get(&sprite, &tmpos)
		if tmpos.Z != pos.Z {
			continue
		}

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(pos.X*assets.TileSize), float64(pos.Y*assets.TileSize))

		sprite.Image.DrawImage(tile.Image, op)

		for _, r := range tile.Resources {
			sprite.Image.DrawImage(getImageForResourceType(r.ResourceTypeEnum), op)
		}

		for _, b := range tile.Buildings {
			sprite.Image.DrawImage(getImageForTileType(b.TileTypeEnum), op)
		}

		for _, i := range tile.Items {
			sprite.Image.DrawImage(getImageForItemType(i.ItemType), op)
		}
		break
	}
}

func DrawImage(w engine.World, pos components.Position, image *ebiten.Image) {
	camera, found := w.View(components.Zoom{}, components.Position{}).Get()
	if !found {
		return
	}
	var camPos *components.Position
	camera.Get(&camPos)

	if pos.Z != camPos.Z {
		return
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(pos.X*assets.TileSize), float64(pos.Y*assets.TileSize))

	// op.GeoM.Scale(zoom.Value, zoom.Value)

	// ww, wh := ebiten.WindowSize()
	// op.GeoM.Translate(-float64(camPos.X-(ww/2)), -float64(camPos.Y-(wh/2)))
	// op.Filter = ebiten.FilterNearest
	rs := GetRenderSingleton(w)
	rs.OffScreen.DrawImage(image, op)
}

// MovingRenderPosition interpolates from the entity's grid cell toward the
// next cell in its path, based on progress toward the energy cost of the
// current step, so movement renders as a slide rather than a per-cell snap.
func MovingRenderPosition(pos *components.Position, move *components.Move, inv *components.Inventory) (float64, float64) {
	x, y := float64(pos.X), float64(pos.Y)

	if len(move.CurrentPaths) == 0 {
		return x, y
	}

	next := move.CurrentPaths[0].Next()
	if next == nil {
		return x, y
	}

	progress := float64(move.CurrentEnergy) / float64(move.TotalEnergy+inv.Weight)
	if progress > 1 {
		progress = 1
	}

	return x + (float64(next.X)-x)*progress, y + (float64(next.Y)-y)*progress
}

func DrawMovingImages(w engine.World, ents engine.View) {
	rs := GetRenderSingleton(w)

	camera, found := w.View(components.Zoom{}, components.Position{}).Get()
	if !found {
		return
	}
	var camPos *components.Position
	camera.Get(&camPos)

	var s *components.Sprite
	var p *components.Position
	var mv *components.Move
	var inv *components.Inventory

	ents.Each(func(e engine.Entity) {
		e.Get(&s, &p, &mv, &inv)

		if p.Z != camPos.Z || !s.Drawn {
			return
		}

		rx, ry := MovingRenderPosition(p, mv, inv)

		op := &ebiten.DrawImageOptions{}
		if mv.FacingLeft {
			op.GeoM.Scale(-1, 1)
			op.GeoM.Translate(float64(assets.TileSize), 0)
		}
		op.GeoM.Translate(rx*float64(assets.TileSize), ry*float64(assets.TileSize))
		rs.OffScreen.DrawImage(s.Image, op)
	})
}

func DrawImages(w engine.World, ents engine.View) {
	rs := GetRenderSingleton(w)

	camera, found := w.View(components.Zoom{}, components.Position{}).Get()
	if !found {
		return
	}
	// var zoom *components.Zoom
	var camPos *components.Position
	camera.Get(&camPos)

	var s *components.Sprite
	var p *components.Position

	ents.Each(func(e engine.Entity) {
		e.Get(&s, &p)

		if p.Z != camPos.Z || !s.Drawn {
			return
		}

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(p.X*assets.TileSize), float64(p.Y*assets.TileSize))
		rs.OffScreen.DrawImage(s.Image, op)
	})
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
// 			op.GeoM.Translate(float64(tile.X*assets.TileSize), float64(tile.Y*assets.TileSize))

// 			switch newTileType {
// 			case enums.TileTypeGrass0:
// 				// r := rand.Intn(3)
// 				tmSprite.Image.DrawImage(assets.OpaqueImages[enums.TileTypeGrass0], op)
// 				cell := gmComp.Grids[tile.Z].Get(tile.X, tile.Y)
// 				cell.Walkable = true
// 			case enums.TileTypeRockFloor:
// 				tmSprite.Image.DrawImage(assets.OpaqueImages[enums.TileTypeRockFloor], op)
// 				cell := gmComp.Grids[tile.Z].Get(tile.X, tile.Y)
// 				cell.Walkable = true
// 				updateAdjacentTiles(w, gmComp, tile, enums.TileTypeRockFloor)
// 			case enums.TileTypeRock:
// 				tmSprite.Image.DrawImage(assets.OpaqueImages[enums.TileTypeRock], op)
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
// 				index, err := GetTileByTypeIndexFromPos(currentTile, gmComp.TilesByType[enums.TileTypeRock])
// 				if err != nil {
// 					log.Printf("failed to find index for %v at %v\n", enums.TileTypeRock, currentTile)
// 				}
// 				UpdateTile(w, enums.TileTypeRock, enums.TileTypeRock, index, gmComp)
// 			}
// 		}
// 	}
// }
