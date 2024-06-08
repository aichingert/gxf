package drawing

import (
    "github.com/aichingert/dxf/pkg/header"
    "github.com/aichingert/dxf/pkg/blocks"
)

type Dxf struct {
    FileName    string
    Header      *header.Header
    Blocks      []*blocks.Block
}

func New(filename string) *Dxf {
    dxf := new (Dxf)

    dxf.FileName = filename
    dxf.Header   = header.New()
    dxf.Blocks   = []*blocks.Block{}

    return dxf
}
