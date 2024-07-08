package gitconfig

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidKey           = errors.New("invalid key format")
	ErrKeyNotFound          = errors.New("could not find the given key")
	ErrInvalidValueType     = errors.New("invalid value type")
	ErrEmptyValue           = errors.New("empty value")
	ErrInvalidSection       = errors.New("illegal characters in section")
	ErrInvalidSubsection    = errors.New("illegal characters in subsection")
	ErrInvalidVariableName  = errors.New("illegal characters in variable name")
	ErrInvalidVariableValue = errors.New("illegal characters in variable value")
	ErrInvalidLine          = errors.New("illegal characters in line")
)

// ParseError returned if there's an error while parsing
type ParseError struct {
	Err        error
	Line       string
	LineNumber int
}

func (pe *ParseError) Error() string {
	return fmt.Sprintf("%s: %s (line %d)", pe.Err.Error(), pe.Line, pe.LineNumber)
}
