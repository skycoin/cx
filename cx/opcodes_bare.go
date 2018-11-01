package base

import (
	
)

var CorePackages = []string{
	// temporary solution until we can implement these packages in pure CX I guess
	"gl", "glfw", "time", "http", "os", "explorer", "aff", "gltext", "serial",
}

// op codes
const (
	OP_IDENTITY = iota
	OP_JMP
	OP_DEBUG

	OP_UND_EQUAL
	OP_UND_UNEQUAL
	OP_UND_BITAND
	OP_UND_BITXOR
	OP_UND_BITOR
	OP_UND_BITCLEAR
	OP_UND_MUL
	OP_UND_DIV
	OP_UND_MOD
	OP_UND_ADD
	OP_UND_SUB
	OP_UND_BITSHL
	OP_UND_BITSHR
	OP_UND_LT
	OP_UND_GT
	OP_UND_LTEQ
	OP_UND_GTEQ
	OP_UND_LEN
	OP_UND_PRINTF
	OP_UND_SPRINTF
	OP_UND_READ

	OP_BOOL_PRINT

	OP_BOOL_EQUAL
	OP_BOOL_UNEQUAL
	OP_BOOL_NOT
	OP_BOOL_OR
	OP_BOOL_AND

	OP_BYTE_BYTE
	OP_BYTE_STR
	OP_BYTE_I32
	OP_BYTE_I64
	OP_BYTE_F32
	OP_BYTE_F64
	
	OP_BYTE_PRINT

	OP_I32_BYTE
	OP_I32_STR
	OP_I32_I32
	OP_I32_I64
	OP_I32_F32
	OP_I32_F64
	
	OP_I32_PRINT
	OP_I32_ADD
	OP_I32_SUB
	OP_I32_MUL
	OP_I32_DIV
	OP_I32_ABS
	OP_I32_POW
	OP_I32_GT
	OP_I32_GTEQ
	OP_I32_LT
	OP_I32_LTEQ
	OP_I32_EQ
	OP_I32_UNEQ
	OP_I32_MOD
	OP_I32_RAND
	OP_I32_BITAND
	OP_I32_BITOR
	OP_I32_BITXOR
	OP_I32_BITCLEAR
	OP_I32_BITSHL
	OP_I32_BITSHR
	OP_I32_SQRT
	OP_I32_LOG
	OP_I32_LOG2
	OP_I32_LOG10

	OP_I32_MAX
	OP_I32_MIN

	OP_I64_BYTE
	OP_I64_STR
	OP_I64_I32
	OP_I64_I64
	OP_I64_F32
	OP_I64_F64

	OP_I64_PRINT
	OP_I64_ADD
	OP_I64_SUB
	OP_I64_MUL
	OP_I64_DIV
	OP_I64_ABS
	OP_I64_POW
	OP_I64_GT
	OP_I64_GTEQ
	OP_I64_LT
	OP_I64_LTEQ
	OP_I64_EQ
	OP_I64_UNEQ
	OP_I64_MOD
	OP_I64_RAND
	OP_I64_BITAND
	OP_I64_BITOR
	OP_I64_BITXOR
	OP_I64_BITCLEAR
	OP_I64_BITSHL
	OP_I64_BITSHR
	OP_I64_SQRT
	OP_I64_LOG
	OP_I64_LOG10
	OP_I64_LOG2
	OP_I64_MAX
	OP_I64_MIN

	OP_F32_BYTE
	OP_F32_STR
	OP_F32_I32
	OP_F32_I64
	OP_F32_F32
	OP_F32_F64
	
	OP_F32_PRINT
	OP_F32_ADD
	OP_F32_SUB
	OP_F32_MUL
	OP_F32_DIV
	OP_F32_ABS
	OP_F32_POW
	OP_F32_GT
	OP_F32_GTEQ
	OP_F32_LT
	OP_F32_LTEQ
	OP_F32_EQ
	OP_F32_UNEQ
	OP_F32_COS
	OP_F32_SIN
	OP_F32_SQRT
	OP_F32_LOG
	OP_F32_LOG2
	OP_F32_LOG10
	OP_F32_MAX
	OP_F32_MIN

	OP_F64_BYTE
	OP_F64_STR
	OP_F64_I32
	OP_F64_I64
	OP_F64_F32
	OP_F64_F64

	OP_F64_PRINT
	OP_F64_ADD
	OP_F64_SUB
	OP_F64_MUL
	OP_F64_DIV
	OP_F64_ABS
	OP_F64_POW
	OP_F64_GT
	OP_F64_GTEQ
	OP_F64_LT
	OP_F64_LTEQ
	OP_F64_EQ
	OP_F64_UNEQ
	OP_F64_COS
	OP_F64_SIN

	OP_F64_SQRT
	OP_F64_LOG
	OP_F64_LOG2
	OP_F64_LOG10
	OP_F64_MAX
	OP_F64_MIN

	OP_STR_PRINT
	OP_STR_CONCAT
	OP_STR_EQ

	OP_STR_BYTE
	OP_STR_STR
	OP_STR_I32
	OP_STR_I64
	OP_STR_F32
	OP_STR_F64
	
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

	OP_ASSERT

	// affordances
	OP_AFF_PRINT
	OP_AFF_QUERY
	OP_AFF_ON
	OP_AFF_OF
	OP_AFF_INFORM
	OP_AFF_REQUEST

	END_OF_BARE_OPS
)

// For the parser. These shouldn't be used in the runtime for performance reasons
var OpNames map[int]string
var OpCodes map[string]int
var Natives map[int]*CXFunction
var execNative func(*CXProgram)

