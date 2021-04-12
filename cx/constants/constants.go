package constants

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

const BOOL_SIZE = 1
const I8_SIZE = 1
const I16_SIZE = 2
const I32_SIZE = 4
const I64_SIZE = 8
const F32_SIZE = 4
const F64_SIZE = 8
const STR_SIZE = 4

const MARK_SIZE = 1
const OBJECT_HEADER_SIZE = 9
const OBJECT_GC_HEADER_SIZE = 5
const FORWARDING_ADDRESS_SIZE = 4
const OBJECT_SIZE = 4

const CALLSTACK_SIZE = 1000

var STACK_SIZE = 1048576     // 1 Mb
var INIT_HEAP_SIZE = 2097152 // 2 Mb
var MAX_HEAP_SIZE = 67108864 // 64 Mb
var MIN_HEAP_FREE_RATIO float32 = 0.4
var MAX_HEAP_FREE_RATIO float32 = 0.7

const NULL_HEAP_ADDRESS_OFFSET = 4
const NULL_HEAP_ADDRESS = 0
const STR_HEADER_SIZE = 4
const TYPE_POINTER_SIZE = 4
const SLICE_HEADER_SIZE = 8

const MAX_UINT32 = ^uint32(0)
const MIN_UINT32 = 0
const MAX_INT32 = int(MAX_UINT32 >> 1)
const MIN_INT32 = -MAX_INT32 - 1

var BASIC_TYPES []string = []string{
	"bool", "str", "i8", "i16", "i32", "i64", "ui8", "ui16", "ui32", "ui64", "f32", "f64",
	"[]bool", "[]str", "[]i8", "[]i16", "[]i32", "[]i64", "[]ui8", "[]ui16", "[]ui32", "[]ui64", "[]f32", "[]f64",
}

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

/*
grep -rn "PASSBY_VALUE" .
./cxparser/actions/functions.go:666:					out.PassBy = cxcore.PASSBY_VALUE
./cxparser/actions/functions.go:712:			out.PassBy = cxcore.PASSBY_VALUE
./cxparser/actions/functions.go:723:				assignElt.PassBy = cxcore.PASSBY_VALUE
./cx/config.go:172:	PASSBY_VALUE = iota
./cx/config.go:180:	PASSBY_VALUE //= iota
./cx/ast.go:1507:./cxparser/actions/functions.go:666:				out.PassBy = PASSBY_VALUE
./cx/ast.go:1509:./cxparser/actions/functions.go:712:			out.PassBy = PASSBY_VALUE
./cx/ast.go:1510:./cxparser/actions/functions.go:723:				assignElt.PassBy = PASSBY_VALUE
./cx/op_misc.go:37:		case PASSBY_VALUE:

grep -rn "PASSBY_REFERENCE" .
./cxparser/actions/misc.go:425:			arg.PassBy = cxcore.PASSBY_REFERENCE
./cxparser/actions/functions.go:678:		if elt.PassBy == cxcore.PASSBY_REFERENCE &&
./cxparser/actions/functions.go:915:			expr.Inputs[0].PassBy = cxcore.PASSBY_REFERENCE
./cxparser/actions/functions.go:1157:						nameFld.PassBy = cxcore.PASSBY_REFERENCE
./cxparser/actions/literals.go:219:				sym.PassBy = cxcore.PASSBY_REFERENCE
./cxparser/actions/expressions.go:336:		baseOut.PassBy = cxcore.PASSBY_REFERENCE
./cxparser/actions/assignment.go:57:		out.PassBy = cxcore.PASSBY_REFERENCE
./cxparser/actions/declarations.go:417:		arg.PassBy = cxcore.PASSBY_REFERENCE
./cxgo/actions/misc.go:425:			arg.PassBy = cxcore.PASSBY_REFERENCE
./cxgo/actions/functions.go:678:		if elt.PassBy == cxcore.PASSBY_REFERENCE &&
./cxgo/actions/functions.go:915:			expr.ProgramInput[0].PassBy = cxcore.PASSBY_REFERENCE
./cxgo/actions/functions.go:1157:						nameFld.PassBy = cxcore.PASSBY_REFERENCE
./cxgo/actions/literals.go:219:				sym.PassBy = cxcore.PASSBY_REFERENCE
./cxgo/actions/expressions.go:336:		baseOut.PassBy = cxcore.PASSBY_REFERENCE
./cxgo/actions/assignment.go:57:		out.PassBy = cxcore.PASSBY_REFERENCE
./cxgo/actions/declarations.go:417:		arg.PassBy = cxcore.PASSBY_REFERENCE
./cx/config.go:173:	PASSBY_REFERENCE
./cx/config.go:181:	PASSBY_REFERENCE
./cx/op_http.go:50:	headerFld.PassBy = PASSBY_REFERENCE
./cx/op_http.go:75:	transferEncodingFld.PassBy = PASSBY_REFERENCE
./cx/execute.go:442:						if inp.PassBy == PASSBY_REFERENCE {
./cx/ast.go:1506:./cxparser/actions/misc.go:425:			arg.PassBy = PASSBY_REFERENCE
./cx/ast.go:1508:./cxparser/actions/functions.go:678:		if elt.PassBy == PASSBY_REFERENCE &&
./cx/ast.go:1511:./cxparser/actions/functions.go:915:			expr.Inputs[0].PassBy = PASSBY_REFERENCE
./cx/ast.go:1513:./cxparser/actions/functions.go:1157:						nameFld.PassBy = PASSBY_REFERENCE
./cx/ast.go:1514:./cxparser/actions/literals.go:219:				sym.PassBy = PASSBY_REFERENCE
./cx/ast.go:1515:./cxparser/actions/expressions.go:336:		baseOut.PassBy = PASSBY_REFERENCE
./cx/ast.go:1516:./cxparser/actions/assignment.go:57:		out.PassBy = PASSBY_REFERENCE
./cx/ast.go:1525:./cxparser/actions/declarations.go:417:		arg.PassBy = PASSBY_REFERENCE
./cx/ast.go:1506:./cxgo/actions/misc.go:425:			arg.PassBy = PASSBY_REFERENCE
./cx/ast.go:1508:./cxgo/actions/functions.go:678:		if elt.PassBy == PASSBY_REFERENCE &&
./cx/ast.go:1511:./cxgo/actions/functions.go:915:			expr.ProgramInput[0].PassBy = PASSBY_REFERENCE
./cx/ast.go:1513:./cxgo/actions/functions.go:1157:						nameFld.PassBy = PASSBY_REFERENCE
./cx/ast.go:1514:./cxgo/actions/literals.go:219:				sym.PassBy = PASSBY_REFERENCE
./cx/ast.go:1515:./cxgo/actions/expressions.go:336:		baseOut.PassBy = PASSBY_REFERENCE
./cx/ast.go:1516:./cxgo/actions/assignment.go:57:		out.PassBy = PASSBY_REFERENCE
./cx/ast.go:1525:./cxgo/actions/declarations.go:417:		arg.PassBy = PASSBY_REFERENCE
./cx/ast.go:1528:./cx/op_http.go:50:	headerFld.PassBy = PASSBY_REFERENCE
./cx/ast.go:1529:./cx/op_http.go:75:	transferEncodingFld.PassBy = PASSBY_REFERENCE
./cx/ast.go:1533:./cx/execute.go:366:				if inp.PassBy == PASSBY_REFERENCE {
./cx/ast.go:1536:./cx/utilities.go:184:	if arg.PassBy == PASSBY_REFERENCE {
./cx/op_misc.go:39:		case PASSBY_REFERENCE:
./cx/utilities.go:182:	if arg.PassBy == PASSBY_REFERENCE {
*/

