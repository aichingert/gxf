package entity

type Circle struct {
    *entity

    Radius      float64
    Coordinates [3]float64
}

func NewCircle(handle uint64, owner uint64) *Circle {
    return &Circle {
        entity: &entity {
            handle:     handle,
            owner:      owner,
            LayerName:  "",
        },
        Radius:         0.0,
        Coordinates:    [3]float64{0.0,0.0,0.0},
    }
}
