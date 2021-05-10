package parser

import (
	"github.com/skycoin/cx/cxparser/ast"
	"github.com/skycoin/cx/cxparser/token"
)

/*
	parseVarStatement returns ast.VarStatement.
*/
func (p *Parser) parseVarStatement() *ast.VarStatement {

	stmt := &ast.VarStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	if fl, ok := stmt.Value.(*ast.FunctionLiteral); ok {
		fl.Name = stmt.Name.Value
	}

	p.nextToken()

	return stmt
}