func init () {
	OpNames = map[int]string{
		OP_IDENTITY:   "identity",
		OP_JMP:        "jmp",
		OP_DEBUG:      "debug",

		OP_UND_EQUAL:    "eq",
		OP_UND_UNEQUAL:  "uneq",
		OP_UND_BITAND:   "bitand",
		OP_UND_BITXOR:   "bitxor",
		OP_UND_BITOR:    "bitor",
		OP_UND_BITCLEAR: "bitclear",
		OP_UND_MUL:      "mul",
		OP_UND_DIV:      "div",
		OP_UND_MOD:      "mod",
		OP_UND_ADD:      "add",
		OP_UND_SUB:      "sub",
		OP_UND_BITSHL:   "bitshl",
		OP_UND_LT:       "lt",
		OP_UND_GT:       "gt",
		OP_UND_LTEQ:     "lteq",
		OP_UND_GTEQ:     "gteq",
		OP_UND_LEN:      "len",
		OP_UND_PRINTF:   "printf",
		OP_UND_SPRINTF:  "sprintf",
		OP_UND_READ:     "read",

		OP_BYTE_BYTE: "byte.byte",
		OP_BYTE_STR: "byte.str",
		OP_BYTE_I32: "byte.i32",
		OP_BYTE_I64: "byte.i64",
		OP_BYTE_F32: "byte.f32",
		OP_BYTE_F64: "byte.f64",
		
		OP_BYTE_PRINT: "byte.print",

		OP_BOOL_PRINT:   "bool.print",
		OP_BOOL_EQUAL:   "bool.eq",
		OP_BOOL_UNEQUAL: "bool.uneq",
		OP_BOOL_NOT:     "bool.not",
		OP_BOOL_OR:      "bool.or",
		OP_BOOL_AND:     "bool.and",
		
		OP_I32_BYTE:     "i32.byte",
		OP_I32_STR:      "i32.str",
		OP_I32_I32:      "i32.i32",
		OP_I32_I64:      "i32.i64",
		OP_I32_F32:      "i32.f32",
		OP_I32_F64:      "i32.f64",

		OP_I32_PRINT:    "i32.print",
		OP_I32_ADD:      "i32.add",
		OP_I32_SUB:      "i32.sub",
		OP_I32_MUL:      "i32.mul",
		OP_I32_DIV:      "i32.div",
		OP_I32_ABS:      "i32.abs",
		OP_I32_POW:      "i32.pow",
		OP_I32_GT:       "i32.gt",
		OP_I32_GTEQ:     "i32.gteq",
		OP_I32_LT:       "i32.lt",
		OP_I32_LTEQ:     "i32.lteq",
		OP_I32_EQ:       "i32.eq",
		OP_I32_UNEQ:     "i32.uneq",
		OP_I32_MOD:      "i32.mod",
		OP_I32_RAND:     "i32.rand",
		OP_I32_BITAND:   "i32.bitand",
		OP_I32_BITOR:    "i32.bitor",
		OP_I32_BITXOR:   "i32.bitxor",
		OP_I32_BITCLEAR: "i32.bitclear",
		OP_I32_BITSHL:   "i32.bitshl",
		OP_I32_BITSHR:   "i32.bitshr",
		OP_I32_SQRT:     "i32.sqrt",
		OP_I32_LOG:      "i32.log",
		OP_I32_LOG2:     "i32.log2",
		OP_I32_LOG10:    "i32.log10",
		OP_I32_MAX:      "i32.max",
		OP_I32_MIN:      "i32.min",

		OP_I64_BYTE:     "i64.byte",
		OP_I64_STR:      "i64.str",
		OP_I64_I32:      "i64.i32",
		OP_I64_I64:      "i64.i64",
		OP_I64_F32:      "i64.f32",
		OP_I64_F64:      "i64.f64",
		
		OP_I64_PRINT:    "i64.print",
		OP_I64_ADD:      "i64.add",
		OP_I64_SUB:      "i64.sub",
		OP_I64_MUL:      "i64.mul",
		OP_I64_DIV:      "i64.div",
		OP_I64_ABS:      "i64.abs",
		OP_I64_POW:      "i64.pow",
		OP_I64_GT:       "i64.gt",
		OP_I64_GTEQ:     "i64.gteq",
		OP_I64_LT:       "i64.lt",
		OP_I64_LTEQ:     "i64.lteq",
		OP_I64_EQ:       "i64.eq",
		OP_I64_UNEQ:     "i64.uneq",
		OP_I64_MOD:      "i64.mod",
		OP_I64_RAND:     "i64.rand",
		OP_I64_BITAND:   "i64.bitand",
		OP_I64_BITOR:    "i64.bitor",
		OP_I64_BITXOR:   "i64.bitxor",
		OP_I64_BITCLEAR: "i64.bitclear",
		OP_I64_BITSHL:   "i64.bitshl",
		OP_I64_BITSHR:   "i64.bitshr",
		OP_I64_SQRT:     "i64.sqrt",
		OP_I64_LOG:      "i64.log",
		OP_I64_LOG2:     "i64.log2",
		OP_I64_LOG10:    "i64.log10",
		OP_I64_MAX:      "i64.max",
		OP_I64_MIN:      "i64.min",
		
		OP_F32_BYTE:     "f32.byte",
		OP_F32_STR:      "f32.str",
		OP_F32_I32:      "f32.i32",
		OP_F32_I64:      "f32.i64",
		OP_F32_F32:      "f32.f32",
		OP_F32_F64:      "f32.f64",
		
		OP_F32_PRINT:    "f32.print",
		OP_F32_ADD:      "f32.add",
		OP_F32_SUB:      "f32.sub",
		OP_F32_MUL:      "f32.mul",
		OP_F32_DIV:      "f32.div",
		OP_F32_ABS:      "f32.abs",
		OP_F32_POW:      "f32.pow",
		OP_F32_GT:       "f32.gt",
		OP_F32_GTEQ:     "f32.gteq",
		OP_F32_LT:       "f32.lt",
		OP_F32_LTEQ:     "f32.lteq",
		OP_F32_EQ:       "f32.eq",
		OP_F32_UNEQ:     "f32.uneq",
		OP_F32_COS:      "f32.cos",
		OP_F32_SIN:      "f32.sin",
		OP_F32_SQRT:     "f32.sqrt",
		OP_F32_LOG:      "f32.log",
		OP_F32_LOG2:     "f32.log2",
		OP_F32_LOG10:    "f32.log10",
		OP_F32_MAX:      "f32.max",
		OP_F32_MIN:      "f32.min",

		OP_F64_BYTE:    "f64.byte",
		OP_F64_STR:    "f64.str",
		OP_F64_I32:    "f64.i32",
		OP_F64_I64:    "f64.i64",
		OP_F64_F32:    "f64.f32",
		OP_F64_F64:    "f64.f64",
		
		OP_F64_PRINT:    "f64.print",
		OP_F64_ADD:      "f64.add",
		OP_F64_SUB:      "f64.sub",
		OP_F64_MUL:      "f64.mul",
		OP_F64_DIV:      "f64.div",
		OP_F64_ABS:      "f64.abs",
		OP_F64_POW:      "f64.pow",
		OP_F64_GT:       "f64.gt",
		OP_F64_GTEQ:     "f64.gteq",
		OP_F64_LT:       "f64.lt",
		OP_F64_LTEQ:     "f64.lteq",
		OP_F64_EQ:       "f64.eq",
		OP_F64_UNEQ:     "f64.uneq",
		OP_F64_COS:      "f64.cos",
		OP_F64_SIN:      "f64.sin",
		OP_F64_SQRT:     "f64.sqrt",
		OP_F64_LOG:      "f64.log",
		OP_F64_LOG2:     "f64.log2",
		OP_F64_LOG10:    "f64.log10",
		OP_F64_MAX:      "f64.max",
		OP_F64_MIN:      "f64.min",

		OP_STR_PRINT: "str.print",
		OP_STR_CONCAT: "str.concat",
		OP_STR_EQ: "str.eq",

		OP_STR_BYTE: "str.byte",
		OP_STR_STR: "str.str",
		OP_STR_I32: "str.i32",
		OP_STR_I64: "str.i64",
		OP_STR_F32: "str.f32",
		OP_STR_F64: "str.f64",

		OP_APPEND: "append",
		OP_ASSERT: "assert",

		// affordances
		OP_AFF_PRINT: "aff.print",
		OP_AFF_QUERY: "aff.query",
		OP_AFF_ON: "aff.on",
		OP_AFF_OF: "aff.of",
		OP_AFF_INFORM: "aff.inform",
		OP_AFF_REQUEST: "aff.request",
	}

	OpCodes = map[string]int{
		"identity": OP_IDENTITY,
		"jmp":      OP_JMP,
		"debug":    OP_DEBUG,

		"eq":       OP_UND_EQUAL,
		"uneq":     OP_UND_UNEQUAL,
		"bitand":   OP_UND_BITAND,
		"bitxor":   OP_UND_BITXOR,
		"bitor":    OP_UND_BITOR,
		"bitclear": OP_UND_BITCLEAR,
		"mul":      OP_UND_MUL,
		"div":      OP_UND_DIV,
		"mod":      OP_UND_MOD,
		"add":      OP_UND_ADD,
		"sub":      OP_UND_SUB,
		"bitshl":   OP_UND_BITSHL,
		"bitshr":   OP_UND_BITSHR,
		"lt":       OP_UND_LT,
		"gt":       OP_UND_GT,
		"lteq":     OP_UND_LTEQ,
		"gteq":     OP_UND_GTEQ,
		"len":      OP_UND_LEN,
		"printf":   OP_UND_PRINTF,
		"sprintf":  OP_UND_SPRINTF,
		"read":     OP_UND_READ,

		"byte.byte":  OP_BYTE_BYTE,
		"byte.str":   OP_BYTE_STR,
		"byte.i32":   OP_BYTE_I32,
		"byte.i64":   OP_BYTE_I64,
		"byte.f32":   OP_BYTE_F32,
		"byte.f64":   OP_BYTE_F64,
		
		"byte.print": OP_BYTE_PRINT,

		"bool.print": OP_BOOL_PRINT,
		"bool.eq":    OP_BOOL_EQUAL,
		"bool.uneq":  OP_BOOL_UNEQUAL,
		"bool.not":   OP_BOOL_NOT,
		"bool.or":    OP_BOOL_OR,
		"bool.and":   OP_BOOL_AND,

		"i32.byte":     OP_I32_BYTE,
		"i32.str":      OP_I32_STR,
		"i32.i32":      OP_I32_I32,
		"i32.i64":      OP_I32_I64,
		"i32.f32":      OP_I32_F32,
		"i32.f64":      OP_I32_F64,
		
		"i32.print":    OP_I32_PRINT,
		"i32.add":      OP_I32_ADD,
		"i32.sub":      OP_I32_SUB,
		"i32.mul":      OP_I32_MUL,
		"i32.div":      OP_I32_DIV,
		"i32.abs":      OP_I32_ABS,
		"i32.pow":      OP_I32_POW,
		"i32.gt":       OP_I32_GT,
		"i32.gteq":     OP_I32_GTEQ,
		"i32.lt":       OP_I32_LT,
		"i32.lteq":     OP_I32_LTEQ,
		"i32.eq":       OP_I32_EQ,
		"i32.uneq":     OP_I32_UNEQ,
		"i32.mod":      OP_I32_MOD,
		"i32.rand":     OP_I32_RAND,
		"i32.bitand":   OP_I32_BITAND,
		"i32.bitor":    OP_I32_BITOR,
		"i32.bitxor":   OP_I32_BITXOR,
		"i32.bitclear": OP_I32_BITCLEAR,
		"i32.bitshl":   OP_I32_BITSHL,
		"i32.bitshr":   OP_I32_BITSHR,
		"i32.sqrt":     OP_I32_SQRT,
		"i32.log":      OP_I32_LOG,
		"i32.log2":     OP_I32_LOG2,
		"i32.log10":    OP_I32_LOG10,
		"i32.max":      OP_I32_MAX,
		"i32.min":      OP_I32_MIN,

		"i64.byte":     OP_I64_BYTE,
		"i64.str":      OP_I64_STR,
		"i64.i32":      OP_I64_I32,
		"i64.i64":      OP_I64_I64,
		"i64.f32":      OP_I64_F32,
		"i64.f64":      OP_I64_F64,
		
		"i64.print":    OP_I64_PRINT,
		"i64.add":      OP_I64_ADD,
		"i64.sub":      OP_I64_SUB,
		"i64.mul":      OP_I64_MUL,
		"i64.div":      OP_I64_DIV,
		"i64.abs":      OP_I64_ABS,
		"i64.pow":      OP_I64_POW,
		"i64.gt":       OP_I64_GT,
		"i64.gteq":     OP_I64_GTEQ,
		"i64.lt":       OP_I64_LT,
		"i64.lteq":     OP_I64_LTEQ,
		"i64.eq":       OP_I64_EQ,
		"i64.uneq":     OP_I64_UNEQ,
		"i64.mod":      OP_I64_MOD,
		"i64.rand":     OP_I64_RAND,
		"i64.bitand":   OP_I64_BITAND,
		"i64.bitor":    OP_I64_BITOR,
		"i64.bitxor":   OP_I64_BITXOR,
		"i64.bitclear": OP_I64_BITCLEAR,
		"i64.bitshl":   OP_I64_BITSHL,
		"i64.bitshr":   OP_I64_BITSHR,
		"i64.sqrt":     OP_I64_SQRT,
		"i64.log":      OP_I64_LOG,
		"i64.log2":     OP_I64_LOG2,
		"i64.log10":    OP_I64_LOG10,
		"i64.max":      OP_I64_MAX,
		"i64.min":      OP_I64_MIN,
		
		"f32.byte":     OP_F32_BYTE,
		"f32.str":      OP_F32_STR,
		"f32.i32":      OP_F32_I32,
		"f32.i64":      OP_F32_I64,
		"f32.f32":      OP_F32_F32,
		"f32.f64":      OP_F32_F64,
		
		"f32.print":    OP_F32_PRINT,
		"f32.add":      OP_F32_ADD,
		"f32.sub":      OP_F32_SUB,
		"f32.mul":      OP_F32_MUL,
		"f32.div":      OP_F32_DIV,
		"f32.abs":      OP_F32_ABS,
		"f32.pow":      OP_F32_POW,
		"f32.gt":       OP_F32_GT,
		"f32.gteq":     OP_F32_GTEQ,
		"f32.lt":       OP_F32_LT,
		"f32.lteq":     OP_F32_LTEQ,
		"f32.eq":       OP_F32_EQ,
		"f32.uneq":     OP_F32_UNEQ,
		"f32.cos":      OP_F32_COS,
		"f32.sin":      OP_F32_SIN,
		"f32.sqrt":     OP_F32_SQRT,
		"f32.log":      OP_F32_LOG,
		"f32.log2":     OP_F32_LOG2,
		"f32.log10":    OP_F32_LOG10,
		"f32.max":      OP_F32_MAX,
		"f32.min":      OP_F32_MIN,

		"f64.byte":     OP_F64_BYTE,
		"f64.str":      OP_F64_STR,
		"f64.i32":      OP_F64_I32,
		"f64.i64":      OP_F64_I64,
		"f64.f32":      OP_F64_F32,
		"f64.f64":      OP_F64_F64,

		"f64.print":    OP_F64_PRINT,
		"f64.add":      OP_F64_ADD,
		"f64.sub":      OP_F64_SUB,
		"f64.mul":      OP_F64_MUL,
		"f64.div":      OP_F64_DIV,
		"f64.abs":      OP_F64_ABS,
		"f64.pow":      OP_F64_POW,
		"f64.gt":       OP_F64_GT,
		"f64.gteq":     OP_F64_GTEQ,
		"f64.lt":       OP_F64_LT,
		"f64.lteq":     OP_F64_LTEQ,
		"f64.eq":       OP_F64_EQ,
		"f64.uneq":     OP_F64_UNEQ,
		"f64.cos":      OP_F64_COS,
		"f64.sin":      OP_F64_SIN,
		"f64.sqrt":     OP_F64_SQRT,
		"f64.log":      OP_F64_LOG,
		"f64.log2":     OP_F64_LOG2,
		"f64.log10":    OP_F64_LOG10,
		"f64.max":      OP_F64_MAX,
		"f64.min":      OP_F64_MIN,

		"str.print": OP_STR_PRINT,
		"str.concat": OP_STR_CONCAT,
		"str.eq": OP_STR_EQ,

		"str.byte": OP_STR_BYTE,
		"str.str": OP_STR_STR,
		"str.i32": OP_STR_I32,
		"str.i64": OP_STR_I64,
		"str.f32": OP_STR_F32,
		"str.f64": OP_STR_F64,

		"append":       OP_APPEND,
		"assert":       OP_ASSERT,

		// affordances
		"aff.print": OP_AFF_PRINT,
		"aff.query": OP_AFF_QUERY,
		"aff.on": OP_AFF_ON,
		"aff.of": OP_AFF_OF,
		"aff.inform": OP_AFF_INFORM,
		"aff.request": OP_AFF_REQUEST,
	}

	Natives = map[int]*CXFunction{
		OP_IDENTITY:     MakeNative(OP_IDENTITY, []int{TYPE_UNDEFINED}, []int{TYPE_UNDEFINED}),
		OP_JMP:          MakeNative(OP_JMP, []int{TYPE_BOOL}, []int{}),
		OP_DEBUG:        MakeNative(OP_DEBUG, []int{}, []int{}),

		OP_UND_EQUAL:    MakeNative(OP_UND_EQUAL, []int{TYPE_UNDEFINED, TYPE_UNDEFINED}, []int{TYPE_BOOL}),
		OP_UND_UNEQUAL:  MakeNative(OP_UND_UNEQUAL, []int{TYPE_UNDEFINED, TYPE_UNDEFINED}, []int{TYPE_BOOL}),
		OP_UND_BITAND:   MakeNative(OP_UND_BITAND, []int{TYPE_UNDEFINED, TYPE_UNDEFINED}, []int{TYPE_UNDEFINED}),
		OP_UND_BITXOR:   MakeNative(OP_UND_BITXOR, []int{TYPE_UNDEFINED, TYPE_UNDEFINED}, []int{TYPE_UNDEFINED}),
		OP_UND_BITOR:    MakeNative(OP_UND_BITOR, []int{TYPE_UNDEFINED, TYPE_UNDEFINED}, []int{TYPE_UNDEFINED}),
		OP_UND_BITCLEAR: MakeNative(OP_UND_BITCLEAR, []int{TYPE_UNDEFINED, TYPE_UNDEFINED}, []int{TYPE_UNDEFINED}),
		OP_UND_MUL:      MakeNative(OP_UND_MUL, []int{TYPE_UNDEFINED, TYPE_UNDEFINED}, []int{TYPE_UNDEFINED}),
		OP_UND_DIV:      MakeNative(OP_UND_DIV, []int{TYPE_UNDEFINED, TYPE_UNDEFINED}, []int{TYPE_UNDEFINED}),
		OP_UND_MOD:      MakeNative(OP_UND_MOD, []int{TYPE_UNDEFINED, TYPE_UNDEFINED}, []int{TYPE_UNDEFINED}),
		OP_UND_ADD:      MakeNative(OP_UND_ADD, []int{TYPE_UNDEFINED, TYPE_UNDEFINED}, []int{TYPE_UNDEFINED}),
		OP_UND_SUB:      MakeNative(OP_UND_SUB, []int{TYPE_UNDEFINED, TYPE_UNDEFINED}, []int{TYPE_UNDEFINED}),
		OP_UND_BITSHL:   MakeNative(OP_UND_BITSHL, []int{TYPE_UNDEFINED, TYPE_UNDEFINED}, []int{TYPE_UNDEFINED}),
		OP_UND_BITSHR:   MakeNative(OP_UND_BITSHR, []int{TYPE_UNDEFINED, TYPE_UNDEFINED}, []int{TYPE_UNDEFINED}),
		OP_UND_LT:       MakeNative(OP_UND_LT, []int{TYPE_UNDEFINED, TYPE_UNDEFINED}, []int{TYPE_BOOL}),
		OP_UND_GT:       MakeNative(OP_UND_GT, []int{TYPE_UNDEFINED, TYPE_UNDEFINED}, []int{TYPE_BOOL}),
		OP_UND_LTEQ:     MakeNative(OP_UND_LTEQ, []int{TYPE_UNDEFINED, TYPE_UNDEFINED}, []int{TYPE_BOOL}),
		OP_UND_GTEQ:     MakeNative(OP_UND_GTEQ, []int{TYPE_UNDEFINED, TYPE_UNDEFINED}, []int{TYPE_BOOL}),
		OP_UND_LEN:      MakeNative(OP_UND_LEN, []int{TYPE_UNDEFINED}, []int{TYPE_I32}),
		OP_UND_PRINTF:   MakeNative(OP_UND_PRINTF, []int{TYPE_UNDEFINED}, []int{}),
		OP_UND_SPRINTF:  MakeNative(OP_UND_SPRINTF, []int{TYPE_UNDEFINED}, []int{TYPE_STR}),
		OP_UND_READ:     MakeNative(OP_UND_READ, []int{}, []int{TYPE_STR}),

		OP_BYTE_BYTE:   MakeNative(OP_BYTE_BYTE, []int{TYPE_BYTE}, []int{TYPE_BYTE}),
		OP_BYTE_STR:    MakeNative(OP_BYTE_STR, []int{TYPE_BYTE}, []int{TYPE_STR}),
		OP_BYTE_I32:    MakeNative(OP_BYTE_I32, []int{TYPE_BYTE}, []int{TYPE_I32}),
		OP_BYTE_I64:    MakeNative(OP_BYTE_I64, []int{TYPE_BYTE}, []int{TYPE_I64}),
		OP_BYTE_F32:    MakeNative(OP_BYTE_F32, []int{TYPE_BYTE}, []int{TYPE_F32}),
		OP_BYTE_F64:    MakeNative(OP_BYTE_F64, []int{TYPE_BYTE}, []int{TYPE_F64}),

		OP_BYTE_PRINT:   MakeNative(OP_BYTE_PRINT, []int{TYPE_BYTE}, []int{}),

		OP_BOOL_PRINT:   MakeNative(OP_BOOL_PRINT, []int{TYPE_BOOL}, []int{}),
		OP_BOOL_EQUAL:   MakeNative(OP_BOOL_EQUAL, []int{TYPE_BOOL, TYPE_BOOL}, []int{TYPE_BOOL}),
		OP_BOOL_UNEQUAL: MakeNative(OP_BOOL_UNEQUAL, []int{TYPE_BOOL, TYPE_BOOL}, []int{TYPE_BOOL}),
		OP_BOOL_NOT:     MakeNative(OP_BOOL_NOT, []int{TYPE_BOOL}, []int{TYPE_BOOL}),
		OP_BOOL_OR:      MakeNative(OP_BOOL_OR, []int{TYPE_BOOL, TYPE_BOOL}, []int{TYPE_BOOL}),
		OP_BOOL_AND:     MakeNative(OP_BOOL_AND, []int{TYPE_BOOL, TYPE_BOOL}, []int{TYPE_BOOL}),

		OP_I32_BYTE:     MakeNative(OP_I32_BYTE, []int{TYPE_I32}, []int{TYPE_BYTE}),
		OP_I32_STR:      MakeNative(OP_I32_STR, []int{TYPE_I32}, []int{TYPE_STR}),
		OP_I32_I32:      MakeNative(OP_I32_I32, []int{TYPE_I32}, []int{TYPE_I32}),
		OP_I32_I64:      MakeNative(OP_I32_I64, []int{TYPE_I32}, []int{TYPE_I64}),
		OP_I32_F32:      MakeNative(OP_I32_F32, []int{TYPE_I32}, []int{TYPE_F32}),
		OP_I32_F64:      MakeNative(OP_I32_F64, []int{TYPE_I32}, []int{TYPE_F64}),

		OP_I32_PRINT:    MakeNative(OP_I32_PRINT, []int{TYPE_I32}, []int{}),
		OP_I32_ADD:      MakeNative(OP_I32_ADD, []int{TYPE_I32, TYPE_I32}, []int{TYPE_I32}),
		OP_I32_SUB:      MakeNative(OP_I32_SUB, []int{TYPE_I32, TYPE_I32}, []int{TYPE_I32}),
		OP_I32_MUL:      MakeNative(OP_I32_MUL, []int{TYPE_I32, TYPE_I32}, []int{TYPE_I32}),
		OP_I32_DIV:      MakeNative(OP_I32_DIV, []int{TYPE_I32, TYPE_I32}, []int{TYPE_I32}),
		OP_I32_ABS:      MakeNative(OP_I32_ABS, []int{TYPE_I32}, []int{TYPE_I32}),
		OP_I32_POW:      MakeNative(OP_I32_POW, []int{TYPE_I32, TYPE_I32}, []int{TYPE_I32}),
		OP_I32_GT:       MakeNative(OP_I32_GT, []int{TYPE_I32, TYPE_I32}, []int{TYPE_BOOL}),
		OP_I32_GTEQ:     MakeNative(OP_I32_GTEQ, []int{TYPE_I32, TYPE_I32}, []int{TYPE_BOOL}),
		OP_I32_LT:       MakeNative(OP_I32_LT, []int{TYPE_I32, TYPE_I32}, []int{TYPE_BOOL}),
		OP_I32_LTEQ:     MakeNative(OP_I32_LTEQ, []int{TYPE_I32, TYPE_I32}, []int{TYPE_BOOL}),
		OP_I32_EQ:       MakeNative(OP_I32_EQ, []int{TYPE_I32, TYPE_I32}, []int{TYPE_BOOL}),
		OP_I32_UNEQ:     MakeNative(OP_I32_UNEQ, []int{TYPE_I32, TYPE_I32}, []int{TYPE_BOOL}),
		OP_I32_MOD:      MakeNative(OP_I32_MOD, []int{TYPE_I32, TYPE_I32}, []int{TYPE_I32}),
		OP_I32_RAND:     MakeNative(OP_I32_RAND, []int{TYPE_I32, TYPE_I32}, []int{TYPE_I32}),
		OP_I32_BITAND:   MakeNative(OP_I32_BITAND, []int{TYPE_I32, TYPE_I32}, []int{TYPE_I32}),
		OP_I32_BITOR:    MakeNative(OP_I32_BITOR, []int{TYPE_I32, TYPE_I32}, []int{TYPE_I32}),
		OP_I32_BITXOR:   MakeNative(OP_I32_BITXOR, []int{TYPE_I32, TYPE_I32}, []int{TYPE_I32}),
		OP_I32_BITCLEAR: MakeNative(OP_I32_BITCLEAR, []int{TYPE_I32, TYPE_I32}, []int{TYPE_I32}),
		OP_I32_BITSHL:   MakeNative(OP_I32_BITSHL, []int{TYPE_I32, TYPE_I32}, []int{TYPE_I32}),
		OP_I32_BITSHR:   MakeNative(OP_I32_BITSHR, []int{TYPE_I32, TYPE_I32}, []int{TYPE_I32}),
		OP_I32_SQRT:     MakeNative(OP_I32_SQRT, []int{TYPE_I32}, []int{TYPE_I32}),
		OP_I32_LOG:      MakeNative(OP_I32_LOG, []int{TYPE_I32}, []int{TYPE_I32}),
		OP_I32_LOG2:     MakeNative(OP_I32_LOG2, []int{TYPE_I32}, []int{TYPE_I32}),
		OP_I32_LOG10:    MakeNative(OP_I32_LOG10, []int{TYPE_I32}, []int{TYPE_I32}),
		OP_I32_MAX:      MakeNative(OP_I32_MAX, []int{TYPE_I32, TYPE_I32}, []int{TYPE_I32}),
		OP_I32_MIN:      MakeNative(OP_I32_MIN, []int{TYPE_I32, TYPE_I32}, []int{TYPE_I32}),

		OP_I64_BYTE:     MakeNative(OP_I64_BYTE, []int{TYPE_I64}, []int{TYPE_BYTE}),
		OP_I64_STR:      MakeNative(OP_I64_STR, []int{TYPE_I64}, []int{TYPE_STR}),
		OP_I64_I32:      MakeNative(OP_I64_I32, []int{TYPE_I64}, []int{TYPE_I32}),
		OP_I64_I64:      MakeNative(OP_I64_I64, []int{TYPE_I64}, []int{TYPE_I64}),
		OP_I64_F32:      MakeNative(OP_I64_F32, []int{TYPE_I64}, []int{TYPE_F32}),
		OP_I64_F64:      MakeNative(OP_I64_F64, []int{TYPE_I64}, []int{TYPE_F64}),
		
		OP_I64_PRINT:    MakeNative(OP_I64_PRINT, []int{TYPE_I64}, []int{}),
		OP_I64_ADD:      MakeNative(OP_I64_ADD, []int{TYPE_I64, TYPE_I64}, []int{TYPE_I64}),
		OP_I64_SUB:      MakeNative(OP_I64_SUB, []int{TYPE_I64, TYPE_I64}, []int{TYPE_I64}),
		OP_I64_MUL:      MakeNative(OP_I64_MUL, []int{TYPE_I64, TYPE_I64}, []int{TYPE_I64}),
		OP_I64_DIV:      MakeNative(OP_I64_DIV, []int{TYPE_I64, TYPE_I64}, []int{TYPE_I64}),
		OP_I64_ABS:      MakeNative(OP_I64_ABS, []int{TYPE_I64}, []int{TYPE_I64}),
		OP_I64_POW:      MakeNative(OP_I64_POW, []int{TYPE_I64, TYPE_I64}, []int{TYPE_I64}),
		OP_I64_GT:       MakeNative(OP_I64_GT, []int{TYPE_I64, TYPE_I64}, []int{TYPE_BOOL}),
		OP_I64_GTEQ:     MakeNative(OP_I64_GTEQ, []int{TYPE_I64, TYPE_I64}, []int{TYPE_BOOL}),
		OP_I64_LT:       MakeNative(OP_I64_LT, []int{TYPE_I64, TYPE_I64}, []int{TYPE_BOOL}),
		OP_I64_LTEQ:     MakeNative(OP_I64_LTEQ, []int{TYPE_I64, TYPE_I64}, []int{TYPE_BOOL}),
		OP_I64_EQ:       MakeNative(OP_I64_EQ, []int{TYPE_I64, TYPE_I64}, []int{TYPE_BOOL}),
		OP_I64_UNEQ:     MakeNative(OP_I64_UNEQ, []int{TYPE_I64, TYPE_I64}, []int{TYPE_BOOL}),
		OP_I64_MOD:      MakeNative(OP_I64_MOD, []int{TYPE_I64, TYPE_I64}, []int{TYPE_I64}),
		OP_I64_RAND:     MakeNative(OP_I64_RAND, []int{TYPE_I64, TYPE_I64}, []int{TYPE_I64}),
		OP_I64_BITAND:   MakeNative(OP_I64_BITAND, []int{TYPE_I64, TYPE_I64}, []int{TYPE_I64}),
		OP_I64_BITOR:    MakeNative(OP_I64_BITOR, []int{TYPE_I64, TYPE_I64}, []int{TYPE_I64}),
		OP_I64_BITXOR:   MakeNative(OP_I64_BITXOR, []int{TYPE_I64, TYPE_I64}, []int{TYPE_I64}),
		OP_I64_BITCLEAR: MakeNative(OP_I64_BITCLEAR, []int{TYPE_I64, TYPE_I64}, []int{TYPE_I64}),
		OP_I64_BITSHL:   MakeNative(OP_I64_BITSHL, []int{TYPE_I64, TYPE_I64}, []int{TYPE_I64}),
		OP_I64_BITSHR:   MakeNative(OP_I64_BITSHR, []int{TYPE_I64, TYPE_I64}, []int{TYPE_I64}),
		OP_I64_SQRT:     MakeNative(OP_I64_SQRT, []int{TYPE_I64}, []int{TYPE_I64}),
		OP_I64_LOG:      MakeNative(OP_I64_LOG, []int{TYPE_I64}, []int{TYPE_I64}),
		OP_I64_LOG2:     MakeNative(OP_I64_LOG2, []int{TYPE_I64}, []int{TYPE_I64}),
		OP_I64_LOG10:    MakeNative(OP_I64_LOG10, []int{TYPE_I64}, []int{TYPE_I64}),
		OP_I64_MAX:      MakeNative(OP_I64_MAX, []int{TYPE_I64, TYPE_I64}, []int{TYPE_I64}),
		OP_I64_MIN:      MakeNative(OP_I64_MIN, []int{TYPE_I64, TYPE_I64}, []int{TYPE_I64}),

		OP_F32_BYTE:     MakeNative(OP_F32_BYTE, []int{TYPE_F32}, []int{TYPE_BYTE}),
		OP_F32_STR:      MakeNative(OP_F32_STR,  []int{TYPE_F32}, []int{TYPE_STR}),
		OP_F32_I32:      MakeNative(OP_F32_I32,  []int{TYPE_F32}, []int{TYPE_I32}),
		OP_F32_I64:      MakeNative(OP_F32_I64,  []int{TYPE_F32}, []int{TYPE_I64}),
		OP_F32_F32:      MakeNative(OP_F32_F32,  []int{TYPE_F32}, []int{TYPE_F32}),
		OP_F32_F64:      MakeNative(OP_F32_F64,  []int{TYPE_F32}, []int{TYPE_F64}),
		
		OP_F32_PRINT:    MakeNative(OP_F32_PRINT, []int{TYPE_F32}, []int{}),
		OP_F32_ADD:      MakeNative(OP_F32_ADD, []int{TYPE_F32, TYPE_F32}, []int{TYPE_F32}),
		OP_F32_SUB:      MakeNative(OP_F32_SUB, []int{TYPE_F32, TYPE_F32}, []int{TYPE_F32}),
		OP_F32_MUL:      MakeNative(OP_F32_MUL, []int{TYPE_F32, TYPE_F32}, []int{TYPE_F32}),
		OP_F32_DIV:      MakeNative(OP_F32_DIV, []int{TYPE_F32, TYPE_F32}, []int{TYPE_F32}),
		OP_F32_ABS:      MakeNative(OP_F32_ABS, []int{TYPE_F32}, []int{TYPE_F32}),
		OP_F32_POW:      MakeNative(OP_F32_POW, []int{TYPE_F32, TYPE_F32}, []int{TYPE_F32}),
		OP_F32_GT:       MakeNative(OP_F32_GT, []int{TYPE_F32, TYPE_F32}, []int{TYPE_BOOL}),
		OP_F32_GTEQ:     MakeNative(OP_F32_GTEQ, []int{TYPE_F32, TYPE_F32}, []int{TYPE_BOOL}),
		OP_F32_LT:       MakeNative(OP_F32_LT, []int{TYPE_F32, TYPE_F32}, []int{TYPE_BOOL}),
		OP_F32_LTEQ:     MakeNative(OP_F32_LTEQ, []int{TYPE_F32, TYPE_F32}, []int{TYPE_BOOL}),
		OP_F32_EQ:       MakeNative(OP_F32_EQ, []int{TYPE_F32, TYPE_F32}, []int{TYPE_BOOL}),
		OP_F32_UNEQ:     MakeNative(OP_F32_UNEQ, []int{TYPE_F32, TYPE_F32}, []int{TYPE_BOOL}),
		OP_F32_COS:      MakeNative(OP_F32_COS, []int{TYPE_F32}, []int{TYPE_F32}),
		OP_F32_SIN:      MakeNative(OP_F32_SIN, []int{TYPE_F32}, []int{TYPE_F32}),
		OP_F32_SQRT:     MakeNative(OP_F32_SQRT, []int{TYPE_F32}, []int{TYPE_F32}),
		OP_F32_LOG:      MakeNative(OP_F32_LOG, []int{TYPE_F32}, []int{TYPE_F32}),
		OP_F32_LOG2:     MakeNative(OP_F32_LOG2, []int{TYPE_F32}, []int{TYPE_F32}),
		OP_F32_LOG10:    MakeNative(OP_F32_LOG10, []int{TYPE_F32}, []int{TYPE_F32}),
		OP_F32_MAX:      MakeNative(OP_F32_MAX, []int{TYPE_F32, TYPE_F32}, []int{TYPE_F32}),
		OP_F32_MIN:      MakeNative(OP_F32_MIN, []int{TYPE_F32, TYPE_F32}, []int{TYPE_F32}),

		OP_F64_BYTE:   MakeNative(OP_F64_BYTE, []int{TYPE_F64}, []int{TYPE_BYTE}),
		OP_F64_STR:    MakeNative(OP_F64_STR, []int{TYPE_F64}, []int{TYPE_STR}),
		OP_F64_I32:    MakeNative(OP_F64_I32, []int{TYPE_F64}, []int{TYPE_I32}),
		OP_F64_I64:    MakeNative(OP_F64_I64, []int{TYPE_F64}, []int{TYPE_I64}),
		OP_F64_F32:    MakeNative(OP_F64_F32, []int{TYPE_F64}, []int{TYPE_F32}),
		OP_F64_F64:    MakeNative(OP_F64_F64, []int{TYPE_F64}, []int{TYPE_F64}),

		OP_F64_PRINT:    MakeNative(OP_F64_PRINT, []int{TYPE_F64}, []int{}),
		OP_F64_ADD:      MakeNative(OP_F64_ADD, []int{TYPE_F64, TYPE_F64}, []int{TYPE_F64}),
		OP_F64_SUB:      MakeNative(OP_F64_SUB, []int{TYPE_F64, TYPE_F64}, []int{TYPE_F64}),
		OP_F64_MUL:      MakeNative(OP_F64_MUL, []int{TYPE_F64, TYPE_F64}, []int{TYPE_F64}),
		OP_F64_DIV:      MakeNative(OP_F64_DIV, []int{TYPE_F64, TYPE_F64}, []int{TYPE_F64}),
		OP_F64_ABS:      MakeNative(OP_F64_ABS, []int{TYPE_F64}, []int{TYPE_F64}),
		OP_F64_POW:      MakeNative(OP_F64_POW, []int{TYPE_F64, TYPE_F64}, []int{TYPE_F64}),
		OP_F64_GT:       MakeNative(OP_F64_GT, []int{TYPE_F64, TYPE_F64}, []int{TYPE_BOOL}),
		OP_F64_GTEQ:     MakeNative(OP_F64_GTEQ, []int{TYPE_F64, TYPE_F64}, []int{TYPE_BOOL}),
		OP_F64_LT:       MakeNative(OP_F64_LT, []int{TYPE_F64, TYPE_F64}, []int{TYPE_BOOL}),
		OP_F64_LTEQ:     MakeNative(OP_F64_LTEQ, []int{TYPE_F64, TYPE_F64}, []int{TYPE_BOOL}),
		OP_F64_EQ:       MakeNative(OP_F64_EQ, []int{TYPE_F64, TYPE_F64}, []int{TYPE_BOOL}),
		OP_F64_UNEQ:     MakeNative(OP_F64_UNEQ, []int{TYPE_F64, TYPE_F64}, []int{TYPE_BOOL}),
		OP_F64_COS:      MakeNative(OP_F64_COS, []int{TYPE_F64}, []int{TYPE_F64}),
		OP_F64_SIN:      MakeNative(OP_F64_SIN, []int{TYPE_F64}, []int{TYPE_F64}),
		OP_F64_SQRT:     MakeNative(OP_F64_SQRT, []int{TYPE_F64}, []int{TYPE_F64}),
		OP_F64_LOG:      MakeNative(OP_F64_LOG, []int{TYPE_F64}, []int{TYPE_F64}),
		OP_F64_LOG2:     MakeNative(OP_F64_LOG2, []int{TYPE_F64}, []int{TYPE_F64}),
		OP_F64_LOG10:    MakeNative(OP_F64_LOG10, []int{TYPE_F64}, []int{TYPE_F64}),
		OP_F64_MIN:      MakeNative(OP_F64_MIN, []int{TYPE_F64, TYPE_F64}, []int{TYPE_F64}),
		OP_F64_MAX:      MakeNative(OP_F64_MAX, []int{TYPE_F64, TYPE_F64}, []int{TYPE_F64}),
		
		OP_STR_PRINT:    MakeNative(OP_STR_PRINT, []int{TYPE_STR}, []int{}),
		OP_STR_CONCAT:   MakeNative(OP_STR_CONCAT, []int{TYPE_STR, TYPE_STR}, []int{TYPE_STR}),
		OP_STR_EQ:       MakeNative(OP_STR_EQ, []int{TYPE_STR, TYPE_STR}, []int{TYPE_BOOL}),

		OP_STR_BYTE:      MakeNative(OP_STR_BYTE,[]int{TYPE_STR}, []int{TYPE_BYTE}),
		OP_STR_STR:       MakeNative(OP_STR_STR, []int{TYPE_STR}, []int{TYPE_STR}),
		OP_STR_I32:       MakeNative(OP_STR_I32, []int{TYPE_STR}, []int{TYPE_I32}),
		OP_STR_I64:       MakeNative(OP_STR_I64, []int{TYPE_STR}, []int{TYPE_I64}),
		OP_STR_F32:       MakeNative(OP_STR_F32, []int{TYPE_STR}, []int{TYPE_F32}),
		OP_STR_F64:       MakeNative(OP_STR_F64, []int{TYPE_STR}, []int{TYPE_F64}),

		OP_APPEND:     MakeNative(OP_APPEND, []int{TYPE_UNDEFINED, TYPE_UNDEFINED}, []int{TYPE_UNDEFINED}),
		OP_ASSERT:     MakeNative(OP_ASSERT, []int{TYPE_UNDEFINED, TYPE_UNDEFINED, TYPE_STR}, []int{TYPE_BOOL}),

		// affordances
		OP_AFF_PRINT:  MakeNative(OP_AFF_PRINT, []int{TYPE_AFF}, []int{}),
		OP_AFF_QUERY:  MakeNative(OP_AFF_QUERY, []int{TYPE_AFF}, []int{TYPE_AFF}),
		OP_AFF_ON: MakeNative(OP_AFF_ON, []int{TYPE_AFF, TYPE_AFF}, []int{}),
		OP_AFF_OF: MakeNative(OP_AFF_OF, []int{TYPE_AFF, TYPE_AFF}, []int{}),
		OP_AFF_INFORM: MakeNative(OP_AFF_INFORM, []int{TYPE_AFF, TYPE_I32, TYPE_AFF}, []int{}),
		OP_AFF_REQUEST: MakeNative(OP_AFF_REQUEST, []int{TYPE_AFF, TYPE_I32, TYPE_AFF}, []int{}),
	}

	execNative = func(prgrm *CXProgram) {
		call := &prgrm.CallStack[prgrm.CallCounter]
		expr := call.Operator.Expressions[call.Line]
		opCode := expr.Operator.OpCode
		fp := call.FramePointer

		switch opCode {
		case OP_IDENTITY:
			op_identity(expr, fp)
		case OP_JMP:
			op_jmp(expr, fp, call)
		case OP_DEBUG:
			prgrm.PrintStack()

		case OP_UND_EQUAL:
			op_equal(expr, fp)
		case OP_UND_UNEQUAL:
			op_unequal(expr, fp)
		case OP_UND_BITAND:
			op_bitand(expr, fp)
		case OP_UND_BITXOR:
			op_bitxor(expr, fp)
		case OP_UND_BITOR:
			op_bitor(expr, fp)
		case OP_UND_BITCLEAR:
			op_bitclear(expr, fp)
		case OP_UND_MUL:
			op_mul(expr, fp)
		case OP_UND_DIV:
			op_div(expr, fp)
		case OP_UND_MOD:
			op_mod(expr, fp)
		case OP_UND_ADD:
			op_add(expr, fp)
		case OP_UND_SUB:
			op_sub(expr, fp)
		case OP_UND_BITSHL:
			op_bitshl(expr, fp)
		case OP_UND_BITSHR:
			op_bitshr(expr, fp)
		case OP_UND_LT:
			op_lt(expr, fp)
		case OP_UND_GT:
			op_gt(expr, fp)
		case OP_UND_LTEQ:
			op_lteq(expr, fp)
		case OP_UND_GTEQ:
			op_gteq(expr, fp)
		case OP_UND_LEN:
			op_len(expr, fp)
		case OP_UND_PRINTF:
			op_printf(expr, fp)
		case OP_UND_SPRINTF:
			op_sprintf(expr, fp)
		case OP_UND_READ:
			op_read(expr, fp)

		case OP_BYTE_BYTE:
			op_byte_byte(expr, fp)
		case OP_BYTE_STR:
			op_byte_byte(expr, fp)
		case OP_BYTE_I32:
			op_byte_byte(expr, fp)
		case OP_BYTE_I64:
			op_byte_byte(expr, fp)
		case OP_BYTE_F32:
			op_byte_byte(expr, fp)
		case OP_BYTE_F64:
			op_byte_byte(expr, fp)

		case OP_BYTE_PRINT:
			op_byte_print(expr, fp)

		case OP_BOOL_PRINT:
			op_bool_print(expr, fp)
		case OP_BOOL_EQUAL:
			op_bool_equal(expr, fp)
		case OP_BOOL_UNEQUAL:
			op_bool_unequal(expr, fp)
		case OP_BOOL_NOT:
			op_bool_not(expr, fp)
		case OP_BOOL_OR:
			op_bool_or(expr, fp)
		case OP_BOOL_AND:
			op_bool_and(expr, fp)

		case OP_I32_BYTE:
			op_i32_i32(expr, fp)
		case OP_I32_STR:
			op_i32_i32(expr, fp)
		case OP_I32_I32:
			op_i32_i32(expr, fp)
		case OP_I32_I64:
			op_i32_i32(expr, fp)
		case OP_I32_F32:
			op_i32_i32(expr, fp)
		case OP_I32_F64:
			op_i32_i32(expr, fp)
			
		case OP_I32_PRINT:
			op_i32_print(expr, fp)
		case OP_I32_ADD:
			op_i32_add(expr, fp)
		case OP_I32_SUB:
			op_i32_sub(expr, fp)
		case OP_I32_MUL:
			op_i32_mul(expr, fp)
		case OP_I32_DIV:
			op_i32_div(expr, fp)
		case OP_I32_ABS:
			op_i32_abs(expr, fp)
		case OP_I32_POW:
			op_i32_pow(expr, fp)
		case OP_I32_GT:
			op_i32_gt(expr, fp)
		case OP_I32_GTEQ:
			op_i32_gteq(expr, fp)
		case OP_I32_LT:
			op_i32_lt(expr, fp)
		case OP_I32_LTEQ:
			op_i32_lteq(expr, fp)
		case OP_I32_EQ:
			op_i32_eq(expr, fp)
		case OP_I32_UNEQ:
			op_i32_uneq(expr, fp)
		case OP_I32_MOD:
			op_i32_mod(expr, fp)
		case OP_I32_RAND:
			op_i32_rand(expr, fp)
		case OP_I32_BITAND:
			op_i32_bitand(expr, fp)
		case OP_I32_BITOR:
			op_i32_bitor(expr, fp)
		case OP_I32_BITXOR:
			op_i32_bitxor(expr, fp)
		case OP_I32_BITCLEAR:
			op_i32_bitclear(expr, fp)
		case OP_I32_BITSHL:
			op_i32_bitshl(expr, fp)
		case OP_I32_BITSHR:
			op_i32_bitshr(expr, fp)
		case OP_I32_SQRT:
			op_i32_sqrt(expr, fp)
		case OP_I32_LOG:
			op_i32_log(expr, fp)
		case OP_I32_LOG2:
			op_i32_log2(expr, fp)
		case OP_I32_LOG10:
			op_i32_log10(expr, fp)

		case OP_I32_MAX:
			op_i32_max(expr, fp)
		case OP_I32_MIN:
			op_i32_min(expr, fp)

		case OP_I64_BYTE:
			op_i64_i64(expr, fp)
		case OP_I64_STR:
			op_i64_i64(expr, fp)
		case OP_I64_I32:
			op_i64_i64(expr, fp)
		case OP_I64_I64:
			op_i64_i64(expr, fp)
		case OP_I64_F32:
			op_i64_i64(expr, fp)
		case OP_I64_F64:
			op_i64_i64(expr, fp)

		case OP_I64_PRINT:
			op_i64_print(expr, fp)
		case OP_I64_ADD:
			op_i64_add(expr, fp)
		case OP_I64_SUB:
			op_i64_sub(expr, fp)
		case OP_I64_MUL:
			op_i64_mul(expr, fp)
		case OP_I64_DIV:
			op_i64_div(expr, fp)
		case OP_I64_ABS:
			op_i64_abs(expr, fp)
		case OP_I64_POW:
			op_i64_pow(expr, fp)
		case OP_I64_GT:
			op_i64_gt(expr, fp)
		case OP_I64_GTEQ:
			op_i64_gteq(expr, fp)
		case OP_I64_LT:
			op_i64_lt(expr, fp)
		case OP_I64_LTEQ:
			op_i64_lteq(expr, fp)
		case OP_I64_EQ:
			op_i64_eq(expr, fp)
		case OP_I64_UNEQ:
			op_i64_uneq(expr, fp)
		case OP_I64_MOD:
			op_i64_mod(expr, fp)
		case OP_I64_RAND:
			op_i64_rand(expr, fp)
		case OP_I64_BITAND:
			op_i64_bitand(expr, fp)
		case OP_I64_BITOR:
			op_i64_bitor(expr, fp)
		case OP_I64_BITXOR:
			op_i64_bitxor(expr, fp)
		case OP_I64_BITCLEAR:
			op_i64_bitclear(expr, fp)
		case OP_I64_BITSHL:
			op_i64_bitshl(expr, fp)
		case OP_I64_BITSHR:
			op_i64_bitshr(expr, fp)
		case OP_I64_SQRT:
			op_i64_sqrt(expr, fp)
		case OP_I64_LOG:
			op_i64_log(expr, fp)
		case OP_I64_LOG2:
			op_i64_log2(expr, fp)
		case OP_I64_LOG10:
			op_i64_log10(expr, fp)
		case OP_I64_MAX:
			op_i64_max(expr, fp)
		case OP_I64_MIN:
			op_i64_min(expr, fp)

		case OP_F32_BYTE:
			op_f32_f32(expr, fp)
		case OP_F32_STR:
			op_f32_f32(expr, fp)
		case OP_F32_I32:
			op_f32_f32(expr, fp)
		case OP_F32_I64:
			op_f32_f32(expr, fp)
		case OP_F32_F32:
			op_f32_f32(expr, fp)
		case OP_F32_F64:
			op_f32_f32(expr, fp)
			
		case OP_F32_PRINT:
			op_f32_print(expr, fp)
		case OP_F32_ADD:
			op_f32_add(expr, fp)
		case OP_F32_SUB:
			op_f32_sub(expr, fp)
		case OP_F32_MUL:
			op_f32_mul(expr, fp)
		case OP_F32_DIV:
			op_f32_div(expr, fp)
		case OP_F32_ABS:
			op_f32_abs(expr, fp)
		case OP_F32_POW:
			op_f32_pow(expr, fp)
		case OP_F32_GT:
			op_f32_gt(expr, fp)
		case OP_F32_GTEQ:
			op_f32_gteq(expr, fp)
		case OP_F32_LT:
			op_f32_lt(expr, fp)
		case OP_F32_LTEQ:
			op_f32_lteq(expr, fp)
		case OP_F32_EQ:
			op_f32_eq(expr, fp)
		case OP_F32_UNEQ:
			op_f32_uneq(expr, fp)
		case OP_F32_COS:
			op_f32_cos(expr, fp)
		case OP_F32_SIN:
			op_f32_sin(expr, fp)
		case OP_F32_SQRT:
			op_f32_sqrt(expr, fp)
		case OP_F32_LOG:
			op_f32_log(expr, fp)
		case OP_F32_LOG2:
			op_f32_log2(expr, fp)
		case OP_F32_LOG10:
			op_f32_log10(expr, fp)
		case OP_F32_MAX:
			op_f32_max(expr, fp)
		case OP_F32_MIN:
			op_f32_min(expr, fp)

		case OP_F64_BYTE:
			op_f64_f64(expr, fp)
		case OP_F64_STR:
			op_f64_f64(expr, fp)
		case OP_F64_I32:
			op_f64_f64(expr, fp)
		case OP_F64_I64:
			op_f64_f64(expr, fp)
		case OP_F64_F32:
			op_f64_f64(expr, fp)
		case OP_F64_F64:
			op_f64_f64(expr, fp)

		case OP_F64_PRINT:
			op_f64_print(expr, fp)
		case OP_F64_ADD:
			op_f64_add(expr, fp)
		case OP_F64_SUB:
			op_f64_sub(expr, fp)
		case OP_F64_MUL:
			op_f64_mul(expr, fp)
		case OP_F64_DIV:
			op_f64_div(expr, fp)
		case OP_F64_ABS:
			op_f64_abs(expr, fp)
		case OP_F64_POW:
			op_f64_pow(expr, fp)
		case OP_F64_GT:
			op_f64_gt(expr, fp)
		case OP_F64_GTEQ:
			op_f64_gteq(expr, fp)
		case OP_F64_LT:
			op_f64_lt(expr, fp)
		case OP_F64_LTEQ:
			op_f64_lteq(expr, fp)
		case OP_F64_EQ:
			op_f64_eq(expr, fp)
		case OP_F64_UNEQ:
			op_f64_uneq(expr, fp)
		case OP_F64_COS:
			op_f64_cos(expr, fp)
		case OP_F64_SIN:
			op_f64_sin(expr, fp)
		case OP_F64_SQRT:
			op_f64_sqrt(expr, fp)
		case OP_F64_LOG:
			op_f64_log(expr, fp)
		case OP_F64_LOG2:
			op_f64_log2(expr, fp)
		case OP_F64_LOG10:
			op_f64_log10(expr, fp)
		case OP_F64_MAX:
			op_f64_max(expr, fp)
		case OP_F64_MIN:
			op_f64_min(expr, fp)
		case OP_STR_PRINT:
			op_str_print(expr, fp)
		case OP_STR_CONCAT:
			op_str_concat(expr, fp)
		case OP_STR_EQ:
			op_str_eq(expr, fp)
			
		case OP_STR_BYTE:
			op_str_str(expr, fp)
		case OP_STR_STR:
			op_str_str(expr, fp)
		case OP_STR_I32:
			op_str_str(expr, fp)
		case OP_STR_I64:
			op_str_str(expr, fp)
		case OP_STR_F32:
			op_str_str(expr, fp)
		case OP_STR_F64:
			op_str_str(expr, fp)

		case OP_MAKE:
		case OP_READ:
		case OP_WRITE:
		case OP_LEN:
		case OP_CONCAT:
		case OP_APPEND:
			op_append(expr, fp)
		case OP_COPY:
		case OP_CAST:
		case OP_EQ:
		case OP_UNEQ:
		case OP_AND:
		case OP_OR:
		case OP_NOT:
		case OP_SLEEP:
		case OP_HALT:
		case OP_GOTO:
		case OP_REMCX:
		case OP_ADDCX:
		case OP_QUERY:
		case OP_EXECUTE:
		case OP_INDEX:
		case OP_NAME:
		case OP_EVOLVE:
		case OP_ASSERT:
			op_assert_value(expr, fp)

			// affordances
		case OP_AFF_PRINT:
			op_aff_print(expr, fp)
		case OP_AFF_QUERY:
			op_aff_query(expr, fp)
		case OP_AFF_ON:
			op_aff_on(expr, fp)
		case OP_AFF_OF:
			op_aff_of(expr, fp)
		case OP_AFF_INFORM:
			op_aff_inform(expr, fp)
		case OP_AFF_REQUEST:
			op_aff_request(expr, fp)
		}
	}
}
