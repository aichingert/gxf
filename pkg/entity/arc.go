package entity 

type Arc struct {
    *entity
    Circle      *Circle

    StartAngle  float64
    EndAngle    float64
}

func NewArc(handle uint64, owner uint64) *Arc {
    return &Arc {
        entity: &entity{
            handle:     handle,
            owner:      owner,
            LayerName:  "",
        },
        Circle:         &Circle{
            Coordinates:    [3]float64{0.0, 0.0, 0.0},
            Radius:         0.0,
        },
        StartAngle:     0.0,
        EndAngle:       0.0,
    }
}
