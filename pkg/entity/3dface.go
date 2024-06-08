package entity

type ThreeDFace struct {
    *entity
}

func (tdf *ThreeDFace) Handle() uint64 {
    return tdf.entity.Handle
}

func (tdf *ThreeDFace) Owner() uint64 {
    return tdf.entity.Owner
}
