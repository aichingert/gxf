package drawing

type Vertex struct {
    X       float32
    Y       float32
    Color   uint8
}

type Mesh struct {
    Vertice []Vertex
    Indices []uint32
}


