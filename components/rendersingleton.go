package components

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/tomknightdev/dwarven-fortresses/assets"
)

type RenderSingleton struct {
	OffScreen *ebiten.Image
}

func NewRenderSingleton() RenderSingleton {
	return RenderSingleton{
		OffScreen: ebiten.NewImage(assets.Settings.MapWidth()*assets.TileSize, assets.Settings.MapHeight()*assets.TileSize),
	}
}
