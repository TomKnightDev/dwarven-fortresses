package components

import (
	"github.com/tomknightdev/dwarven-fortresses/enums"
	"github.com/yohamta/furex/v2"
)

type Gui struct {
	X, Y     int
	Position enums.GuiPositionEnum
	Scale    float64
	Action   enums.GuiActionEnum
	FUI      *furex.View
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
