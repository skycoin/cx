package cxcore

import (
	"github.com/skycoin/cx/cx/constants"
	"github.com/skycoin/cx/cx/helper"
)

// constant codes
// nolint golint
const (
	// cx
	CONST_CX_SUCCESS = iota
	CONST_CX_COMPILATION_ERROR
	CONST_CX_PANIC
	CONST_CX_INTERNAL_ERROR
	CONST_CX_ASSERT
	CONST_CX_RUNTIME_ERROR
	CONST_CX_RUNTIME_STACK_OVERFLOW_ERROR
	CONST_CX_RUNTIME_HEAP_EXHAUSTED_ERROR
	CONST_CX_RUNTIME_INVALID_ARGUMENT
	CONST_CX_RUNTIME_SLICE_INDEX_OUT_OF_RANGE
	CONST_CX_RUNTIME_NOT_IMPLEMENTED
)

// For the cxgo. These shouldn't be used in the runtime for performance reasons
var (
	ConstNames = map[int]string{}
	ConstCodes = map[string]int{}
	Constants  = map[int]CXConstant{}
)

// AddConstCode ...
func AddConstCode(code int, name string, typ int, value []byte) {
	ConstNames[code] = name
	ConstCodes[name] = code
	Constants[code] = CXConstant{Type: typ, Value: value}
}

// ConstI32 ...
func ConstI32(code int, name string, value int32) {
	AddConstCode(code, name, constants.TYPE_I32, helper.FromI32(value))
}

// nolint typecheck
func init() {
	// cx
	ConstI32(CONST_CX_SUCCESS, "cx.SUCCESS", constants.CX_SUCCESS)
	ConstI32(CONST_CX_COMPILATION_ERROR, "cx.COMPILATION_ERROR", constants.CX_COMPILATION_ERROR)
	ConstI32(CONST_CX_PANIC, "cx.PANIC", constants.CX_PANIC)
	ConstI32(CONST_CX_INTERNAL_ERROR, "cx.INTERNAL_ERROR", constants.CX_INTERNAL_ERROR)
	ConstI32(CONST_CX_ASSERT, "cx.ASSERT", constants.CX_ASSERT)
	ConstI32(CONST_CX_RUNTIME_ERROR, "cx.RUNTIME_ERROR", constants.CX_RUNTIME_ERROR)
	ConstI32(CONST_CX_RUNTIME_STACK_OVERFLOW_ERROR, "cx.RUNTIME_STACK_OVERFLOW_ERROR", constants.CX_RUNTIME_STACK_OVERFLOW_ERROR)
	ConstI32(CONST_CX_RUNTIME_HEAP_EXHAUSTED_ERROR, "cx.RUNTIME_HEAP_EXHAUSTED_ERROR", constants.CX_RUNTIME_HEAP_EXHAUSTED_ERROR)
	ConstI32(CONST_CX_RUNTIME_INVALID_ARGUMENT, "cx.RUNTIME_INVALID_ARGUMENT", constants.CX_RUNTIME_INVALID_ARGUMENT)
	ConstI32(CONST_CX_RUNTIME_SLICE_INDEX_OUT_OF_RANGE, "cx.RUNTIME_SLICE_INDEX_OUT_OF_RANGE", constants.CX_RUNTIME_SLICE_INDEX_OUT_OF_RANGE)
	ConstI32(CONST_CX_RUNTIME_NOT_IMPLEMENTED, "cx.RUNTIME_NOT_INPLEMENTED", constants.CX_RUNTIME_NOT_IMPLEMENTED)
}
