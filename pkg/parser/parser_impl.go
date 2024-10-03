package parser

import "github.com/aichingert/gxf/pkg/drawing"

var Line uint64 = 0

const (
    decRadix = 10
)

func newParser(impl reader) *parser {
    return &parser{
        pErr: nil,
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
            // TODO:
        case "EOF":
            break L
        default:
            if p.pErr != nil {
                return nil, p.pErr
            }

            p.consumeUntil("ENDSEC")
        }
    }

    return gxf, nil
}

func (p *parser) consume() {
    if p.pErr != nil {
        return
    }

    if err := p.impl.consumeCode(&p.code); err != nil {
        p.pErr = err
        return
    }
    if err := p.impl.consumeLine(&p.line); err != nil {
        p.pErr = err
    }
}

func (p *parser) consumeNext() string {
    if p.pErr != nil {
        return p.pErr.Error()
    }
    defer p.consume()

    value := p.line
    return value
}

func (p *parser) consumeUntil(label string) {
    for p.consumeNext() != label {}
}

func (p *parser) setErr(err error) {
    p.pErr = err
}

func (p *parser) err() error {
    return p.pErr
}
