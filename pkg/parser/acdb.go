package parser

import (
	"log"

	"github.com/aichingert/dxf/pkg/blocks"
	"github.com/aichingert/dxf/pkg/entity"
)

func ParseAcDbEntity(r *Reader, entity entity.Entity) error {
	r.ConsumeNumber(5, HexRadix, "handle", entity.GetHandle())

	// TODO: set hard owner/handle to owner dictionary
	if r.ConsumeStrIf(102, nil) { // consumeIf => ex. {ACAD_XDICTIONARY
		r.ConsumeStr(nil) // 360 => hard owner
		for r.ConsumeNumberIf(330, HexRadix, "soft owner", nil) {
		}
		r.ConsumeStr(nil) // 102 }
	}

	if r.ConsumeStrIf(102, nil) { // consumeIf => ex. {ACAD_XDICTIONARY
		r.ConsumeStr(nil) // 360 => hard owner
		for r.ConsumeNumberIf(330, HexRadix, "soft owner", nil) {
		}
		r.ConsumeStr(nil) // 102 }
	}

	r.ConsumeNumber(330, HexRadix, "owner ptr", entity.GetOwner())

	if r.AssertNextLine("AcDbEntity") != nil {
		return r.Err()
	}

	// TODO: think about paper space visibility
	r.ConsumeStrIf(67, nil)
	r.ConsumeStr(entity.GetLayerName())

	r.ConsumeStrIf(6, nil) // ByBlock
	r.ConsumeNumberIf(62, DecRadix, "color number (present if not bylayer)", nil)
	r.ConsumeFloatIf(48, "linetype scale", nil)
	r.ConsumeNumberIf(60, DecRadix, "object visibility", entity.GetVisibility())

	r.ConsumeNumberIf(420, DecRadix, "24-bit color value", nil)
	r.ConsumeNumberIf(440, DecRadix, "transparency value", nil)
	r.ConsumeNumberIf(370, DecRadix, "not documented", nil)

	return r.Err()
}

func ParseAcDbLine(r *Reader, line *entity.Line) error {
	if r.AssertNextLine("AcDbLine") != nil {
		return r.Err()
	}

	r.ConsumeFloatIf(39, "thickness", nil)
	r.ConsumeCoordinates(line.Src[:])
	r.ConsumeCoordinates(line.Dst[:])

	return r.Err()
}

func ParseAcDbPolyline(r *Reader, polyline *entity.Polyline) error {
	if r.AssertNextLine("AcDbPolyline") != nil {
		return r.Err()
	}

	r.ConsumeNumber(90, DecRadix, "number of vertices", &polyline.Vertices)
	r.ConsumeNumber(70, DecRadix, "polyline flag", &polyline.Flag)

	if !r.ConsumeFloatIf(43, "line width for each vertex", nil) {
		//r.ConsumeFloat(43, "", nil)
		//log.Fatal("[ENTITIES(", Line, ")] TODO: implement line width for each vertex")
	}

	for i := int64(0); i < polyline.Vertices; i++ {
		bulge := 0.0
		coords2D := [2]float64{0.0, 0.0}

		r.ConsumeCoordinates(coords2D[:])
		r.ConsumeFloatIf(42, "expected bulge", &bulge)
		r.ConsumeNumberIf(91, DecRadix, "vertex identifier", nil)

		if r.Err() != nil {
			return r.Err()
		}

		polyline.AppendPLine(coords2D, bulge)
	}

	return r.Err()
}

