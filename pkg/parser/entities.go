package parser

import (
    "log"

    "github.com/aichingert/dxf/pkg/entity"
    "github.com/aichingert/dxf/pkg/drawing"
)

func ParseEntities(r *Reader, dxf *drawing.Dxf) error {
    for r.ScanDxfLine() {
        switch r.DxfLine().Line {
        case "LINE":
            Wrap(ParseLine, r, dxf);
        case "LWPOLYLINE":
            Wrap(ParsePolyline, r, dxf);
        case "ARC":
            Wrap(ParseArc, r, dxf);
        case "CIRCLE":
            Wrap(ParseCircle, r, dxf);
        case "TEXT":
            Wrap(ParseText, r, dxf);
        case "MTEXT":
            Wrap(ParseMText, r, dxf);
        case "HATCH":
            Wrap(ParseHatch, r, dxf);
        case "ELLIPSE":
            Wrap(ParseEllipse, r, dxf);
        case "POINT":
            Wrap(ParsePoint, r, dxf);
        case "INSERT":
            Wrap(ParseInsert, r, dxf);
        default:
            log.Println("[ENTITIES] ", Line, ": ", r.DxfLine().Line)
            return NewParseError("unknown entity")
        }

        if WrappedErr != nil {
            return WrappedErr
        }
    }

    return r.Err()
}

func parseAcDbEntityE(r *Reader, entity entity.Entity) error {
    r.ConsumeStr(nil) // AcDbEntity
    if r.ConsumeStrIf(67, entity.GetLayerName()) {
        return r.Err()
    }

    // TODO: could lead to bug with start and end layername 
    // * seems like it is always the same
    r.ConsumeStr(entity.GetLayerName())
    return r.Err()
}

func extractHandleAndOwner(r *Reader, handle *uint64, owner *uint64) error {
    r.ConsumeNumber(5, HEX_RADIX, "handle", handle) 

    // TODO: set hard owner/handle to owner dictionary
    if r.ConsumeStrIf(102, nil) { // consumeIf => ex. {ACAD_XDICTIONARY
        r.ConsumeStr(nil) // 360 => hard owner
        r.ConsumeStr(nil) // 102 }
    }

    r.ConsumeNumber(330, HEX_RADIX, "owner ptr", owner)
    return r.Err()
}

// 656790

func ParseLine(r *Reader, dxf *drawing.Dxf) error {
    if extractHandleAndOwner(r, &Handle, &Owner) != nil {
        return r.Err()
    }

    line := entity.NewLine(Handle, Owner)

    if parseAcDbEntityE(r, line) != nil || r.AssertNextLine("AcDbLine") != nil {
        return r.Err()
    }

    r.ConsumeCoordinates(line.Src[:])
    r.ConsumeCoordinates(line.Dst[:])

    dxf.Lines = append(dxf.Lines, line)
    return r.Err()
}

func ParsePolyline(r *Reader, dxf *drawing.Dxf) error {
    if extractHandleAndOwner(r, &Handle, &Owner) != nil {
        return r.Err()
    }

    polyline := entity.NewPolyline(Handle, Owner)

    if parseAcDbEntityE(r, polyline) != nil || r.AssertNextLine("AcDbPolyline") != nil {
        return r.Err()
    }

    r.ConsumeNumber(90, DEC_RADIX, "number of vertices", &polyline.Vertices)
    r.ConsumeNumber(70, DEC_RADIX, "polyline flag", &polyline.Flag)

    // expecting code 43 
    line, err := r.ConsumeDxfLine()
    if err != nil { return err }
    if line.Code != 43 {
        log.Fatal("[ENTITIES] TODO: implement line width for each vertex")
    }

    for i := uint64(0); i < polyline.Vertices; i++ {
        // TODO: sometimes there is a bulge value for a vertex
        bulge := 0.0
        coords := [2]float64{0.0,0.0}

        r.ConsumeCoordinates(coords[:])
        r.ConsumeFloatIf(42, "expected, bulge", &bulge)

        if r.Err() != nil {
            return r.Err()
        }

        polyline.PolylineAppendCoordinate(coords, bulge)
    }

    dxf.Polylines = append(dxf.Polylines, polyline)
    return r.Err()
}

