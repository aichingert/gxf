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
            ParseAcDbEntity(r, block.Entity)
        case "ATTDEF":
            Wrap(ParseAttDef, r, dxf)
        case "REGION":
            Wrap(ParseRegion, r, dxf)
        case "AcDbBlockEnd":
            dxf.Blocks = append(dxf.Blocks, block)
        default:
            log.Fatal("[Block(", Line, ")] subclass ", r.DxfLine().Line)
        }

        if WrappedErr != nil || r.DxfLine().Line == "AcDbBlockEnd" {
            return WrappedErr
        }

        r.ConsumeStrIf(1001, nil)
        r.ConsumeNumberIf(1070, DEC_RADIX, "not sure", nil)
        r.ConsumeNumberIf(1071, DEC_RADIX, "not sure", nil)

        r.ConsumeStrIf(1000, nil)
        r.ConsumeStrIf(1000, nil)

        r.ConsumeNumberIf(1005, HEX_RADIX, "not sure", nil)
        r.ConsumeStrIf(1001, nil)
        r.ConsumeNumberIf(1070, DEC_RADIX, "not sure", nil)
        r.ConsumeStrIf(1000, nil)
        r.ConsumeStrIf(1002, nil)
        r.ConsumeNumberIf(1070, DEC_RADIX, "not sure", nil)
        r.ConsumeNumberIf(1070, DEC_RADIX, "not sure", nil)
        r.ConsumeStrIf(1002, nil) 

        r.ConsumeStrIf(1001, nil)
        r.ConsumeNumberIf(1010, DEC_RADIX, "not sure", nil)
        r.ConsumeNumberIf(1020, DEC_RADIX, "not sure", nil)
        r.ConsumeNumberIf(1030, DEC_RADIX, "not sure", nil)

        r.ConsumeStrIf(1001, nil)
        r.ConsumeNumberIf(1070, DEC_RADIX, "not sure", nil)
        r.ConsumeNumberIf(1071, DEC_RADIX, "not sure", nil)
        r.ConsumeNumberIf(1005, HEX_RADIX, "not sure", nil)

        r.ConsumeStrIf(1001, nil)
        r.ConsumeNumberIf(1010, DEC_RADIX, "not sure", nil)
        r.ConsumeNumberIf(1020, DEC_RADIX, "not sure", nil)
        r.ConsumeNumberIf(1030, DEC_RADIX, "not sure", nil)
    }

    return r.Err()
}
