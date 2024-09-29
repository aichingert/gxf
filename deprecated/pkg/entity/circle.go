package entity

type Circle struct {
	Entity *EntityData

	Radius      float64
	Coordinates [3]float64
}

func NewCircle() *Circle {
	return &Circle{
		Entity:      NewEntityData(),
		Radius:      0.0,
		Coordinates: [3]float64{0.0, 0.0, 0.0},
	}
}

func (e *EntitiesData) AppendCircle(circle *Circle) {
	e.Circles = append(e.Circles, circle)
}
