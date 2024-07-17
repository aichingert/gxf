package parser

import (
    "github.com/aichingert/dxf/pkg/drawing"
    "log"
)

var (
	Handle uint64 = 0
	Owner  uint64 = 0
	Line   uint64 = 0

	WrappedErr error
)

type ParseFunction func(*Reader, *drawing.Dxf) error

func Wrap(fn ParseFunction, r *Reader, dxf *drawing.Dxf) {
	if WrappedErr != nil {
		return
	}

	WrappedErr = fn(r, dxf)
}

func FromFile(filename string) (*drawing.Dxf, error) {
	dxf := drawing.New(filename)
	reader, file, err := NewReader(filename)

	if err != nil {
		return nil, err
	}
	defer file.Close()

	for reader.ScanDxfLine() {
		switch reader.DxfLine().Line {
		case "SECTION":
			section, err := reader.ConsumeDxfLine()
			if err != nil {
				return dxf, err
			}

            log.Println(section)

			switch section.Line {
			case "HEADER":
				Wrap(ParseHeader, reader, dxf)
           case "TABLES":
               Wrap(ParseTables, reader, dxf)
			case "BLOCKS":
				Wrap(ParseBlocks, reader, dxf)
			case "ENTITIES":
				if err := ParseEntities(reader, dxf.EntitiesData); err != nil {
					return dxf, err
				}
			default:
				reader.SkipToLabel("ENDSEC")
			}
		case "EOF":
			return dxf, reader.Err()
		default:
			return nil, NewParseError("unexpected")
		}

		if WrappedErr != nil {
			return dxf, WrappedErr
		}
	}

	return dxf, reader.Err()
}
