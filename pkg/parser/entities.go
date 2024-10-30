package parser

import (
    "strings"

    "github.com/aichingert/gxf/pkg/drawing"
)

func (p *parser) parseEntities(layers map[string][]uint8, blocks map[string]*drawing.Block) *drawing.Block {
    block := drawing.NewBlock()

    for {
        switch p.consumeNext() {
        case "LINE":
            p.consumeLine(layers, block.Lines, block.Bounds)
        case "LWPOLYLINE":
            p.consumePolyline(layers, block.Lines, block.Bounds)
        case "ENDSEC":
            return block
        default:
        }

        if p.err != nil {
            return nil
        }

        for p.code != 0 {
            p.consume()
        }
    }
}

func (p *parser) parseEntity() string {
    for !strings.HasPrefix(p.consumeNext(), "AcDb") {
    }

    p.discardIf(67)
    line := p.consumeNext()

    for !strings.HasPrefix(p.consumeNext(), "AcDb") {
    }

    return line
}

func (p *parser) consumeLine(layers map[string][]uint8, lines *drawing.Mesh, bnds *drawing.Bounds) {
    layer := p.parseEntity()
    p.discardIf(39)

    srcX := p.expectNextFloat(10)
    srcY := p.expectNextFloat(20)
    p.discardIf(30)

    dstX := p.expectNextFloat(11)
    dstY := p.expectNextFloat(21)
    p.discardIf(31)

    lines.Vertices = append(lines.Vertices, drawing.NewVertex(srcX, srcY, layers[layer]))
    lines.Vertices = append(lines.Vertices, drawing.NewVertex(dstX, dstY, layers[layer]))

    bnds.UpdateX([]float32{srcX, dstX})
    bnds.UpdateY([]float32{srcY, dstY})
}

func (p *parser) consumePolyline(layers map[string][]uint8, lines *drawing.Mesh, bnds *drawing.Bounds) {
    layer := p.parseEntity()

    vertices := p.expectNextInt(90, decRadix)
    if vertices < 0 { return }

    flag := p.expectNextInt(70, decRadix)
    p.discardIf(43) // width for each vertex

    xs := []float32{}
    ys := []float32{}
    l := 0

    for i := uint32(0); i < vertices; i++ { 
        xs = append(xs, p.expectNextFloat(10))
        ys = append(ys, p.expectNextFloat(20))
        p.discardIf(30)

        p.discardIf(40)
        p.discardIf(41)

        // TODO: calculate points for bulge
        p.discardIf(42) 
        p.discardIf(91)

        l = len(xs)

        if l > 1 {
            lines.Vertices = append(lines.Vertices, drawing.NewVertex(xs[l - 2], ys[l - 2], layers[layer]))
            lines.Vertices = append(lines.Vertices, drawing.NewVertex(xs[l - 1], ys[l - 1], layers[layer]))
        }
    }

    if flag & 1 == 1 {
        lines.Vertices = append(lines.Vertices, drawing.NewVertex(xs[l - 1], ys[l - 1], layers[layer]))
        lines.Vertices = append(lines.Vertices, drawing.NewVertex(xs[0], ys[0], layers[layer]))
    }

    bnds.UpdateX(xs)
    bnds.UpdateY(ys)
}

func (p *parser) consumeInsert(block *drawing.Block, blocks map[string]*drawing.Block) {
    _ = p.parseEntity()

    p.discardIf(66)
    name := p.consumeNext()
    x := p.expectNextFloat(10)
    y := p.expectNextFloat(20)
    p.discardIf(30)

    xs := p.expectNextFloat(41)
    ys := p.expectNextFloat(42)
    p.discardIf(43)

    rot := p.expectNextFloat(50)

    _ = name
    _ = x
    _ = y
    _ = xs
    _ = ys
    _ = rot

    _ = block
    _ = blocks
}
