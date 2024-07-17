package parser

import (
	"fmt"
	"log"

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

func ParseEntities(r *Reader, entities entity.Entities) error {
	for r.ScanDxfLine() {
		switch r.DxfLine().Line {
		case "LINE":
			WrapEntity(ParseLine, r, entities)
		case "POLYLINE":
			WrapEntity(ParsePolyline, r, entities)
		case "LWPOLYLINE":
			WrapEntity(ParseLwPolyline, r, entities)
		case "ARC":
			WrapEntity(ParseArc, r, entities)
		case "CIRCLE":
			WrapEntity(ParseCircle, r, entities)
		case "TEXT":
			WrapEntity(ParseText, r, entities)
		case "MTEXT":
			WrapEntity(ParseMText, r, entities)
		case "HATCH":
			WrapEntity(ParseHatch, r, entities)
		case "ELLIPSE":
			WrapEntity(ParseEllipse, r, entities)
		case "SOLID":
			WrapEntity(ParseSolid, r, entities)
		case "POINT":
			WrapEntity(ParsePoint, r, entities)
		case "DIMENSION":
			WrapEntity(ParseDimension, r, entities)
		case "REGION":
			WrapEntity(ParseRegion, r, entities)
		case "VIEWPORT":
			WrapEntity(ParseViewport, r, entities)
		case "ATTDEF":
			WrapEntity(ParseAttdef, r, entities)
		case "INSERT":
			WrapEntity(ParseInsert, r, entities)
		case "ENDSEC":
			fallthrough
		case "ENDBLK":
			return r.Err()
		default:
			log.Println("[ENTITIES] ", Line, ": ", r.DxfLine().Line)
			return NewParseError(fmt.Sprintf("unknown entity: %s", r.DxfLine().Line))
		}

		peek, err := r.PeekCode()
		for err == nil && peek != 0 {
			r.ConsumeStr(nil)
			peek, err = r.PeekCode()
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

// TODO: create polyline and lwpolyline
func ParsePolyline(r *Reader, entities entity.Entities) error {
	polyline := entity.NewPolyline()

	if ParseAcDbEntity(r, polyline.Entity) != nil ||
		ParseAcDb2dPolyline(r, polyline) != nil {
		return r.Err()
	}

	peek, err := r.PeekCode()
	for err == nil && peek != 0 {
		r.ConsumeStr(nil)
		peek, err = r.PeekCode()
	}

	for r.ScanDxfLine() {
		switch r.DxfLine().Line {
		case "VERTEX":
			WrapEntity(ParseVertex, r, entities)
		case "SEQEND":
			// marks end of insert
			ParseAcDbEntity(r, polyline.Entity)
			return r.Err()
		default:
			log.Fatal("[", Line, "] Invalid entity: ", r.DxfLine().Line)
		}

		peek, err := r.PeekCode()
		for err == nil && peek != 0 {
			r.ConsumeStr(nil)
			peek, err = r.PeekCode()
		}
	}

	//entities.AppendPolyline(polyline)
	return r.Err()
}

func ParseLwPolyline(r *Reader, entities entity.Entities) error {
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

func ParseSolid(r *Reader, entities entity.Entities) error {
	solid := entity.NewMText()

	if ParseAcDbEntity(r, solid.Entity) != nil ||
		ParseAcDbTrace(r, solid) != nil {
		return r.Err()
	}

	return r.Err()
}

// TODO: create vertex
func ParseVertex(r *Reader, entities entity.Entities) error {
	vertex := entity.NewMText()

	if ParseAcDbEntity(r, vertex.Entity) != nil ||
		ParseAcDbVertex(r, vertex) != nil {
		return r.Err()
	}

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

Att:
	for insert.AttributesFollow == 1 && r.ScanDxfLine() {
		switch r.DxfLine().Line {
		case "ATTRIB":
			if ParseAttrib(r, insert) != nil {
				return r.Err()
			}
		case "SEQEND":
			// marks end of attributes
			ParseAcDbEntity(r, insert.Entity)
			break Att
		default:
			log.Fatal("[INSERT(", Line, ")] invalid subclass marker ", r.DxfLine().Line)
		}

		if WrappedEntityErr != nil {
			return WrappedEntityErr
		}
	}

	entities.AppendInsert(insert)
	return r.Err()
}

func ParseDimension(r *Reader, entities entity.Entities) error {
	throwAway := entity.NewAttdef()

	if ParseAcDbEntity(r, throwAway.Entity) != nil ||
		ParseAcDbDimension(r, throwAway) != nil {
		return r.Err()
	}

	r.ConsumeNumberIf(290, DecRadix, "not documented", nil)
	r.ConsumeStrIf(2, nil)

	return r.Err()
}

func ParseRegion(r *Reader, entities entity.Entities) error {
	throwAway := entity.NewMText()

	if ParseAcDbEntity(r, throwAway.Entity) != nil ||
		r.AssertNextLine("AcDbModelerGeometry") != nil {
		return r.Err()
	}

	r.ConsumeNumberIf(290, DecRadix, "not documented", nil)
	r.ConsumeStrIf(2, nil)

	return r.Err()
}

func ParseViewport(r *Reader, entities entity.Entities) error {
	throwAway := entity.NewMText()

	if ParseAcDbEntity(r, throwAway.Entity) != nil ||
		ParseAcDbViewport(r, throwAway) != nil {
		return r.Err()
	}

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
