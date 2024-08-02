package parser

import (
	"log"

	"github.com/aichingert/dxf/pkg/blocks"
	"github.com/aichingert/dxf/pkg/entity"
)

// TODO: replace with actual entities values
var (
	coords2D = [2]float64{0, 0}
	coords3D = [3]float64{0, 0, 0}
)

func ParseAcDbEntity(r Reader, entity entity.Entity) {
	r.consumeNumber(5, HexRadix, "handle", entity.GetHandle())

	// TODO: set hard owner/handle to owner dictionary
	if r.consumeStrIf(102, nil) { // consumeIf => ex. {ACAD_XDICTIONARY
		r.consumeStr(nil) // 360 => hard owner
		for r.consumeNumberIf(330, HexRadix, "soft owner", nil) {
		}
		r.consumeStr(nil) // 102 }
	}

	if r.consumeStrIf(102, nil) { // consumeIf => ex. {ACAD_XDICTIONARY
		r.consumeStr(nil) // 360 => hard owner
		for r.consumeNumberIf(330, HexRadix, "soft owner", nil) {
		}
		r.consumeStr(nil) // 102 }
	}

	r.consumeNumber(330, HexRadix, "owner ptr", entity.GetOwner())

	if r.assertNextLine("AcDbEntity") != nil {
		return
	}

	// TODO: think about paper space visibility
	r.consumeStrIf(67, nil)
	r.consumeStr(entity.GetLayerName())

	r.consumeStrIf(6, nil) // ByBlock
	r.consumeNumberIf(62, DecRadix, "color number (present if not bylayer)", nil)
	r.consumeFloatIf(48, "linetype scale", nil)
	r.consumeNumberIf(60, DecRadix, "object visibility", entity.GetVisibility())

	r.consumeNumberIf(420, DecRadix, "24-bit color value", nil)
	r.consumeNumberIf(440, DecRadix, "transparency value", nil)
	r.consumeNumberIf(370, DecRadix, "not documented", nil)
}

func ParseAcDbLine(r Reader, line *entity.Line) {
	if r.assertNextLine("AcDbLine") != nil {
		return
	}

	r.consumeFloatIf(39, "thickness", nil)
	r.consumeCoordinates(line.Src[:])
	r.consumeCoordinates(line.Dst[:])
}

func ParseAcDbPolyline(r Reader, polyline *entity.Polyline) {
	if r.assertNextLine("AcDbPolyline") != nil {
		return
	}

	vertices := int64(0)
	r.consumeNumber(90, DecRadix, "number of vertices", &vertices)
	r.consumeNumber(70, DecRadix, "polyline flag", &polyline.Flag)
	r.consumeFloatIf(43, "line width for each vertex", nil)

	for i := int64(0); i < vertices; i++ {
		bulge := 0.0

		r.consumeCoordinates(coords2D[:])

		r.consumeFloatIf(40, "default start width", nil)
		r.consumeFloatIf(41, "default end width", nil)

		r.consumeFloatIf(42, "expected bulge", &bulge)
		r.consumeNumberIf(91, DecRadix, "vertex identifier", nil)

		if r.Err() != nil {
			return
		}

		polyline.AppendPLine(coords2D, bulge)
	}
}

func ParseAcDb2dPolyline(r Reader, _ *entity.Polyline) {
	if r.assertNextLine("AcDb2dPolyline") != nil {
		return
	}

	r.consumeNumberIf(66, DecRadix, "obsolete", nil)
	r.consumeCoordinates(coords3D[:])
	r.consumeFloatIf(39, "thickness", nil)
	r.consumeNumberIf(70, DecRadix, "polyline flag", nil)

	r.consumeFloatIf(40, "start width default 0", nil)
	r.consumeFloatIf(41, "end width default 0", nil)
	r.consumeFloatIf(71, "mesh M vertex count", nil)
	r.consumeFloatIf(72, "mesh N vertex count", nil)
	r.consumeFloatIf(73, "smooth surface M density", nil)
	r.consumeFloatIf(74, "smooth surface N density", nil)
	r.consumeNumberIf(75, DecRadix, "curves and smooth surface default 0", nil)

	r.consumeCoordinatesIf(210, coords3D[:])
}

func ParseAcDbCircle(r Reader, circle *entity.Circle) {
	if r.assertNextLine("AcDbCircle") != nil {
		return
	}

	r.consumeFloatIf(39, "thickness", nil)
	r.consumeCoordinates(circle.Coordinates[:])
	r.consumeFloat(40, "expected radius", &circle.Radius)
}

