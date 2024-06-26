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

	if ParseAcDbEntity(r, insert.Entity) != nil || r.AssertNextLine("AcDbBlockReference") != nil {
		return r.Err()
	}

	coord3D := [3]float64{0.0, 0.0, 0.0}

	// Variable attributes-follow flag default = 0
	r.ConsumeStrIf(66, nil)
	r.ConsumeStr(nil)                // Block name
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
