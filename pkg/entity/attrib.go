package entity

type Attrib struct {
	Entity *EntityData

	*Text
	Tag   string
	Flags uint64
}

func NewAttrib() *Attrib {
	return &Attrib{
		Entity: NewEntityData(),
		Text:   NewText(),

		Tag:   "",
		Flags: 0,
	}
}
