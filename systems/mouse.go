package systems

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sedyh/mizu/pkg/engine"
	"github.com/tomknightdev/dwarven-fortresses/assets"
	"github.com/tomknightdev/dwarven-fortresses/components"
	"github.com/tomknightdev/dwarven-fortresses/enums"
	"github.com/tomknightdev/dwarven-fortresses/helpers"
)

type Mouse struct {
}

func NewMouse() *Mouse {
	return &Mouse{}
}

func (m *Mouse) Update(w engine.World) {
	var inputSingleton *components.InputSingleton
	is, found := w.View(components.InputSingleton{}).Get()
	if !found {
		panic("input singleton was not found")
	}
	is.Get(&inputSingleton)

	inputSingleton.LeftClickedTiles = nil
	inputSingleton.RightClickedTiles = nil
	if inputSingleton.InGui {
		return
	}

	mouse, found := w.View(components.Mouse{}, components.Position{}, components.Sprite{}).Get()
	if !found {
		return
	}
	var mComp *components.Mouse
	var p *components.Position
	var s *components.Sprite
	mouse.Get(&mComp, &p, &s)

	camera, found := w.View(components.Zoom{}, components.Position{}).Get()
	if !found {
		return
	}
	var zoom *components.Zoom
	var camPos *components.Position
	camera.Get(&zoom, &camPos)

	// Mouse tile position
	ww, wh := ebiten.WindowSize()
	inputSingleton.MouseWorldPosX = inputSingleton.MousePosX + (camPos.X - (ww / 2))
	inputSingleton.MouseWorldPosY = inputSingleton.MousePosY + (camPos.Y - (wh / 2))
	p.X = int((float64(inputSingleton.MouseWorldPosX) / zoom.Value) / float64(assets.TileSize))
	p.Y = int((float64(inputSingleton.MouseWorldPosY) / zoom.Value) / float64(assets.TileSize))
	p.Z = camPos.Z

	if inputSingleton.IsResetInputModePressed {
		inputSingleton.InputMode = enums.InputModeNone
	}

	if inputSingleton.MouseWorldPosX < 0 || inputSingleton.MouseWorldPosY < 0 {
		return
	}

	if inputSingleton.InputMode != enums.InputModeNone {
		if inputSingleton.IsMouseLeftPressed || inputSingleton.IsMouseRightPressed {
			mComp.MouseStart = *p
		}

		if inputSingleton.MouseLeftPressDuration || inputSingleton.MouseRightPressDuration {
			mComp.LeftClickedTiles = nil
			mComp.RightClickedTiles = nil
			startX := mComp.MouseStart.X
			startY := mComp.MouseStart.Y
			endX := p.X
			endY := p.Y

			if endX < startX {
				startX = endX + startX
				endX = startX - endX
				startX = startX - endX
			}

			if endY < startY {
				startY = endY + startY
				endY = startY - endY
				startY = startY - endY
			}

			for mx := startX; mx <= endX; mx++ {
				for my := startY; my <= endY; my++ {
					if inputSingleton.MouseLeftPressDuration {
						mComp.LeftClickedTiles = append(mComp.LeftClickedTiles, components.NewPosition(mx, my, p.Z))
					} else {
						mComp.RightClickedTiles = append(mComp.RightClickedTiles, components.NewPosition(mx, my, p.Z))
					}
				}
			}
		}

		if inputSingleton.IsMouseLeftReleased {
			inputSingleton.LeftClickedTiles = mComp.LeftClickedTiles
			mComp.LeftClickedTiles = nil
		} else if inputSingleton.IsMouseRightReleased {
			inputSingleton.RightClickedTiles = mComp.RightClickedTiles
			mComp.RightClickedTiles = nil
		}
	}
}

func (m *Mouse) Draw(w engine.World, screen *ebiten.Image) {
	var inputSingleton *components.InputSingleton
	is, found := w.View(components.InputSingleton{}).Get()
	if !found {
		panic("input singleton was not found")
	}
	is.Get(&inputSingleton)

	if inputSingleton.InGui {
		return
	}

	if inputSingleton.InputMode == enums.InputModeNone {
		return
	}

	mouseInput, found := w.View(components.Mouse{}, components.Position{}, components.Sprite{}).Get()
	if !found {
		return
	}
	var mComp *components.Mouse
	var mousePos *components.Position
	var mouseSprite *components.Sprite
	mouseInput.Get(&mComp, &mousePos, &mouseSprite)

	if inputSingleton.MouseWorldPosX < 0 || inputSingleton.MouseWorldPosY < 0 {
		return
	}

	switch inputSingleton.InputMode {
	case enums.InputModeNone:
		mouseSprite.Image = assets.OpaqueImages[enums.TileSpriteEmpty]
	case enums.InputModeGather:
		mouseSprite.Image = assets.OpaqueImages[enums.TileSpriteCursor]
	case enums.InputModeBuild:
		mouseSprite.Image = assets.OpaqueImages[enums.TileSpriteStairDown]
	case enums.InputModeChop:
		mouseSprite.Image = assets.OpaqueImages[enums.TileSpriteCursor]
	case enums.InputModeMine:
		mouseSprite.Image = assets.OpaqueImages[enums.TileSpritePickaxe]
	case enums.InputModeStockpile:
		mouseSprite.Image = assets.OpaqueImages[enums.TileSpriteStockpile]
	}

	helpers.DrawImage(w, screen, *mousePos, mouseSprite.Image)

	if inputSingleton.MouseLeftPressDuration || inputSingleton.MouseRightPressDuration {
		startX := mComp.MouseStart.X
		startY := mComp.MouseStart.Y
		endX := mousePos.X
		endY := mousePos.Y

		if endX < startX {
			startX = endX + startX
			endX = startX - endX
			startX = startX - endX
		}

		if endY < startY {
			startY = endY + startY
			endY = startY - endY
			startY = startY - endY
		}

		for mx := startX; mx <= endX; mx++ {
			for my := startY; my <= endY; my++ {
				helpers.DrawImage(w, screen, components.NewPosition(mx, my, mousePos.Z), mouseSprite.Image)
			}
		}
	}
}
