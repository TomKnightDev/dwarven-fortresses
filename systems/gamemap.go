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
	"github.com/tomknightdev/dwarven-fortresses/helpers"
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
	ents := w.View(components.Position{}, components.Sprite{}, components.TileMap{})
	helpers.DrawImages(w, ents)
}

var grassVariants = []enums.TileTypeEnum{
	enums.TileTypeGrass0, enums.TileTypeGrass0, enums.TileTypeGrass0,
	enums.TileTypeGrass0, enums.TileTypeGrass0,
	enums.TileTypeGrass1, enums.TileTypeGrass1,
	enums.TileTypeGrass2,
}

func generateWorld(w engine.World, gms *components.GameMapSingleton) {
	rand.Seed(time.Now().UnixNano())

	// Large-scale noise defines where forest regions are.
	forestNoise := worldgen.New()
	forestNoise.Octaves = 4
	forestNoise.Scale = 20.0
	forestNoise.Lacunarity = 2.0
	forestNoise.Persistance = 0.5

	// Fine-scale noise breaks up forest edges and creates internal gaps.
	scatterNoise := worldgen.New()
	scatterNoise.Octaves = 3
	scatterNoise.Scale = 4.0
	scatterNoise.Lacunarity = 2.0
	scatterNoise.Persistance = 0.6

	// Setup world tiles
	for z := 1; z <= assets.WorldLevels; z++ {
		g := paths.NewGrid(assets.WorldWidth, assets.WorldHeight, assets.TileSize, assets.TileSize)
		gms.Grids[z] = g

		tmImage := ebiten.NewImage(assets.WorldWidth*assets.TileSize, assets.WorldHeight*assets.TileSize)

		w.AddEntities(&entities.TileMap{
			Sprite:   components.NewSprite(tmImage),
			Position: components.NewPosition(0, 0, z),
			TileMap:  components.NewTileMap(),
		})

		for x := 0; x < assets.WorldWidth; x++ {
			for y := 0; y < assets.WorldHeight; y++ {
				c := g.Get(x, y)

				tile := components.Tile{
					Position: components.NewPosition(x, y, z),
				}

				gms.Tiles[tile.Position] = &tile

				if z == assets.Groundlevel {
					grassType := grassVariants[rand.Intn(len(grassVariants))]
					tile.TileTypeEnum = grassType
					tile.Image = assets.OpaqueImages[grassType]

					forest := forestNoise.GenerateXYTile(x, y)
					scatter := scatterNoise.GenerateXYTile(x, y)

					inForest := forest > 1.15 && scatter > 1.05
					isolated := rand.Float64() < 0.015

					if inForest || isolated {
						tile.Resources = append(tile.Resources, components.NewResource(enums.ResourceTypeTree))
						c.Walkable = false
					}
				} else if z < assets.Groundlevel {
					tile.TileTypeEnum = enums.TileTypeRock
					tile.Image = assets.OpaqueImages[enums.TileTypeRock]
					c.Walkable = false
				} else {
					tile.TileTypeEnum = enums.TileTypeEmpty
					c.Walkable = false
				}

				helpers.RenderTile(w, tile.Position)
			}
		}

	}

	// for z := 0; z < assets.WorldLevels; z++ {
	// 	// Tiles
	// 	tmImage := ebiten.NewImage(assets.WorldWidth*assets.TileSize, assets.WorldHeight*assets.TileSize)
	// 	tiles := gms.GetTileForZ(z)
	// 	for _, t := range tiles {
	// 		op := &ebiten.DrawImageOptions{}
	// 		op.GeoM.Translate(float64(t.X*assets.TileSize), float64(t.Y*assets.TileSize))

	// 		for _, i := range t.Sprites {
	// 			tmImage.DrawImage(i, op)
	// 		}
	// 	}

	// 	// Resources
	// 	for _, r := range gms.ResourcesByZ[z] {
	// 		w.AddEntities(&entities.Tree{
	// 			// Sprite:    r.Sprite,
	// 			Position:  r.Position,
	// 			Resource:  components.NewResource(),
	// 			Choppable: components.NewChoppable(),
	// 			Drops:     components.NewDrops(enums.DropTypeLog, 1),
	// 			Nature:    components.NewNature(),
	// 		})

	// 		op := &ebiten.DrawImageOptions{}
	// 		op.GeoM.Translate(float64(r.X*assets.TileSize), float64(r.Y*assets.TileSize))

	// 		if r.Image != nil {
	// 			tmImage.DrawImage(r.Image, op)
	// 		}
	// 	}

	// 	w.AddEntities(&entities.TileMap{
	// 		Sprite:   components.NewSprite(tmImage),
	// 		Position: components.NewPosition(0, 0, z),
	// 		TileMap:  components.NewTileMap(),
	// 	})

	// }
}
