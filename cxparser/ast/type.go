package ast

import "github.com/skycoin/cx/cxparser/token"

/*
	StringLiteral represents StringLiteral of cx prograamming language.
*/
type StringLiteral struct {
	Token token.Token
	Value string
}

/*
	IntegerLiteral represents cx programming language Integer Literal.
*/
type IntegerLiteral struct {
	Token token.Token
	Value int64
}

/*
	Boolean represents cx programming language Boolean.
*/
type Boolean struct {
	Token token.Token
	Value bool
}

/*
	Identifier represents cx programming language Identifier.
*/
type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string
}
