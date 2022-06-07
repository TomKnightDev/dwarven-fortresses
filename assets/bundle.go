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
	TileSize      int
	TilesetImages = make(map[string]*ebiten.Image)
	OpaqueImages  = make(map[enums.TileTypeEnum]*ebiten.Image)
	TransImages   = make(map[enums.TileTypeEnum]*ebiten.Image)
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

	TileSize = tilesetDef.TileSize

	for _, t := range tilesetDef.Tiles {
		OpaqueImages[enums.TileTypeEnum(t.Id)] = TilesetImages[t.OpaqueFileName].SubImage(image.Rect(t.X*TileSize, t.Y*TileSize, (t.X+1)*TileSize, (t.Y+1)*TileSize)).(*ebiten.Image)
		TransImages[enums.TileTypeEnum(t.Id)] = TilesetImages[t.TransparentFileName].SubImage(image.Rect(t.X*TileSize, t.Y*TileSize, (t.X+1)*TileSize, (t.Y+1)*TileSize)).(*ebiten.Image)
	}
}
