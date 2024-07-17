package entity

type Attdef struct {
	Entity *EntityData

	*Text

	Tag    string
	Prompt string

	Flags            int64
	AttachmentPoint  int64
	DrawingDirection int64

	Direction [3]float64
}

func NewAttdef() *Attdef {
	return &Attdef{
		Entity: NewEntityData(),
		Text:   NewText(),

		Tag:              "",
		Flags:            0,
		AttachmentPoint:  0,
		DrawingDirection: 0,
		Direction:        [3]float64{0.0, 0.0, 0.0},
	}
}