func ParseArc(r *Reader, dxf *drawing.Dxf) error {
    if extractHandleAndOwner(r, &Handle, &Owner) != nil {
        return r.Err()
    }

    arc := entity.NewArc(Handle, Owner)

    if parseAcDbEntityE(r, arc) != nil || r.AssertNextLine("AcDbCircle") != nil {
        return r.Err()
    }

    r.ConsumeCoordinates(arc.Circle.Coordinates[:])
    r.ConsumeFloat(40, "expected radius", &arc.Circle.Radius)

    if r.AssertNextLine("AcDbArc") != nil {
        return r.Err()
    }

    r.ConsumeFloat(50, "expected startAngle", &arc.StartAngle)
    r.ConsumeFloat(51, "expected endAngle", &arc.EndAngle)

    dxf.Arcs = append(dxf.Arcs, arc)
    return r.Err()
}

func ParseCircle(r *Reader, dxf *drawing.Dxf) error {
    if extractHandleAndOwner(r, &Handle, &Owner) != nil {
        return r.Err()
    }

    circle := entity.NewCircle(Handle, Owner)

    if parseAcDbEntityE(r, circle) != nil || r.AssertNextLine("AcDbCircle") != nil {
        return r.Err()
    }

    r.ConsumeCoordinates(circle.Coordinates[:])
    r.ConsumeFloat(40, "expected radius", &circle.Radius)

    dxf.Circles = append(dxf.Circles, circle)
    return r.Err()
}

func ParseText(r *Reader, dxf *drawing.Dxf) error {
    if extractHandleAndOwner(r, &Handle, &Owner) != nil {
        return r.Err()
    }
    
    text := entity.NewMText(Handle, Owner)

    if parseAcDbEntityE(r, text) != nil || r.AssertNextLine("AcDbText") != nil {
        return r.Err()
    }

    r.ConsumeFloatIf(39, "expected thickness", nil)

    // first alignment point
    coords := [3]float64{0.0,0.0,0.0}
    r.ConsumeCoordinates(coords[:])
    r.ConsumeCoordinates(coords[:])

    r.ConsumeFloat(40, "expected text height", nil)
    r.ConsumeStr(nil) // "[1] default value the string itself"

    r.ConsumeNumberIf(50, DEC_RADIX, "text rotation default 0", nil)

    r.ConsumeFloatIf(41, "relative x scale factor default 1", nil)
    r.ConsumeFloatIf(51, "oblique angle default 0", nil)

    r.ConsumeStrIf(7, nil) // "text style name default STANDARD"

    r.ConsumeNumberIf(71, DEC_RADIX, "text generation flags default 0", nil)
    r.ConsumeNumberIf(72, DEC_RADIX, "horizontal text justification default 0", nil)

    r.ConsumeCoordinatesIf(11, coords[:])
    // XYZ extrusion direction
    r.ConsumeCoordinatesIf(210, coords[:]) // optional default 0,0,1

    if r.AssertNextLine("AcDbText") != nil {
        return r.Err()
    }

    // Group 72 and 73 integer codes 
    // https://help.autodesk.com/view/OARX/2024/ENU/?guid=GUID-62E5383D-8A14-47B4-BFC4-35824CAE8363

    r.ConsumeNumberIf(73, DEC_RADIX, "vertical text justification type default 0", nil)

    _ = dxf
    _ = coords
    return r.Err()
}

func ParseMText(r *Reader, dxf *drawing.Dxf) error {
    if extractHandleAndOwner(r, &Handle, &Owner) != nil {
        return r.Err()
    }

    mText   := entity.NewMText(Handle, Owner)

    if parseAcDbEntityE(r, mText) != nil || r.AssertNextLine("AcDbMText") != nil {
        return r.Err()
    }

    r.ConsumeCoordinates(mText.Coordinates[:])
    r.ConsumeFloat(40, "expected text height", &mText.TextHeight)

    // TODO: https://ezdxf.readthedocs.io/en/stable/dxfinternals/entities/mtext.html
    r.ConsumeFloat(41, "rectangle width", nil)
    r.ConsumeFloat(46, "column height", nil)

    r.ConsumeNumber(71, DEC_RADIX, "attachment point", &mText.Layout)
    r.ConsumeNumber(72, DEC_RADIX, "direction (ex: left to right)", &mText.Direction)

    // TODO: implement more helper :smelting:
    code, err := r.PeekCode()
    if err != nil { return err }

    for code == 1 || code == 3 {
        line, err := r.ConsumeDxfLine()
        if err != nil { return err }

        mText.Text          = append(mText.Text, line.Line)

        code, err = r.PeekCode()
        if err != nil { return err }
    }

    r.ConsumeStr(&mText.TextStyle)
    r.ConsumeCoordinates(mText.Vector[:])

    r.ConsumeNumber(73, DEC_RADIX, "line spacing", &mText.LineSpacing)
    r.ConsumeFloat(44, "line spacing factor", nil)

    dxf.MTexts = append(dxf.MTexts, mText)
    return r.Err()
}

