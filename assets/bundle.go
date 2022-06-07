package assets

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/tomknightdev/dwarven-fortresses/enums"
)

// This package is for loading all the images and storing world information

const (
	WorldWidth  = 200
	WorldHeight = 200
	WorldLevels = 10

	StartingDwarfCount = 7
)

var (
	CellSize int

	// TilesetFileNames []string

	// Tileset       []byte
	TilesetImages = make(map[string]*ebiten.Image)

	OpaqueImages = make(map[enums.TileTypeEnum]*ebiten.Image)
	TransImages  = make(map[enums.TileTypeEnum]*ebiten.Image)
)

type TilesetDefinition struct {
	TilesetName      string   `json:"tilesetName"`
	TilesetFileNames []string `json:"tilesetFileNames"`
	TileSize         int      `json:"tileSize"`
	Tiles            []Tile   `json:"tiles"`
}

type Tile struct {
	Name                string `json:"name"`
	Id                  int    `json:"id`
	OpaqueFileName      string `json:"opaqueFileName"`
	TransparentFileName string `json:"transparentFileName"`
	X                   int    `json:"x"`
	Y                   int    `json:"y"`
}

func init() {
	abs, err := filepath.Abs("./assets/images/tileset.json")
	if err == nil {
		fmt.Println("Absolute:", abs)
	}

	fmt.Println(abs)
	// Read tileset json file
	tsd, err := ioutil.ReadFile(abs)
	if err != nil {
		log.Fatal(err)
	}

	tilesetDef := TilesetDefinition{}
	err = json.Unmarshal(tsd, &tilesetDef)
	if err != nil {
		log.Fatal(err)
	}

	for _, ts := range tilesetDef.TilesetFileNames {
		// Read tile sets from definition
		abs, err = filepath.Abs("./assets/images/" + ts)
		if err == nil {
			fmt.Println("Absolute:", abs)
		}

		tileset, err := os.ReadFile(abs)
		if err != nil {
			log.Fatal(err)
		}

		img, err := png.Decode(bytes.NewReader(tileset))
		if err != nil {
			log.Fatal(err)
		}
		TilesetImages[ts] = ebiten.NewImageFromImage(img)
	}

	// abs, err = filepath.Abs("./assets/images/" + tilesetDef.TransparentTilesetFileName)
	// if err == nil {
	// 	fmt.Println("Absolute:", abs)
	// }
	// TransTileset, err := os.ReadFile(abs)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// img, err = png.Decode(bytes.NewReader(TransTileset))
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// TransTilesetImage = ebiten.NewImageFromImage(img)

	CellSize = tilesetDef.TileSize

	for _, t := range tilesetDef.Tiles {
		OpaqueImages[enums.TileTypeEnum(t.Id)] = TilesetImages[t.OpaqueFileName].SubImage(image.Rect(t.X*CellSize, t.Y*CellSize, (t.X+1)*CellSize, (t.Y+1)*CellSize)).(*ebiten.Image)
		TransImages[enums.TileTypeEnum(t.Id)] = TilesetImages[t.TransparentFileName].SubImage(image.Rect(t.X*CellSize, t.Y*CellSize, (t.X+1)*CellSize, (t.Y+1)*CellSize)).(*ebiten.Image)
	}

	// Images["empty"] = TransTilesetImage.SubImage(image.Rect(0*CellSize, 0*CellSize, 1*CellSize, 1*CellSize)).(*ebiten.Image)
	// Images["cursor"] = TransTilesetImage.SubImage(image.Rect(28*CellSize, 14*CellSize, 29*CellSize, 15*CellSize)).(*ebiten.Image)

	// Images["stockpile"] = TransTilesetImage.SubImage(image.Rect(29*CellSize, 14*CellSize, 30*CellSize, 15*CellSize)).(*ebiten.Image)

	// Images["rock"] = TilesetImage.SubImage(image.Rect(19*CellSize, 1*CellSize, 20*CellSize, 2*CellSize)).(*ebiten.Image)
	// Images["rockfloor"] = TilesetImage.SubImage(image.Rect(2*CellSize, 0*CellSize, 3*CellSize, 1*CellSize)).(*ebiten.Image)
	// Images["rocks"] = TransTilesetImage.SubImage(image.Rect(5*CellSize, 2*CellSize, 6*CellSize, 3*CellSize)).(*ebiten.Image)

	// Images["dirt0"] = TilesetImage.SubImage(image.Rect(0*CellSize, 0*CellSize, 1*CellSize, 1*CellSize)).(*ebiten.Image)
	// Images["dirt1"] = TilesetImage.SubImage(image.Rect(1*CellSize, 0*CellSize, 2*CellSize, 1*CellSize)).(*ebiten.Image)

	// Images["grass0"] = TilesetImage.SubImage(image.Rect(5*CellSize, 0*CellSize, 6*CellSize, 1*CellSize)).(*ebiten.Image)
	// Images["grass1"] = TilesetImage.SubImage(image.Rect(6*CellSize, 0*CellSize, 7*CellSize, 1*CellSize)).(*ebiten.Image)
	// Images["grass2"] = TilesetImage.SubImage(image.Rect(7*CellSize, 0*CellSize, 8*CellSize, 1*CellSize)).(*ebiten.Image)

	// Images["tree0"] = TransTilesetImage.SubImage(image.Rect(0*CellSize, 1*CellSize, 1*CellSize, 2*CellSize)).(*ebiten.Image)
	// Images["tree1"] = TransTilesetImage.SubImage(image.Rect(1*CellSize, 1*CellSize, 2*CellSize, 2*CellSize)).(*ebiten.Image)
	// Images["tree2"] = TransTilesetImage.SubImage(image.Rect(2*CellSize, 1*CellSize, 3*CellSize, 2*CellSize)).(*ebiten.Image)
	// Images["tree3"] = TransTilesetImage.SubImage(image.Rect(3*CellSize, 1*CellSize, 4*CellSize, 2*CellSize)).(*ebiten.Image)
	// Images["tree4"] = TransTilesetImage.SubImage(image.Rect(4*CellSize, 1*CellSize, 5*CellSize, 2*CellSize)).(*ebiten.Image)
	// Images["tree5"] = TransTilesetImage.SubImage(image.Rect(5*CellSize, 1*CellSize, 6*CellSize, 2*CellSize)).(*ebiten.Image)

	// Images["log0"] = TransTilesetImage.SubImage(image.Rect(20*CellSize, 6*CellSize, 21*CellSize, 7*CellSize)).(*ebiten.Image)

	// Images["water"] = TransTilesetImage.SubImage(image.Rect(14*CellSize, 5*CellSize, 15*CellSize, 6*CellSize)).(*ebiten.Image)

	// Images["dwarf"] = TransTilesetImage.SubImage(image.Rect(25*CellSize, 0*CellSize, 26*CellSize, 1*CellSize)).(*ebiten.Image)

	// Images["stairdown"] = TilesetImage.SubImage(image.Rect(3*CellSize, 6*CellSize, 4*CellSize, 7*CellSize)).(*ebiten.Image)
	// Images["stairup"] = TilesetImage.SubImage(image.Rect(2*CellSize, 6*CellSize, 3*CellSize, 7*CellSize)).(*ebiten.Image)

	// Images["pickaxe"] = TilesetImage.SubImage(image.Rect(11*CellSize, 27*CellSize, 12*CellSize, 28*CellSize)).(*ebiten.Image)

}
