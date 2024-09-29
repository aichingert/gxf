package dxf

import (
	"github.com/aichingert/dxf/pkg/drawing"
	"github.com/aichingert/dxf/pkg/parser"
)

func Open(filename string) (*drawing.Dxf, error) {
	return parser.FromFile(filename)
}

func Parse(filename string, buffer []byte) (*drawing.Dxf, error) {
	return parser.FromBuffer(filename, buffer)
}
