package parser

import "fmt"

type byteReader struct {
    bytes []byte
}

func (b *byteReader) consumeCode(code *uint16, err *error) {
    if err != nil {
        return
    }

    fmt.Println(code)
}

func (b *byteReader) consumeLine(line *string, err *error) {
    if err != nil {
        return
    }

    fmt.Println(line)
}


