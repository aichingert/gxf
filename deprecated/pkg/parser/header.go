package parser

import (
	"github.com/aichingert/dxf/pkg/drawing"
	"strings"
)

func ParseHeader(r Reader, dxf *drawing.Dxf) {
	for {
		line := r.consumeNext()

		if strings.HasSuffix(line, "MODE") {
			dxf.Header.Modes[r.line()] = r.consumeNext()
			continue
		}

		if line == "$CUSTOMPROPERTYTAG" {
			dxf.Header.CustomProperties[r.line()] = r.consumeNext()
			continue
		}

		switch line {
		case "$INSBASE":
			r.consumeCoordinates(dxf.Header.InsBase[:])
		case "$EXTMAX":
			r.consumeCoordinates(dxf.Header.ExtMax[:])
		case "$EXTMIN":
			r.consumeCoordinates(dxf.Header.ExtMin[:])
		case "$LIMMIN":
			r.consumeCoordinates(dxf.Header.LimMin[:])
		case "$LIMMAX":
			r.consumeCoordinates(dxf.Header.LimMax[:])
		case "ENDSEC":
			return
		default:
			if r.Err() != nil {
				return
			}

			dxf.Header.Variables[r.line()] = r.consumeNext()
		}
	}
}
