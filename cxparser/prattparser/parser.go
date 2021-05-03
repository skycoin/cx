package parser

import (
	"fmt"

	"github.com/skycoin/cx/cxparser/ast"
	"github.com/skycoin/cx/cxparser/token"

	"strconv"
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

//parseVarStatement returns ast.VarStatement.

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

//parseReturnStatement returns ast.ReturnStatement

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {

	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	stmt.ReturnValue = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
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

/*
	peekPrecedence method returns the precedence associates
	with the token type  p.peekToken.Type
*/
func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}

	return LOWEST
}

/*
	curPrecedence method returns the precedence associates
	with the token type  p.curToken.Type
*/
func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

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

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

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

/*
	parseBoolean parse boolean expression return ast.Expression.
*/
func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return exp
}

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

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.curToken, Function: function}
	exp.Arguments = p.parseExpressionList(token.RPAREN)
	return exp
}

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
