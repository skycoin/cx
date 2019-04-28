package cxcore

// CorePackages ...
var CorePackages = []string{
	// temporary solution until we can implement these packages in pure CX I guess
	"gl", "glfw", "time", "http", "os", "explorer", "aff", "gltext", "cx",
}

// op codes
// nolint golint
const (
	OP_IDENTITY = iota + 1
	OP_JMP
	OP_DEBUG

	OP_SERIALIZE
	OP_DESERIALIZE

	START_UND_OPS

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

	END_UND_OPS

	// parse ops

	START_PARSE_OPS

	OP_I32_STR
	OP_I32_I32
	OP_I32_I64
	OP_I32_BYTE
	OP_I32_UI8
	OP_I32_UI16
	OP_I32_UI32
	OP_I32_UI64
	OP_I32_F32
	OP_I32_F64

	OP_I64_STR
	OP_I64_I32
	OP_I64_I64
	OP_I64_BYTE
	OP_I64_UI8
	OP_I64_UI16
	OP_I64_UI32
	OP_I64_UI64
	OP_I64_F32
	OP_I64_F64

	OP_BYTE_STR
	OP_BYTE_I32
	OP_BYTE_I64
	OP_BYTE_BYTE
	OP_BYTE_UI8
	OP_BYTE_UI16
	OP_BYTE_UI32
	OP_BYTE_UI64
	OP_BYTE_F32
	OP_BYTE_F64
	OP_BYTE_PRINT

	OP_UI8_STR
	OP_UI8_I32
	OP_UI8_I64
	OP_UI8_UI16
	OP_UI8_UI32
	OP_UI8_UI64
	OP_UI8_F32
	OP_UI8_F64

	OP_UI16_STR
	OP_UI16_I32
	OP_UI16_I64
	OP_UI16_UI8
	OP_UI16_UI32
	OP_UI16_UI64
	OP_UI16_F32
	OP_UI16_F64

	OP_UI32_STR
	OP_UI32_I32
	OP_UI32_I64
	OP_UI32_UI8
	OP_UI32_UI16
	OP_UI32_UI64
	OP_UI32_F32
	OP_UI32_F64

	OP_UI64_STR
	OP_UI64_I32
	OP_UI64_I64
	OP_UI64_UI8
	OP_UI64_UI16
	OP_UI64_UI32
	OP_UI64_F32
	OP_UI64_F64

	OP_F32_STR
	OP_F32_I32
	OP_F32_I64
	OP_F32_BYTE
	OP_F32_UI8
	OP_F32_UI16
	OP_F32_UI32
	OP_F32_UI64
	OP_F32_F32
	OP_F32_F64

	OP_F64_STR
	OP_F64_I32
	OP_F64_I64
	OP_F64_BYTE
	OP_F64_UI8
	OP_F64_UI16
	OP_F64_UI32
	OP_F64_UI64
	OP_F64_F32
	OP_F64_F64

	OP_STR_STR
	OP_STR_I32
	OP_STR_I64
	OP_STR_BYTE
	OP_STR_UI8
	OP_STR_UI16
	OP_STR_UI32
	OP_STR_UI64
	OP_STR_F32
	OP_STR_F64

	END_PARSE_OPS

	OP_BOOL_PRINT

	OP_BOOL_EQUAL
	OP_BOOL_UNEQUAL
	OP_BOOL_NOT
	OP_BOOL_OR
	OP_BOOL_AND

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

	OP_UI8_ADD
	OP_UI8_SUB
	OP_UI8_MUL
	OP_UI8_DIV
	OP_UI8_GT
	OP_UI8_GTEQ
	OP_UI8_LT
	OP_UI8_LTEQ
	OP_UI8_EQ
	OP_UI8_UNEQ
	OP_UI8_MOD
	OP_UI8_RAND
	OP_UI8_BITAND
	OP_UI8_BITOR
	OP_UI8_BITXOR
	OP_UI8_BITCLEAR
	OP_UI8_BITSHL
	OP_UI8_BITSHR
	OP_UI8_MAX
	OP_UI8_MIN

	OP_UI16_ADD
	OP_UI16_SUB
	OP_UI16_MUL
	OP_UI16_DIV
	OP_UI16_GT
	OP_UI16_GTEQ
	OP_UI16_LT
	OP_UI16_LTEQ
	OP_UI16_EQ
	OP_UI16_UNEQ
	OP_UI16_MOD
	OP_UI16_RAND
	OP_UI16_BITAND
	OP_UI16_BITOR
	OP_UI16_BITXOR
	OP_UI16_BITCLEAR
	OP_UI16_BITSHL
	OP_UI16_BITSHR
	OP_UI16_MAX
	OP_UI16_MIN

	OP_UI32_ADD
	OP_UI32_SUB
	OP_UI32_MUL
	OP_UI32_DIV
	OP_UI32_GT
	OP_UI32_GTEQ
	OP_UI32_LT
	OP_UI32_LTEQ
	OP_UI32_EQ
	OP_UI32_UNEQ
	OP_UI32_MOD
	OP_UI32_RAND
	OP_UI32_BITAND
	OP_UI32_BITOR
	OP_UI32_BITXOR
	OP_UI32_BITCLEAR
	OP_UI32_BITSHL
	OP_UI32_BITSHR
	OP_UI32_MAX
	OP_UI32_MIN

	OP_UI64_ADD
	OP_UI64_SUB
	OP_UI64_MUL
	OP_UI64_DIV
	OP_UI64_GT
	OP_UI64_GTEQ
	OP_UI64_LT
	OP_UI64_LTEQ
	OP_UI64_EQ
	OP_UI64_UNEQ
	OP_UI64_MOD
	OP_UI64_RAND
	OP_UI64_BITAND
	OP_UI64_BITOR
	OP_UI64_BITXOR
	OP_UI64_BITCLEAR
	OP_UI64_BITSHL
	OP_UI64_BITSHR
	OP_UI64_MAX
	OP_UI64_MIN

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
	OP_F32_IS_NAN

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

	OP_APPEND
	OP_RESIZE
	OP_INSERT
	OP_REMOVE
	OP_COPY

	OP_MAKE
	OP_READ
	OP_WRITE
	OP_LEN
	OP_CONCAT
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
	OP_STRERROR

	// affordances
	OP_AFF_PRINT
	OP_AFF_QUERY
	OP_AFF_ON
	OP_AFF_OF
	OP_AFF_INFORM
	OP_AFF_REQUEST

	OP_UND_NEG
	OP_I32_NEG
	OP_I64_NEG
	OP_F32_NEG
	OP_F64_NEG
	END_OF_CORE_OPS
)

// For the parser. These shouldn't be used in the runtime for performance reasons
var (
	OpNames        = map[int]string{}
	OpCodes        = map[string]int{}
	Natives        = map[int]*CXFunction{}
	execNativeCore func(*CXProgram)
	execNative     func(*CXProgram)
)

