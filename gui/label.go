package gui

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/yohamta/furex/v2"
	"golang.org/x/image/font"
)

type Label struct {
	Text string
	Font font.Face
}

var _ furex.DrawHandler = (*Label)(nil)

func (l *Label) HandleDraw(screen *ebiten.Image, frame image.Rectangle) {
	rect := text.BoundString(l.Font, l.Text)
	text.Draw(screen, l.Text, l.Font, frame.Min.X+((frame.Dx()-rect.Dx())/2), frame.Min.Y+(frame.Dy()+rect.Dy())/2, color.White)
}