func ParseHatch(r *Reader, dxf *drawing.Dxf) error {
    if extractHandleAndOwner(r, &Handle, &Owner) != nil {
        return r.Err()
    }

    hatch := entity.NewMText(Handle, Owner) // TODO: hatch

    if parseAcDbEntityE(r, hatch) != nil || r.AssertNextLine("AcDbHatch") != nil {
        return r.Err()
    }

    // 10,20,30
    coords := [3]float64{0.0,0.0,0.0}
    r.ConsumeCoordinates(coords[:])

    // TODO: [210/220/230] Extrustion direction (only need 2D maybe later)
    r.ConsumeCoordinates(coords[:])

    r.ConsumeStr(nil) // pattern name
    r.ConsumeNumber(70, DEC_RADIX, "solid fill flag", nil)
    r.ConsumeNumber(71, DEC_RADIX, "associativity flag", nil)

    // Number of boundary paths?
    r.ConsumeNumber(91, DEC_RADIX, "boundary paths", nil)

    pt := uint64(0)
    r.ConsumeNumber(92, DEC_RADIX, "boundary path type", &pt)

    if pt & 2 > 0 {
        r.ConsumeNumber(72, DEC_RADIX, "has bulge flag", nil)
        r.ConsumeNumber(73, DEC_RADIX, "is closed flag", nil)

        n := uint64(0)
        r.ConsumeNumber(93, DEC_RADIX, "number fo polyline vertices", &n)

        coord2D := [2]float64{0.0,0.0}

        for i := uint64(0); i < n; i++ {
            r.ConsumeCoordinates(coord2D[:])
            r.ConsumeFloatIf(42, "expected bulge", nil)
        }
    } else {
        n, t := uint64(0), uint64(0)

        r.ConsumeNumber(93, DEC_RADIX, "number of edges in this boundary path", &n)
        r.ConsumeNumber(72, DEC_RADIX, "edge type data", &t)

        coord2D := [2]float64{0.0,0.0}

        switch t {
        case 1:
            // Parse Line
            for i := uint64(0); i < n; i++ {
                r.ConsumeCoordinates(coord2D[:])
            }
        case 2:
            // Circular arc
            for i := uint64(0); i < n; i++ {
                r.ConsumeCoordinates(coord2D[:])
            }

            r.ConsumeNumber(40, DEC_RADIX, "radius", nil)
            r.ConsumeNumber(50, DEC_RADIX, "start angle", nil)
            r.ConsumeNumber(51, DEC_RADIX, "end angle", nil)
            r.ConsumeNumber(73, DEC_RADIX, "is counterclockwise flag", nil)
        case 3:
            // Elliptic arc
            for i := uint64(0); i < n; i++ { 
                r.ConsumeCoordinates(coord2D[:])
            }

            r.ConsumeNumber(40, DEC_RADIX, "length of minor axis", nil)
            r.ConsumeNumber(50, DEC_RADIX, "start angle", nil)
            r.ConsumeNumber(51, DEC_RADIX, "end angle", nil)
            r.ConsumeNumber(73, DEC_RADIX, "is counterclockwise flag", nil)
        case 4:
            // Spine
            r.ConsumeNumber(94, DEC_RADIX, "degree", nil)
            r.ConsumeNumber(73, DEC_RADIX, "rational", nil)
            r.ConsumeNumber(74, DEC_RADIX, "periodic", nil)
            k := uint64(0)
            r.ConsumeNumber(95, DEC_RADIX, "number of knots", &k)
            r.ConsumeNumber(96, DEC_RADIX, "number of control points", nil)

            for i := uint64(0); i < k; i++ {
                r.ConsumeNumber(40, DEC_RADIX, "knot values", nil)
                r.ConsumeCoordinates(coord2D[:])
            }

            r.ConsumeNumberIf(42, DEC_RADIX, "weights", nil) // optional 1

            r.ConsumeNumber(97, DEC_RADIX, "number of fit data", nil)
            r.ConsumeNumber(11, DEC_RADIX, "X fit datum value", nil)
            r.ConsumeNumber(21, DEC_RADIX, "Y fit datum value", nil)
            r.ConsumeNumber(12, DEC_RADIX, "X start tangent", nil)
            r.ConsumeNumber(22, DEC_RADIX, "Y start tangent", nil)
            r.ConsumeNumber(13, DEC_RADIX, "X end tangent", nil)
            r.ConsumeNumber(23, DEC_RADIX, "Y end tangent", nil)
        default:
            log.Println("[ENTITIES(HATCH - ", Line, " )] invalid edge type data: ", t)
            return NewParseError("invalid edge type data")
        }
    }

    bo, sp := uint64(0), uint64(0)

    r.ConsumeNumber(97, DEC_RADIX, "number of source boundary objects", &bo)

    for i := uint64(0); i < bo; i++ {
        r.ConsumeNumber(330, DEC_RADIX, "reference to source boundary objects", nil)
    }

    r.ConsumeNumber(75, DEC_RADIX, "hatch style", nil)
    r.ConsumeNumber(76, DEC_RADIX, "hatch pattern type", nil)
    r.ConsumeNumber(98, DEC_RADIX, "number of seed points", &sp)

    if sp != 0 {
        log.Fatal("TODO(", Line, "): hatch implement seed points")
    }

    _ = dxf
    return r.Err()
}

