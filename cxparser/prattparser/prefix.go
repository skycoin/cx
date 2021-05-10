package parser

import "github.com/skycoin/cx/cxparser/ast"

/*
	parsePrefixExpression parses Prefix Expression
	and create ast.Expression.
*/
func (p *Parser) parsePrefixExpression() ast.Expression {

	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()

	//we got the expression and now parseExpression with precedence PREFIX.
	expression.Right = p.parseExpression(PREFIX)

	return expression
}
