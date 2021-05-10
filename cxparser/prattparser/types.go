package parser

import (
	"fmt"
	"strconv"

	"github.com/skycoin/cx/cxparser/ast"
	"github.com/skycoin/cx/cxparser/token"
)

/*
	parseIdentifier return  ast.Expression of Identifier.
*/
func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

/*
	parseIntegerLiteral return  ast.Expression of Integer.
*/
func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value

	return lit
}

/*
	parseStringLiteral return  ast.Expression of Integer.
*/
func (p *Parser) parseStringLiteral() ast.Expression {

	return &ast.StringLiteral{
		Token: p.curToken,
		Value: p.curToken.Literal}
}

/*
	parseBoolean parse boolean expression return ast.Expression.
*/
func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
}
