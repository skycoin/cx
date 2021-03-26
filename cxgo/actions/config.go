package actions

import (
	"github.com/skycoin/cx/cx"
)

//Move out actions/interactive to own module?

var PRGRM *cxcore.CXProgram
var DataOffset int = cxcore.STACK_SIZE //Heap Offset is Stack Size

var CurrentFile string
var LineNo int

var SysInitExprs []*cxcore.CXExpression

var InFn bool = false

const (
	SEL_RESERVED = iota
	SEL_ELSEIF
	SEL_ELSEIFELSE
)
