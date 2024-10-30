package drawing

import "math"

type Block struct {
    Lines       *Mesh
    Triangles   *Mesh

    Bounds      *Bounds
}

type Mesh struct {
    Indices     []uint32
    Vertices    []Vertex
}

type Vertex struct {
    X       float32
    Y       float32
    R       float32
    G       float32
    B       float32
}

type Bounds struct {
    minX float32
    maxX float32

    minY float32
    maxY float32
}

func NewBlock() *Block {
    return &Block {
        Lines:      NewMesh(),
        Triangles:  NewMesh(),
        Bounds:     NewBounds(),
    }
}

func NewMesh() *Mesh {
    return &Mesh {
        Indices:    []uint32{},
        Vertices:   []Vertex{},
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
        maxX: -1_000_000.0,

        minY: math.MaxFloat32,
        maxY: -1_000_000.0,
    }
}

func (b *Bounds) UpdateWithScale(other *Bounds, sx float32, sy float32) {
    b.UpdateX([]float32{other.minX * sx, other.maxX * sx})
    b.UpdateY([]float32{other.minY * sy, other.maxY * sy})
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

func NewVertex(x float32, y float32, rgb []uint8) Vertex {
    return Vertex{
        X: x,
        Y: y,
        R: float32(rgb[0]) / 255.,
        G: float32(rgb[1]) / 255.,
        B: float32(rgb[2]) / 255.,
    }
}
