package components

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/tomknightdev/dwarven-fortresses/assets"
)

type NatureSingleton struct {
	GrowTimer        int
	CurrentGrowTimer int

	OffScreen *ebiten.Image
}

func NewNatureSingleton() NatureSingleton {
	return NatureSingleton{
		GrowTimer: 100,
		OffScreen: ebiten.NewImage(assets.WorldWidth*assets.TileSize, assets.WorldHeight*assets.TileSize),
	}
}
