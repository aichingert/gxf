package parser

import (
    "log"

    "github.com/aichingert/dxf/pkg/entity"
    "github.com/aichingert/dxf/pkg/drawing"
)

func ParseEntities(r *Reader, dxf *drawing.Dxf) error {
    for {
        line, err := r.ConsumeDxfLine()
        if err != nil { return err }

        switch line.Line {
        case "LINE":
            if err = ParseLine(r, dxf);     err != nil { return err }
        case "LWPOLYLINE":
            if err = ParsePolyline(r, dxf); err != nil { return err }
        case "ARC":
            if err = ParseArc(r, dxf);      err != nil { return err }
        case "CIRCLE":
            if err = ParseCircle(r, dxf);   err != nil { return err }
        case "TEXT":
            if err = ParseText(r, dxf);     err != nil { return err }
        case "MTEXT":
            if err = ParseMText(r, dxf);    err != nil { return err }
        case "HATCH":
            if err = ParseHatch(r, dxf);    err != nil { return err }
        case "ELLIPSE":
            if err = ParseEllipse(r, dxf);  err != nil { return err }
        case "POINT":
            if err = ParsePoint(r, dxf);    err != nil { return err }
        case "INSERT":
            if err = ParseInsert(r, dxf);   err != nil { return err }
        default:
            log.Println("[ENTITIES] ", Line, ": ", line)
            return NewParseError("unknown entity")
        }
    }
}

func parseAcDbEntityE(r *Reader, entity entity.Entity) error {
    _, err := r.ConsumeDxfLine()
    if err != nil { return err }
    optional, err := r.ConsumeDxfLine()
    if err != nil { return err }

    // TODO: think about paper space visibility
    if optional.Code != 67 {
        entity.SetLayerName(optional.Line)
        return nil
    }

    // TODO: could lead to bug with start and end layername - seems like it is always the same
    layerName, err := r.ConsumeDxfLine()
    if err != nil { return err }
    entity.SetLayerName(layerName.Line)
    return nil
}

func extractHandleAndOwner(r *Reader) (*[2]uint64, error) {
    handle, err := r.ConsumeNumber(5, 16, "handle")
    if err != nil { return nil, err }

    // TODO: set hard owner/handle to owner dictionary
    code, err := r.PeekCode()
    if err != nil { return nil, err }
    if code == 102 {
        _, err = r.ConsumeDxfLine()
        if err != nil { return nil, err }
        _, err = r.ConsumeDxfLine()
        if err != nil { return nil, err }
        _, err = r.ConsumeDxfLine()
        if err != nil { return nil, err }
    }

    owner, err := r.ConsumeNumber(330, 16, "owner ptr")
    if err != nil { return nil, err }
    return &[2]uint64{handle, owner}, nil
}

func ParseLine(r *Reader, dxf *drawing.Dxf) error {
    result, err := extractHandleAndOwner(r)
    if err != nil { return err }
    line := entity.NewLine(result[0], result[1])

    err = parseAcDbEntityE(r, line)
    if err != nil { return err }

    check, err := r.ConsumeDxfLine()
    if err != nil { return err }

    if check.Line != "AcDbLine" {
        log.Println("[ENTITIES(", Line, ")] Expected AcDbLine got ", check)
        return NewParseError("expected AcDbLine")
    }

    src, err := r.ConsumeCoordinates3D()
    if err != nil { return err }
    dst, err := r.ConsumeCoordinates3D()
    if err != nil { return err }

    line.Src = src
    line.Dst = dst

    dxf.Lines = append(dxf.Lines, line)
    return nil
}

