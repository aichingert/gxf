package parser

import (
	"log"

	"github.com/aichingert/dxf/pkg/drawing"
	"github.com/aichingert/dxf/pkg/entity"
)

type ParseEntityFunction func(*Reader, entity.Entities) error

var WrappedEntityErr error

func WrapEntity(fn ParseEntityFunction, r *Reader, entities entity.Entities) {
	if WrappedEntityErr != nil {
		return
	}

	WrappedEntityErr = fn(r, entities)
}

func ParseEntities(r *Reader, dxf *drawing.Dxf) error {
	for r.ScanDxfLine() {
		switch r.DxfLine().Line {
		case "LINE":
			WrapEntity(ParseLine, r, dxf.EntitiesData)
		case "LWPOLYLINE":
			WrapEntity(ParsePolyline, r, dxf.EntitiesData)
		case "ARC":
			WrapEntity(ParseArc, r, dxf.EntitiesData)
		case "CIRCLE":
			WrapEntity(ParseCircle, r, dxf.EntitiesData)
		case "TEXT":
			WrapEntity(ParseText, r, dxf.EntitiesData)
		case "MTEXT":
			WrapEntity(ParseMText, r, dxf.EntitiesData)
		case "HATCH":
			WrapEntity(ParseHatch, r, dxf.EntitiesData)
		case "ELLIPSE":
			WrapEntity(ParseEllipse, r, dxf.EntitiesData)
		case "POINT":
			WrapEntity(ParsePoint, r, dxf.EntitiesData)
		case "INSERT":
			WrapEntity(ParseInsert, r, dxf.EntitiesData)
		case "ENDSEC":
			return r.Err()
		default:
			log.Println("[ENTITIES] ", Line, ": ", r.DxfLine().Line)
			return NewParseError("unknown entity")
		}

		if WrappedEntityErr != nil {
			return WrappedEntityErr
		}
	}

	return r.Err()
}

func ParseLine(r *Reader, entities entity.Entities) error {
	line := entity.NewLine()

	if ParseAcDbEntity(r, line.Entity) != nil ||
		ParseAcDbLine(r, line) != nil {
		return r.Err()
	}

	entities.AppendLine(line)
	return r.Err()
}

func ParsePolyline(r *Reader, entities entity.Entities) error {
	polyline := entity.NewPolyline()

	if ParseAcDbEntity(r, polyline.Entity) != nil ||
		ParseAcDbPolyline(r, polyline) != nil {
		return r.Err()
	}

	entities.AppendPolyline(polyline)
	return r.Err()
}

func ParseArc(r *Reader, entities entity.Entities) error {
	arc := entity.NewArc()

	if ParseAcDbEntity(r, arc.Entity) != nil ||
		ParseAcDbCircle(r, arc.Circle) != nil ||
		ParseAcDbArc(r, arc) != nil {
		return r.Err()
	}

	entities.AppendArc(arc)
	return r.Err()
}

func ParseCircle(r *Reader, entities entity.Entities) error {
	circle := entity.NewCircle()

	if ParseAcDbEntity(r, circle.Entity) != nil ||
		ParseAcDbCircle(r, circle) != nil {
		return r.Err()
	}

	entities.AppendCircle(circle)
	return r.Err()
}

func ParseText(r *Reader, entities entity.Entities) error {
	text := entity.NewText()

	if ParseAcDbEntity(r, text.Entity) != nil ||
		ParseAcDbText(r, text) != nil {
		return r.Err()
	}

	entities.AppendText(text)
	return r.Err()
}

func ParseMText(r *Reader, entities entity.Entities) error {
	mText := entity.NewMText()

	if ParseAcDbEntity(r, mText.Entity) != nil ||
		ParseAcDbMText(r, mText) != nil {
		return r.Err()
	}

	entities.AppendMText(mText)
	return r.Err()
}

func ParseHatch(r *Reader, entities entity.Entities) error {
	hatch := entity.NewHatch()

	if ParseAcDbEntity(r, hatch.Entity) != nil ||
		ParseAcDbHatch(r, hatch) != nil {
		return r.Err()
	}

	entities.AppendHatch(hatch)
	return r.Err()
}

func ParseEllipse(r *Reader, entities entity.Entities) error {
	ellipse := entity.NewEllipse()

	if ParseAcDbEntity(r, ellipse.Entity) != nil ||
		ParseAcDbEllipse(r, ellipse) != nil {
		return r.Err()
	}

	entities.AppendEllipse(ellipse)
	return r.Err()
}

// TODO: create entity point
func ParsePoint(r *Reader, entities entity.Entities) error {
	point := entity.NewMText()

	if ParseAcDbEntity(r, point.Entity) != nil ||
		ParseAcDbPoint(r, point) != nil {
		return r.Err()
	}

	return r.Err()
}

func ParseInsert(r *Reader, entities entity.Entities) error {
	insert := entity.NewInsert()

	if ParseAcDbEntity(r, insert.Entity) != nil ||
		ParseAcDbBlockReference(r, insert) != nil {
		return r.Err()
	}

	for insert.AttributesFollow == 1 && r.ScanDxfLine() {
		switch r.DxfLine().Line {
		case "ATTRIB":
			if ParseAttrib(r, insert) != nil {
				return r.Err()
			}
		case "SEQEND":
			// marks end of insert
			ParseAcDbEntity(r, insert.Entity)
			entities.AppendInsert(insert)
			return r.Err()
		default:
			log.Fatal("[INSERT(", Line, ")] invalid subclass marker ", r.DxfLine().Line)
		}

		if WrappedEntityErr != nil {
			return WrappedEntityErr
		}
	}

	return r.Err()
}

func ParseRegion(r *Reader, entities entity.Entities) error {
	throwAway := entity.NewMText()

	if ParseAcDbEntity(r, throwAway.Entity) != nil ||
		r.AssertNextLine("AcDbModelerGeometry") != nil {
		return r.Err()
	}

	r.ConsumeNumberIf(290, DEC_RADIX, "not documented", nil)
	r.ConsumeStrIf(2, nil)

	return r.Err()
}

func ParseAttrib(r *Reader, appender entity.AttribAppender) error {
	attrib := entity.NewAttrib()

	if ParseAcDbEntity(r, attrib.Entity) != nil ||
		ParseAcDbText(r, attrib.Text) != nil ||
		ParseAcDbAttribute(r, attrib) != nil {
		return r.Err()
	}

	appender.AppendAttrib(attrib)
	return r.Err()
}

func ParseAttdef(r *Reader, entities entity.Entities) error {
	attdef := entity.NewAttdef()

	if ParseAcDbEntity(r, attdef.Entity) != nil ||
		ParseAcDbText(r, attdef.Text) != nil ||
		ParseAcDbAttributeDefinition(r, attdef) != nil {
		return r.Err()
	}

	return r.Err()
}
