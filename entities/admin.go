package entities

import "github.com/tomknightdev/dwarven-fortresses/components"

type Admin struct {
	components.RenderSingleton
	components.InputSingleton
	components.GameMapSingleton
	components.NatureSingleton
}
