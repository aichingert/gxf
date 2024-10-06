package drawing

import "math"

type Gxf struct {
    BorderX [2]float32
    BorderY [2]float32

    Lines Mesh
    Polygons Mesh
}

func NewGxf() *Gxf {
    return &Gxf{
        BorderX: [2]float32{ math.MaxFloat32, -50_000.0 },
        BorderY: [2]float32{ math.MaxFloat32, -50_000.0 },

        Lines: Mesh{},
        Polygons: Mesh{},
    }
}
