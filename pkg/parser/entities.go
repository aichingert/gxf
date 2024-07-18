package parser

import (
	"fmt"
	"log"

	"github.com/aichingert/dxf/pkg/entity"
)

func ParseEntities(r *Reader, entities entity.Entities) {
	for r.ScanDxfLine() {
		switch r.DxfLine().Line {
		case "LINE":
			ParseLine(r, entities)
		case "POLYLINE":
			ParsePolyline(r, entities)
		case "LWPOLYLINE":
			ParseLwPolyline(r, entities)
		case "ARC":
			ParseArc(r, entities)
		case "CIRCLE":
			ParseCircle(r, entities)
		case "TEXT":
			ParseText(r, entities)
		case "MTEXT":
			ParseMText(r, entities)
		case "HATCH":
			ParseHatch(r, entities)
		case "ELLIPSE":
			ParseEllipse(r, entities)
		case "SPLINE":
			ParseSpline(r, entities)
		case "SOLID":
			ParseSolid(r, entities)
		case "POINT":
			ParsePoint(r, entities)
		case "DIMENSION":
			ParseDimension(r, entities)
		case "REGION":
			ParseRegion(r, entities)
		case "VIEWPORT":
			ParseViewport(r, entities)
		case "ATTDEF":
			ParseAttdef(r, entities)
		case "INSERT":
			ParseInsert(r, entities)
		case "ENDSEC":
			fallthrough
		case "ENDBLK":
			return
		default:
			r.err = NewParseError(fmt.Sprintf("unknown entity: %s", r.DxfLine().Line))
			return
		}

		peek, err := r.PeekCode()
		for err == nil && peek != 0 {
			r.ConsumeStr(nil)
			peek, err = r.PeekCode()
		}
	}
}

func ParseLine(r *Reader, entities entity.Entities) {
	line := entity.NewLine()

	ParseAcDbEntity(r, line.Entity)
	ParseAcDbLine(r, line)

	entities.AppendLine(line)
}

