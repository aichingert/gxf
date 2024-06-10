package parser

import (
    "log"

    "github.com/aichingert/dxf/pkg/drawing"
)

func parseCustomProperty(r *Reader, dxf *drawing.Dxf) {
    tag      := r.ConsumeDxfLine()
    property := r.ConsumeDxfLine()
    dxf.Header.CustomProperties[tag.Line] = property.Line
}

func ParseHeader(r *Reader, dxf *drawing.Dxf) {
    for {
        switch variable := r.ConsumeDxfLine(); variable.Line {
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
            data := r.ConsumeDxfLine()
            dxf.Header.Variables[variable.Line] = data.Line
        case "$ORTHOMODE":          fallthrough
        case "$REGENMODE":          fallthrough
        case "$FILLMODE":           fallthrough
        case "$QTEXTMODE":          fallthrough
        case "$MIRRTEXT":           fallthrough
        case "$ATTMODE":
            data := r.ConsumeDxfLine()
            dxf.Header.Modes[variable.Line] = data.Line
        case "$CUSTOMPROPERTYTAG":
            parseCustomProperty(r, dxf)
        case "$INSBASE":
            dxf.Header.InsBase = r.ConsumeCoordinates3D()
        case "$EXTMAX":
            dxf.Header.ExtMax  = r.ConsumeCoordinates3D()
        case "$EXTMIN":
            dxf.Header.ExtMin  = r.ConsumeCoordinates3D()
        case "$LIMMIN":
            log.Println(variable)
            dxf.Header.LimMin  = r.ConsumeCoordinates2D()
        case "$LIMMAX":
            dxf.Header.LimMax  = r.ConsumeCoordinates2D()
        case "ENDSEC":
            return
        default:
            log.Println("[HEADER] Warning [NOT IMPLEMENTED]: ", variable)
        }
    }
}
