package parser

import (
	"log"

	"github.com/aichingert/dxf/pkg/blocks"
	"github.com/aichingert/dxf/pkg/drawing"
	_ "github.com/aichingert/dxf/pkg/entity"
)

func ParseBlocks(r *Reader, dxf *drawing.Dxf) error {
	for r.ScanDxfLine() {
		switch r.DxfLine().Line {
		case "BLOCK":
			Wrap(ParseBlock, r, dxf)
		case "ENDSEC":
			return r.Err()
		default:
			log.Println("Block(", Line, "): ", r.DxfLine().Line)
			return r.Err()
		}

		if WrappedErr != nil {
			return WrappedErr
		}
	}

	return r.Err()
}

func ParseBlock(r *Reader, dxf *drawing.Dxf) error {
	block := blocks.NewBlock()

	if ParseAcDbEntity(r, block.Entity) != nil ||
		ParseAcDbBlockBegin(r, block) != nil {
		return r.Err()
	}

    

    if err := ParseEntities(r, block.EntitiesData); err != nil {
        return err
    }

    Wrap(ParseBlockEnd, r, dxf)

    dxf.Blocks[block.BlockName] = block
    return WrappedEntityErr
}

// TODO: maybe pass block to function
func ParseBlockEnd(r *Reader, dxf *drawing.Dxf) error {
	endblk := blocks.NewBlock()

	if ParseAcDbEntity(r, endblk.Entity) != nil ||
		r.AssertNextLine("AcDbBlockEnd") != nil {
		return r.Err()
	}

	return r.Err()
}