func ParseAcDbArc(r Reader, arc *entity.Arc) {
	if r.assertNextLine("AcDbArc") != nil {
		return
	}

	r.consumeFloat(50, "expected startAngle", &arc.StartAngle)
	r.consumeFloat(51, "expected endAngle", &arc.EndAngle)
}

func ParseAcDbText(r Reader, text *entity.Text) {
	if r.assertNextLine("AcDbText") != nil {
		return
	}

	r.consumeFloatIf(39, "expected thickness", &text.Thickness)
	r.consumeCoordinates(text.Coordinates[:])

	r.consumeFloat(40, "expected text height", &text.Height)
	r.consumeStr(&text.Text) // [1] default value of the string itself

	r.consumeFloatIf(50, "text rotation default 0", &text.Rotation)
	r.consumeFloatIf(41, "relative x scale factor default 1", &text.XScale)
	r.consumeFloatIf(51, "oblique angle default 0", &text.Oblique)

	r.consumeStrIf(7, &text.Style) // text style name default STANDARD

	r.consumeNumberIf(71, DecRadix, "text generation flags default 0", &text.Flags)
	r.consumeNumberIf(72, DecRadix, "horizontal text justification", &text.HJustification)

	r.consumeCoordinatesIf(11, text.Vector[:])
	r.consumeCoordinatesIf(210, text.Vector[:])

	if r.line() == "AcDbText" {
		r.consumeNext()
	}

	// group 72 and 73 integer codes
	// https://help.autodesk.com/view/OARX/2024/ENU/?guid=GUID-62E5383D-8A14-47B4-BFC4-35824CAE8363

	r.consumeNumberIf(73, DecRadix, "vertical text justification", &text.VJustification)
}

func ParseAcDbMText(r Reader, mText *entity.MText) {
	if r.assertNextLine("AcDbMText") != nil {
		return
	}

	r.consumeCoordinates(mText.Coordinates[:])
	r.consumeFloat(40, "expected text height", &mText.TextHeight)

	// TODO: https://ezdxf.readthedocs.io/en/stable/dxfinternals/entities/mtext.html
	r.consumeFloat(41, "rectangle width", nil)
	r.consumeFloat(46, "column height", nil)

	r.consumeNumber(71, DecRadix, "attachment point", &mText.Layout)
	r.consumeNumber(72, DecRadix, "direction (ex: left to right)", &mText.Direction)

	for r.code() == 1 || r.code() == 3 {
		line := r.consumeNext()

		if r.Err() != nil {
			return
		}

		mText.Text = append(mText.Text, line)
	}

	r.consumeStrIf(7, &mText.TextStyle)
	r.consumeCoordinatesIf(11, mText.Vector[:])

	r.consumeNumber(73, DecRadix, "line spacing", &mText.LineSpacing)
	r.consumeFloat(44, "line spacing factor", nil)
	HelperParseEmbeddedObject(r)
}

func HelperParseEmbeddedObject(r Reader) {
	// Embedded Object
	if r.consumeStrIf(101, nil) {
		r.consumeNumberIf(70, DecRadix, "not documented", nil)
		r.consumeCoordinates(coords3D[:])
		r.consumeCoordinatesIf(11, coords3D[:])

		r.consumeFloatIf(40, "not documented", nil)
		r.consumeFloatIf(41, "not documented", nil)
		r.consumeFloatIf(42, "not documented", nil)
		r.consumeFloatIf(43, "not documented", nil)
		r.consumeFloatIf(46, "not documented", nil)

		r.consumeNumberIf(71, DecRadix, "not documented", nil)
		r.consumeNumberIf(72, DecRadix, "not documented", nil)
		r.consumeStrIf(1, nil)

		r.consumeFloatIf(44, "not documented", nil)
		r.consumeFloatIf(45, "not documented", nil)

		r.consumeNumberIf(73, DecRadix, "not documented", nil)
		r.consumeNumberIf(74, DecRadix, "not documented", nil)

		r.consumeFloatIf(44, "not documented", nil)
		r.consumeFloatIf(46, "not documented", nil)
	}
}

