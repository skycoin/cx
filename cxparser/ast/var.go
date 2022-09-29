package ast

import "github.com/skycoin/cx/cxparser/token"

/*
	VarStatement reprsents var statment except global variable.
*/
type VarStatement struct {
	Token token.Token // the token.VAR token
	Name  *Identifier
	Value Expression
}
