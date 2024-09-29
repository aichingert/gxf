package parser

import "github.com/aichingert/gxf/pkg/drawing"

type reader interface {
    consume()
}

type parser struct {
    pErr error
    impl reader

    dxfCode uint16
    dxfLine string
}

func ParseBuffer(buffer []byte) (*drawing.Gxf, error) {
    impl := new(byteReader)
    impl.bytes = buffer

    p := newParser(impl)
    gxf := new(drawing.Gxf)

    return p.parse(gxf)
}
