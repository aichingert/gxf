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

	for r.ScanDxfLine() {
		switch r.DxfLine().Line {
		case "INSERT":
			Wrap(ParseInsert, r, dxf)
		case "LINE":
			Wrap(ParseLine, r, dxf)
		case "LWPOLYLINE":
			Wrap(ParsePolyline, r, dxf)
		case "MTEXT":
			Wrap(ParseMText, r, dxf)
		case "ARC":
			Wrap(ParseArc, r, dxf)
		case "CIRCLE":
			Wrap(ParseCircle, r, dxf)
		case "HATCH":
			Wrap(ParseHatch, r, dxf)
		case "ENDBLK":
			Wrap(ParseBlockEnd, r, dxf)
			dxf.Blocks = append(dxf.Blocks, block)
			return WrappedErr
		case "ATTDEF":
			Wrap(ParseAttdef, r, dxf)
		case "REGION":
			Wrap(ParseRegion, r, dxf)
		default:
			log.Fatal("[Block(", Line, ")] invalid subclass ", r.DxfLine().Line)
		}

		if WrappedErr != nil {
			return WrappedErr
		}

		// TODO: parse XDATA
		code, err := r.PeekCode()
		log.Println(code)
		for code != 0 && err == nil {
			r.ConsumeStr(nil)
			code, err = r.PeekCode()
		}
		if err != nil {
			return err
		}
	}

	return r.Err()
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
