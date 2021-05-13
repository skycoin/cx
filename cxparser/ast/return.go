package ast

import "github.com/skycoin/cx/cxparser/token"

/*
	ReturnStatement reprsents return  statement.
*/
type ReturnStatement struct {
	Token       token.Token // the 'return' token
	ReturnValue Expression
}
