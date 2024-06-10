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
    return [2]uint64{
        r.ConsumeNumber(5, 16, "handle"),
        r.ConsumeNumber(330, 16, "owner ptr"),
    }
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
