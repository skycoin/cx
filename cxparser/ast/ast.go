package ast

/*
 The base Node interface
	it has two methods
	TokenLiteral()
	String ()
*/
type Node interface {
	TokenLiteral() string
	String() string
}

/*
 statement interface reperents cx programming language Statement.

*/
type Statement interface {
	Node
	//	statementNode()
}

/*
 Expression interface reperents cx programming language Expression.

*/
type Expression interface {
	Node
	//	expressionNode()
}
