package entity

type Entity interface {
    Handle()    uint64
    Owner()     uint64
}

type entity struct {
    Handle      uint64
    Owner       uint64
}
