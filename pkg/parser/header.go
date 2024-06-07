package parser

import (
    "log"
    "bufio"

    "github.com/aichingert/dxf/pkg/drawing"
)

func parseCustomProperty(sc *bufio.Scanner, dxf *drawing.Dxf) {
    tag      := ExtractCodeAndValue(sc)
    property := ExtractCodeAndValue(sc)
    dxf.Header.CustomProperties[tag[1]] = property[1]
}

func ParseHeader(sc *bufio.Scanner, dxf *drawing.Dxf) {
    for true {
        variable := ExtractCodeAndValue(sc)

        switch variable[1] {
        case "$ACADVER":            fallthrough
        case "$ACADMAINTVER":       fallthrough
        case "$DWGCODEPAGE":        fallthrough
        case "$REQUIREDVERSIONS":   fallthrough
        case "$LASTSAVEDBY":
            data := ExtractCodeAndValue(sc)
            dxf.Header.Variables[variable[1]] = data[1]
        case "$CUSTOMPROPERTYTAG":  
            parseCustomProperty(sc, dxf)
        case "$INSBASE":
            dxf.Header.InsBase = ExtractCoordinates(sc)
        case "$EXTMAX":
            dxf.Header.ExtMax  = ExtractCoordinates(sc)
        case "$EXTMIN":
            dxf.Header.ExtMin  = ExtractCoordinates(sc)
        case "$LIMMIN":
            dxf.Header.LimMin  = ExtractCoordinates(sc)
        case "$LIMMAX":
            dxf.Header.LimMax  = ExtractCoordinates(sc)
        default:
            log.Fatal(variable[1])
        }
    }

}
