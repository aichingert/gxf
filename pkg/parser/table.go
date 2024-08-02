package parser

import (
	"fmt"
	"github.com/aichingert/dxf/pkg/drawing"
	_ "github.com/aichingert/dxf/pkg/entity"
	"github.com/aichingert/dxf/pkg/table"
)

func ParseTables(r Reader, dxf *drawing.Dxf) {
	for {
		switch r.consumeNext() {
		case "TABLE":
			parseTable(r, dxf)
		case "ENDSEC":
			return
		default:
			if r.Err() == nil {
				r.setErr(NewParseError(fmt.Sprintf("Table(%d): %s", Line, r.line())))
			}

			return
		}
	}
}

func parseTable(r Reader, dxf *drawing.Dxf) {
	for {
		switch r.consumeNext() {
		case "LAYER":
			parseLayerTable(r, dxf)
		case "ENDTAB":
			return
		default:
			if r.Err() != nil {
				return
			}

			r.consumeUntil("ENDTAB")
			return
		}

		for r.code() != 0 {
			r.consume()
		}
	}
}

// TODO: maybe put this in acdb as well
func parseLayerTable(r Reader, dxf *drawing.Dxf) {
	r.consumeNumber(5, HexRadix, "handle", nil)

	// TODO: set hard owner/handle to owner dictionary
	if r.consumeStrIf(102, nil) { // consumeIf => ex. {ACAD_XDICTIONARY
		r.consumeStr(nil) // 360 => hard owner
		for r.consumeNumberIf(330, HexRadix, "soft owner", nil) {
		}
		r.consumeStr(nil) // 102 }
	}

	if r.consumeStrIf(102, nil) { // consumeIf => ex. {ACAD_XDICTIONARY
		r.consumeStr(nil) // 360 => hard owner
		for r.consumeNumberIf(330, HexRadix, "soft owner", nil) {
		}
		r.consumeStr(nil) // 102 }
	}
	r.consumeNumber(330, HexRadix, "owner ptr", nil)

	if r.consumeNext() == "AcDbSymbolTableRecord" {
		parseAcDbLayerTableRecord(r, dxf)
	} else {
		r.consumeNumber(70, DecRadix, "standard flag", nil)
	}
}

func parseAcDbLayerTableRecord(r Reader, dxf *drawing.Dxf) {
	if r.assertNextLine("AcDbLayerTableRecord") != nil {
		return
	}

	layer := table.NewLayer()
	layerName := ""

	r.consumeStr(&layerName) // [2]
	r.consumeNumber(70, DecRadix, "standard flag", nil)
	r.consumeNumber(62, DecRadix, "Layer color", &layer.Color)
	r.consumeNumberIf(420, DecRadix, "layer true color", &layer.TrueColor)
	r.consumeStr(&layer.LineType)

	r.consumeNumberIf(290, DecRadix, "Plotting flag", nil)
	r.consumeNumber(370, DecRadix, "line weight enum value", nil)
	r.consumeNumber(390, DecRadix, "hard-pointer Id/handle to plot style name object", nil)
	r.consumeNumber(347, DecRadix, "hard-pointer Id/handle to material object", nil)
	r.consumeNumber(348, DecRadix, "not documented", nil)

	dxf.Layers[layerName] = layer
}
