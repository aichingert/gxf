package header

type Header struct {
    // Stores the insertion base point set by BASE, which gets expressed as a UCS coordinate for the current space. 
    InsBase             [3]float64
    // 
    ExtMin              [3]float64
    // 
    ExtMax              [3]float64
    // 
    LimMin              [3]float64
    // 
    LimMax              [3]float64


    Modes               map[string]string
    Variables           map[string]string
    CustomProperties    map[string]string
}

func New() *Header {
    header := new (Header)

    header.ExtMin  = [3]float64{0.0, 0.0, 0.0}
    header.ExtMax  = [3]float64{0.0, 0.0, 0.0}
    header.LimMax  = [3]float64{0.0, 0.0, 0.0}
    header.LimMin  = [3]float64{0.0, 0.0, 0.0}
    header.InsBase = [3]float64{0.0, 0.0, 0.0}

    header.Modes = make(map[string]string)
    header.Variables = make(map[string]string)
    header.CustomProperties = make(map[string]string)

    return header
}
