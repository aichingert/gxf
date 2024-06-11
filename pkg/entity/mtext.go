package entity

type MText struct {
    *entity

    // [71] Attachment Point
    Layout      uint8
    // [72] Drawing direction
    Direction   uint8
    // [73] Line spacing (optional)
    LineSpacing uint8

    // [1]/[3] 1 default text and 3 for additional (up to 250 characters)
    Text        []string
    // [7] (optional)
    TextStyle   string
    // [40] initial text height
    TextHeight  float64

    // [50] Rotation Angel in radians
    // TODO: this are actually degrees and autocad does set
    // the 50 group code
    // https://forums.autodesk.com/t5/visual-lisp-autolisp-and-general/dxf-format-for-mtext-angle/td-p/2632521

    Vector      [3]float64
    Coordinates [3]float64
}

func NewMText(handle uint64, owner uint64) *MText {
    return &MText {
        entity: &entity {
            handle:     handle,
            owner:      owner,
            LayerName:  "",
        },
        Layout:         0,
        Direction:      0,
        LineSpacing:    1,
        Text:           nil,
        TextStyle:      "",
        TextHeight:     0,
        Vector:         [3]float64{0.0,0.0,0.0},
        Coordinates:    [3]float64{0.0,0.0,0.0},
    }
}
