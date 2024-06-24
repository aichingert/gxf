package test

import (
    "testing"

    "github.com/aichingert/dxf"
)

// TODO: implement proper tests
// with minimal things

func TestOpen(t *testing.T) {
    drawing, err := dxf.Open("test.dxf")

    _ = drawing
    _ = err
}
