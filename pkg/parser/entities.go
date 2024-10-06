package parser

import (
    "strings"

    "github.com/aichingert/gxf/pkg/drawing"
)

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

func (p *parser) parseEntities(gxf *drawing.Gxf) {
    for {
        switch p.consumeNext() {
        case "LINE":
            p.consumeLine(gxf)
        case "POLYLINE":
            p.parseEntity()
        case "ENDSEC":
            return
        default:
        }

        for p.code != 0 {
            p.consume()
        }
    }
}

func max(a float32, b float32) float32 {
    if a > b {
        return a
    } else {
        return b
    }
}

func min(a float32, b float32) float32 {
    if a < b {
        return a
    } else {
        return b
    }
}

func (p *parser) consumeLine(gxf *drawing.Gxf) {
    p.parseEntity()
    // NOTE(code 39): not contained by any files I tested, stands for thickness
    if p.code == 39 {
        p.consume()
    }

    srcX := p.expectNextFloat(10)
    srcY := p.expectNextFloat(20)
    _ = p.expectNextFloat(30)

    dstX := p.expectNextFloat(11)
    dstY := p.expectNextFloat(21)
    _ = p.expectNextFloat(31)

    gxf.BorderX[0] = min(gxf.BorderX[0], srcX)
    gxf.BorderX[0] = min(gxf.BorderX[0], dstX)
    gxf.BorderX[1] = max(gxf.BorderX[1], srcX)
    gxf.BorderX[1] = max(gxf.BorderX[1], dstX)

    gxf.BorderY[0] = min(gxf.BorderY[0], srcY)
    gxf.BorderY[0] = min(gxf.BorderY[0], dstY)
    gxf.BorderY[1] = max(gxf.BorderY[1], srcY)
    gxf.BorderY[1] = max(gxf.BorderY[1], dstY)

    gxf.Lines.Vertices = append(gxf.Lines.Vertices, drawing.Vertex{ X: srcX, Y: srcY })
    gxf.Lines.Vertices = append(gxf.Lines.Vertices, drawing.Vertex{ X: dstX, Y: dstY })
}
