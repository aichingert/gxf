package parser

import (
    "os"

    "log"
    "bufio"

    "github.com/aichingert/dxf/pkg/drawing"
)

func FromFile(filename string) *drawing.Dxf {
    file, err := os.Open(filename)
    defer file.Close()

    if err != nil {
        log.Fatal("Failed to open file: ", err)
    }

    scanner := bufio.NewScanner(file)
    dxf     := drawing.New(filename)

    ParseHeader(scanner, dxf)

    return dxf
}
