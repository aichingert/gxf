package header

type Header struct {
    Version string
}

func New() *Header {
    header := new (Header)

    header.Version = "TODO"

    return header
}
