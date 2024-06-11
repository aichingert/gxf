package parser

import (
    "log"

    "github.com/aichingert/dxf/pkg/entity"
    "github.com/aichingert/dxf/pkg/drawing"
)

func ParseEntities(r *Reader, dxf *drawing.Dxf) {
    for {
        switch variable := r.ConsumeDxfLine(); variable.Line {
        case "LINE":
            ParseLine(r, dxf)
        case "LWPOLYLINE":
            ParsePolyline(r, dxf)
        case "ARC":
            ParseArc(r, dxf)
        case "CIRCLE":
            ParseCircle(r, dxf)
        case "MTEXT":
            ParseMText(r, dxf)
        default:
            log.Fatal("[ENTITIES] ", Line, ": ", variable)
        }
    }
}

func parseAcDbEntityE(r *Reader, entity entity.Entity) {
    _ = r.ConsumeDxfLine()
    optional := r.ConsumeDxfLine()

    // TODO: think about paper space visibility
    if optional.Code != 67 {
        entity.SetLayerName(optional.Line)
        return
    }

    // TODO: could lead to bug with start and end layername - seems like it is always the same
    layerName := r.ConsumeDxfLine()
    entity.SetLayerName(layerName.Line)
}

func extractHandleAndOwner(r *Reader) [2]uint64 {
    handle  := r.ConsumeNumber(5, 16, "handle")

    // TODO: set hard owner/handle to owner dictionary
    if r.PeekCode() == 102 {
        _ = r.ConsumeDxfLine()
        _ = r.ConsumeDxfLine()
        _ = r.ConsumeDxfLine()
    }

    owner   := r.ConsumeNumber(330, 16, "owner ptr")

    return [2]uint64{handle, owner}
}

func ParseLine(r *Reader, dxf *drawing.Dxf) {
    result := extractHandleAndOwner(r)
    line   := entity.NewLine(result[0], result[1])

    parseAcDbEntityE(r, line)

    check := r.ConsumeDxfLine()

    if check.Line != "AcDbLine" {
        log.Fatal("[ENTITIES(", Line, ")] Expected AcDbLine got ", check)
    }

    line.Src = r.ConsumeCoordinates3D()
    line.Dst = r.ConsumeCoordinates3D()

    dxf.Lines = append(dxf.Lines, line)
}

func ParsePolyline(r *Reader, dxf *drawing.Dxf) {
    result   := extractHandleAndOwner(r)
    polyline := entity.NewPolyline(result[0], result[1])

    parseAcDbEntityE(r, polyline)
    check := r.ConsumeDxfLine()

    if check.Line != "AcDbPolyline" {
        log.Fatal("[ENTITIES(", Line, ")] Expected AcDbPolyline got ", check)
    }

    polyline.Vertices = r.ConsumeNumber(90, 10, "number of vertices")
    polyline.Flag = r.ConsumeNumber(70, 10, "polyline flag")

    // expecting code 43 
    if r.ConsumeDxfLine().Code != 43 {
        log.Fatal("[ENTITIES] TODO: implement line width for each vertex")
    }

    for i := uint64(0); i < polyline.Vertices; i++ {
        // TODO: sometimes there is a bulge value for a vertex
        coords := r.ConsumeCoordinates2D()

        // bulge = groupcode 42
        if r.PeekCode() == 42 {
            bulge := ParseFloat(r.ConsumeDxfLine().Line)
            polyline.PolylineAppendCoordinate(coords, bulge)
            continue
        }

        polyline.PolylineAppendCoordinate(coords, 0.0)
    }

    dxf.Polylines = append(dxf.Polylines, polyline)
}

func ParseArc(r *Reader, dxf *drawing.Dxf) {
    result  := extractHandleAndOwner(r)
    arc     := entity.NewArc(result[0], result[1])

    parseAcDbEntityE(r, arc)
    check := r.ConsumeDxfLine()

    if check.Line != "AcDbCircle" {
        log.Fatal("[ENTITIES(", Line, ")] Expected AcDbCircle got ", check)
    }

    arc.Circle = &entity.Circle {
        Coordinates:    r.ConsumeCoordinates3D(),
        Radius:         ParseFloat(r.ConsumeDxfLine().Line),
    }

    check = r.ConsumeDxfLine()

    if check.Line != "AcDbArc" {
        log.Fatal("[ENTITIES(", Line, ")] Expected AcDbArc got ", check)
    }

    arc.StartAngle  = ParseFloat(r.ConsumeDxfLine().Line)
    arc.EndAngle    = ParseFloat(r.ConsumeDxfLine().Line)

    dxf.Arcs = append(dxf.Arcs, arc)
}

func ParseCircle(r *Reader, dxf *drawing.Dxf) {
    result := extractHandleAndOwner(r)
    circle := entity.NewCircle(result[0], result[1])

    parseAcDbEntityE(r, circle)
    check := r.ConsumeDxfLine()

    if check.Line != "AcDbCircle" {
        log.Fatal("[ENTITIES(", Line, ")] Expected AcDbCircle got ", check)
    }

    circle.Coordinates = r.ConsumeCoordinates3D()
    circle.Radius      = ParseFloat(r.ConsumeDxfLine().Line)

    dxf.Circles = append(dxf.Circles, circle)
}

func ParseMText(r *Reader, dxf *drawing.Dxf) {
    result  := extractHandleAndOwner(r)
    mText   := entity.NewMText(result[0], result[1])

    parseAcDbEntityE(r, mText)
    check := r.ConsumeDxfLine()

    if check.Line != "AcDbMText" {
        log.Fatal("[ENTITIES(", Line, ")] Expected AcDbMText got ", check)
    }

    mText.Coordinates   = r.ConsumeCoordinates3D()
    mText.TextHeight    = ParseFloat(r.ConsumeDxfLine().Line)
    _                   = ParseFloat(r.ConsumeDxfLine().Line)

    // TODO: https://ezdxf.readthedocs.io/en/stable/dxfinternals/entities/mtext.html
    _                   = ParseFloat(r.ConsumeDxfLine().Line)

    mText.Layout        = uint8(r.ConsumeNumber(71, 10, "attachment point"))
    mText.Direction     = uint8(r.ConsumeNumber(72, 10, "direction (ex: left to right)"))

    for code := r.PeekCode(); code == 1 || code == 3; code = r.PeekCode() {
        mText.Text          = append(mText.Text, r.ConsumeDxfLine().Line)
    }

    mText.TextStyle     = r.ConsumeDxfLine().Line
    mText.Vector        = r.ConsumeCoordinates3D()
    mText.LineSpacing   = uint8(r.ConsumeNumber(73, 10, "line spacing"))
    // [44] LineSpacingFactor
    _                   = r.ConsumeDxfLine()

    // 691748


    dxf.MTexts = append(dxf.MTexts, mText)

}