func ParseEllipse(r *Reader, dxf *drawing.Dxf) error {
    if extractHandleAndOwner(r, &Handle, &Owner) != nil {
        return r.Err()
    }

    ellipse := entity.NewMText(Handle, Owner) // TODO: ellipse

    if parseAcDbEntityE(r, ellipse) != nil || r.AssertNextLine("AcDbEllipse") != nil {
        return r.Err()
    }

    coord3D := [3]float64{0.0,0.0,0.0}

    r.ConsumeCoordinates(coord3D[:]) // Center point
    r.ConsumeCoordinates(coord3D[:]) // endpoint of major axis

    // optional default = 0, 0, 1
    // XYZ extrusion direction
    r.ConsumeCoordinatesIf(210, coord3D[:])

    r.ConsumeFloat(40, "ratio of minor axis to major axis", nil)
    r.ConsumeFloat(41, "start parameter", nil)
    r.ConsumeFloat(42, "end parameter", nil)

    _ = dxf
    return r.Err()
}

func ParsePoint(r *Reader, dxf *drawing.Dxf) error {
    if extractHandleAndOwner(r, &Handle, &Owner) != nil {
        return r.Err()
    }

    point := entity.NewMText(Handle, Owner) // TODO: point

    if parseAcDbEntityE(r, point) != nil || r.AssertNextLine("AcDbPoint") != nil {
        return r.Err()
    }

    coord3D := [3]float64{0.0,0.0,0.0}

    r.ConsumeCoordinates(coord3D[:]) // Point location
    r.ConsumeNumberIf(39, DEC_RADIX, "thickness", nil)

    // optional default = 0, 0, 1
    // XYZ extrusion direction
    r.ConsumeCoordinatesIf(210, coord3D[:]) // Point location
    r.ConsumeFloatIf(50, "angle of the x axis", nil)

    _ = dxf
    return r.Err()
}

// TODO: have to implement block section first
func ParseInsert(r *Reader, dxf *drawing.Dxf) error {
    if extractHandleAndOwner(r, &Handle, &Owner) != nil {
        return r.Err()
    }

    insert := entity.NewMText(Handle, Owner) // TODO: insert

    if parseAcDbEntityE(r, insert) != nil || r.AssertNextLine("AcDbBlockReference") != nil {
        return r.Err()
    }

    coord3D := [3]float64{0.0,0.0,0.0}

    // Variable attributes-follow flag default = 0
    r.ConsumeStrIf(66, nil)
    r.ConsumeStr(nil) // Block name
    r.ConsumeCoordinates(coord3D[:]) // insertion point

    r.ConsumeFloatIf(41, "x scale factor default 1", nil)
    r.ConsumeFloatIf(42, "x scale factor default 1", nil)
    r.ConsumeFloatIf(43, "x scale factor default 1", nil)

    r.ConsumeFloatIf(50, "rotation angle default 0", nil)
    r.ConsumeFloatIf(70, "column count default 1", nil)
    r.ConsumeFloatIf(71, "row count default 1", nil)

    r.ConsumeFloatIf(44, "column spacing default 0", nil)
    r.ConsumeFloatIf(45, "row spacing default 0", nil)

    // optional default = 0, 0, 1
    // XYZ extrusion direction
    r.ConsumeCoordinatesIf(210, coord3D[:])

    // TODO: parse insert
    // attrib =>  987212

    _ = dxf
    return nil
}
