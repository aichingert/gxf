package entity

type MText struct {
	Entity *EntityData

	// [71] Attachment Point
	Layout int64
	// [72] Drawing direction
	Direction int64
	// [73] Line spacing (optional)
	LineSpacing int64

	// [1]/[3] 1 default text and 3 for additional (up to 250 characters)
	Text []string
	// [7] (optional)
	TextStyle string
	// [40] initial text height
	TextHeight float64

	// [50] Rotation Angel in radians
	// TODO: this are actually degrees and autocad does set
	// the 50 group code
	// https://forums.autodesk.com/t5/visual-lisp-autolisp-and-general/dxf-format-for-mtext-angle/td-p/2632521

	Vector      [3]float64
	Coordinates [3]float64
}

func NewMText() *MText {
	return &MText{
		Entity:      NewEntityData(),
		Layout:      0,
		Direction:   0,
		LineSpacing: 1,
		Text:        nil,
		TextStyle:   "",
		TextHeight:  0,
		Vector:      [3]float64{0.0, 0.0, 0.0},
		Coordinates: [3]float64{0.0, 0.0, 0.0},
	}
}

func (e *EntitiesData) AppendMText(mtext *MText) {
	e.MTexts = append(e.MTexts, mtext)
}
