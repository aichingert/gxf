package table

import "github.com/aichingert/dxf/pkg/entity"

type Layer struct {
	Entity *entity.EntityData

	Color     int64
	TrueColor int64 // 24-bit color
	LineType  string
}

func NewLayer() *Layer {
	return &Layer{
		Entity:    entity.NewEntityData(),
		Color:     0,
		TrueColor: 0,
		LineType:  "",
	}
}
