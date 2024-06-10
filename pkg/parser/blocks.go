package parser

import (
    "log"

    "github.com/aichingert/dxf/pkg/blocks"
    "github.com/aichingert/dxf/pkg/drawing"
)

func ParseBlocks(r *Reader, dxf *drawing.Dxf) {
    for {
        switch variable := r.ConsumeDxfLine(); variable.Line {
        case "BLOCK":
            parseBlock(r, dxf)
        case "ENDSEC":
            return
        default:
            log.Println("[BLOCK] Warning not implemented: ", variable)
        }

    }
}

func parseBlock(r *Reader, dxf *drawing.Dxf) {
    block := new (blocks.Block) 

    block.Handle = r.ConsumeHex(5, "handle")
    block.Owner = r.ConsumeHex(330, "owner")

    for parseSubClass(r, block) {}

    dxf.Blocks = append(dxf.Blocks, block)
}

func parseSubClass(r *Reader, block *blocks.Block) bool {
    switch variable := r.ConsumeDxfLine(); variable.Line {
    case "AcDbEntity":
        parseAcDbEntity(r, block)
    case "AcDbBlockBegin":
        parseAcDbBlockBegin(r, block)
    case "ENDBLK":
        parseEndblk(r, block)

    // TODO: parse entities
    case "ATTDEF":      fallthrough
    case "LWPOLYLINE":  fallthrough
    case "HATCH":       fallthrough
    case "INSERT":      fallthrough
    case "MTEXT":       fallthrough
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

func parseAcDbEntity(r *Reader, block *blocks.Block) {
    optional := r.ConsumeDxfLine()

    // TODO: think about paper space visibility
    if optional.Code != 67 {
        block.LayerName = optional.Line
        return
    }

    // TODO: could lead to bug with start and end layername - seems like it is always the same
    layerName := r.ConsumeDxfLine()
    block.LayerName = layerName.Line
}

func parseAcDbBlockBegin(r *Reader, block *blocks.Block) {
    block.BlockName = r.ConsumeDxfLine().Line
    block.Flag = r.ConsumeDxfLine().Line
    block.Coordinates = r.ConsumeCoordinates3D()

    // assumption is that this is the blockName again
    validate := r.ConsumeDxfLine()

    if block.BlockName != validate.Line {
        log.Fatal("[BLOCK] Invalid assumption different block names")
    }

    block.XrefPath = r.ConsumeDxfLine().Line

    // TODO: use bufio.Reader to be able to peek at possible description
}

func parseEndblk(r *Reader, block *blocks.Block) {
    block.EndHandle = r.ConsumeHex(5, "end handle")

    if block.Owner != r.ConsumeHex(330, "end owner") {
        log.Fatal("[BLOCK] Invalid assumption different end owners")
    }
}

// TODO: implement it (currently skips to next)
func parseAttDef(r *Reader, block *blocks.Block) {
    r.SkipToLabel("ENDBLK")
}