func ParsePolyline(r *Reader, dxf *drawing.Dxf) error {
    result, err := extractHandleAndOwner(r)
    if err != nil { return err }
    polyline := entity.NewPolyline(result[0], result[1])

    err = parseAcDbEntityE(r, polyline)
    if err != nil { return err }
    check, err := r.ConsumeDxfLine()
    if err != nil { return err }

    if check.Line != "AcDbPolyline" {
        log.Println("[ENTITIES(", Line, ")] Expected AcDbPolyline got ", check)
        return NewParseError("expected AcDbPolyline")
    }

    vertices, err := r.ConsumeNumber(90, 10, "number of vertices")
    if err != nil { return err }
    flag, err := r.ConsumeNumber(70, 10, "polyline flag")
    if err != nil { return err }

    polyline.Vertices = vertices
    polyline.Flag = flag

    // expecting code 43 
    line, err := r.ConsumeDxfLine()
    if err != nil { return err }
    if line.Code != 43 {
        log.Fatal("[ENTITIES] TODO: implement line width for each vertex")
    }

    for i := uint64(0); i < polyline.Vertices; i++ {
        // TODO: sometimes there is a bulge value for a vertex
        coords, err := r.ConsumeCoordinates2D()
        if err != nil { return err }
        bulge := 0.0

        // bulge = groupcode 42
        code, err := r.PeekCode()
        if err != nil { return err }

        if code == 42 {
            bulge, err = r.ConsumeFloat(42, "expected bulge")
            if err != nil { return err }
        }

        polyline.PolylineAppendCoordinate(coords, bulge)
    }

    dxf.Polylines = append(dxf.Polylines, polyline)
    return nil
}

func ParseArc(r *Reader, dxf *drawing.Dxf) error {
    result, err := extractHandleAndOwner(r)
    if err != nil { return err }
    arc     := entity.NewArc(result[0], result[1])

    err = parseAcDbEntityE(r, arc)
    if err != nil { return err }
    check, err := r.ConsumeDxfLine()
    if err != nil { return err }

    if check.Line != "AcDbCircle" {
        log.Println("[ENTITIES(", Line, ")] Expected AcDbCircle got ", check)
        return NewParseError("expected AcDbCircle ")
    }

    coords, err := r.ConsumeCoordinates3D()
    if err != nil { return err }
    radius, err := r.ConsumeFloat(40, "expected radius")
    if err != nil { return err }

    arc.Circle = &entity.Circle {
        Coordinates:    coords,
        Radius:         radius,
    }

    check, err = r.ConsumeDxfLine()
    if err != nil { return err }

    if check.Line != "AcDbArc" {
        log.Println("[ENTITIES(", Line, ")] Expected AcDbArc got ", check)
        return NewParseError("expected AcDbArc")
    }

    startAngle, err := r.ConsumeFloat(50, "expected startAngle")
    if err != nil { return err }
    endAngle, err := r.ConsumeFloat(51, "expected endAngle")
    if err != nil { return err }

    arc.StartAngle  = startAngle
    arc.EndAngle    = endAngle

    dxf.Arcs = append(dxf.Arcs, arc)
    return nil
}

func ParseCircle(r *Reader, dxf *drawing.Dxf) error {
    result, err := extractHandleAndOwner(r)
    if err != nil { return err }
    circle := entity.NewCircle(result[0], result[1])

    err = parseAcDbEntityE(r, circle)
    if err != nil { return err }
    check, err := r.ConsumeDxfLine()
    if err != nil { return err }

    if check.Line != "AcDbCircle" {
        log.Println("[ENTITIES(", Line, ")] Expected AcDbCircle got ", check)
        return NewParseError("expected AcDbCircle")
    }

    coords, err := r.ConsumeCoordinates3D()
    if err != nil { return err }
    radius, err := r.ConsumeFloat(40, "expected radius")
    if err != nil { return err }

    circle.Coordinates = coords
    circle.Radius      = radius

    dxf.Circles = append(dxf.Circles, circle)
    return nil
}

