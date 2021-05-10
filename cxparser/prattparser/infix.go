package parser

import "github.com/skycoin/cx/cxparser/ast"

/*
	parseInfixExpression parses infix epxression
*/
func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {

	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left, //left ast.Expression
	}

	//get the precence of current token
	precedence := p.curPrecedence()

	//move  one token
	p.nextToken()

	//this time passing in precedence of the operator token.
	expression.Right = p.parseExpression(precedence)

	return expression
}
