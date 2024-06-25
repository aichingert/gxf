package parser

import (
	"log"

	"github.com/aichingert/dxf/pkg/entity"
)

func ParseAcDbEntity(r *Reader, entity entity.Entity) error {
	r.ConsumeNumber(5, HEX_RADIX, "handle", entity.GetHandle())

	// TODO: set hard owner/handle to owner dictionary
	if r.ConsumeStrIf(102, nil) { // consumeIf => ex. {ACAD_XDICTIONARY
		r.ConsumeStr(nil) // 360 => hard owner
		r.ConsumeStr(nil) // 102 }
	}

	r.ConsumeNumber(330, HEX_RADIX, "owner ptr", entity.GetOwner())

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

	r.ConsumeFloatIf(39, "expected thickness", nil)

	coords3D := [3]float64{0.0, 0.0, 0.0}
	r.ConsumeCoordinates(coords3D[:])
	r.ConsumeCoordinates(coords3D[:])

	r.ConsumeFloat(40, "expected text height", nil)
	r.ConsumeStr(nil) // [1] default value of the string itself

	r.ConsumeNumberIf(50, DEC_RADIX, "text rotation default 0", nil)
	r.ConsumeFloatIf(41, "relative x scale factor default 1", nil)
	r.ConsumeFloatIf(51, "oblique angle default 0", nil)

	r.ConsumeStrIf(7, nil) // text style name default STANDARD

	r.ConsumeNumberIf(71, DEC_RADIX, "text generation flags default 0", nil)
	r.ConsumeNumberIf(72, DEC_RADIX, "horizontal text justification default 0", nil)

	r.ConsumeCoordinatesIf(11, coords3D[:])
	// XYZ extrusion direction
	// optional default 0, 0, 1
	r.ConsumeCoordinatesIf(210, coords3D[:])

	if r.AssertNextLine("AcDbText") != nil {
		return r.Err()
	}

	// group 72 and 73 integer codes
	// https://help.autodesk.com/view/OARX/2024/ENU/?guid=GUID-62E5383D-8A14-47B4-BFC4-35824CAE8363

	r.ConsumeNumberIf(73, DEC_RADIX, "vertical text justification type default 0", nil)

	_ = text
	_ = coords3D
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
	if err != nil {
		return err
	}

	for code == 1 || code == 3 {
		line, err := r.ConsumeDxfLine()
		if err != nil {
			return err
		}

		mText.Text = append(mText.Text, line.Line)

		code, err = r.PeekCode()
		if err != nil {
			return err
		}
	}

	r.ConsumeStr(&mText.TextStyle)
	r.ConsumeCoordinates(mText.Vector[:])

	r.ConsumeNumber(73, DEC_RADIX, "line spacing", &mText.LineSpacing)
	r.ConsumeFloat(44, "line spacing factor", nil)

	return r.Err()
}

// TODO: replace with entity.Hatch
func ParseAcDbHatch(r *Reader, hatch *entity.MText) error {
	if r.AssertNextLine("AcDbHatch") != nil {
		return r.Err()
	}

	coords3D := [3]float64{0.0, 0.0, 0.0}
	r.ConsumeCoordinates(coords3D[:])

	// TODO: [210/220/230] extrusion direction (only need 2d maybe later)
	r.ConsumeCoordinates(coords3D[:])

	r.ConsumeStr(nil) // pattern name
	r.ConsumeNumber(70, DEC_RADIX, "solid fill flag", nil)
	r.ConsumeNumber(71, DEC_RADIX, "associativity flag", nil)

	// number of boundary paths?
	r.ConsumeNumber(91, DEC_RADIX, "boundary paths", nil)

	pathTypeFlag := uint64(0)
	// [92] Boundary path type flag (bit coded):
	// 0 = Default | 1 = External | 2  = Polyline
	// 4 = Derived | 8 = Textbox  | 16 = Outermost
	r.ConsumeNumber(92, DEC_RADIX, "boundary path type flag", &pathTypeFlag)

	if pathTypeFlag&2 == 2 {
		r.ConsumeNumber(72, DEC_RADIX, "has bulge flag", nil)
		r.ConsumeNumber(73, DEC_RADIX, "is closed flag", nil)

		vertices := uint64(0)
		coord2D := [2]float64{0.0, 0.0}

		r.ConsumeNumber(93, DEC_RADIX, "number of polyline vertices", &vertices)

		for vertex := uint64(0); vertex < vertices; vertex++ {
			r.ConsumeCoordinates(coord2D[:])
			r.ConsumeFloatIf(42, "expected bulge", nil)
		}
	} else {
		edges, edgeType, coord2D := uint64(0), uint64(0), [2]float64{0.0, 0.0}

		r.ConsumeNumber(93, DEC_RADIX, "number of edges in this boundary path", &edges)
		r.ConsumeNumber(72, DEC_RADIX, "edge type data", &edgeType)

		switch edgeType {
		case 1: // Line
			for edge := uint64(0); edge < edges; edge++ {
				r.ConsumeCoordinates(coord2D[:])
			}
		case 2: // Circular arc
			// TODO
			log.Fatal("hatch circular arc")
		case 3: // Elliptic arc
			log.Fatal("hatch elliptic arc")
		case 4: // Spine
			log.Fatal("hatch spine")
		default:
			log.Println("[AcDbHatch(", Line, ")] invalid edge type data", edgeType)
			return NewParseError("invalid edge type data")
		}
	}

	boundaryObjectSize, boundaryObjectRef := uint64(0), uint64(0)

	r.ConsumeNumber(97, DEC_RADIX, "number of source boundary objects", &boundaryObjectSize)
	for i := uint64(0); i < boundaryObjectSize; i++ {
		r.ConsumeNumber(330, DEC_RADIX, "reference to source object", &boundaryObjectRef)
	}

	r.ConsumeNumber(75, DEC_RADIX, "hatch style", nil)
	r.ConsumeNumber(76, DEC_RADIX, "hatch pattern type", nil)
	r.ConsumeNumber(98, DEC_RADIX, "number of seed points", nil)

	_ = hatch
	return r.Err()
}

// TODO: implement entity entity.Ellipse
func ParseAcDbEllipse(r *Reader, ellipse *entity.MText) error {
	if r.AssertNextLine("AcDbEllipse") != nil {
		return r.Err()
	}

	coord3D := [3]float64{0.0, 0.0, 0.0}

	r.ConsumeCoordinates(coord3D[:]) // Center point
	r.ConsumeCoordinates(coord3D[:]) // Endpoint of major axis

	// XYZ extrusion direction
	// optional default = 0, 0, 1
	r.ConsumeCoordinatesIf(210, coord3D[:])

	r.ConsumeFloat(40, "ratio of minor axis to major axis", nil)
	r.ConsumeFloat(41, "start parameter", nil)
	r.ConsumeFloat(42, "end parameter", nil)

	_ = ellipse
	return r.Err()
}

// TODO: implement entity entity.Point
func ParseAcDbPoint(r *Reader, point *entity.MText) error {
	if r.AssertNextLine("AcDbPoint") != nil {
		return r.Err()
	}

	coord3D := [3]float64{0.0, 0.0, 0.0}

	r.ConsumeCoordinates(coord3D[:])
	r.ConsumeNumberIf(39, DEC_RADIX, "thickness", nil)

	// XYZ extrusion direction
	// optional default 0, 0, 1
	r.ConsumeCoordinatesIf(210, coord3D[:])
	r.ConsumeFloatIf(50, "angle of the x axis", nil)

	_ = point
	return r.Err()
}
