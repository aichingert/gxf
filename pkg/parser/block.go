package parser

import (
    "fmt"

    "github.com/aichingert/gxf/pkg/drawing"
)

/// parseBlock returns nothing. It appends the needed vertices into the gxfs two
/// primary buffers used for storing all vertices and updates the blockOffset
/// array with the values of this block as well as setting the blockOffset for 
/// the block name in another map only used during parsing
func (p *parser) parseBlocks(gxf *drawing.Gxf, layers map[string][]uint8) {
    for {
        switch p.consumeNext() {
        case "BLOCK":
            _ = p.parseEntity()

            name := p.consumeNext()
            p.consumeNext()

            p.discardIf(10)
            p.discardIf(20)
            p.discardIf(30)

            p.discardIf(3)
            p.discardIf(1)

            p.parseEntities(gxf, layers)

            lLen := uint32(len(gxf.Data.Lines.Vertices))
            tLen := uint32(len(gxf.Data.Triangles.Vertices))
            gxf.BlockOffsets = append(gxf.BlockOffsets, [2]uint32{lLen, tLen})
            gxf.BlockNameRes[name] = uint16(len(gxf.BlockOffsets))

            p.consumeUntil("AcDbBlockEnd")
        case "ENDSEC":
            return 
        default:
            p.err = NewParseError(fmt.Sprintf("invalid block value %s", p.line))
        }

        if p.err != nil {
            return
        }

        for p.code != 0 {
            p.consume()
        }
    }
}
