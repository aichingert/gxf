package parser

import (
    "strings"

    "github.com/aichingert/gxf/pkg/drawing"
)

func (p *parser) parseEntities(lines *drawing.Mesh, polygon *drawing.Mesh) {
    for {
        switch p.consumeNext() {
        case "LINE":
            p.consumeLine(lines)
        case "LWPOLYLINE":
            p.consumePolyline(lines)
        case "ENDSEC":
            return
        default:
        }

        if p.err != nil {
            return
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

func (p *parser) consumeLine(lines *drawing.Mesh) {
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

    gxf.UpdateBorder(srcX, dstX, srcY, dstY)

    lines.Vertices = append(lines.Vertices, drawing.Vertex{ X: srcX, Y: srcY })
    lines.Vertices = append(lines.Vertices, drawing.Vertex{ X: dstX, Y: dstY })
}

func (p *parser) consumePolyline(lines *drawing.Mesh) {
    p.parseEntity()

    vertices := p.expectNextInt(90, decRadix)
    flag     := p.expectNextInt(70, decRadix)
    p.discardIf(43) // width for each vertex

    srcX := float32(0)
    srcY := float32(0)
    
    if vertices > 0 {
        srcX = p.expectNextFloat(10)
        srcY = p.expectNextFloat(20)
        p.discardIf(30) // z
        p.discardIf(40) // start width
        p.discardIf(41) // end   width
        p.discardIf(42) // TODO: calculate points for bulge
        p.discardIf(91) // vertex ident
    }

    prvX := srcX
    prvY := srcY

    for i := uint32(1); i < vertices; i++ { 
        nxtX := p.expectNextFloat(10)
        nxtY := p.expectNextFloat(20)
        p.discardIf(30) // z

        gxf.UpdateBorder(prvX, nxtX, prvY, nxtY)

        p.discardIf(40) // start width
        p.discardIf(41) // end   width
        p.discardIf(42) // TODO: calculate points for bulge

        p.discardIf(91) // vertex ident

        lines.Vertices = append(lines.Vertices, drawing.Vertex{ X: prvX, Y: prvY })
        lines.Vertices = append(lines.Vertices, drawing.Vertex{ X: nxtX, Y: nxtY })

        prvX = nxtX
        prvY = nxtY
    }

    if flag & 1 == 1 {
        lines.Vertices = append(lines.Vertices, drawing.Vertex{ X: prvX, Y: prvY })
        lines.Vertices = append(lines.Vertices, drawing.Vertex{ X: srcX, Y: srcY })
    }
}
