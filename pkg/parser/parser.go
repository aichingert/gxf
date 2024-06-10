package parser

import (
    "log"

    "github.com/aichingert/dxf/pkg/drawing"
)

var Line int64

func FromFile(filename string) *drawing.Dxf {
    dxf     := drawing.New(filename)
    reader, file := NewReader(filename)
    defer file.Close()

    for {
        switch data := reader.ConsumeDxfLine(); data.Line {
        case "SECTION":
            switch section := reader.ConsumeDxfLine(); section.Line {
            case "HEADER":
                ParseHeader(reader, dxf)
            case "BLOCKS":
                ParseBlocks(reader, dxf) 
            case "ENTITIES":
                ParseEntities(reader, dxf)
            default:
                log.Println("WARNING: section not implemented: ", section)
                reader.SkipToLabel("ENDSEC")
            }
        case "EOF":
            return dxf
        default:
            log.Fatal("HERE", data.Line,"HERE")
        }
    }
}
