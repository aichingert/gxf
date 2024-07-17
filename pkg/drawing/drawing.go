package drawing

import (
	"github.com/aichingert/dxf/pkg/blocks"
	"github.com/aichingert/dxf/pkg/entity"
	"github.com/aichingert/dxf/pkg/header"
	"github.com/aichingert/dxf/pkg/table"
)

type Dxf struct {
	FileName string
	Header   *header.Header
	Blocks   map[string]*blocks.Block
	Layers   map[string]*table.Layer
	*entity.EntitiesData
}

func New(filename string) *Dxf {
	dxf := new(Dxf)

	dxf.FileName = filename
	dxf.Header = header.New()
	dxf.Blocks = make(map[string]*blocks.Block)
	dxf.Layers = make(map[string]*table.Layer)
	dxf.EntitiesData = entity.New()

	return dxf
}
