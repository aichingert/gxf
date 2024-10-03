package parser

import "github.com/aichingert/gxf/pkg/drawing"

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
    for true {
        break L
    }

    return gxf, nil
}

func (p *parser) consume() {
    if p.err != nil {
        return
    }

    p.impl.consumeCode(&p.code, &p.pErr)
    p.impl.consumeLine(&p.line, &p.pErr)
}

func (p *parser) setErr(err error) {
    p.pErr = err
}

func (p *parser) err() error {
    return p.pErr
}
