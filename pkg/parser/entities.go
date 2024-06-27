package parser

import (
	"log"

	"github.com/aichingert/dxf/pkg/drawing"
	"github.com/aichingert/dxf/pkg/entity"
)

func ParseEntities(r *Reader, dxf *drawing.Dxf) error {
	for r.ScanDxfLine() {
		switch r.DxfLine().Line {
		case "LINE":
			Wrap(ParseLine, r, dxf)
		case "LWPOLYLINE":
			Wrap(ParsePolyline, r, dxf)
		case "ARC":
			Wrap(ParseArc, r, dxf)
		case "CIRCLE":
			Wrap(ParseCircle, r, dxf)
		case "TEXT":
			Wrap(ParseText, r, dxf)
		case "MTEXT":
			Wrap(ParseMText, r, dxf)
		case "HATCH":
			Wrap(ParseHatch, r, dxf)
		case "ELLIPSE":
			Wrap(ParseEllipse, r, dxf)
		case "POINT":
			Wrap(ParsePoint, r, dxf)
		case "INSERT":
			Wrap(ParseInsert, r, dxf)
        case "ENDSEC":
            return r.Err()
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

func ParseLine(r *Reader, dxf *drawing.Dxf) error {
	line := entity.NewLine()

	if ParseAcDbEntity(r, line.Entity) != nil ||
		ParseAcDbLine(r, line) != nil {
		return r.Err()
	}

	dxf.Lines = append(dxf.Lines, line)
	return r.Err()
}

func ParsePolyline(r *Reader, dxf *drawing.Dxf) error {
	polyline := entity.NewPolyline()

	if ParseAcDbEntity(r, polyline.Entity) != nil ||
		ParseAcDbPolyline(r, polyline) != nil {
		return r.Err()
	}

	dxf.Polylines = append(dxf.Polylines, polyline)
	return r.Err()
}

func ParseArc(r *Reader, dxf *drawing.Dxf) error {
	arc := entity.NewArc()

	if ParseAcDbEntity(r, arc.Entity) != nil ||
		ParseAcDbCircle(r, arc.Circle) != nil ||
		ParseAcDbArc(r, arc) != nil {
		return r.Err()
	}

	dxf.Arcs = append(dxf.Arcs, arc)
	return r.Err()
}

func ParseCircle(r *Reader, dxf *drawing.Dxf) error {
	circle := entity.NewCircle()

	if ParseAcDbEntity(r, circle.Entity) != nil ||
		ParseAcDbCircle(r, circle) != nil {
		return r.Err()
	}

	dxf.Circles = append(dxf.Circles, circle)
	return r.Err()
}

// TODO: create entity text
func ParseText(r *Reader, dxf *drawing.Dxf) error {
	text := entity.NewMText()

	if ParseAcDbEntity(r, text.Entity) != nil ||
		ParseAcDbText(r, text) != nil {
		return r.Err()
	}

	_ = dxf
	return r.Err()
}

func ParseMText(r *Reader, dxf *drawing.Dxf) error {
	mText := entity.NewMText()

	if ParseAcDbEntity(r, mText.Entity) != nil ||
		ParseAcDbMText(r, mText) != nil {
		return r.Err()
	}

	dxf.MTexts = append(dxf.MTexts, mText)
	return r.Err()
}

// TODO: create hatch
func ParseHatch(r *Reader, dxf *drawing.Dxf) error {
	hatch := entity.NewMText()

	if ParseAcDbEntity(r, hatch.Entity) != nil ||
		ParseAcDbHatch(r, hatch) != nil {
		return r.Err()
	}

	_ = dxf
	return r.Err()
}

// TODO: create entity ellipse
func ParseEllipse(r *Reader, dxf *drawing.Dxf) error {
	ellipse := entity.NewMText()

	if ParseAcDbEntity(r, ellipse.Entity) != nil ||
		ParseAcDbEllipse(r, ellipse) != nil {
		return r.Err()
	}

	_ = dxf
	return r.Err()
}

// TODO: create entity point
func ParsePoint(r *Reader, dxf *drawing.Dxf) error {
	point := entity.NewMText()

	if ParseAcDbEntity(r, point.Entity) != nil ||
		ParseAcDbPoint(r, point) != nil {
		return r.Err()
	}

	_ = dxf
	return r.Err()
}

// TODO: have to implement block section first
func ParseInsert(r *Reader, dxf *drawing.Dxf) error {
	insert := entity.NewMText() // TODO: insert

	if ParseAcDbEntity(r, insert.Entity) != nil || 
        ParseAcDbBlockReference(r, insert) != nil {
		return r.Err()
	}

    for r.ScanDxfLine() {
        switch r.DxfLine().Line {
        case "ATTRIB":
            Wrap(ParseAttrib, r, dxf)
        case "ATTDEF":
            Wrap(ParseAttDef, r, dxf)
        case "HATCH":
            Wrap(ParseHatch, r, dxf)
        case "MTEXT":
            Wrap(ParseMText, r, dxf)
        case "LINE":
            Wrap(ParseLine, r, dxf)
        case "LWPOLYLINE":
            Wrap(ParsePolyline, r, dxf)
        case "CIRCLE":
            Wrap(ParseCircle, r, dxf)
        case "INSERT":
            Wrap(ParseInsert, r, dxf)
        case "ENDBLK":
            ParseAcDbEntity(r, insert.Entity) // insert does not end with seqend
            return r.Err()
        case "AcDbBlockEnd":
            return r.Err()
        case "SEQEND":
            ParseAcDbEntity(r, insert.Entity) // marks end of insert 
            return r.Err()
        default:
            log.Fatal("[INSERT(", Line, ")] invalid subclass marker ", r.DxfLine().Line)
        }

        if WrappedErr != nil {
            return WrappedErr
        }
    }

	return r.Err()
}

// TODO: implement attrib
func ParseAttrib(r *Reader, dxf *drawing.Dxf) error {
    attrib := entity.NewMText() 

    if ParseAcDbEntity(r, attrib.Entity) != nil ||
        ParseAcDbText(r, attrib)         != nil ||
        ParseAcDbAttribute(r, attrib)    != nil {
        return r.Err()
    }

    _ = dxf
    return r.Err()
}

// TODO: implement type
func ParseAttDef(r *Reader, attdef *drawing.Dxf) error {
    // TODO: include this in attdef type
    text := entity.NewMText() // this is text not mtext

    if ParseAcDbEntity(r, text.Entity) != nil ||
        ParseAcDbText(r, text) != nil || 
        ParseAcDbAttributeDefinition(r, text) != nil {
        return r.Err()
    }

    return r.Err()
}