func ParseAcDbHatch(r Reader, hatch *entity.Hatch) {
	if r.assertNextLine("AcDbHatch") != nil {
		return
	}

	r.consumeCoordinates(coords3D[:]) // elevation
	r.consumeCoordinates(coords3D[:])

	r.consumeStr(&hatch.PatternName)
	r.consumeNumber(70, DecRadix, "solid fill flag", &hatch.SolidFill)
	r.consumeNumber(71, DecRadix, "associativity flag", &hatch.Associative)

	boundaryPaths := int64(0)
	r.consumeNumber(91, DecRadix, "boundary paths", &boundaryPaths)
	for i := int64(0); i < boundaryPaths; i++ {
		hatch.BoundaryPaths = append(hatch.BoundaryPaths, ParseBoundaryPath(r))
	}

	r.consumeNumber(75, DecRadix, "hatch style", &hatch.Style)
	r.consumeNumber(76, DecRadix, "hatch pattern type", &hatch.Pattern)
	r.consumeFloatIf(52, "hatch pattern angle", &hatch.Angle)
	r.consumeFloatIf(41, "hatch pattern scale or spacing", &hatch.Scale)
	r.consumeNumberIf(77, DecRadix, "hatch pattern double flag", &hatch.Double)

	patternDefinitions := int64(0)

	r.consumeNumberIf(78, DecRadix, "number of pattern definition lines", &patternDefinitions)

	for i := int64(0); i < patternDefinitions; i++ {
		base, offset, angle := [2]float64{0.0, 0.0}, [2]float64{0.0, 0.0}, 0.0
		var dashes []float64
		dashLen := 0.0

		r.consumeFloat(53, "pattern line angle", &angle)
		r.consumeFloat(43, "pattern line base point x", &base[0])
		r.consumeFloat(44, "pattern line base point y", &base[1])
		r.consumeFloat(45, "pattern line offset x", &offset[0])
		r.consumeFloat(46, "pattern line offset y", &offset[1])

		dashLengths := int64(0)
		r.consumeNumber(79, DecRadix, "number of dash length items", &dashLengths)

		for j := int64(0); j < dashLengths; j++ {
			r.consumeFloat(49, "dash length", &dashLen)
			dashes = append(dashes, dashLen)
		}

		hatch.AppendPatternLine(angle, base, offset, dashes)
	}

	r.consumeFloatIf(47, "pixel size used to determine the density", &hatch.PixelSize)

	seedPoints, nColors := int64(0), int64(0)
	r.consumeNumber(98, DecRadix, "number of seed points", &seedPoints)

	for seedPoint := int64(0); seedPoint < seedPoints; seedPoint++ {
		r.consumeCoordinates(hatch.SeedPoint[:2])
	}

	r.consumeNumberIf(450, DecRadix, "indicates solid hatch or gradient", nil)
	r.consumeNumberIf(451, DecRadix, "zero is reserved for future use", nil)
	r.consumeFloatIf(460, "rotation angle in radians for gradients", nil)
	r.consumeFloatIf(461, "gradient definition", nil)
	r.consumeNumberIf(452, DecRadix, "records how colors were defined", nil)
	r.consumeFloatIf(462, "color tint value used by dialog", nil)
	r.consumeNumberIf(453, DecRadix, "number of colors", &nColors)
	for color := int64(0); color < nColors; color++ {
		r.consumeFloatIf(463, "reserved for future use", nil)
		r.consumeNumberIf(63, DecRadix, "not documented", nil)
		r.consumeNumberIf(421, DecRadix, "not documented", nil)
	}
	r.consumeStrIf(470, nil) // string default = LINEAR
}

