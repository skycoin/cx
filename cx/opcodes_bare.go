package base

import ()

var CorePackages = []string{
	// temporary solution until we can implement these packages in pure CX I guess
	"gl", "glfw", "time", "http", "os", "explorer", "aff", "gltext", "cx",
}

// op codes
const (
	OP_IDENTITY = iota
	OP_JMP
	OP_DEBUG

	OP_SERIALIZE
	OP_DESERIALIZE

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

	OP_F32_IS_NAN
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
	OP_STR_SUBSTR
	OP_STR_INDEX
	OP_STR_TRIM_SPACE
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
	OP_TEST
	OP_PANIC

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
var OpNames map[int]string = map[int]string{}
var OpCodes map[string]int = map[string]int{}
var Natives map[int]*CXFunction = map[int]*CXFunction{}
var execNativeBare func(*CXProgram)
var execNative func(*CXProgram)

func AddOpCode(code int, name string, inputs []int, outputs []int) {
	OpNames[code] = name
	OpCodes[name] = code
	Natives[code] = MakeNative(code, inputs, outputs)
}

/*
// debug helper
func DumpOpCodes(opCode int) () {
	var keys []int
	for k := range OpNames {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, k := range keys {
		fmt.Printf("%5d : %s\n", k, OpNames[k])
	}

	fmt.Printf("opCode : %d\n", opCode)
}*/

func init() {
	AddOpCode(OP_IDENTITY, "identity", []int{TYPE_UNDEFINED}, []int{TYPE_UNDEFINED})
	AddOpCode(OP_JMP, "jmp", []int{TYPE_BOOL}, []int{})
	AddOpCode(OP_DEBUG, "debug", []int{}, []int{})

	AddOpCode(OP_SERIALIZE, "serialize", []int{TYPE_AFF}, []int{TYPE_BYTE})
	AddOpCode(OP_DESERIALIZE, "deserialize", []int{TYPE_BYTE}, []int{})

	AddOpCode(OP_UND_EQUAL, "eq", []int{TYPE_UNDEFINED, TYPE_UNDEFINED}, []int{TYPE_BOOL})
	AddOpCode(OP_UND_UNEQUAL, "uneq", []int{TYPE_UNDEFINED, TYPE_UNDEFINED}, []int{TYPE_BOOL})
	AddOpCode(OP_UND_BITAND, "bitand", []int{TYPE_UNDEFINED, TYPE_UNDEFINED}, []int{TYPE_UNDEFINED})
	AddOpCode(OP_UND_BITXOR, "bitxor", []int{TYPE_UNDEFINED, TYPE_UNDEFINED}, []int{TYPE_UNDEFINED})
	AddOpCode(OP_UND_BITOR, "bitor", []int{TYPE_UNDEFINED, TYPE_UNDEFINED}, []int{TYPE_UNDEFINED})
	AddOpCode(OP_UND_BITCLEAR, "bitclear", []int{TYPE_UNDEFINED, TYPE_UNDEFINED}, []int{TYPE_UNDEFINED})
	AddOpCode(OP_UND_MUL, "mul", []int{TYPE_UNDEFINED, TYPE_UNDEFINED}, []int{TYPE_UNDEFINED})
	AddOpCode(OP_UND_DIV, "div", []int{TYPE_UNDEFINED, TYPE_UNDEFINED}, []int{TYPE_UNDEFINED})
	AddOpCode(OP_UND_MOD, "mod", []int{TYPE_UNDEFINED, TYPE_UNDEFINED}, []int{TYPE_UNDEFINED})
	AddOpCode(OP_UND_ADD, "add", []int{TYPE_UNDEFINED, TYPE_UNDEFINED}, []int{TYPE_UNDEFINED})
	AddOpCode(OP_UND_SUB, "sub", []int{TYPE_UNDEFINED, TYPE_UNDEFINED}, []int{TYPE_UNDEFINED})
	AddOpCode(OP_UND_BITSHL, "bitshl", []int{TYPE_UNDEFINED, TYPE_UNDEFINED}, []int{TYPE_UNDEFINED})
	AddOpCode(OP_UND_BITSHR, "bitshr", []int{TYPE_UNDEFINED, TYPE_UNDEFINED}, []int{TYPE_UNDEFINED})
	AddOpCode(OP_UND_LT, "lt", []int{TYPE_UNDEFINED, TYPE_UNDEFINED}, []int{TYPE_BOOL})
	AddOpCode(OP_UND_GT, "gt", []int{TYPE_UNDEFINED, TYPE_UNDEFINED}, []int{TYPE_BOOL})
	AddOpCode(OP_UND_LTEQ, "lteq", []int{TYPE_UNDEFINED, TYPE_UNDEFINED}, []int{TYPE_BOOL})
	AddOpCode(OP_UND_GTEQ, "gteq", []int{TYPE_UNDEFINED, TYPE_UNDEFINED}, []int{TYPE_BOOL})
	AddOpCode(OP_UND_LEN, "len", []int{TYPE_UNDEFINED}, []int{TYPE_I32})
	AddOpCode(OP_UND_PRINTF, "printf", []int{TYPE_UNDEFINED}, []int{})
	AddOpCode(OP_UND_SPRINTF, "sprintf", []int{TYPE_UNDEFINED}, []int{TYPE_STR})
	AddOpCode(OP_UND_READ, "read", []int{}, []int{TYPE_STR})

	AddOpCode(OP_BYTE_BYTE, "byte.byte", []int{TYPE_BYTE}, []int{TYPE_BYTE})
	AddOpCode(OP_BYTE_STR, "byte.str", []int{TYPE_BYTE}, []int{TYPE_STR})
	AddOpCode(OP_BYTE_I32, "byte.i32", []int{TYPE_BYTE}, []int{TYPE_I32})
	AddOpCode(OP_BYTE_I64, "byte.i64", []int{TYPE_BYTE}, []int{TYPE_I64})
	AddOpCode(OP_BYTE_F32, "byte.f32", []int{TYPE_BYTE}, []int{TYPE_F32})
	AddOpCode(OP_BYTE_F64, "byte.f64", []int{TYPE_BYTE}, []int{TYPE_F64})

	AddOpCode(OP_BYTE_PRINT, "byte.print", []int{TYPE_BYTE}, []int{})

	AddOpCode(OP_BOOL_PRINT, "bool.print", []int{TYPE_BOOL}, []int{})
	AddOpCode(OP_BOOL_EQUAL, "bool.eq", []int{TYPE_BOOL, TYPE_BOOL}, []int{TYPE_BOOL})
	AddOpCode(OP_BOOL_UNEQUAL, "bool.uneq", []int{TYPE_BOOL, TYPE_BOOL}, []int{TYPE_BOOL})
	AddOpCode(OP_BOOL_NOT, "bool.not", []int{TYPE_BOOL}, []int{TYPE_BOOL})
	AddOpCode(OP_BOOL_OR, "bool.or", []int{TYPE_BOOL, TYPE_BOOL}, []int{TYPE_BOOL})
	AddOpCode(OP_BOOL_AND, "bool.and", []int{TYPE_BOOL, TYPE_BOOL}, []int{TYPE_BOOL})

	AddOpCode(OP_I32_BYTE, "i32.byte", []int{TYPE_I32}, []int{TYPE_BYTE})
	AddOpCode(OP_I32_STR, "i32.str", []int{TYPE_I32}, []int{TYPE_STR})
	AddOpCode(OP_I32_I32, "i32.i32", []int{TYPE_I32}, []int{TYPE_I32})
	AddOpCode(OP_I32_I64, "i32.i64", []int{TYPE_I32}, []int{TYPE_I64})
	AddOpCode(OP_I32_F32, "i32.f32", []int{TYPE_I32}, []int{TYPE_F32})
	AddOpCode(OP_I32_F64, "i32.f64", []int{TYPE_I32}, []int{TYPE_F64})

	AddOpCode(OP_I32_PRINT, "i32.print", []int{TYPE_I32}, []int{})
	AddOpCode(OP_I32_ADD, "i32.add", []int{TYPE_I32, TYPE_I32}, []int{TYPE_I32})
	AddOpCode(OP_I32_SUB, "i32.sub", []int{TYPE_I32, TYPE_I32}, []int{TYPE_I32})
	AddOpCode(OP_I32_MUL, "i32.mul", []int{TYPE_I32, TYPE_I32}, []int{TYPE_I32})
	AddOpCode(OP_I32_DIV, "i32.div", []int{TYPE_I32, TYPE_I32}, []int{TYPE_I32})
	AddOpCode(OP_I32_ABS, "i32.abs", []int{TYPE_I32}, []int{TYPE_I32})
	AddOpCode(OP_I32_POW, "i32.pow", []int{TYPE_I32, TYPE_I32}, []int{TYPE_I32})
	AddOpCode(OP_I32_GT, "i32.gt", []int{TYPE_I32, TYPE_I32}, []int{TYPE_BOOL})
	AddOpCode(OP_I32_GTEQ, "i32.gteq", []int{TYPE_I32, TYPE_I32}, []int{TYPE_BOOL})
	AddOpCode(OP_I32_LT, "i32.lt", []int{TYPE_I32, TYPE_I32}, []int{TYPE_BOOL})
	AddOpCode(OP_I32_LTEQ, "i32.lteq", []int{TYPE_I32, TYPE_I32}, []int{TYPE_BOOL})
	AddOpCode(OP_I32_EQ, "i32.eq", []int{TYPE_I32, TYPE_I32}, []int{TYPE_BOOL})
	AddOpCode(OP_I32_UNEQ, "i32.uneq", []int{TYPE_I32, TYPE_I32}, []int{TYPE_BOOL})
	AddOpCode(OP_I32_MOD, "i32.mod", []int{TYPE_I32, TYPE_I32}, []int{TYPE_I32})
	AddOpCode(OP_I32_RAND, "i32.rand", []int{TYPE_I32, TYPE_I32}, []int{TYPE_I32})
	AddOpCode(OP_I32_BITAND, "i32.bitand", []int{TYPE_I32, TYPE_I32}, []int{TYPE_I32})
	AddOpCode(OP_I32_BITOR, "i32.bitor", []int{TYPE_I32, TYPE_I32}, []int{TYPE_I32})
	AddOpCode(OP_I32_BITXOR, "i32.bitxor", []int{TYPE_I32, TYPE_I32}, []int{TYPE_I32})
	AddOpCode(OP_I32_BITCLEAR, "i32.bitclear", []int{TYPE_I32, TYPE_I32}, []int{TYPE_I32})
	AddOpCode(OP_I32_BITSHL, "i32.bitshl", []int{TYPE_I32, TYPE_I32}, []int{TYPE_I32})
	AddOpCode(OP_I32_BITSHR, "i32.bitshr", []int{TYPE_I32, TYPE_I32}, []int{TYPE_I32})
	AddOpCode(OP_I32_SQRT, "i32.sqrt", []int{TYPE_I32}, []int{TYPE_I32})
	AddOpCode(OP_I32_LOG, "i32.log", []int{TYPE_I32}, []int{TYPE_I32})
	AddOpCode(OP_I32_LOG2, "i32.log2", []int{TYPE_I32}, []int{TYPE_I32})
	AddOpCode(OP_I32_LOG10, "i32.log10", []int{TYPE_I32}, []int{TYPE_I32})
	AddOpCode(OP_I32_MAX, "i32.max", []int{TYPE_I32, TYPE_I32}, []int{TYPE_I32})
	AddOpCode(OP_I32_MIN, "i32.min", []int{TYPE_I32, TYPE_I32}, []int{TYPE_I32})

	AddOpCode(OP_I64_BYTE, "i64.byte", []int{TYPE_I64}, []int{TYPE_BYTE})
	AddOpCode(OP_I64_STR, "i64.str", []int{TYPE_I64}, []int{TYPE_STR})
	AddOpCode(OP_I64_I32, "i64.i32", []int{TYPE_I64}, []int{TYPE_I32})
	AddOpCode(OP_I64_I64, "i64.i64", []int{TYPE_I64}, []int{TYPE_I64})
	AddOpCode(OP_I64_F32, "i64.f32", []int{TYPE_I64}, []int{TYPE_F32})
	AddOpCode(OP_I64_F64, "i64.f64", []int{TYPE_I64}, []int{TYPE_F64})

	AddOpCode(OP_I64_PRINT, "i64.print", []int{TYPE_I64}, []int{})
	AddOpCode(OP_I64_ADD, "i64.add", []int{TYPE_I64, TYPE_I64}, []int{TYPE_I64})
	AddOpCode(OP_I64_SUB, "i64.sub", []int{TYPE_I64, TYPE_I64}, []int{TYPE_I64})
	AddOpCode(OP_I64_MUL, "i64.mul", []int{TYPE_I64, TYPE_I64}, []int{TYPE_I64})
	AddOpCode(OP_I64_DIV, "i64.div", []int{TYPE_I64, TYPE_I64}, []int{TYPE_I64})
	AddOpCode(OP_I64_ABS, "i64.abs", []int{TYPE_I64}, []int{TYPE_I64})
	AddOpCode(OP_I64_POW, "i64.pow", []int{TYPE_I64, TYPE_I64}, []int{TYPE_I64})
	AddOpCode(OP_I64_GT, "i64.gt", []int{TYPE_I64, TYPE_I64}, []int{TYPE_BOOL})
	AddOpCode(OP_I64_GTEQ, "i64.gteq", []int{TYPE_I64, TYPE_I64}, []int{TYPE_BOOL})
	AddOpCode(OP_I64_LT, "i64.lt", []int{TYPE_I64, TYPE_I64}, []int{TYPE_BOOL})
	AddOpCode(OP_I64_LTEQ, "i64.lteq", []int{TYPE_I64, TYPE_I64}, []int{TYPE_BOOL})
	AddOpCode(OP_I64_EQ, "i64.eq", []int{TYPE_I64, TYPE_I64}, []int{TYPE_BOOL})
	AddOpCode(OP_I64_UNEQ, "i64.uneq", []int{TYPE_I64, TYPE_I64}, []int{TYPE_BOOL})
	AddOpCode(OP_I64_MOD, "i64.mod", []int{TYPE_I64, TYPE_I64}, []int{TYPE_I64})
	AddOpCode(OP_I64_RAND, "i64.rand", []int{TYPE_I64, TYPE_I64}, []int{TYPE_I64})
	AddOpCode(OP_I64_BITAND, "i64.bitand", []int{TYPE_I64, TYPE_I64}, []int{TYPE_I64})
	AddOpCode(OP_I64_BITOR, "i64.bitor", []int{TYPE_I64, TYPE_I64}, []int{TYPE_I64})
	AddOpCode(OP_I64_BITXOR, "i64.bitxor", []int{TYPE_I64, TYPE_I64}, []int{TYPE_I64})
	AddOpCode(OP_I64_BITCLEAR, "i64.bitclear", []int{TYPE_I64, TYPE_I64}, []int{TYPE_I64})
	AddOpCode(OP_I64_BITSHL, "i64.bitshl", []int{TYPE_I64, TYPE_I64}, []int{TYPE_I64})
	AddOpCode(OP_I64_BITSHR, "i64.bitshr", []int{TYPE_I64, TYPE_I64}, []int{TYPE_I64})
	AddOpCode(OP_I64_SQRT, "i64.sqrt", []int{TYPE_I64}, []int{TYPE_I64})
	AddOpCode(OP_I64_LOG, "i64.log", []int{TYPE_I64}, []int{TYPE_I64})
	AddOpCode(OP_I64_LOG2, "i64.log2", []int{TYPE_I64}, []int{TYPE_I64})
	AddOpCode(OP_I64_LOG10, "i64.log10", []int{TYPE_I64}, []int{TYPE_I64})
	AddOpCode(OP_I64_MAX, "i64.max", []int{TYPE_I64, TYPE_I64}, []int{TYPE_I64})
	AddOpCode(OP_I64_MIN, "i64.min", []int{TYPE_I64, TYPE_I64}, []int{TYPE_I64})

	AddOpCode(OP_F32_IS_NAN, "f32.isnan", []int{TYPE_F32}, []int{TYPE_BOOL})
	AddOpCode(OP_F32_BYTE, "f32.byte", []int{TYPE_F32}, []int{TYPE_BYTE})
	AddOpCode(OP_F32_STR, "f32.str", []int{TYPE_F32}, []int{TYPE_STR})
	AddOpCode(OP_F32_I32, "f32.i32", []int{TYPE_F32}, []int{TYPE_I32})
	AddOpCode(OP_F32_I64, "f32.i64", []int{TYPE_F32}, []int{TYPE_I64})
	AddOpCode(OP_F32_F32, "f32.f32", []int{TYPE_F32}, []int{TYPE_F32})
	AddOpCode(OP_F32_F64, "f32.f64", []int{TYPE_F32}, []int{TYPE_F64})
	AddOpCode(OP_F32_PRINT, "f32.print", []int{TYPE_F32}, []int{})
	AddOpCode(OP_F32_ADD, "f32.add", []int{TYPE_F32, TYPE_F32}, []int{TYPE_F32})
	AddOpCode(OP_F32_SUB, "f32.sub", []int{TYPE_F32, TYPE_F32}, []int{TYPE_F32})
	AddOpCode(OP_F32_MUL, "f32.mul", []int{TYPE_F32, TYPE_F32}, []int{TYPE_F32})
	AddOpCode(OP_F32_DIV, "f32.div", []int{TYPE_F32, TYPE_F32}, []int{TYPE_F32})
	AddOpCode(OP_F32_ABS, "f32.abs", []int{TYPE_F32}, []int{TYPE_F32})
	AddOpCode(OP_F32_POW, "f32.pow", []int{TYPE_F32, TYPE_F32}, []int{TYPE_F32})
	AddOpCode(OP_F32_GT, "f32.gt", []int{TYPE_F32, TYPE_F32}, []int{TYPE_BOOL})
	AddOpCode(OP_F32_GTEQ, "f32.gteq", []int{TYPE_F32, TYPE_F32}, []int{TYPE_BOOL})
	AddOpCode(OP_F32_LT, "f32.lt", []int{TYPE_F32, TYPE_F32}, []int{TYPE_BOOL})
	AddOpCode(OP_F32_LTEQ, "f32.lteq", []int{TYPE_F32, TYPE_F32}, []int{TYPE_BOOL})
	AddOpCode(OP_F32_EQ, "f32.eq", []int{TYPE_F32, TYPE_F32}, []int{TYPE_BOOL})
	AddOpCode(OP_F32_UNEQ, "f32.uneq", []int{TYPE_F32, TYPE_F32}, []int{TYPE_BOOL})
	AddOpCode(OP_F32_COS, "f32.cos", []int{TYPE_F32}, []int{TYPE_F32})
	AddOpCode(OP_F32_SIN, "f32.sin", []int{TYPE_F32}, []int{TYPE_F32})
	AddOpCode(OP_F32_SQRT, "f32.sqrt", []int{TYPE_F32}, []int{TYPE_F32})
	AddOpCode(OP_F32_LOG, "f32.log", []int{TYPE_F32}, []int{TYPE_F32})
	AddOpCode(OP_F32_LOG2, "f32.log2", []int{TYPE_F32}, []int{TYPE_F32})
	AddOpCode(OP_F32_LOG10, "f32.log10", []int{TYPE_F32}, []int{TYPE_F32})
	AddOpCode(OP_F32_MAX, "f32.max", []int{TYPE_F32, TYPE_F32}, []int{TYPE_F32})
	AddOpCode(OP_F32_MIN, "f32.min", []int{TYPE_F32, TYPE_F32}, []int{TYPE_F32})

	AddOpCode(OP_F64_BYTE, "f64.byte", []int{TYPE_F64}, []int{TYPE_BYTE})
	AddOpCode(OP_F64_STR, "f64.str", []int{TYPE_F64}, []int{TYPE_STR})
	AddOpCode(OP_F64_I32, "f64.i32", []int{TYPE_F64}, []int{TYPE_I32})
	AddOpCode(OP_F64_I64, "f64.i64", []int{TYPE_F64}, []int{TYPE_I64})
	AddOpCode(OP_F64_F32, "f64.f32", []int{TYPE_F64}, []int{TYPE_F32})
	AddOpCode(OP_F64_F64, "f64.f64", []int{TYPE_F64}, []int{TYPE_F64})

	AddOpCode(OP_F64_PRINT, "f64.print", []int{TYPE_F64}, []int{})
	AddOpCode(OP_F64_ADD, "f64.add", []int{TYPE_F64, TYPE_F64}, []int{TYPE_F64})
	AddOpCode(OP_F64_SUB, "f64.sub", []int{TYPE_F64, TYPE_F64}, []int{TYPE_F64})
	AddOpCode(OP_F64_MUL, "f64.mul", []int{TYPE_F64, TYPE_F64}, []int{TYPE_F64})
	AddOpCode(OP_F64_DIV, "f64.div", []int{TYPE_F64, TYPE_F64}, []int{TYPE_F64})
	AddOpCode(OP_F64_ABS, "f64.abs", []int{TYPE_F64}, []int{TYPE_F64})
	AddOpCode(OP_F64_POW, "f64.pow", []int{TYPE_F64, TYPE_F64}, []int{TYPE_F64})
	AddOpCode(OP_F64_GT, "f64.gt", []int{TYPE_F64, TYPE_F64}, []int{TYPE_BOOL})
	AddOpCode(OP_F64_GTEQ, "f64.gteq", []int{TYPE_F64, TYPE_F64}, []int{TYPE_BOOL})
	AddOpCode(OP_F64_LT, "f64.lt", []int{TYPE_F64, TYPE_F64}, []int{TYPE_BOOL})
	AddOpCode(OP_F64_LTEQ, "f64.lteq", []int{TYPE_F64, TYPE_F64}, []int{TYPE_BOOL})
	AddOpCode(OP_F64_EQ, "f64.eq", []int{TYPE_F64, TYPE_F64}, []int{TYPE_BOOL})
	AddOpCode(OP_F64_UNEQ, "f64.uneq", []int{TYPE_F64, TYPE_F64}, []int{TYPE_BOOL})
	AddOpCode(OP_F64_COS, "f64.cos", []int{TYPE_F64}, []int{TYPE_F64})
	AddOpCode(OP_F64_SIN, "f64.sin", []int{TYPE_F64}, []int{TYPE_F64})
	AddOpCode(OP_F64_SQRT, "f64.sqrt", []int{TYPE_F64}, []int{TYPE_F64})
	AddOpCode(OP_F64_LOG, "f64.log", []int{TYPE_F64}, []int{TYPE_F64})
	AddOpCode(OP_F64_LOG2, "f64.log2", []int{TYPE_F64}, []int{TYPE_F64})
	AddOpCode(OP_F64_LOG10, "f64.log10", []int{TYPE_F64}, []int{TYPE_F64})
	AddOpCode(OP_F64_MAX, "f64.max", []int{TYPE_F64, TYPE_F64}, []int{TYPE_F64})
	AddOpCode(OP_F64_MIN, "f64.min", []int{TYPE_F64, TYPE_F64}, []int{TYPE_F64})

	AddOpCode(OP_STR_PRINT, "str.print", []int{TYPE_STR}, []int{})
	AddOpCode(OP_STR_CONCAT, "str.concat", []int{TYPE_STR, TYPE_STR}, []int{TYPE_STR})
	AddOpCode(OP_STR_SUBSTR, "str.substr", []int{TYPE_STR, TYPE_I32, TYPE_I32}, []int{TYPE_STR})
	AddOpCode(OP_STR_INDEX, "str.index", []int{TYPE_STR, TYPE_STR}, []int{TYPE_I32})
	AddOpCode(OP_STR_TRIM_SPACE, "str.trimspace", []int{TYPE_STR}, []int{TYPE_STR})
	AddOpCode(OP_STR_EQ, "str.eq", []int{TYPE_STR, TYPE_STR}, []int{TYPE_BOOL})

	AddOpCode(OP_STR_BYTE, "str.byte", []int{TYPE_STR}, []int{TYPE_BYTE})
	AddOpCode(OP_STR_STR, "str.str", []int{TYPE_STR}, []int{TYPE_STR})
	AddOpCode(OP_STR_I32, "str.i32", []int{TYPE_STR}, []int{TYPE_I32})
	AddOpCode(OP_STR_I64, "str.i64", []int{TYPE_STR}, []int{TYPE_I64})
	AddOpCode(OP_STR_F32, "str.f32", []int{TYPE_STR}, []int{TYPE_F32})
	AddOpCode(OP_STR_F64, "str.f64", []int{TYPE_STR}, []int{TYPE_F64})

	AddOpCode(OP_APPEND, "append", []int{TYPE_UNDEFINED, TYPE_UNDEFINED}, []int{TYPE_UNDEFINED})
	AddOpCode(OP_ASSERT, "assert", []int{TYPE_UNDEFINED, TYPE_UNDEFINED, TYPE_STR}, []int{TYPE_BOOL})
	AddOpCode(OP_TEST, "test", []int{TYPE_UNDEFINED, TYPE_UNDEFINED, TYPE_STR}, []int{})
	AddOpCode(OP_PANIC, "panic", []int{TYPE_UNDEFINED, TYPE_UNDEFINED, TYPE_STR}, []int{})

	// affordances
	AddOpCode(OP_AFF_PRINT, "aff.print", []int{TYPE_AFF}, []int{})
	AddOpCode(OP_AFF_QUERY, "aff.query", []int{TYPE_AFF}, []int{TYPE_AFF})
	AddOpCode(OP_AFF_ON, "aff.on", []int{TYPE_AFF, TYPE_AFF}, []int{})
	AddOpCode(OP_AFF_OF, "aff.of", []int{TYPE_AFF, TYPE_AFF}, []int{})
	AddOpCode(OP_AFF_INFORM, "aff.inform", []int{TYPE_AFF, TYPE_I32, TYPE_AFF}, []int{})
	AddOpCode(OP_AFF_REQUEST, "aff.request", []int{TYPE_AFF, TYPE_I32, TYPE_AFF}, []int{})

	// exec
	execNativeBare = func(prgrm *CXProgram) {
		defer RuntimeError()
		call := &prgrm.CallStack[prgrm.CallCounter]
		expr := call.Operator.Expressions[call.Line]
		opCode := expr.Operator.OpCode
		fp := call.FramePointer

		switch opCode {
		case OP_IDENTITY:
			opIdentity(expr, fp)
		case OP_JMP:
			opJmp(expr, fp, call)
		case OP_DEBUG:
			prgrm.PrintStack()

		case OP_SERIALIZE:
			op_serialize(expr, fp)
		case OP_DESERIALIZE:
			op_deserialize(expr, fp)

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
			opByteByte(expr, fp)
		case OP_BYTE_STR:
			opByteByte(expr, fp)
		case OP_BYTE_I32:
			opByteByte(expr, fp)
		case OP_BYTE_I64:
			opByteByte(expr, fp)
		case OP_BYTE_F32:
			opByteByte(expr, fp)
		case OP_BYTE_F64:
			opByteByte(expr, fp)

		case OP_BYTE_PRINT:
			opBytePrint(expr, fp)

		case OP_BOOL_PRINT:
			opBoolPrint(expr, fp)
		case OP_BOOL_EQUAL:
			opBoolEqual(expr, fp)
		case OP_BOOL_UNEQUAL:
			opBoolUnequal(expr, fp)
		case OP_BOOL_NOT:
			opBoolNot(expr, fp)
		case OP_BOOL_OR:
			opBoolOr(expr, fp)
		case OP_BOOL_AND:
			opBoolAnd(expr, fp)

		case OP_I32_BYTE:
			opI32I32(expr, fp)
		case OP_I32_STR:
			opI32I32(expr, fp)
		case OP_I32_I32:
			opI32I32(expr, fp)
		case OP_I32_I64:
			opI32I32(expr, fp)
		case OP_I32_F32:
			opI32I32(expr, fp)
		case OP_I32_F64:
			opI32I32(expr, fp)

		case OP_I32_PRINT:
			opI32Print(expr, fp)
		case OP_I32_ADD:
			opI32Add(expr, fp)
		case OP_I32_SUB:
			opI32Sub(expr, fp)
		case OP_I32_MUL:
			opI32Mul(expr, fp)
		case OP_I32_DIV:
			opI32Div(expr, fp)
		case OP_I32_ABS:
			opI32Abs(expr, fp)
		case OP_I32_POW:
			opI32Pow(expr, fp)
		case OP_I32_GT:
			opI32Gt(expr, fp)
		case OP_I32_GTEQ:
			opI32Gteq(expr, fp)
		case OP_I32_LT:
			opI32Lt(expr, fp)
		case OP_I32_LTEQ:
			opI32Lteq(expr, fp)
		case OP_I32_EQ:
			opI32Eq(expr, fp)
		case OP_I32_UNEQ:
			opI32Uneq(expr, fp)
		case OP_I32_MOD:
			opI32Mod(expr, fp)
		case OP_I32_RAND:
			opI32Rand(expr, fp)
		case OP_I32_BITAND:
			opI32Bitand(expr, fp)
		case OP_I32_BITOR:
			opI32Bitor(expr, fp)
		case OP_I32_BITXOR:
			opI32Bitxor(expr, fp)
		case OP_I32_BITCLEAR:
			opI32Bitclear(expr, fp)
		case OP_I32_BITSHL:
			opI32Bitshl(expr, fp)
		case OP_I32_BITSHR:
			opI32Bitshr(expr, fp)
		case OP_I32_SQRT:
			opI32Sqrt(expr, fp)
		case OP_I32_LOG:
			opI32Log(expr, fp)
		case OP_I32_LOG2:
			opI32Log2(expr, fp)
		case OP_I32_LOG10:
			opI32Log10(expr, fp)

		case OP_I32_MAX:
			opI32Max(expr, fp)
		case OP_I32_MIN:
			opI32Min(expr, fp)

		case OP_I64_BYTE:
			opI64I64(expr, fp)
		case OP_I64_STR:
			opI64I64(expr, fp)
		case OP_I64_I32:
			opI64I64(expr, fp)
		case OP_I64_I64:
			opI64I64(expr, fp)
		case OP_I64_F32:
			opI64I64(expr, fp)
		case OP_I64_F64:
			opI64I64(expr, fp)

		case OP_I64_PRINT:
			opI64Print(expr, fp)
		case OP_I64_ADD:
			opI64Add(expr, fp)
		case OP_I64_SUB:
			opI64Sub(expr, fp)
		case OP_I64_MUL:
			opI64Mul(expr, fp)
		case OP_I64_DIV:
			opI64Div(expr, fp)
		case OP_I64_ABS:
			opI64Abs(expr, fp)
		case OP_I64_POW:
			opI64Pow(expr, fp)
		case OP_I64_GT:
			opI64Gt(expr, fp)
		case OP_I64_GTEQ:
			opI64Gteq(expr, fp)
		case OP_I64_LT:
			opI64Lt(expr, fp)
		case OP_I64_LTEQ:
			opI64Lteq(expr, fp)
		case OP_I64_EQ:
			opI64Eq(expr, fp)
		case OP_I64_UNEQ:
			opI64Uneq(expr, fp)
		case OP_I64_MOD:
			opI64Mod(expr, fp)
		case OP_I64_RAND:
			opI64Rand(expr, fp)
		case OP_I64_BITAND:
			opI64Bitand(expr, fp)
		case OP_I64_BITOR:
			opI64Bitor(expr, fp)
		case OP_I64_BITXOR:
			opI64Bitxor(expr, fp)
		case OP_I64_BITCLEAR:
			opI64Bitclear(expr, fp)
		case OP_I64_BITSHL:
			opI64Bitshl(expr, fp)
		case OP_I64_BITSHR:
			opI64Bitshr(expr, fp)
		case OP_I64_SQRT:
			opI64Sqrt(expr, fp)
		case OP_I64_LOG:
			opI64Log(expr, fp)
		case OP_I64_LOG2:
			opI64Log2(expr, fp)
		case OP_I64_LOG10:
			opI64Log10(expr, fp)
		case OP_I64_MAX:
			opI64Max(expr, fp)
		case OP_I64_MIN:
			opI64Min(expr, fp)

		case OP_F32_IS_NAN:
			opF32Isnan(expr, fp)
		case OP_F32_BYTE:
			opF32F32(expr, fp)
		case OP_F32_STR:
			opF32F32(expr, fp)
		case OP_F32_I32:
			opF32F32(expr, fp)
		case OP_F32_I64:
			opF32F32(expr, fp)
		case OP_F32_F32:
			opF32F32(expr, fp)
		case OP_F32_F64:
			opF32F32(expr, fp)
		case OP_F32_PRINT:
			opF32Print(expr, fp)
		case OP_F32_ADD:
			opF32Add(expr, fp)
		case OP_F32_SUB:
			opF32Sub(expr, fp)
		case OP_F32_MUL:
			opF32Mul(expr, fp)
		case OP_F32_DIV:
			opF32Div(expr, fp)
		case OP_F32_ABS:
			opF32Abs(expr, fp)
		case OP_F32_POW:
			opF32Pow(expr, fp)
		case OP_F32_GT:
			opF32Gt(expr, fp)
		case OP_F32_GTEQ:
			opF32Gteq(expr, fp)
		case OP_F32_LT:
			opF32Lt(expr, fp)
		case OP_F32_LTEQ:
			opF32Lteq(expr, fp)
		case OP_F32_EQ:
			opF32Eq(expr, fp)
		case OP_F32_UNEQ:
			opF32Uneq(expr, fp)
		case OP_F32_COS:
			opF32Cos(expr, fp)
		case OP_F32_SIN:
			opF32Sin(expr, fp)
		case OP_F32_SQRT:
			opF32Sqrt(expr, fp)
		case OP_F32_LOG:
			opF32Log(expr, fp)
		case OP_F32_LOG2:
			opF32Log2(expr, fp)
		case OP_F32_LOG10:
			opF32Log10(expr, fp)
		case OP_F32_MAX:
			opF32Max(expr, fp)
		case OP_F32_MIN:
			opF32Min(expr, fp)

		case OP_F64_BYTE:
			opF64F64(expr, fp)
		case OP_F64_STR:
			opF64F64(expr, fp)
		case OP_F64_I32:
			opF64F64(expr, fp)
		case OP_F64_I64:
			opF64F64(expr, fp)
		case OP_F64_F32:
			opF64F64(expr, fp)
		case OP_F64_F64:
			opF64F64(expr, fp)

		case OP_F64_PRINT:
			opF64Print(expr, fp)
		case OP_F64_ADD:
			opF64Add(expr, fp)
		case OP_F64_SUB:
			opF64Sub(expr, fp)
		case OP_F64_MUL:
			opF64Mul(expr, fp)
		case OP_F64_DIV:
			opF64Div(expr, fp)
		case OP_F64_ABS:
			opF64Abs(expr, fp)
		case OP_F64_POW:
			opF64Pow(expr, fp)
		case OP_F64_GT:
			opF64Gt(expr, fp)
		case OP_F64_GTEQ:
			opF64Gteq(expr, fp)
		case OP_F64_LT:
			opF64Lt(expr, fp)
		case OP_F64_LTEQ:
			opF64Lteq(expr, fp)
		case OP_F64_EQ:
			opF64Eq(expr, fp)
		case OP_F64_UNEQ:
			opF64Uneq(expr, fp)
		case OP_F64_COS:
			opF64Cos(expr, fp)
		case OP_F64_SIN:
			opF64Sin(expr, fp)
		case OP_F64_SQRT:
			opF64Sqrt(expr, fp)
		case OP_F64_LOG:
			opF64Log(expr, fp)
		case OP_F64_LOG2:
			opF64Log2(expr, fp)
		case OP_F64_LOG10:
			opF64Log10(expr, fp)
		case OP_F64_MAX:
			opF64Max(expr, fp)
		case OP_F64_MIN:
			opF64Min(expr, fp)

		case OP_STR_PRINT:
			opStrPrint(expr, fp)
		case OP_STR_EQ:
			opStrEq(expr, fp)
		case OP_STR_CONCAT:
			opStrConcat(expr, fp)
		case OP_STR_SUBSTR:
			opStrSubstr(expr, fp)
		case OP_STR_INDEX:
			opStrIndex(expr, fp)
		case OP_STR_TRIM_SPACE:
			opStrTrimSpace(expr, fp)

		case OP_STR_BYTE:
			opStrStr(expr, fp)
		case OP_STR_STR:
			opStrStr(expr, fp)
		case OP_STR_I32:
			opStrStr(expr, fp)
		case OP_STR_I64:
			opStrStr(expr, fp)
		case OP_STR_F32:
			opStrStr(expr, fp)
		case OP_STR_F64:
			opStrStr(expr, fp)

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
		case OP_TEST:
			op_test(expr, fp)
		case OP_PANIC:
			op_panic(expr, fp)

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
		default:
			// DumpOpCodes(opCode) // debug helper
			panic("invalid bare opcode")
		}
	}

	execNative = execNativeBare
}
