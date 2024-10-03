package parser

import (
    "strings"
    "strconv"
)

type byteReader struct {
    index int
    bytes []byte
}

func (b *byteReader) consumeCode(code *uint16) error {
    buf, err := b.readUntil('\n')

    if err != nil {
        return err
    }

    bCode, err := strconv.ParseUint(strings.TrimSpace(string(buf)), decRadix, 16)
    *code = uint16(bCode)

    return err
}

func (b *byteReader) consumeLine(line *string) error {
    buf, err := b.readUntil('\n')

    if err != nil {
        return err
    }

    offset := 1

    // NOTE: 13 => \r
    if len(buf) > 1 && buf[len(buf) - 2] == 13 {
        offset++
    }

    *line = string(buf[:len(buf) - offset])
    return nil
}

func (b *byteReader) readUntil(delim byte) ([]byte, error) {
    if b.index >= len(b.bytes) {
        return nil, NewParseError("EOF")
    }

    src := b.index

    for b.index < len(b.bytes) && b.bytes[b.index] != delim {
        if b.bytes[b.index] == '\n' {
            Line++
        }

        b.index++
    }

    if b.index >= len(b.bytes) {
        return nil, NewParseError("EOF")
    }

    b.index++
    return b.bytes[src:b.index], nil
}
