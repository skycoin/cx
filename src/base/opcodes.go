package base

// op codes
const (
	OP_IDENTITY = iota
	OP_JMP
	
	OP_I32_ADD
	OP_I32_SUB
	OP_I32_MUL
	OP_I32_DIV
	OP_I32_ABS
	OP_I32_MOD
	OP_I32_POW

	OP_I64_ADD
	OP_I64_SUB
	OP_I64_MUL
	OP_I64_DIV
	OP_I64_ABS
	OP_I64_MOD
	OP_I64_POW
	
	OP_F32_SIN
	OP_F32_COS
	
	OP_F64_SIN
	OP_F64_COS
	
	OP_BITAND
	OP_BITOR
	OP_BITXOR
	OP_BITCLEAR
	OP_BITSHL
	OP_BITSHR
	OP_I32_PRINT
	OP_I64_PRINT
	OP_BOOL_PRINT
	OP_MAKE
	OP_READ
	OP_WRITE
	OP_LEN
	OP_CONCAT
	OP_APPEND
	OP_COPY
	OP_CAST
	OP_EQ
	OP_UNEQ
	OP_I32_LT
	OP_I32_GT
	OP_I32_LTEQ
	OP_I32_GTEQ
	OP_RAND
	OP_AND
	OP_OR
	OP_NOT
	OP_SLEEP
	OP_HALT
	OP_GOTO
	OP_REMCX
	OP_ADDCX
	OP_QUERY
	OP_EXECUTE
	OP_INDEX
	OP_NAME
	OP_EVOLVE
	OP_TEST_START
	OP_TEST_STOP
	OP_TEST_ERROR
	OP_TEST
	OP_TIME_UNIX
	OP_TIME_UNIX_MILLI
	OP_TIME_UNIX_NANO
)

func execNative (prgrm *CXProgram) {
	call := &prgrm.CallStack[prgrm.CallCounter]
	stack := &prgrm.Stacks[0]
	expr := call.Operator.Expressions[call.Line]
	opCode := expr.Operator.OpCode
	fp := call.FramePointer
	
	switch opCode {
	case OP_IDENTITY: identity(expr, stack, fp)
	case OP_JMP: jmp(expr, stack, fp, call)
		
	case OP_I32_ADD: i32_add(expr, stack, fp)
	case OP_I32_SUB: i32_sub(expr, stack, fp)
	case OP_I32_MUL:
	case OP_I32_DIV:
	case OP_I32_ABS:
	case OP_I32_MOD:
	case OP_I32_POW:

	case OP_I64_ADD: i64_add(expr, stack, fp)
	case OP_I64_SUB: i64_sub(expr, stack, fp)
	case OP_I64_MUL:
	case OP_I64_DIV:
	case OP_I64_ABS:
	case OP_I64_MOD:
	case OP_I64_POW:
		
	case OP_F32_COS:
	case OP_F32_SIN:
		
	case OP_BITAND:
	case OP_BITOR:
	case OP_BITXOR:
	case OP_BITCLEAR:
	case OP_BITSHL:
	case OP_BITSHR:
	case OP_I32_PRINT: i32_print(expr, stack, fp)
	case OP_I64_PRINT: i64_print(expr, stack, fp)
	case OP_BOOL_PRINT: bool_print(expr, stack, fp)
	case OP_MAKE:
	case OP_READ:
	case OP_WRITE:
	case OP_LEN:
	case OP_CONCAT:
	case OP_APPEND:
	case OP_COPY:
	case OP_CAST:
	case OP_EQ:
	case OP_UNEQ:
	case OP_I32_LT: i32_lt(expr, stack, fp)
	case OP_I32_GT: i32_gt(expr, stack, fp)
	case OP_I32_LTEQ:
	case OP_I32_GTEQ:
	case OP_RAND:
	case OP_AND:
	case OP_OR:
	case OP_NOT:
	case OP_SLEEP:
	case OP_HALT:
	case OP_GOTO: goTo(expr, call)
	case OP_REMCX:
	case OP_ADDCX:
	case OP_QUERY:
	case OP_EXECUTE:
	case OP_INDEX:
	case OP_NAME:
	case OP_EVOLVE:
	case OP_TEST_START:
	case OP_TEST_STOP:
	case OP_TEST_ERROR:
	case OP_TEST:
	case OP_TIME_UNIX:
	case OP_TIME_UNIX_MILLI: time_UnixMilli(expr, stack, fp)
	case OP_TIME_UNIX_NANO:
	}
}

// For the parser. These shouldn't be used in the runtime for performance reasons
var OpNames map[int]string = map[int]string{
	OP_IDENTITY: "identity",
	OP_JMP: "jmp",
	OP_I32_ADD: "i32.add",
	OP_I32_SUB: "i32.sub",
	OP_I32_PRINT: "i32.print",
	OP_I32_GT: "i32.gt",
	OP_I32_LT: "i32.lt",

	OP_I64_ADD: "i64.add",
	OP_I64_SUB: "i64.sub",
	OP_I64_PRINT: "i64.print",

	OP_TIME_UNIX_MILLI: "time.UnixMilli",

	OP_BOOL_PRINT: "bool.print",
}

// For the parser. These shouldn't be used in the runtime for performance reasons
var OpCodes map[string]int = map[string]int{
	"identity": OP_IDENTITY,
	"jmp": OP_JMP,
	"i32.add": OP_I32_ADD,
	"i32.sub": OP_I32_SUB,
	"i32.print": OP_I32_PRINT,
	"i32.gt": OP_I32_GT,
	"i32.lt": OP_I32_LT,

	"i64.add": OP_I64_ADD,
	"i64.sub": OP_I64_SUB,
	"i64.print": OP_I64_PRINT,

	"time.UnixMilli": OP_TIME_UNIX_MILLI,

	"bool.print": OP_BOOL_PRINT,
}

var Natives map[int]*CXFunction = map[int]*CXFunction{
	OP_IDENTITY: MakeNative(OP_IDENTITY, []int{TYPE_UNDEFINED}, []int{TYPE_UNDEFINED}),
	OP_I32_ADD: MakeNative(OP_I32_ADD, []int{TYPE_I32, TYPE_I32}, []int{TYPE_I32}),
	OP_I32_SUB: MakeNative(OP_I32_SUB, []int{TYPE_I32, TYPE_I32}, []int{TYPE_I32}),
	OP_I32_PRINT: MakeNative(OP_I32_PRINT, []int{TYPE_I32}, []int{}),
	OP_I32_LT: MakeNative(OP_I32_LT, []int{TYPE_I32, TYPE_I32}, []int{TYPE_BOOL}),
	OP_I32_GT: MakeNative(OP_I32_GT, []int{TYPE_I32, TYPE_I32}, []int{TYPE_BOOL}),

	OP_I64_ADD: MakeNative(OP_I64_ADD, []int{TYPE_I64, TYPE_I64}, []int{TYPE_I64}),
	OP_I64_SUB: MakeNative(OP_I64_SUB, []int{TYPE_I64, TYPE_I64}, []int{TYPE_I64}),
	OP_I64_PRINT: MakeNative(OP_I64_PRINT, []int{TYPE_I64}, []int{}),

	OP_BOOL_PRINT: MakeNative(OP_BOOL_PRINT, []int{TYPE_I32}, []int{}),

	OP_TIME_UNIX_MILLI: MakeNative(OP_TIME_UNIX_MILLI, []int{}, []int{TYPE_I64}),
	OP_JMP: MakeNative(OP_JMP, []int{TYPE_BOOL, TYPE_I32, TYPE_I32}, []int{}),
}
