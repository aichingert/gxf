package parser

import (
	"fmt"
	"os"

	"bufio"
	"strconv"
	"strings"
)

const (
	DecRadix = 10
	HexRadix = 16
)

type Reader struct {
	err     error
	reader  *bufio.Reader
	dxfLine *DxfLine
}

type DxfLine struct {
	Code uint16
	Line string
}

func NewReader(filename string) (*Reader, *os.File) {
	file, err := os.Open(filename)

	return &Reader{
		err:     err,
		reader:  bufio.NewReader(file),
		dxfLine: &DxfLine{},
	}, file
}

// Consume reads the next code and line of the dxf file
func (r *Reader) consume() {
	if r.err != nil {
		return
	}

	r.dxfLine.Code, r.err = r.consumeCode()
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

	r.dxfLine.Line = string(buf[:len(buf)-offset])
}

// ConsumeNext returns the current value and reads the next one already
func (r *Reader) consumeNext() string {
	if r.err != nil {
		return r.err.Error()
	}
	defer r.consume()

	value := r.dxfLine.Line
	return value
}

func (r *Reader) consumeUntil(label string) {
	for r.consumeNext() != label {
	}
}

func (r *Reader) assertNextCode(expected uint16) error {
	if r.err != nil {
		return r.err
	}

	if r.dxfLine.Code != expected {
		r.dxfLine.Code = 1001
		r.err = NewParseError(fmt.Sprintf("Invalid group code expected %d", expected))
	}

	return r.err
}

func (r *Reader) assertNextLine(expected string) error {
	if r.err != nil {
		return r.err
	}

	if r.dxfLine.Line != expected {
		r.err = NewParseError(fmt.Sprintf("[%d] expected %s", Line, expected))
	}

	r.consume()
	return r.err
}

func (r *Reader) consumeNumber(code uint16, radix int, description string, n *int64) {
	if r.assertNextCode(code) != nil {
		return
	}

	if n != nil {
		*n, r.err = strconv.ParseInt(strings.TrimSpace(r.dxfLine.Line), radix, 64)

		if r.err != nil {
			r.err = NewParseError(fmt.Sprintf("%d %s because %s", Line, r.err.Error(), description))
		}
	}

	r.consume()
}

func (r *Reader) consumeNumberIf(code uint16, radix int, description string, n *int64) bool {
	if r.err != nil || r.dxfLine.Code != code {
		return false
	}

	r.consumeNumber(code, radix, description, n)
	return r.err == nil
}

func (r *Reader) consumeFloat(code uint16, description string, f *float64) {
	if r.assertNextCode(code) != nil {
		fmt.Println("[TO_FLOAT(", Line, ")] failed: with invalid group code expected ", code, " got ", r.dxfLine.Code)
		return
	}

	if f != nil {
		*f, r.err = strconv.ParseFloat(r.dxfLine.Line, 64)

		if r.err != nil {
			fmt.Println("[READER] ConsumeFloat expected number got ", r.dxfLine.Line)
			r.err = NewParseError(description)
		}
	}

	r.consume()
}

func (r *Reader) consumeFloatIf(code uint16, description string, f *float64) bool {
	if r.dxfLine.Code != code {
		return false
	}

	r.consumeFloat(code, description, f)
	return r.err == nil
}

func (r *Reader) consumeStr(s *string) {
	if r.err != nil {
		return
	}

	if s != nil {
		*s = r.dxfLine.Line
	}

	r.consume()
}

func (r *Reader) consumeStrIf(code uint16, s *string) bool {
	if r.dxfLine.Code != code {
		return false
	}

	r.consumeStr(s)
	return r.err == nil
}

func (r *Reader) consumeCoordinates(coords []float64) {
	for i := 0; i < len(coords); i++ {
		if r.err != nil {
			return
		}

		index := r.dxfLine.Code%100/10 - 1

		if index >= uint16(len(coords)) {
			r.err = NewParseError(fmt.Sprintf("cords out of bounds: len is %d index was %d", len(coords), index))
			return
		}

		coords[index], r.err = strconv.ParseFloat(r.dxfLine.Line, 64)
		r.consumeNext()
	}
}

func (r *Reader) consumeCoordinatesIf(code uint16, coords []float64) {
	if r.dxfLine.Code != code {
		return
	}

	r.consumeCoordinates(coords)
}

func (r *Reader) consumeCode() (uint16, error) {
	line, err := r.reader.ReadBytes('\n')
	Line++

	if r.err = err; err != nil {
		return 0, r.err
	}

	code, err := strconv.ParseUint(strings.TrimSpace(string(line)), DecRadix, 16)

	if r.err = err; r.err != nil {
		return 0, r.err
	}

	return uint16(code), nil
}
