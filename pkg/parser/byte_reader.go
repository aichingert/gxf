package parser

import "fmt"

type byteReader struct {
    bytes []byte
}

func (b *byteReader) consume() {
    fmt.Println("implements reader!")
}
