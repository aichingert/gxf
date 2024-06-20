package parser

import "fmt"

type ParseError struct {
    msg     string
}

func NewParseError(msg string) *ParseError {
    return &ParseError {
        msg,
    }
}

func (p *ParseError) Error() string {
    return fmt.Sprintf("%d: %s", Line, p.msg)
}
