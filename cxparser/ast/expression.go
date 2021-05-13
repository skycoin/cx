package ast

import "github.com/skycoin/cx/cxparser/token"

/*
	ExpressionStatement struct represents ExpressionStatement of cx programming language.
*/
type ExpressionStatement struct {
	Token      token.Token // the first token of the expression
	Expression Expression
}

func (es *ExpressionStatement) statementNode() {}

/*
	IfExpression represents cx programming language IfExpression.
*/
type IfExpression struct {
	Token       token.Token // The 'if' token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

/*
	BlockStatement struct represents BlockStatement of cx programming language.
*/
type BlockStatement struct {
	Token      token.Token // the { token
	Statements []Statement
}

/*
	PrefixExpression represents cx programming language Prefix Expression.
*/
type PrefixExpression struct {
	Token    token.Token // The prefix token, e.g. !
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode() {}

/*
	InfixExpression represents cx programming language Infix Expression.
*/
type InfixExpression struct {
	Token    token.Token // The operator token, e.g. +
	Left     Expression
	Operator string
	Right    Expression
}
