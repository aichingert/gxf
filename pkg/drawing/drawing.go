package drawing

import (
	"github.com/aichingert/dxf/pkg/blocks"
	"github.com/aichingert/dxf/pkg/entity"
	"github.com/aichingert/dxf/pkg/header"
)

type Dxf struct {
	FileName string
	Header   *header.Header
	Blocks   []*blocks.Block
	*entity.EntitiesData
}

func New(filename string) *Dxf {
	dxf := new(Dxf)

	dxf.FileName = filename
	dxf.Header = header.New()
	dxf.EntitiesData = entity.New()

	return dxf
}
