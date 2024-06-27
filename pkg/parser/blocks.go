package parser

import (
    "log"

	"github.com/aichingert/dxf/pkg/blocks"
	_ "github.com/aichingert/dxf/pkg/entity"
	"github.com/aichingert/dxf/pkg/drawing"
)

func ParseBlocks(r *Reader, dxf *drawing.Dxf) error {

    for r.ScanDxfLine() {
        switch r.DxfLine().Line {
        case "BLOCK":
            Wrap(ParseBlock, r, dxf)
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
        case "CIRCLE":
            Wrap(ParseCircle, r, dxf)
        case "HATCH":
            Wrap(ParseHatch, r, dxf)
        case "ENDBLK":
            ParseAcDbEntity(r, block.Entity)
        case "ATTDEF":
            Wrap(ParseAttDef, r, dxf)
        case "AcDbBlockEnd":
            dxf.Blocks = append(dxf.Blocks, block)
            return r.Err()
        default:
            log.Fatal("[Block(", Line, ")] subclass ", r.DxfLine().Line)
        }

        if WrappedErr != nil {
            return WrappedErr
        }
    }

    return r.Err()
}
