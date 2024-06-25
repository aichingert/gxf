package entity

type Circle struct {
	Entity *EntityData

	Radius      float64
	Coordinates [3]float64
}

func NewCircle() *Circle {
	return &Circle{
		Entity: &EntityData{
			Handle:    0,
			Owner:     0,
			LayerName: "",
		},
		Radius:      0.0,
		Coordinates: [3]float64{0.0, 0.0, 0.0},
	}
}
