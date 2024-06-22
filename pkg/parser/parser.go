package parser

import (
    "log"

    "github.com/aichingert/dxf/pkg/drawing"
)

var Line int64

func FromFile(filename string) (*drawing.Dxf, error) {
    dxf     := drawing.New(filename)
    reader, file, err := NewReader(filename)

    if err != nil { return nil, err }
    defer file.Close()

    for {
        data, err := reader.ConsumeDxfLine()
        if err != nil { return dxf, err }

        switch data.Line {
        case "SECTION":
            section, err := reader.ConsumeDxfLine()
            if err != nil { return dxf, err }

            switch section.Line {
            case "HEADER":
                if err = ParseHeader(reader, dxf); err != nil { return dxf, err }
            case "BLOCKS":
                if err = ParseBlocks(reader, dxf); err != nil { return dxf, err }
            case "ENTITIES":
                if err = ParseEntities(reader, dxf); err != nil { return dxf, err }
            default:
                log.Println("WARNING: section not implemented: ", section)
                reader.SkipToLabel("ENDSEC")
            }
        case "EOF":
            return dxf, nil
        default:
            return nil, NewParseError("unexpected")
        }
    }
}
