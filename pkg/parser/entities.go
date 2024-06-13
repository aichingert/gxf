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
        case "TEXT":
            ParseText(r, dxf)
        case "MTEXT":
            ParseMText(r, dxf)
        case "HATCH":
            ParseHatch(r, dxf)
        case "ELLIPSE":
            ParseEllipse(r, dxf)
        case "POINT":
            ParsePoint(r, dxf)
        case "INSERT":
            ParseInsert(r, dxf)
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

func ParseText(r *Reader, dxf *drawing.Dxf) {
    result  := extractHandleAndOwner(r)
    text    := entity.NewMText(result[0], result[1])

    parseAcDbEntityE(r, text)
    check := r.ConsumeDxfLine()

    if check.Line != "AcDbText" {
        log.Fatal("[ENTITIES(", Line, ")] Expected AcDbText got ", check)
    }

    if r.PeekCode() == 39 {
        // thickness
        _ = ParseFloat(r.ConsumeDxfLine().Line)
    }

    // first alignment point
    _ = r.ConsumeCoordinates3D()
    _ = ParseFloat(r.ConsumeDxfLine().Line) // [40] text height

    _ = r.ConsumeDxfLine() // [1] default value the string itself

    if r.PeekCode() == 50 {
        _ = r.ConsumeDxfLine() // text rotation default 0
    }
    if r.PeekCode() == 41 {
        _ = r.ConsumeDxfLine() // relative x scale factor default 1
    }

    if r.PeekCode() == 51 {
        _ = r.ConsumeDxfLine() // oblique angle default 0
    }

    if r.PeekCode() == 7 {
        _ = r.ConsumeDxfLine() // text style name default STANDARD
    }

    if r.PeekCode() == 71 {
        _ = r.ConsumeDxfLine() // text generation flags default 0
    }

    if r.PeekCode() == 72 {
        _ = r.ConsumeDxfLine() // horizontal text justification default 0
    }
    
    if r.PeekCode() == 11 {
        // second alignment point
        _ = r.ConsumeCoordinates3D()
    }

    // optional default = 0, 0, 1
    if r.PeekCode() == 210 {
        // XYZ extrusion direction
        _ = ParseFloat(r.ConsumeDxfLine().Line)
        _ = ParseFloat(r.ConsumeDxfLine().Line)
        _ = ParseFloat(r.ConsumeDxfLine().Line)
    }

    check = r.ConsumeDxfLine()

    if check.Line != "AcDbText" {
        log.Fatal("[ENTITIES(", Line, ")] Expected AcDbText got ", check)
    }

    // Group 72 and 73 integer codes 
    // https://help.autodesk.com/view/OARX/2024/ENU/?guid=GUID-62E5383D-8A14-47B4-BFC4-35824CAE8363

    if r.PeekCode() == 73 {
        _ = r.ConsumeDxfLine() // Vertical text justification type default 0
    }

    _ = dxf
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

    dxf.MTexts = append(dxf.MTexts, mText)
}

