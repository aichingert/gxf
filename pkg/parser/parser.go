package parser

import (
    "os"

    "log"
    "bufio"
    "strconv"

    "github.com/aichingert/dxf/pkg/drawing"
)

func FromFile(filename string) *drawing.Dxf {
    file, err := os.Open(filename)
    defer file.Close()

    if err != nil {
        log.Fatal("Failed to open file: ", err)
    }

    scanner := bufio.NewScanner(file)
    dxf     := drawing.New(filename)

    for true {
        switch data := ExtractCodeAndValue(scanner); data[1] {
        case "SECTION":
            switch section := ExtractCodeAndValue(scanner); section[1] {
            case "HEADER":
                ParseHeader(scanner, dxf)
            case "BLOCKS":
                ParseBlocks(scanner, dxf)
            case "CLASSES": fallthrough
            case "TABLES": fallthrough
            default:
                log.Println("WARNING: section not implemented: ", section)
                SkipToNextSection(scanner)
            }
        default:
            log.Fatal(data)
        }
    }

    return dxf
}

func SkipToNextSection(sc *bufio.Scanner) {
    for sc.Scan() && sc.Text() != "ENDSEC" {}
}

func ExtractCodeAndValue(sc *bufio.Scanner) [2]string {
    data := [2]string{}

    for line := 0; line < 2 && sc.Scan(); line++ {
        data[line] = sc.Text()
    }

    return data
}

func ExtractCoordinates(sc *bufio.Scanner) [3]float64 {
    coords := [3]float64{0.0, 0.0, 0.0}
    
    for i := 0; i < 3; i++ {
        coord := ExtractCodeAndValue(sc)
        axis, aerr := strconv.Atoi(coord[0])
        val,  verr := strconv.ParseFloat(coord[1], 64)

        if aerr != nil { panic(aerr) }
        if verr != nil { panic(verr) }

        coords[axis / 10 - 1] = val
    }

    return coords
}
