package parser

import (
    "fmt"

    "github.com/aichingert/gxf/pkg/color"
    "github.com/aichingert/gxf/pkg/drawing"
)

func (p *parser) parseTables(gxf *drawing.Gxf) map[string][]uint8 {
    layers := make(map[string][]uint8)

    for {
        switch p.consumeNext() {
        case "TABLE":
            if p.consumeNext() != "LAYER" {
                p.consumeUntil("ENDTAB")
                continue
            }

            for {
                p.consumeUntilPrefix("AcDbSymbolTable")
                if p.consumeNext() == "AcDbLayerTableRecord" {
                    layerName := p.consumeNext()
                    p.discardIf(70)
                    colorIdx := p.expectNextInt(62, decRadix)
                    
                    layers[layerName] = color.DxfColorToRGB[colorIdx]
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
            p.err = NewParseError(fmt.Sprintf("invalid table value %s", p.line))
        }

        if p.err != nil {
            return layers
        }

        for p.code != 0 {
            p.consume()
        }
    }
}
