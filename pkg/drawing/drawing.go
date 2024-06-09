package drawing

import (
    "github.com/aichingert/dxf/pkg/header"
    "github.com/aichingert/dxf/pkg/blocks"
    "github.com/aichingert/dxf/pkg/entity"
)

type Dxf struct {
    FileName    string
    Header      *header.Header
    Blocks      []*blocks.Block

    Lines       []*entity.Line
    Polylines   []*entity.Polyline
}

func New(filename string) *Dxf {
    dxf := new (Dxf)

    dxf.FileName = filename
    dxf.Header   = header.New()
    dxf.Blocks   = []*blocks.Block{}
    dxf.Lines    = []*entity.Line{}
    dxf.Polylines    = []*entity.Polyline{}

    return dxf
}

