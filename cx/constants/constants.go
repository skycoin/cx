package constants

import (
	"github.com/skycoin/cx/cx/types"
)

// var COREPATH = ""

const STACK_OVERFLOW_ERROR = "stack overflow"
const HEAP_EXHAUSTED_ERROR = "heap exhausted"
const MAIN_FUNC = "main"
const SYS_INIT_FUNC = "*init"
const MAIN_PKG = "main"
const STDLIB_PKG = "stdlib"
const OS_PKG = "os"
const OS_ARGS = "Args"

const NON_ASSIGN_PREFIX = "nonAssign"
const LOCAL_PREFIX = "*tmp"
const LABEL_PREFIX = "*lbl"

// Used in `PrintProgram` to represent literals (`CXArgument`s with no name).
const LITERAL_PLACEHOLDER = "*lit"
const ID_FN = "identity"
const INIT_FN = "initDef"

const CALLSTACK_SIZE types.Pointer = 1000

var STACK_SIZE types.Pointer = 1048576     // 1 Mb
var INIT_HEAP_SIZE types.Pointer = 2097152 // 2 Mb
var MAX_HEAP_SIZE types.Pointer = 67108864 // 64 Mb
var MIN_HEAP_FREE_RATIO float32 = 0.4
var MAX_HEAP_FREE_RATIO float32 = 0.7

const NULL_HEAP_ADDRESS_OFFSET = types.POINTER_SIZE
const NULL_HEAP_ADDRESS = types.Pointer(0)
const SLICE_HEADER_SIZE = 2 * types.POINTER_SIZE

//VERY WEIRD
//gives error, "cx" not found, even if it exists when changed

/*
grep -rn "PARAM_DEFAULT" .
./cx/config.go:87:	PARAM_DEFAULT = iota
./cx/config.go:96:	PARAM_DEFAULT
./cx/opcodes.go:843:	case PARAM_DEFAULT:

grep -rn "PARAM_SLICE" .
./cx/config.go:88:	PARAM_SLICE
./cx/config.go:97:	PARAM_SLICE
./cx/opcodes.go:845:	case PARAM_SLICE:

grep -rn "PARAM_STRUCT" .
./cx/config.go:89:	PARAM_STRUCT
./cx/config.go:98:	PARAM_STRUCT
./cx/opcodes.go:847:	case PARAM_STRUCT:

*/

//works
// BUG
const (
	PARAM_DEFAULT = iota
	PARAM_SLICE
	PARAM_STRUCT
)

/*
//doesnt work
const (
	PARAM_UNUSED = iota
	PARAM_DEFAULT
	PARAM_SLICE
	PARAM_STRUCT
)
*/

const (
	CX_SUCCESS = iota //zero can be success
	CX_COMPILATION_ERROR
	CX_PANIC // 2
	CX_INTERNAL_ERROR
	CX_ASSERT
	CX_RUNTIME_ERROR
	CX_RUNTIME_STACK_OVERFLOW_ERROR
	CX_RUNTIME_HEAP_EXHAUSTED_ERROR
	CX_RUNTIME_INVALID_ARGUMENT
	CX_RUNTIME_SLICE_INDEX_OUT_OF_RANGE
	CX_RUNTIME_NOT_IMPLEMENTED
	CX_RUNTIME_INVALID_CAST
)

var ErrorStrings map[int]string = map[int]string{
	CX_SUCCESS:                          "CX_SUCCESS",
	CX_COMPILATION_ERROR:                "CX_COMPILATION_ERROR",
	CX_PANIC:                            "CX_PANIC",
	CX_INTERNAL_ERROR:                   "CX_INTERNAL_ERROR",
	CX_ASSERT:                           "CX_ASSERT",
	CX_RUNTIME_ERROR:                    "CX_RUNTIME_ERROR",
	CX_RUNTIME_STACK_OVERFLOW_ERROR:     "CX_RUNTIME_STACK_OVERFLOW_ERROR",
	CX_RUNTIME_HEAP_EXHAUSTED_ERROR:     "CX_RUNTIME_HEAP_EXHAUSTED_ERROR",
	CX_RUNTIME_INVALID_ARGUMENT:         "CX_RUNTIME_INVALID_ARGUMENT",
	CX_RUNTIME_SLICE_INDEX_OUT_OF_RANGE: "CX_RUNTIME_SLICE_INDEX_OUT_OF_RANGE",
	CX_RUNTIME_NOT_IMPLEMENTED:          "CX_RUNTIME_NOT_IMPLEMENTED",
	CX_RUNTIME_INVALID_CAST:             "CX_RUNTIME_INVALID_CAST",
}

const (
	DECL_UNUSED   = iota //if this value appears, program should crash; is error
	DECL_POINTER         // 1
	DECL_DEREF           // 2
	DECL_ARRAY           // 3
	DECL_SLICE           // 4
	DECL_STRUCT          // 5
	DECL_INDEXING        // 6
	DECL_BASIC           // 7
	DECL_FUNC            // 8
)

//ERROR: see below,
const (
	//PASSBY_UNUSED = iota
	PASSBY_VALUE = iota
	PASSBY_REFERENCE
)

const (
	DEREF_UNUSED  = iota //reserve zero value, if this value appears, program should crash; should be assert
	DEREF_ARRAY          // 1
	DEREF_FIELD          // 2
	DEREF_POINTER        // 3
	DEREF_DEREF          // 4
	DEREF_SLICE          // 5
)
