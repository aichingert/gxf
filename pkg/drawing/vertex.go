package drawing

type Vertex struct {
    X       float32
    Y       float32
    // TODO: 
    // Color uint8
}

type Mesh struct {
    Indices     []uint32
    Vertices    []Vertex
}


