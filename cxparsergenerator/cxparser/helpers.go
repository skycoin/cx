package cxparser

import (
	"bytes"

	cxparsingcompletor "github.com/skycoin/cx/cxparsergenerator/cxparsingcompletor"
	cxpartialparsing "github.com/skycoin/cx/cxparsergenerator/cxpartialparsing"
)

func passone(source string) int {

	parseErrors := cxpartialparsing.Parse(source)

	return parseErrors
}

func passtwo(b *bytes.Buffer) int {

	parseErrors := cxparsingcompletor.Parse(cxparsingcompletor.NewLexer(b))

	return parseErrors
}