func ParseText(r *Reader, dxf *drawing.Dxf) error {
    result, err := extractHandleAndOwner(r)
    if err != nil { return err }
    text    := entity.NewMText(result[0], result[1])

    err = parseAcDbEntityE(r, text)
    if err != nil { return err }
    check, err := r.ConsumeDxfLine()
    if err != nil { return err }

    if check.Line != "AcDbText" {
        log.Println("[ENTITIES(", Line, ")] Expected AcDbText got ", check)
        return NewParseError("expected AcDbText")
    }

    code, err := r.PeekCode()
    if err != nil { return err }
    if code == 39 {
        // thickness
        _, err = r.ConsumeFloat(39, "expected thickness")
        if err != nil { return err }
    }

    // first alignment point
    _, err = r.ConsumeCoordinates3D()
    if err != nil { return err }

    _, err = r.ConsumeFloat(40, "expected text height")
    if err != nil { return err }

    _, err = r.ConsumeDxfLine() // [1] default value the string itself
    if err != nil { return err }

    code, err = r.PeekCode()
    if err != nil { return err }

    if code == 50 {
        _, err = r.ConsumeDxfLine() // text rotation default 0
        if err != nil { return err }
    }

    code, err = r.PeekCode()
    if err != nil { return err }
    if code == 41 {
        _, err = r.ConsumeDxfLine() // relative x scale factor default 1
        if err != nil { return err }
    }

    code, err = r.PeekCode()
    if err != nil { return err }
    if code == 51 {
        _, err = r.ConsumeDxfLine() // oblique angle default 0
        if err != nil { return err }
    }

    code, err = r.PeekCode()
    if err != nil { return err }
    if code == 7 {
        _, err = r.ConsumeDxfLine() // text style name default STANDARD
        if err != nil { return err }
    }

    code, err = r.PeekCode()
    if err != nil { return err }
    if code == 71 {
        _, err = r.ConsumeDxfLine() // text generation flags default 0
        if err != nil { return err }
    }

    code, err = r.PeekCode()
    if err != nil { return err }
    if code == 72 {
        _, err = r.ConsumeDxfLine() // horizontal text justification default 0
        if err != nil { return err }
    }

    code, err = r.PeekCode()
    if err != nil { return err }
    if code == 11 {
        // second alignment point
        _, err = r.ConsumeCoordinates3D()
        if err != nil { return err }
    }

    // optional default = 0, 0, 1
    code, err = r.PeekCode()
    if err != nil { return err }

    if code == 210 {
        // XYZ extrusion direction
        _, err = r.ConsumeCoordinates3D()
    }

    check, err = r.ConsumeDxfLine()
    if err != nil { return err }

    if check.Line != "AcDbText" {
        log.Println("[ENTITIES(", Line, ")] Expected AcDbText got ", check)
        return NewParseError("expected AcDbText")
    }

    // Group 72 and 73 integer codes 
    // https://help.autodesk.com/view/OARX/2024/ENU/?guid=GUID-62E5383D-8A14-47B4-BFC4-35824CAE8363

    code, err = r.PeekCode()
    if err != nil { return err }
    if code == 73 {
        _, err = r.ConsumeDxfLine() // Vertical text justification type default 0
        if err != nil { return err }
    }

    _ = dxf
    return nil
}

func ParseMText(r *Reader, dxf *drawing.Dxf) error {
    result, err := extractHandleAndOwner(r)
    if err != nil { return err }
    mText   := entity.NewMText(result[0], result[1])

    err = parseAcDbEntityE(r, mText)
    if err != nil { return err }
    check, err := r.ConsumeDxfLine()
    if err != nil { return err }

    if check.Line != "AcDbMText" {
        log.Println("[ENTITIES(", Line, ")] Expected AcDbMText got ", check)
        return NewParseError("expected AcDbMText")
    }

    coords, err := r.ConsumeCoordinates3D()
    if err != nil { return err }
    textHeight, err := r.ConsumeFloat(40, "expected text height")
    if err != nil { return err }

    // TODO: https://ezdxf.readthedocs.io/en/stable/dxfinternals/entities/mtext.html
    _, err = r.ConsumeFloat(41, "rectangle width")
    if err != nil { return err }

    _, err = r.ConsumeFloat(46, "column height")
    if err != nil { return err }

    layout, err := r.ConsumeNumber(71, 10, "attachment point")
    if err != nil { return err }
    direction, err := r.ConsumeNumber(72, 10, "direction (ex: left to right)")
    if err != nil { return err }

    mText.Coordinates   = coords
    mText.TextHeight    = textHeight

    mText.Layout        = uint8(layout)
    mText.Direction     = uint8(direction)

    code, err := r.PeekCode()
    if err != nil { return err }

    for code == 1 || code == 3 {
        line, err := r.ConsumeDxfLine()
        if err != nil { return err }
        mText.Text          = append(mText.Text, line.Line)

        code, err = r.PeekCode()
        if err != nil { return err }
    }

    line, err := r.ConsumeDxfLine()
    if err != nil { return err }

    vector, err := r.ConsumeCoordinates3D()
    if err != nil { return err }
    spacing, err := r.ConsumeNumber(73, 10, "line spacing")

    mText.TextStyle     = line.Line
    mText.Vector        = vector
    mText.LineSpacing   = uint8(spacing)

    // [44] LineSpacingFactor
    _, err              = r.ConsumeDxfLine()
    if err != nil { return err }

    dxf.MTexts = append(dxf.MTexts, mText)
    return nil
}

