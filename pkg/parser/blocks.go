package parser

import (
	"log"

	"github.com/aichingert/dxf/pkg/blocks"
	"github.com/aichingert/dxf/pkg/drawing"
)

func ParseBlocks(r *Reader, dxf *drawing.Dxf) error {
	for r.ScanDxfLine() {
		switch r.DxfLine().Line {
		case "BLOCK":
			Wrap(parseBlock, r, dxf)
		case "ENDSEC":
			return nil
		default:
			log.Println("[BLOCK] Warning not implemented: ", r.DxfLine().Line)
		}

		if WrappedErr != nil {
			return WrappedErr
		}
	}

	return r.Err()
}

func parseBlock(r *Reader, dxf *drawing.Dxf) error {
	block := blocks.NewBlock()

	if ParseAcDbEntity(r, block.Entity) != nil {
		return r.Err()
	}

	for parseSubClass(r, block) {
	}

	dxf.Blocks = append(dxf.Blocks, block)
	return nil
}

// TODO: refactor this thing as a whole
func parseSubClass(r *Reader, block *blocks.Block) bool {
	switch variable, _ := r.ConsumeDxfLine(); variable.Line {
	case "AcDbBlockBegin":
		parseAcDbBlockBegin(r, block)
	case "ENDBLK":
		parseEndblk(r, block)

	// TODO: parse entities
	case "ATTDEF":
		fallthrough
	case "LWPOLYLINE":
		fallthrough
	case "HATCH":
		fallthrough
	case "INSERT":
		fallthrough
	case "MTEXT":
		fallthrough
	case "LINE":
		// TODO: currently skips to ENDBLK
		parseAttDef(r, block)

		parseEndblk(r, block)
	case "AcDbBlockEnd":
		return false
	default:
		log.Fatal("[BLOCK] Failed to parse subClass: ", variable, " ", Line)
	}

	return true
}

/*
func parseAcDbEntity(r *Reader, block *blocks.Block) {
	optional, _ := r.ConsumeDxfLine()

	// TODO: think about paper space visibility
	if optional.Code != 67 {
		block.LayerName = optional.Line
		return
	}

	// TODO: could lead to bug with start and end layername - seems like it is always the same
	layerName, _ := r.ConsumeDxfLine()
	block.LayerName = layerName.Line
}
*/

func parseAcDbBlockBegin(r *Reader, block *blocks.Block) error {
	r.ConsumeStr(&block.BlockName)
	r.ConsumeNumber(0, DEC_RADIX, "ato", &block.Flag)
	r.ConsumeCoordinates(block.Coordinates[:])

	// assumption is that this is the blockName again
	if r.AssertNextLine(block.BlockName) != nil {
		log.Fatal("[BLOCK(", Line, ")] Invalid assumption different block names ", block.BlockName)
	}

	r.ConsumeStr(&block.XrefPath)

	// TODO: use bufio.Reader to be able to peek at possible description
	return r.Err()
}

func parseEndblk(r *Reader, block *blocks.Block) error {
	r.ConsumeNumber(5, 16, "end handle", &block.EndHandle)
	r.ConsumeNumber(330, 16, "end owner", &Owner)

	if block.Entity.Owner != Owner {
		log.Fatal("[BLOCK] Invalid assumption different end owners")
	}

	return r.Err()
}

// TODO: implement it (currently skips to next)
func parseAttDef(r *Reader, block *blocks.Block) {
	_ = r.SkipToLabel("ENDBLK")
}
