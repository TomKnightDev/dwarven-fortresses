package components

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/tomknightdev/dwarven-fortresses/enums"
)

type Gui struct {
	X, Y     int
	Position enums.GuiPositionEnum
	Scale    float64
	Action   enums.GuiActionEnum
	UIUpdate func(*ebiten.Image)
}

func NewGui(x, y int, position enums.GuiPositionEnum, scale float64, action enums.GuiActionEnum) Gui {
	return Gui{
		X:        x,
		Y:        y,
		Position: position,
		Scale:    scale,
		Action:   action,
	}
}