func ParseHatch(r *Reader, dxf *drawing.Dxf) error {
    result, err := extractHandleAndOwner(r)
    if err != nil { return err }
    hatch   := entity.NewMText(result[0], result[1])

    err = parseAcDbEntityE(r, hatch)
    if err != nil { return err }
    check, err := r.ConsumeDxfLine()
    if err != nil { return err }

    if check.Line != "AcDbHatch" {
        log.Println("[ENTITIES(", Line, ")] Expected AcDbHatch got ", check)
        return NewParseError("expected AcDbHatch")
    }

    // 10,20,30
    _, err = r.ConsumeCoordinates3D()
    if err != nil { return err }

    // TODO: [210/220/230] Extrustion direction (only need 2D maybe later)
    _, err = r.ConsumeDxfLine()
    if err != nil { return err }
    _, err = r.ConsumeDxfLine()
    if err != nil { return err }
    _, err = r.ConsumeDxfLine()
    if err != nil { return err }

    patternName, err := r.ConsumeDxfLine()
    if err != nil { return err }
    solidFillFlag, err := r.ConsumeNumber(70, 10, "solid fill flag")
    if err != nil { return err }
    associativityFlag, err := r.ConsumeNumber(71, 10, "associativity flag")
    if err != nil { return err }

    // Number of boundary paths?
    _, err = r.ConsumeNumber(91, 10, "boundary paths")
    if err != nil { return err }
    pt, err := r.ConsumeNumber(92, 10, "boundary path type")
    if err != nil { return err }

    if pt & 2 > 0 {
        b, err := r.ConsumeNumber(72, 10, "has bulge flag")
        if err != nil { return err }
        _, err = r.ConsumeNumber(73, 10, "is closed flag")
        if err != nil { return err }
        n, err := r.ConsumeNumber(93, 10, "number fo polyline vertices")
        if err != nil { return err }

        for i := uint64(0); i < n; i++ {
            _, err = r.ConsumeCoordinates2D()
            if err != nil { return err }
            code, err := r.PeekCode()
            if err != nil { return err }
            if code == 42 {
                _ ,err = r.ConsumeFloat(42, "expected bulge")
                if err != nil { return err }
            }
        }

        _ = b
    } else {
        n, err := r.ConsumeNumber(93, 10, "number of edges in this boundary path")
        if err != nil { return err }
        t, err := r.ConsumeNumber(72, 10, "edge type data")
        if err != nil { return err }

        switch t {
        case 1:
            // Parse Line
            for i := uint64(0); i < n; i++ {
                _, err = r.ConsumeCoordinates2D()
                if err != nil { return err }
            }
        case 2:
            // Circular arc
            for i := uint64(0); i < n; i++ {
                _,err = r.ConsumeCoordinates2D() 
                if err != nil { return err }
            }

            _, err = r.ConsumeNumber(40, 10, "radius")
            if err != nil { return err }
            _, err = r.ConsumeNumber(50, 10, "start angle")
            if err != nil { return err }
            _, err = r.ConsumeNumber(51, 10, "end angle")
            if err != nil { return err }
            _, err = r.ConsumeNumber(73, 10, "is counterclockwise flag")
            if err != nil { return err }
        case 3:
            // Elliptic arc
            for i := uint64(0); i < n; i++ { 
                _, err = r.ConsumeCoordinates2D() 
                if err != nil { return err }
            }

            _, err = r.ConsumeNumber(40, 10, "length of minor axis")
            if err != nil { return err }
            _, err = r.ConsumeNumber(50, 10, "start angle")
            if err != nil { return err }
            _, err = r.ConsumeNumber(51, 10, "end angle")
            if err != nil { return err }
            _, err = r.ConsumeNumber(73, 10, "is counterclockwise flag")
            if err != nil { return err }
        case 4:
            // Spine
            _, err = r.ConsumeNumber(94, 10, "degree")
            if err != nil { return err }
            _, err = r.ConsumeNumber(73, 10, "rational")
            if err != nil { return err }
            _, err = r.ConsumeNumber(74, 10, "periodic")
            if err != nil { return err }
            k, err := r.ConsumeNumber(95, 10, "number of knots")
            if err != nil { return err }
            _, err = r.ConsumeNumber(96, 10, "number of control points")
            if err != nil { return err }

            for i := uint64(0); i < k; i++ {
                _, err = r.ConsumeNumber(40, 10, "knot values")
                if err != nil { return err }
                _, err = r.ConsumeCoordinates2D()
                if err != nil { return err }
            }

            code, err := r.PeekCode()
            if err != nil { return err }
            if code == 42 {
                _, err = r.ConsumeNumber(42, 10, "weights") // optional 1
                if err != nil { return err }
            }

            _, err = r.ConsumeNumber(97, 10, "number of fit data")
            if err != nil { return err }
            _, err = r.ConsumeNumber(11, 10, "X fit datum value")
            if err != nil { return err }
            _, err = r.ConsumeNumber(21, 10, "Y fit datum value")
            if err != nil { return err }
            _, err = r.ConsumeNumber(12, 10, "X start tangent")
            if err != nil { return err }
            _, err = r.ConsumeNumber(22, 10, "Y start tangent")
            if err != nil { return err }
            _, err = r.ConsumeNumber(13, 10, "X end tangent")
            if err != nil { return err }
            _, err = r.ConsumeNumber(23, 10, "Y end tangent")
            if err != nil { return err }
        default:
            log.Println("[ENTITIES(HATCH - ", Line, " )] invalid edge type data: ", t)
            return NewParseError("invalid edge type data")
        }
    }

    bo, err := r.ConsumeNumber(97, 10, "number of source boundary objects")
    if err != nil { return err }

    for i := uint64(0); i < bo; i++ {
        _, err = r.ConsumeNumber(330, 10, "reference to source boundary objects")
        if err != nil { return err }
    }

    _, err = r.ConsumeNumber(75, 10, "hatch style")
    if err != nil { return err }
    _, err = r.ConsumeNumber(76, 10, "hatch pattern type")
    if err != nil { return err }

    sp, err := r.ConsumeNumber(98, 10, "number of seed points")
    if err != nil { return err }

    if sp != 0 {
        log.Fatal("TODO(", Line, "): hatch implement seed points")
    }

    _ = patternName
    _ = solidFillFlag
    _ = associativityFlag
    _ = dxf
    return nil
}