func ParseAcDb2dPolyline(r *Reader, _ *entity.Polyline) error {
	if r.AssertNextLine("AcDb2dPolyline") != nil {
		return r.Err()
	}

	r.ConsumeNumberIf(66, DecRadix, "obsolete", nil)

	coords3D := [3]float64{0.0, 0.0, 0.0}
	r.ConsumeCoordinates(coords3D[:])

	r.ConsumeFloatIf(39, "thickness", nil)
	r.ConsumeNumberIf(70, DecRadix, "polyline flag", nil)

	r.ConsumeFloatIf(40, "start width default 0", nil)
	r.ConsumeFloatIf(41, "end width default 0", nil)
	r.ConsumeFloatIf(71, "mesh M vertex count", nil)
	r.ConsumeFloatIf(72, "mesh N vertex count", nil)
	r.ConsumeFloatIf(73, "smooth surface M density", nil)
	r.ConsumeFloatIf(74, "smooth surface N density", nil)
	r.ConsumeNumberIf(75, DecRadix, "curves and smooth surface default 0", nil)

	r.ConsumeCoordinatesIf(210, coords3D[:])

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

	r.ConsumeNumberIf(71, DecRadix, "text generation flags default 0", &text.Flags)
	r.ConsumeNumberIf(72, DecRadix, "horizontal text justification", &text.HJustification)

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

	r.ConsumeNumberIf(73, DecRadix, "vertical text justification", &text.VJustification)

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

	r.ConsumeNumber(71, DecRadix, "attachment point", &mText.Layout)
	r.ConsumeNumber(72, DecRadix, "direction (ex: left to right)", &mText.Direction)

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

	r.ConsumeNumber(73, DecRadix, "line spacing", &mText.LineSpacing)
	r.ConsumeFloat(44, "line spacing factor", nil)

	return HelperParseEmbeddedObject(r)
}

