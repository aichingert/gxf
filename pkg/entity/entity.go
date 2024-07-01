package entity

type Entities interface {
	AppendArc(arc *Arc)
	AppendCircle(circle *Circle)
	AppendEllipse(ellipse *Ellipse)
	AppendLine(lines *Line)
	AppendPolyline(polyline *Polyline)
	AppendText(text *Text)
	AppendMText(mtext *MText)
	AppendHatch(hatch *Hatch)
	AppendInsert(insert *Insert)
}

type EntitiesData struct {
	Arcs      []*Arc
	Circles   []*Circle
	Ellipses  []*Ellipse
	Lines     []*Line
	Polylines []*Polyline
	Texts     []*Text
	MTexts    []*MText
	Hatches   []*Hatch
	Inserts   []*Insert
}

func New() *EntitiesData {
	return new(EntitiesData)
}

type Entity interface {
	GetHandle() *uint64
	GetOwner() *uint64
	GetVisibility() *uint64
	GetLayerName() *string
}

type EntityData struct {
	Handle     uint64
	Owner      uint64
	Visibility uint64

	LayerName string
}

func NewEntityData() *EntityData {
	return &EntityData{
		Handle:     0,
		Owner:      0,
		Visibility: 0,
		LayerName:  "",
	}
}

func (e *EntityData) GetHandle() *uint64 {
	return &e.Handle
}

func (e *EntityData) GetOwner() *uint64 {
	return &e.Owner
}

func (e *EntityData) GetVisibility() *uint64 {
	return &e.Visibility
}

func (e *EntityData) GetLayerName() *string {
	return &e.LayerName
}
