package gitconfig

import (
	"io"
	"unicode"
)

type lineType int

const (
	unknown lineType = 1 << iota
	comment
	section
	variable
	end
)

func Parse(in []byte) (*GitConfig, error) {
	c := new(configFile)
	c.init(in)
	gc, err := c.parse()
	if err != nil {
		return &GitConfig{}, err
	}

	return gc, nil
}

type configFile struct {
	data, buff []byte
	lstart     int  // start of the line
	eof        bool // end of file
	off        int  // curr position
	n          int  // num of chars
	cline      int  // current line number
}

func (c *configFile) init(data []byte) {
	c.data = data
	c.n = len(c.data)
	c.off = 0
	if c.n > 0 {
		c.buff = make([]byte, 0, 128)
	}
	c.cline = 1
}

func (c *configFile) lineString() string {
	for c.nextCh() != '\n' && !c.eof {
		_, _ = c.readCh()
	}

	if c.data[c.off-1] == '\r' {
		c.off--
	}

	return string(c.data[c.lstart:c.off])
}

func (c *configFile) readCh() (byte, error) {
	if c.off < c.n {
		cur := c.off
		c.off++
		if c.data[cur] == '\n' {
			c.lstart = c.off
			c.cline++
		}
		return c.data[cur], nil
	}
	c.eof = true
	return 0, io.EOF
}

func (c *configFile) nextCh() byte {
	if c.off < c.n {
		return c.data[c.off]
	}
	return 0
}

func (c *configFile) toEndOfLine() error {
	for c.nextCh() != '\n' {
		_, err := c.readCh()
		if err != nil {
			return err
		}

	}
	return nil
}

func (c *configFile) trimSpaceLeft() {
	for unicode.IsSpace(rune(c.nextCh())) {
		_, err := c.readCh()
		if err != nil {
			break
		}
	}
}

func (c *configFile) trimSpaceRight() {
	i := len(c.buff) - 1
	for ; i >= 0 && unicode.IsSpace(rune(c.buff[i])); i-- {
	}
	c.buff = c.buff[:i+1]
}

func (c configFile) getType() lineType {
	if isAlpha(c.nextCh()) {
		return variable
	}
	switch c.nextCh() {
	case '[':
		return section
	case ';', '#':
		return comment
	case 0:
		return end
	}
	return unknown
}

func (c *configFile) removeCarriageReturn() {
	if len(c.buff) > 0 && c.buff[len(c.buff)-1] == '\r' {
		c.buff = c.buff[:len(c.buff)-1]
	}
}

func (c *configFile) parse() (*GitConfig, error) {
	var (
		err error
		sec Section
	)
	gc := New()

	for !c.eof {
		c.trimSpaceLeft()
		switch c.getType() {
		case section:
			sec, err = c.parseSection()
			if err != nil {
				return nil, &ParseError{
					Err:        err,
					Line:       c.lineString(),
					LineNumber: c.cline,
				}
			}
		case variable:
			name, err := c.parseVariable()
			if err != nil {
				return nil, &ParseError{
					Err:        err,
					Line:       c.lineString(),
					LineNumber: c.cline,
				}
			}
			gc.add(sec, name, Value{string(c.buff)})
		case comment, end:
			err = c.toEndOfLine()
			if err != nil {
				break
			}
		default:
			return nil, &ParseError{
				Err:        ErrInvalidLine,
				Line:       c.lineString(),
				LineNumber: c.cline,
			}
		}
	}

	return gc, nil
}

func (c *configFile) parseSection() (Section, error) {
	c.buff = c.buff[:0]
	// first char is a '[', drop it
	_, _ = c.readCh()
loop:
	for {
		ch, err := c.readCh()
		if err != nil {
			break
		}
		switch ch {
		case ']': // end of section
			break loop
		case ' ': // probably have subsection
			ch, err := c.readCh()
			if err != nil {
				break loop
			}
			if ch != '"' {
				return Section{}, ErrInvalidLine
			}
			c.buff = append(c.buff, '.')
			err = c.parseSubsection()
			if err != nil {
				return Section{}, err
			}
			break loop
		}
		c.buff = append(c.buff, ch)
	}

	c.removeCarriageReturn()
	sec, err := NewSection(string(c.buff))
	if err != nil {
		return Section{}, err
	}

	_ = c.toEndOfLine()

	return sec, nil
}

// any characters except newline && null byte allowed
func (c *configFile) parseSubsection() error {
	for {
		ch, err := c.readCh()
		if err != nil {
			break
		}
		if ch == '\\' && c.nextCh() == 'n' {
			return ErrInvalidKey
		}
		c.buff = append(c.buff, ch)
		if c.nextCh() == '"' && ch != '\\' {
			break
		}
	}
	return nil
}

func (c *configFile) parseVariable() (VariableName, error) {
	var spaceFound bool
	c.buff = c.buff[:0]
	for { // parse variable name, only allows alphanumeric and '-'
		ch, err := c.readCh()
		if err != nil {
			break
		}
		if spaceFound && (isAlnum(ch) || ch == '-') {
			return "", ErrInvalidVariableName
		}
		if ch == '=' {
			break
		}
		if unicode.IsSpace(rune(ch)) {
			spaceFound = true
			continue
		}
		if !isAlnum(ch) && ch != '-' {
			return "", ErrInvalidVariableName
		}
		c.buff = append(c.buff, ch)
	}
	name := VariableName(string(c.buff))
	if !name.isValid() {
		return "", ErrInvalidVariableName
	}
	c.trimSpaceLeft()

	err := c.parseValue()
	if err != nil {
		return "", err
	}

	c.removeCarriageReturn()
	return name, nil
}

// allow any chars, '\' denotes value continues on the next line
func (c *configFile) parseValue() error {
	var (
		isQuoted      bool
		sameLineAsKey = true
	)
	c.buff = c.buff[:0]
loop:
	for {
		ch, err := c.readCh()
		if err != nil {
			break
		}

		if ch == '\n' {
			break
		}
		if !isQuoted && (ch == ';' || ch == '#') {
			_ = c.toEndOfLine()
			break
		}
		if ch == '\\' {
			switch c.nextCh() {
			case 0:
				break loop
			case '\n', '\r':
				err := c.toEndOfLine()
				if err != nil {
					break loop
				}
				_, _ = c.readCh()
				sameLineAsKey = false
				goto loop
			case 't', 'b', 'n', '\\', '"':
				if !isQuoted && c.nextCh() == 'n' {
					c.cline--
				}
				c.buff = append(c.buff, ch)
				ch, _ = c.readCh()
				goto add
			default:
				if sameLineAsKey {
					return ErrInvalidVariableValue
				} else {
					return ErrInvalidLine
				}
			}
		}

		if ch == '"' {
			isQuoted = !isQuoted
		}
	add:
		c.buff = append(c.buff, ch)
	}
	c.trimSpaceRight()

	return nil
}
