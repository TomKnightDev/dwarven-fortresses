package systems

import (
	"github.com/sedyh/mizu/pkg/engine"
	"github.com/tomknightdev/dwarven-fortresses/assets"
	"github.com/tomknightdev/dwarven-fortresses/components"
)

type Camera struct {
}

func NewCamera() *Camera {
	return &Camera{}
}

func (c *Camera) Update(w engine.World) {
	var inputSingleton *components.InputSingleton
	is, found := w.View(components.InputSingleton{}).Get()
	if !found {
		panic("input singleton was not found")
	}
	is.Get(&inputSingleton)

	camera, found := w.View(components.Zoom{}, components.Position{}).Get()
	if !found {
		return
	}
	var zoom *components.Zoom
	var camPos *components.Position
	camera.Get(&zoom, &camPos)

	if inputSingleton.IsCameraRightPressed && camPos.X < assets.WorldWidth*assets.TileSize {
		camPos.X += assets.TileSize / 2
	} else if inputSingleton.IsCameraLeftPressed && camPos.X > 0 {
		camPos.X -= assets.TileSize / 2
	}
	if inputSingleton.IsCameraDownPressed && camPos.Y < assets.WorldHeight*assets.TileSize {
		camPos.Y += assets.TileSize / 2
	} else if inputSingleton.IsCameraUpPressed && camPos.Y > 0 {
		camPos.Y -= assets.TileSize / 2
	}

	// Zoom the camera
	if inputSingleton.MouseWheel > 0 && zoom.Value < 10 {
		zoom.Value += 0.1
	} else if inputSingleton.MouseWheel < 0 && zoom.Value > 0.2 {
		zoom.Value -= 0.1
	}

	// Adjust camera z level
	if inputSingleton.IsCameraLowerPressed && camPos.Z > 1 {
		camPos.Z--
	} else if inputSingleton.IsCameraRaisePressed && camPos.Z < assets.WorldLevels {
		camPos.Z++
	}

}
