package test

import (
    "testing"

    "github.com/aichingert/dxf"
)

func TestOpen(t *testing.T) {
    drawing := dxf.Open("test.dxf")

    t.Fatalf(drawing)
}
