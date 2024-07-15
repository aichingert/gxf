package parser

import (
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
		case "VPORT":
			fallthrough
		case "LTYPE":
			fallthrough
		case "LAYER":
			fallthrough
		case "STYLE":
			fallthrough
		case "VIEW":
			fallthrough
		case "UCS":
			fallthrough
		case "APPID":
			fallthrough
		case "DIMSTYLE":
			r.SkipToLabel("ENDTAB")
			return r.Err()

			log.Fatal(r.DxfLine())
		default:
			log.Fatal(Line, r.DxfLine())
		}

	}

	return r.Err()
}
