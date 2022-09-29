package ast

import (
	"bytes"
	"strings"

	"github.com/skycoin/cx/cxparser/token"
)

type ArrayLiteral struct {
	Token    token.Token
	Elements []Expression
}

/*
	This ArrayLiteral struct method returns TokenLiteral as a string of ArrayLiteral.
*/
func (al *ArrayLiteral) TokenLiteral() string {

	return al.Token.Literal
}

/*
	This ArrayLiteral's struct method to represents ArrayLiteral into string.
*/
func (al *ArrayLiteral) String() string {

	var out bytes.Buffer

	elements := []string{}

	for _, el := range al.Elements {
		elements = append(elements, el.String())
	}

	out.WriteString("[")

	out.WriteString(strings.Join(elements, ", "))

	out.WriteString("]")

	return out.String()
}

type IndexExpression struct {
	Token token.Token // '[' token
	Left  Expression
	Index Expression
}

/*
	This ArrayLiteral struct method returns TokenLiteral as a string of ArrayLiteral.
*/
func (ie *IndexExpression) TokenLiteral() string {

	return ie.Token.Literal
}

/*
	This IndexExpression's struct method to represents IndexExpression into string.
*/

func (ie *IndexExpression) String() string {

	var out bytes.Buffer

	out.WriteString("(")

	out.WriteString(ie.Left.String())

	out.WriteString("[")

	out.WriteString(ie.Index.String())

	out.WriteString("]")

	return out.String()
}
