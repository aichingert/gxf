package entity

type Arc struct {
	Entity *EntityData
	Circle *Circle

	StartAngle float64
	EndAngle   float64
    Counterclockwise    uint64
}

func NewArc() *Arc {
	return &Arc{
		Entity:     NewEntityData(),
		Circle:     NewCircle(),
		StartAngle: 0.0,
		EndAngle:   0.0,
        Counterclockwise: 0,
	}
}
