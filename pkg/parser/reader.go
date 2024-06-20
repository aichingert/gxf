package parser

import (
    "os"
    "log"

    "bufio"
    "strconv"
    "strings"
)

type Reader struct {
    reader  *bufio.Reader
}

type DxfLine struct {
    Code    uint16
    Line    string
}

func NewReader(filename string) (*Reader, *os.File, error) {
    file, err := os.Open(filename)

    if err != nil {
        log.Println("[READER] Unable to open file: ", err)
        return nil, nil, err
    }

    return &Reader {
        reader: bufio.NewReader(file),
    }, file, nil
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

const dxfCodeBytes = 4

func (r *Reader) PeekCode() (uint16, error) {
    line, err := r.reader.Peek(dxfCodeBytes)

    if err != nil {
        log.Println("[READER] unexpected eof ", err)
        return 0, err
    }

    code, err := strconv.ParseUint(strings.TrimSpace(string(line)), 10, 16)

    if err != nil {
        log.Println("[READER] unable to convert code to int ", err)
        return 0, err
    }

    return uint16(code), nil
}

func (r *Reader) consumeCode() (uint16, error) {
    line, err := r.reader.ReadBytes('\n')
    Line++

    if err != nil {
        log.Println("[READER(",Line,")] Corrupt Dxf file: ", err)
        return 0, err
    }

    code, err := strconv.ParseUint(strings.TrimSpace(string(line)), 10, 16)

    if err != nil {
        log.Println("[READER(",Line,")] Corrupt Dxf file: expected code got: ", err)
        return 0, err
    }

    return uint16(code), nil
}

func (r *Reader) ConsumeDxfLine() (*DxfLine, error) {
    code, err := r.consumeCode()

    if err != nil {
        return nil, err
    }

    line, err := r.reader.ReadBytes('\n')

    if err != nil {
        log.Println("[READER] Unexpected eof: ", err)
        return nil, err
    }

    offset    := 1
    Line++ 

    // \r\n
    if len(line) > 1 && line[len(line) - 2] == 13 {
        offset++
    }
    
    return &DxfLine {
        Code: code,
        Line: string(line[:len(line) - offset]),
    }, nil
}

func (r *Reader) ConsumeNumber(code uint16, radix int, description string) (uint64, error) {
    line, err := r.ConsumeDxfLine()

    if err != nil {
        return 0, err
    }

    if line.Code != code {
        log.Println("[TO_NUMBER(", Line, ")] failed: with invalid group code expected ", code, " got ", line)
        return 0, NewParseError("Invalid group code expected") 
    }

    val, err := strconv.ParseUint(strings.TrimSpace(line.Line), radix, 64)

    if err != nil {
        log.Println("[TO_NUMBER(", Line, ")] failed: should be ", description, " got (", line, ")")
        return 0, NewParseError(description)
    }

    return val, nil
}

func (r *Reader) ConsumeCoordinates3D() ([3]float64, error) {
    coords  := [3]float64{0.0, 0.0, 0.0}
    err     := r.consumeCoordinates(coords[:], len(coords))

    if err != nil { return coords, err }
    return coords, nil
}

func (r *Reader) ConsumeCoordinates2D() ([2]float64, error) {
    coords  := [2]float64{0.0, 0.0}
    err     := r.consumeCoordinates(coords[:], len(coords))

    if err != nil { return coords, err }
    return coords, nil
}

// TODO: take DxfLine to check for code and add description like ParseNumber
func ParseFloat(value string) (float64, error) {
    val, err := strconv.ParseFloat(value, 64)

    if err != nil { 
        log.Println("[READER] parseFloat expected number got ", err) 
        return 0.0, err
    }

    return val, nil
}

func (r *Reader) consumeCoordinates(coords []float64, len int) error {
    for i := 0; i < len; i++ {
        coord, err := r.ConsumeDxfLine()

        if err != nil { return err }

        switch coord.Code {
        case 10: fallthrough
        case 11:
            coords[0], err = ParseFloat(coord.Line)
            if err != nil { return err }
        case 20: fallthrough
        case 21:
            coords[1], err = ParseFloat(coord.Line)
            if err != nil { return err }
        case 30: fallthrough
        case 31:
            coords[2], err = ParseFloat(coord.Line)
            if err != nil { return err }
        default:
            log.Println("[READER(",Line,")] extract coordinates invalid index: ", coord)
            return NewParseError("Invalid group code in coordinates")
        }
    }

    return nil
}
