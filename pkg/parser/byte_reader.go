package parser

import (
	"fmt"
	"strconv"
	"strings"
)

type byteReader struct {
	err     error
	idx     int
	bytes   []byte
	dxfLine *dxfLine
}

func NewByteReader(bytes []byte) *byteReader {
	return &byteReader{
		err:     nil,
		bytes:   bytes,
		dxfLine: &dxfLine{},
	}
}

func (r *byteReader) read(delimiter byte) ([]byte, error) {
	if r.idx >= len(r.bytes) {
		return nil, NewParseError("EOF")
	}

	start := r.idx
	Line++

	for r.idx < len(r.bytes) && r.bytes[r.idx] != delimiter {
		r.idx++
	}

	if r.idx >= len(r.bytes) {
		return nil, NewParseError("EOF")
	}

	r.idx++
	return r.bytes[start:r.idx], nil
}

func (r *byteReader) line() string {
	return r.dxfLine.line
}

func (r *byteReader) code() uint16 {
	return r.dxfLine.code
}

func (r *byteReader) consumeCode() {
	if r.err != nil {
		return
	}

	buf, err := r.read('\n')

	if r.err = err; err != nil {
		return
	}

	code, err := strconv.ParseUint(strings.TrimSpace(string(buf)), DecRadix, 16)

	if r.err = err; err != nil {
		return
	}

	r.dxfLine.code = uint16(code)
}

func (r *byteReader) consumeLine() {
	if r.err != nil {
		return
	}

	buf, err := r.read('\n')

	if r.err = err; err != nil {
		return
	}

	offset := 1

	if len(buf) > 1 && buf[len(buf)-2] == 13 {
		offset++
	}

	r.dxfLine.line = string(buf[:len(buf)-offset])
}

func (r *byteReader) consume() {
	if r.err != nil {
		return
	}

	r.consumeCode()
	r.consumeLine()
}

func (r *byteReader) consumeNext() string {
	if r.err != nil {
		return r.err.Error()
	}
	defer r.consume()

	value := r.dxfLine.line
	return value
}

func (r *byteReader) consumeUntil(label string) {
	for r.consumeNext() != label {
	}
}

func (r *byteReader) assertNextCode(expected uint16) error {
	if r.err != nil {
		return r.err
	}

	if r.dxfLine.code != expected {
		r.dxfLine.code = 1001
		r.err = NewParseError(fmt.Sprintf("Invalid group code expected %d", expected))
	}

	return r.err
}

func (r *byteReader) assertNextLine(expected string) error {
	if r.err != nil {
		return r.err
	}

	if r.dxfLine.line != expected {
		r.err = NewParseError(fmt.Sprintf("[%d] expected %s", Line, expected))
	}

	r.consume()
	return r.err
}

func (r *byteReader) consumeNumber(code uint16, radix int, description string, n *int64) {
	if r.assertNextCode(code) != nil {
		return
	}

	if n != nil {
		*n, r.err = strconv.ParseInt(strings.TrimSpace(r.dxfLine.line), radix, 64)

		if r.err != nil {
			r.err = NewParseError(fmt.Sprintf("%d %s because %s", Line, r.err.Error(), description))
		}
	}

	r.consume()
}

func (r *byteReader) consumeNumberIf(code uint16, radix int, description string, n *int64) bool {
	if r.err != nil || r.dxfLine.code != code {
		return false
	}

	r.consumeNumber(code, radix, description, n)
	return r.err == nil
}

func (r *byteReader) consumeFloat(code uint16, description string, f *float64) {
	if r.assertNextCode(code) != nil {
		fmt.Println("[TO_FLOAT(", Line, ")] failed: with invalid group code expected ", code, " got ", r.dxfLine.code)
		return
	}

	if f != nil {
		*f, r.err = strconv.ParseFloat(r.dxfLine.line, 64)

		if r.err != nil {
			fmt.Println("[READER] ConsumeFloat expected number got ", r.dxfLine.line)
			r.err = NewParseError(description)
		}
	}

	r.consume()
}

func (r *byteReader) consumeFloatIf(code uint16, description string, f *float64) bool {
	if r.dxfLine.code != code {
		return false
	}

	r.consumeFloat(code, description, f)
	return r.err == nil
}

func (r *byteReader) consumeStr(s *string) {
	if r.err != nil {
		return
	}

	if s != nil {
		*s = r.dxfLine.line
	}

	r.consume()
}

func (r *byteReader) consumeStrIf(code uint16, s *string) bool {
	if r.dxfLine.code != code {
		return false
	}

	r.consumeStr(s)
	return r.err == nil
}

func (r *byteReader) consumeCoordinates(cords []float64) {
	for i := 0; i < len(cords); i++ {
		if r.err != nil {
			return
		}

		index := r.dxfLine.code%100/10 - 1

		if index >= uint16(len(cords)) {
			r.err = NewParseError(fmt.Sprintf("cords out of bounds: len is %d index was %d", len(cords), index))
			return
		}

		cords[index], r.err = strconv.ParseFloat(r.dxfLine.line, 64)
		r.consumeNext()
	}
}

func (r *byteReader) consumeCoordinatesIf(code uint16, coords []float64) bool {
	if r.dxfLine.code != code {
		return false
	}

	r.consumeCoordinates(coords)
	return r.err == nil
}

func (r *byteReader) setErr(err error) {
	r.err = err
}

func (r *byteReader) Err() error {
	return r.err
}
