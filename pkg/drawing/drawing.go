package drawing

/// Gxf represents the whole drawing with only one line and triangle buffer
/// these two buffers store the entire dxf file, blocks are included within
/// them. The data is being used by the instances and the last bit by the 
/// default entities instances use the block offsets and a transformation
/// matrix to reduce the memory footprint.
type Gxf struct {
    Data            *Obj

    BlockOffsets    [][2]uint32
    BlockNameRes    map[string]uint16
    InstanceData    map[string][][4]float32
}

func NewGxf() *Gxf {
    return &Gxf {
        Data:           NewObj(),
        BlockOffsets:   [][2]uint32{},
        BlockNameRes:   make(map[string]uint16),
        InstanceData:   make(map[string][][4]float32),
    }
}
