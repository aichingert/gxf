package parser

import (
    "log"
    "bufio" 

    "github.com/aichingert/dxf/pkg/blocks"
    "github.com/aichingert/dxf/pkg/drawing"
)

func ParseBlocks(sc *bufio.Scanner, dxf *drawing.Dxf) {
    for {
        variable := ExtractCodeAndValue(sc)

        switch variable[1] {
        case "BLOCK":
            parseBlock(sc, dxf)
        case "ENDSEC":
            return
        default:
            if sc.Err != nil {
                log.Fatal("[BLOCK] Scanner Failed: ", sc.Err)
            }
            log.Println("[BLOCK] Warning not implemented: ", variable)
        }

    }
}

func parseBlock(sc *bufio.Scanner, dxf *drawing.Dxf) {
    block := new (blocks.Block) 

    block.Handle = ExtractHex(sc, "5", "handle")
    block.Owner = ExtractHex(sc, "330", "owner")

    for parseSubClass(sc, block) {}

    dxf.Blocks = append(dxf.Blocks, block)
}

func parseSubClass(sc *bufio.Scanner, block *blocks.Block) bool {
    variable := ExtractCodeAndValue(sc)

    switch variable[1] {
    case "AcDbEntity":
        parseAcDbEntity(sc, block)
    case "AcDbBlockBegin":
        parseAcDbBlockBegin(sc, block)
    case "ENDBLK":
        parseEndblk(sc, block)

    // TODO: parse entities
    case "ATTDEF":      fallthrough
    case "LWPOLYLINE":  fallthrough
    case "HATCH":       fallthrough
    case "INSERT":      fallthrough
    case "MTEXT":       fallthrough
    case "LINE":
        // TODO: currently skips to ENDBLK
        parseAttDef(sc, block)

        parseEndblk(sc, block)
    case "AcDbBlockEnd":
        return false
    default:
        log.Fatal("[BLOCK] Failed to parse subClass: ", variable, " ", Line)
    }

    return true
}

func parseAcDbEntity(sc *bufio.Scanner, block *blocks.Block) {
    optional := ExtractCodeAndValue(sc)

    // TODO: think about paper space visibility
    if optional[0] != " 67" {
        block.LayerName = optional[1]
        return
    }

    // TODO: could lead to bug with start and end layername - seems like it is always the same
    layerName := ExtractCodeAndValue(sc)
    block.LayerName = layerName[1]
}

func parseAcDbBlockBegin(sc *bufio.Scanner, block *blocks.Block) {
    block.BlockName = ExtractCodeAndValue(sc)[1]
    block.Flag = ExtractCodeAndValue(sc)[1]
    block.Coordinates = ExtractCoordinates3D(sc)

    // assumption is that this is the blockName again
    validate := ExtractCodeAndValue(sc)

    if block.BlockName != validate[1] {
        log.Fatal("[BLOCK] Invalid assumption different block names")
    }

    block.XrefPath = ExtractCodeAndValue(sc)[1]

    // TODO: use bufio.Reader to be able to peek at possible description
}

func parseEndblk(sc *bufio.Scanner, block *blocks.Block) {
    block.EndHandle = ExtractHex(sc, "5", "end handle")

    if block.Owner != ExtractHex(sc, "330", "end owner") {
        log.Fatal("[BLOCK] Invalid assumption different end owners")
    }
}

// TODO: implement it (currently skips to next)
func parseAttDef(sc *bufio.Scanner, block *blocks.Block) {
    for sc.Scan() && sc.Text() != "ENDBLK" {
        Line++
    }
}
