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
	Ellipses []*entity.Ellipse

	Lines     []*entity.Line
	Polylines []*entity.Polyline

	Texts  []*entity.Text
	MTexts []*entity.MText

    Hatches []*entity.Hatch
}

func New(filename string) *Dxf {
	dxf := new(Dxf)

	dxf.FileName = filename
	dxf.Header = header.New()

	return dxf
}
