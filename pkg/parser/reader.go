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

func NewReader(filename string) (*Reader, *os.File) {
    file, err := os.Open(filename)

    if err != nil {
        log.Fatal("[READER] Unable to open file: ", err)
    }

    return &Reader {
        reader: bufio.NewReader(file),
    }, file
}

func (r *Reader) SkipToLabel(label string) {
    for r.ConsumeDxfLine().Line != label {}
}

const dxfCodeBytes = 4

func (r *Reader) PeekCode() uint16 {
    line, err := r.reader.Peek(dxfCodeBytes)

    if err != nil {
        log.Fatal("[READER] unexpected eof ", err)
    }

    code, err := strconv.ParseUint(strings.TrimSpace(string(line)), 10, 16)

    if err != nil {
        log.Fatal("[READER] unable to convert code to int ", err)
    }

    return uint16(code)
}

func (r *Reader) consumeCode() uint16 {
    line, err := r.reader.ReadBytes('\n')
    Line++

    if err != nil {
        log.Fatal("[READER] Corrupt Dxf file: ", err)
    }

    code, err := strconv.ParseUint(strings.TrimSpace(string(line)), 10, 16)

    if err != nil {
        log.Fatal("[READER] Corrupt Dxf file: expected code got: ", err)
    }

    return uint16(code)
}

func (r *Reader) ConsumeDxfLine() DxfLine {
    code      := r.consumeCode()
    line, err := r.reader.ReadBytes('\n')
    offset    := 1
    Line++

    if err != nil {
        log.Fatal("[READER] Unexpected eof: ", err)
    }

    // \r\n
    if len(line) > 1 && line[len(line) - 2] == 13 {
        offset++
    }
    
    return DxfLine {
        Code: code,
        Line: string(line[:len(line) - offset]),
    }
}

func (r *Reader) ConsumeNumber(code uint16, radix int, description string) uint64 {
    line := r.ConsumeDxfLine()

    if line.Code != code {
        log.Fatal("[TO_HEX] failed: with invalid group code expected ", code, " got ", line)
    }

    val, err := strconv.ParseUint(strings.TrimSpace(line.Line), radix, 64)

    if err != nil {
        log.Fatal("[TO_HEX] failed: should be", description, " got (", line, ")")
    }

    return val
}

func (r *Reader) ConsumeCoordinates3D() [3]float64 {
    coords := [3]float64{0.0, 0.0, 0.0}
    r.consumeCoordinates(coords[:], len(coords))
    return coords
}

func (r *Reader) ConsumeCoordinates2D() [2]float64 {
    coords := [2]float64{0.0, 0.0}
    r.consumeCoordinates(coords[:], len(coords))
    return coords
}

func ParseFloat(value string) float64 {
    val, err := strconv.ParseFloat(value, 64)

    if err != nil { 
        log.Fatal("[READER] parseFloat expected number got ", err) 
    }

    return val
}

func (r *Reader) consumeCoordinates(coords []float64, len int) {
    for i := 0; i < len; i++ {
        switch coord := r.ConsumeDxfLine(); coord.Code {
        case 10: fallthrough
        case 11:
            coords[0] = ParseFloat(coord.Line)
        case 20: fallthrough
        case 21:
            coords[1] = ParseFloat(coord.Line)
        case 30: fallthrough
        case 31:
            coords[2] = ParseFloat(coord.Line)
        default:
            log.Fatal("[READER(",Line,")] extract coordinates invalid index: ", coord)
        }
    }
}
