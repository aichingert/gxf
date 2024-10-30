package drawing

type Gxf struct {
    Lines Mesh
    Polygons Mesh

    Blocks map[string]*Mesh
    Layers map[string][]uint8
}

func NewGxf() *Gxf {
    return &Gxf{
        Lines: Mesh{},
        Polygons: Mesh{}, 

        Blocks: make(map[string]*Mesh),
        Layers: make(map[string][]uint8),
    }
}
