package parser

import "github.com/aichingert/gxf/pkg/drawing"

type reader interface {
    consumeCode(code *uint16) error
    consumeLine(line *string) error
}

type parser struct {
    err error
    impl reader

    code uint16
    line string

}

func ParseBuffer(buffer []byte) (*drawing.Gxf, error) {
    impl := new(byteReader)
    impl.bytes = buffer

    p := newParser(impl)

    // Idea:
    // Using instancing with the gpu to not store to not make the buffers too large. Have to research
    // if it is costly to instance many different meshes because most of the times there are many blocks.

    // OPTIONAL: using mesh optimizing algorithms to reduce size

    // TODO: after that implement hatch entity

    // figure out how to connect triangles and consider actually using indices

    // TODO: circle smoothness

    // Idea 1:
    // create a parse config struct to specify how smooth corners of a circle should be and polylines with
    // a bulge as well

    // OPTIONAL: just predefine it.

    gxf := drawing.NewGxf()

    return p.parse(gxf)
}
