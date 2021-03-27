package actions

import (
	"github.com/skycoin/cx/cx"
	"github.com/skycoin/cx/cx/constants"
)

//Move out actions/interactive to own module?

var PRGRM *cxcore.CXProgram
var DataOffset int = constants.STACK_SIZE //Heap Offset is Stack Size

//!!!
//Why cxcore.STACK_SIZE and not PRGRM.STACK_SIZE

var CurrentFile string
var LineNo int

const (
	SEL_RESERVED = iota
	SEL_ELSEIF
	SEL_ELSEIFELSE
)
