package entity

type Text struct {
	Entity *EntityData

	Text  string
	Style string

	Flags int64
	// Horizontal
	HJustification int64
	// Vertical
	VJustification int64

	Rotation  float64
	Thickness float64
	XScale    float64
	Height    float64
	Oblique   float64

	Vector      [3]float64
	Coordinates [3]float64
}

func NewText() *Text {
	return &Text{
		Entity: NewEntityData(),

		Text:  "",
		Style: "STANDARD",

		Flags:          0,
		HJustification: 0,
		VJustification: 0,

		Rotation:  0.0,
		XScale:    1.0,
		Thickness: 1.0,
		Height:    1.0,
		Oblique:   0.0,

		Vector:      [3]float64{0.0, 0.0, 0.0},
		Coordinates: [3]float64{0.0, 0.0, 0.0},
	}
}

func (e *EntitiesData) AppendText(text *Text) {
	e.Texts = append(e.Texts, text)
}
