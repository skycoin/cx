package actions

import (
	"github.com/skycoin/cx/cx"
)

/* [2020 Jun 07 (ReewassSquared)] we should add verbose compilation options */

var PRGRM *cxcore.CXProgram
var DataOffset int = cxcore.STACK_SIZE

var CurrentFile string
var LineNo int
var ReplTargetFn string = ""
var ReplTargetStrct string = ""
var ReplTargetMod string = ""

var SysInitExprs []*cxcore.CXExpression

// var dStack bool = false
var InFn bool = false

// var tag string = ""
// var asmNL = "\n"
// var fileName string

const (
	// type of selector
	SELECT_TYP_PKG = iota
	SELECT_TYP_FUNC
	SELECT_TYP_STRCT
)

const (
	SEL_ELSEIF = iota
	SEL_ELSEIFELSE
)
