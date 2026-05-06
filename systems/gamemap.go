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

func generateWorld(w engine.World, gms *components.GameMapSingleton) {
	rand.Seed(time.Now().UnixNano())

	treeGen := worldgen.New()
	treeGen.Octaves = 5
	treeGen.Scale = 5.0
	treeGen.Lacunarity = 2.0
	treeGen.Persistance = 0.1

	// Pre-compute water tiles for the ground level.
	waterTiles := buildWaterTiles()

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
					if waterTiles[components.NewPosition(x, y, z)] {
						tile.TileTypeEnum = enums.TileTypeWater
						tile.Image = assets.OpaqueImages[enums.TileTypeWater]
						c.Walkable = false
					} else {
						tile.TileTypeEnum = enums.TileTypeGrass0
						tile.Image = assets.OpaqueImages[enums.TileTypeGrass0]

						resourceSeed := treeGen.GenerateXYTile(tile.X, tile.Y)
						if int(resourceSeed) == 1 {
							tile.Resources = append(tile.Resources, components.NewResource(enums.ResourceTypeTree))
							c.Walkable = false
						}
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

}

// buildWaterTiles generates the set of ground-level positions that should be
// water. It produces a river crossing the map from top to bottom with a gentle
// meander, plus scattered lakes and ponds using a second noise field.
func buildWaterTiles() map[components.Position]bool {
	water := make(map[components.Position]bool)

	// --- Lakes and ponds ---
	// Large-scale noise (Scale=50) produces big blobs; low values become water.
	lakeGen := worldgen.New()
	lakeGen.Octaves = 2
	lakeGen.Scale = 50.0
	lakeGen.Lacunarity = 2.0
	lakeGen.Persistance = 0.5

	for x := 0; x < assets.WorldWidth; x++ {
		for y := 0; y < assets.WorldHeight; y++ {
			v := lakeGen.GenerateXYTile(x, y)
			if v < 0.45 {
				water[components.NewPosition(x, y, assets.Groundlevel)] = true
			}
		}
	}

	// --- River ---
	// Starts at a random point on the top edge and walks to the bottom,
	// meandering left/right guided by a noise field.
	riverGen := worldgen.New()
	riverGen.Octaves = 3
	riverGen.Scale = 30.0
	riverGen.Lacunarity = 2.0
	riverGen.Persistance = 0.5

	rx := assets.WorldWidth/4 + rand.Intn(assets.WorldWidth/2)

	for y := 0; y < assets.WorldHeight; y++ {
		v := riverGen.GenerateXYTile(rx, y)
		// v is roughly in [0, 2]; shift to ~[-1, 1] then scale
		delta := int((v - 1.0) * 2.5)
		rx += delta
		if rx < 1 {
			rx = 1
		}
		if rx >= assets.WorldWidth-1 {
			rx = assets.WorldWidth - 2
		}

		// Carve a 3-tile-wide channel
		for dx := -1; dx <= 1; dx++ {
			nx := rx + dx
			if nx >= 0 && nx < assets.WorldWidth {
				water[components.NewPosition(nx, y, assets.Groundlevel)] = true
			}
		}
	}

	return water
}
