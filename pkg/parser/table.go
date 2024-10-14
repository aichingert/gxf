package parser

import (
    "fmt"

    "github.com/aichingert/gxf/pkg/color"
    "github.com/aichingert/gxf/pkg/drawing"
)

func (p *parser) parseTables(gxf *drawing.Gxf) {
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
                    colorIdx := p.expectNextInt(62, decRadix)
                    
                    gxf.Layers[layerName] = color.DxfColorToRGB[colorIdx]
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
