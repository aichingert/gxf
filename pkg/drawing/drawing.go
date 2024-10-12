package drawing

type Gxf struct {
    Lines Mesh
    Polygons Mesh
}

func NewGxf() *Gxf {
    return &Gxf{
        Lines: Mesh{},
        Polygons: Mesh{}, 
    }
}
