package parser

import (
	"fmt"
	"os"

	"bufio"
	"strconv"
	"strings"
)

const (
    DXF_CODE_LINE_SIZE_IN_BYTES = 4
    DEC_RADIX = 10
    HEX_RADIX = 16
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

func NewReader(filename string) (*Reader, *os.File, error) {
	file, err := os.Open(filename)

	if err != nil {
		return nil, nil, err
	}

	return &Reader{
		err:     nil,
		reader:  bufio.NewReader(file),
		dxfLine: nil,
	}, file, nil
}

func (r *Reader) ScanDxfLine() bool {
	r.dxfLine, r.err = r.ConsumeDxfLine()
	return r.err == nil
}

func (r *Reader) DxfLine() *DxfLine {
	return r.dxfLine
}

func (r *Reader) AssertNext(code uint16) (*DxfLine, error) {
	if r.err != nil {
		return nil, r.err
	}

	line, err := r.ConsumeDxfLine()
	r.err = err

	if r.err != nil {
		return nil, r.err
	}

	if line.Code != code {
		line = nil
		r.err = NewParseError(fmt.Sprintf("Invalid group code expected %d", code))
	}

	return line, r.err
}

func (r *Reader) AssertNextLine(line string) error {
	if r.err != nil {
		return r.err
	}

	check, err := r.ConsumeDxfLine()
	r.err = err

	if r.err != nil {
		return r.err
	}

	if check.Line != line {
		fmt.Println("[", Line, "] expected ", line, " got ", *check)
		r.err = NewParseError(fmt.Sprintf("[%d] expected %s", Line, line))
	}

	return r.err
}

func (r *Reader) ConsumeNumber(code uint16, radix int, description string, n *uint64) {
	if r.err != nil {
		return
	}

	line, err := r.AssertNext(code)
	r.err = err

	if r.err != nil {
		fmt.Println("[TO_NUMBER(", Line, ")] failed: with invalid group code expected ", code, " got ", line)
		r.err = NewParseError(description)
		return
	}

	if n != nil {
		*n, r.err = strconv.ParseUint(strings.TrimSpace(line.Line), radix, 64)
	}
}

func (r *Reader) ConsumeNumberIf(code uint16, radix int, description string, n *uint64) bool {
	if r.err != nil {
		return false
	}

	check, err := r.PeekCode()
	r.err = err

	if r.err != nil || check != code {
		return false
	}

	r.ConsumeNumber(code, radix, description, n)
	return r.err == nil
}

func (r *Reader) ConsumeFloat(code uint16, description string, f *float64) {
	if r.err != nil {
		return
	}

	line, err := r.AssertNext(code)
	r.err = err

	if r.err != nil {
		fmt.Println("[TO_FLOAT(", Line, ")] failed: with invalid group code expected ", code, " got ", line)
		return
	}

	if f != nil {
		*f, r.err = strconv.ParseFloat(line.Line, 64)
	}

	if r.err != nil {
		r.err = NewParseError(description)
		fmt.Println("[READER] ConsumeFloat expected number got ", err)
	}
}

func (r *Reader) ConsumeFloatIf(code uint16, description string, f *float64) bool {
	if r.err != nil {
		return false
	}

	check, err := r.PeekCode()
	r.err = err

	if r.err != nil || check != code {
		return false
	}

	r.ConsumeFloat(code, description, f)
	return true
}

func (r *Reader) ConsumeStr(s *string) {
	if r.err != nil {
		return
	}

	line, err := r.ConsumeDxfLine()
	r.err = err

	if r.err != nil {
		return
	}

	if s != nil {
		*s = line.Line
	}
}

func (r *Reader) ConsumeStrIf(code uint16, s *string) bool {
	if r.err != nil {
		return false
	}

	check, err := r.PeekCode()
	r.err = err

	if r.err != nil || check != code {
		return false
	}

	r.ConsumeStr(s)
	return r.Err() == nil
}

func (r *Reader) ConsumeCoordinates(coords []float64) {
	for i := 0; i < len(coords) && r.ScanDxfLine(); i++ {
		if r.err != nil {
			return
		}

		switch coord := r.DxfLine(); coord.Code {
		case 10:
			fallthrough
		case 11:
			fallthrough
		case 210:
			coords[0], r.err = strconv.ParseFloat(coord.Line, 64)
		case 20:
			fallthrough
		case 21:
			fallthrough
		case 220:
			coords[1], r.err = strconv.ParseFloat(coord.Line, 64)
		case 30:
			fallthrough
		case 31:
			fallthrough
		case 230:
			coords[2], r.err = strconv.ParseFloat(coord.Line, 64)
		default:
			fmt.Println("[READER(", Line, ")] extract coordinates invalid index: ", coord)
			r.err = NewParseError("Invalid group code in coordinates")
		}
	}
}

func (r *Reader) ConsumeCoordinatesIf(code uint16, coords []float64) {
	if r.err != nil {
		return
	}

	check, err := r.PeekCode()
	r.err = err

	if r.err != nil || check != code {
		return
	}

	r.ConsumeCoordinates(coords)
}

func (r *Reader) Err() error {
	return r.err
}

func (r *Reader) SkipToLabel(label string) error {
	for {
		line, err := r.ConsumeDxfLine()

		if err != nil {
			return err
		}

		if line.Line == label {
			return nil
		}
	}
}

func (r *Reader) PeekCode() (uint16, error) {
	if r.err != nil {
		return 0, r.err
	}

	line, err := r.reader.Peek(DXF_CODE_LINE_SIZE_IN_BYTES)

	if err != nil {
		fmt.Println("[READER] unexpected eof ", err)
		return 0, err
	}

	code, err := strconv.ParseUint(strings.TrimSpace(string(line)), DEC_RADIX, 16)

	if err != nil {
		fmt.Println("[READER] unable to convert code to int ", err)
		return 0, err
	}

	return uint16(code), nil
}

// TODO: this is bad change this
// not sure but to be refactored
func (r *Reader) PeekLine() (string, error) {
    if r.err != nil {
        return "", r.err
    }

    offset := DXF_CODE_LINE_SIZE_IN_BYTES + 2
    line, err := r.reader.Peek(offset)
    r.err = err

    if r.err != nil {
        return "", r.err
    }

    for len(line) < 1 || line[len(line) - 1] != '\n' {
        line, r.err = r.reader.Peek(offset)

        if r.err != nil {
            return "", r.err
        }

        offset += 1
    }

    return string(line[DXF_CODE_LINE_SIZE_IN_BYTES + 1:len(line) - 2]), r.err
}

func (r *Reader) consumeCode() (uint16, error) {
	line, err := r.reader.ReadBytes('\n')
	Line++

	if err != nil {
		fmt.Println("[READER(", Line, ")] Corrupt Dxf file: ", err)
		return 0, err
	}

	code, err := strconv.ParseUint(strings.TrimSpace(string(line)), DEC_RADIX, 16)

	if err != nil {
		fmt.Println("[READER(", Line, ")] Corrupt Dxf file: expected code got: ", err)
		return 0, err
	}

	return uint16(code), nil
}

func (r *Reader) ConsumeDxfLine() (*DxfLine, error) {
	if r.err != nil {
		return nil, r.err
	}

	code, err := r.consumeCode()
	if err != nil {
		return nil, err
	}

	line, err := r.reader.ReadBytes('\n')

	if err != nil {
		fmt.Println("[READER] Unexpected eof: ", err)
		return nil, err
	}

	offset := 1
	Line++

	// \r\n
	if len(line) > 1 && line[len(line)-2] == 13 {
		offset++
	}

	return &DxfLine{
		Code: code,
		Line: string(line[:len(line)-offset]),
	}, nil
}
