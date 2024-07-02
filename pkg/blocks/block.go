package blocks

import (
	"github.com/aichingert/dxf/pkg/entity"
)

type Block struct {
	Entity *entity.EntityData
	*entity.EntitiesData

	Attdefs []*entity.Attdef

	BlockName string
	OtherName string
	EndHandle uint64
	Flag      uint64

	XRefPath    string
	Description string
	Coordinates [3]float64
}

func NewBlock() *Block {
	return &Block{
		Entity: &entity.EntityData{
			Handle:    0,
			Owner:     0,
			LayerName: "",
		},
		EntitiesData: entity.New(),
		BlockName:    "",
		EndHandle:    0,
		Flag:         0,

		XRefPath:    "",
		Description: "",
		Coordinates: [3]float64{0.0, 0.0, 0.0},
	}
}
