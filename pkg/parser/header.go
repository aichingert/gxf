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
    for {
        switch variable := ExtractCodeAndValue(sc); variable[1] {
        case "$ACADVER":            fallthrough
        case "$ACADMAINTVER":       fallthrough
        case "$DWGCODEPAGE":        fallthrough
        case "$REQUIREDVERSIONS":   fallthrough

        // TODO: maybe make it into a seperate variables
        case "$LTSCALE":            fallthrough
        case "$TEXTSIZE":           fallthrough
        case "$TRACEWID":           fallthrough
        case "$TEXTSTYLE":          fallthrough
        case "$CLAYER":             fallthrough
        case "$CELTYPE":            fallthrough
        case "$CECOLOR":            fallthrough
        case "$CELTSCALE":          fallthrough
        case "$DISPSILH":           fallthrough
        case "$DIMSCALE":           fallthrough
        case "$DIMASZ":             fallthrough
        case "$DIMEXO":             fallthrough
        case "$DIMDLI":             fallthrough
        case "$DIMRND":             fallthrough
        case "$DIMDLE":             fallthrough
        case "$DIMEXE":             fallthrough
        case "$DIMTP":              fallthrough
        case "$DIMTM":              fallthrough
        case "$DITXT":              fallthrough
        case "$DIMCEN":             fallthrough
        case "$DIMTSZ":             fallthrough
        case "$DIMTOL":             fallthrough
        case "$DIMLIM":             fallthrough

        case "$LASTSAVEDBY":
            data := ExtractCodeAndValue(sc)
            dxf.Header.Variables[variable[1]] = data[1]
        case "$ORTHOMODE":          fallthrough
        case "$REGENMODE":          fallthrough
        case "$FILLMODE":           fallthrough
        case "$QTEXTMODE":          fallthrough
        case "$MIRRTEXT":           fallthrough
        case "$ATTMODE":
            data := ExtractCodeAndValue(sc)
            dxf.Header.Modes[variable[1]] = data[1]
        case "$CUSTOMPROPERTYTAG":
            parseCustomProperty(sc, dxf)
        case "$INSBASE":
            dxf.Header.InsBase = ExtractCoordinates3D(sc)
        case "$EXTMAX":
            dxf.Header.ExtMax  = ExtractCoordinates3D(sc)
        case "$EXTMIN":
            dxf.Header.ExtMin  = ExtractCoordinates3D(sc)
        case "$LIMMIN":
            log.Println(variable)
            dxf.Header.LimMin  = ExtractCoordinates2D(sc)
        case "$LIMMAX":
            dxf.Header.LimMax  = ExtractCoordinates2D(sc)
        case "ENDSEC":
            return
        default:
            log.Println("[HEADER] Warning [NOT IMPLEMENTED]: ", variable)
        }
    }

}
