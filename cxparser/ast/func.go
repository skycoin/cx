package ast

import "github.com/skycoin/cx/cxparser/token"

/*
	FunctionLiteral represents cx programming language Function Literal.
*/
type FunctionLiteral struct {
	Token      token.Token // The 'func' token
	Parameters []*Identifier
	Body       *BlockStatement
	Name       string
}

/*
	CallExpression represents CallExpression of cx prograamming language.
*/
type CallExpression struct {
	Token     token.Token // The '(' token
	Function  Expression  // Identifier or FunctionLiteral
	Arguments []Expression
}
