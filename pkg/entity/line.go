package entity

type Line struct {
    *entity

    Src [3]float64
    Dst [3]float64
}

func NewLine(handle uint64, owner uint64) *Line {
    return &Line {
        entity: &entity{
            handle: handle,
            owner: owner,
            LayerName: "",
        },
        Src:    [3]float64{0.0,0.0,0.0},
        Dst:    [3]float64{0.0,0.0,0.0},
    }
}
