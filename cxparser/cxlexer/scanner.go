package cxlexer

import (
	"bufio"
	"bytes"
	"io"
)

type Scanner struct {
	r *bufio.Reader
}

func NewScanner(r io.Reader) *Scanner {

	return &Scanner{r: bufio.NewReader(r)}

}

func (s *Scanner) Scan() (tok token, lit string) {

	ch := s.read()

	if isWhitespace(ch) {
		s.unread()
		return s.scanWhitespace()
	} else if isLetter(ch) {
		s.unread()

	}

	buf.WriteRune(s.read())

}

func (s *Scanner) scanWhitespace() (tok token, lit string) {

	//Create a buffer and read the current character into it.

	var buf bytes.Buffer

	buf.WriteRune(s.read())

	//Read every subsequent whitespace characher into buffer.
	//Non-whitespace characters and EOF will cause the loop to exit.

	for {

		if ch := s.read(); ch == eof {

			break

		} else if !isWhitespace(ch) {

			s.unread()
			break

		} else {
			buf.WriteRune(ch)
		}

	}
}

// read reads the next rune from the buffered reader.
//Returns the rune(0) if an error occurs (or io.EOF is returned)
func (s *Scanner) read() rune {

	ch, _, err := s.r.ReadRune()

	if err != nil {

		return eof
	}

	return ch
}

func (s *Scanner) isWhitespace(ch rune) bool {

	return (ch == ' ' || ch == '\t' || ch == '\n')
}

//unread places the previously read rune back on the reader.
func (s *Scanner) unread() {

	_ = s.r.UnreadRune()
}

//isdigit returns true if the rune is a letter.
func isLetter(ch rune) bool {

	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

//isDigit returns true if thr rune is digit.
func isDigit(ch rune) bool {

	return (ch >= '0' && ch <= '9')
}

//eof represents a marker rune for the end of the reader.
var eof = rune(0)
