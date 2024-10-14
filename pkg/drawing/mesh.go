package drawing

import "math"

type Vertex struct {
    X       float32
    Y       float32
    R       float32
    G       float32
    B       float32
}

type Mesh struct {
    Indices     []uint32
    Vertices    []Vertex
}

type Bounds struct {
    minX float32
    maxX float32

    minY float32
    maxY float32
}

func NewMesh() *Mesh {
    return &Mesh {
        Indices: []uint32{},
        Vertices: []Vertex{},
    }
}

func (m *Mesh) Scale(bounds *Bounds) {
    xDenom := (bounds.maxX - bounds.minX) / 2
    yDenom := (bounds.maxY - bounds.minY) / 2

    for i := range m.Vertices {
        m.Vertices[i].X = (m.Vertices[i].X - bounds.minX) / xDenom - 1.0
        m.Vertices[i].Y = (m.Vertices[i].Y - bounds.minY) / yDenom - 1.0
    }
}

func NewBounds() *Bounds {
    return &Bounds{
        minX: math.MaxFloat32,
        maxX: -50_000.0,

        minY: math.MaxFloat32,
        maxY: -50_000.0,
    }
}

func (b *Bounds) UpdateX(xs []float32) {
    for _, x := range xs {
        if b.minX > x { b.minX = x }
        if b.maxX < x { b.maxX = x }
    }
}

func (b *Bounds) UpdateY(ys []float32) {
    for _, y := range ys {
        if b.minY > y { b.minY = y }
        if b.maxY < y { b.maxY = y }
    }
}

func NewVertex(x float32, y float32) Vertex {
    return Vertex{
        X: x,
        Y: y,
        R: 0.65,
        G: 0.65,
        B: 0.65,
    }
}
