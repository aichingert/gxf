package parser

import "github.com/aichingert/dxf/pkg/drawing"

func parseCustomProperty(r *Reader, dxf *drawing.Dxf) {
	tag, property := "", ""
	r.ConsumeStr(&tag)
	r.ConsumeStr(&property)

	dxf.Header.CustomProperties[tag] = property
}

func ParseHeader(r *Reader, dxf *drawing.Dxf) {
	for r.ScanDxfLine() {
		switch r.DxfLine().Line {
		case "$ACADVER":
			fallthrough
		case "$ACADMAINTVER":
			fallthrough
		case "$DWGCODEPAGE":
			fallthrough
		case "$REQUIREDVERSIONS":
			fallthrough

		// TODO: maybe make it into a seperate variables
		case "$LTSCALE":
			fallthrough
		case "$TEXTSIZE":
			fallthrough
		case "$TRACEWID":
			fallthrough
		case "$TEXTSTYLE":
			fallthrough
		case "$CLAYER":
			fallthrough
		case "$CELTYPE":
			fallthrough
		case "$CECOLOR":
			fallthrough
		case "$CELTSCALE":
			fallthrough
		case "$DISPSILH":
			fallthrough
		case "$DIMSCALE":
			fallthrough
		case "$DIMASZ":
			fallthrough
		case "$DIMEXO":
			fallthrough
		case "$DIMDLI":
			fallthrough
		case "$DIMRND":
			fallthrough
		case "$DIMDLE":
			fallthrough
		case "$DIMEXE":
			fallthrough
		case "$DIMTP":
			fallthrough
		case "$DIMTM":
			fallthrough
		case "$DITXT":
			fallthrough
		case "$DIMCEN":
			fallthrough
		case "$DIMTSZ":
			fallthrough
		case "$DIMTOL":
			fallthrough
		case "$DIMLIM":
			fallthrough
		case "$LASTSAVEDBY":
			savedBy := ""
			r.ConsumeStr(&savedBy)
			dxf.Header.Variables[r.DxfLine().Line] = savedBy
		case "$ORTHOMODE":
			fallthrough
		case "$REGENMODE":
			fallthrough
		case "$FILLMODE":
			fallthrough
		case "$QTEXTMODE":
			fallthrough
		case "$MIRRTEXT":
			fallthrough
		case "$ATTMODE":
			mode := ""
			r.ConsumeStr(&mode)
			dxf.Header.Modes[r.DxfLine().Line] = mode
		case "$CUSTOMPROPERTYTAG":
			parseCustomProperty(r, dxf)
		case "$INSBASE":
			r.ConsumeCoordinates(dxf.Header.InsBase[:])
		case "$EXTMAX":
			r.ConsumeCoordinates(dxf.Header.ExtMax[:])
		case "$EXTMIN":
			r.ConsumeCoordinates(dxf.Header.ExtMin[:])
		case "$LIMMIN":
			r.ConsumeCoordinates(dxf.Header.LimMin[:])
		case "$LIMMAX":
			r.ConsumeCoordinates(dxf.Header.LimMax[:])
		case "ENDSEC":
			return
		}
	}
}
