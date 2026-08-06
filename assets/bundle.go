package assets

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"image"
	"image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
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

//go:embed images/tileset.json
var tilesetJSON []byte

//go:embed images/tiny_dungeon_interface.png
var imgInterface []byte

//go:embed images/tiny_dungeon_world.png
var imgWorld []byte

//go:embed images/tiny_dungeon_monsters.png
var imgMonsters []byte

//go:embed images/tiny_dungeon_items.png
var imgItems []byte

//go:embed audio/scenes/main.mp3
var audioMain []byte

//go:embed fonts/manaspc.ttf
var fontMain []byte

var embeddedImages = map[string][]byte{}

var (
	TileSize      int
	OpaqueImages  = make(map[enums.TileTypeEnum]*ebiten.Image)
	TransImages   = make(map[enums.TileTypeEnum]*ebiten.Image)
	MainAudio     *mp3.Stream
	MainFont12    text.Face
	MainFont24    text.Face
	MainFont36    text.Face
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
	embeddedImages = map[string][]byte{
		"tiny_dungeon_interface.png": imgInterface,
		"tiny_dungeon_world.png":     imgWorld,
		"tiny_dungeon_monsters.png":  imgMonsters,
		"tiny_dungeon_items.png":     imgItems,
	}
	LoadImages()
	LoadAudio()
	LoadFonts()
}

func LoadImages() {
	tilesetDef := TilesetDefinition{}
	if err := json.Unmarshal(tilesetJSON, &tilesetDef); err != nil {
		log.Fatal(err)
	}

	tilesetImages := make(map[string]*ebiten.Image)
	for _, ts := range tilesetDef.TilesetFileNames {
		data, ok := embeddedImages[ts]
		if !ok {
			log.Fatalf("embedded image not found: %s", ts)
		}
		img, err := png.Decode(bytes.NewReader(data))
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
	stream, err := mp3.DecodeWithSampleRate(44100, bytes.NewReader(audioMain))
	if err != nil {
		log.Fatal(err)
	}
	MainAudio = stream
}

func LoadFonts() {
	tt, err := opentype.Parse(fontMain)
	if err != nil {
		log.Fatal(err)
	}

	const dpi = 72
	face12, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    12,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	MainFont12 = text.NewGoXFace(face12)

	face24, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    24,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	MainFont24 = text.NewGoXFace(face24)

	face36, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    36,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	MainFont36 = text.NewGoXFace(face36)
}