// TODO: create polyline and lwpolyline
func ParsePolyline(r *Reader, entities entity.Entities) {
	polyline := entity.NewPolyline()

	ParseAcDbEntity(r, polyline.Entity)
	ParseAcDb2dPolyline(r, polyline)

	peek, err := r.PeekCode()
	for err == nil && peek != 0 {
		r.ConsumeStr(nil)
		peek, err = r.PeekCode()
	}

	for r.ScanDxfLine() {
		switch r.DxfLine().Line {
		case "VERTEX":
			ParseVertex(r, entities)
		case "SEQEND":
			// marks end of insert
			ParseAcDbEntity(r, polyline.Entity)
			return
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
}

func ParseLwPolyline(r *Reader, entities entity.Entities) {
	polyline := entity.NewPolyline()

	ParseAcDbEntity(r, polyline.Entity)
	ParseAcDbPolyline(r, polyline)

	entities.AppendPolyline(polyline)
}

func ParseArc(r *Reader, entities entity.Entities) {
	arc := entity.NewArc()

	ParseAcDbEntity(r, arc.Entity)
	ParseAcDbCircle(r, arc.Circle)
	ParseAcDbArc(r, arc)

	entities.AppendArc(arc)
}

func ParseCircle(r *Reader, entities entity.Entities) {
	circle := entity.NewCircle()

	ParseAcDbEntity(r, circle.Entity)
	ParseAcDbCircle(r, circle)

	entities.AppendCircle(circle)
}

func ParseText(r *Reader, entities entity.Entities) {
	text := entity.NewText()

	ParseAcDbEntity(r, text.Entity)
	ParseAcDbText(r, text)

	entities.AppendText(text)
}

func ParseMText(r *Reader, entities entity.Entities) {
	mText := entity.NewMText()

	ParseAcDbEntity(r, mText.Entity)
	ParseAcDbMText(r, mText)

	entities.AppendMText(mText)
}

func ParseHatch(r *Reader, entities entity.Entities) {
	hatch := entity.NewHatch()

	ParseAcDbEntity(r, hatch.Entity)
	ParseAcDbHatch(r, hatch)

	entities.AppendHatch(hatch)
}

func ParseEllipse(r *Reader, entities entity.Entities) {
	ellipse := entity.NewEllipse()

	ParseAcDbEntity(r, ellipse.Entity)
	ParseAcDbEllipse(r, ellipse)

	entities.AppendEllipse(ellipse)
}

// TODO: create entity spline
func ParseSpline(r *Reader, _ entity.Entities) {
	spline := entity.NewMText()

	ParseAcDbEntity(r, spline.Entity)
	ParseAcDbSpline(r, spline)
}

// TODO: create entity solid
func ParseSolid(r *Reader, _ entity.Entities) {
	solid := entity.NewMText()

	ParseAcDbEntity(r, solid.Entity)
	ParseAcDbTrace(r, solid)
}

// TODO: create entity vertex
func ParseVertex(r *Reader, _ entity.Entities) {
	vertex := entity.NewMText()

	ParseAcDbEntity(r, vertex.Entity)
	ParseAcDbVertex(r, vertex)
}

// TODO: create entity point
func ParsePoint(r *Reader, _ entity.Entities) {
	point := entity.NewMText()

	ParseAcDbEntity(r, point.Entity)
	ParseAcDbPoint(r, point)
}

func ParseInsert(r *Reader, entities entity.Entities) {
	insert := entity.NewInsert()

	ParseAcDbEntity(r, insert.Entity)
	ParseAcDbBlockReference(r, insert)

	peek, err := r.PeekCode()
	for err == nil && peek != 0 {
		r.ConsumeStr(nil)
		peek, err = r.PeekCode()
	}

Att:
	for insert.AttributesFollow == 1 && r.ScanDxfLine() {
		switch r.DxfLine().Line {
		case "ATTRIB":
			ParseAttrib(r, insert)
		case "SEQEND":
			// marks end of attributes
			ParseAcDbEntity(r, insert.Entity)
			break Att
		default:
			log.Fatal("[INSERT(", Line, ")] invalid subclass marker ", r.DxfLine().Line)
		}
	}

	entities.AppendInsert(insert)
}

// TODO: create entity DIMENSION
func ParseDimension(r *Reader, _ entity.Entities) {
	throwAway := entity.NewAttdef()

	ParseAcDbEntity(r, throwAway.Entity)
	ParseAcDbDimension(r, throwAway)

	r.ConsumeNumberIf(290, DecRadix, "not documented", nil)
	r.ConsumeStrIf(2, nil)
}

// TODO: create entity region
func ParseRegion(r *Reader, _ entity.Entities) {
	throwAway := entity.NewMText()

	ParseAcDbEntity(r, throwAway.Entity)

	if r.AssertNextLine("AcDbModelerGeometry") != nil {
		return
	}

	r.ConsumeNumberIf(290, DecRadix, "not documented", nil)
	r.ConsumeStrIf(2, nil)
}

// TODO: create entity viewport
func ParseViewport(r *Reader, _ entity.Entities) {
	throwAway := entity.NewMText()

	ParseAcDbEntity(r, throwAway.Entity)
	ParseAcDbViewport(r, throwAway)
}

func ParseAttrib(r *Reader, appender entity.AttribAppender) {
	attrib := entity.NewAttrib()

	ParseAcDbEntity(r, attrib.Entity)
	ParseAcDbText(r, attrib.Text)
	ParseAcDbAttribute(r, attrib)

	appender.AppendAttrib(attrib)
}

func ParseAttdef(r *Reader, _ entity.Entities) {
	attdef := entity.NewAttdef()

	ParseAcDbEntity(r, attdef.Entity)
	ParseAcDbText(r, attdef.Text)
	ParseAcDbAttributeDefinition(r, attdef)
}
