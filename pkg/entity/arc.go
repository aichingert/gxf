package entity

type Arc struct {
	Entity *EntityData
	Circle *Circle

	StartAngle float64
	EndAngle   float64
}

func NewArc() *Arc {
	return &Arc{
		Entity: &EntityData{
			Handle:    0,
			Owner:     0,
			LayerName: "",
		},
		Circle:     NewCircle(),
		StartAngle: 0.0,
		EndAngle:   0.0,
	}
}