func ParseBoundaryPath(r Reader) *entity.BoundaryPath {
	path := &entity.BoundaryPath{}

	// [92] Boundary path type flag (bit coded):
	// 0 = Default | 1 = External | 2  = Polyline
	// 4 = Derived | 8 = Textbox  | 16 = Outermost
	r.consumeNumber(92, DecRadix, "boundary path type flag", &path.Flag)

	if path.Flag&2 == 2 {
		path.Polyline = &entity.Polyline{}
		hasBulge, bulge, vertices := int64(0), 0.0, int64(0)

		r.consumeNumber(72, DecRadix, "has bulge flag", &hasBulge)
		r.consumeNumber(73, DecRadix, "is closed flag", &path.Polyline.Flag)
		r.consumeNumber(93, DecRadix, "number of polyline vertices", &vertices)

		for vertex := int64(0); vertex < vertices; vertex++ {
			r.consumeCoordinates(coords2D[:])
			if hasBulge == 1 {
				r.consumeFloat(42, "expected bulge", &bulge)
			}
			path.Polyline.AppendPLine(coords2D, bulge)
		}
	} else {
		edges, edgeType := int64(0), int64(0)

		r.consumeNumber(93, DecRadix, "number of edges in this boundary path", &edges)

		for edge := int64(0); edge < edges; edge++ {
			r.consumeNumber(72, DecRadix, "edge type data", &edgeType)

			switch edgeType {
			case 1: // Line
				line := entity.NewLine()
				line.Entity = nil
				r.consumeCoordinates(line.Src[:2])
				r.consumeCoordinates(line.Dst[:2])
				path.Lines = append(path.Lines, line)
			case 2: // Circular arc
				arc := entity.NewArc()
				arc.Entity = nil
				r.consumeCoordinates(arc.Circle.Coordinates[:2])
				r.consumeFloat(40, "radius", &arc.Circle.Radius)

				r.consumeFloat(50, "start angle", &arc.StartAngle)
				r.consumeFloat(51, "end angle", &arc.EndAngle)
				r.consumeNumber(73, DecRadix, "is counterclockwise", &arc.Counterclockwise)
				path.Arcs = append(path.Arcs, arc)
			case 3: // Elliptic arc
				ellipse := entity.NewEllipse()
				ellipse.Entity = nil

				r.consumeCoordinates(ellipse.Center[:2])
				r.consumeCoordinates(ellipse.EndPoint[:2])
				r.consumeFloat(40, "length of minor axis", &ellipse.Ratio)
				r.consumeFloat(50, "start angle", &ellipse.Start)
				r.consumeFloat(51, "end angle", &ellipse.End)
				r.consumeFloat(73, "is counterclockwise", nil)

				path.Ellipses = append(path.Ellipses, ellipse)
			case 4: // Spine
				log.Fatal("[AcDbHatch(", Line, ")] TODO: implement boundary path spline")
			default:
				r.setErr(NewParseError("invalid edge type data"))
				return path
			}
		}
	}

	boundaryObjectSize, boundaryObjectRef := int64(0), int64(0)
	r.consumeNumber(97, DecRadix, "number of source boundary objects", &boundaryObjectSize)
	for i := int64(0); i < boundaryObjectSize; i++ {
		r.consumeNumber(330, HexRadix, "reference to source object", &boundaryObjectRef)
	}

	return path
}

func ParseAcDbEllipse(r Reader, ellipse *entity.Ellipse) {
	if r.assertNextLine("AcDbEllipse") != nil {
		return
	}

	r.consumeCoordinates(ellipse.Center[:])   // Center point
	r.consumeCoordinates(ellipse.EndPoint[:]) // Endpoint of major axis

	r.consumeCoordinatesIf(210, coords3D[:])

	r.consumeFloat(40, "ratio of minor axis to major axis", &ellipse.Ratio)
	r.consumeFloat(41, "start parameter", &ellipse.Start)
	r.consumeFloat(42, "end parameter", &ellipse.End)
}

func ParseAcDbSpline(r Reader, _ *entity.MText) {
	if r.assertNextLine("AcDbSpline") != nil {
		return
	}

	knots, controlPoints, fitPoints := int64(0), int64(0), int64(0)

	r.consumeCoordinates(coords3D[:])
	r.consumeNumber(70, DecRadix, "spline flag", nil)
	r.consumeNumber(71, DecRadix, "degree of the spline curve", nil)
	r.consumeNumber(72, DecRadix, "number of knots", &knots)
	r.consumeNumber(73, DecRadix, "number of control points", &controlPoints)
	r.consumeNumber(74, DecRadix, "number of fit points", &fitPoints)
	r.consumeFloatIf(42, "knot tolerance default 0.0000001", nil)
	r.consumeFloatIf(43, "control point tolerance 0.0000001", nil)
	r.consumeFloatIf(44, "fit tolerance default 0.0000001", nil)

	for i := int64(0); i < knots; i++ {
		r.consumeFloat(40, "knot value", nil)
	}
	for i := int64(0); i < controlPoints; i++ {
		r.consumeCoordinates(coords3D[:]) // start tangent - may be omitted
	}
	for i := int64(0); i < fitPoints; i++ {
		r.consumeCoordinates(coords3D[:]) // end tangent   - may be omitted
	}
}

