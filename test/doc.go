package test

import (
    "testing"

    "github.com/aichingert/dxf"
)

func TestOpen() (t *testing.T) {
    dxf.Open("test.dxf")
}
