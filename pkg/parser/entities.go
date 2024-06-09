package parser

import (
    "log"
    "bufio"

    "github.com/aichingert/dxf/pkg/entity"
    "github.com/aichingert/dxf/pkg/drawing"
)

func ParseEntities(sc *bufio.Scanner, dxf *drawing.Dxf) {
    for {
        switch variable := ExtractCodeAndValue(sc); variable[1] {
        case "LINE":
            ParseLine(sc, dxf)
        case "LWPOLYLINE":
            ParsePolyline(sc, dxf)
        default:
            log.Fatal("[ENTITIES] ", Line, ": ", variable)
        }
    }
    _ = dxf
}

func parseAcDbEntityE(sc *bufio.Scanner, entity entity.Entity) {
    _ = ExtractCodeAndValue(sc)
    optional := ExtractCodeAndValue(sc)

    // TODO: think about paper space visibility
    if optional[0] != " 67" {
        entity.SetLayerName(optional[1])
        return
    }

    // TODO: could lead to bug with start and end layername - seems like it is always the same
    layerName := ExtractCodeAndValue(sc)
    entity.SetLayerName(layerName[1])
}

func extractHandleAndOwner(sc *bufio.Scanner) [2]uint64 {
    return [2]uint64{
        ExtractHex(sc, "5", "handle"),
        ExtractHex(sc, "330", "owner ptr"),
    }
}

func ParseLine(sc *bufio.Scanner, dxf *drawing.Dxf) {
    result := extractHandleAndOwner(sc)
    line   := entity.NewLine(result[0], result[1])

    parseAcDbEntityE(sc, line)

    check := ExtractCodeAndValue(sc)

    if check[1] != "AcDbLine" {
        log.Fatal("[ENTITIES(", Line, ")] Expected AcDbLine got ", check)
    }

    line.Src = ExtractCoordinates3D(sc)
    line.Dst = ExtractCoordinates3D(sc)

    dxf.Lines = append(dxf.Lines, line)
}

func ParsePolyline(sc *bufio.Scanner, dxf *drawing.Dxf) {
    result   := extractHandleAndOwner(sc)
    polyline := entity.NewPolyline(result[0], result[1])

    parseAcDbEntityE(sc, polyline)
    check := ExtractCodeAndValue(sc)

    if check[1] != "AcDbPolyline" {
        log.Fatal("[ENTITIES(", Line, ")] Expected AcDbPolyline got ", check)
    }

    polyline.Vertices = ExtractHex(sc, "90", "number of vertices")
    polyline.Flag = ExtractHex(sc, "70", "polyline flag")

    // expecting code 43
    if ExtractCodeAndValue(sc)[0] != " 43" {
        log.Fatal("[ENTITIES] TODO: implement line width for each vertex")
    }

    for i := uint64(0); i < polyline.Vertices; i++ {
        line := ExtractCoordinates2D(sc)
        polyline.Coordinates = append(polyline.Coordinates, line)

        // TODO: sometimes there is a bulge value for a vertex
    }

    dxf.Polylines = append(dxf.Polylines, polyline)
}
