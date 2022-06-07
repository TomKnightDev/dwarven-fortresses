package enums

type TileTypeEnum int
type InputModeEnum int
type GuiActionEnum int
type ResourceTypeEnum int
type DropTypeEnum int
type TaskTypeEnum int
type DesignationTypeEnum int
type ItemTypeEnum int

const (
	TileTypeEmpty TileTypeEnum = iota
	TileTypeCursor
	TileTypeStockpile
	TileTypeRockWallHz
	TileTypeRockFloor
	TileTypeRocks
	TileTypeDirt0
	TileTypeDirt1
	TileTypeGrass0
	TileTypeGrass1
	TileTypeGrass2
	TileTypeTree0
	TileTypeTree1
	TileTypeTree2
	TileTypeTree3
	TileTypeTree4
	TileTypeTree5
	TileTypeLog0
	TileTypeWater
	TileTypeDwarf
	TileTypeStairDown
	TileTypeStairUp
	TileTypePickaxe
	TileTypeRockWallVt
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
