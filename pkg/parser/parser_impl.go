package parser

import (
    "fmt"
    "strconv"
    "strings"

    "github.com/aichingert/gxf/pkg/drawing"
)

var Line uint64 = 0

const (
    decRadix = 10
)

func newParser(impl reader) *parser {
    return &parser{
        err: nil,
        impl: impl,

        code: 0,
        line: "",
    }
}

func (p *parser) parse(gxf *drawing.Gxf) (*drawing.Gxf, error) {
    var layers map[string][]uint8

L:
    for {
        switch p.consumeNext() {
        case "SECTION":
        case "TABLES":
            layers = p.parseTables(gxf)
        case "BLOCKS":
            p.parseBlocks(gxf, layers)
        case "ENTITIES":
            // TODO: check if parser has error first 
            p.parseEntities(gxf, layers)
        case "EOF":
            break L
        default:
            if p.err != nil {
                return nil, p.err
            }

            p.consumeUntil("ENDSEC")
        }
    }

    offset := 0
    if len(gxf.BlockOffsets) > 0 {
        offset = int(gxf.BlockOffsets[len(gxf.BlockOffsets) - 1][0])
    }

    for l := range (len(gxf.Data.Lines.Vertices) - offset) / 5 {
        if gxf.Data.MinX > gxf.Data.Lines.Vertices[l * 5 + offset] {
            gxf.Data.MinX = gxf.Data.Lines.Vertices[l * 5 + offset] 
        }
        if gxf.Data.MaxX < gxf.Data.Lines.Vertices[l * 5 + offset] {
            gxf.Data.MaxX = gxf.Data.Lines.Vertices[l * 5 + offset] 
        }

        if gxf.Data.MinY > gxf.Data.Lines.Vertices[l * 5 + offset + 1] {
            gxf.Data.MinY = gxf.Data.Lines.Vertices[l * 5 + offset + 1] 
        }
        if gxf.Data.MaxY < gxf.Data.Lines.Vertices[l * 5 + offset + 1] {
            gxf.Data.MaxY = gxf.Data.Lines.Vertices[l * 5 + offset + 1] 
        }
    }

    return gxf, nil
}

func (p *parser) consume() {
    if p.err != nil {
        return
    }

    if err := p.impl.consumeCode(&p.code); err != nil {
        p.err = err
        return
    }

    if err := p.impl.consumeLine(&p.line); err != nil {
        p.err = err
    }
}

func (p *parser) consumeNext() string {
    if p.err != nil {
        return p.err.Error()
    }
    defer p.consume()

    value := p.line
    return value
}

func (p *parser) consumeUntil(label string) {
    for p.err == nil && p.consumeNext() != label {}
}

func (p *parser) consumeUntilPrefix(label string) {
    for !strings.HasPrefix(p.consumeNext(), label) {}
}

func (p *parser) expectNextFloat(code uint16) float32 {
    if p.err != nil {
        return 0.0
    }
    defer p.consume()

    if p.code != code {
        p.err = NewParseError(fmt.Sprintf("Float: expected %d got %d", code, p.code))
        return 0.0
    }

    f32, err := strconv.ParseFloat(p.line, 32)
    p.err = err
    return float32(f32)
}

func (p *parser) consumeFloatIf(code uint16, def float32) float32 {
    if p.err != nil || p.code != code {
        return def
    }
    defer p.consume()

    f32, err := strconv.ParseFloat(p.line, 32)
    p.err = err
    return float32(f32)
}

func (p *parser) expectNextInt(code uint16, radix int) uint32 {
    if p.err != nil {
        return 0
    }
    defer p.consume()

    if p.code != code {
        p.err = NewParseError(fmt.Sprintf("Expect int(invalid code): expected %d got %d", code, p.code))
        return 0
    }

    u32, err := strconv.ParseInt(strings.TrimSpace(p.line), radix, 32)
    p.err = err

    return uint32(u32)
}

func (p *parser) discardIf(code uint16) {
    if p.err != nil {
        return
    }

    if p.code == code {
        p.consume()
    }
}

