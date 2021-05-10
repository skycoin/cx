package parser

import (
	"github.com/skycoin/cx/cxparser/ast"
	"github.com/skycoin/cx/cxparser/token"
)

func (p *Parser) parseArrayLiteral() ast.Expression {

	array := &ast.ArrayLiteral{Token: p.curToken}

	array.Elements = p.parseExpressionList(token.RBRACKET)

	return array
}

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {

	exp := &ast.IndexExpression{Token: p.curToken, Left: left}

	p.nextToken()

	exp.Index = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RBRACKET) {
		return nil
	}

	return exp
}
