package parser

import (
    "os"
    "fmt"

    "bufio"
    "strconv"
    "strings"
)

const dxfCodeLineSizeInBytes = 4

type Reader struct {
    err     error
    reader  *bufio.Reader
    dxfLine *DxfLine
}

type DxfLine struct {
    Code    uint16
    Line    string
}

func NewReader(filename string) (*Reader, *os.File, error) {
    file, err := os.Open(filename)

    if err != nil {
        return nil, nil, err 
    }

    return &Reader {
        err: nil,
        reader: bufio.NewReader(file),
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
    line, err := r.ConsumeDxfLine()

    if err != nil { return nil, err }

    if line.Code != code {
        return nil, NewParseError(fmt.Sprintf("Invalid group code expected %d", code)) 
    }

    return line, nil
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
    line, err := r.reader.Peek(dxfCodeLineSizeInBytes)

    if err != nil {
        fmt.Println("[READER] unexpected eof ", err)
        return 0, err
    }

    code, err := strconv.ParseUint(strings.TrimSpace(string(line)), 10, 16)

    if err != nil {
        fmt.Println("[READER] unable to convert code to int ", err)
        return 0, err
    }

    return uint16(code), nil
}

func (r *Reader) consumeCode() (uint16, error) {
    line, err := r.reader.ReadBytes('\n')
    Line++

    if err != nil {
        fmt.Println("[READER(",Line,")] Corrupt Dxf file: ", err)
        return 0, err
    }

    code, err := strconv.ParseUint(strings.TrimSpace(string(line)), 10, 16)

    if err != nil {
        fmt.Println("[READER(",Line,")] Corrupt Dxf file: expected code got: ", err)
        return 0, err
    }

    return uint16(code), nil
}

func (r *Reader) ConsumeDxfLine() (*DxfLine, error) {
    code, err := r.consumeCode()
    if err != nil { return nil, err }

    line, err := r.reader.ReadBytes('\n')

    if err != nil {
        fmt.Println("[READER] Unexpected eof: ", err)
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
    line, err := r.AssertNext(code)

    if err != nil {
        fmt.Println("[TO_NUMBER(", Line, ")] failed: with invalid group code expected ", code, " got ", line)
        return 0, err
    }

    val, err := strconv.ParseUint(strings.TrimSpace(line.Line), radix, 64)

    if err != nil {
        fmt.Println("[TO_NUMBER(", Line, ")] failed: should be ", description, " got (", line, ")")
        return 0, NewParseError(description)
    }

    return val, nil
}

func (r *Reader) ConsumeFloat(code uint16, description string) (float64, error) {
    line, err := r.AssertNext(code)

    if err != nil {
        fmt.Println("[TO_FLOAT(", Line, ")] failed: with invalid group code expected ", code, " got ", line)
        return 0, err
    }

    val, err := strconv.ParseFloat(line.Line, 64)

    if err != nil { 
        fmt.Println("[READER] ConsumeFloat expected number got ", err) 
        return 0.0, NewParseError(description)
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

func (r *Reader) consumeCoordinates(coords []float64, len int) error {
    var err error
    for i := 0; i < len && r.ScanDxfLine(); i++ {
        switch coord := r.DxfLine(); coord.Code {
        case 10: fallthrough
        case 11: fallthrough
        case 210:
            if coords[0], err = strconv.ParseFloat(coord.Line, 64); err != nil {return err }
        case 20: fallthrough
        case 21: fallthrough
        case 220:
            if coords[1], err = strconv.ParseFloat(coord.Line, 64); err != nil { return err }
        case 30: fallthrough
        case 31: fallthrough
        case 230:
            if coords[2], err = strconv.ParseFloat(coord.Line, 64); err != nil { return err }
        default:
            fmt.Println("[READER(",Line,")] extract coordinates invalid index: ", coord)
            return NewParseError("Invalid group code in coordinates")
        }
    }

    return r.Err()
}
