package parser

import (
    "os"

    "log"
    "bufio"
    "strconv"
    "strings"

    "github.com/aichingert/dxf/pkg/drawing"
)

var Line int64

func FromFile(filename string) *drawing.Dxf {
    file, err := os.Open(filename)
    defer file.Close()

    if err != nil {
        log.Fatal("Failed to open file: ", err)
    }

    scanner := bufio.NewScanner(file)
    dxf     := drawing.New(filename)

    for {
        switch data := ExtractCodeAndValue(scanner); data[1] {
        case "SECTION":
            switch section := ExtractCodeAndValue(scanner); section[1] {
            case "HEADER":
                ParseHeader(scanner, dxf)
            case "BLOCKS":
                ParseBlocks(scanner, dxf) 
            case "ENTITIES":
                ParseEntities(scanner, dxf)
            default:
                log.Println("WARNING: section not implemented: ", section)
                SkipToNextLabel(scanner, "ENDSEC")
            }
        case "EOF":
            return dxf
        default:
            log.Fatal(data)
        }
    }
}

func SkipToNextLabel(sc *bufio.Scanner, label string) {
    for sc.Scan() && sc.Text() != label {
        Line++
    }
}

func ExtractCodeAndValue(sc *bufio.Scanner) [2]string {
    data := [2]string{}

    for line := 0; line < 2 && sc.Scan(); line++ {
        data[line] = sc.Text()
        Line++
    }

    return data
}

func ExtractCoordinates3D(sc *bufio.Scanner) [3]float64 {
    coords := [3]float64{0.0, 0.0, 0.0}
    extractCoordinates(sc, coords[:], len(coords))
    return coords
}

func ExtractCoordinates2D(sc *bufio.Scanner) [2]float64 {
    coords := [2]float64{0.0, 0.0}
    extractCoordinates(sc, coords[:], len(coords))
    return coords
}

func extractCoordinates(sc *bufio.Scanner, coords []float64, len int) {
    for i := 0; i < len; i++ {
        coord := ExtractCodeAndValue(sc)
        axis, aerr := strconv.Atoi(strings.TrimSpace(coord[0]))
        val,  verr := strconv.ParseFloat(coord[1], 64)

        if aerr != nil { panic(aerr) }
        if verr != nil { panic(verr) }

        coords[axis / 10 - 1] = val
    }
}
