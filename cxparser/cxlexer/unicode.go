package cxlexer

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"unicode"
)

type Lexer struct {
	r *bufio.Reader
	Position
}

type Position struct {
	line   int
	column int
}

//NewLexer returns a new instance of Lexer.
func NewLexer(r io.Reader) *Lexer {

	return &Lexer{r: bufio.NewReader(r)}

}

func (l *Lexer) Lex(Token string) {

	var tok Token

	for {

		r := l.readrune()

		switch r {

		case '=':

			tok = Token(ASSIGN, r)

		case '/':

			tok = Token(SLASH, r)

		case '+':

			tok = Token(PLUS, r)

		case '-':

			tok = Token(MINUS, r)

		case '*':

			tok = Token(ASTERISK, r)

		case '<':

			tok = Token(LT, r)

		case '>':

			tok = Token(GT, r)

		case '(':

			tok = Token(LPAREN, r)

		case ')':

			tok = Token(RPAREN, r)

		case '{':

			tok = Token(LBRACE, r)

		case '}':

			tok = Token(RBRACE, r)

		case '[':

			tok = Token(LBRACKET, r)

		case '[':

			tok = Token(LBRACKET, r)

		case '!':

			peek := l.readrune()

			if peek == '=' {
				tok = newToken(NOT_EQL, r)
			} else {
				tok = newToken(BANG, r)
			}

		case ':':

			tok = newToken(COLON, r)

		case ';':

			tok = newToken(SEMICOLON, r)

		case '=':
			peek := l.readrune()
			if peek == '=' {
				tok = newToken(EQL, r)
			} else {

				tok = newToken(EQL, r)
			}

		default:
			if unicode.IsSpace(r) {
				//todo
				continue

			} else if unicode.IsDigit(r) {
				l.unreadrune()
				l.scanint()

			} else if unicode.IsLetter(r) {

				l.unreadrune()
				l.scanIdent()
			}

		}
	}

}

func (l *Lexer) scanIdent() string {

	var buf bytes.Buffer

	r := l.readrune()

	buf.WriteRune(r)

	for {

		if ch := l.readrune(); ch == eof {
			break
		} else if unicode.IsDigit(ch) && unicode.IsLetter(ch) && ch != '_' {
			l.unreadrune()
			break
		} else {

			_, err := buf.WriteRune(r)

			if err != nil {
				log.Fatalln("Lex scanIdent", err)
			}
		}

	}

	keyword := buf.String()

	switch keyword {

	case "var":
		return keyword

		tok = newToken(VAR, r)

	case "true":
		tok = newToken(TRUE, r)

	case "false":
		tok = newToken(FALSE, r)
	case "if":

		tok = newToken(IF, r)

	case "else":

		tok = newToken(ELSE, r)

	case "return":

		tok = newToken(RETURN, r)

	case "func":

		tok = newToken(FUNC, r)

	}

	return tok
}

func (l *Lexer) scanint() string {

	var lit string

	for {

		r := l.readrune()

		if unicode.IsDigit(r) {

			lit = lit + string(r)
		} else {

			l.unreadrune()
			return lit
		}

	}

}

// read reads the next rune from the buffered reader.
//Returns the rune(0) if an error occurs (or io.EOF is returned)
func (l *Lexer) readrune() rune {

	ch, _, err := l.r.ReadRune()

	if err != nil {

		return eof
	}

	return ch
}

//unread places the previously read rune back on the reader.
func (s *Lexer) unreadrune() {

	err := s.r.UnreadRune()

	if err != nil {
		log.Fatalln("Lex Error : UnreadRune ", err)
	}
}

//eof represents a marker rune for the end of the reader.
var eof = rune(0)

func newToken(TokenType TokenType, ch rune) Token {

	return Token{
		Type:    TokenType,
		Literal: string(ch),
	}

}
