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

func ParseText(r *Reader, dxf *drawing.Dxf) error {
	text := entity.NewText()

	if ParseAcDbEntity(r, text.Entity) != nil ||
		ParseAcDbText(r, text) != nil {
		return r.Err()
	}

	dxf.Texts = append(dxf.Texts, text)
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

func ParseHatch(r *Reader, dxf *drawing.Dxf) error {
	hatch := entity.NewHatch()

	if ParseAcDbEntity(r, hatch.Entity) != nil ||
		ParseAcDbHatch(r, hatch) != nil {
		return r.Err()
	}

	dxf.Hatches = append(dxf.Hatches, hatch)
	return r.Err()
}

func ParseEllipse(r *Reader, dxf *drawing.Dxf) error {
	ellipse := entity.NewEllipse()

	if ParseAcDbEntity(r, ellipse.Entity) != nil ||
		ParseAcDbEllipse(r, ellipse) != nil {
		return r.Err()
	}

	dxf.Ellipses = append(dxf.Ellipses, ellipse)
	return r.Err()
}

// TODO: create entity point
func ParsePoint(r *Reader, dxf *drawing.Dxf) error {
	point := entity.NewMText()

	if ParseAcDbEntity(r, point.Entity) != nil ||
		ParseAcDbPoint(r, point) != nil {
		return r.Err()
	}

	return r.Err()
}

func ParseInsert(r *Reader, dxf *drawing.Dxf) error {
	insert := entity.NewInsert()

	if ParseAcDbEntity(r, insert.Entity) != nil ||
		ParseAcDbBlockReference(r, insert) != nil {
		return r.Err()
	}

	for insert.AttributesFollow == 1 && r.ScanDxfLine() {
		switch r.DxfLine().Line {
		case "ATTRIB":
			Wrap(ParseAttrib, r, dxf)
		case "SEQEND":
			// marks end of insert
			ParseAcDbEntity(r, insert.Entity)
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

func ParseRegion(r *Reader, dxf *drawing.Dxf) error {
	throwAway := entity.NewMText()

	if ParseAcDbEntity(r, throwAway.Entity) != nil ||
		r.AssertNextLine("AcDbModelerGeometry") != nil {
		return r.Err()
	}

	r.ConsumeNumberIf(290, DEC_RADIX, "not documented", nil)
	r.ConsumeStrIf(2, nil)

	return r.Err()
}

func ParseAttrib(r *Reader, dxf *drawing.Dxf) error {
	attrib := entity.NewAttrib()

	if ParseAcDbEntity(r, attrib.Entity) != nil ||
		ParseAcDbText(r, attrib.Text) != nil ||
		ParseAcDbAttribute(r, attrib) != nil {
		return r.Err()
	}

	return r.Err()
}

func ParseAttdef(r *Reader, dxf *drawing.Dxf) error {
	attdef := entity.NewAttdef()

	if ParseAcDbEntity(r, attdef.Entity) != nil ||
		ParseAcDbText(r, attdef.Text) != nil ||
		ParseAcDbAttributeDefinition(r, attdef) != nil {
		return r.Err()
	}

	return r.Err()
}
