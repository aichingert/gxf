package test

import (
    "os"
    "fmt"
    "testing"

    "github.com/stretchr/testify/require"

    "github.com/aichingert/gxf"
)

func TestMain(m *testing.M) {
    _, err := os.ReadFile("test.dxf")


    if err != nil {
        panic("Test file not found!")
    }

    m.Run()
}

func Test_ParsingDxfFile_ToGxf(t *testing.T) {
    bytes, err := os.ReadFile("test.dxf")
    require.Nil(t, err)

    plan, err := gxf.Parse(bytes)

    _ = plan
    fmt.Println(err)
}