// ParseAcDbTrace implement AcDbPoint
func ParseAcDbTrace(r Reader, _ *entity.MText) {
	if r.assertNextLine("AcDbTrace") != nil {
		return
	}

	r.consumeCoordinates(coords3D[:])
	r.consumeCoordinates(coords3D[:])
	r.consumeCoordinates(coords3D[:])
	r.consumeCoordinates(coords3D[:])

	r.consumeNumberIf(39, DecRadix, "thickness", nil)
	r.consumeCoordinatesIf(210, coords3D[:])
	r.consumeFloatIf(50, "angle of the x axis", nil)
}

// ParseAcDbVertex implement entity entity.Vertex
func ParseAcDbVertex(r Reader, _ *entity.MText) {
	if r.assertNextLine("AcDbVertex") != nil {
		return
	}

	next := ""
	r.consumeStr(&next) // AcDb2dVertex or AcDb3dPolylineVertex

	r.consumeCoordinates(coords3D[:])
	r.consumeFloatIf(40, "starting width", nil)
	r.consumeFloatIf(41, "end width", nil)
	r.consumeFloatIf(42, "bulge", nil)

	r.consumeNumberIf(70, DecRadix, "vertex flags", nil)
	r.consumeFloatIf(50, "curve fit tangent direction", nil)

	r.consumeFloatIf(71, "polyface mesh vertex index", nil)
	r.consumeFloatIf(72, "polyface mesh vertex index", nil)
	r.consumeFloatIf(73, "polyface mesh vertex index", nil)
	r.consumeFloatIf(74, "polyface mesh vertex index", nil)

	r.consumeNumberIf(91, DecRadix, "vertex identifier", nil)
}

// ParseAcDbPoint implement entity entity.Point
func ParseAcDbPoint(r Reader, _ *entity.MText) {
	if r.assertNextLine("AcDbPoint") != nil {
		return
	}

	r.consumeCoordinates(coords3D[:])
	r.consumeNumberIf(39, DecRadix, "thickness", nil)

	// XYZ extrusion direction
	// optional default 0, 0, 1
	r.consumeCoordinatesIf(210, coords3D[:])
	r.consumeFloatIf(50, "angle of the x axis", nil)
}

func ParseAcDbBlockReference(r Reader, insert *entity.Insert) {
	line := ""
	r.consumeStr(&line)
	if r.Err() != nil || !(line == "AcDbBlockReference" || line == "AcDbMInsertBlock") {
		return
	}

	r.consumeNumberIf(66, DecRadix, "attributes follow", &insert.AttributesFollow)
	r.consumeStr(&insert.BlockName)
	r.consumeCoordinates(insert.Coordinates[:])

	r.consumeFloatIf(41, "x scale factor", &insert.Scale[0])
	r.consumeFloatIf(42, "y scale factor", &insert.Scale[1])
	r.consumeFloatIf(43, "z scale factor", &insert.Scale[2])

	r.consumeFloatIf(50, "rotation angle", &insert.Rotation)
	r.consumeNumberIf(70, DecRadix, "column count", &insert.ColCount)
	r.consumeNumberIf(71, DecRadix, "row count", &insert.RowCount)

	r.consumeFloatIf(44, "column spacing", &insert.ColSpacing)
	r.consumeFloatIf(45, "row spacing", &insert.RowSpacing)

	// optional default = 0, 0, 1
	// XYZ extrusion direction
	r.consumeCoordinatesIf(210, coords3D[:])
}

func ParseAcDbBlockBegin(r Reader, block *blocks.Block) {
	if r.assertNextLine("AcDbBlockBegin") != nil {
		return
	}

	r.consumeStr(&block.BlockName) // [2] block name
	r.consumeNumber(70, DecRadix, "block-type flag", &block.Flag)
	r.consumeCoordinates(block.Coordinates[:])

	r.consumeStr(&block.OtherName) // [3] block name
	r.consumeStr(&block.XRefPath)  // [1] Xref path name
}

