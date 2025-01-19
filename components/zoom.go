package components

type Zoom struct {
	Value float64
}

func NewZoom() Zoom {
	return Zoom{
		Value: 2,
	}
}
