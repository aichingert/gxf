package parser

import (
    //"os"
    "fmt"

    "github.com/aichingert/gxf/pkg/drawing"
)

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
