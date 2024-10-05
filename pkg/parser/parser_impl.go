package parser

import (
    "fmt"
    "strconv"

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
L:
    for {
        switch p.consumeNext() {
        case "SECTION":
        case "HEADER":
            p.consumeUntil("ENDSEC")
        case "TABLES":
            p.consumeUntil("ENDSEC")
        case "BLOCKS":
            p.consumeUntil("ENDSEC")
        case "ENTITIES":
            p.parseEntities(gxf)
        case "EOF":
            break L
        default:
            if p.err != nil {
                return nil, p.err
            }

            p.consumeUntil("ENDSEC")
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
    for p.consumeNext() != label {}
}

func (p *parser) expectNextFloat(code uint16) float32 {
    if p.err != nil {
        return 0.0
    }
    defer p.consume()

    if p.code != code {
        p.err = NewParseError(fmt.Sprintf("Expect float invalid code: expected %d got %d", code, p.code))
        return 0.0
    }

    f32, err := strconv.ParseFloat(p.line, 32)
    p.err = err

    return float32(f32)
}
