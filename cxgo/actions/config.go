package actions

import (
	. "github.com/SkycoinProject/cx/cx"
)

/* [2020 Jun 07 (ReewassSquared)] we should add verbose compilation options */

var PRGRM *CXProgram
var DataOffset int = STACK_SIZE

var CurrentFile string
var LineNo int
var ReplTargetFn = ""
var ReplTargetStrct = ""
var ReplTargetMod = ""

var SysInitExprs []*CXExpression

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
