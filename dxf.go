package dxf

import (
    "log"

    reader "github.com/aichingert/dxf/pkg/reader"
)

func Open(filename string) {
    reader.Open(filename)
}