// AddOpCode ...
func AddOpCode(code int, name string, inputs []*CXArgument, outputs []*CXArgument) {
	OpNames[code] = name
	OpCodes[name] = code
	Natives[code] = MakeNativeFunction(code, inputs, outputs)
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

// Helper function for creating parameters for standard library operators.
// The current standard library only uses basic types and slices. If more options are needed, modify this function
func newOpPar(typCode int, isSlice bool) *CXArgument {
	arg := MakeArgument("", "", -1).AddType(TypeNames[typCode])
	if isSlice {
		arg.IsSlice = true
		arg.DeclarationSpecifiers = append(arg.DeclarationSpecifiers, DECL_SLICE)
	}
	return arg
}

func init() {
	AddOpCode(OP_IDENTITY, "identity",
		[]*CXArgument{newOpPar(TYPE_UNDEFINED, false)},
		[]*CXArgument{newOpPar(TYPE_UNDEFINED, false)})
	AddOpCode(OP_JMP, "jmp",
		[]*CXArgument{newOpPar(TYPE_BOOL, false)},
		[]*CXArgument{}) // newOpPar(TYPE_UNDEFINED, false) to allow 0 inputs (goto)
	AddOpCode(OP_DEBUG, "debug",
		[]*CXArgument{},
		[]*CXArgument{})

	AddOpCode(OP_SERIALIZE, "serialize",
		[]*CXArgument{newOpPar(TYPE_AFF, false)},
		[]*CXArgument{newOpPar(TYPE_BYTE, false)})
	AddOpCode(OP_DESERIALIZE, "deserialize",
		[]*CXArgument{newOpPar(TYPE_BYTE, false)},
		[]*CXArgument{})

	AddOpCode(OP_UND_EQUAL, "eq",
		[]*CXArgument{newOpPar(TYPE_UNDEFINED, false), newOpPar(TYPE_UNDEFINED, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_UND_UNEQUAL, "uneq",
		[]*CXArgument{newOpPar(TYPE_UNDEFINED, false), newOpPar(TYPE_UNDEFINED, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_UND_BITAND, "bitand",
		[]*CXArgument{newOpPar(TYPE_UNDEFINED, false), newOpPar(TYPE_UNDEFINED, false)},
		[]*CXArgument{newOpPar(TYPE_UNDEFINED, false)})
	AddOpCode(OP_UND_BITXOR, "bitxor",
		[]*CXArgument{newOpPar(TYPE_UNDEFINED, false), newOpPar(TYPE_UNDEFINED, false)},
		[]*CXArgument{newOpPar(TYPE_UNDEFINED, false)})
	AddOpCode(OP_UND_BITOR, "bitor",
		[]*CXArgument{newOpPar(TYPE_UNDEFINED, false), newOpPar(TYPE_UNDEFINED, false)},
		[]*CXArgument{newOpPar(TYPE_UNDEFINED, false)})
	AddOpCode(OP_UND_BITCLEAR, "bitclear",
		[]*CXArgument{newOpPar(TYPE_UNDEFINED, false), newOpPar(TYPE_UNDEFINED, false)},
		[]*CXArgument{newOpPar(TYPE_UNDEFINED, false)})
	AddOpCode(OP_UND_MUL, "mul",
		[]*CXArgument{newOpPar(TYPE_UNDEFINED, false), newOpPar(TYPE_UNDEFINED, false)},
		[]*CXArgument{newOpPar(TYPE_UNDEFINED, false)})
	AddOpCode(OP_UND_DIV, "div",
		[]*CXArgument{newOpPar(TYPE_UNDEFINED, false), newOpPar(TYPE_UNDEFINED, false)},
		[]*CXArgument{newOpPar(TYPE_UNDEFINED, false)})
	AddOpCode(OP_UND_MOD, "mod",
		[]*CXArgument{newOpPar(TYPE_UNDEFINED, false), newOpPar(TYPE_UNDEFINED, false)},
		[]*CXArgument{newOpPar(TYPE_UNDEFINED, false)})
	AddOpCode(OP_UND_ADD, "add",
		[]*CXArgument{newOpPar(TYPE_UNDEFINED, false), newOpPar(TYPE_UNDEFINED, false)},
		[]*CXArgument{newOpPar(TYPE_UNDEFINED, false)})
	AddOpCode(OP_UND_SUB, "sub",
		[]*CXArgument{newOpPar(TYPE_UNDEFINED, false), newOpPar(TYPE_UNDEFINED, false)},
		[]*CXArgument{newOpPar(TYPE_UNDEFINED, false)})
	AddOpCode(OP_UND_NEG, "neg",
		[]*CXArgument{newOpPar(TYPE_UNDEFINED, false)},
		[]*CXArgument{newOpPar(TYPE_UNDEFINED, false)})
	AddOpCode(OP_UND_BITSHL, "bitshl",
		[]*CXArgument{newOpPar(TYPE_UNDEFINED, false), newOpPar(TYPE_UNDEFINED, false)},
		[]*CXArgument{newOpPar(TYPE_UNDEFINED, false)})
	AddOpCode(OP_UND_BITSHR, "bitshr",
		[]*CXArgument{newOpPar(TYPE_UNDEFINED, false), newOpPar(TYPE_UNDEFINED, false)},
		[]*CXArgument{newOpPar(TYPE_UNDEFINED, false)})
	AddOpCode(OP_UND_LT, "lt",
		[]*CXArgument{newOpPar(TYPE_UNDEFINED, false), newOpPar(TYPE_UNDEFINED, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_UND_GT, "gt",
		[]*CXArgument{newOpPar(TYPE_UNDEFINED, false), newOpPar(TYPE_UNDEFINED, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_UND_LTEQ, "lteq",
		[]*CXArgument{newOpPar(TYPE_UNDEFINED, false), newOpPar(TYPE_UNDEFINED, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_UND_GTEQ, "gteq",
		[]*CXArgument{newOpPar(TYPE_UNDEFINED, false), newOpPar(TYPE_UNDEFINED, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_UND_LEN, "len",
		[]*CXArgument{newOpPar(TYPE_UNDEFINED, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false)})
	AddOpCode(OP_UND_PRINTF, "printf",
		[]*CXArgument{newOpPar(TYPE_UNDEFINED, false)},
		[]*CXArgument{})
	AddOpCode(OP_UND_SPRINTF, "sprintf",
		[]*CXArgument{newOpPar(TYPE_UNDEFINED, false)},
		[]*CXArgument{newOpPar(TYPE_STR, false)})
	AddOpCode(OP_UND_READ, "read",
		[]*CXArgument{},
		[]*CXArgument{newOpPar(TYPE_STR, false)})

	AddOpCode(OP_BOOL_PRINT, "bool.print",
		[]*CXArgument{newOpPar(TYPE_BOOL, false)},
		[]*CXArgument{})
	AddOpCode(OP_BOOL_EQUAL, "bool.eq",
		[]*CXArgument{newOpPar(TYPE_BOOL, false), newOpPar(TYPE_BOOL, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_BOOL_UNEQUAL, "bool.uneq",
		[]*CXArgument{newOpPar(TYPE_BOOL, false), newOpPar(TYPE_BOOL, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_BOOL_NOT, "bool.not",
		[]*CXArgument{newOpPar(TYPE_BOOL, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_BOOL_OR, "bool.or",
		[]*CXArgument{newOpPar(TYPE_BOOL, false), newOpPar(TYPE_BOOL, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_BOOL_AND, "bool.and",
		[]*CXArgument{newOpPar(TYPE_BOOL, false), newOpPar(TYPE_BOOL, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})

	AddOpCode(OP_I32_STR, "i32.str",
		[]*CXArgument{newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_STR, false)})
	AddOpCode(OP_I32_I32, "i32.i32",
		[]*CXArgument{newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false)})
	AddOpCode(OP_I32_I64, "i32.i64",
		[]*CXArgument{newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_I64, false)})
	AddOpCode(OP_I32_BYTE, "i32.byte",
		[]*CXArgument{newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_BYTE, false)})
	AddOpCode(OP_I32_UI8, "i32.ui8",
		[]*CXArgument{newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_UI8, false)})
	AddOpCode(OP_I32_UI16, "i32.ui16",
		[]*CXArgument{newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_UI16, false)})
	AddOpCode(OP_I32_UI32, "i32.ui32",
		[]*CXArgument{newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_UI32, false)})
	AddOpCode(OP_I32_UI64, "i32.ui64",
		[]*CXArgument{newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_UI64, false)})
	AddOpCode(OP_I32_F32, "i32.f32",
		[]*CXArgument{newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_F32, false)})
	AddOpCode(OP_I32_F64, "i32.f64",
		[]*CXArgument{newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_F64, false)})
	AddOpCode(OP_I32_PRINT, "i32.print",
		[]*CXArgument{newOpPar(TYPE_I32, false)},
		[]*CXArgument{})
	AddOpCode(OP_I32_ADD, "i32.add",
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false)})
	AddOpCode(OP_I32_SUB, "i32.sub",
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false)})
	AddOpCode(OP_I32_NEG, "i32.neg",
		[]*CXArgument{newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false)})
	AddOpCode(OP_I32_MUL, "i32.mul",
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false)})
	AddOpCode(OP_I32_DIV, "i32.div",
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false)})
	AddOpCode(OP_I32_ABS, "i32.abs",
		[]*CXArgument{newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false)})
	AddOpCode(OP_I32_POW, "i32.pow",
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false)})
	AddOpCode(OP_I32_GT, "i32.gt",
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_I32_GTEQ, "i32.gteq",
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_I32_LT, "i32.lt",
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_I32_LTEQ, "i32.lteq",
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_I32_EQ, "i32.eq",
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_I32_UNEQ, "i32.uneq",
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_I32_MOD, "i32.mod",
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false)})
	AddOpCode(OP_I32_RAND, "i32.rand",
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false)})
	AddOpCode(OP_I32_BITAND, "i32.bitand",
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false)})
	AddOpCode(OP_I32_BITOR, "i32.bitor",
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false)})
	AddOpCode(OP_I32_BITXOR, "i32.bitxor",
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false)})
	AddOpCode(OP_I32_BITCLEAR, "i32.bitclear",
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false)})
	AddOpCode(OP_I32_BITSHL, "i32.bitshl",
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false)})
	AddOpCode(OP_I32_BITSHR, "i32.bitshr",
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false)})
	AddOpCode(OP_I32_SQRT, "i32.sqrt",
		[]*CXArgument{newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false)})
	AddOpCode(OP_I32_LOG, "i32.log",
		[]*CXArgument{newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false)})
	AddOpCode(OP_I32_LOG2, "i32.log2",
		[]*CXArgument{newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false)})
	AddOpCode(OP_I32_LOG10, "i32.log10",
		[]*CXArgument{newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false)})
	AddOpCode(OP_I32_MAX, "i32.max",
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false)})
	AddOpCode(OP_I32_MIN, "i32.min",
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false)})

	AddOpCode(OP_I64_STR, "i64.str",
		[]*CXArgument{newOpPar(TYPE_I64, false)},
		[]*CXArgument{newOpPar(TYPE_STR, false)})
	AddOpCode(OP_I64_I32, "i64.i32",
		[]*CXArgument{newOpPar(TYPE_I64, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false)})
	AddOpCode(OP_I64_I64, "i64.i64",
		[]*CXArgument{newOpPar(TYPE_I64, false)},
		[]*CXArgument{newOpPar(TYPE_I64, false)})
	AddOpCode(OP_I64_BYTE, "i64.byte",
		[]*CXArgument{newOpPar(TYPE_I64, false)},
		[]*CXArgument{newOpPar(TYPE_BYTE, false)})
	AddOpCode(OP_I64_UI8, "i64.ui8",
		[]*CXArgument{newOpPar(TYPE_I64, false)},
		[]*CXArgument{newOpPar(TYPE_UI8, false)})
	AddOpCode(OP_I64_UI16, "i64.ui16",
		[]*CXArgument{newOpPar(TYPE_I64, false)},
		[]*CXArgument{newOpPar(TYPE_UI16, false)})
	AddOpCode(OP_I64_UI32, "i64.ui32",
		[]*CXArgument{newOpPar(TYPE_I64, false)},
		[]*CXArgument{newOpPar(TYPE_UI32, false)})
	AddOpCode(OP_I64_UI64, "i64.ui64",
		[]*CXArgument{newOpPar(TYPE_I64, false)},
		[]*CXArgument{newOpPar(TYPE_UI64, false)})
	AddOpCode(OP_I64_F32, "i64.f32",
		[]*CXArgument{newOpPar(TYPE_I64, false)},
		[]*CXArgument{newOpPar(TYPE_F32, false)})
	AddOpCode(OP_I64_F64, "i64.f64",
		[]*CXArgument{newOpPar(TYPE_I64, false)},
		[]*CXArgument{newOpPar(TYPE_F64, false)})
	AddOpCode(OP_I64_PRINT, "i64.print",
		[]*CXArgument{newOpPar(TYPE_I64, false)},
		[]*CXArgument{})
	AddOpCode(OP_I64_ADD, "i64.add",
		[]*CXArgument{newOpPar(TYPE_I64, false), newOpPar(TYPE_I64, false)},
		[]*CXArgument{newOpPar(TYPE_I64, false)})
	AddOpCode(OP_I64_SUB, "i64.sub",
		[]*CXArgument{newOpPar(TYPE_I64, false), newOpPar(TYPE_I64, false)},
		[]*CXArgument{newOpPar(TYPE_I64, false)})
	AddOpCode(OP_I64_NEG, "i64.neg",
		[]*CXArgument{newOpPar(TYPE_I64, false)},
		[]*CXArgument{newOpPar(TYPE_I64, false)})
	AddOpCode(OP_I64_MUL, "i64.mul",
		[]*CXArgument{newOpPar(TYPE_I64, false), newOpPar(TYPE_I64, false)},
		[]*CXArgument{newOpPar(TYPE_I64, false)})
	AddOpCode(OP_I64_DIV, "i64.div",
		[]*CXArgument{newOpPar(TYPE_I64, false), newOpPar(TYPE_I64, false)},
		[]*CXArgument{newOpPar(TYPE_I64, false)})
	AddOpCode(OP_I64_ABS, "i64.abs",
		[]*CXArgument{newOpPar(TYPE_I64, false)},
		[]*CXArgument{newOpPar(TYPE_I64, false)})
	AddOpCode(OP_I64_POW, "i64.pow",
		[]*CXArgument{newOpPar(TYPE_I64, false), newOpPar(TYPE_I64, false)},
		[]*CXArgument{newOpPar(TYPE_I64, false)})
	AddOpCode(OP_I64_GT, "i64.gt",
		[]*CXArgument{newOpPar(TYPE_I64, false), newOpPar(TYPE_I64, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_I64_GTEQ, "i64.gteq",
		[]*CXArgument{newOpPar(TYPE_I64, false), newOpPar(TYPE_I64, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_I64_LT, "i64.lt",
		[]*CXArgument{newOpPar(TYPE_I64, false), newOpPar(TYPE_I64, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_I64_LTEQ, "i64.lteq",
		[]*CXArgument{newOpPar(TYPE_I64, false), newOpPar(TYPE_I64, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_I64_EQ, "i64.eq",
		[]*CXArgument{newOpPar(TYPE_I64, false), newOpPar(TYPE_I64, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_I64_UNEQ, "i64.uneq",
		[]*CXArgument{newOpPar(TYPE_I64, false), newOpPar(TYPE_I64, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_I64_MOD, "i64.mod",
		[]*CXArgument{newOpPar(TYPE_I64, false), newOpPar(TYPE_I64, false)},
		[]*CXArgument{newOpPar(TYPE_I64, false)})
	AddOpCode(OP_I64_RAND, "i64.rand",
		[]*CXArgument{newOpPar(TYPE_I64, false), newOpPar(TYPE_I64, false)},
		[]*CXArgument{newOpPar(TYPE_I64, false)})
	AddOpCode(OP_I64_BITAND, "i64.bitand",
		[]*CXArgument{newOpPar(TYPE_I64, false), newOpPar(TYPE_I64, false)},
		[]*CXArgument{newOpPar(TYPE_I64, false)})
	AddOpCode(OP_I64_BITOR, "i64.bitor",
		[]*CXArgument{newOpPar(TYPE_I64, false), newOpPar(TYPE_I64, false)},
		[]*CXArgument{newOpPar(TYPE_I64, false)})
	AddOpCode(OP_I64_BITXOR, "i64.bitxor",
		[]*CXArgument{newOpPar(TYPE_I64, false), newOpPar(TYPE_I64, false)},
		[]*CXArgument{newOpPar(TYPE_I64, false)})
	AddOpCode(OP_I64_BITCLEAR, "i64.bitclear",
		[]*CXArgument{newOpPar(TYPE_I64, false), newOpPar(TYPE_I64, false)},
		[]*CXArgument{newOpPar(TYPE_I64, false)})
	AddOpCode(OP_I64_BITSHL, "i64.bitshl",
		[]*CXArgument{newOpPar(TYPE_I64, false), newOpPar(TYPE_I64, false)},
		[]*CXArgument{newOpPar(TYPE_I64, false)})
	AddOpCode(OP_I64_BITSHR, "i64.bitshr",
		[]*CXArgument{newOpPar(TYPE_I64, false), newOpPar(TYPE_I64, false)},
		[]*CXArgument{newOpPar(TYPE_I64, false)})
	AddOpCode(OP_I64_SQRT, "i64.sqrt",
		[]*CXArgument{newOpPar(TYPE_I64, false)},
		[]*CXArgument{newOpPar(TYPE_I64, false)})
	AddOpCode(OP_I64_LOG, "i64.log",
		[]*CXArgument{newOpPar(TYPE_I64, false)},
		[]*CXArgument{newOpPar(TYPE_I64, false)})
	AddOpCode(OP_I64_LOG2, "i64.log2",
		[]*CXArgument{newOpPar(TYPE_I64, false)},
		[]*CXArgument{newOpPar(TYPE_I64, false)})
	AddOpCode(OP_I64_LOG10, "i64.log10",
		[]*CXArgument{newOpPar(TYPE_I64, false)},
		[]*CXArgument{newOpPar(TYPE_I64, false)})
	AddOpCode(OP_I64_MAX, "i64.max",
		[]*CXArgument{newOpPar(TYPE_I64, false), newOpPar(TYPE_I64, false)},
		[]*CXArgument{newOpPar(TYPE_I64, false)})
	AddOpCode(OP_I64_MIN, "i64.min",
		[]*CXArgument{newOpPar(TYPE_I64, false), newOpPar(TYPE_I64, false)},
		[]*CXArgument{newOpPar(TYPE_I64, false)})

	AddOpCode(OP_BYTE_STR, "byte.str",
		[]*CXArgument{newOpPar(TYPE_BYTE, false)},
		[]*CXArgument{newOpPar(TYPE_STR, false)})
	AddOpCode(OP_BYTE_I32, "byte.i32",
		[]*CXArgument{newOpPar(TYPE_BYTE, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false)})
	AddOpCode(OP_BYTE_I64, "byte.i64",
		[]*CXArgument{newOpPar(TYPE_BYTE, false)},
		[]*CXArgument{newOpPar(TYPE_I64, false)})
	AddOpCode(OP_BYTE_BYTE, "byte.byte",
		[]*CXArgument{newOpPar(TYPE_BYTE, false)},
		[]*CXArgument{newOpPar(TYPE_BYTE, false)})
	AddOpCode(OP_BYTE_UI8, "byte.ui8",
		[]*CXArgument{newOpPar(TYPE_BYTE, false)},
		[]*CXArgument{newOpPar(TYPE_UI8, false)})
	AddOpCode(OP_BYTE_UI16, "byte.ui16",
		[]*CXArgument{newOpPar(TYPE_BYTE, false)},
		[]*CXArgument{newOpPar(TYPE_UI16, false)})
	AddOpCode(OP_BYTE_UI32, "byte.ui32",
		[]*CXArgument{newOpPar(TYPE_BYTE, false)},
		[]*CXArgument{newOpPar(TYPE_UI32, false)})
	AddOpCode(OP_BYTE_UI64, "byte.ui64",
		[]*CXArgument{newOpPar(TYPE_BYTE, false)},
		[]*CXArgument{newOpPar(TYPE_UI64, false)})
	AddOpCode(OP_BYTE_F32, "byte.f32",
		[]*CXArgument{newOpPar(TYPE_BYTE, false)},
		[]*CXArgument{newOpPar(TYPE_F32, false)})
	AddOpCode(OP_BYTE_F64, "byte.f64",
		[]*CXArgument{newOpPar(TYPE_BYTE, false)},
		[]*CXArgument{newOpPar(TYPE_F64, false)})
	AddOpCode(OP_BYTE_PRINT, "byte.print",
		[]*CXArgument{newOpPar(TYPE_BYTE, false)},
		[]*CXArgument{})

	AddOpCode(OP_UI8_STR, "ui8.str",
		[]*CXArgument{newOpPar(TYPE_UI8, false)},
		[]*CXArgument{newOpPar(TYPE_STR, false)})
	AddOpCode(OP_UI8_I32, "ui8.i32",
		[]*CXArgument{newOpPar(TYPE_UI8, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false)})
	AddOpCode(OP_UI8_I64, "ui8.i64",
		[]*CXArgument{newOpPar(TYPE_UI8, false)},
		[]*CXArgument{newOpPar(TYPE_I64, false)})
	AddOpCode(OP_UI8_UI16, "ui8.ui16",
		[]*CXArgument{newOpPar(TYPE_UI8, false)},
		[]*CXArgument{newOpPar(TYPE_UI16, false)})
	AddOpCode(OP_UI8_UI32, "ui8.ui32",
		[]*CXArgument{newOpPar(TYPE_UI8, false)},
		[]*CXArgument{newOpPar(TYPE_UI32, false)})
	AddOpCode(OP_UI8_UI64, "ui8.ui64",
		[]*CXArgument{newOpPar(TYPE_UI8, false)},
		[]*CXArgument{newOpPar(TYPE_UI64, false)})
	AddOpCode(OP_UI8_F32, "ui8.f32",
		[]*CXArgument{newOpPar(TYPE_UI8, false)},
		[]*CXArgument{newOpPar(TYPE_F32, false)})
	AddOpCode(OP_UI8_F64, "ui8.f64",
		[]*CXArgument{newOpPar(TYPE_UI8, false)},
		[]*CXArgument{newOpPar(TYPE_F64, false)})
	AddOpCode(OP_UI8_ADD, "ui8.add",
		[]*CXArgument{newOpPar(TYPE_UI8, false), newOpPar(TYPE_UI8, false)},
		[]*CXArgument{newOpPar(TYPE_UI8, false)})
	AddOpCode(OP_UI8_SUB, "ui8.sub",
		[]*CXArgument{newOpPar(TYPE_UI8, false), newOpPar(TYPE_UI8, false)},
		[]*CXArgument{newOpPar(TYPE_UI8, false)})
	AddOpCode(OP_UI8_MUL, "ui8.mul",
		[]*CXArgument{newOpPar(TYPE_UI8, false), newOpPar(TYPE_UI8, false)},
		[]*CXArgument{newOpPar(TYPE_UI8, false)})
	AddOpCode(OP_UI8_DIV, "ui8.div",
		[]*CXArgument{newOpPar(TYPE_UI8, false), newOpPar(TYPE_UI8, false)},
		[]*CXArgument{newOpPar(TYPE_UI8, false)})
	AddOpCode(OP_UI8_GT, "ui8.gt",
		[]*CXArgument{newOpPar(TYPE_UI8, false), newOpPar(TYPE_UI8, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_UI8_GTEQ, "ui8.gteq",
		[]*CXArgument{newOpPar(TYPE_UI8, false), newOpPar(TYPE_UI8, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_UI8_LT, "ui8.lt",
		[]*CXArgument{newOpPar(TYPE_UI8, false), newOpPar(TYPE_UI8, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_UI8_LTEQ, "ui8.lteq",
		[]*CXArgument{newOpPar(TYPE_UI8, false), newOpPar(TYPE_UI8, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_UI8_EQ, "ui8.eq",
		[]*CXArgument{newOpPar(TYPE_UI8, false), newOpPar(TYPE_UI8, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_UI8_UNEQ, "ui8.uneq",
		[]*CXArgument{newOpPar(TYPE_UI8, false), newOpPar(TYPE_UI8, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_UI8_MOD, "ui8.mod",
		[]*CXArgument{newOpPar(TYPE_UI8, false), newOpPar(TYPE_UI8, false)},
		[]*CXArgument{newOpPar(TYPE_UI8, false)})
	AddOpCode(OP_UI8_RAND, "ui8.rand",
		[]*CXArgument{},
		[]*CXArgument{newOpPar(TYPE_UI8, false)})
	AddOpCode(OP_UI8_BITAND, "ui8.bitand",
		[]*CXArgument{newOpPar(TYPE_UI8, false), newOpPar(TYPE_UI8, false)},
		[]*CXArgument{newOpPar(TYPE_UI8, false)})
	AddOpCode(OP_UI8_BITOR, "ui8.bitor",
		[]*CXArgument{newOpPar(TYPE_UI8, false), newOpPar(TYPE_UI8, false)},
		[]*CXArgument{newOpPar(TYPE_UI8, false)})
	AddOpCode(OP_UI8_BITXOR, "ui8.bitxor",
		[]*CXArgument{newOpPar(TYPE_UI8, false), newOpPar(TYPE_UI8, false)},
		[]*CXArgument{newOpPar(TYPE_UI8, false)})
	AddOpCode(OP_UI8_BITCLEAR, "ui8.bitclear",
		[]*CXArgument{newOpPar(TYPE_UI8, false), newOpPar(TYPE_UI8, false)},
		[]*CXArgument{newOpPar(TYPE_UI8, false)})
	AddOpCode(OP_UI8_BITSHL, "ui8.bitshl",
		[]*CXArgument{newOpPar(TYPE_UI8, false), newOpPar(TYPE_UI8, false)},
		[]*CXArgument{newOpPar(TYPE_UI8, false)})
	AddOpCode(OP_UI8_BITSHR, "ui8.bitshr",
		[]*CXArgument{newOpPar(TYPE_UI8, false), newOpPar(TYPE_UI8, false)},
		[]*CXArgument{newOpPar(TYPE_UI8, false)})
	AddOpCode(OP_UI8_MAX, "ui8.max",
		[]*CXArgument{newOpPar(TYPE_UI8, false), newOpPar(TYPE_UI8, false)},
		[]*CXArgument{newOpPar(TYPE_UI8, false)})
	AddOpCode(OP_UI8_MIN, "ui8.min",
		[]*CXArgument{newOpPar(TYPE_UI8, false), newOpPar(TYPE_UI8, false)},
		[]*CXArgument{newOpPar(TYPE_UI8, false)})

	AddOpCode(OP_UI16_STR, "ui16.str",
		[]*CXArgument{newOpPar(TYPE_UI16, false)},
		[]*CXArgument{newOpPar(TYPE_STR, false)})
	AddOpCode(OP_UI16_I32, "ui16.i32",
		[]*CXArgument{newOpPar(TYPE_UI16, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false)})
	AddOpCode(OP_UI16_I64, "ui16.i64",
		[]*CXArgument{newOpPar(TYPE_UI16, false)},
		[]*CXArgument{newOpPar(TYPE_I64, false)})
	AddOpCode(OP_UI16_UI8, "ui16.ui8",
		[]*CXArgument{newOpPar(TYPE_UI16, false)},
		[]*CXArgument{newOpPar(TYPE_UI8, false)})
	AddOpCode(OP_UI16_UI32, "ui16.ui32",
		[]*CXArgument{newOpPar(TYPE_UI16, false)},
		[]*CXArgument{newOpPar(TYPE_UI32, false)})
	AddOpCode(OP_UI16_UI64, "ui16.ui64",
		[]*CXArgument{newOpPar(TYPE_UI16, false)},
		[]*CXArgument{newOpPar(TYPE_UI64, false)})
	AddOpCode(OP_UI16_F32, "ui16.f32",
		[]*CXArgument{newOpPar(TYPE_UI16, false)},
		[]*CXArgument{newOpPar(TYPE_F32, false)})
	AddOpCode(OP_UI16_F64, "ui16.f64",
		[]*CXArgument{newOpPar(TYPE_UI16, false)},
		[]*CXArgument{newOpPar(TYPE_F64, false)})
	AddOpCode(OP_UI16_ADD, "ui16.add",
		[]*CXArgument{newOpPar(TYPE_UI16, false), newOpPar(TYPE_UI16, false)},
		[]*CXArgument{newOpPar(TYPE_UI16, false)})
	AddOpCode(OP_UI16_SUB, "ui16.sub",
		[]*CXArgument{newOpPar(TYPE_UI16, false), newOpPar(TYPE_UI16, false)},
		[]*CXArgument{newOpPar(TYPE_UI16, false)})
	AddOpCode(OP_UI16_MUL, "ui16.mul",
		[]*CXArgument{newOpPar(TYPE_UI16, false), newOpPar(TYPE_UI16, false)},
		[]*CXArgument{newOpPar(TYPE_UI16, false)})
	AddOpCode(OP_UI16_DIV, "ui16.div",
		[]*CXArgument{newOpPar(TYPE_UI16, false), newOpPar(TYPE_UI16, false)},
		[]*CXArgument{newOpPar(TYPE_UI16, false)})
	AddOpCode(OP_UI16_GT, "ui16.gt",
		[]*CXArgument{newOpPar(TYPE_UI16, false), newOpPar(TYPE_UI16, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_UI16_GTEQ, "ui16.gteq",
		[]*CXArgument{newOpPar(TYPE_UI16, false), newOpPar(TYPE_UI16, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_UI16_LT, "ui16.lt",
		[]*CXArgument{newOpPar(TYPE_UI16, false), newOpPar(TYPE_UI16, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_UI16_LTEQ, "ui16.lteq",
		[]*CXArgument{newOpPar(TYPE_UI16, false), newOpPar(TYPE_UI16, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_UI16_EQ, "ui16.eq",
		[]*CXArgument{newOpPar(TYPE_UI16, false), newOpPar(TYPE_UI16, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_UI16_UNEQ, "ui16.uneq",
		[]*CXArgument{newOpPar(TYPE_UI16, false), newOpPar(TYPE_UI16, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_UI16_MOD, "ui16.mod",
		[]*CXArgument{newOpPar(TYPE_UI16, false), newOpPar(TYPE_UI16, false)},
		[]*CXArgument{newOpPar(TYPE_UI16, false)})
	AddOpCode(OP_UI16_RAND, "ui16.rand",
		[]*CXArgument{},
		[]*CXArgument{newOpPar(TYPE_UI16, false)})
	AddOpCode(OP_UI16_BITAND, "ui16.bitand",
		[]*CXArgument{newOpPar(TYPE_UI16, false), newOpPar(TYPE_UI16, false)},
		[]*CXArgument{newOpPar(TYPE_UI16, false)})
	AddOpCode(OP_UI16_BITOR, "ui16.bitor",
		[]*CXArgument{newOpPar(TYPE_UI16, false), newOpPar(TYPE_UI16, false)},
		[]*CXArgument{newOpPar(TYPE_UI16, false)})
	AddOpCode(OP_UI16_BITXOR, "ui16.bitxor",
		[]*CXArgument{newOpPar(TYPE_UI16, false), newOpPar(TYPE_UI16, false)},
		[]*CXArgument{newOpPar(TYPE_UI16, false)})
	AddOpCode(OP_UI16_BITCLEAR, "ui16.bitclear",
		[]*CXArgument{newOpPar(TYPE_UI16, false), newOpPar(TYPE_UI16, false)},
		[]*CXArgument{newOpPar(TYPE_UI16, false)})
	AddOpCode(OP_UI16_BITSHL, "ui16.bitshl",
		[]*CXArgument{newOpPar(TYPE_UI16, false), newOpPar(TYPE_UI16, false)},
		[]*CXArgument{newOpPar(TYPE_UI16, false)})
	AddOpCode(OP_UI16_BITSHR, "ui16.bitshr",
		[]*CXArgument{newOpPar(TYPE_UI16, false), newOpPar(TYPE_UI16, false)},
		[]*CXArgument{newOpPar(TYPE_UI16, false)})
	AddOpCode(OP_UI16_MAX, "ui16.max",
		[]*CXArgument{newOpPar(TYPE_UI16, false), newOpPar(TYPE_UI16, false)},
		[]*CXArgument{newOpPar(TYPE_UI16, false)})
	AddOpCode(OP_UI16_MIN, "ui16.min",
		[]*CXArgument{newOpPar(TYPE_UI16, false), newOpPar(TYPE_UI16, false)},
		[]*CXArgument{newOpPar(TYPE_UI16, false)})

	AddOpCode(OP_UI32_STR, "ui32.str",
		[]*CXArgument{newOpPar(TYPE_UI32, false)},
		[]*CXArgument{newOpPar(TYPE_STR, false)})
	AddOpCode(OP_UI32_I32, "ui32.i32",
		[]*CXArgument{newOpPar(TYPE_UI32, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false)})
	AddOpCode(OP_UI32_I64, "ui32.i64",
		[]*CXArgument{newOpPar(TYPE_UI32, false)},
		[]*CXArgument{newOpPar(TYPE_I64, false)})
	AddOpCode(OP_UI32_UI8, "ui32.ui8",
		[]*CXArgument{newOpPar(TYPE_UI32, false)},
		[]*CXArgument{newOpPar(TYPE_UI8, false)})
	AddOpCode(OP_UI32_UI16, "ui32.ui16",
		[]*CXArgument{newOpPar(TYPE_UI32, false)},
		[]*CXArgument{newOpPar(TYPE_UI16, false)})
	AddOpCode(OP_UI32_UI64, "ui32.ui64",
		[]*CXArgument{newOpPar(TYPE_UI32, false)},
		[]*CXArgument{newOpPar(TYPE_UI64, false)})
	AddOpCode(OP_UI32_F32, "ui32.f32",
		[]*CXArgument{newOpPar(TYPE_UI32, false)},
		[]*CXArgument{newOpPar(TYPE_F32, false)})
	AddOpCode(OP_UI32_F64, "ui32.f64",
		[]*CXArgument{newOpPar(TYPE_UI32, false)},
		[]*CXArgument{newOpPar(TYPE_F64, false)})
	AddOpCode(OP_UI32_ADD, "ui32.add",
		[]*CXArgument{newOpPar(TYPE_UI32, false), newOpPar(TYPE_UI32, false)},
		[]*CXArgument{newOpPar(TYPE_UI32, false)})
	AddOpCode(OP_UI32_SUB, "ui32.sub",
		[]*CXArgument{newOpPar(TYPE_UI32, false), newOpPar(TYPE_UI32, false)},
		[]*CXArgument{newOpPar(TYPE_UI32, false)})
	AddOpCode(OP_UI32_MUL, "ui32.mul",
		[]*CXArgument{newOpPar(TYPE_UI32, false), newOpPar(TYPE_UI32, false)},
		[]*CXArgument{newOpPar(TYPE_UI32, false)})
	AddOpCode(OP_UI32_DIV, "ui32.div",
		[]*CXArgument{newOpPar(TYPE_UI32, false), newOpPar(TYPE_UI32, false)},
		[]*CXArgument{newOpPar(TYPE_UI32, false)})
	AddOpCode(OP_UI32_GT, "ui32.gt",
		[]*CXArgument{newOpPar(TYPE_UI32, false), newOpPar(TYPE_UI32, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_UI32_GTEQ, "ui32.gteq",
		[]*CXArgument{newOpPar(TYPE_UI32, false), newOpPar(TYPE_UI32, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_UI32_LT, "ui32.lt",
		[]*CXArgument{newOpPar(TYPE_UI32, false), newOpPar(TYPE_UI32, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_UI32_LTEQ, "ui32.lteq",
		[]*CXArgument{newOpPar(TYPE_UI32, false), newOpPar(TYPE_UI32, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_UI32_EQ, "ui32.eq",
		[]*CXArgument{newOpPar(TYPE_UI32, false), newOpPar(TYPE_UI32, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_UI32_UNEQ, "ui32.uneq",
		[]*CXArgument{newOpPar(TYPE_UI32, false), newOpPar(TYPE_UI32, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_UI32_MOD, "ui32.mod",
		[]*CXArgument{newOpPar(TYPE_UI32, false), newOpPar(TYPE_UI32, false)},
		[]*CXArgument{newOpPar(TYPE_UI32, false)})
	AddOpCode(OP_UI32_RAND, "ui32.rand",
		[]*CXArgument{},
		[]*CXArgument{newOpPar(TYPE_UI32, false)})
	AddOpCode(OP_UI32_BITAND, "ui32.bitand",
		[]*CXArgument{newOpPar(TYPE_UI32, false), newOpPar(TYPE_UI32, false)},
		[]*CXArgument{newOpPar(TYPE_UI32, false)})
	AddOpCode(OP_UI32_BITOR, "ui32.bitor",
		[]*CXArgument{newOpPar(TYPE_UI32, false), newOpPar(TYPE_UI32, false)},
		[]*CXArgument{newOpPar(TYPE_UI32, false)})
	AddOpCode(OP_UI32_BITXOR, "ui32.bitxor",
		[]*CXArgument{newOpPar(TYPE_UI32, false), newOpPar(TYPE_UI32, false)},
		[]*CXArgument{newOpPar(TYPE_UI32, false)})
	AddOpCode(OP_UI32_BITCLEAR, "ui32.bitclear",
		[]*CXArgument{newOpPar(TYPE_UI32, false), newOpPar(TYPE_UI32, false)},
		[]*CXArgument{newOpPar(TYPE_UI32, false)})
	AddOpCode(OP_UI32_BITSHL, "ui32.bitshl",
		[]*CXArgument{newOpPar(TYPE_UI32, false), newOpPar(TYPE_UI32, false)},
		[]*CXArgument{newOpPar(TYPE_UI32, false)})
	AddOpCode(OP_UI32_BITSHR, "ui32.bitshr",
		[]*CXArgument{newOpPar(TYPE_UI32, false), newOpPar(TYPE_UI32, false)},
		[]*CXArgument{newOpPar(TYPE_UI32, false)})
	AddOpCode(OP_UI32_MAX, "ui32.max",
		[]*CXArgument{newOpPar(TYPE_UI32, false), newOpPar(TYPE_UI32, false)},
		[]*CXArgument{newOpPar(TYPE_UI32, false)})
	AddOpCode(OP_UI32_MIN, "ui32.min",
		[]*CXArgument{newOpPar(TYPE_UI32, false), newOpPar(TYPE_UI32, false)},
		[]*CXArgument{newOpPar(TYPE_UI32, false)})

	AddOpCode(OP_UI64_STR, "ui64.str",
		[]*CXArgument{newOpPar(TYPE_UI64, false)},
		[]*CXArgument{newOpPar(TYPE_STR, false)})
	AddOpCode(OP_UI64_I32, "ui64.i32",
		[]*CXArgument{newOpPar(TYPE_UI64, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false)})
	AddOpCode(OP_UI64_I64, "ui64.i64",
		[]*CXArgument{newOpPar(TYPE_UI64, false)},
		[]*CXArgument{newOpPar(TYPE_I64, false)})
	AddOpCode(OP_UI64_UI8, "ui64.ui8",
		[]*CXArgument{newOpPar(TYPE_UI64, false)},
		[]*CXArgument{newOpPar(TYPE_UI8, false)})
	AddOpCode(OP_UI64_UI16, "ui64.ui16",
		[]*CXArgument{newOpPar(TYPE_UI64, false)},
		[]*CXArgument{newOpPar(TYPE_UI16, false)})
	AddOpCode(OP_UI64_UI32, "ui64.ui32",
		[]*CXArgument{newOpPar(TYPE_UI64, false)},
		[]*CXArgument{newOpPar(TYPE_UI32, false)})
	AddOpCode(OP_UI64_F32, "ui64.f32",
		[]*CXArgument{newOpPar(TYPE_UI64, false)},
		[]*CXArgument{newOpPar(TYPE_F32, false)})
	AddOpCode(OP_UI64_F64, "ui64.f64",
		[]*CXArgument{newOpPar(TYPE_UI64, false)},
		[]*CXArgument{newOpPar(TYPE_F64, false)})
	AddOpCode(OP_UI64_ADD, "ui64.add",
		[]*CXArgument{newOpPar(TYPE_UI64, false), newOpPar(TYPE_UI64, false)},
		[]*CXArgument{newOpPar(TYPE_UI64, false)})
	AddOpCode(OP_UI64_SUB, "ui64.sub",
		[]*CXArgument{newOpPar(TYPE_UI64, false), newOpPar(TYPE_UI64, false)},
		[]*CXArgument{newOpPar(TYPE_UI64, false)})
	AddOpCode(OP_UI64_MUL, "ui64.mul",
		[]*CXArgument{newOpPar(TYPE_UI64, false), newOpPar(TYPE_UI64, false)},
		[]*CXArgument{newOpPar(TYPE_UI64, false)})
	AddOpCode(OP_UI64_DIV, "ui64.div",
		[]*CXArgument{newOpPar(TYPE_UI64, false), newOpPar(TYPE_UI64, false)},
		[]*CXArgument{newOpPar(TYPE_UI64, false)})
	AddOpCode(OP_UI64_GT, "ui64.gt",
		[]*CXArgument{newOpPar(TYPE_UI64, false), newOpPar(TYPE_UI64, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_UI64_GTEQ, "ui64.gteq",
		[]*CXArgument{newOpPar(TYPE_UI64, false), newOpPar(TYPE_UI64, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_UI64_LT, "ui64.lt",
		[]*CXArgument{newOpPar(TYPE_UI64, false), newOpPar(TYPE_UI64, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_UI64_LTEQ, "ui64.lteq",
		[]*CXArgument{newOpPar(TYPE_UI64, false), newOpPar(TYPE_UI64, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_UI64_EQ, "ui64.eq",
		[]*CXArgument{newOpPar(TYPE_UI64, false), newOpPar(TYPE_UI64, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_UI64_UNEQ, "ui64.uneq",
		[]*CXArgument{newOpPar(TYPE_UI64, false), newOpPar(TYPE_UI64, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_UI64_MOD, "ui64.mod",
		[]*CXArgument{newOpPar(TYPE_UI64, false), newOpPar(TYPE_UI64, false)},
		[]*CXArgument{newOpPar(TYPE_UI64, false)})
	AddOpCode(OP_UI64_RAND, "ui64.rand",
		[]*CXArgument{},
		[]*CXArgument{newOpPar(TYPE_UI64, false)})
	AddOpCode(OP_UI64_BITAND, "ui64.bitand",
		[]*CXArgument{newOpPar(TYPE_UI64, false), newOpPar(TYPE_UI64, false)},
		[]*CXArgument{newOpPar(TYPE_UI64, false)})
	AddOpCode(OP_UI64_BITOR, "ui64.bitor",
		[]*CXArgument{newOpPar(TYPE_UI64, false), newOpPar(TYPE_UI64, false)},
		[]*CXArgument{newOpPar(TYPE_UI64, false)})
	AddOpCode(OP_UI64_BITXOR, "ui64.bitxor",
		[]*CXArgument{newOpPar(TYPE_UI64, false), newOpPar(TYPE_UI64, false)},
		[]*CXArgument{newOpPar(TYPE_UI64, false)})
	AddOpCode(OP_UI64_BITCLEAR, "ui64.bitclear",
		[]*CXArgument{newOpPar(TYPE_UI64, false), newOpPar(TYPE_UI64, false)},
		[]*CXArgument{newOpPar(TYPE_UI64, false)})
	AddOpCode(OP_UI64_BITSHL, "ui64.bitshl",
		[]*CXArgument{newOpPar(TYPE_UI64, false), newOpPar(TYPE_UI64, false)},
		[]*CXArgument{newOpPar(TYPE_UI64, false)})
	AddOpCode(OP_UI64_BITSHR, "ui64.bitshr",
		[]*CXArgument{newOpPar(TYPE_UI64, false), newOpPar(TYPE_UI64, false)},
		[]*CXArgument{newOpPar(TYPE_UI64, false)})
	AddOpCode(OP_UI64_MAX, "ui64.max",
		[]*CXArgument{newOpPar(TYPE_UI64, false), newOpPar(TYPE_UI64, false)},
		[]*CXArgument{newOpPar(TYPE_UI64, false)})
	AddOpCode(OP_UI64_MIN, "ui64.min",
		[]*CXArgument{newOpPar(TYPE_UI64, false), newOpPar(TYPE_UI64, false)},
		[]*CXArgument{newOpPar(TYPE_UI64, false)})

	AddOpCode(OP_F32_STR, "f32.str",
		[]*CXArgument{newOpPar(TYPE_F32, false)},
		[]*CXArgument{newOpPar(TYPE_STR, false)})
	AddOpCode(OP_F32_I32, "f32.i32",
		[]*CXArgument{newOpPar(TYPE_F32, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false)})
	AddOpCode(OP_F32_I64, "f32.i64",
		[]*CXArgument{newOpPar(TYPE_F32, false)},
		[]*CXArgument{newOpPar(TYPE_I64, false)})
	AddOpCode(OP_F32_BYTE, "f32.byte",
		[]*CXArgument{newOpPar(TYPE_F32, false)},
		[]*CXArgument{newOpPar(TYPE_BYTE, false)})
	AddOpCode(OP_F32_UI8, "f32.ui8",
		[]*CXArgument{newOpPar(TYPE_F32, false)},
		[]*CXArgument{newOpPar(TYPE_UI8, false)})
	AddOpCode(OP_F32_UI16, "f32.ui16",
		[]*CXArgument{newOpPar(TYPE_F32, false)},
		[]*CXArgument{newOpPar(TYPE_UI16, false)})
	AddOpCode(OP_F32_UI32, "f32.ui32",
		[]*CXArgument{newOpPar(TYPE_F32, false)},
		[]*CXArgument{newOpPar(TYPE_UI32, false)})
	AddOpCode(OP_F32_UI64, "f32.ui64",
		[]*CXArgument{newOpPar(TYPE_F32, false)},
		[]*CXArgument{newOpPar(TYPE_UI64, false)})
	AddOpCode(OP_F32_F32, "f32.f32",
		[]*CXArgument{newOpPar(TYPE_F32, false)},
		[]*CXArgument{newOpPar(TYPE_F32, false)})
	AddOpCode(OP_F32_F64, "f32.f64",
		[]*CXArgument{newOpPar(TYPE_F32, false)},
		[]*CXArgument{newOpPar(TYPE_F64, false)})
	AddOpCode(OP_F32_PRINT, "f32.print",
		[]*CXArgument{newOpPar(TYPE_F32, false)},
		[]*CXArgument{})
	AddOpCode(OP_F32_ADD, "f32.add",
		[]*CXArgument{newOpPar(TYPE_F32, false), newOpPar(TYPE_F32, false)},
		[]*CXArgument{newOpPar(TYPE_F32, false)})
	AddOpCode(OP_F32_SUB, "f32.sub",
		[]*CXArgument{newOpPar(TYPE_F32, false), newOpPar(TYPE_F32, false)},
		[]*CXArgument{newOpPar(TYPE_F32, false)})
	AddOpCode(OP_F32_NEG, "f32.neg",
		[]*CXArgument{newOpPar(TYPE_F32, false)},
		[]*CXArgument{newOpPar(TYPE_F32, false)})
	AddOpCode(OP_F32_MUL, "f32.mul",
		[]*CXArgument{newOpPar(TYPE_F32, false), newOpPar(TYPE_F32, false)},
		[]*CXArgument{newOpPar(TYPE_F32, false)})
	AddOpCode(OP_F32_DIV, "f32.div",
		[]*CXArgument{newOpPar(TYPE_F32, false), newOpPar(TYPE_F32, false)},
		[]*CXArgument{newOpPar(TYPE_F32, false)})
	AddOpCode(OP_F32_ABS, "f32.abs",
		[]*CXArgument{newOpPar(TYPE_F32, false)},
		[]*CXArgument{newOpPar(TYPE_F32, false)})
	AddOpCode(OP_F32_POW, "f32.pow",
		[]*CXArgument{newOpPar(TYPE_F32, false), newOpPar(TYPE_F32, false)},
		[]*CXArgument{newOpPar(TYPE_F32, false)})
	AddOpCode(OP_F32_GT, "f32.gt",
		[]*CXArgument{newOpPar(TYPE_F32, false), newOpPar(TYPE_F32, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_F32_GTEQ, "f32.gteq",
		[]*CXArgument{newOpPar(TYPE_F32, false), newOpPar(TYPE_F32, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_F32_LT, "f32.lt",
		[]*CXArgument{newOpPar(TYPE_F32, false), newOpPar(TYPE_F32, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_F32_LTEQ, "f32.lteq",
		[]*CXArgument{newOpPar(TYPE_F32, false), newOpPar(TYPE_F32, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_F32_EQ, "f32.eq",
		[]*CXArgument{newOpPar(TYPE_F32, false), newOpPar(TYPE_F32, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_F32_UNEQ, "f32.uneq",
		[]*CXArgument{newOpPar(TYPE_F32, false), newOpPar(TYPE_F32, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_F32_COS, "f32.cos",
		[]*CXArgument{newOpPar(TYPE_F32, false)},
		[]*CXArgument{newOpPar(TYPE_F32, false)})
	AddOpCode(OP_F32_SIN, "f32.sin",
		[]*CXArgument{newOpPar(TYPE_F32, false)},
		[]*CXArgument{newOpPar(TYPE_F32, false)})
	AddOpCode(OP_F32_SQRT, "f32.sqrt",
		[]*CXArgument{newOpPar(TYPE_F32, false)},
		[]*CXArgument{newOpPar(TYPE_F32, false)})
	AddOpCode(OP_F32_LOG, "f32.log",
		[]*CXArgument{newOpPar(TYPE_F32, false)},
		[]*CXArgument{newOpPar(TYPE_F32, false)})
	AddOpCode(OP_F32_LOG2, "f32.log2",
		[]*CXArgument{newOpPar(TYPE_F32, false)},
		[]*CXArgument{newOpPar(TYPE_F32, false)})
	AddOpCode(OP_F32_LOG10, "f32.log10",
		[]*CXArgument{newOpPar(TYPE_F32, false)},
		[]*CXArgument{newOpPar(TYPE_F32, false)})
	AddOpCode(OP_F32_MAX, "f32.max",
		[]*CXArgument{newOpPar(TYPE_F32, false), newOpPar(TYPE_F32, false)},
		[]*CXArgument{newOpPar(TYPE_F32, false)})
	AddOpCode(OP_F32_MIN, "f32.min",
		[]*CXArgument{newOpPar(TYPE_F32, false), newOpPar(TYPE_F32, false)},
		[]*CXArgument{newOpPar(TYPE_F32, false)})
	AddOpCode(OP_F32_IS_NAN, "f32.isnan",
		[]*CXArgument{newOpPar(TYPE_F32, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})

	AddOpCode(OP_F64_STR, "f64.str",
		[]*CXArgument{newOpPar(TYPE_F64, false)},
		[]*CXArgument{newOpPar(TYPE_STR, false)})
	AddOpCode(OP_F64_I32, "f64.i32",
		[]*CXArgument{newOpPar(TYPE_F64, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false)})
	AddOpCode(OP_F64_I64, "f64.i64",
		[]*CXArgument{newOpPar(TYPE_F64, false)},
		[]*CXArgument{newOpPar(TYPE_I64, false)})
	AddOpCode(OP_F64_BYTE, "f64.byte",
		[]*CXArgument{newOpPar(TYPE_F64, false)},
		[]*CXArgument{newOpPar(TYPE_BYTE, false)})
	AddOpCode(OP_F64_UI8, "f64.ui8",
		[]*CXArgument{newOpPar(TYPE_F64, false)},
		[]*CXArgument{newOpPar(TYPE_UI8, false)})
	AddOpCode(OP_F64_UI16, "f64.ui16",
		[]*CXArgument{newOpPar(TYPE_F64, false)},
		[]*CXArgument{newOpPar(TYPE_UI16, false)})
	AddOpCode(OP_F64_UI32, "f64.ui32",
		[]*CXArgument{newOpPar(TYPE_F64, false)},
		[]*CXArgument{newOpPar(TYPE_UI32, false)})
	AddOpCode(OP_F64_UI64, "f64.ui64",
		[]*CXArgument{newOpPar(TYPE_F64, false)},
		[]*CXArgument{newOpPar(TYPE_UI64, false)})
	AddOpCode(OP_F64_F32, "f64.f32",
		[]*CXArgument{newOpPar(TYPE_F64, false)},
		[]*CXArgument{newOpPar(TYPE_F32, false)})
	AddOpCode(OP_F64_F64, "f64.f64",
		[]*CXArgument{newOpPar(TYPE_F64, false)},
		[]*CXArgument{newOpPar(TYPE_F64, false)})
	AddOpCode(OP_F64_PRINT, "f64.print",
		[]*CXArgument{newOpPar(TYPE_F64, false)},
		[]*CXArgument{})
	AddOpCode(OP_F64_ADD, "f64.add",
		[]*CXArgument{newOpPar(TYPE_F64, false), newOpPar(TYPE_F64, false)},
		[]*CXArgument{newOpPar(TYPE_F64, false)})
	AddOpCode(OP_F64_SUB, "f64.sub",
		[]*CXArgument{newOpPar(TYPE_F64, false), newOpPar(TYPE_F64, false)},
		[]*CXArgument{newOpPar(TYPE_F64, false)})
	AddOpCode(OP_F64_NEG, "f64.neg",
		[]*CXArgument{newOpPar(TYPE_F64, false)},
		[]*CXArgument{newOpPar(TYPE_F64, false)})
	AddOpCode(OP_F64_MUL, "f64.mul",
		[]*CXArgument{newOpPar(TYPE_F64, false), newOpPar(TYPE_F64, false)},
		[]*CXArgument{newOpPar(TYPE_F64, false)})
	AddOpCode(OP_F64_DIV, "f64.div",
		[]*CXArgument{newOpPar(TYPE_F64, false), newOpPar(TYPE_F64, false)},
		[]*CXArgument{newOpPar(TYPE_F64, false)})
	AddOpCode(OP_F64_ABS, "f64.abs",
		[]*CXArgument{newOpPar(TYPE_F64, false)},
		[]*CXArgument{newOpPar(TYPE_F64, false)})
	AddOpCode(OP_F64_POW, "f64.pow",
		[]*CXArgument{newOpPar(TYPE_F64, false), newOpPar(TYPE_F64, false)},
		[]*CXArgument{newOpPar(TYPE_F64, false)})
	AddOpCode(OP_F64_GT, "f64.gt",
		[]*CXArgument{newOpPar(TYPE_F64, false), newOpPar(TYPE_F64, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_F64_GTEQ, "f64.gteq",
		[]*CXArgument{newOpPar(TYPE_F64, false), newOpPar(TYPE_F64, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_F64_LT, "f64.lt",
		[]*CXArgument{newOpPar(TYPE_F64, false), newOpPar(TYPE_F64, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_F64_LTEQ, "f64.lteq",
		[]*CXArgument{newOpPar(TYPE_F64, false), newOpPar(TYPE_F64, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_F64_EQ, "f64.eq",
		[]*CXArgument{newOpPar(TYPE_F64, false), newOpPar(TYPE_F64, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_F64_UNEQ, "f64.uneq",
		[]*CXArgument{newOpPar(TYPE_F64, false), newOpPar(TYPE_F64, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_F64_COS, "f64.cos",
		[]*CXArgument{newOpPar(TYPE_F64, false)},
		[]*CXArgument{newOpPar(TYPE_F64, false)})
	AddOpCode(OP_F64_SIN, "f64.sin",
		[]*CXArgument{newOpPar(TYPE_F64, false)},
		[]*CXArgument{newOpPar(TYPE_F64, false)})
	AddOpCode(OP_F64_SQRT, "f64.sqrt",
		[]*CXArgument{newOpPar(TYPE_F64, false)},
		[]*CXArgument{newOpPar(TYPE_F64, false)})
	AddOpCode(OP_F64_LOG, "f64.log",
		[]*CXArgument{newOpPar(TYPE_F64, false)},
		[]*CXArgument{newOpPar(TYPE_F64, false)})
	AddOpCode(OP_F64_LOG2, "f64.log2",
		[]*CXArgument{newOpPar(TYPE_F64, false)},
		[]*CXArgument{newOpPar(TYPE_F64, false)})
	AddOpCode(OP_F64_LOG10, "f64.log10",
		[]*CXArgument{newOpPar(TYPE_F64, false)},
		[]*CXArgument{newOpPar(TYPE_F64, false)})
	AddOpCode(OP_F64_MAX, "f64.max",
		[]*CXArgument{newOpPar(TYPE_F64, false), newOpPar(TYPE_F64, false)},
		[]*CXArgument{newOpPar(TYPE_F64, false)})
	AddOpCode(OP_F64_MIN, "f64.min",
		[]*CXArgument{newOpPar(TYPE_F64, false), newOpPar(TYPE_F64, false)},
		[]*CXArgument{newOpPar(TYPE_F64, false)})

	AddOpCode(OP_STR_PRINT, "str.print",
		[]*CXArgument{newOpPar(TYPE_STR, false)},
		[]*CXArgument{})
	AddOpCode(OP_STR_CONCAT, "str.concat",
		[]*CXArgument{newOpPar(TYPE_STR, false), newOpPar(TYPE_STR, false)},
		[]*CXArgument{newOpPar(TYPE_STR, false)})
	AddOpCode(OP_STR_SUBSTR, "str.substr",
		[]*CXArgument{newOpPar(TYPE_STR, false), newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_STR, false)})
	AddOpCode(OP_STR_INDEX, "str.index",
		[]*CXArgument{newOpPar(TYPE_STR, false), newOpPar(TYPE_STR, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false)})
	AddOpCode(OP_STR_TRIM_SPACE, "str.trimspace",
		[]*CXArgument{newOpPar(TYPE_STR, false)},
		[]*CXArgument{newOpPar(TYPE_STR, false)})
	AddOpCode(OP_STR_EQ, "str.eq",
		[]*CXArgument{newOpPar(TYPE_STR, false), newOpPar(TYPE_STR, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})

	AddOpCode(OP_STR_STR, "str.str",
		[]*CXArgument{newOpPar(TYPE_STR, false)},
		[]*CXArgument{newOpPar(TYPE_STR, false)})
	AddOpCode(OP_STR_I32, "str.i32",
		[]*CXArgument{newOpPar(TYPE_STR, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false)})
	AddOpCode(OP_STR_I64, "str.i64",
		[]*CXArgument{newOpPar(TYPE_STR, false)},
		[]*CXArgument{newOpPar(TYPE_I64, false)})
	AddOpCode(OP_STR_BYTE, "str.byte",
		[]*CXArgument{newOpPar(TYPE_STR, false)},
		[]*CXArgument{newOpPar(TYPE_BYTE, false)})
	AddOpCode(OP_STR_UI8, "str.ui8",
		[]*CXArgument{newOpPar(TYPE_STR, false)},
		[]*CXArgument{newOpPar(TYPE_UI8, false)})
	AddOpCode(OP_STR_UI16, "str.ui16",
		[]*CXArgument{newOpPar(TYPE_STR, false)},
		[]*CXArgument{newOpPar(TYPE_UI16, false)})
	AddOpCode(OP_STR_UI32, "str.ui32",
		[]*CXArgument{newOpPar(TYPE_STR, false)},
		[]*CXArgument{newOpPar(TYPE_UI32, false)})
	AddOpCode(OP_STR_UI64, "str.ui64",
		[]*CXArgument{newOpPar(TYPE_STR, false)},
		[]*CXArgument{newOpPar(TYPE_UI64, false)})
	AddOpCode(OP_STR_F32, "str.f32",
		[]*CXArgument{newOpPar(TYPE_STR, false)},
		[]*CXArgument{newOpPar(TYPE_F32, false)})
	AddOpCode(OP_STR_F64, "str.f64",
		[]*CXArgument{newOpPar(TYPE_STR, false)},
		[]*CXArgument{newOpPar(TYPE_F64, false)})

	AddOpCode(OP_APPEND, "append",
		[]*CXArgument{newOpPar(TYPE_UNDEFINED, true), newOpPar(TYPE_UNDEFINED, true)},
		[]*CXArgument{newOpPar(TYPE_UNDEFINED, true)})
	AddOpCode(OP_RESIZE, "resize",
		[]*CXArgument{newOpPar(TYPE_UNDEFINED, true), newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_UNDEFINED, true)})
	AddOpCode(OP_INSERT, "insert",
		[]*CXArgument{newOpPar(TYPE_UNDEFINED, true), newOpPar(TYPE_UNDEFINED, true)},
		[]*CXArgument{newOpPar(TYPE_UNDEFINED, true)})
	AddOpCode(OP_REMOVE, "remove",
		[]*CXArgument{newOpPar(TYPE_UNDEFINED, true), newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_UNDEFINED, true)})
	AddOpCode(OP_COPY, "copy",
		[]*CXArgument{newOpPar(TYPE_UNDEFINED, true), newOpPar(TYPE_UNDEFINED, true)},
		[]*CXArgument{newOpPar(TYPE_I32, false)})

	AddOpCode(OP_ASSERT, "assert",
		[]*CXArgument{newOpPar(TYPE_UNDEFINED, false), newOpPar(TYPE_UNDEFINED, false), newOpPar(TYPE_STR, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_TEST, "test",
		[]*CXArgument{newOpPar(TYPE_UNDEFINED, false), newOpPar(TYPE_UNDEFINED, false), newOpPar(TYPE_STR, false)},
		[]*CXArgument{})
	AddOpCode(OP_PANIC, "panic",
		[]*CXArgument{newOpPar(TYPE_UNDEFINED, false), newOpPar(TYPE_UNDEFINED, false), newOpPar(TYPE_STR, false)},
		[]*CXArgument{})
	AddOpCode(OP_STRERROR, "strerror",
		[]*CXArgument{newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_STR, false)})
	// affordances
	AddOpCode(OP_AFF_PRINT, "aff.print",
		[]*CXArgument{newOpPar(TYPE_AFF, false)},
		[]*CXArgument{})
	AddOpCode(OP_AFF_QUERY, "aff.query",
		[]*CXArgument{newOpPar(TYPE_AFF, false)},
		[]*CXArgument{newOpPar(TYPE_AFF, false)})
	AddOpCode(OP_AFF_ON, "aff.on",
		[]*CXArgument{newOpPar(TYPE_AFF, false), newOpPar(TYPE_AFF, false)},
		[]*CXArgument{})
	AddOpCode(OP_AFF_OF, "aff.of",
		[]*CXArgument{newOpPar(TYPE_AFF, false), newOpPar(TYPE_AFF, false)},
		[]*CXArgument{})
	AddOpCode(OP_AFF_INFORM, "aff.inform",
		[]*CXArgument{newOpPar(TYPE_AFF, false), newOpPar(TYPE_I32, false), newOpPar(TYPE_AFF, false)},
		[]*CXArgument{})
	AddOpCode(OP_AFF_REQUEST, "aff.request",
		[]*CXArgument{newOpPar(TYPE_AFF, false), newOpPar(TYPE_I32, false), newOpPar(TYPE_AFF, false)},
		[]*CXArgument{})

	// exec
	execNativeCore = func(prgrm *CXProgram) {
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
			opSerialize(expr, fp)
		case OP_DESERIALIZE:
			opDeserialize(expr, fp)

		case OP_UND_EQUAL:
			opEqual(expr, fp)
		case OP_UND_UNEQUAL:
			opUnequal(expr, fp)
		case OP_UND_BITAND:
			opBitand(expr, fp)
		case OP_UND_BITXOR:
			opBitxor(expr, fp)
		case OP_UND_BITOR:
			opBitor(expr, fp)
		case OP_UND_BITCLEAR:
			opBitclear(expr, fp)
		case OP_UND_MUL:
			opMul(expr, fp)
		case OP_UND_DIV:
			opDiv(expr, fp)
		case OP_UND_MOD:
			opMod(expr, fp)
		case OP_UND_ADD:
			opAdd(expr, fp)
		case OP_UND_SUB, OP_UND_NEG:
			opSub(expr, fp)
		case OP_UND_BITSHL:
			opBitshl(expr, fp)
		case OP_UND_BITSHR:
			opBitshr(expr, fp)
		case OP_UND_LT:
			opLt(expr, fp)
		case OP_UND_GT:
			opGt(expr, fp)
		case OP_UND_LTEQ:
			opLteq(expr, fp)
		case OP_UND_GTEQ:
			opGteq(expr, fp)
		case OP_UND_LEN:
			opLen(expr, fp)
		case OP_UND_PRINTF:
			opPrintf(expr, fp)
		case OP_UND_SPRINTF:
			opSprintf(expr, fp)
		case OP_UND_READ:
			opRead(expr, fp)

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

		case OP_I32_STR:
			opI32Cast(expr, fp)
		case OP_I32_I32:
			opI32Cast(expr, fp)
		case OP_I32_I64:
			opI32Cast(expr, fp)
		case OP_I32_BYTE:
			opI32Cast(expr, fp)
		case OP_I32_UI8:
			opI32Cast(expr, fp)
		case OP_I32_UI16:
			opI32Cast(expr, fp)
		case OP_I32_UI32:
			opI32Cast(expr, fp)
		case OP_I32_UI64:
			opI32Cast(expr, fp)
		case OP_I32_F32:
			opI32Cast(expr, fp)
		case OP_I32_F64:
			opI32Cast(expr, fp)
		case OP_I32_PRINT:
			opI32Print(expr, fp)
		case OP_I32_ADD:
			opI32Add(expr, fp)
		case OP_I32_SUB, OP_I32_NEG:
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

		case OP_I64_STR:
			opI64Cast(expr, fp)
		case OP_I64_I32:
			opI64Cast(expr, fp)
		case OP_I64_I64:
			opI64Cast(expr, fp)
		case OP_I64_BYTE:
			opI64Cast(expr, fp)
		case OP_I64_UI8:
			opI64Cast(expr, fp)
		case OP_I64_UI16:
			opI64Cast(expr, fp)
		case OP_I64_UI32:
			opI64Cast(expr, fp)
		case OP_I64_UI64:
			opI64Cast(expr, fp)
		case OP_I64_F32:
			opI64Cast(expr, fp)
		case OP_I64_F64:
			opI64Cast(expr, fp)
		case OP_I64_PRINT:
			opI64Print(expr, fp)
		case OP_I64_ADD:
			opI64Add(expr, fp)
		case OP_I64_SUB, OP_I64_NEG:
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

		case OP_BYTE_STR:
			opByteByte(expr, fp)
		case OP_BYTE_I32:
			opByteByte(expr, fp)
		case OP_BYTE_I64:
			opByteByte(expr, fp)
		case OP_BYTE_BYTE:
			opByteByte(expr, fp)
		case OP_BYTE_UI8:
			opByteByte(expr, fp)
		case OP_BYTE_UI16:
			opByteByte(expr, fp)
		case OP_BYTE_UI32:
			opByteByte(expr, fp)
		case OP_BYTE_UI64:
			opByteByte(expr, fp)
		case OP_BYTE_F32:
			opByteByte(expr, fp)
		case OP_BYTE_F64:
			opByteByte(expr, fp)
		case OP_BYTE_PRINT:
			opBytePrint(expr, fp)

		case OP_UI8_STR:
			opUI8ToStr(expr, fp)
		case OP_UI8_I32:
			opUI8ToI32(expr, fp)
		case OP_UI8_I64:
			opUI8ToI64(expr, fp)
		case OP_UI8_UI16:
			opUI8ToUI16(expr, fp)
		case OP_UI8_UI32:
			opUI8ToUI32(expr, fp)
		case OP_UI8_UI64:
			opUI8ToUI64(expr, fp)
		case OP_UI8_F32:
			opUI8ToF32(expr, fp)
		case OP_UI8_F64:
			opUI8ToF64(expr, fp)
		case OP_UI8_ADD:
			opUI8Add(expr, fp)
		case OP_UI8_SUB:
			opUI8Sub(expr, fp)
		case OP_UI8_MUL:
			opUI8Mul(expr, fp)
		case OP_UI8_DIV:
			opUI8Div(expr, fp)
		case OP_UI8_GT:
			opUI8Gt(expr, fp)
		case OP_UI8_GTEQ:
			opUI8Gteq(expr, fp)
		case OP_UI8_LT:
			opUI8Lt(expr, fp)
		case OP_UI8_LTEQ:
			opUI8Lteq(expr, fp)
		case OP_UI8_EQ:
			opUI8Eq(expr, fp)
		case OP_UI8_UNEQ:
			opUI8Uneq(expr, fp)
		case OP_UI8_MOD:
			opUI8Mod(expr, fp)
		case OP_UI8_RAND:
			opUI8Rand(expr, fp)
		case OP_UI8_BITAND:
			opUI8Bitand(expr, fp)
		case OP_UI8_BITOR:
			opUI8Bitor(expr, fp)
		case OP_UI8_BITXOR:
			opUI8Bitxor(expr, fp)
		case OP_UI8_BITCLEAR:
			opUI8Bitclear(expr, fp)
		case OP_UI8_BITSHL:
			opUI8Bitshl(expr, fp)
		case OP_UI8_BITSHR:
			opUI8Bitshr(expr, fp)
		case OP_UI8_MAX:
			opUI8Max(expr, fp)
		case OP_UI8_MIN:
			opUI8Min(expr, fp)

		case OP_UI16_STR:
			opUI16ToStr(expr, fp)
		case OP_UI16_I32:
			opUI16ToI32(expr, fp)
		case OP_UI16_I64:
			opUI16ToI64(expr, fp)
		case OP_UI16_UI8:
			opUI16ToUI8(expr, fp)
		case OP_UI16_UI32:
			opUI16ToUI32(expr, fp)
		case OP_UI16_UI64:
			opUI16ToUI64(expr, fp)
		case OP_UI16_F32:
			opUI16ToF32(expr, fp)
		case OP_UI16_F64:
			opUI16ToF64(expr, fp)
		case OP_UI16_ADD:
			opUI16Add(expr, fp)
		case OP_UI16_SUB:
			opUI16Sub(expr, fp)
		case OP_UI16_MUL:
			opUI16Mul(expr, fp)
		case OP_UI16_DIV:
			opUI16Div(expr, fp)
		case OP_UI16_GT:
			opUI16Gt(expr, fp)
		case OP_UI16_GTEQ:
			opUI16Gteq(expr, fp)
		case OP_UI16_LT:
			opUI16Lt(expr, fp)
		case OP_UI16_LTEQ:
			opUI16Lteq(expr, fp)
		case OP_UI16_EQ:
			opUI16Eq(expr, fp)
		case OP_UI16_UNEQ:
			opUI16Uneq(expr, fp)
		case OP_UI16_MOD:
			opUI16Mod(expr, fp)
		case OP_UI16_RAND:
			opUI16Rand(expr, fp)
		case OP_UI16_BITAND:
			opUI16Bitand(expr, fp)
		case OP_UI16_BITOR:
			opUI16Bitor(expr, fp)
		case OP_UI16_BITXOR:
			opUI16Bitxor(expr, fp)
		case OP_UI16_BITCLEAR:
			opUI16Bitclear(expr, fp)
		case OP_UI16_BITSHL:
			opUI16Bitshl(expr, fp)
		case OP_UI16_BITSHR:
			opUI16Bitshr(expr, fp)
		case OP_UI16_MAX:
			opUI16Max(expr, fp)
		case OP_UI16_MIN:
			opUI16Min(expr, fp)

		case OP_UI32_STR:
			opUI32ToStr(expr, fp)
		case OP_UI32_I32:
			opUI32ToI32(expr, fp)
		case OP_UI32_I64:
			opUI32ToI64(expr, fp)
		case OP_UI32_UI8:
			opUI32ToUI8(expr, fp)
		case OP_UI32_UI16:
			opUI32ToUI16(expr, fp)
		case OP_UI32_UI64:
			opUI32ToUI64(expr, fp)
		case OP_UI32_F32:
			opUI32ToF32(expr, fp)
		case OP_UI32_F64:
			opUI32ToF64(expr, fp)
		case OP_UI32_ADD:
			opUI32Add(expr, fp)
		case OP_UI32_SUB:
			opUI32Sub(expr, fp)
		case OP_UI32_MUL:
			opUI32Mul(expr, fp)
		case OP_UI32_DIV:
			opUI32Div(expr, fp)
		case OP_UI32_GT:
			opUI32Gt(expr, fp)
		case OP_UI32_GTEQ:
			opUI32Gteq(expr, fp)
		case OP_UI32_LT:
			opUI32Lt(expr, fp)
		case OP_UI32_LTEQ:
			opUI32Lteq(expr, fp)
		case OP_UI32_EQ:
			opUI32Eq(expr, fp)
		case OP_UI32_UNEQ:
			opUI32Uneq(expr, fp)
		case OP_UI32_MOD:
			opUI32Mod(expr, fp)
		case OP_UI32_RAND:
			opUI32Rand(expr, fp)
		case OP_UI32_BITAND:
			opUI32Bitand(expr, fp)
		case OP_UI32_BITOR:
			opUI32Bitor(expr, fp)
		case OP_UI32_BITXOR:
			opUI32Bitxor(expr, fp)
		case OP_UI32_BITCLEAR:
			opUI32Bitclear(expr, fp)
		case OP_UI32_BITSHL:
			opUI32Bitshl(expr, fp)
		case OP_UI32_BITSHR:
			opUI32Bitshr(expr, fp)
		case OP_UI32_MAX:
			opUI32Max(expr, fp)
		case OP_UI32_MIN:
			opUI32Min(expr, fp)

		case OP_UI64_STR:
			opUI64ToStr(expr, fp)
		case OP_UI64_I32:
			opUI64ToI32(expr, fp)
		case OP_UI64_I64:
			opUI64ToI64(expr, fp)
		case OP_UI64_UI8:
			opUI64ToUI8(expr, fp)
		case OP_UI64_UI16:
			opUI64ToUI16(expr, fp)
		case OP_UI64_UI32:
			opUI64ToUI32(expr, fp)
		case OP_UI64_F32:
			opUI64ToF32(expr, fp)
		case OP_UI64_F64:
			opUI64ToF64(expr, fp)
		case OP_UI64_ADD:
			opUI64Add(expr, fp)
		case OP_UI64_SUB:
			opUI64Sub(expr, fp)
		case OP_UI64_MUL:
			opUI64Mul(expr, fp)
		case OP_UI64_DIV:
			opUI64Div(expr, fp)
		case OP_UI64_GT:
			opUI64Gt(expr, fp)
		case OP_UI64_GTEQ:
			opUI64Gteq(expr, fp)
		case OP_UI64_LT:
			opUI64Lt(expr, fp)
		case OP_UI64_LTEQ:
			opUI64Lteq(expr, fp)
		case OP_UI64_EQ:
			opUI64Eq(expr, fp)
		case OP_UI64_UNEQ:
			opUI64Uneq(expr, fp)
		case OP_UI64_MOD:
			opUI64Mod(expr, fp)
		case OP_UI64_RAND:
			opUI64Rand(expr, fp)
		case OP_UI64_BITAND:
			opUI64Bitand(expr, fp)
		case OP_UI64_BITOR:
			opUI64Bitor(expr, fp)
		case OP_UI64_BITXOR:
			opUI64Bitxor(expr, fp)
		case OP_UI64_BITCLEAR:
			opUI64Bitclear(expr, fp)
		case OP_UI64_BITSHL:
			opUI64Bitshl(expr, fp)
		case OP_UI64_BITSHR:
			opUI64Bitshr(expr, fp)
		case OP_UI64_MAX:
			opUI64Max(expr, fp)
		case OP_UI64_MIN:
			opUI64Min(expr, fp)

		case OP_F32_STR:
			opF32Cast(expr, fp)
		case OP_F32_I32:
			opF32Cast(expr, fp)
		case OP_F32_I64:
			opF32Cast(expr, fp)
		case OP_F32_BYTE:
			opF32Cast(expr, fp)
		case OP_F32_UI8:
			opF32Cast(expr, fp)
		case OP_F32_UI16:
			opF32Cast(expr, fp)
		case OP_F32_UI32:
			opF32Cast(expr, fp)
		case OP_F32_UI64:
			opF32Cast(expr, fp)
		case OP_F32_F32:
			opF32Cast(expr, fp)
		case OP_F32_F64:
			opF32Cast(expr, fp)
		case OP_F32_PRINT:
			opF32Print(expr, fp)
		case OP_F32_ADD:
			opF32Add(expr, fp)
		case OP_F32_SUB, OP_F32_NEG:
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
		case OP_F32_IS_NAN:
			opF32Isnan(expr, fp)

		case OP_F64_STR:
			opF64Cast(expr, fp)
		case OP_F64_I32:
			opF64Cast(expr, fp)
		case OP_F64_I64:
			opF64Cast(expr, fp)
		case OP_F64_BYTE:
			opF64Cast(expr, fp)
		case OP_F64_UI8:
			opF64Cast(expr, fp)
		case OP_F64_UI16:
			opF64Cast(expr, fp)
		case OP_F64_UI32:
			opF64Cast(expr, fp)
		case OP_F64_UI64:
			opF64Cast(expr, fp)
		case OP_F64_F32:
			opF64Cast(expr, fp)
		case OP_F64_F64:
			opF64Cast(expr, fp)
		case OP_F64_PRINT:
			opF64Print(expr, fp)
		case OP_F64_ADD:
			opF64Add(expr, fp)
		case OP_F64_SUB, OP_F64_NEG:
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

		case OP_STR_STR:
			opStrStr(expr, fp)
		case OP_STR_I32:
			opStrStr(expr, fp)
		case OP_STR_I64:
			opStrStr(expr, fp)
		case OP_STR_BYTE:
			opStrStr(expr, fp)
		case OP_STR_UI8:
			opStrStr(expr, fp)
		case OP_STR_UI16:
			opStrStr(expr, fp)
		case OP_STR_UI32:
			opStrStr(expr, fp)
		case OP_STR_UI64:
			opStrStr(expr, fp)
		case OP_STR_F32:
			opStrStr(expr, fp)
		case OP_STR_F64:
			opStrStr(expr, fp)
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

		case OP_APPEND:
			opAppend(expr, fp)
		case OP_RESIZE:
			opResize(expr, fp)
		case OP_INSERT:
			opInsert(expr, fp)
		case OP_REMOVE:
			opRemove(expr, fp)
		case OP_COPY:
			opCopy(expr, fp)

		case OP_MAKE:
		case OP_READ:
		case OP_WRITE:
		case OP_LEN:
		case OP_CONCAT:
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
			opAssertValue(expr, fp)
		case OP_TEST:
			opTest(expr, fp)
		case OP_PANIC:
			opPanic(expr, fp)
		case OP_STRERROR:
			opStrError(expr, fp)

		// affordances
		case OP_AFF_PRINT:
			opAffPrint(expr, fp)
		case OP_AFF_QUERY:
			opAffQuery(expr, fp)
		case OP_AFF_ON:
			opAffOn(expr, fp)
		case OP_AFF_OF:
			opAffOf(expr, fp)
		case OP_AFF_INFORM:
			opAffInform(expr, fp)
		case OP_AFF_REQUEST:
			opAffRequest(expr, fp)
		default:
			// DumpOpCodes(opCode) // debug helper
			panic("invalid bare opcode")
		}
	}

	execNative = execNativeCore
}
