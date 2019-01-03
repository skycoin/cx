package actions

import (
	. "github.com/skycoin/cx/cx"
)

var PRGRM *CXProgram
var DataOffset int = STACK_SIZE + TYPE_POINTER_SIZE // to be able to handle nil pointers

var CurrentFile string
var LineNo int
var WebMode bool
var IdeMode bool
var WebPersistantMode bool
var BaseOutput bool
var ReplMode bool
var HelpMode bool
var InterpretMode bool
var CompileMode bool
var ReplTargetFn string = ""
var ReplTargetStrct string = ""
var ReplTargetMod string = ""

var FoundCompileErrors bool

var InREPL bool = false

var SysInitExprs []*CXExpression

var dStack bool = false
var InFn bool = false
var tag string = ""
var asmNL = "\n"
var fileName string

// to decide what shorthand op to use
const (
	OP_EQUAL = iota
	OP_UNEQUAL

	OP_BITAND
	OP_BITXOR
	OP_BITOR
	OP_BITCLEAR

	OP_MUL
	OP_DIV
	OP_MOD
	OP_ADD
	OP_SUB
	OP_BITSHL
	OP_BITSHR
	OP_LT
	OP_GT
	OP_LTEQ
	OP_GTEQ
)

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
