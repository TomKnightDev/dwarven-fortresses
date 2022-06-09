package components

type TileInfo struct {
	Position
	TileType
	Sprite
}

func NewTileInfo(pos Position, tileType TileType, sprite Sprite) TileInfo {
	return TileInfo{
		Position: pos,
		TileType: tileType,
		Sprite:   sprite,
	}
}
