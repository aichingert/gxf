package dxf

import (
	"github.com/aichingert/dxf/pkg/drawing"
	"github.com/aichingert/dxf/pkg/parser"
)

func Open(filename string) (*drawing.Dxf, error) {
	return parser.FromFile(filename)
}
