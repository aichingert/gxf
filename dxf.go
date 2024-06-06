package dxf

import (
    "github.com/aichingert/dxf/pkg/parser"
    "github.com/aichingert/dxf/pkg/drawing"
)

func Open(filename string) *drawing.Dxf {
    return parser.FromFile(filename)
}
