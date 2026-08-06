package assets

type mapSizePreset struct {
	Name   string
	Width  int
	Height int
}

var mapSizePresets = []mapSizePreset{
	{Name: "Small", Width: 100, Height: 100},
	{Name: "Medium", Width: 250, Height: 250},
	{Name: "Large", Width: 400, Height: 400},
}

// GameSettings holds values the player can configure before starting a game.
type GameSettings struct {
	DwarfCount   int
	MapSizeIndex int
}

func (s *GameSettings) MapWidth() int      { return mapSizePresets[s.MapSizeIndex].Width }
func (s *GameSettings) MapHeight() int     { return mapSizePresets[s.MapSizeIndex].Height }
func (s *GameSettings) MapSizeName() string { return mapSizePresets[s.MapSizeIndex].Name }
func (s *GameSettings) MapSizeCount() int  { return len(mapSizePresets) }

// Settings is the active configuration applied when a new game starts.
var Settings = GameSettings{
	DwarfCount:   StartingDwarfCount,
	MapSizeIndex: 1, // Medium
}
