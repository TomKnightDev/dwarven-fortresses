package gui

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/yohamta/furex/v2"
)

type Label struct {
	Text string
	Font text.Face
	color.Color

	Button
}

var _ furex.DrawHandler = (*Label)(nil)

func (l *Label) HandleDraw(screen *ebiten.Image, frame image.Rectangle) {
	width, height := text.Measure(l.Text, l.Font, 0)

	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(frame.Min.X)+(float64(frame.Dx())-width)/2, float64(frame.Min.Y)+(float64(frame.Dy())-height)/2)
	op.ColorScale.ScaleWithColor(l.Color)

	text.Draw(screen, l.Text, l.Font, op)
}
