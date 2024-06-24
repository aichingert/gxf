package parser

import (
    "log"

    "github.com/aichingert/dxf/pkg/entity"
)

func ParseAcDbEntity(r *Reader, entity entity.Entity) error {
    if r.AssertNextLine("AcDbEntity") != nil {
        return r.Err()
    }

    // TODO: think about paper space visibility
    r.ConsumeStrIf(67, nil)

    r.ConsumeStr(entity.GetLayerName())
    return r.Err()
}

func ParseAcDbLine(r *Reader, line *entity.Line) error {
    if r.AssertNextLine("AcDbLine") != nil {
        return r.Err()
    }

    r.ConsumeCoordinates(line.Src[:])
    r.ConsumeCoordinates(line.Dst[:])

    return r.Err()
}

func ParseAcDbPolyline(r *Reader, polyline *entity.Polyline) error {
    if r.AssertNextLine("AcDbPolyline") != nil {
        return r.Err()
    }

    r.ConsumeNumber(90, DEC_RADIX, "number of vertices", &polyline.Vertices)
    r.ConsumeNumber(70, DEC_RADIX, "polyline flag", &polyline.Flag)

    if r.ConsumeFloatIf(43, "line width for each vertex", nil) {
        log.Fatal("[ENTITIES] TODO: implement line width for each vertex")
    }

    for i := uint64(0); i < polyline.Vertices; i++ {
        bulge := 0.0
        coords2D := [2]float64{0.0, 0.0}

        r.ConsumeCoordinates(coords2D[:])
        r.ConsumeFloatIf(42, "expected bulge", &bulge)

        if r.Err() != nil {
            return r.Err()
        }

        polyline.PolylineAppendCoordinate(coords2D, bulge)
    }

    return r.Err()
}

func ParseAcDbCircle(r *Reader, circle *entity.Circle) error {
    if r.AssertNextLine("AcDbCircle") != nil {
        return r.Err()
    }

    r.ConsumeCoordinates(circle.Coordinates[:])
    r.ConsumeFloat(40, "expected radius", &circle.Radius)

    return r.Err()
}

func ParseAcDbArc(r *Reader, arc *entity.Arc) error {
    if r.AssertNextLine("AcDbArc") != nil {
        return r.Err()
    }

    r.ConsumeFloat(50, "expected startAngle", &arc.StartAngle)
    r.ConsumeFloat(51, "expected endAngle", &arc.EndAngle)

    return r.Err()
}

// TODO: change to text entity
func ParseAcDbText(r *Reader, text *entity.MText) error {
    if r.AssertNextLine("AcDbText") != nil {
        return r.Err()
    }

    log.Fatal("TODO: implement text entity")

    return r.Err()
}

func ParseAcDbMText(r *Reader, mText *entity.MText) error {
    if r.AssertNextLine("AcDbMText") != nil {
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

    return r.Err()
}
