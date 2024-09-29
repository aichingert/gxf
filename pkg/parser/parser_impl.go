package parser

import "github.com/aichingert/gxf/pkg/drawing"

func newParser(impl reader) *parser {
    return &parser{
        pErr: nil,
        impl: impl,

        dxfCode: 0,
        dxfLine: "",
    }
}

func (p *parser) parse(gxf *drawing.Gxf) (*drawing.Gxf, error) {
    return gxf, nil
}

func (p *parser) code() uint16 {
    return p.dxfCode
}

func (p *parser) line() string {
    return p.dxfLine
}

func (p *parser) setErr(err error) {
    p.pErr = err
}

func (p *parser) err() error {
    return p.pErr
}
