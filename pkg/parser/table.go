package parser

import (
	"fmt"
	"github.com/aichingert/dxf/pkg/table"
	"log"

	"github.com/aichingert/dxf/pkg/drawing"
	_ "github.com/aichingert/dxf/pkg/entity"
)

func ParseTables(r *Reader, dxf *drawing.Dxf) {
	for r.ScanDxfLine() {
		switch r.DxfLine().Line {
		case "TABLE":
			parseTable(r, dxf)
		case "ENDSEC":
			return
		default:
			r.err = NewParseError(fmt.Sprintf("Table(%d): %s", Line, r.DxfLine().Line))
			return
		}
	}
}

func parseTable(r *Reader, dxf *drawing.Dxf) {
	for r.ScanDxfLine() {
		switch r.DxfLine().Line {
		case "LAYER":
			parseLayerTable(r, dxf)
		case "VPORT":
			fallthrough
		case "LTYPE":
			fallthrough
		case "STYLE":
			fallthrough
		case "VIEW":
			fallthrough
		case "UCS":
			fallthrough
		case "APPID":
			fallthrough
		case "BLOCK_RECORD":
			fallthrough
		case "DIMSTYLE":
			r.SkipToLabel("ENDTAB")
			return
		case "ENDTAB":
			return
		default:
			log.Fatal(Line, r.DxfLine())
		}

		peek, err := r.PeekCode()
		for err == nil && peek != 0 {
			r.ConsumeStr(nil)
			peek, err = r.PeekCode()
		}
	}
}

// TODO: maybe put this in acdb as well
func parseLayerTable(r *Reader, dxf *drawing.Dxf) {
	r.ConsumeNumber(5, HexRadix, "handle", nil)
	// TODO: set hard owner/handle to owner dictionary
	if r.ConsumeStrIf(102, nil) { // consumeIf => ex. {ACAD_XDICTIONARY
		r.ConsumeStr(nil) // 360 => hard owner
		for r.ConsumeNumberIf(330, HexRadix, "soft owner", nil) {
		}
		r.ConsumeStr(nil) // 102 }
	}

	if r.ConsumeStrIf(102, nil) { // consumeIf => ex. {ACAD_XDICTIONARY
		r.ConsumeStr(nil) // 360 => hard owner
		for r.ConsumeNumberIf(330, HexRadix, "soft owner", nil) {
		}
		r.ConsumeStr(nil) // 102 }
	}
	r.ConsumeNumber(330, HexRadix, "owner ptr", nil)

	next := ""
	r.ConsumeStr(&next)

	if next == "AcDbSymbolTableRecord" {
		parseAcDbLayerTableRecord(r, dxf)
	} else {
		r.ConsumeNumber(70, DecRadix, "standard flag", nil)
	}
}

func parseAcDbLayerTableRecord(r *Reader, dxf *drawing.Dxf) {
	if r.AssertNextLine("AcDbLayerTableRecord") != nil {
		return
	}

	layer := table.NewLayer()
	layerName := ""

	r.ConsumeStr(&layerName) // [2]
	r.ConsumeNumber(70, DecRadix, "standard flag", nil)
	r.ConsumeNumber(62, DecRadix, "Layer color", &layer.Color)
	r.ConsumeNumberIf(420, DecRadix, "layer true color", &layer.TrueColor)
	r.ConsumeStr(&layer.LineType)

	r.ConsumeNumberIf(290, DecRadix, "Plotting flag", nil)
	r.ConsumeNumber(370, DecRadix, "line weight enum value", nil)
	r.ConsumeNumber(390, DecRadix, "hard-pointer Id/handle to plot style name object", nil)
	r.ConsumeNumber(347, DecRadix, "hard-pointer Id/handle to material object", nil)
	r.ConsumeNumber(348, DecRadix, "not documented", nil)

	dxf.Layers[layerName] = layer
}
