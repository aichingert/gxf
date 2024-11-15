package drawing

type Obj struct {
    Lines       *Mesh
    Triangles   *Mesh
}

func NewObj() *Obj {
    return &Obj {
        Lines:      NewMesh(),
        Triangles:  NewMesh(),
    }
}

/// Mesh is used to store vertex and index data of the dxf file
/// the vertex data is in the representation:
/// 
/// type Vertex struct {
///     X       float32
///     Y       float32
///     R       float32
///     G       float32
///     B       float32
/// }
/// 
/// indices are used to make the vertex buffer smaller
type Mesh struct {
    Indices     []uint16
    Vertices    []float32
}

func NewMesh() *Mesh {
    return &Mesh {
        Indices:    []uint16{},
        Vertices:   []float32{},
    }
}

/// AddVertex returns nothing and takes a topology to discern between lines
/// and triangles it then pushes the vertex into the meshes corresponding
/// vertex buffer. Here are the topology to buffer mappings:
///
///     topology: 1 => Lines
///     topology: 2 => Triangles
func (o *Obj) AddVertex(topology uint8, x float32, y float32, rgb []uint8) {
    var buffer *Mesh

    switch (topology) {
    case 1: 
        buffer = o.Lines
    case 2:
        buffer = o.Triangles
    default:
        return
    }

    buffer.Vertices = append(buffer.Vertices, x)
    buffer.Vertices = append(buffer.Vertices, y)
    buffer.Vertices = append(buffer.Vertices, float32(rgb[0]) / 255.)
    buffer.Vertices = append(buffer.Vertices, float32(rgb[1]) / 255.)
    buffer.Vertices = append(buffer.Vertices, float32(rgb[2]) / 255.)
}

//func (m *Mesh) Scale(bounds *Bounds) {
//    xDenom := (bounds.maxX - bounds.minX) / 2
//    yDenom := (bounds.maxY - bounds.minY) / 2
//
//    for i := range m.Vertices {
//        m.Vertices[i].X = (m.Vertices[i].X - bounds.minX) / xDenom - 1.0
//        m.Vertices[i].Y = (m.Vertices[i].Y - bounds.minY) / yDenom - 1.0
//    }
//}
