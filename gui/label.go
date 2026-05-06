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
	Text    string
	GetText func() string // if set, called each frame instead of Text
	Font    font.Face
	color.Color

	Button
}

var _ furex.DrawHandler = (*Label)(nil)

func (l *Label) HandleDraw(screen *ebiten.Image, frame image.Rectangle) {
	t := l.Text
	if l.GetText != nil {
		t = l.GetText()
	}
	rect := text.BoundString(l.Font, t)
	text.Draw(screen, t, l.Font, frame.Min.X+((frame.Dx()-rect.Dx())/2), frame.Min.Y+(frame.Dy()+rect.Dy())/2, l.Color)
}
