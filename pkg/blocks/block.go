package blocks

type Block struct {
    BlockName   string
    LayerName   string

    Owner       uint64
    Handle      uint64
    EndHandle   uint64

    // TODO: make this an int
    Flag        string

    XrefPath    string
    Description string
    Coordinates [3]float64
}

func New() *Block {
    return new (Block)
}
