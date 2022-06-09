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
	"github.com/tomknightdev/dwarven-fortresses/worldgen"
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

func generateWorld(w engine.World, gms *components.GameMapSingleton) {
	// Setup world tiles
	for z := 1; z <= assets.WorldLevels; z++ {
		g := paths.NewGrid(assets.WorldWidth, assets.WorldHeight, assets.TileSize, assets.TileSize)
		for x := 0; x < assets.WorldWidth; x++ {
			for y := 0; y < assets.WorldHeight; y++ {
				c := g.Get(x, y)
				t := components.TileInfo{
					Position: components.NewPosition(x, y, z),
				}

				if z == 5 {
					t.TileType = components.NewTileType(enums.TileTypeTerrain, enums.TileSpriteGrass0)
					t.Image = assets.OpaqueImages[enums.TileSpriteGrass0]
				} else if z < 5 {
					t.TileType = components.NewTileType(enums.TileTypeRock, enums.TileSpriteRock)
					// t.Image = assets.Images["rock"]
					c.Walkable = false
				} else {
					t.TileType = components.NewTileType(enums.TileTypeRock, enums.TileSpriteEmpty)
					c.Walkable = false
				}

				gms.Tiles[t.Position] = t
			}
		}
		gms.Grids[z] = g
	}

	// Setup resource tiles
	rand.Seed(time.Now().UnixNano())

	wG := worldgen.New()
	wG.Octaves = 5
	wG.Scale = 5.0
	wG.Lacunarity = 2.0
	wG.Persistance = 0.1

	for _, tile := range gms.Tiles {
		if tile.TileSpriteEnum != enums.TileSpriteGrass0 {
			continue
		}

		var resourceType enums.TileSpriteEnum
		resourceSeed := wG.GenerateXYTile(tile.X, tile.Y)
		switch int(resourceSeed) {
		case 0:
			continue
		case 1:
			resourceType = enums.TileSpriteTree0
		}
		g := gms.Grids[tile.Z]
		c := g.Get(tile.X, tile.Y)
		c.Walkable = false

		t := components.TileInfo{
			Position: components.NewPosition(tile.X, tile.Y, tile.Z),
			TileType: components.NewTileType(enums.TileTypeResource, enums.TileSpriteTree0),
			Sprite:   components.NewSprite(assets.OpaqueImages[resourceType]),
		}

		gms.Resources[tile.Position] = append(gms.Resources[tile.Position], t)
	}

	for z := 0; z < assets.WorldLevels; z++ {
		// Tiles
		tmImage := ebiten.NewImage(assets.WorldWidth*assets.TileSize, assets.WorldHeight*assets.TileSize)
		for _, t := range gms.Tiles {
			if t.Z != z {
				continue
			}
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
		for _, r := range gms.Resources {
			for _, ti := range r {
				if ti.Z != z {
					continue
				}

				w.AddEntities(&entities.Tree{
					Sprite:    ti.Sprite,
					Position:  ti.Position,
					Resource:  components.NewResource(),
					Choppable: components.NewChoppable(),
					Drops:     components.NewDrops(enums.DropTypeLog, 3),
					Nature:    components.NewNature(),
				})
			}
		}
	}
}
