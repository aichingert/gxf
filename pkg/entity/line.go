package entity

type Line struct {
	Entity *EntityData

	Src [3]float64
	Dst [3]float64
}

func NewLine() *Line {
	return &Line{
		Entity: NewEntityData(),
		Src:    [3]float64{0.0, 0.0, 0.0},
		Dst:    [3]float64{0.0, 0.0, 0.0},
	}
}
