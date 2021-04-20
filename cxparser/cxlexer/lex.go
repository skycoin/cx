package cxlexer

import "go/token"

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

func New(input string) *Lexer {

	l := &Lexer{input: input}
	l.readChar()

	return l
}

func (l *Lexer) NextToken() token.Token {

	var tok token.Token

	l.skipWhitespace()

	switch l.ch {

		case l.ch {
	



	}


	}
}


//isdigit returns true if the byte is a letter.
func isLetter(ch byte) bool {

	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

//isDigit returns true if thr byte is digit.
func isDigit(ch byte) bool {

	return (ch >= '0' && ch <= '9')
}



func newToken(TokenType token.TokenType ch byte) token.Token{

	return token.Token{

		Type : tokTokenType,
		Literal : string(ch)
	}

}


func (l *Lexer) readNumber() string{

	position := l.position

	for isDigit(l.ch){
		l.readChar()
	}

	return l.input[position:l.position]
}


func (l *Lexer) readString() string{

	position := l.position +1

	for{

		l.readChar()

		if l.ch == '"' || l.ch == 0{
				break
		}
	}

	return l.input[position:l.position]

}



func (l *Lexer) readIdentifier () string{

	position := l.position
	for isLetter(l.ch){

		l.readChar()
	}

	return l.input[position:l.position]
}
