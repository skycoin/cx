package cxparser

import (
	"bytes"

	cxparsingcompletor "github.com/skycoin/cx/cxparsergenerator/cxparsingcompletor"
	cxpartialparsing "github.com/skycoin/cx/cxparsergenerator/cxpartialparsing"
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
