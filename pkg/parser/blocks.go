package parser

import (

	_ "github.com/aichingert/dxf/pkg/blocks"
	"github.com/aichingert/dxf/pkg/drawing"
)

func ParseBlocks(r *Reader, dxf *drawing.Dxf) error {
    return r.Err()
}
