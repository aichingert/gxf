package entity

type Polyline struct {
    *entity

    Flag        uint64
    Vertices    uint64
    Coordinates []line
}

type line struct {
    x       float64
    y       float64
    bulge   float64
}

func NewPolyline(handle uint64, owner uint64) *Polyline {
    return &Polyline {
        entity: &entity{
            handle: handle,
            owner: owner,
            LayerName: "",
        },
        Flag: 0,
        Vertices: 0,
        Coordinates: nil,
    }
}

func (p *Polyline) PolylineAppendCoordinate(coords2D [2]float64, bulge float64) {
    line := line{
        x: coords2D[0],
        y: coords2D[1],
        bulge: bulge,
    }

    p.Coordinates = append(p.Coordinates, line)
}
