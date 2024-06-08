package parser

import (
    "log"
    "bufio"
    "strconv"

    "github.com/aichingert/dxf/pkg/blocks"
    "github.com/aichingert/dxf/pkg/drawing"
)

func ParseBlocks(sc *bufio.Scanner, dxf *drawing.Dxf) {
    for true {
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
    block := blocks.New()

    block.Handle = extractHex(sc, "5", "handle")
    block.Owner = extractHex(sc, "330", "owner")

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
    case "AcDbBlockEnd":
        return false
    default:
        log.Fatal("[BLOCK] Failed to parse subClass: ", variable)
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
    block.Coordinates = ExtractCoordinates(sc)

    // assumption is that this is the blockName again
    validate := ExtractCodeAndValue(sc)

    if block.BlockName != validate[1] {
        log.Fatal("[BLOCK] Invalid assumption different block names")
    }

    block.XrefPath = ExtractCodeAndValue(sc)[1]

    // TODO: use bufio.Reader to be able to peek at possible description
}

func parseEndblk(sc *bufio.Scanner, block *blocks.Block) {
    block.EndHandle = extractHex(sc, "5", "end handle")

    if block.Owner != extractHex(sc, "330", "end owner") {
        log.Fatal("[BLOCK] Invalid assumption different end owners")
    }
}

func extractHex(sc *bufio.Scanner, code string, description string) uint64 {
    value := ExtractCodeAndValue(sc)

    if code != value[0] {
        log.Fatal("[BLOCK] parseBlock failed invalid group code: expected ", code, " got ", value)
    }

    val, err := strconv.ParseUint(value[1], 16, 64)

    if err != nil {
        log.Fatal("[BLOCK] parseBlock failed invalid ", description, ": (", value, ")")
    }

    return val
}
