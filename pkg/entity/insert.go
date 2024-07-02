package entity

type Insert struct {
	Entity           *EntityData
	BlockName        string
	AttributesFollow uint64

	Rotation float64

	RowCount   uint64
	ColCount   uint64
	RowSpacing float64
	ColSpacing float64

	Scale       [3]float64
	Coordinates [3]float64
	Attributes  []*Attrib
}

func NewInsert() *Insert {
	return &Insert{
		Entity:           NewEntityData(),
		BlockName:        "",
		AttributesFollow: 0,
		RowCount:         1,
		ColCount:         1,
		Scale:            [3]float64{1.0, 1.0, 1.0},
		Coordinates:      [3]float64{0.0, 0.0, 0.0},
	}
}

func (e *EntitiesData) AppendInsert(insert *Insert) {
	e.Inserts = append(e.Inserts, insert)
}
