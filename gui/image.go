package gui

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/tomknightdev/dwarven-fortresses/assets"
	"github.com/yohamta/furex/v2"
)

type Image struct {
	Image *ebiten.Image
	Scale float64
}

var _ furex.DrawHandler = (*Image)(nil)

func (i *Image) HandleDraw(screen *ebiten.Image, frame image.Rectangle) {
	op := &ebiten.DrawImageOptions{}

	//offsetX, offsetY := i.Image.Size()
	op.GeoM.Translate(float64(frame.Min.X+frame.Dx()/2-assets.TileSize*int(i.Scale)/2)/i.Scale, float64(frame.Min.Y+frame.Dy()/2-assets.TileSize*int(i.Scale)/2)/i.Scale)

	op.GeoM.Scale(i.Scale, i.Scale)

	screen.DrawImage(i.Image, op)
}
