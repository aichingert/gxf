package parser

import (
	"github.com/aichingert/dxf/pkg/table"
	"log"

	"github.com/aichingert/dxf/pkg/drawing"
	_ "github.com/aichingert/dxf/pkg/entity"
)

func ParseTables(r *Reader, dxf *drawing.Dxf) error {
	for r.ScanDxfLine() {
		switch r.DxfLine().Line {
		case "TABLE":
			Wrap(parseTable, r, dxf)
		case "ENDSEC":
			return r.Err()
		default:
			log.Println("Table(", Line, "): ", r.DxfLine().Line)
			return r.Err()
		}

		if WrappedErr != nil {
			return WrappedErr
		}
	}

	return r.Err()
}

func parseTable(r *Reader, dxf *drawing.Dxf) error {
	for r.ScanDxfLine() {
		switch r.DxfLine().Line {
		case "LAYER":
			if parseLayerTable(r, dxf) != nil {
				return r.Err()
			}
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
			_ = r.SkipToLabel("ENDTAB")
			return r.Err()
		case "ENDTAB":
			return r.Err()
		default:
			log.Fatal(Line, r.DxfLine())
		}

		peek, err := r.PeekCode()
		for err == nil && peek != 0 {
			r.ConsumeStr(nil)
			peek, err = r.PeekCode()
		}
	}

	return r.Err()
}

// TODO: maybe put this in acdb as well
func parseLayerTable(r *Reader, dxf *drawing.Dxf) error {
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
		return parseAcDbLayerTableRecord(r, dxf)
	} else {
		r.ConsumeNumber(70, DecRadix, "standard flag", nil)
	}

	return r.Err()
}

func parseAcDbLayerTableRecord(r *Reader, dxf *drawing.Dxf) error {
	if r.AssertNextLine("AcDbLayerTableRecord") != nil {
		return r.Err()
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

	if r.Err() == nil {
		dxf.Layers[layerName] = layer
	}

	return r.Err()
}
