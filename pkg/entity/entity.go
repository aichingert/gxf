package entity

type Entity interface {
	GetHandle() *uint64
	GetOwner() *uint64
	GetVisibility() *uint64
	GetLayerName() *string
}

type EntityData struct {
	Handle uint64
	Owner  uint64
    Visibility uint64

	LayerName string
}

func NewEntityData() *EntityData {
    return &EntityData {
        Handle: 0,
        Owner:  0,
        Visibility: 0,
        LayerName: "",
    }
}

func (e *EntityData) GetHandle() *uint64 {
	return &e.Handle
}

func (e *EntityData) GetOwner() *uint64 {
	return &e.Owner
}

func (e *EntityData) GetVisibility() *uint64 {
	return &e.Visibility
}

func (e *EntityData) GetLayerName() *string {
	return &e.LayerName
}
