package parser

import "github.com/aichingert/gxf/pkg/drawing"

type reader interface {
    consumeCode(code *uint16) error
    consumeLine(line *string) error
}

type parser struct {
    pErr error
    impl reader

    code uint16
    line string
}

func ParseBuffer(buffer []byte) (*drawing.Gxf, error) {
    impl := new(byteReader)
    impl.bytes = buffer

    p := newParser(impl)
    gxf := new(drawing.Gxf)

    return p.parse(gxf)
}