//ERROR: see below,
const (
	//PASSBY_UNUSED = iota
	PASSBY_VALUE = iota
	PASSBY_REFERENCE
)

/*
// massive problem
const (
	PASSBY_UNUSED = iota //if this value appears, program should crash; should be assert
	PASSBY_VALUE
	PASSBY_REFERENCE
)
*/

const (
	DEREF_UNUSED  = iota //reserve zero value, if this value appears, program should crash; should be assert
	DEREF_ARRAY          // 1
	DEREF_FIELD          // 2
	DEREF_POINTER        // 3
	DEREF_DEREF          // 4
	DEREF_SLICE          // 5
)

const (
	TPYE_UNUSED = iota //reserve zero value, if this value appears, program should crash; should be assert
	TYPE_UNDEFINED
	TYPE_AFF
	TYPE_BOOL
	TYPE_STR
	TYPE_F32
	TYPE_F64
	TYPE_I8
	TYPE_I16
	TYPE_I32
	TYPE_I64
	TYPE_UI8
	TYPE_UI16
	TYPE_UI32
	TYPE_UI64
	TYPE_FUNC

	TYPE_CUSTOM
	TYPE_POINTER
	TYPE_ARRAY
	TYPE_SLICE
	TYPE_IDENTIFIER
	TYPE_COUNT
)

var TypeCounter int
var TypeCodes map[string]int = map[string]int{
	"unused": TPYE_UNUSED, //if this appears, should be error; program should crash
	"ident":  TYPE_IDENTIFIER,
	"aff":    TYPE_AFF,
	"bool":   TYPE_BOOL,
	"str":    TYPE_STR,
	"f32":    TYPE_F32,
	"f64":    TYPE_F64,
	"i8":     TYPE_I8,
	"i16":    TYPE_I16,
	"i32":    TYPE_I32,
	"i64":    TYPE_I64,
	"ui8":    TYPE_UI8,
	"ui16":   TYPE_UI16,
	"ui32":   TYPE_UI32,
	"ui64":   TYPE_UI64,
	"und":    TYPE_UNDEFINED,
	"func":   TYPE_FUNC,
}

var TypeNames map[int]string = map[int]string{
	TPYE_UNUSED:     "UNUSED", //if this appears, should trigger asset; is error
	TYPE_IDENTIFIER: "ident",
	TYPE_AFF:        "aff",
	TYPE_BOOL:       "bool",
	TYPE_STR:        "str",
	TYPE_F32:        "f32",
	TYPE_F64:        "f64",
	TYPE_I8:         "i8",
	TYPE_I16:        "i16",
	TYPE_I32:        "i32",
	TYPE_I64:        "i64",
	TYPE_UI8:        "ui8",
	TYPE_UI16:       "ui16",
	TYPE_UI32:       "ui32",
	TYPE_UI64:       "ui64",
	TYPE_FUNC:       "func",
	TYPE_UNDEFINED:  "und",
}

// memory locations
const (
	MEM_UNUSED = iota //if this value appears, program should crash; should be assert
	MEM_STACK
	MEM_HEAP
	MEM_DATA
)

// GetArgSize ...
func GetArgSize(typ int) int {
	switch typ {
	case TYPE_BOOL, TYPE_I8, TYPE_UI8:
		return 1
	case TYPE_I16, TYPE_UI16:
		return 2
	case TYPE_STR, TYPE_I32, TYPE_UI32, TYPE_F32, TYPE_AFF:
		return 4
	case TYPE_I64, TYPE_UI64, TYPE_F64:
		return 8
	default:
		return 4
		//return -1 // should be panic
		//panic(CX_INTERNAL_ERROR)
	}
}
