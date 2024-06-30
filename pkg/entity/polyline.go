package entity

type Polyline struct {
	Entity *EntityData

	Flag        uint64
	Vertices    uint64
	Coordinates []PLine
}

type PLine struct {
	X     float64
	Y     float64
	Bulge float64
}

func NewPolyline() *Polyline {
	return &Polyline{
		Entity:      NewEntityData(),
		Flag:        0,
		Vertices:    0,
		Coordinates: nil,
	}
}

func (p *Polyline) PolylineAppendCoordinate(coords2D [2]float64, bulge float64) {
	line := PLine{
		X:     coords2D[0],
		Y:     coords2D[1],
		Bulge: bulge,
	}

	p.Coordinates = append(p.Coordinates, line)
}
