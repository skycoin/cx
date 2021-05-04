package ast

import (
	"bytes"
	"fmt"
	"strings"
)

/*
	This VarStatement struct methiodd returns TokenLiteral  as a string of VarStatement.
*/
func (ls *VarStatement) TokenLiteral() string {

	return ls.Token.Literal
}

/*
	This VarStatement struct method to represents VarStatement into string.
*/
func (ls *VarStatement) String() string {

	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")

	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}

	out.WriteString(";")

	return out.String()
}

/*
	This ReturnStatement struct method returns TokenLiteral as a string of ReturnStatement.
*/
func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}

/*
	This ReturnStatement struct method to represents ReturnStatement into string.
*/
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")

	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}

	out.WriteString(";")

	return out.String()
}

/*
	This ExpressionStatement struct method returns TokenLiteral as a string of ExpressionStatement.
*/
func (es *ExpressionStatement) TokenLiteral() string {
	return es.Token.Literal
}

/*
	This ExpressionStatement struct method to represents ExpressionStatement into string.
*/
func (es *ExpressionStatement) String() string {

	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

/*
	This PrefixExpression struct method returns TokenLiteral as a string of PrefixExpression.
*/
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }

/*
	This PrefixExpression struct method to represents PrefixExpression into string.
*/
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}

/*
	This InfixExpression struct method returns TokenLiteral as a string of InfixExpression.
*/
func (oe *InfixExpression) TokenLiteral() string {
	return oe.Token.Literal
}

/*
	This InfixExpression struct method to represents InfixExpression into string.
*/
func (oe *InfixExpression) String() string {

	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(oe.Left.String())
	out.WriteString(" " + oe.Operator + " ")
	out.WriteString(oe.Right.String())
	out.WriteString(")")

	return out.String()
}

/*
	This IfExpression struct method returns TokenLiteral as a string of IfExpression.
*/
func (ie *IfExpression) TokenLiteral() string {
	return ie.Token.Literal
}

/*
	This IfExpression struct method to represents IfExpression into string.
*/
func (ie *IfExpression) String() string {

	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(ie.Condition.String())

	out.WriteString(" ")

	out.WriteString(ie.Consequence.String())

	if ie.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(ie.Alternative.String())
	}

	return out.String()
}

/*
	This BlockStatement's struct method returns TokenLiteral as a string of BlockStatement.
*/
func (bs *BlockStatement) TokenLiteral() string {
	return bs.Token.Literal
}

/*
	This BlockStatement struct method to represents BlockStatement into string.
*/
func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

/*
	This FunctionLiteral's struct method returns TokenLiteral as a string of FunctionLiteral.
*/
func (fl *FunctionLiteral) TokenLiteral() string {
	return fl.Token.Literal
}

/*
	This FunctionLiteral's struct method to represents FunctionLiteral into string.
*/
func (fl *FunctionLiteral) String() string {

	var out bytes.Buffer

	params := []string{}

	for _, p := range fl.Parameters {

		params = append(params, p.String())
	}

	out.WriteString(fl.TokenLiteral())

	if fl.Name != "" {
		out.WriteString(fmt.Sprintf("<%s>", fl.Name))
	}

	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(fl.Body.String())

	return out.String()
}

/*
	This CallExpression struct method returns TokenLiteral as a string of CallExpression.
*/
func (ce *CallExpression) TokenLiteral() string {
	return ce.Token.Literal
}

/*
	This CallExpression's struct method to represents CallExpression into string.
*/
func (ce *CallExpression) String() string {

	var out bytes.Buffer

	args := []string{}

	for _, a := range ce.Arguments {

		args = append(args, a.String())
	}

	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}

/*
	This StringLiteral struct method returns TokenLiteral as a string of StringLiteral.
*/
func (sl *StringLiteral) TokenLiteral() string {

	return sl.Token.Literal
}

/*
	This StringLiteral struct method to represents StringLiteral into string.
*/
func (sl *StringLiteral) String() string {

	return sl.Token.Literal
}

/*
	This IntegerLiteral struct method returns TokenLiteral as a string of IntegerLiteral.
*/
func (il *IntegerLiteral) TokenLiteral() string {
	return il.Token.Literal
}

/*
	This IntegerLiteral struct method to represents IntegerLiteral into string.
*/
func (il *IntegerLiteral) String() string {
	return il.Token.Literal
}

func (b *Boolean) TokenLiteral() string { return b.Token.Literal }

/*
	This Boolean struct method to represents Boolean into string.
*/
func (b *Boolean) String() string {
	return b.Token.Literal
}

/*
	This Identifier's struct method to represents Identifier into string.
*/
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

/*
	This Identifier's struct method to represents Identifier into string.
*/
func (i *Identifier) String() string { return i.Value }
