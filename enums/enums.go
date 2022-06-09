package enums

type TileSpriteEnum int
type TileTypeEnum int
type InputModeEnum int
type GuiActionEnum int
type ResourceTypeEnum int
type DropTypeEnum int
type TaskTypeEnum int
type DesignationTypeEnum int
type ItemTypeEnum int

const (
	TileSpriteEmpty TileSpriteEnum = iota
	TileSpriteCursor
	TileSpriteStockpile
	TileSpriteRock
	TileSpriteRockFloor
	TileSpriteRocks
	TileSpriteDirt0
	TileSpriteDirt1
	TileSpriteGrass0
	TileSpriteGrass1
	TileSpriteGrass2
	TileSpriteTree0
	TileSpriteTree1
	TileSpriteTree2
	TileSpriteTree3
	TileSpriteTree4
	TileSpriteTree5
	TileSpriteLog0
	TileSpriteWater
	TileSpriteDwarf
	TileSpriteStairDown
	TileSpriteStairUp
	TileSpritePickaxe
)

const (
	TileTypeNone TileTypeEnum = iota
	TileTypeRock
	TileTypeTerrain
	TileTypeResource
	TileTypeItem
	TileTypeBuilding
	TileTypeTraverseUp
	TileTypeTraverseDown
)

const (
	InputModeNone InputModeEnum = iota
	InputModeBuild
	InputModeGather
	InputModeChop
	InputModeMine
	InputModeStockpile
	InputModeHaul
)

const (
	GuiActionNone GuiActionEnum = iota
	GuiActionChop
	GuiActionStair
	GuiActionMine
	GuiActionStockpile
)

const (
	ResourceTypeNone ResourceTypeEnum = iota
	ResourceTypeTree
)

const (
	DropTypeNone DropTypeEnum = iota
	DropTypeLog
)

const (
	TaskTypeNone TaskTypeEnum = iota
	TaskTypePickUp
	TaskTypeHaul
	TaskTypeDrop
	TaskTypeChop
	TaskTypeMine
	TaskTypeBuild
	TaskTypeAddToStockpile
)

const (
	DesignationTypeNone DesignationTypeEnum = iota
	DesignationTypeStockpile
)

const (
	ItemTypeNone ItemTypeEnum = iota
	ItemTypeLog
	ItemTypeStone
)
