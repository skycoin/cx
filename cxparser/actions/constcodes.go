package actions

import (
	"github.com/skycoin/cx/cx/constants"
	"github.com/skycoin/cx/cx/types"
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
	ConstCodes = map[string]int{}     //Todo: Used only once in actions/postfix.go
	Constants  = map[int]CXConstant{} //Todo: Only used once in actions/postfix.go
)

// nolint typecheck
func init() {
	// cx
	AddConstI32(CONST_CX_SUCCESS, "cx.SUCCESS", constants.CX_SUCCESS)
	AddConstI32(CONST_CX_COMPILATION_ERROR, "cx.COMPILATION_ERROR", constants.CX_COMPILATION_ERROR)
	AddConstI32(CONST_CX_PANIC, "cx.PANIC", constants.CX_PANIC)
	AddConstI32(CONST_CX_INTERNAL_ERROR, "cx.INTERNAL_ERROR", constants.CX_INTERNAL_ERROR)
	AddConstI32(CONST_CX_ASSERT, "cx.ASSERT", constants.CX_ASSERT)
	AddConstI32(CONST_CX_RUNTIME_ERROR, "cx.RUNTIME_ERROR", constants.CX_RUNTIME_ERROR)
	AddConstI32(CONST_CX_RUNTIME_STACK_OVERFLOW_ERROR, "cx.RUNTIME_STACK_OVERFLOW_ERROR", constants.CX_RUNTIME_STACK_OVERFLOW_ERROR)
	AddConstI32(CONST_CX_RUNTIME_HEAP_EXHAUSTED_ERROR, "cx.RUNTIME_HEAP_EXHAUSTED_ERROR", constants.CX_RUNTIME_HEAP_EXHAUSTED_ERROR)
	AddConstI32(CONST_CX_RUNTIME_INVALID_ARGUMENT, "cx.RUNTIME_INVALID_ARGUMENT", constants.CX_RUNTIME_INVALID_ARGUMENT)
	AddConstI32(CONST_CX_RUNTIME_SLICE_INDEX_OUT_OF_RANGE, "cx.RUNTIME_SLICE_INDEX_OUT_OF_RANGE", constants.CX_RUNTIME_SLICE_INDEX_OUT_OF_RANGE)
	AddConstI32(CONST_CX_RUNTIME_NOT_IMPLEMENTED, "cx.RUNTIME_NOT_INPLEMENTED", constants.CX_RUNTIME_NOT_IMPLEMENTED)
	AddConstI32(types.Cast_ptr_to_int(types.POINTER_SIZE), "cx.POINTER_SIZE", types.Cast_ptr_to_i32(types.POINTER_SIZE))
}

// AddConstCode ...
func AddConstCode(code int, name string, typeCode types.Code, value []byte) {
	ConstCodes[name] = code
	Constants[code] = CXConstant{Type: typeCode, Value: value}
}

// AddConstI32 ...
func AddConstI32(code int, name string, value int32) {
	var memory [4]byte
	types.Write_i32(memory[:], 0, value)
	AddConstCode(code, name, types.I32, memory[:])
}

// CXConstant ...
type CXConstant struct {
	// native constants. only used for pre-packaged constants (e.g. math package's PI)
	// these fields are used to feed WritePrimary
	Value []byte
	Type  types.Code
}
