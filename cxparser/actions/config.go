package actions

import (
	"github.com/skycoin/cx/cx/ast"
)

var AST *ast.CXProgram

var CurrentFile string
var LineNo int
var LineStr string

const (
	SEL_RESERVED = iota
	SEL_ELSEIF
	SEL_ELSEIFELSE
)
