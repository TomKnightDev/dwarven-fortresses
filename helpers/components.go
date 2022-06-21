package helpers

import (
	"github.com/sedyh/mizu/pkg/engine"
	"github.com/tomknightdev/dwarven-fortresses/components"
)

func GetGameMapSingleton(w engine.World) *components.GameMapSingleton {
	gmsent, found := w.View(components.GameMapSingleton{}).Get()
	if !found {
		panic("game map singleton not found")
	}
	var gms *components.GameMapSingleton
	gmsent.Get(&gms)
	return gms
}

func GetRenderSingleton(w engine.World) *components.RenderSingleton {
	rsent, found := w.View(components.RenderSingleton{}).Get()
	if !found {
		panic("render singleton not found")
	}
	var rs *components.RenderSingleton
	rsent.Get(&rs)
	return rs
}

func GetInputSingleton(w engine.World) *components.InputSingleton {
	isent, found := w.View(components.InputSingleton{}).Get()
	if !found {
		panic("input singletone not found")
	}
	var is *components.InputSingleton
	isent.Get(&is)
	return is
}
