package blocks

import (
	"github.com/aichingert/dxf/pkg/entity"
)

type Block struct {
	Entity *entity.EntityData

	BlockName string
	EndHandle uint64
	Flag      uint64

	XrefPath    string
	Description string
	Coordinates [3]float64

	// TODO: ATTDEF
}

func NewBlock() *Block {
	return &Block{
		Entity: &entity.EntityData{
			Handle:    0,
			Owner:     0,
			LayerName: "",
		},
		BlockName: "",
		EndHandle: 0,
		Flag:      0,

		XrefPath:    "",
		Description: "",
		Coordinates: [3]float64{0.0, 0.0, 0.0},
	}
}
