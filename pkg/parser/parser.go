package parser

import (
	"github.com/aichingert/dxf/pkg/drawing"
	"os"
)

var Line uint64 = 0

func FromFile(filename string) (*drawing.Dxf, error) {
	dxf := drawing.New(filename)
	reader, file := NewReader(filename)
	reader.consume()
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

L:
	for {
		switch reader.consumeNext() {
		case "SECTION":
		case "HEADER":
			ParseHeader(reader, dxf)
		case "TABLES":
			ParseTables(reader, dxf)
		case "BLOCKS":
			ParseBlocks(reader, dxf)
		case "ENTITIES":
			ParseEntities(reader, dxf)
		case "EOF":
			break L
		default:
			if reader.err != nil {
				return nil, reader.err
			}

			reader.consumeUntil("ENDSEC")
		}
	}

	return dxf, nil
}
