package entity

type Entity interface {
	GetHandle() *uint64
	GetOwner() *uint64
	GetLayerName() *string
}

type EntityData struct {
	Handle uint64
	Owner  uint64

	LayerName string
}

func (e *EntityData) GetHandle() *uint64 {
	return &e.Handle
}

func (e *EntityData) GetOwner() *uint64 {
	return &e.Owner
}

func (e *EntityData) GetLayerName() *string {
	return &e.LayerName
}
