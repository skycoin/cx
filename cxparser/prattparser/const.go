package parser

import (
	"github.com/skycoin/cx/cxparser/ast"
	"github.com/skycoin/cx/cxparser/cxlexer"
	"github.com/skycoin/cx/cxparser/token"
)

const (
	_ int = iota
	LOWEST

	// ==
	EQUALS

	// > or <
	LESSGREATER

	// +
	SUM

	// *
	PRODUCT

	// -value or !=condition
	PREFIX

	//add(1,2)
	CALL

	INDEX
)

//precedences set the precedence of token type.

var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
	token.LPAREN:   CALL,
	token.LBRACKET: INDEX,
}

type (

	/*
		prefixParseFn gets callled when we encounter the
		associated token type in prefix position.
	*/
	prefixParseFn func() ast.Expression

	/*
		infixParseFn gets callled when we encounter the
		associated token type in infix position.
	*/

	infixParseFn func(ast.Expression) ast.Expression
)

/*

Parser represents parser for cx programing language.

	 *cxlexer.Lexer which provide token.

	each token type can have up to two parsing functions assocated with it,
	depends on whether the token is found in a prefix position or infix position.

*/

type Parser struct {
	l      *cxlexer.Lexer
	errors []string

	//curToken points to current token
	curToken token.Token

	//peekToken token point to next current token
	peekToken token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}
