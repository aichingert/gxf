package entity

type Ellipse struct {
	Entity *EntityData

	Ratio float64
	Start float64
	End   float64

	Center   [3]float64
	EndPoint [3]float64
}

func NewEllipse() *Ellipse {
	return &Ellipse{
		Entity: NewEntityData(),

		Center:   [3]float64{0.0, 0.0, 0.0},
		EndPoint: [3]float64{0.0, 0.0, 0.0},
	}
}