func ParseHatch(r *Reader, dxf *drawing.Dxf) {
    result  := extractHandleAndOwner(r)
    hatch   := entity.NewMText(result[0], result[1])

    parseAcDbEntityE(r, hatch)
    check := r.ConsumeDxfLine()

    if check.Line != "AcDbHatch" {
        log.Fatal("[ENTITIES(", Line, ")] Expected AcDbHatch got ", check)
    }

    // 10,20,30
    _ = r.ConsumeCoordinates3D()

    // TODO: [210/220/230] Extrustion direction (only need 2D maybe later)
    _ = r.ConsumeDxfLine()
    _ = r.ConsumeDxfLine()
    _ = r.ConsumeDxfLine()

    patternName := r.ConsumeDxfLine()
    solidFillFlag := r.ConsumeNumber(70, 10, "solid fill flag")
    associativityFlag := r.ConsumeNumber(71, 10, "associativity flag")

    // Number of boundary paths?
    _ = r.ConsumeNumber(91, 10, "boundary paths")
    pt := r.ConsumeNumber(92, 10, "boundary path type")

    if pt & 2 > 0 {
        b := r.ConsumeNumber(72, 10, "has bulge flag")
        _ = r.ConsumeNumber(73, 10, "is closed flag")
        n := r.ConsumeNumber(93, 10, "number fo polyline vertices")

        for i := uint64(0); i < n; i++ {
            _ = r.ConsumeCoordinates2D()
            if r.PeekCode() == 42 {
                _ = ParseFloat(r.ConsumeDxfLine().Line) // bulge
            }
        }

        _ = b
    } else {
        n := r.ConsumeNumber(93, 10, "number of edges in this boundary path")
        t := r.ConsumeNumber(72, 10, "edge type data")

        switch t {
        case 1:
            // Parse Line
            for i := uint64(0); i < n; i++ {
                _ = r.ConsumeCoordinates2D()
            }
        case 2:
            // Circular arc
            for i := uint64(0); i < n; i++ {
                _ = r.ConsumeCoordinates2D() 
            }

            _ = r.ConsumeNumber(40, 10, "radius")
            _ = r.ConsumeNumber(50, 10, "start angle")
            _ = r.ConsumeNumber(51, 10, "end angle")
            _ = r.ConsumeNumber(73, 10, "is counterclockwise flag")
        case 3:
            // Elliptic arc
            for i := uint64(0); i < n; i++ { _ = r.ConsumeCoordinates2D() }

            _ = r.ConsumeNumber(40, 10, "length of minor axis")
            _ = r.ConsumeNumber(50, 10, "start angle")
            _ = r.ConsumeNumber(51, 10, "end angle")
            _ = r.ConsumeNumber(73, 10, "is counterclockwise flag")
        case 4:
            // Spine
            _ = r.ConsumeNumber(94, 10, "degree")
            _ = r.ConsumeNumber(73, 10, "rational")
            _ = r.ConsumeNumber(74, 10, "periodic")
            k := r.ConsumeNumber(95, 10, "number of knots")
            _ = r.ConsumeNumber(96, 10, "number of control points")

            for i := uint64(0); i < k; i++ {
                _ = r.ConsumeNumber(40, 10, "knot values")
                _ = r.ConsumeCoordinates2D()
            }

            if r.PeekCode() == 42 {
                _ = r.ConsumeNumber(42, 10, "weights") // optional 1
            }

            _ = r.ConsumeNumber(97, 10, "number of fit data")
            _ = r.ConsumeNumber(11, 10, "X fit datum value")
            _ = r.ConsumeNumber(21, 10, "Y fit datum value")
            _ = r.ConsumeNumber(12, 10, "X start tangent")
            _ = r.ConsumeNumber(22, 10, "Y start tangent")
            _ = r.ConsumeNumber(13, 10, "X end tangent")
            _ = r.ConsumeNumber(23, 10, "Y end tangent")
        default:
            log.Fatal("[ENTITIES(HATCH - ", Line, " )] invalid edge type data: ", t)
        }
    }

    bo := r.ConsumeNumber(97, 10, "number of source boundary objects")

    for i := uint64(0); i < bo; i++ {
        _ = r.ConsumeNumber(330, 10, "reference to source boundary objects")
    }

    _ = r.ConsumeNumber(75, 10, "hatch style")
    _ = r.ConsumeNumber(76, 10, "hatch pattern type")

    sp := r.ConsumeNumber(98, 10, "number of seed points")

    if sp != 0 {
        log.Fatal("TODO(", Line, "): hatch implement seed points")
    }

    _ = patternName
    _ = solidFillFlag
    _ = associativityFlag
    _ = dxf
}

