package cxparsering

import (
	"bytes"

	cxparsingcompletor "github.com/skycoin/cx/cxparser/cxparsingcompletor"
	cxpartialparsing "github.com/skycoin/cx/cxparser/cxpartialparsing"
)

/*
passone perform cxpartial parsing using grammmer.y
*/
func passone(source string) int {

	parseErrors := cxpartialparsing.Parse(source)

	return parseErrors
}

/*
passtwo perform cxparsingcompletor parsing using lexer.y
*/

func passtwo(b *bytes.Buffer) int {

	parseErrors := cxparsingcompletor.Parse(cxparsingcompletor.NewLexer(b))

	return parseErrors
}
