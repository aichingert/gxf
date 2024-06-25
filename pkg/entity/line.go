package entity

type Line struct {
	Entity *EntityData

	Src [3]float64
	Dst [3]float64
}

func NewLine() *Line {
	return &Line{
		Entity: &EntityData{
			Handle:    0,
			Owner:     0,
			LayerName: "",
		},
		Src: [3]float64{0.0, 0.0, 0.0},
		Dst: [3]float64{0.0, 0.0, 0.0},
	}
}
