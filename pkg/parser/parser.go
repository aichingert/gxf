package parser

import (
	"github.com/aichingert/dxf/pkg/drawing"
	"log"
	"os"
)

var Line uint64 = 0

func FromFile(filename string) (*drawing.Dxf, error) {
	dxf := drawing.New(filename)
	reader, file, err := NewReader(filename)

	if err != nil {
		return nil, err
	}

	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	for reader.ScanDxfLine() {
		switch reader.DxfLine().Line {
		case "SECTION":
			section := reader.ConsumeDxfLine()
			if reader.err != nil {
				return dxf, err
			}

			log.Println(section)

			switch section.Line {
			case "HEADER":
				ParseHeader(reader, dxf)
			case "TABLES":
				ParseTables(reader, dxf)
			case "BLOCKS":
				ParseBlocks(reader, dxf)
			case "ENTITIES":
				ParseEntities(reader, dxf.EntitiesData)
			default:
				reader.SkipToLabel("ENDSEC")
			}
		case "EOF":
			return dxf, reader.Err()
		default:
			return nil, NewParseError("unexpected")
		}
	}

	return dxf, reader.Err()
}
