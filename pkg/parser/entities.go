package parser

import (
    "strings"

    "github.com/aichingert/gxf/pkg/drawing"
)

func (p *parser) parseEntities(gxf *drawing.Gxf, layers map[string][]uint8) {
    for {
        switch p.consumeNext() {
        case "LINE":
            p.consumeLine(gxf, layers)
        case "LWPOLYLINE":
            p.consumePolyline(gxf, layers)
        case "INSERT":
            p.consumeInsert(gxf)
        case "ENDSEC":
            return
        case "ENDBLK":
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

func (p *parser) parseEntity() string {
    for !strings.HasPrefix(p.consumeNext(), "AcDb") {
    }

    p.discardIf(67)
    line := p.consumeNext()

    for !strings.HasPrefix(p.consumeNext(), "AcDb") {
    }

    return line
}

func (p *parser) consumeLine(gxf *drawing.Gxf, layers map[string][]uint8) {
    layer := p.parseEntity()
    p.discardIf(39)

    gxf.Data.AddVertex(1, p.expectNextFloat(10), p.expectNextFloat(20), layers[layer])
    p.discardIf(30)

    gxf.Data.AddVertex(1, p.expectNextFloat(11), p.expectNextFloat(21), layers[layer])
    p.discardIf(31)
}

func (p *parser) consumePolyline(gxf *drawing.Gxf, layers map[string][]uint8) {
    layer := p.parseEntity()

    vertices := p.expectNextInt(90, decRadix)
    if vertices < 0 { return }

    flag := p.expectNextInt(70, decRadix)
    p.discardIf(43)

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
            gxf.Data.AddVertex(1, xs[l - 2], ys[l - 2], layers[layer])
            gxf.Data.AddVertex(1, xs[l - 1], ys[l - 1], layers[layer])
        }
    }

    if flag & 1 == 1 {
        gxf.Data.AddVertex(1, xs[l - 1], ys[l - 1], layers[layer])
        gxf.Data.AddVertex(1, xs[0], ys[0], layers[layer])
    }
}

func (p *parser) consumeInsert(gxf *drawing.Gxf) {
    _ = p.parseEntity()
    p.discardIf(66)

    name := p.consumeNext()

    x := p.expectNextFloat(10)
    y := p.expectNextFloat(20)
    p.discardIf(30)

    sx := p.consumeFloatIf(41, 1.0)
    sy := p.consumeFloatIf(42, 1.0)
    p.discardIf(43)

    rot := p.consumeFloatIf(50, 0.0)
    instance := gxf.BlockNameRes[name]

    gxf.InstanceData[instance] = append(gxf.InstanceData[instance], [4]float32{x, y, sx, sy})

    _ = rot
}
