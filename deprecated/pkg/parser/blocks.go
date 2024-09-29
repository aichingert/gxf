package parser

import (
	"fmt"
	"github.com/aichingert/dxf/pkg/blocks"
	"github.com/aichingert/dxf/pkg/drawing"
	_ "github.com/aichingert/dxf/pkg/entity"
)

func ParseBlocks(r Reader, dxf *drawing.Dxf) {
	for {
		switch r.consumeNext() {
		case "BLOCK":
			ParseBlock(r, dxf)
		case "ENDSEC":
			return
		default:
			r.setErr(NewParseError(fmt.Sprintf("Block(%d): %s", Line, r.line())))
			return
		}
	}
}

func ParseBlock(r Reader, dxf *drawing.Dxf) {
	block := blocks.NewBlock()

	ParseAcDbEntity(r, block.Entity)
	ParseAcDbBlockBegin(r, block)

	ParseEntities(r, block.EntitiesData)
	ParseBlockEnd(r, dxf)

	dxf.Blocks[block.BlockName] = block
}

// TODO: maybe pass block to function
func ParseBlockEnd(r Reader, _ *drawing.Dxf) {
	endblk := blocks.NewBlock()

	ParseAcDbEntity(r, endblk.Entity)
	_ = r.assertNextLine("AcDbBlockEnd")
}
