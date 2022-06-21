package systems

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sedyh/mizu/pkg/engine"
	"github.com/tomknightdev/dwarven-fortresses/components"
)

type Render struct{}

func NewRender() *Render {
	return &Render{}
}

func (r *Render) Draw(w engine.World, screen *ebiten.Image) {
	renderEnt, found := w.View(components.RenderSingleton{}).Get()
	if !found {
		panic("could not find render singleton")
	}

	var rs *components.RenderSingleton
	renderEnt.Get(&rs)

	camera, found := w.View(components.Zoom{}, components.Position{}).Get()
	if !found {
		return
	}
	var zoom *components.Zoom
	var camPos *components.Position
	camera.Get(&zoom, &camPos)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(0, 0)

	op.GeoM.Scale(zoom.Value, zoom.Value)

	ww, wh := ebiten.WindowSize()
	op.GeoM.Translate(-float64(camPos.X-(ww/2)), -float64(camPos.Y-(wh/2)))
	// op.Filter = ebiten.FilterNearest
	screen.DrawImage(rs.OffScreen, op)
	rs.OffScreen.Clear()
}
