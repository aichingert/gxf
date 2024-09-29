package entity

type AttribAppender interface {
	AppendAttrib(attrib *Attrib)
}

type Attrib struct {
	Entity *EntityData

	*Text
	Tag   string
	Flags int64
}

func NewAttrib() *Attrib {
	return &Attrib{
		Entity: NewEntityData(),
		Text:   NewText(),

		Tag:   "",
		Flags: 0,
	}
}

func (i *Insert) AppendAttrib(attrib *Attrib) {
	i.Attributes = append(i.Attributes, attrib)
}
