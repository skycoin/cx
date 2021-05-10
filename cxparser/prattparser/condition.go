package parser

import (
	"github.com/skycoin/cx/cxparser/ast"
	"github.com/skycoin/cx/cxparser/token"
)

/*
parseIfExpression parse parseIfExpression return ast.Expression.
*/
func (p *Parser) parseIfExpression() ast.Expression {

	expression := &ast.IfExpression{Token: p.curToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.nextToken()

	expression.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	expression.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(token.ELSE) {

		p.nextToken()

		if !p.expectPeek(token.LBRACE) {
			return nil
		}

		expression.Alternative = p.parseBlockStatement()
	}

	return expression
}

/*
	parseBlockStatement parse BlockStatement return ast.Expression.
*/
func (p *Parser) parseBlockStatement() *ast.BlockStatement {

	block := &ast.BlockStatement{Token: p.curToken}

	block.Statements = []ast.Statement{}

	p.nextToken()

	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {

		stmt := p.parseStatement()

		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}

/*
parseGroupedExpression parse GroupedExpression return ast.Expression.
*/
func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return exp
}
