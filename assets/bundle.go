package assets

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"image"
	"image/png"
	"log"
	"os"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/tomknightdev/dwarven-fortresses/enums"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

// This package is for loading all the images and storing world information

const (
	WorldWidth  = 250
	WorldHeight = 250
	WorldLevels = 10
	Groundlevel = 5

	StartingDwarfCount = 10
)

var (
	TileSize      int
	tilesetImages = make(map[string]*ebiten.Image)
	OpaqueImages  = make(map[enums.TileTypeEnum]*ebiten.Image)
	TransImages   = make(map[enums.TileTypeEnum]*ebiten.Image)
	MainAudio     *mp3.Stream
	MainFont12    font.Face
	MainFont24    font.Face
	MainFont36    font.Face
)

type TilesetDefinition struct {
	TilesetName      string   `json:"tilesetName"`
	TilesetFileNames []string `json:"tilesetFileNames"`
	TileSize         int      `json:"tileSize"`
	Tiles            []Tile   `json:"tiles"`
}

type Tile struct {
	Name                string `json:"name"`
	Id                  int    `json:"id"`
	OpaqueFileName      string `json:"opaqueFileName"`
	TransparentFileName string `json:"transparentFileName"`
	X                   int    `json:"x"`
	Y                   int    `json:"y"`
}

func init() {
	LoadImages()
	LoadAudio()
	LoadFonts()
}

func LoadImages() {
	abs, err := filepath.Abs("./assets/images/tileset.json")
	if err != nil {
		log.Fatal(err)
	}

	// Read tileset json file
	tsd, err := os.ReadFile(abs)
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
		if err != nil {
			log.Fatal(err)
		}

		tileset, err := os.ReadFile(abs)
		if err != nil {
			log.Fatal(err)
		}

		img, err := png.Decode(bytes.NewReader(tileset))
		if err != nil {
			log.Fatal(err)
		}
		tilesetImages[ts] = ebiten.NewImageFromImage(img)
	}

	TileSize = tilesetDef.TileSize

	for _, t := range tilesetDef.Tiles {
		OpaqueImages[enums.TileTypeEnum(t.Id)] = tilesetImages[t.OpaqueFileName].SubImage(image.Rect(t.X*TileSize, t.Y*TileSize, (t.X+1)*TileSize, (t.Y+1)*TileSize)).(*ebiten.Image)
		TransImages[enums.TileTypeEnum(t.Id)] = tilesetImages[t.TransparentFileName].SubImage(image.Rect(t.X*TileSize, t.Y*TileSize, (t.X+1)*TileSize, (t.Y+1)*TileSize)).(*ebiten.Image)
	}

}

func LoadAudio() {
	audio, err := ebitenutil.OpenFile("./assets/audio/scenes/main.mp3")
	if err != nil {
		log.Fatal(err)
	}
	MainAudio, err = mp3.DecodeWithSampleRate(44100, audio)
	if err != nil {
		log.Fatal(err)
	}
}

func LoadFonts() {
	loc, err := filepath.Abs("./assets/fonts/manaspc.ttf")
	if err != nil {
		log.Fatal(err)
	}

	fontFile, err := os.ReadFile(loc)
	if err != nil {
		log.Fatal(err)
	}

	tt, err := opentype.Parse(fontFile)
	if err != nil {
		log.Fatal(err)
	}

	const dpi = 72
	MainFont12, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    12,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}

	MainFont24, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    24,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}

	MainFont36, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    36,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
}
