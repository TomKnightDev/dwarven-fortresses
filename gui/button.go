package gui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/furex/v2"
)

type Button struct {
	pressed bool
	OnClick func()
}

var _ furex.ButtonHandler = (*Button)(nil)

func (b *Button) HandlePress(x, y int, t ebiten.TouchID) {
	b.pressed = true
}

func (b *Button) HandleRelease(x, y int, isCancel bool) {
	b.pressed = false
	if !isCancel {
		if b.OnClick != nil {
			b.OnClick()
		}
	}
}
