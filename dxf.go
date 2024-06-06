package dxf

import (
    reader "github.com/aichingert/dxf/pkg/reader"
)

func Open(filename string) {
    reader.Open(filename)
}
