package drawing

import (
    "github.com/aichingert/dxf/pkg/header"
)

type Dxf struct {
    FileName    string
    Header      *header.Header
}

func New(filename string) *Dxf {
    dxf := new (Dxf)

    dxf.FileName = filename
    dxf.Header   = header.New()

    return dxf
}