func ParseEllipse(r *Reader, dxf *drawing.Dxf) error {
    result, err := extractHandleAndOwner(r)
    if err != nil { return err }
    ellipse := entity.NewMText(result[0], result[1]) // todo ellipse

    err = parseAcDbEntityE(r, ellipse)
    if err != nil { return err }
    check, err := r.ConsumeDxfLine()
    if err != nil { return err }

    if check.Line != "AcDbEllipse" {
        log.Println("[ENTITIES(", Line, ")] Expected AcDbEllipse got ", check)
        return NewParseError("expected AcDbEllipse")
    }

    _, err = r.ConsumeCoordinates3D() // Center point
    if err != nil { return err }
    _, err = r.ConsumeCoordinates3D() // endpoint of major axis
    if err != nil { return err }

    // optional default = 0, 0, 1
    code, err := r.PeekCode()
    if err != nil { return err }
    if code == 210 {
        // XYZ extrusion direction
        _, err = r.ConsumeCoordinates3D()
        if err != nil { return err }
    }

    _, err = r.ConsumeFloat(40, "ratio of minor axis to major axis")
    if err != nil { return err }

    _, err = r.ConsumeFloat(41, "start parameter")
    if err != nil { return err }

    _, err = r.ConsumeFloat(42, "end parameter")
    if err != nil { return err }

    _ = dxf
    return nil
}

