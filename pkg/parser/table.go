package parser

import (
    "fmt"

    _ "github.com/aichingert/gxf/pkg/drawing"
)

func (p *parser) parseTables() map[string]uint8 {
    layers := make(map[string]uint8)

    for {
        switch p.consumeNext() {
        case "TABLE":
            tableType := p.consumeNext()

            if tableType != "LAYER" {
                p.consumeUntil("ENDTAB")
                continue
            }

            for {
                p.consumeUntilPrefix("AcDbSymbolTable")
                if p.consumeNext() == "AcDbLayerTableRecord" {
                    layerName := p.consumeNext()
                    p.discardIf(70)
                    color := p.expectNextInt(62, decRadix)

                    layers[layerName] = uint8(color)
                }

                for p.code != 0 {
                    p.consume()
                }

                if p.line != "LAYER" {
                    break
                }
            }

            p.consumeUntil("ENDTAB")
        case "ENDSEC":
            return layers
        default:
            p.err = NewParseError(fmt.Sprintf("invalid block value %s", p.line))
        }

        if p.err != nil {
            return layers
        }

        for p.code != 0 {
            p.consume()
        }
    }
}
