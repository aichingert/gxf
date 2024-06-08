package test

import (
    "log"
    "testing"

    "github.com/aichingert/dxf"
)

func TestOpen(t *testing.T) {
    drawing := dxf.Open("test.dxf")

    log.Println(drawing)

}
