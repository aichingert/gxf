package parser

import (
    "log"

    "github.com/aichingert/dxf/pkg/drawing"
)

func parseCustomProperty(r *Reader, dxf *drawing.Dxf) error {
    tag, err := r.ConsumeDxfLine()
    if err != nil { return err }
    property, err := r.ConsumeDxfLine()
    if err != nil { return err }

    dxf.Header.CustomProperties[tag.Line] = property.Line
    return nil
}

func ParseHeader(r *Reader, dxf *drawing.Dxf) error {
    for {
        line, err := r.ConsumeDxfLine()
        if err != nil { return err }

        switch line.Line {
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
            data, err := r.ConsumeDxfLine()
            if err != nil { return err }
            dxf.Header.Variables[line.Line] = data.Line
        case "$ORTHOMODE":          fallthrough
        case "$REGENMODE":          fallthrough
        case "$FILLMODE":           fallthrough
        case "$QTEXTMODE":          fallthrough
        case "$MIRRTEXT":           fallthrough
        case "$ATTMODE":
            data, err := r.ConsumeDxfLine()
            if err != nil { return err }
            dxf.Header.Modes[line.Line] = data.Line
        case "$CUSTOMPROPERTYTAG":
            if err := parseCustomProperty(r, dxf); err != nil { return err }
        case "$INSBASE":
            insBase, err := r.ConsumeCoordinates3D()
            if err != nil { return err }
            dxf.Header.InsBase = insBase
        case "$EXTMAX":
            extMax, err := r.ConsumeCoordinates3D()
            if err != nil { return err }
            dxf.Header.ExtMax  = extMax
        case "$EXTMIN":
            extMin, err := r.ConsumeCoordinates3D()
            if err != nil { return err }
            dxf.Header.ExtMin = extMin
        case "$LIMMIN":
            limMin, err := r.ConsumeCoordinates2D()
            if err != nil { return err }
            dxf.Header.LimMin  = limMin 
        case "$LIMMAX":
            limMax, err := r.ConsumeCoordinates2D()
            if err != nil { return err }
            dxf.Header.LimMax  = limMax
        case "ENDSEC":
            return nil
        default:
            log.Println("[HEADER] Warning [NOT IMPLEMENTED]: ", line)
        }
    }
}
