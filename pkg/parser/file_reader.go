package parser

import (
	"fmt"
	"os"

	"bufio"
	"strconv"
	"strings"
)

type fileReader struct {
	err     error
	reader  *bufio.Reader
	dxfLine *dxfLine
}

func NewFileReader(filename string) (*fileReader, *os.File) {
	file, err := os.Open(filename)

	return &fileReader{
		err:     err,
		reader:  bufio.NewReader(file),
		dxfLine: &dxfLine{},
	}, file
}

func (r *fileReader) line() string {
	return r.dxfLine.line
}

func (r *fileReader) code() uint16 {
	return r.dxfLine.code
}

// Consume reads the next code and line of the dxf file
func (r *fileReader) consume() {
	if r.err != nil {
		return
	}

	r.consumeCode()
	r.consumeLine()
}

// ConsumeNext returns the current value and reads the next one already
func (r *fileReader) consumeNext() string {
	if r.err != nil {
		return r.err.Error()
	}
	defer r.consume()

	value := r.dxfLine.line
	return value
}

func (r *fileReader) consumeUntil(label string) {
	for r.consumeNext() != label {
	}
}

func (r *fileReader) assertNextCode(expected uint16) error {
	if r.err != nil {
		return r.err
	}

	if r.dxfLine.code != expected {
		r.dxfLine.code = 1001
		r.err = NewParseError(fmt.Sprintf("Invalid group code expected %d", expected))
	}

	return r.err
}

func (r *fileReader) assertNextLine(expected string) error {
	if r.err != nil {
		return r.err
	}

	if r.dxfLine.line != expected {
		r.err = NewParseError(fmt.Sprintf("[%d] expected %s", Line, expected))
	}

	r.consume()
	return r.err
}

func (r *fileReader) consumeNumber(code uint16, radix int, description string, n *int64) {
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

func (r *fileReader) consumeNumberIf(code uint16, radix int, description string, n *int64) bool {
	if r.err != nil || r.dxfLine.code != code {
		return false
	}

	r.consumeNumber(code, radix, description, n)
	return r.err == nil
}

func (r *fileReader) consumeFloat(code uint16, description string, f *float64) {
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

func (r *fileReader) consumeFloatIf(code uint16, description string, f *float64) bool {
	if r.dxfLine.code != code {
		return false
	}

	r.consumeFloat(code, description, f)
	return r.err == nil
}

func (r *fileReader) consumeStr(s *string) {
	if r.err != nil {
		return
	}

	if s != nil {
		*s = r.dxfLine.line
	}

	r.consume()
}

func (r *fileReader) consumeStrIf(code uint16, s *string) bool {
	if r.dxfLine.code != code {
		return false
	}

	r.consumeStr(s)
	return r.err == nil
}

func (r *fileReader) consumeCoordinates(coords []float64) {
	for i := 0; i < len(coords); i++ {
		if r.err != nil {
			return
		}

		index := r.dxfLine.code%100/10 - 1

		if index >= uint16(len(coords)) {
			r.err = NewParseError(fmt.Sprintf("cords out of bounds: len is %d index was %d", len(coords), index))
			return
		}

		coords[index], r.err = strconv.ParseFloat(r.dxfLine.line, 64)
		r.consumeNext()
	}
}

func (r *fileReader) consumeCoordinatesIf(code uint16, coords []float64) bool {
	if r.dxfLine.code != code {
		return false
	}

	r.consumeCoordinates(coords)
	return r.err == nil
}

func (r *fileReader) consumeCode() {
	line, err := r.reader.ReadBytes('\n')
	Line++

	if r.err = err; err != nil {
		return
	}

	code, err := strconv.ParseUint(strings.TrimSpace(string(line)), DecRadix, 16)
	if r.err = err; err != nil {
		return
	}

	r.dxfLine.code = uint16(code)
}

func (r *fileReader) consumeLine() {
	if r.err != nil {
		return
	}

	buf, err := r.reader.ReadBytes('\n')

	if err != nil {
		r.err = NewParseError("[READER] invalid file")
		return
	}

	offset := 1
	Line++

	// \r\n
	if len(buf) > 1 && buf[len(buf)-2] == 13 {
		offset++
	}

	r.dxfLine.line = string(buf[:len(buf)-offset])
}

func (r *fileReader) setErr(err error) {
	r.err = err
}

func (r *fileReader) Err() error {
	return r.err
}
