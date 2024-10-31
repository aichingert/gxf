package parser

import (
    "fmt"

    "github.com/aichingert/gxf/pkg/drawing"
)

// block coordinates
// 
// dx dy dz
//  |  |  |
// ix iy iz
//  |  |  |
// sx sy sz

// 1. Parse block -> then scale it without anything else | this is then considered as the base block
// 2. When parsing entities -> 
//                              calculate the position with offset + scale
//                              then store this data in the entities array, with their position
//                              using this to scale the entities
// 3. Remove all blocks and compare their positions with the base block to create instance data

func (p *parser) parseBlocks(gxf *drawing.Gxf) {
    for {
        switch p.consumeNext() {
        case "BLOCK":
            // parse block begin
            _ = p.parseEntity()

            name := p.consumeNext()
            p.consumeNext() // flags 

            p.discardIf(10) // anchor x
            p.discardIf(20) // anchor y
            p.discardIf(30)

            p.discardIf(3)
            p.discardIf(1)

            // parse entities
            gxf.Blocks[name] = p.parseEntities(gxf.Layers, gxf.Blocks)

            // parse block end
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