func ParseAcDbAttribute(r Reader, attrib *entity.Attrib) {
	if r.assertNextLine("AcDbAttribute") != nil {
		return
	}

	r.consumeStr(&attrib.Tag) // [2] Attribute tag
	r.consumeNumber(70, DecRadix, "attribute flags", &attrib.Flags)
	r.consumeNumberIf(74, DecRadix, "vertical text justification", &attrib.Text.VJustification) // group code 73 TEXT
	r.consumeNumberIf(280, DecRadix, "version number", nil)

	r.consumeNumberIf(73, DecRadix, "field length", nil) // not currently used
	r.consumeFloatIf(50, "text rotation", &attrib.Text.Rotation)
	r.consumeFloatIf(41, "relative x scale factor (width)", &attrib.Text.XScale) // adjusted when fit-type text is used
	r.consumeFloatIf(51, "oblique angle", &attrib.Text.Oblique)
	r.consumeStrIf(7, &attrib.Text.Style) // text style name default STANDARD
	r.consumeNumberIf(71, DecRadix, "text generation flags", &attrib.Text.Flags)
	r.consumeNumberIf(72, DecRadix, "horizontal text justification", &attrib.Text.HJustification)

	r.consumeCoordinatesIf(11, attrib.Text.Vector[:])
	r.consumeCoordinatesIf(210, attrib.Text.Vector[:])

	// TODO: parse XDATA
	for r.code() != 0 {
		r.consumeNext()
	}
}

func ParseAcDbAttributeDefinition(r Reader, attdef *entity.Attdef) {
	if r.assertNextLine("AcDbAttributeDefinition") != nil {
		return
	}

	r.consumeStr(&attdef.Prompt) // [3] prompt string
	r.consumeStr(&attdef.Tag)    // [2] tag string
	r.consumeNumber(70, DecRadix, "attribute flags", &attdef.Flags)
	r.consumeFloatIf(73, "field length", nil)
	r.consumeNumberIf(74, DecRadix, "vertical text justification", &attdef.Text.VJustification)

	r.consumeNumber(280, DecRadix, "lock position flag", nil)

	r.consumeNumberIf(71, DecRadix, "attachment point", &attdef.AttachmentPoint)
	r.consumeNumberIf(72, DecRadix, "drawing direction", &attdef.DrawingDirection)

	r.consumeCoordinatesIf(11, attdef.Direction[:])
	HelperParseEmbeddedObject(r)
}

func ParseAcDbDimension(r Reader, _ *entity.Attdef) {
	if r.assertNextLine("AcDbDimension") != nil {
		return
	}

	r.consumeNumber(280, DecRadix, "version number", nil)
	r.consumeStr(nil) // name of the block

	r.consumeCoordinates(coords3D[:])
	r.consumeCoordinates(coords3D[:])
	r.consumeCoordinatesIf(12, coords3D[:])

	r.consumeNumberIf(70, DecRadix, "dimension type", nil)
	r.consumeNumberIf(71, DecRadix, "attachment point", nil)
	r.consumeNumberIf(72, DecRadix, "dimension text-line spacing", nil)

	r.consumeNumberIf(41, DecRadix, "dimension text-line factor", nil)
	r.consumeNumberIf(42, DecRadix, "actual measurement", nil)

	r.consumeNumberIf(73, DecRadix, "not documented", nil)
	r.consumeNumberIf(74, DecRadix, "not documented", nil)
	r.consumeNumberIf(75, DecRadix, "not documented", nil)

	r.consumeStrIf(1, nil) // dimension text
	r.consumeFloatIf(53, "roation angle of the dimension", nil)
	r.consumeFloatIf(51, "horizontal direction", nil)

	r.consumeNumberIf(71, DecRadix, "attachment point", nil)
	r.consumeNumberIf(42, DecRadix, "actual measurement", nil)
	r.consumeNumberIf(73, DecRadix, "not documented", nil)
	r.consumeNumberIf(74, DecRadix, "not documented", nil)
	r.consumeNumberIf(75, DecRadix, "not documented", nil)

	r.consumeCoordinatesIf(210, coords3D[:])
	r.consumeStrIf(3, nil) // [3] dimension style name

	dim := ""
	r.consumeStr(&dim)

	switch dim {
	// should be acdb3pointangulardimension
	case "AcDb2LineAngularDimension":
		r.consumeCoordinates(coords3D[:]) // point for linear and angular dimension
		r.consumeCoordinates(coords3D[:]) // point for linear and angular dimension
		r.consumeCoordinates(coords3D[:]) // point for diameter, radius, and angular dimension
	case "AcDbAlignedDimension":
		r.consumeCoordinatesIf(12, coords3D[:]) // insertion point for clones of a dimension
		r.consumeCoordinates(coords3D[:])       // definition point for linear and angular dimensions
		r.consumeCoordinates(coords3D[:])       // definition point for linear and angular dimensions
		r.consumeFloatIf(50, "angle of rotated, horizontal, or vertical dimensions", nil)
		r.consumeFloatIf(52, "oblique angle", nil)
		r.consumeStrIf(100, nil) // subclass marker AcDbRotatedDimension
	default:
		log.Fatal("Dimension(", Line, ")", dim)
	}
}

