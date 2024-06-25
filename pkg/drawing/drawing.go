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

	Arcs    []*entity.Arc
	Circles []*entity.Circle

	Lines     []*entity.Line
	Polylines []*entity.Polyline

	MTexts []*entity.MText
}

func New(filename string) *Dxf {
	dxf := new(Dxf)

	dxf.FileName = filename
	dxf.Header = header.New()
	dxf.Blocks = []*blocks.Block{}

	dxf.Arcs = []*entity.Arc{}
	dxf.Circles = []*entity.Circle{}

	dxf.Lines = []*entity.Line{}
	dxf.Polylines = []*entity.Polyline{}

	dxf.MTexts = []*entity.MText{}

	return dxf
}
