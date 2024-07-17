package entity

type Hatch struct {
	Entity *EntityData

	PatternName string
	SolidFill   int64
	Associative int64

	Style   int64
	Pattern int64
	Angle   float64
	Scale   float64
	Double  int64

	PatternLines []*PatternLine

	Lines     []*Line
	Polylines []*Polyline
	Arcs      []*Arc
	Ellipses  []*Ellipse

	// TODO: spline
}

type PatternLine struct {
	Base        [2]float64
	Offset      [2]float64
	Angle       float64
	DashLengths []float64
}

func NewHatch() *Hatch {
	return &Hatch{
		Entity: NewEntityData(),

		PatternName: "",
		SolidFill:   0,
		Associative: 0,

		Lines:     nil,
		Polylines: nil,
		Arcs:      nil,
		Ellipses:  nil,
	}
}

func (h *Hatch) AppendPatternLine(
	angle float64,
	base [2]float64,
	offset [2]float64,
	dashes []float64) {
	patternLine := &PatternLine{
		Angle:       angle,
		Base:        base,
		Offset:      offset,
		DashLengths: dashes,
	}

	h.PatternLines = append(h.PatternLines, patternLine)
}

func (e *EntitiesData) AppendHatch(hatch *Hatch) {
	e.Hatches = append(e.Hatches, hatch)
}