func ParseAcDbViewport(r Reader, _ *entity.MText) {
	if r.assertNextLine("AcDbViewport") != nil {
		return
	}

	r.consumeCoordinates(coords3D[:])
	r.consumeFloat(40, "width in paper space units", nil)
	r.consumeFloat(41, "height in paper space units", nil)

	// => -1 0 On, 0 = Off
	r.consumeFloatIf(68, "viewport status field", nil)

	r.consumeNumber(69, DecRadix, "viewport id", nil)

	r.consumeCoordinates(coords2D[:]) // center point
	r.consumeCoordinates(coords2D[:]) // snap base point
	r.consumeCoordinates(coords2D[:]) // snap spacing point
	r.consumeCoordinates(coords2D[:]) // grid spacing point
	r.consumeCoordinates(coords3D[:]) // view direction vector
	r.consumeCoordinates(coords3D[:]) // view target point

	r.consumeFloat(42, "perspective lens length", nil)
	r.consumeFloat(43, "front clip plane z value", nil)
	r.consumeFloat(44, "back clip plane z value", nil)
	r.consumeFloat(45, "view height", nil)
	r.consumeFloat(50, "snap angle", nil)
	r.consumeFloat(51, "view twist angle", nil)

	r.consumeFloat(72, "circle zoom percent", nil)

	for r.code() == 331 {
		//r.consumeNumber(331, DecRadix, "frozen layer object Id/handle", nil)
		r.consumeNext()
	}

	r.consumeNumber(90, HexRadix, "viewport status bit-coded flags", nil)
	r.consumeNumberIf(340, DecRadix, "hard-pointer id/handle to entity that serves as the viewports clipping boundary", nil)
	r.consumeStr(nil) // [1]
	r.consumeNumber(281, DecRadix, "render mode", nil)
	r.consumeNumber(71, DecRadix, "ucs per viewport flag", nil)
	r.consumeNumber(74, DecRadix, "display ucs icon at ucs origin flag", nil)

	r.consumeCoordinates(coords3D[:]) // ucs origin
	r.consumeCoordinates(coords3D[:]) // ucs x-axis
	r.consumeCoordinates(coords3D[:]) // ucs y-axis

	r.consumeNumberIf(345, DecRadix, "id/handle of AcDbUCSTableRecord if UCS is a named ucs", nil)
	r.consumeNumberIf(346, DecRadix, "id/handle of AcDbUCSTableRecord of base ucs", nil)
	r.consumeNumber(79, DecRadix, "Orthographic type of UCS", nil)
	r.consumeFloat(146, "elevation", nil)
	r.consumeNumber(170, DecRadix, "ShadePlot mode", nil)
	r.consumeNumber(61, DecRadix, "frequency of major grid lines compared to minor grid lines", nil)

	r.consumeNumberIf(332, DecRadix, "background id/handle", nil)
	r.consumeNumberIf(333, DecRadix, "shade plot id/handle", nil)
	r.consumeNumberIf(348, DecRadix, "visual style id/handle", nil)

	r.consumeNumber(292, DecRadix, "default lighting type on when no use lights are specified", nil)
	r.consumeNumber(282, DecRadix, "default lighting type", nil)
	r.consumeFloat(141, "view brightness", nil)
	r.consumeFloat(142, "view contrast", nil)

	r.consumeFloatIf(63, "ambient light color only if not black", nil)
	r.consumeFloatIf(421, "ambient light color only if not black", nil)
	r.consumeFloatIf(431, "ambient light color only if not black", nil)

	r.consumeNumberIf(361, DecRadix, "sun id/handle", nil)
	r.consumeNumberIf(335, DecRadix, "soft pointer reference to id/handle", nil)
	r.consumeNumberIf(343, DecRadix, "soft pointer reference to id/handle", nil)
	r.consumeNumberIf(344, DecRadix, "soft pointer reference to id/handle", nil)
	r.consumeNumberIf(91, DecRadix, "soft pointer reference to id/handle", nil)
}
