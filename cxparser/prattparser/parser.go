package parser

import (
	"github.com/skycoin/cx/cxparser/ast"
	"github.com/skycoin/cx/cxparser/token"
)

/*
	ParseProgram is entry point fuction for parser which construct
	root node of AST *ast.Program.


	It then builds the child nodes, statements, by calling
	other functions that know which AST node to constuct based on current token.

*/
func (p *Parser) ParseProgram() *ast.Program {

	program := &ast.Program{}

	program.Statements = []ast.Statement{}

	//parser token til token.EOF
	for !p.curTokenIs(token.EOF) {

		stmt := p.parseStatement()

		if stmt != nil {

			program.Statements = append(program.Statements, stmt)

		}

		//update the currenttoken and peektoken
		p.nextToken()
	}

	return program
}

/*
	parseStatement take decision base on token to delegate parser
	to either
	i) var statements
   ii) return statment
  iii) expression

*/
func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.VAR:
		return p.parseVarStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

//parseExpressionStatement return *ast.ExpressionStatement
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {

	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

/*
	parseExpression parse Expression based on precedence and creates
	ast.Expression.
*/
func (p *Parser) parseExpression(precedence int) ast.Expression {

	/*
		step 1 : check parsing function is associated
				 with p.curToken.Type in prefix position.
	*/

	prefix := p.prefixParseFns[p.curToken.Type]

	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}

	/*
		step 2 :
		we have parsing function available so call it
		it will return leftExp.

	*/

	leftExp := prefix()

	/*
		check the precedence of token
		and call again
		infix(leftExp) till infix == nil
	*/
	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {

		//we took for infixParse function  for p.peekToken.Type
		// and if we found infixfn
		infixfn := p.infixParseFns[p.peekToken.Type]

		if infixfn == nil {
			return leftExp
		}

		//move ahead by one token
		p.nextToken()

		//call that infixfn on leftExp
		leftExp = infixfn(leftExp)
	}

	return leftExp
}
