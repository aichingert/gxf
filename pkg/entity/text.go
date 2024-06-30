package entity

type Text struct {
	Entity *EntityData

	Text  string
	Style string

	Flags uint64
	// Horizontal
	HJustification uint64
	// Vertical
	VJustification uint64

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
