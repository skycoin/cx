package parser

import (
	"github.com/skycoin/cx/cxparser/ast"
	"github.com/skycoin/cx/cxparser/token"
)

/*
	parseFunctionLiteral parses Function Literal.
*/
func (p *Parser) parseFunctionLiteral() ast.Expression {

	lit := &ast.FunctionLiteral{Token: p.curToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	lit.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	lit.Body = p.parseBlockStatement()

	return lit
}

/*
	parseFunctionParameters  parse FunctionParameters returns []*ast.Identifier
*/
func (p *Parser) parseFunctionParameters() []*ast.Identifier {

	identifiers := []*ast.Identifier{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return identifiers
	}

	p.nextToken()

	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	identifiers = append(identifiers, ident)

	for p.peekTokenIs(token.COMMA) {

		p.nextToken()

		p.nextToken()

		ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

		identifiers = append(identifiers, ident)
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return identifiers
}

/* parseCallExpression parses Call Expression. */
func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {

	exp := &ast.CallExpression{Token: p.curToken, Function: function}

	exp.Arguments = p.parseExpressionList(token.RPAREN)

	return exp
}

/* parseExpressionList  parses ExpressionList. */
func (p *Parser) parseExpressionList(end token.TokenType) []ast.Expression {

	list := []ast.Expression{}

	if p.peekTokenIs(end) {
		p.nextToken()
		return list
	}

	p.nextToken()

	list = append(list, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		list = append(list, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(end) {
		return nil
	}

	return list
}