func ParsePoint(r *Reader, dxf *drawing.Dxf) error {
    result, err := extractHandleAndOwner(r)
    if err != nil { return err }
    point   := entity.NewMText(result[0], result[1]) // TODO: point

    err = parseAcDbEntityE(r, point)
    if err != nil { return err }
    check, err := r.ConsumeDxfLine()
    if err != nil { return err }

    if check.Line != "AcDbPoint" {
        log.Println("[ENTITIES(", Line, ")] Expected AcDbPoint got ", check)
        return NewParseError("expected AcDbPoint")
    }

    _, err = r.ConsumeCoordinates3D() // Point location
    if err != nil { return err }
    code, err := r.PeekCode()
    if err != nil { return err }
    if code == 39 {
        _, err = r.ConsumeNumber(39, 10, "thickness")
        if err != nil { return err }
    }

    // optional default = 0, 0, 1
    code, err = r.PeekCode()
    if err != nil { return err }
    if code == 210 {
        // XYZ extrusion direction
        _, err = r.ConsumeCoordinates3D()
        if err != nil { return err }
    }

    code, err = r.PeekCode()
    if err != nil { return err }
    if code == 50 {
        _, err = r.ConsumeFloat(50, "angle of the x axis")
        if err != nil { return err }
    }
   
    _ = dxf
    return nil
}

// TODO: have to implement block section first
func ParseInsert(r *Reader, dxf *drawing.Dxf) error {
    result, err := extractHandleAndOwner(r)
    if err != nil { return err }
    insert  := entity.NewMText(result[0], result[1]) // TODO: insert

    err = parseAcDbEntityE(r, insert)
    if err != nil { return err }
    check, err := r.ConsumeDxfLine()
    if err != nil { return err }

    if check.Line != "AcDbBlockReference" {
        log.Println("[ENTITIES(", Line, ")] Expected AcDbBlockReference got ", check)
        return NewParseError("expected AcDbBlockReference")
    }

    code, err := r.PeekCode()
    if err != nil { return err }
    if code == 66 {
        _, err = r.ConsumeDxfLine() // Variable attributes-follow flag default = 0
        if err != nil { return err }
    }

    _, err = r.ConsumeDxfLine() // Block name
    if err != nil { return err }
    _, err = r.ConsumeCoordinates3D() // insertion point
    if err != nil { return err }

    code, err = r.PeekCode()
    if err != nil { return err }
    if code == 41 {
        _, err = r.ConsumeDxfLine() // xyz scale factors default 1
        if err != nil { return err }
        _, err = r.ConsumeDxfLine() // xyz scale factors default 1
        if err != nil { return err }
        _, err = r.ConsumeDxfLine() // xyz scale factors default 1
        if err != nil { return err }
    }

    code, err = r.PeekCode()
    if err != nil { return err }
    if code == 50 {
        _, err = r.ConsumeDxfLine() // rotation angle default = 0
        if err != nil { return err }
    }

    code, err = r.PeekCode()
    if err != nil { return err }
    if code == 70 {
        _, err = r.ConsumeDxfLine() // column count default = 1
        if err != nil { return err }
    }
    code, err = r.PeekCode()
    if err != nil { return err }
    if code == 71 {
        _, err = r.ConsumeDxfLine() // row count default = 1
        if err != nil { return err }
    }

    code, err = r.PeekCode()
    if err != nil { return err }
    if code == 44 {
        _, err = r.ConsumeDxfLine() // column spacing default = 0
        if err != nil { return err }
    }
    code, err = r.PeekCode()
    if err != nil { return err }
    if code == 45 {
        _, err = r.ConsumeDxfLine() // row spacing default = 0
        if err != nil { return err }
    }

    // optional default = 0, 0, 1
    code, err = r.PeekCode()
    if err != nil { return err }
    if code == 210 {
        // XYZ extrusion direction
        _, err = r.ConsumeCoordinates3D()
        if err != nil { return err }
    }

    // TODO: parse insert
    // attrib =>  987212

    _ = dxf
    return nil
}
