package drawing

import "math"

type Gxf struct {
    Lines Mesh
    Polygons Mesh

    MinX float32
    MinY float32

    MaxX float32
    MaxY float32
}

func NewGxf() *Gxf {
    return &Gxf{
        Lines: Mesh{},
        Polygons: Mesh{},

        MinX: math.MaxFloat32,
        MaxX: -50_000.0,

        MinY: math.MaxFloat32,
        MaxY: -50_000.0,
    }
}

func max(a float32, b float32) float32 {
    if a > b {
        return a
    } else {
        return b
    }
}

func min(a float32, b float32) float32 {
    if a < b {
        return a
    } else {
        return b
    }
}

func (g *Gxf) UpdateBorder(srcX float32, dstX float32, srcY float32, dstY float32) {
    g.MinX = min(g.MinX, srcX)
    g.MinX = min(g.MinX, dstX)
    g.MaxX = max(g.MaxX, srcX)
    g.MaxX = max(g.MaxX, dstX)

    g.MinY = min(g.MinY, srcY)
    g.MinY = min(g.MinY, dstY)
    g.MaxY = max(g.MaxY, srcY)
    g.MaxY = max(g.MaxY, dstY)
}
