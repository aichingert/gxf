package drawing

type Gxf struct {
    Plan    *Block

    Blocks  map[string]*Block
    Layers  map[string][]uint8
}

func NewGxf() *Gxf {
    return &Gxf{
        Plan:   nil,
        Blocks: make(map[string]*Block),
        Layers: make(map[string][]uint8),
    }
}
