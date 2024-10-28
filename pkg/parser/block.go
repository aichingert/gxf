package parser

import (
    //"os"
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

// k times f = p

// how do you scale the scale?

// 10 20 - 20 30
// 10 30 - 10 60

// 2

// 20 40 - 40 60
// 10 30 - 20 120

// 10, 40
// 30, 120

// xDenom := (bounds.maxX - bounds.minX) / 2
// m.Vertices[i].X = (m.Vertices[i].X - bounds.minX) / xDenom - 1.0

// Block scaling probably will work by parsing the block then look at the inserts and then recalculate the 
// ..., I don't see it right now hopefully this will be clear soon
     

func (p *parser) parseBlocks(gxf *drawing.Gxf) {
    for {
        switch p.consumeNext() {
        case "BLOCK":
            _ = p.parseEntity()

            name := p.consumeNext()
            p.consumeNext() // Block Flag
            
            p.discardIf(10) //anchorX := p.expectNextFloat(10) 
            p.discardIf(20) //anchorY := p.expectNextFloat(20)
            p.discardIf(30)

            p.discardIf(3)
            p.discardIf(1)

            fmt.Println(name)
            fmt.Println(p.code)

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
