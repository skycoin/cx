package ast

import "bytes"

/*
	This method returns TokenLiteral return string.
*/
func (p *Program) TokenLiteral() string {

	if len(p.Statements) > 0 {

		return p.Statements[0].TokenLiteral()

	} else {
		return ""
	}
}

/*
	This method returns Program as string.
*/
func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

func (ls *VarStatement) statementNode() {}

func (oe *InfixExpression) expressionNode() {}

func (ie *IfExpression) expressionNode() {}

func (bs *BlockStatement) statementNode() {}

func (i *Identifier) expressionNode() {}

func (b *Boolean) expressionNode() {}

func (il *IntegerLiteral) expressionNode() {}

func (rs *ReturnStatement) statementNode() {}
