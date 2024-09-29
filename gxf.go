package gxf

import (
    "github.com/aichingert/gxf/pkg/drawing"
    "github.com/aichingert/gxf/pkg/parser"
)

func Parse(buffer []byte) (*drawing.Gxf, error) {
    return parser.ParseBuffer(buffer)
}
