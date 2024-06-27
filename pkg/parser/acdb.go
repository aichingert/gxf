package parser

import (
	"log"

	"github.com/aichingert/dxf/pkg/entity"
	"github.com/aichingert/dxf/pkg/blocks"
)

func ParseAcDbEntity(r *Reader, entity entity.Entity) error {
	r.ConsumeNumber(5, HEX_RADIX, "handle", entity.GetHandle())

	// TODO: set hard owner/handle to owner dictionary
	if r.ConsumeStrIf(102, nil) { // consumeIf => ex. {ACAD_XDICTIONARY
		r.ConsumeStr(nil) // 360 => hard owner
        for r.ConsumeNumberIf(330, HEX_RADIX, "soft owner", nil) {}
		r.ConsumeStr(nil) // 102 }
	}

	r.ConsumeNumber(330, HEX_RADIX, "owner ptr", entity.GetOwner())

	if r.AssertNextLine("AcDbEntity") != nil {
		return r.Err()
	}

	// TODO: think about paper space visibility
	r.ConsumeStrIf(67, nil)

	r.ConsumeStr(entity.GetLayerName())
	r.ConsumeNumberIf(60, DEC_RADIX, "object visibility", nil)

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

	if !r.ConsumeFloatIf(43, "line width for each vertex", nil) {
		log.Fatal("[ENTITIES(", Line, ")] TODO: implement line width for each vertex")
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

    r.ConsumeStrIf(1001, nil) 
    r.ConsumeStrIf(1070, nil) 
    r.ConsumeStrIf(1071, nil) 
    r.ConsumeStrIf(1005, nil) 

	return r.Err()
}

func ParseAcDbCircle(r *Reader, circle *entity.Circle) error {
	if r.AssertNextLine("AcDbCircle") != nil {
		return r.Err()
	}

	r.ConsumeCoordinates(circle.Coordinates[:])
	r.ConsumeFloat(40, "expected radius", &circle.Radius)

    r.ConsumeStrIf(1001, nil)
    r.ConsumeNumberIf(1070, DEC_RADIX, "not sure", nil)
    r.ConsumeNumberIf(1071, DEC_RADIX, "not sure", nil)
    r.ConsumeNumberIf(1005, DEC_RADIX, "not sure", nil)

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

	r.ConsumeFloat(40, "expected text height", nil)
	r.ConsumeStr(nil) // [1] default value of the string itself

	r.ConsumeFloatIf(50, "text rotation default 0", nil)
	r.ConsumeFloatIf(41, "relative x scale factor default 1", nil)
	r.ConsumeFloatIf(51, "oblique angle default 0", nil)

	r.ConsumeStrIf(7, nil) // text style name default STANDARD

	r.ConsumeNumberIf(71, DEC_RADIX, "text generation flags default 0", nil)
	r.ConsumeNumberIf(72, DEC_RADIX, "horizontal text justification default 0", nil)

	r.ConsumeCoordinatesIf(11, coords3D[:])
	// XYZ extrusion direction
	// optional default 0, 0, 1
	r.ConsumeCoordinatesIf(210, coords3D[:])

    line, _ := r.PeekLine()
    if line == "AcDbText" {
        r.ConsumeStr(nil) // second AcDbText (optional)
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
    boundaryPaths := uint64(0)
	r.ConsumeNumber(91, DEC_RADIX, "boundary paths", &boundaryPaths)

    for i := uint64(0); i < boundaryPaths; i++ {
        pathTypeFlag := uint64(0)
        // [92] Boundary path type flag (bit coded):
        // 0 = Default | 1 = External | 2  = Polyline
        // 4 = Derived | 8 = Textbox  | 16 = Outermost
        r.ConsumeNumber(92, DEC_RADIX, "boundary path type flag", &pathTypeFlag)

        if pathTypeFlag&2 == 2 {

            r.ConsumeNumber(72, DEC_RADIX, "has bulge flag", nil)
            r.ConsumeNumber(73, DEC_RADIX, "is closed flag", nil)

            vertices := uint64(0)
            r.ConsumeNumber(93, DEC_RADIX, "number of polyline vertices", &vertices)

            coord2D := [2]float64{0.0, 0.0}

            for vertex := uint64(0); vertex < vertices; vertex++ {
                r.ConsumeCoordinates(coord2D[:])
                r.ConsumeFloatIf(42, "expected bulge", nil)
            }
        } else {
            edges, edgeType, coord2D := uint64(0), uint64(0), [2]float64{0.0, 0.0}

            r.ConsumeNumber(93, DEC_RADIX, "number of edges in this boundary path", &edges)

            for edge := uint64(0); edge < edges; edge++ {
                r.ConsumeNumber(72, DEC_RADIX, "edge type data", &edgeType)

                switch edgeType {
                case 1: // Line
                    r.ConsumeCoordinates(coord2D[:])
                    r.ConsumeCoordinates(coord2D[:])
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
        }

        boundaryObjectSize, boundaryObjectRef := uint64(0), uint64(0)

        r.ConsumeNumber(97, DEC_RADIX, "number of source boundary objects", &boundaryObjectSize)
        for i := uint64(0); i < boundaryObjectSize; i++ {
            r.ConsumeNumber(330, HEX_RADIX, "reference to source object", &boundaryObjectRef)
        }
    }

    r.ConsumeNumber(75, DEC_RADIX, "hatch style", nil)
    r.ConsumeNumber(76, DEC_RADIX, "hatch pattern type", nil)
    r.ConsumeFloatIf(47, "pixel size used to determine density to perform ray casting", nil)

    seedPoints := uint64(0)
    r.ConsumeNumber(98, DEC_RADIX, "number of seed points", &seedPoints)

    coord2D := [2]float64{0.0, 0.0}

    for seedPoint := uint64(0); seedPoint < seedPoints; seedPoint++ {
        r.ConsumeCoordinates(coord2D[:])
    }

    r.ConsumeNumberIf(450, DEC_RADIX, "indicates solid hatch or gradient", nil)
    r.ConsumeNumberIf(451, DEC_RADIX, "zero is reserved for future use", nil)
    
    // default 0,0
    r.ConsumeFloatIf(460, "rotation angle in radians for gradients", nil)
    r.ConsumeFloatIf(461, "gradient definition", nil)
    r.ConsumeNumberIf(452, DEC_RADIX, "records how colors were defined", nil)
    r.ConsumeFloatIf(462, "color tint value used by dialog", nil)

    nColors := uint64(0)
    r.ConsumeNumberIf(453, DEC_RADIX, "number of colors", &nColors)

    for color := uint64(0); color < nColors; color++ {
        r.ConsumeFloatIf(463, "reserved for future use", nil)
        r.ConsumeNumberIf(63, DEC_RADIX, "not documented", nil)
        r.ConsumeNumberIf(421, DEC_RADIX, "not documented", nil)
    }

    r.ConsumeStrIf(470, nil) // string default = LINEAR

    r.ConsumeStrIf(1001, nil)
    r.ConsumeNumberIf(1070, DEC_RADIX, "not sure", nil)
    r.ConsumeStrIf(1001, nil)
    r.ConsumeNumberIf(1070, DEC_RADIX, "not sure", nil)

    r.ConsumeStrIf(1001, nil) // acad
    r.ConsumeFloatIf(1010, "not sure", nil)
    r.ConsumeFloatIf(1020, "not sure", nil)
    r.ConsumeFloatIf(1030, "not sure", nil)

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

// TODO: implement block reference
func ParseAcDbBlockReference(r *Reader, reference *entity.MText) error {
    if r.AssertNextLine("AcDbBlockReference") != nil {
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

    _ = reference
    return r.Err()
}

func ParseAcDbBlockBegin(r *Reader, block *blocks.Block) error {
    if r.AssertNextLine("AcDbBlockBegin") != nil {
        return r.Err()
    }

    r.ConsumeStr(nil) // [2] block name
    r.ConsumeNumber(70, DEC_RADIX, "block-type flag", nil)

    coords3D := [3]float64{0.0, 0.0, 0.0}
    r.ConsumeCoordinates(coords3D[:])

    r.ConsumeStr(nil) // [3] block name
    r.ConsumeStr(nil) // [1] Xref path name

    return r.Err()
}

// TODO: implement attribute
func ParseAcDbAttribute(r *Reader, attribute *entity.MText) error {
    if r.AssertNextLine("AcDbAttribute") != nil {
        return r.Err()
    }

    r.ConsumeStr(nil) // [2] Attribute tag
    r.ConsumeNumber(70, DEC_RADIX, "attribute flag", nil)
    r.ConsumeNumberIf(74, DEC_RADIX, "vertical text justification type default 0", nil) // group code 73 TEXT
    r.ConsumeNumberIf(280, DEC_RADIX, "version number", nil)

    r.ConsumeNumberIf(73, DEC_RADIX, "field length", nil)
    r.ConsumeNumberIf(50, DEC_RADIX, "text rotation default 0", nil)
    r.ConsumeNumberIf(41, DEC_RADIX, "relative x scale factor (width) default 0", nil) // adjusted when fit-type text is used
    r.ConsumeNumberIf(51, DEC_RADIX, "oblique angle default 0", nil)
    r.ConsumeStrIf(7, nil) // text style name default STANDARD
    r.ConsumeNumberIf(71, DEC_RADIX, "text generation flags default 0", nil)
    r.ConsumeNumberIf(72, DEC_RADIX, "horizontal text justification type default 0", nil)

    coord3D := [3]float64{0.0, 0.0, 0.0}
    r.ConsumeCoordinatesIf(11, coord3D[:])
    r.ConsumeCoordinatesIf(210, coord3D[:])

    // TODO: maybe continues?
    // not documented

    r.ConsumeStrIf(1001, nil) // AcadAnnotative
    r.ConsumeStrIf(1000, nil) // AnnotativeData
    r.ConsumeStrIf(1002, nil) // {
    r.ConsumeNumberIf(1070, DEC_RADIX, "not sure", nil)
    r.ConsumeNumberIf(1070, DEC_RADIX, "not sure", nil)
    r.ConsumeNumberIf(1002, DEC_RADIX, "not sure", nil)
    // }

    r.ConsumeStrIf(1001, nil) // AcDbBlockRepETag
    r.ConsumeNumberIf(1070, DEC_RADIX, "not sure", nil)
    r.ConsumeNumberIf(1071, DEC_RADIX, "not sure", nil)
    r.ConsumeNumberIf(1005, DEC_RADIX, "not sure", nil)

    return r.Err()
}

func ParseAcDbAttributeDefinition(r *Reader, attdef *entity.MText) error {
    if r.AssertNextLine("AcDbAttributeDefinition") != nil {
        return r.Err()
    }

    r.ConsumeStr(nil) // [3] prompt string
    r.ConsumeStr(nil) // [2] tag string
    r.ConsumeNumber(70, DEC_RADIX, "attribute flags", nil)
    r.ConsumeFloatIf(73, "field length", nil)
    r.ConsumeFloatIf(74, "vertical text justification type default 0", nil)

    r.ConsumeNumber(280, DEC_RADIX, "lock position flag", nil) 

    r.ConsumeNumberIf(71, DEC_RADIX, "not documented", nil)
    r.ConsumeNumberIf(72, DEC_RADIX, "not documented", nil)
    coords3D := [3]float64{0.0, 0.0, 0.0}

    r.ConsumeCoordinatesIf(11, coords3D[:])

    r.ConsumeStrIf(1001, nil) // AcadAnnotative
    r.ConsumeStrIf(1000, nil) // AnnotativeData
    r.ConsumeStrIf(1002, nil) // {
    r.ConsumeNumberIf(1070, DEC_RADIX, "not sure", nil)
    r.ConsumeNumberIf(1070, DEC_RADIX, "not sure", nil)
    r.ConsumeNumberIf(1002, DEC_RADIX, "not sure", nil)
    // }

    r.ConsumeStrIf(1001, nil) // AcDbBlockRepETag
    r.ConsumeNumberIf(1070, DEC_RADIX, "not sure", nil)
    r.ConsumeNumberIf(1071, DEC_RADIX, "not sure", nil)
    r.ConsumeNumberIf(1005, DEC_RADIX, "not sure", nil)

    return r.Err()
}
