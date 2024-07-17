package entity

type Polyline struct {
	Entity *EntityData

	Flag        int64
	Vertices    int64
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

func (p *Polyline) AppendPLine(coords2D [2]float64, bulge float64) {
	line := PLine{
		X:     coords2D[0],
		Y:     coords2D[1],
		Bulge: bulge,
	}

	p.Coordinates = append(p.Coordinates, line)
}

func (e *EntitiesData) AppendPolyline(polyline *Polyline) {
	e.Polylines = append(e.Polylines, polyline)
}
