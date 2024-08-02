package parser

import (
	"github.com/aichingert/dxf/pkg/drawing"
	"os"
)

var Line uint64 = 0

const (
	DecRadix = 10
	HexRadix = 16
)

type dxfLine struct {
	code uint16
	line string
}

// FIXME: make reader struct which has most of these functions
// and add a dxfLine extractor interface which will be the
// difference from the file reader and the byte reader

type Reader interface {
	line() string
	code() uint16

	consume()
	consumeNext() string
	consumeUntil(string)
	assertNextLine(string) error
	consumeStr(*string)
	consumeFloat(uint16, string, *float64)
	consumeNumber(uint16, int, string, *int64)
	consumeCoordinates([]float64)

	consumeStrIf(uint16, *string) bool
	consumeFloatIf(uint16, string, *float64) bool
	consumeNumberIf(uint16, int, string, *int64) bool
	consumeCoordinatesIf(uint16, []float64) bool

	setErr(error)
	Err() error
}

func FromFile(filename string) (*drawing.Dxf, error) {
	dxf := drawing.New(filename)
	reader, file := NewFileReader(filename)
	reader.consume()
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	return parse(reader, dxf)
}

func FromBuffer(filename string, buffer []byte) (*drawing.Dxf, error) {
	dxf := drawing.New(filename)
	reader := NewByteReader(buffer)
	reader.consume()

	return parse(reader, dxf)
}

func parse(r Reader, dxf *drawing.Dxf) (*drawing.Dxf, error) {
L:
	for {
		switch r.consumeNext() {
		case "SECTION":
		case "HEADER":
			ParseHeader(r, dxf)
		case "TABLES":
			ParseTables(r, dxf)
		case "BLOCKS":
			ParseBlocks(r, dxf)
		case "ENTITIES":
			ParseEntities(r, dxf)
		case "EOF":
			break L
		default:
			if r.Err() != nil {
				return nil, r.Err()
			}

			r.consumeUntil("ENDSEC")
		}
	}

	return dxf, nil
}
