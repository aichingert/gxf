package parser

import (
    "strings"

    "github.com/aichingert/gxf/pkg/drawing"
)

func (p *parser) parseEntities() *drawing.Mesh {
    mesh := drawing.NewMesh()
    bnds := drawing.NewBounds()

    for {
        switch p.consumeNext() {
        case "LINE":
            p.consumeLine(mesh, bnds)
        case "LWPOLYLINE":
            p.consumePolyline(mesh, bnds)
        case "ENDSEC":
            mesh.Scale(bnds)
            return mesh
        default:
        }

        if p.err != nil {
            return mesh
        }

        for p.code != 0 {
            p.consume()
        }
    }
}

func (p *parser) parseEntity() uint8 {
    for !strings.HasPrefix(p.consumeNext(), "AcDb") {
    }
    // NOTE(code 67): is an optional paper space visibility and is not used so we skip it
    if p.code == 67 {
        p.consume()
    }

    line := p.consumeNext()

    // NOTE(code 62): defines the color number (present if not by layer)
    // NOTE(code 48): defines the linetype scale (maybe we find files that actualy use it)

    for !strings.HasPrefix(p.consumeNext(), "AcDb") {
    }

    _ = line
    return 0
}

func (p *parser) consumeLine(lines *drawing.Mesh, bnds *drawing.Bounds) {
    p.parseEntity()
    // NOTE(code 39): not contained by any files I tested, stands for thickness
    if p.code == 39 {
        p.consume()
    }

    srcX := p.expectNextFloat(10)
    srcY := p.expectNextFloat(20)
    p.discardIf(30) // z

    dstX := p.expectNextFloat(11)
    dstY := p.expectNextFloat(21)
    p.discardIf(31) // z

    lines.Vertices = append(lines.Vertices, drawing.NewVertex(srcX, srcY))
    lines.Vertices = append(lines.Vertices, drawing.NewVertex(dstX, dstY))
    bnds.UpdateX([]float32{srcX, dstX})
    bnds.UpdateY([]float32{srcY, dstY})
}

func (p *parser) consumePolyline(lines *drawing.Mesh, bnds *drawing.Bounds) {
    p.parseEntity()

    vertices := p.expectNextInt(90, decRadix)
    if vertices < 0 { return }

    flag := p.expectNextInt(70, decRadix)
    p.discardIf(43) // width for each vertex

    xs := []float32{}
    ys := []float32{}

    for i := uint32(0); i < vertices; i++ { 
        xs = append(xs, p.expectNextFloat(10))
        ys = append(ys, p.expectNextFloat(20))
        p.discardIf(30) // z

        p.discardIf(40) // start width
        p.discardIf(41) // end   width
        p.discardIf(42) // TODO: calculate points for bulge

        p.discardIf(91) // vertex ident
        l := len(xs)

        if l > 1 {
            lines.Vertices = append(lines.Vertices, drawing.NewVertex(xs[l - 2], ys[l - 2]))
            lines.Vertices = append(lines.Vertices, drawing.NewVertex(xs[l - 1], ys[l - 1]))
        }
    }

    if flag & 1 == 1 {
        lines.Vertices = append(lines.Vertices, drawing.NewVertex(xs[len(xs) - 2], ys[len(xs) - 2]))
        lines.Vertices = append(lines.Vertices, drawing.NewVertex(xs[0], ys[0]))
    }

    bnds.UpdateX(xs)
    bnds.UpdateY(ys)
}
