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
			WrapEntity(ParseInsert, r, block.EntitiesData)
		case "LINE":
			WrapEntity(ParseLine, r, block.EntitiesData)
		case "LWPOLYLINE":
			WrapEntity(ParseLwPolyline, r, block.EntitiesData)
		case "POINT":
			WrapEntity(ParsePoint, r, block.EntitiesData)
		case "TEXT":
			WrapEntity(ParseText, r, block.EntitiesData)
		case "MTEXT":
			WrapEntity(ParseMText, r, block.EntitiesData)
		case "ARC":
			WrapEntity(ParseArc, r, block.EntitiesData)
		case "CIRCLE":
			WrapEntity(ParseCircle, r, block.EntitiesData)
		case "HATCH":
			WrapEntity(ParseHatch, r, block.EntitiesData)
		case "DIMENSION":
			WrapEntity(ParseDimension, r, block.EntitiesData)
		case "REGION":
			WrapEntity(ParseRegion, r, block.EntitiesData)
		case "ATTDEF":
			WrapEntity(ParseAttdef, r, block.EntitiesData)
		case "ENDBLK":
			Wrap(ParseBlockEnd, r, dxf)
			dxf.Blocks[block.BlockName] = block
			return WrappedEntityErr
		default:
			log.Fatal("[Block(", Line, ")] invalid subclass ", r.DxfLine().Line)
		}

		if WrappedEntityErr != nil {
			return WrappedEntityErr
		}

		// TODO: parse XDATA
		code, err := r.PeekCode()
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