func HelperParseEmbeddedObject(r *Reader) error {
	// Embedded Object
	if r.ConsumeStrIf(101, nil) {
		r.ConsumeNumberIf(70, DecRadix, "not documented", nil)
		coords3D := [3]float64{0.0, 0.0, 0.0}
		r.ConsumeCoordinates(coords3D[:])
		r.ConsumeCoordinatesIf(11, coords3D[:])

		r.ConsumeFloatIf(40, "not documented", nil)
		r.ConsumeFloatIf(41, "not documented", nil)
		r.ConsumeFloatIf(42, "not documented", nil)
		r.ConsumeFloatIf(43, "not documented", nil)
		r.ConsumeFloatIf(46, "not documented", nil)

		r.ConsumeNumberIf(71, DecRadix, "not documented", nil)
		r.ConsumeNumberIf(72, DecRadix, "not documented", nil)
		r.ConsumeStrIf(1, nil)

		r.ConsumeFloatIf(44, "not documented", nil)
		r.ConsumeFloatIf(45, "not documented", nil)

		r.ConsumeNumberIf(73, DecRadix, "not documented", nil)
		r.ConsumeNumberIf(74, DecRadix, "not documented", nil)

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
	r.ConsumeNumber(70, DecRadix, "solid fill flag", &hatch.SolidFill)
	r.ConsumeNumber(71, DecRadix, "associativity flag", &hatch.Associative)

	boundaryPaths := int64(0)
	r.ConsumeNumber(91, DecRadix, "boundary paths", &boundaryPaths)

	for i := int64(0); i < boundaryPaths; i++ {
		pathTypeFlag := int64(0)
		// [92] Boundary path type flag (bit coded):
		// 0 = Default | 1 = External | 2  = Polyline
		// 4 = Derived | 8 = Textbox  | 16 = Outermost
		r.ConsumeNumber(92, DecRadix, "boundary path type flag", &pathTypeFlag)

		if pathTypeFlag&2 == 2 {
			polyline := entity.NewPolyline()

			// maybe consider?
			r.ConsumeNumber(72, DecRadix, "has bulge flag", nil)
			r.ConsumeNumber(73, DecRadix, "is closed flag", nil)
			r.ConsumeNumber(93, DecRadix, "number of polyline vertices", &polyline.Vertices)

			for vertex := int64(0); vertex < polyline.Vertices; vertex++ {
				coord2D, bulge := [2]float64{0.0, 0.0}, 0.0
				r.ConsumeCoordinates(coord2D[:])
				r.ConsumeFloatIf(42, "expected bulge", &bulge)
				polyline.AppendPLine(coord2D, bulge)
			}

			hatch.Polylines = append(hatch.Polylines, polyline)
		} else {
			edges, edgeType := int64(0), int64(0)

			r.ConsumeNumber(93, DecRadix, "number of edges in this boundary path", &edges)

			for edge := int64(0); edge < edges; edge++ {
				r.ConsumeNumber(72, DecRadix, "edge type data", &edgeType)

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
					r.ConsumeNumber(73, DecRadix, "is counterclockwise", &arc.Counterclockwise)
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

		boundaryObjectSize, boundaryObjectRef := int64(0), int64(0)

		r.ConsumeNumber(97, DecRadix, "number of source boundary objects", &boundaryObjectSize)
		for i := int64(0); i < boundaryObjectSize; i++ {
			r.ConsumeNumber(330, HexRadix, "reference to source object", &boundaryObjectRef)
		}
	}

	r.ConsumeNumber(75, DecRadix, "hatch style", &hatch.Style)
	r.ConsumeNumber(76, DecRadix, "hatch pattern type", &hatch.Pattern)
	r.ConsumeFloatIf(52, "hatch pattern angle", &hatch.Angle)
	r.ConsumeFloatIf(41, "hatch pattern scale or spacing", &hatch.Scale)
	r.ConsumeNumberIf(77, DecRadix, "hatch pattern double flag", &hatch.Double)

	patternDefinitions := int64(0)

	r.ConsumeNumberIf(78, DecRadix, "number of pattern definition lines", &patternDefinitions)

	for i := int64(0); i < patternDefinitions; i++ {
		base, offset, angle := [2]float64{0.0, 0.0}, [2]float64{0.0, 0.0}, 0.0
		dashes, dashLen := []float64{}, 0.0

		r.ConsumeFloat(53, "pattern line angle", &angle)
		r.ConsumeFloat(43, "pattern line base point x", &base[0])
		r.ConsumeFloat(44, "pattern line base point y", &base[1])
		r.ConsumeFloat(45, "pattern line offset x", &offset[0])
		r.ConsumeFloat(46, "pattern line offset y", &offset[1])

		dashLengths := int64(0)
		r.ConsumeNumber(79, DecRadix, "number of dash length items", &dashLengths)

		for j := int64(0); j < dashLengths; j++ {
			r.ConsumeFloat(49, "dash length", &dashLen)
			dashes = append(dashes, dashLen)
		}

		hatch.AppendPatternLine(angle, base, offset, dashes)
	}

	r.ConsumeFloatIf(47, "pixel size used to determine density to perform ray casting", nil)

	seedPoints := int64(0)
	r.ConsumeNumber(98, DecRadix, "number of seed points", &seedPoints)

	coord2D := [2]float64{0.0, 0.0}

	for seedPoint := int64(0); seedPoint < seedPoints; seedPoint++ {
		r.ConsumeCoordinates(coord2D[:])
	}

	r.ConsumeNumberIf(450, DecRadix, "indicates solid hatch or gradient", nil)
	r.ConsumeNumberIf(451, DecRadix, "zero is reserved for future use", nil)

	// default 0,0
	r.ConsumeFloatIf(460, "rotation angle in radians for gradients", nil)
	r.ConsumeFloatIf(461, "gradient definition", nil)
	r.ConsumeNumberIf(452, DecRadix, "records how colors were defined", nil)
	r.ConsumeFloatIf(462, "color tint value used by dialog", nil)

	nColors := int64(0)
	r.ConsumeNumberIf(453, DecRadix, "number of colors", &nColors)

	for color := int64(0); color < nColors; color++ {
		r.ConsumeFloatIf(463, "reserved for future use", nil)
		r.ConsumeNumberIf(63, DecRadix, "not documented", nil)
		r.ConsumeNumberIf(421, DecRadix, "not documented", nil)
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

// AcDbPoint
func ParseAcDbTrace(r *Reader, _ *entity.MText) error {
	if r.AssertNextLine("AcDbTrace") != nil {
		return r.Err()
	}

	coord3D := [3]float64{0.0, 0.0, 0.0}

	r.ConsumeCoordinates(coord3D[:])
	r.ConsumeCoordinates(coord3D[:])
	r.ConsumeFloat(12, "", nil)
	r.ConsumeFloat(22, "", nil)
	r.ConsumeFloat(32, "", nil)
	r.ConsumeFloat(13, "", nil)
	r.ConsumeFloat(23, "", nil)
	r.ConsumeFloat(33, "", nil)

	r.ConsumeNumberIf(39, DecRadix, "thickness", nil)

	// XYZ extrusion direction
	// optional default 0, 0, 1
	r.ConsumeCoordinatesIf(210, coord3D[:])
	r.ConsumeFloatIf(50, "angle of the x axis", nil)

	return r.Err()
}

// TODO: implement entity entity.Vertex
func ParseAcDbVertex(r *Reader, _ *entity.MText) error {
	if r.AssertNextLine("AcDbVertex") != nil {
		return r.Err()
	}

	next := ""
	r.ConsumeStr(&next) // AcDb2dVertex or AcDb3dPolylineVertex

	coord3D := [3]float64{0.0, 0.0, 0.0}

	r.ConsumeCoordinates(coord3D[:])
	r.ConsumeFloatIf(40, "starting width", nil)
	r.ConsumeFloatIf(41, "end width", nil)
	r.ConsumeFloatIf(42, "bulge", nil)

	r.ConsumeNumberIf(70, DecRadix, "vertex flags", nil)
	r.ConsumeFloatIf(50, "curve fit tangent direction", nil)

	r.ConsumeFloatIf(71, "polyface mesh vertex index", nil)
	r.ConsumeFloatIf(72, "polyface mesh vertex index", nil)
	r.ConsumeFloatIf(73, "polyface mesh vertex index", nil)
	r.ConsumeFloatIf(74, "polyface mesh vertex index", nil)

	r.ConsumeNumberIf(91, DecRadix, "vertex identifier", nil)

	return r.Err()
}

// TODO: implement entity entity.Point
func ParseAcDbPoint(r *Reader, _ *entity.MText) error {
	if r.AssertNextLine("AcDbPoint") != nil {
		return r.Err()
	}

	coord3D := [3]float64{0.0, 0.0, 0.0}

	r.ConsumeCoordinates(coord3D[:])
	r.ConsumeNumberIf(39, DecRadix, "thickness", nil)

	// XYZ extrusion direction
	// optional default 0, 0, 1
	r.ConsumeCoordinatesIf(210, coord3D[:])
	r.ConsumeFloatIf(50, "angle of the x axis", nil)

	return r.Err()
}

func ParseAcDbBlockReference(r *Reader, insert *entity.Insert) error {
	if r.AssertNextLine("AcDbBlockReference") != nil {
		return r.Err()
	}

	r.ConsumeNumberIf(66, DecRadix, "attributes follow", &insert.AttributesFollow)
	r.ConsumeStr(&insert.BlockName)
	r.ConsumeCoordinates(insert.Coordinates[:])

	r.ConsumeFloatIf(41, "x scale factor", &insert.Scale[0])
	r.ConsumeFloatIf(42, "y scale factor", &insert.Scale[1])
	r.ConsumeFloatIf(43, "z scale factor", &insert.Scale[2])

	r.ConsumeFloatIf(50, "rotation angle", &insert.Rotation)
	r.ConsumeNumberIf(70, DecRadix, "column count", &insert.ColCount)
	r.ConsumeNumberIf(71, DecRadix, "row count", &insert.RowCount)

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
	r.ConsumeNumber(70, DecRadix, "block-type flag", &block.Flag)
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
	r.ConsumeNumber(70, DecRadix, "attribute flags", &attrib.Flags)
	r.ConsumeNumberIf(74, DecRadix, "vertical text justification", &attrib.Text.VJustification) // group code 73 TEXT
	r.ConsumeNumberIf(280, DecRadix, "version number", nil)

	r.ConsumeNumberIf(73, DecRadix, "field length", nil) // not currently used
	r.ConsumeFloatIf(50, "text rotation", &attrib.Text.Rotation)
	r.ConsumeFloatIf(41, "relative x scale factor (width)", &attrib.Text.XScale) // adjusted when fit-type text is used
	r.ConsumeFloatIf(51, "oblique angle", &attrib.Text.Oblique)
	r.ConsumeStrIf(7, &attrib.Text.Style) // text style name default STANDARD
	r.ConsumeNumberIf(71, DecRadix, "text generation flags", &attrib.Text.Flags)
	r.ConsumeNumberIf(72, DecRadix, "horizontal text justification", &attrib.Text.HJustification)

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
	r.ConsumeNumber(70, DecRadix, "attribute flags", &attdef.Flags)
	r.ConsumeFloatIf(73, "field length", nil)
	r.ConsumeNumberIf(74, DecRadix, "vertical text justification", &attdef.Text.VJustification)

	r.ConsumeNumber(280, DecRadix, "lock position flag", nil)

	r.ConsumeNumberIf(71, DecRadix, "attachment point", &attdef.AttachmentPoint)
	r.ConsumeNumberIf(72, DecRadix, "drawing direction", &attdef.DrawingDirection)

	r.ConsumeCoordinatesIf(11, attdef.Direction[:])

	return HelperParseEmbeddedObject(r)
}

func ParseAcDbDimension(r *Reader, _ *entity.Attdef) error {
	if r.AssertNextLine("AcDbDimension") != nil {
		return r.Err()
	}

	r.ConsumeNumber(280, DecRadix, "version number", nil)
	r.ConsumeStr(nil) // name of the block

	coords3D := [3]float64{0.0, 0.0, 0.0}

	r.ConsumeCoordinates(coords3D[:])
	r.ConsumeCoordinates(coords3D[:])

	r.ConsumeNumber(70, DecRadix, "dimension type", nil)
	r.ConsumeNumber(71, DecRadix, "attachment point", nil)
	r.ConsumeNumberIf(72, DecRadix, "dimension text-line spacing", nil)

	r.ConsumeNumberIf(41, DecRadix, "dimension text-line factor", nil)
	r.ConsumeNumberIf(42, DecRadix, "actual measurement", nil)

	r.ConsumeNumberIf(73, DecRadix, "not documented", nil)
	r.ConsumeNumberIf(74, DecRadix, "not documented", nil)
	r.ConsumeNumberIf(75, DecRadix, "not documented", nil)

	r.ConsumeStrIf(1, nil) // dimension text
	r.ConsumeFloatIf(53, "roation angle of the dimension", nil)
	r.ConsumeFloatIf(51, "horizontal direction", nil)

	r.ConsumeCoordinatesIf(210, coords3D[:])
	r.ConsumeStrIf(3, nil) // [3] dimension style name

	dim := ""
	r.ConsumeStr(&dim)

	switch dim {
	// should be acdb3pointangulardimension
	case "AcDb2LineAngularDimension":
		r.ConsumeFloat(13, "point for linear and angular dimension", nil)
		r.ConsumeFloat(23, "point for linear and angular dimension", nil)
		r.ConsumeFloat(33, "point for linear and angular dimension", nil)
		r.ConsumeFloat(14, "point for linear and angular dimension", nil)
		r.ConsumeFloat(24, "point for linear and angular dimension", nil)
		r.ConsumeFloat(34, "point for linear and angular dimension", nil)
		r.ConsumeFloat(15, "point for diameter, radius, and angular dimension", nil)
		r.ConsumeFloat(25, "point for diameter, radius, and angular dimension", nil)
		r.ConsumeFloat(35, "point for diameter, radius, and angular dimension", nil)
	case "AcDbAlignedDimension":
		r.ConsumeFloatIf(12, "insertion point for clones of a dimension", nil)
		r.ConsumeFloatIf(22, "insertion point for clones of a dimension", nil)
		r.ConsumeFloatIf(32, "insertion point for clones of a dimension", nil)
		r.ConsumeFloat(13, "definition point for linear and angular dimensions", nil)
		r.ConsumeFloat(23, "definition point for linear and angular dimensions", nil)
		r.ConsumeFloat(33, "definition point for linear and angular dimensions", nil)
		r.ConsumeFloat(14, "definition point for linear and angular dimensions", nil)
		r.ConsumeFloat(24, "definition point for linear and angular dimensions", nil)
		r.ConsumeFloat(34, "definition point for linear and angular dimensions", nil)
		r.ConsumeFloatIf(50, "angle of rotated, horizontal, or vertical dimensions", nil)
		r.ConsumeFloatIf(52, "oblique angle", nil)
		r.ConsumeStrIf(100, nil) // subclass marker AcDbRotatedDimension
	default:
		log.Fatal(dim)
	}

	return r.Err()
}

func ParseAcDbViewport(r *Reader, _ *entity.MText) error {
	if r.AssertNextLine("AcDbViewport") != nil {
		return r.Err()
	}

	coords3D := [3]float64{0, 0, 0}
	r.ConsumeCoordinates(coords3D[:])

	r.ConsumeFloat(40, "width in paper space units", nil)
	r.ConsumeFloat(41, "height in paper space units", nil)

	r.ConsumeFloat(68, "viewport status field", nil)
	// => -1 0 On, 0 = Off

	r.ConsumeNumber(69, DecRadix, "viewport id", nil)

	r.ConsumeFloat(12, "center point x", nil)
	r.ConsumeFloat(22, "center point y", nil)
	r.ConsumeFloat(13, "snap base point x", nil)
	r.ConsumeFloat(23, "snap base point y", nil)
	r.ConsumeFloat(14, "snap spacing x", nil)
	r.ConsumeFloat(24, "snap spacing y", nil)
	r.ConsumeFloat(15, "grid spacing x", nil)
	r.ConsumeFloat(25, "grid spacing y", nil)

	r.ConsumeFloat(16, "view direction vector x", nil)
	r.ConsumeFloat(26, "view direction vector y", nil)
	r.ConsumeFloat(36, "view direction vector z", nil)

	r.ConsumeFloat(17, "view target point y", nil)
	r.ConsumeFloat(27, "view target point x", nil)
	r.ConsumeFloat(37, "view target point z", nil)

	r.ConsumeFloat(42, "perspective lens length", nil)
	r.ConsumeFloat(43, "front clip plane z value", nil)
	r.ConsumeFloat(44, "back clip plane z value", nil)
	r.ConsumeFloat(45, "view height", nil)
	r.ConsumeFloat(50, "snap angle", nil)
	r.ConsumeFloat(51, "view twist angle", nil)

	r.ConsumeFloat(72, "circle zoom percent", nil)

	code, err := r.PeekCode()
	for err != nil && code == 331 {
		r.ConsumeNumber(331, DecRadix, "frozen layer object Id/handle", nil)
	}

	r.ConsumeNumber(90, HexRadix, "viewport status bit-coded flags", nil)
	r.ConsumeNumberIf(340, DecRadix, "hard-pointer id/handle to entity that serves as the viewports clipping boundary", nil)
	r.ConsumeStr(nil) // [1]
	r.ConsumeNumber(281, DecRadix, "render mode", nil)
	r.ConsumeNumber(71, DecRadix, "ucs per viewport flag", nil)
	r.ConsumeNumber(74, DecRadix, "display ucs icon at ucs origin flag", nil)
	r.ConsumeFloat(110, "ucs origin x", nil)
	r.ConsumeFloat(120, "ucs origin y", nil)
	r.ConsumeFloat(130, "ucs origin z", nil)

	r.ConsumeFloat(111, "ucs x-axis x", nil)
	r.ConsumeFloat(121, "ucs x-axis y", nil)
	r.ConsumeFloat(131, "ucs x-axis z", nil)

	r.ConsumeFloat(112, "ucs y-axis x", nil)
	r.ConsumeFloat(122, "ucs y-axis y", nil)
	r.ConsumeFloat(132, "ucs y-axis z", nil)

	r.ConsumeNumberIf(345, DecRadix, "id/handle of AcDbUCSTableRecord if UCS is a named ucs", nil)
	r.ConsumeNumberIf(346, DecRadix, "id/handle of AcDbUCSTableRecord of base ucs", nil)
	r.ConsumeNumber(79, DecRadix, "Orthographic type of UCS", nil)
	r.ConsumeFloat(146, "elevation", nil)
	r.ConsumeNumber(170, DecRadix, "ShadePlot mode", nil)
	r.ConsumeNumber(61, DecRadix, "frequency of major grid lines compared to minor grid lines", nil)

	r.ConsumeNumberIf(332, DecRadix, "background id/handle", nil)
	r.ConsumeNumberIf(333, DecRadix, "shade plot id/handle", nil)
	r.ConsumeNumberIf(348, DecRadix, "visual style id/handle", nil)

	r.ConsumeNumber(292, DecRadix, "default lighting type on when no use lights are specified", nil)
	r.ConsumeNumber(282, DecRadix, "default lighting type", nil)
	r.ConsumeFloat(141, "view brightness", nil)
	r.ConsumeFloat(142, "view contrast", nil)

	r.ConsumeFloatIf(63, "ambient light color only if not black", nil)
	r.ConsumeFloatIf(421, "ambient light color only if not black", nil)
	r.ConsumeFloatIf(431, "ambient light color only if not black", nil)

	r.ConsumeNumberIf(361, DecRadix, "sun id/handle", nil)
	r.ConsumeNumberIf(335, DecRadix, "soft pointer reference to id/handle", nil)
	r.ConsumeNumberIf(343, DecRadix, "soft pointer reference to id/handle", nil)
	r.ConsumeNumberIf(344, DecRadix, "soft pointer reference to id/handle", nil)
	r.ConsumeNumberIf(91, DecRadix, "soft pointer reference to id/handle", nil)

	return r.Err()
}
