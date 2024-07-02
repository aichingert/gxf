package parser

import (
	"log"

	"github.com/aichingert/dxf/pkg/blocks"
	"github.com/aichingert/dxf/pkg/entity"
)

func ParseAcDbEntity(r *Reader, entity entity.Entity) error {
	r.ConsumeNumber(5, HEX_RADIX, "handle", entity.GetHandle())

	// TODO: set hard owner/handle to owner dictionary
	if r.ConsumeStrIf(102, nil) { // consumeIf => ex. {ACAD_XDICTIONARY
		r.ConsumeStr(nil) // 360 => hard owner
		for r.ConsumeNumberIf(330, HEX_RADIX, "soft owner", nil) {
		}
		r.ConsumeStr(nil) // 102 }
	}

	if r.ConsumeStrIf(102, nil) { // consumeIf => ex. {ACAD_XDICTIONARY
		r.ConsumeStr(nil) // 360 => hard owner
		for r.ConsumeNumberIf(330, HEX_RADIX, "soft owner", nil) {
		}
		r.ConsumeStr(nil) // 102 }
	}

	r.ConsumeNumber(330, HEX_RADIX, "owner ptr", entity.GetOwner())

	if r.AssertNextLine("AcDbEntity") != nil {
		return r.Err()
	}

	// TODO: think about paper space visibility
	r.ConsumeStrIf(67, nil)
	r.ConsumeStr(entity.GetLayerName())

	r.ConsumeStrIf(6, nil) // ByBlock
	r.ConsumeNumberIf(62, DEC_RADIX, "color number (present if not bylayer)", nil)
	r.ConsumeFloatIf(48, "linetype scale", nil)
	r.ConsumeNumberIf(60, DEC_RADIX, "object visibility", entity.GetVisibility())

	r.ConsumeNumberIf(420, DEC_RADIX, "24-bit color value", nil)
	r.ConsumeNumberIf(440, DEC_RADIX, "transparency value", nil)
	r.ConsumeNumberIf(370, DEC_RADIX, "not documented", nil)

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
		r.ConsumeNumberIf(91, DEC_RADIX, "vertex identifier", nil)

		if r.Err() != nil {
			return r.Err()
		}

		polyline.AppendPLine(coords2D, bulge)
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

func ParseAcDbText(r *Reader, text *entity.Text) error {
	if r.AssertNextLine("AcDbText") != nil {
		return r.Err()
	}

	r.ConsumeFloatIf(39, "expected thickness", &text.Thickness)
	r.ConsumeCoordinates(text.Coordinates[:])

	r.ConsumeFloat(40, "expected text height", &text.Height)
	r.ConsumeStr(&text.Text) // [1] default value of the string itself

	r.ConsumeFloatIf(50, "text rotation default 0", &text.Rotation)
	r.ConsumeFloatIf(41, "relative x scale factor default 1", &text.XScale)
	r.ConsumeFloatIf(51, "oblique angle default 0", &text.Oblique)

	r.ConsumeStrIf(7, &text.Style) // text style name default STANDARD

	r.ConsumeNumberIf(71, DEC_RADIX, "text generation flags default 0", &text.Flags)
	r.ConsumeNumberIf(72, DEC_RADIX, "horizontal text justification", &text.HJustification)

	r.ConsumeCoordinatesIf(11, text.Vector[:])
	// XYZ extrusion direction
	// optional default 0, 0, 1
	// TODO: support extrusion?
	r.ConsumeCoordinatesIf(210, text.Vector[:])

	line, _ := r.PeekLine()
	if line == "AcDbText" {
		r.ConsumeStr(nil) // second AcDbText (optional)
	}

	// group 72 and 73 integer codes
	// https://help.autodesk.com/view/OARX/2024/ENU/?guid=GUID-62E5383D-8A14-47B4-BFC4-35824CAE8363

	r.ConsumeNumberIf(73, DEC_RADIX, "vertical text justification", &text.VJustification)

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

	r.ConsumeStrIf(7, &mText.TextStyle)
	r.ConsumeCoordinatesIf(11, mText.Vector[:])

	r.ConsumeNumber(73, DEC_RADIX, "line spacing", &mText.LineSpacing)
	r.ConsumeFloat(44, "line spacing factor", nil)

	return HelperParseEmbeddedObject(r)
}

func HelperParseEmbeddedObject(r *Reader) error {
	// Embedded Object
	if r.ConsumeStrIf(101, nil) {
		r.ConsumeNumberIf(70, DEC_RADIX, "not documented", nil)
		coords3D := [3]float64{0.0, 0.0, 0.0}
		r.ConsumeCoordinates(coords3D[:])
		r.ConsumeCoordinatesIf(11, coords3D[:])

		r.ConsumeFloatIf(40, "not documented", nil)
		r.ConsumeFloatIf(41, "not documented", nil)
		r.ConsumeFloatIf(42, "not documented", nil)
		r.ConsumeFloatIf(43, "not documented", nil)
		r.ConsumeFloatIf(46, "not documented", nil)

		r.ConsumeNumberIf(71, DEC_RADIX, "not documented", nil)
		r.ConsumeNumberIf(72, DEC_RADIX, "not documented", nil)
		r.ConsumeStrIf(1, nil)

		r.ConsumeFloatIf(44, "not documented", nil)
		r.ConsumeFloatIf(45, "not documented", nil)

		r.ConsumeNumberIf(73, DEC_RADIX, "not documented", nil)
		r.ConsumeNumberIf(74, DEC_RADIX, "not documented", nil)

		r.ConsumeFloatIf(44, "not documented", nil)
		r.ConsumeFloatIf(46, "not documented", nil)
	}

	return r.Err()
}

func ParseAcDbHatch(r *Reader, hatch *entity.Hatch) error {
	if r.AssertNextLine("AcDbHatch") != nil {
		return r.Err()
	}

	coords3D := [3]float64{0.0, 0.0, 0.0}
	// TODO: elevation ignored since 2d
	r.ConsumeCoordinates(coords3D[:])
	// TODO: [210/220/230] extrusion direction (only need 2d maybe later)
	r.ConsumeCoordinates(coords3D[:])

	r.ConsumeStr(&hatch.PatternName)
	r.ConsumeNumber(70, DEC_RADIX, "solid fill flag", &hatch.SolidFill)
	r.ConsumeNumber(71, DEC_RADIX, "associativity flag", &hatch.Associative)

	boundaryPaths := uint64(0)
	r.ConsumeNumber(91, DEC_RADIX, "boundary paths", &boundaryPaths)

	for i := uint64(0); i < boundaryPaths; i++ {
		pathTypeFlag := uint64(0)
		// [92] Boundary path type flag (bit coded):
		// 0 = Default | 1 = External | 2  = Polyline
		// 4 = Derived | 8 = Textbox  | 16 = Outermost
		r.ConsumeNumber(92, DEC_RADIX, "boundary path type flag", &pathTypeFlag)

		if pathTypeFlag&2 == 2 {
			polyline := entity.NewPolyline()

			// maybe consider?
			r.ConsumeNumber(72, DEC_RADIX, "has bulge flag", nil)
			r.ConsumeNumber(73, DEC_RADIX, "is closed flag", nil)
			r.ConsumeNumber(93, DEC_RADIX, "number of polyline vertices", &polyline.Vertices)

			for vertex := uint64(0); vertex < polyline.Vertices; vertex++ {
				coord2D, bulge := [2]float64{0.0, 0.0}, 0.0
				r.ConsumeCoordinates(coord2D[:])
				r.ConsumeFloatIf(42, "expected bulge", &bulge)
				polyline.AppendPLine(coord2D, bulge)
			}

			hatch.Polylines = append(hatch.Polylines, polyline)
		} else {
			edges, edgeType := uint64(0), uint64(0)

			r.ConsumeNumber(93, DEC_RADIX, "number of edges in this boundary path", &edges)

			for edge := uint64(0); edge < edges; edge++ {
				r.ConsumeNumber(72, DEC_RADIX, "edge type data", &edgeType)

				switch edgeType {
				case 1: // Line
					line := entity.NewLine()
					r.ConsumeCoordinates(line.Src[:2])
					r.ConsumeCoordinates(line.Dst[:2])
					hatch.Lines = append(hatch.Lines, line)
				case 2: // Circular arc
					arc := entity.NewArc()
					r.ConsumeCoordinates(arc.Circle.Coordinates[:2])
					r.ConsumeFloat(40, "radius", &arc.Circle.Radius)

					r.ConsumeFloat(50, "start angle", &arc.StartAngle)
					r.ConsumeFloat(51, "end angle", &arc.EndAngle)
					r.ConsumeNumber(73, DEC_RADIX, "is counterclockwise", &arc.Counterclockwise)
					hatch.Arcs = append(hatch.Arcs, arc)
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

	r.ConsumeNumber(75, DEC_RADIX, "hatch style", &hatch.Style)
	r.ConsumeNumber(76, DEC_RADIX, "hatch pattern type", &hatch.Pattern)
	r.ConsumeFloatIf(52, "hatch pattern angle", &hatch.Angle)
	r.ConsumeFloatIf(41, "hatch pattern scale or spacing", &hatch.Scale)
	r.ConsumeNumberIf(77, DEC_RADIX, "hatch pattern double flag", &hatch.Double)

	patternDefinitions := uint64(0)

	r.ConsumeNumberIf(78, DEC_RADIX, "number of pattern definition lines", &patternDefinitions)

	for i := uint64(0); i < patternDefinitions; i++ {
		base, offset, angle := [2]float64{0.0, 0.0}, [2]float64{0.0, 0.0}, 0.0
		dashes, dashLen := []float64{}, 0.0

		r.ConsumeFloat(53, "pattern line angle", &angle)
		r.ConsumeFloat(43, "pattern line base point x", &base[0])
		r.ConsumeFloat(44, "pattern line base point y", &base[1])
		r.ConsumeFloat(45, "pattern line offset x", &offset[0])
		r.ConsumeFloat(46, "pattern line offset y", &offset[1])

		dashLengths := uint64(0)
		r.ConsumeNumber(79, DEC_RADIX, "number of dash length items", &dashLengths)

		for j := uint64(0); j < dashLengths; j++ {
			r.ConsumeFloat(49, "dash length", &dashLen)
			dashes = append(dashes, dashLen)
		}

		hatch.AppendPatternLine(angle, base, offset, dashes)
	}

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

	return r.Err()
}

func ParseAcDbEllipse(r *Reader, ellipse *entity.Ellipse) error {
	if r.AssertNextLine("AcDbEllipse") != nil {
		return r.Err()
	}

	r.ConsumeCoordinates(ellipse.Center[:])   // Center point
	r.ConsumeCoordinates(ellipse.EndPoint[:]) // Endpoint of major axis

	// XYZ extrusion direction
	// optional default = 0, 0, 1
	coord3D := [3]float64{0.0, 0.0, 0.0}
	r.ConsumeCoordinatesIf(210, coord3D[:])

	r.ConsumeFloat(40, "ratio of minor axis to major axis", &ellipse.Ratio)
	r.ConsumeFloat(41, "start parameter", &ellipse.Start)
	r.ConsumeFloat(42, "end parameter", &ellipse.End)

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

func ParseAcDbBlockReference(r *Reader, insert *entity.Insert) error {
	if r.AssertNextLine("AcDbBlockReference") != nil {
		return r.Err()
	}

	r.ConsumeNumberIf(66, DEC_RADIX, "attributes follow", &insert.AttributesFollow)
	r.ConsumeStr(&insert.BlockName)
	r.ConsumeCoordinates(insert.Coordinates[:])

	r.ConsumeFloatIf(41, "x scale factor", &insert.Scale[0])
	r.ConsumeFloatIf(42, "y scale factor", &insert.Scale[1])
	r.ConsumeFloatIf(43, "z scale factor", &insert.Scale[2])

	r.ConsumeFloatIf(50, "rotation angle", &insert.Rotation)
	r.ConsumeNumberIf(70, DEC_RADIX, "column count", &insert.ColCount)
	r.ConsumeNumberIf(71, DEC_RADIX, "row count", &insert.RowCount)

	r.ConsumeFloatIf(44, "column spacing", &insert.ColSpacing)
	r.ConsumeFloatIf(45, "row spacing", &insert.RowSpacing)

	// optional default = 0, 0, 1
	// XYZ extrusion direction
	coord3D := [3]float64{0.0, 0.0, 0.0}
	r.ConsumeCoordinatesIf(210, coord3D[:])

	return r.Err()
}

func ParseAcDbBlockBegin(r *Reader, block *blocks.Block) error {
	if r.AssertNextLine("AcDbBlockBegin") != nil {
		return r.Err()
	}

	r.ConsumeStr(&block.BlockName) // [2] block name
	r.ConsumeNumber(70, DEC_RADIX, "block-type flag", &block.Flag)
	r.ConsumeCoordinates(block.Coordinates[:])

	r.ConsumeStr(&block.OtherName) // [3] block name
	r.ConsumeStr(&block.XRefPath)  // [1] Xref path name

	return r.Err()
}

func ParseAcDbAttribute(r *Reader, attrib *entity.Attrib) error {
	if r.AssertNextLine("AcDbAttribute") != nil {
		return r.Err()
	}

	r.ConsumeStr(&attrib.Tag) // [2] Attribute tag
	r.ConsumeNumber(70, DEC_RADIX, "attribute flags", &attrib.Flags)
	r.ConsumeNumberIf(74, DEC_RADIX, "vertical text justification", &attrib.Text.VJustification) // group code 73 TEXT
	r.ConsumeNumberIf(280, DEC_RADIX, "version number", nil)

	r.ConsumeNumberIf(73, DEC_RADIX, "field length", nil) // not currently used
	r.ConsumeFloatIf(50, "text rotation", &attrib.Text.Rotation)
	r.ConsumeFloatIf(41, "relative x scale factor (width)", &attrib.Text.XScale) // adjusted when fit-type text is used
	r.ConsumeFloatIf(51, "oblique angle", &attrib.Text.Oblique)
	r.ConsumeStrIf(7, &attrib.Text.Style) // text style name default STANDARD
	r.ConsumeNumberIf(71, DEC_RADIX, "text generation flags", &attrib.Text.Flags)
	r.ConsumeNumberIf(72, DEC_RADIX, "horizontal text justification", &attrib.Text.HJustification)

	r.ConsumeCoordinatesIf(11, attrib.Text.Vector[:])
	r.ConsumeCoordinatesIf(210, attrib.Text.Vector[:])

	// TODO: parse XDATA
	code, err := r.PeekCode()
	for code != 0 && err == nil {
		r.ConsumeStr(nil)
		code, err = r.PeekCode()
	}
	if err != nil {
		return err
	}

	return r.Err()
}

func ParseAcDbAttributeDefinition(r *Reader, attdef *entity.Attdef) error {
	if r.AssertNextLine("AcDbAttributeDefinition") != nil {
		return r.Err()
	}

	r.ConsumeStr(&attdef.Prompt) // [3] prompt string
	r.ConsumeStr(&attdef.Tag)    // [2] tag string
	r.ConsumeNumber(70, DEC_RADIX, "attribute flags", &attdef.Flags)
	r.ConsumeFloatIf(73, "field length", nil)
	r.ConsumeNumberIf(74, DEC_RADIX, "vertical text justification", &attdef.Text.VJustification)

	r.ConsumeNumber(280, DEC_RADIX, "lock position flag", nil)

	r.ConsumeNumberIf(71, DEC_RADIX, "attachment point", &attdef.AttachmentPoint)
	r.ConsumeNumberIf(72, DEC_RADIX, "drawing direction", &attdef.DrawingDirection)

	r.ConsumeCoordinatesIf(11, attdef.Direction[:])

	return HelperParseEmbeddedObject(r)
}
