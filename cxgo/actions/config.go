package actions

import (
	"github.com/skycoin/cx/cx/ast"
)

//Move out actions/interactive to own module?

var AST *ast.CXProgram

//TODO: THIS IS WRONG
//USE AST.HeapStartsAt

// var DataOffset int = constants.STACK_SIZE //Heap Offset is Stack Size

//var DataOffset int = constants.STACK_SIZE //Heap Offset is Stack Size

//!!!
//Why cxcore.STACK_SIZE and not AST.STACK_SIZE

var CurrentFile string
var LineNo int

const (
	SEL_RESERVED = iota
	SEL_ELSEIF
	SEL_ELSEIFELSE
)