func ParseEllipse(r *Reader, dxf *drawing.Dxf) {
    result  := extractHandleAndOwner(r)
    ellipse := entity.NewMText(result[0], result[1]) // todo ellipse

    parseAcDbEntityE(r, ellipse)
    check := r.ConsumeDxfLine()

    if check.Line != "AcDbEllipse" {
        log.Fatal("[ENTITIES(", Line, ")] Expected AcDbEllipse got ", check)
    }

    _ = r.ConsumeCoordinates3D() // Center point
    _ = r.ConsumeCoordinates3D() // endpoint of major axis

    // optional default = 0, 0, 1
    if r.PeekCode() == 210 {
        // XYZ extrusion direction
        _ = ParseFloat(r.ConsumeDxfLine().Line)
        _ = ParseFloat(r.ConsumeDxfLine().Line)
        _ = ParseFloat(r.ConsumeDxfLine().Line)
    }

    _ = ParseFloat(r.ConsumeDxfLine().Line) // 40, 10, "ratio of minor axis to major axis"
    _ = ParseFloat(r.ConsumeDxfLine().Line) // 41, 10, "start parameter"
    _ = ParseFloat(r.ConsumeDxfLine().Line) // 42, 10, "end parameter"

    _ = dxf
}

func ParsePoint(r *Reader, dxf *drawing.Dxf) {
    result  := extractHandleAndOwner(r)
    point   := entity.NewMText(result[0], result[1]) // TODO: point

    parseAcDbEntityE(r, point)
    check := r.ConsumeDxfLine()

    if check.Line != "AcDbPoint" {
        log.Fatal("[ENTITIES(", Line, ")] Expected AcDbPoint got ", check)
    }

    _ = r.ConsumeCoordinates3D() // Point location
    if r.PeekCode() == 39 {
        _ = r.ConsumeNumber(39, 10, "thickness")
    }

    // optional default = 0, 0, 1
    if r.PeekCode() == 210 {
        // XYZ extrusion direction
        _ = ParseFloat(r.ConsumeDxfLine().Line)
        _ = ParseFloat(r.ConsumeDxfLine().Line)
        _ = ParseFloat(r.ConsumeDxfLine().Line)
    }

    if r.PeekCode() == 50 {
        _ = ParseFloat(r.ConsumeDxfLine().Line) // angle of the x axis
    }
   
    _ = dxf
}

// TODO: have to implement block section first
func ParseInsert(r *Reader, dxf *drawing.Dxf) {
    result  := extractHandleAndOwner(r)
    insert  := entity.NewMText(result[0], result[1]) // TODO: insert

    // 987176
    parseAcDbEntityE(r, insert)
    check := r.ConsumeDxfLine()

    if check.Line != "AcDbBlockReference" {
        log.Fatal("[ENTITIES(", Line, ")] Expected AcDbBlockReference got ", check)
    }

    if r.PeekCode() == 66 {
        _ = r.ConsumeDxfLine() // Variable attributes-follow flag default = 0
    }

    _ = r.ConsumeDxfLine() // Block name
    _ = r.ConsumeCoordinates3D() // insertion point

    if r.PeekCode() == 41 {
        _ = r.ConsumeDxfLine() // xyz scale factors default 1
        _ = r.ConsumeDxfLine() // xyz scale factors default 1
        _ = r.ConsumeDxfLine() // xyz scale factors default 1
    }

    if r.PeekCode() == 50 {
        _ = r.ConsumeDxfLine() // rotation angle default = 0
    }

    if r.PeekCode() == 70 {
        _ = r.ConsumeDxfLine() // column count default = 1
    }
    if r.PeekCode() == 71 {
        _ = r.ConsumeDxfLine() // row count default = 1
    }

    if r.PeekCode() == 44 {
        _ = r.ConsumeDxfLine() // column spacing default = 0
    }
    if r.PeekCode() == 45 {
        _ = r.ConsumeDxfLine() // row spacing default = 0
    }

    // optional default = 0, 0, 1
    if r.PeekCode() == 210 {
        // XYZ extrusion direction
        _ = ParseFloat(r.ConsumeDxfLine().Line)
        _ = ParseFloat(r.ConsumeDxfLine().Line)
        _ = ParseFloat(r.ConsumeDxfLine().Line)
    }

    _ = dxf
}
