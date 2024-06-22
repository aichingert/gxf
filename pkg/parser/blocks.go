package parser

import (
    "log"

    "github.com/aichingert/dxf/pkg/blocks"
    "github.com/aichingert/dxf/pkg/drawing"
)

func ParseBlocks(r *Reader, dxf *drawing.Dxf) error {
    for {
        line, err := r.ConsumeDxfLine()
        if err != nil { return err }

        switch line.Line {
        case "BLOCK":
            Wrap(parseBlock, r, dxf)
        case "ENDSEC":
            return nil
        default:
            log.Println("[BLOCK] Warning not implemented: ", line)
        }

        if WrappedErr != nil {
            return WrappedErr
        }
    }
}

func parseBlock(r *Reader, dxf *drawing.Dxf) error {
    block := new (blocks.Block) 

    handle, err := r.ConsumeNumber(5, 16, "handle")
    if err != nil { return err }
    owner,  err := r.ConsumeNumber(330, 16, "owner")
    if err != nil { return err }

    block.Handle = handle
    block.Owner  = owner

    for parseSubClass(r, block) {}

    dxf.Blocks = append(dxf.Blocks, block)
    return nil
}

// TODO: refactor this thing as a whole
func parseSubClass(r *Reader, block *blocks.Block) bool {
    switch variable, _ := r.ConsumeDxfLine(); variable.Line {
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

func parseAcDbBlockBegin(r *Reader, block *blocks.Block) {
    blockName, _ := r.ConsumeDxfLine()
    flag, _ := r.ConsumeDxfLine()
    coordinates, _ := r.ConsumeCoordinates3D()

    block.BlockName = blockName.Line
    block.Flag = flag.Line
    block.Coordinates = coordinates

    // assumption is that this is the blockName again
    validate, _ := r.ConsumeDxfLine()

    if block.BlockName != validate.Line {
        log.Println(validate)
        log.Fatal("[BLOCK(", Line, ")] Invalid assumption different block names")
    }

    xrefPath, _ := r.ConsumeDxfLine()
    block.XrefPath = xrefPath.Line

    // TODO: use bufio.Reader to be able to peek at possible description
}

func parseEndblk(r *Reader, block *blocks.Block) {
    block.EndHandle, _ = r.ConsumeNumber(5, 16, "end handle")

    owner, _ := r.ConsumeNumber(330, 16, "end owner") 
    if block.Owner != owner {
        log.Fatal("[BLOCK] Invalid assumption different end owners")
    }
}

// TODO: implement it (currently skips to next)
func parseAttDef(r *Reader, block *blocks.Block) {
    _ = r.SkipToLabel("ENDBLK")
}
