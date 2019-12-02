package cxcore

// CorePackages ...
var CorePackages = []string{
	// temporary solution until we can implement these packages in pure CX I guess
	"gl", "glfw", "time", "http", "os", "explorer", "aff", "gltext", "cx", "json",
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
	OP_UND_NEG
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

	OP_I8_STR
	OP_I8_I16
	OP_I8_I32
	OP_I8_I64
	OP_I8_UI8
	OP_I8_UI16
	OP_I8_UI32
	OP_I8_UI64
	OP_I8_F32
	OP_I8_F64

	OP_I16_STR
	OP_I16_I8
	OP_I16_I32
	OP_I16_I64
	OP_I16_UI8
	OP_I16_UI16
	OP_I16_UI32
	OP_I16_UI64
	OP_I16_F32
	OP_I16_F64

	OP_I32_STR
	OP_I32_I8
	OP_I32_I16
	OP_I32_I64
	OP_I32_UI8
	OP_I32_UI16
	OP_I32_UI32
	OP_I32_UI64
	OP_I32_F32
	OP_I32_F64

	OP_I64_STR
	OP_I64_I8
	OP_I64_I16
	OP_I64_I32
	OP_I64_UI8
	OP_I64_UI16
	OP_I64_UI32
	OP_I64_UI64
	OP_I64_F32
	OP_I64_F64

	OP_UI8_STR
	OP_UI8_I8
	OP_UI8_I16
	OP_UI8_I32
	OP_UI8_I64
	OP_UI8_UI16
	OP_UI8_UI32
	OP_UI8_UI64
	OP_UI8_F32
	OP_UI8_F64

	OP_UI16_STR
	OP_UI16_I8
	OP_UI16_I16
	OP_UI16_I32
	OP_UI16_I64
	OP_UI16_UI8
	OP_UI16_UI32
	OP_UI16_UI64
	OP_UI16_F32
	OP_UI16_F64

	OP_UI32_STR
	OP_UI32_I8
	OP_UI32_I16
	OP_UI32_I32
	OP_UI32_I64
	OP_UI32_UI8
	OP_UI32_UI16
	OP_UI32_UI64
	OP_UI32_F32
	OP_UI32_F64

	OP_UI64_STR
	OP_UI64_I8
	OP_UI64_I16
	OP_UI64_I32
	OP_UI64_I64
	OP_UI64_UI8
	OP_UI64_UI16
	OP_UI64_UI32
	OP_UI64_F32
	OP_UI64_F64

	OP_F32_STR
	OP_F32_I8
	OP_F32_I16
	OP_F32_I32
	OP_F32_I64
	OP_F32_UI8
	OP_F32_UI16
	OP_F32_UI32
	OP_F32_UI64
	OP_F32_F32
	OP_F32_F64

	OP_F64_STR
	OP_F64_I8
	OP_F64_I16
	OP_F64_I32
	OP_F64_I64
	OP_F64_UI8
	OP_F64_UI16
	OP_F64_UI32
	OP_F64_UI64
	OP_F64_F32
	OP_F64_F64

	OP_STR_STR
	OP_STR_I8
	OP_STR_I16
	OP_STR_I32
	OP_STR_I64
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

	OP_I8_PRINT
	OP_I8_ADD
	OP_I8_SUB
	OP_I8_NEG
	OP_I8_MUL
	OP_I8_DIV
	OP_I8_MOD
	OP_I8_ABS
	OP_I8_GT
	OP_I8_GTEQ
	OP_I8_LT
	OP_I8_LTEQ
	OP_I8_EQ
	OP_I8_UNEQ
	OP_I8_BITAND
	OP_I8_BITOR
	OP_I8_BITXOR
	OP_I8_BITCLEAR
	OP_I8_BITSHL
	OP_I8_BITSHR
	OP_I8_MAX
	OP_I8_MIN
	OP_I8_RAND

	OP_I16_PRINT
	OP_I16_ADD
	OP_I16_SUB
	OP_I16_NEG
	OP_I16_MUL
	OP_I16_DIV
	OP_I16_MOD
	OP_I16_ABS
	OP_I16_GT
	OP_I16_GTEQ
	OP_I16_LT
	OP_I16_LTEQ
	OP_I16_EQ
	OP_I16_UNEQ
	OP_I16_BITAND
	OP_I16_BITOR
	OP_I16_BITXOR
	OP_I16_BITCLEAR
	OP_I16_BITSHL
	OP_I16_BITSHR
	OP_I16_MAX
	OP_I16_MIN
	OP_I16_RAND

	OP_I32_PRINT
	OP_I32_ADD
	OP_I32_SUB
	OP_I32_NEG
	OP_I32_MUL
	OP_I32_DIV
	OP_I32_MOD
	OP_I32_ABS
	OP_I32_GT
	OP_I32_GTEQ
	OP_I32_LT
	OP_I32_LTEQ
	OP_I32_EQ
	OP_I32_UNEQ
	OP_I32_BITAND
	OP_I32_BITOR
	OP_I32_BITXOR
	OP_I32_BITCLEAR
	OP_I32_BITSHL
	OP_I32_BITSHR
	OP_I32_MAX
	OP_I32_MIN
	OP_I32_RAND

	OP_I64_PRINT
	OP_I64_ADD
	OP_I64_SUB
	OP_I64_NEG
	OP_I64_MUL
	OP_I64_DIV
	OP_I64_MOD
	OP_I64_ABS
	OP_I64_GT
	OP_I64_GTEQ
	OP_I64_LT
	OP_I64_LTEQ
	OP_I64_EQ
	OP_I64_UNEQ
	OP_I64_BITAND
	OP_I64_BITOR
	OP_I64_BITXOR
	OP_I64_BITCLEAR
	OP_I64_BITSHL
	OP_I64_BITSHR
	OP_I64_MAX
	OP_I64_MIN
	OP_I64_RAND

	OP_UI8_PRINT
	OP_UI8_ADD
	OP_UI8_SUB
	OP_UI8_MUL
	OP_UI8_DIV
	OP_UI8_MOD
	OP_UI8_GT
	OP_UI8_GTEQ
	OP_UI8_LT
	OP_UI8_LTEQ
	OP_UI8_EQ
	OP_UI8_UNEQ
	OP_UI8_BITAND
	OP_UI8_BITOR
	OP_UI8_BITXOR
	OP_UI8_BITCLEAR
	OP_UI8_BITSHL
	OP_UI8_BITSHR
	OP_UI8_MAX
	OP_UI8_MIN
	OP_UI8_RAND

	OP_UI16_PRINT
	OP_UI16_ADD
	OP_UI16_SUB
	OP_UI16_MUL
	OP_UI16_DIV
	OP_UI16_MOD
	OP_UI16_GT
	OP_UI16_GTEQ
	OP_UI16_LT
	OP_UI16_LTEQ
	OP_UI16_EQ
	OP_UI16_UNEQ
	OP_UI16_BITAND
	OP_UI16_BITOR
	OP_UI16_BITXOR
	OP_UI16_BITCLEAR
	OP_UI16_BITSHL
	OP_UI16_BITSHR
	OP_UI16_MAX
	OP_UI16_MIN
	OP_UI16_RAND

	OP_UI32_PRINT
	OP_UI32_ADD
	OP_UI32_SUB
	OP_UI32_MUL
	OP_UI32_DIV
	OP_UI32_MOD
	OP_UI32_GT
	OP_UI32_GTEQ
	OP_UI32_LT
	OP_UI32_LTEQ
	OP_UI32_EQ
	OP_UI32_UNEQ
	OP_UI32_BITAND
	OP_UI32_BITOR
	OP_UI32_BITXOR
	OP_UI32_BITCLEAR
	OP_UI32_BITSHL
	OP_UI32_BITSHR
	OP_UI32_MAX
	OP_UI32_MIN
	OP_UI32_RAND

	OP_UI64_PRINT
	OP_UI64_ADD
	OP_UI64_SUB
	OP_UI64_MUL
	OP_UI64_DIV
	OP_UI64_MOD
	OP_UI64_GT
	OP_UI64_GTEQ
	OP_UI64_LT
	OP_UI64_LTEQ
	OP_UI64_EQ
	OP_UI64_UNEQ
	OP_UI64_BITAND
	OP_UI64_BITOR
	OP_UI64_BITXOR
	OP_UI64_BITCLEAR
	OP_UI64_BITSHL
	OP_UI64_BITSHR
	OP_UI64_MAX
	OP_UI64_MIN
	OP_UI64_RAND

	OP_F32_IS_NAN
	OP_F32_PRINT
	OP_F32_ADD
	OP_F32_SUB
	OP_F32_NEG
	OP_F32_MUL
	OP_F32_DIV
	OP_F32_MOD
	OP_F32_ABS
	OP_F32_POW
	OP_F32_GT
	OP_F32_GTEQ
	OP_F32_LT
	OP_F32_LTEQ
	OP_F32_EQ
	OP_F32_UNEQ
	OP_F32_ACOS
	OP_F32_COS
	OP_F32_ASIN
	OP_F32_SIN
	OP_F32_SQRT
	OP_F32_LOG
	OP_F32_LOG2
	OP_F32_LOG10
	OP_F32_MAX
	OP_F32_MIN
	OP_F32_RAND

	OP_F64_IS_NAN
	OP_F64_PRINT
	OP_F64_ADD
	OP_F64_SUB
	OP_F64_NEG
	OP_F64_MUL
	OP_F64_DIV
	OP_F64_MOD
	OP_F64_ABS
	OP_F64_POW
	OP_F64_GT
	OP_F64_GTEQ
	OP_F64_LT
	OP_F64_LTEQ
	OP_F64_EQ
	OP_F64_UNEQ
	OP_F64_ACOS
	OP_F64_COS
	OP_F64_ASIN
	OP_F64_SIN
	OP_F64_SQRT
	OP_F64_LOG
	OP_F64_LOG2
	OP_F64_LOG10
	OP_F64_MAX
	OP_F64_MIN
	OP_F64_RAND

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

	END_OF_CORE_OPS
)

// For the parser. These shouldn't be used in the runtime for performance reasons
var (
	OpNames              = map[int]string{}
	OpCodes              = map[string]int{}
	Natives              = map[int]*CXFunction{}
	opcodeHandlerFinders []opcodeHandlerFinder
)

type opcodeHandler func(prgrm *CXProgram)
type opcodeHandlerFinder func(opCode int) opcodeHandler

func execNative(prgrm *CXProgram) {
	defer RuntimeError()
	opCode := prgrm.GetOpCode()

	var handler opcodeHandler
	for _, f := range opcodeHandlerFinders {
		if x := f(opCode); x != nil {
			if handler != nil {
				// This means 2 or more opcodeHandlerFinder functions handle
				// the same opcode. They must handle unique opcodes only.
				panic("multiple opcode handler finders returned an opcode handler")
			}
			handler = x
		}
	}

	if handler == nil {
		panic("invalid bare opcode")
	}

	handler(prgrm)
}

// AddOpCode ...
func AddOpCode(code int, name string, inputs []*CXArgument, outputs []*CXArgument) {
	OpNames[code] = name
	OpCodes[name] = code
	Natives[code] = MakeNativeFunction(code, inputs, outputs)
}

/*
// debug helper
func DumpOpCodes(opCode int) {
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
		[]*CXArgument{newOpPar(TYPE_UI8, false)})
	AddOpCode(OP_DESERIALIZE, "deserialize",
		[]*CXArgument{newOpPar(TYPE_UI8, false)},
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

	AddOpCode(OP_I8_STR, "i8.str",
		[]*CXArgument{newOpPar(TYPE_I8, false)},
		[]*CXArgument{newOpPar(TYPE_STR, false)})
	AddOpCode(OP_I8_I16, "i8.i16",
		[]*CXArgument{newOpPar(TYPE_I8, false)},
		[]*CXArgument{newOpPar(TYPE_I16, false)})
	AddOpCode(OP_I8_I32, "i8.i32",
		[]*CXArgument{newOpPar(TYPE_I8, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false)})
	AddOpCode(OP_I8_I64, "i8.i64",
		[]*CXArgument{newOpPar(TYPE_I8, false)},
		[]*CXArgument{newOpPar(TYPE_I64, false)})
	AddOpCode(OP_I8_UI8, "i8.ui8",
		[]*CXArgument{newOpPar(TYPE_I8, false)},
		[]*CXArgument{newOpPar(TYPE_UI8, false)})
	AddOpCode(OP_I8_UI16, "i8.ui16",
		[]*CXArgument{newOpPar(TYPE_I8, false)},
		[]*CXArgument{newOpPar(TYPE_UI16, false)})
	AddOpCode(OP_I8_UI32, "i8.ui32",
		[]*CXArgument{newOpPar(TYPE_I8, false)},
		[]*CXArgument{newOpPar(TYPE_UI32, false)})
	AddOpCode(OP_I8_UI64, "i8.ui64",
		[]*CXArgument{newOpPar(TYPE_I8, false)},
		[]*CXArgument{newOpPar(TYPE_UI64, false)})
	AddOpCode(OP_I8_F32, "i8.f32",
		[]*CXArgument{newOpPar(TYPE_I8, false)},
		[]*CXArgument{newOpPar(TYPE_F32, false)})
	AddOpCode(OP_I8_F64, "i8.f64",
		[]*CXArgument{newOpPar(TYPE_I8, false)},
		[]*CXArgument{newOpPar(TYPE_F64, false)})
	AddOpCode(OP_I8_PRINT, "i8.print",
		[]*CXArgument{newOpPar(TYPE_I8, false)},
		[]*CXArgument{})
	AddOpCode(OP_I8_ADD, "i8.add",
		[]*CXArgument{newOpPar(TYPE_I8, false), newOpPar(TYPE_I8, false)},
		[]*CXArgument{newOpPar(TYPE_I8, false)})
	AddOpCode(OP_I8_SUB, "i8.sub",
		[]*CXArgument{newOpPar(TYPE_I8, false), newOpPar(TYPE_I8, false)},
		[]*CXArgument{newOpPar(TYPE_I8, false)})
	AddOpCode(OP_I8_NEG, "i8.neg",
		[]*CXArgument{newOpPar(TYPE_I8, false)},
		[]*CXArgument{newOpPar(TYPE_I8, false)})
	AddOpCode(OP_I8_MUL, "i8.mul",
		[]*CXArgument{newOpPar(TYPE_I8, false), newOpPar(TYPE_I8, false)},
		[]*CXArgument{newOpPar(TYPE_I8, false)})
	AddOpCode(OP_I8_DIV, "i8.div",
		[]*CXArgument{newOpPar(TYPE_I8, false), newOpPar(TYPE_I8, false)},
		[]*CXArgument{newOpPar(TYPE_I8, false)})
	AddOpCode(OP_I8_MOD, "i8.mod",
		[]*CXArgument{newOpPar(TYPE_I8, false), newOpPar(TYPE_I8, false)},
		[]*CXArgument{newOpPar(TYPE_I8, false)})
	AddOpCode(OP_I8_ABS, "i8.abs",
		[]*CXArgument{newOpPar(TYPE_I8, false)},
		[]*CXArgument{newOpPar(TYPE_I8, false)})
	AddOpCode(OP_I8_GT, "i8.gt",
		[]*CXArgument{newOpPar(TYPE_I8, false), newOpPar(TYPE_I8, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_I8_GTEQ, "i8.gteq",
		[]*CXArgument{newOpPar(TYPE_I8, false), newOpPar(TYPE_I8, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_I8_LT, "i8.lt",
		[]*CXArgument{newOpPar(TYPE_I8, false), newOpPar(TYPE_I8, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_I8_LTEQ, "i8.lteq",
		[]*CXArgument{newOpPar(TYPE_I8, false), newOpPar(TYPE_I8, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_I8_EQ, "i8.eq",
		[]*CXArgument{newOpPar(TYPE_I8, false), newOpPar(TYPE_I8, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_I8_UNEQ, "i8.uneq",
		[]*CXArgument{newOpPar(TYPE_I8, false), newOpPar(TYPE_I8, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_I8_BITAND, "i8.bitand",
		[]*CXArgument{newOpPar(TYPE_I8, false), newOpPar(TYPE_I8, false)},
		[]*CXArgument{newOpPar(TYPE_I8, false)})
	AddOpCode(OP_I8_BITOR, "i8.bitor",
		[]*CXArgument{newOpPar(TYPE_I8, false), newOpPar(TYPE_I8, false)},
		[]*CXArgument{newOpPar(TYPE_I8, false)})
	AddOpCode(OP_I8_BITXOR, "i8.bitxor",
		[]*CXArgument{newOpPar(TYPE_I8, false), newOpPar(TYPE_I8, false)},
		[]*CXArgument{newOpPar(TYPE_I8, false)})
	AddOpCode(OP_I8_BITCLEAR, "i8.bitclear",
		[]*CXArgument{newOpPar(TYPE_I8, false), newOpPar(TYPE_I8, false)},
		[]*CXArgument{newOpPar(TYPE_I8, false)})
	AddOpCode(OP_I8_BITSHL, "i8.bitshl",
		[]*CXArgument{newOpPar(TYPE_I8, false), newOpPar(TYPE_I8, false)},
		[]*CXArgument{newOpPar(TYPE_I8, false)})
	AddOpCode(OP_I8_BITSHR, "i8.bitshr",
		[]*CXArgument{newOpPar(TYPE_I8, false), newOpPar(TYPE_I8, false)},
		[]*CXArgument{newOpPar(TYPE_I8, false)})
	AddOpCode(OP_I8_MAX, "i8.max",
		[]*CXArgument{newOpPar(TYPE_I8, false), newOpPar(TYPE_I8, false)},
		[]*CXArgument{newOpPar(TYPE_I8, false)})
	AddOpCode(OP_I8_MIN, "i8.min",
		[]*CXArgument{newOpPar(TYPE_I8, false), newOpPar(TYPE_I8, false)},
		[]*CXArgument{newOpPar(TYPE_I8, false)})
	AddOpCode(OP_I8_RAND, "i8.rand",
		[]*CXArgument{newOpPar(TYPE_I8, false), newOpPar(TYPE_I8, false)},
		[]*CXArgument{newOpPar(TYPE_I8, false)})

	AddOpCode(OP_I16_STR, "i16.str",
		[]*CXArgument{newOpPar(TYPE_I16, false)},
		[]*CXArgument{newOpPar(TYPE_STR, false)})
	AddOpCode(OP_I16_I8, "i16.i8",
		[]*CXArgument{newOpPar(TYPE_I16, false)},
		[]*CXArgument{newOpPar(TYPE_I8, false)})
	AddOpCode(OP_I16_I32, "i16.i32",
		[]*CXArgument{newOpPar(TYPE_I16, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false)})
	AddOpCode(OP_I16_I64, "i16.i64",
		[]*CXArgument{newOpPar(TYPE_I16, false)},
		[]*CXArgument{newOpPar(TYPE_I64, false)})
	AddOpCode(OP_I16_UI8, "i16.ui8",
		[]*CXArgument{newOpPar(TYPE_I16, false)},
		[]*CXArgument{newOpPar(TYPE_UI8, false)})
	AddOpCode(OP_I16_UI16, "i16.ui16",
		[]*CXArgument{newOpPar(TYPE_I16, false)},
		[]*CXArgument{newOpPar(TYPE_UI16, false)})
	AddOpCode(OP_I16_UI32, "i16.ui32",
		[]*CXArgument{newOpPar(TYPE_I16, false)},
		[]*CXArgument{newOpPar(TYPE_UI32, false)})
	AddOpCode(OP_I16_UI64, "i16.ui64",
		[]*CXArgument{newOpPar(TYPE_I16, false)},
		[]*CXArgument{newOpPar(TYPE_UI64, false)})
	AddOpCode(OP_I16_F32, "i16.f32",
		[]*CXArgument{newOpPar(TYPE_I16, false)},
		[]*CXArgument{newOpPar(TYPE_F32, false)})
	AddOpCode(OP_I16_F64, "i16.f64",
		[]*CXArgument{newOpPar(TYPE_I16, false)},
		[]*CXArgument{newOpPar(TYPE_F64, false)})
	AddOpCode(OP_I16_PRINT, "i16.print",
		[]*CXArgument{newOpPar(TYPE_I16, false)},
		[]*CXArgument{})
	AddOpCode(OP_I16_ADD, "i16.add",
		[]*CXArgument{newOpPar(TYPE_I16, false), newOpPar(TYPE_I16, false)},
		[]*CXArgument{newOpPar(TYPE_I16, false)})
	AddOpCode(OP_I16_SUB, "i16.sub",
		[]*CXArgument{newOpPar(TYPE_I16, false), newOpPar(TYPE_I16, false)},
		[]*CXArgument{newOpPar(TYPE_I16, false)})
	AddOpCode(OP_I16_NEG, "i16.neg",
		[]*CXArgument{newOpPar(TYPE_I16, false)},
		[]*CXArgument{newOpPar(TYPE_I16, false)})
	AddOpCode(OP_I16_MUL, "i16.mul",
		[]*CXArgument{newOpPar(TYPE_I16, false), newOpPar(TYPE_I16, false)},
		[]*CXArgument{newOpPar(TYPE_I16, false)})
	AddOpCode(OP_I16_DIV, "i16.div",
		[]*CXArgument{newOpPar(TYPE_I16, false), newOpPar(TYPE_I16, false)},
		[]*CXArgument{newOpPar(TYPE_I16, false)})
	AddOpCode(OP_I16_MOD, "i16.mod",
		[]*CXArgument{newOpPar(TYPE_I16, false), newOpPar(TYPE_I16, false)},
		[]*CXArgument{newOpPar(TYPE_I16, false)})
	AddOpCode(OP_I16_ABS, "i16.abs",
		[]*CXArgument{newOpPar(TYPE_I16, false)},
		[]*CXArgument{newOpPar(TYPE_I16, false)})
	AddOpCode(OP_I16_GT, "i16.gt",
		[]*CXArgument{newOpPar(TYPE_I16, false), newOpPar(TYPE_I16, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_I16_GTEQ, "i16.gteq",
		[]*CXArgument{newOpPar(TYPE_I16, false), newOpPar(TYPE_I16, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_I16_LT, "i16.lt",
		[]*CXArgument{newOpPar(TYPE_I16, false), newOpPar(TYPE_I16, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_I16_LTEQ, "i16.lteq",
		[]*CXArgument{newOpPar(TYPE_I16, false), newOpPar(TYPE_I16, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_I16_EQ, "i16.eq",
		[]*CXArgument{newOpPar(TYPE_I16, false), newOpPar(TYPE_I16, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_I16_UNEQ, "i16.uneq",
		[]*CXArgument{newOpPar(TYPE_I16, false), newOpPar(TYPE_I16, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_I16_BITAND, "i16.bitand",
		[]*CXArgument{newOpPar(TYPE_I16, false), newOpPar(TYPE_I16, false)},
		[]*CXArgument{newOpPar(TYPE_I16, false)})
	AddOpCode(OP_I16_BITOR, "i16.bitor",
		[]*CXArgument{newOpPar(TYPE_I16, false), newOpPar(TYPE_I16, false)},
		[]*CXArgument{newOpPar(TYPE_I16, false)})
	AddOpCode(OP_I16_BITXOR, "i16.bitxor",
		[]*CXArgument{newOpPar(TYPE_I16, false), newOpPar(TYPE_I16, false)},
		[]*CXArgument{newOpPar(TYPE_I16, false)})
	AddOpCode(OP_I16_BITCLEAR, "i16.bitclear",
		[]*CXArgument{newOpPar(TYPE_I16, false), newOpPar(TYPE_I16, false)},
		[]*CXArgument{newOpPar(TYPE_I16, false)})
	AddOpCode(OP_I16_BITSHL, "i16.bitshl",
		[]*CXArgument{newOpPar(TYPE_I16, false), newOpPar(TYPE_I16, false)},
		[]*CXArgument{newOpPar(TYPE_I16, false)})
	AddOpCode(OP_I16_BITSHR, "i16.bitshr",
		[]*CXArgument{newOpPar(TYPE_I16, false), newOpPar(TYPE_I16, false)},
		[]*CXArgument{newOpPar(TYPE_I16, false)})
	AddOpCode(OP_I16_MAX, "i16.max",
		[]*CXArgument{newOpPar(TYPE_I16, false), newOpPar(TYPE_I16, false)},
		[]*CXArgument{newOpPar(TYPE_I16, false)})
	AddOpCode(OP_I16_MIN, "i16.min",
		[]*CXArgument{newOpPar(TYPE_I16, false), newOpPar(TYPE_I16, false)},
		[]*CXArgument{newOpPar(TYPE_I16, false)})
	AddOpCode(OP_I16_RAND, "i16.rand",
		[]*CXArgument{newOpPar(TYPE_I16, false), newOpPar(TYPE_I16, false)},
		[]*CXArgument{newOpPar(TYPE_I16, false)})

	AddOpCode(OP_I32_STR, "i32.str",
		[]*CXArgument{newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_STR, false)})
	AddOpCode(OP_I32_I8, "i32.i8",
		[]*CXArgument{newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_I8, false)})
	AddOpCode(OP_I32_I16, "i32.i16",
		[]*CXArgument{newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_I16, false)})
	AddOpCode(OP_I32_I64, "i32.i64",
		[]*CXArgument{newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_I64, false)})
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
	AddOpCode(OP_I32_MOD, "i32.mod",
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false)})
	AddOpCode(OP_I32_ABS, "i32.abs",
		[]*CXArgument{newOpPar(TYPE_I32, false)},
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
	AddOpCode(OP_I32_MAX, "i32.max",
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false)})
	AddOpCode(OP_I32_MIN, "i32.min",
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false)})
	AddOpCode(OP_I32_RAND, "i32.rand",
		[]*CXArgument{newOpPar(TYPE_I32, false), newOpPar(TYPE_I32, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false)})

	AddOpCode(OP_I64_STR, "i64.str",
		[]*CXArgument{newOpPar(TYPE_I64, false)},
		[]*CXArgument{newOpPar(TYPE_STR, false)})
	AddOpCode(OP_I64_I8, "i64.i8",
		[]*CXArgument{newOpPar(TYPE_I64, false)},
		[]*CXArgument{newOpPar(TYPE_I8, false)})
	AddOpCode(OP_I64_I16, "i64.i16",
		[]*CXArgument{newOpPar(TYPE_I64, false)},
		[]*CXArgument{newOpPar(TYPE_I16, false)})
	AddOpCode(OP_I64_I32, "i64.i32",
		[]*CXArgument{newOpPar(TYPE_I64, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false)})
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
	AddOpCode(OP_I64_MOD, "i64.mod",
		[]*CXArgument{newOpPar(TYPE_I64, false), newOpPar(TYPE_I64, false)},
		[]*CXArgument{newOpPar(TYPE_I64, false)})
	AddOpCode(OP_I64_ABS, "i64.abs",
		[]*CXArgument{newOpPar(TYPE_I64, false)},
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
	AddOpCode(OP_I64_MAX, "i64.max",
		[]*CXArgument{newOpPar(TYPE_I64, false), newOpPar(TYPE_I64, false)},
		[]*CXArgument{newOpPar(TYPE_I64, false)})
	AddOpCode(OP_I64_MIN, "i64.min",
		[]*CXArgument{newOpPar(TYPE_I64, false), newOpPar(TYPE_I64, false)},
		[]*CXArgument{newOpPar(TYPE_I64, false)})
	AddOpCode(OP_I64_RAND, "i64.rand",
		[]*CXArgument{newOpPar(TYPE_I64, false), newOpPar(TYPE_I64, false)},
		[]*CXArgument{newOpPar(TYPE_I64, false)})

	AddOpCode(OP_UI8_STR, "ui8.str",
		[]*CXArgument{newOpPar(TYPE_UI8, false)},
		[]*CXArgument{newOpPar(TYPE_STR, false)})
	AddOpCode(OP_UI8_I8, "ui8.i8",
		[]*CXArgument{newOpPar(TYPE_UI8, false)},
		[]*CXArgument{newOpPar(TYPE_I8, false)})
	AddOpCode(OP_UI8_I16, "ui8.i16",
		[]*CXArgument{newOpPar(TYPE_UI8, false)},
		[]*CXArgument{newOpPar(TYPE_I16, false)})
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
	AddOpCode(OP_UI8_PRINT, "ui8.print",
		[]*CXArgument{newOpPar(TYPE_UI8, false)},
		[]*CXArgument{})
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
	AddOpCode(OP_UI8_MOD, "ui8.mod",
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
	AddOpCode(OP_UI8_RAND, "ui8.rand",
		[]*CXArgument{},
		[]*CXArgument{newOpPar(TYPE_UI8, false)})

	AddOpCode(OP_UI16_STR, "ui16.str",
		[]*CXArgument{newOpPar(TYPE_UI16, false)},
		[]*CXArgument{newOpPar(TYPE_STR, false)})
	AddOpCode(OP_UI16_I8, "ui16.i8",
		[]*CXArgument{newOpPar(TYPE_UI16, false)},
		[]*CXArgument{newOpPar(TYPE_I8, false)})
	AddOpCode(OP_UI16_I16, "ui16.i16",
		[]*CXArgument{newOpPar(TYPE_UI16, false)},
		[]*CXArgument{newOpPar(TYPE_I16, false)})
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
	AddOpCode(OP_UI16_PRINT, "ui16.print",
		[]*CXArgument{newOpPar(TYPE_UI16, false)},
		[]*CXArgument{})
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
	AddOpCode(OP_UI16_MOD, "ui16.mod",
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
	AddOpCode(OP_UI16_RAND, "ui16.rand",
		[]*CXArgument{},
		[]*CXArgument{newOpPar(TYPE_UI16, false)})

	AddOpCode(OP_UI32_STR, "ui32.str",
		[]*CXArgument{newOpPar(TYPE_UI32, false)},
		[]*CXArgument{newOpPar(TYPE_STR, false)})
	AddOpCode(OP_UI32_I8, "ui32.i8",
		[]*CXArgument{newOpPar(TYPE_UI32, false)},
		[]*CXArgument{newOpPar(TYPE_I8, false)})
	AddOpCode(OP_UI32_I16, "ui32.i16",
		[]*CXArgument{newOpPar(TYPE_UI32, false)},
		[]*CXArgument{newOpPar(TYPE_I16, false)})
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
	AddOpCode(OP_UI32_PRINT, "ui32.print",
		[]*CXArgument{newOpPar(TYPE_UI32, false)},
		[]*CXArgument{})
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
	AddOpCode(OP_UI32_MOD, "ui32.mod",
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
	AddOpCode(OP_UI32_RAND, "ui32.rand",
		[]*CXArgument{},
		[]*CXArgument{newOpPar(TYPE_UI32, false)})

	AddOpCode(OP_UI64_STR, "ui64.str",
		[]*CXArgument{newOpPar(TYPE_UI64, false)},
		[]*CXArgument{newOpPar(TYPE_STR, false)})
	AddOpCode(OP_UI64_I8, "ui64.i8",
		[]*CXArgument{newOpPar(TYPE_UI64, false)},
		[]*CXArgument{newOpPar(TYPE_I8, false)})
	AddOpCode(OP_UI64_I16, "ui64.i16",
		[]*CXArgument{newOpPar(TYPE_UI64, false)},
		[]*CXArgument{newOpPar(TYPE_I16, false)})
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
	AddOpCode(OP_UI64_PRINT, "ui64.print",
		[]*CXArgument{newOpPar(TYPE_UI64, false)},
		[]*CXArgument{})
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
	AddOpCode(OP_UI64_MOD, "ui64.mod",
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
	AddOpCode(OP_UI64_RAND, "ui64.rand",
		[]*CXArgument{},
		[]*CXArgument{newOpPar(TYPE_UI64, false)})

	AddOpCode(OP_F32_IS_NAN, "f32.isnan",
		[]*CXArgument{newOpPar(TYPE_F32, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_F32_STR, "f32.str",
		[]*CXArgument{newOpPar(TYPE_F32, false)},
		[]*CXArgument{newOpPar(TYPE_STR, false)})
	AddOpCode(OP_F32_I8, "f32.i8",
		[]*CXArgument{newOpPar(TYPE_F32, false)},
		[]*CXArgument{newOpPar(TYPE_I8, false)})
	AddOpCode(OP_F32_I16, "f32.i16",
		[]*CXArgument{newOpPar(TYPE_F32, false)},
		[]*CXArgument{newOpPar(TYPE_I16, false)})
	AddOpCode(OP_F32_I32, "f32.i32",
		[]*CXArgument{newOpPar(TYPE_F32, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false)})
	AddOpCode(OP_F32_I64, "f32.i64",
		[]*CXArgument{newOpPar(TYPE_F32, false)},
		[]*CXArgument{newOpPar(TYPE_I64, false)})
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
	AddOpCode(OP_F32_MOD, "f32.mod",
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
	AddOpCode(OP_F32_ACOS, "f32.acos",
		[]*CXArgument{newOpPar(TYPE_F32, false)},
		[]*CXArgument{newOpPar(TYPE_F32, false)})
	AddOpCode(OP_F32_COS, "f32.cos",
		[]*CXArgument{newOpPar(TYPE_F32, false)},
		[]*CXArgument{newOpPar(TYPE_F32, false)})
	AddOpCode(OP_F32_ASIN, "f32.asin",
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
	AddOpCode(OP_F32_RAND, "f32.rand",
		[]*CXArgument{},
		[]*CXArgument{newOpPar(TYPE_F32, false)})

	AddOpCode(OP_F64_IS_NAN, "f64.isnan",
		[]*CXArgument{newOpPar(TYPE_F64, false)},
		[]*CXArgument{newOpPar(TYPE_BOOL, false)})
	AddOpCode(OP_F64_STR, "f64.str",
		[]*CXArgument{newOpPar(TYPE_F64, false)},
		[]*CXArgument{newOpPar(TYPE_STR, false)})
	AddOpCode(OP_F64_I8, "f64.i8",
		[]*CXArgument{newOpPar(TYPE_F64, false)},
		[]*CXArgument{newOpPar(TYPE_I8, false)})
	AddOpCode(OP_F64_I16, "f64.i16",
		[]*CXArgument{newOpPar(TYPE_F64, false)},
		[]*CXArgument{newOpPar(TYPE_I16, false)})
	AddOpCode(OP_F64_I32, "f64.i32",
		[]*CXArgument{newOpPar(TYPE_F64, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false)})
	AddOpCode(OP_F64_I64, "f64.i64",
		[]*CXArgument{newOpPar(TYPE_F64, false)},
		[]*CXArgument{newOpPar(TYPE_I64, false)})
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
	AddOpCode(OP_F64_MOD, "f32.mod",
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
	AddOpCode(OP_F64_ACOS, "f64.acos",
		[]*CXArgument{newOpPar(TYPE_F64, false)},
		[]*CXArgument{newOpPar(TYPE_F64, false)})
	AddOpCode(OP_F64_COS, "f64.cos",
		[]*CXArgument{newOpPar(TYPE_F64, false)},
		[]*CXArgument{newOpPar(TYPE_F64, false)})
	AddOpCode(OP_F64_ASIN, "f64.asin",
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
	AddOpCode(OP_F64_RAND, "f64.rand",
		[]*CXArgument{},
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
	AddOpCode(OP_STR_I8, "str.i8",
		[]*CXArgument{newOpPar(TYPE_STR, false)},
		[]*CXArgument{newOpPar(TYPE_I8, false)})
	AddOpCode(OP_STR_I16, "str.i16",
		[]*CXArgument{newOpPar(TYPE_STR, false)},
		[]*CXArgument{newOpPar(TYPE_I16, false)})
	AddOpCode(OP_STR_I32, "str.i32",
		[]*CXArgument{newOpPar(TYPE_STR, false)},
		[]*CXArgument{newOpPar(TYPE_I32, false)})
	AddOpCode(OP_STR_I64, "str.i64",
		[]*CXArgument{newOpPar(TYPE_STR, false)},
		[]*CXArgument{newOpPar(TYPE_I64, false)})
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
	handleOpcode := func(opCode int) opcodeHandler {
		switch opCode {
		case OP_IDENTITY:
			return opIdentity
		case OP_JMP:
			return opJmp
		case OP_DEBUG:
			return func(prgrm *CXProgram) {
				prgrm.PrintStack()
			}

		case OP_SERIALIZE:
			return opSerialize
		case OP_DESERIALIZE:
			return opDeserialize

		case OP_UND_EQUAL:
			return opEqual
		case OP_UND_UNEQUAL:
			return opUnequal
		case OP_UND_BITAND:
			return opBitand
		case OP_UND_BITXOR:
			return opBitxor
		case OP_UND_BITOR:
			return opBitor
		case OP_UND_BITCLEAR:
			return opBitclear
		case OP_UND_MUL:
			return opMul
		case OP_UND_DIV:
			return opDiv
		case OP_UND_MOD:
			return opMod
		case OP_UND_ADD:
			return opAdd
		case OP_UND_SUB:
			return opSub
		case OP_UND_NEG:
			return opNeg
		case OP_UND_BITSHL:
			return opBitshl
		case OP_UND_BITSHR:
			return opBitshr
		case OP_UND_LT:
			return opLt
		case OP_UND_GT:
			return opGt
		case OP_UND_LTEQ:
			return opLteq
		case OP_UND_GTEQ:
			return opGteq
		case OP_UND_LEN:
			return opLen
		case OP_UND_PRINTF:
			return opPrintf
		case OP_UND_SPRINTF:
			return opSprintf
		case OP_UND_READ:
			return opRead

		case OP_BOOL_PRINT:
			return opBoolPrint
		case OP_BOOL_EQUAL:
			return opBoolEqual
		case OP_BOOL_UNEQUAL:
			return opBoolUnequal
		case OP_BOOL_NOT:
			return opBoolNot
		case OP_BOOL_OR:
			return opBoolOr
		case OP_BOOL_AND:
			return opBoolAnd

		case OP_I8_STR:
			return opI8ToStr
		case OP_I8_I16:
			return opI8ToI16
		case OP_I8_I32:
			return opI8ToI32
		case OP_I8_I64:
			return opI8ToI64
		case OP_I8_UI8:
			return opI8ToUI8
		case OP_I8_UI16:
			return opI8ToUI16
		case OP_I8_UI32:
			return opI8ToUI32
		case OP_I8_UI64:
			return opI8ToUI64
		case OP_I8_F32:
			return opI8ToF32
		case OP_I8_F64:
			return opI8ToF64
		case OP_I8_PRINT:
			return opI8Print
		case OP_I8_ADD:
			return opI8Add
		case OP_I8_SUB:
			return opI8Sub
		case OP_I8_NEG:
			return opI8Neg
		case OP_I8_MUL:
			return opI8Mul
		case OP_I8_DIV:
			return opI8Div
		case OP_I8_MOD:
			return opI8Mod
		case OP_I8_ABS:
			return opI8Abs
		case OP_I8_GT:
			return opI8Gt
		case OP_I8_GTEQ:
			return opI8Gteq
		case OP_I8_LT:
			return opI8Lt
		case OP_I8_LTEQ:
			return opI8Lteq
		case OP_I8_EQ:
			return opI8Eq
		case OP_I8_UNEQ:
			return opI8Uneq
		case OP_I8_BITAND:
			return opI8Bitand
		case OP_I8_BITOR:
			return opI8Bitor
		case OP_I8_BITXOR:
			return opI8Bitxor
		case OP_I8_BITCLEAR:
			return opI8Bitclear
		case OP_I8_BITSHL:
			return opI8Bitshl
		case OP_I8_BITSHR:
			return opI8Bitshr
		case OP_I8_MAX:
			return opI8Max
		case OP_I8_MIN:
			return opI8Min
		case OP_I8_RAND:
			return opI8Rand

		case OP_I16_STR:
			return opI16ToStr
		case OP_I16_I8:
			return opI16ToI8
		case OP_I16_I32:
			return opI16ToI32
		case OP_I16_I64:
			return opI16ToI64
		case OP_I16_UI8:
			return opI16ToUI8
		case OP_I16_UI16:
			return opI16ToUI16
		case OP_I16_UI32:
			return opI16ToUI32
		case OP_I16_UI64:
			return opI16ToUI64
		case OP_I16_F32:
			return opI16ToF32
		case OP_I16_F64:
			return opI16ToF64
		case OP_I16_PRINT:
			return opI16Print
		case OP_I16_ADD:
			return opI16Add
		case OP_I16_SUB:
			return opI16Sub
		case OP_I16_NEG:
			return opI16Neg
		case OP_I16_MUL:
			return opI16Mul
		case OP_I16_DIV:
			return opI16Div
		case OP_I16_MOD:
			return opI16Mod
		case OP_I16_ABS:
			return opI16Abs
		case OP_I16_GT:
			return opI16Gt
		case OP_I16_GTEQ:
			return opI16Gteq
		case OP_I16_LT:
			return opI16Lt
		case OP_I16_LTEQ:
			return opI16Lteq
		case OP_I16_EQ:
			return opI16Eq
		case OP_I16_UNEQ:
			return opI16Uneq
		case OP_I16_BITAND:
			return opI16Bitand
		case OP_I16_BITOR:
			return opI16Bitor
		case OP_I16_BITXOR:
			return opI16Bitxor
		case OP_I16_BITCLEAR:
			return opI16Bitclear
		case OP_I16_BITSHL:
			return opI16Bitshl
		case OP_I16_BITSHR:
			return opI16Bitshr
		case OP_I16_MAX:
			return opI16Max
		case OP_I16_MIN:
			return opI16Min
		case OP_I16_RAND:
			return opI16Rand

		case OP_I32_STR:
			return opI32ToStr
		case OP_I32_I8:
			return opI32ToI8
		case OP_I32_I16:
			return opI32ToI16
		case OP_I32_I64:
			return opI32ToI64
		case OP_I32_UI8:
			return opI32ToUI8
		case OP_I32_UI16:
			return opI32ToUI16
		case OP_I32_UI32:
			return opI32ToUI32
		case OP_I32_UI64:
			return opI32ToUI64
		case OP_I32_F32:
			return opI32ToF32
		case OP_I32_F64:
			return opI32ToF64
		case OP_I32_PRINT:
			return opI32Print
		case OP_I32_ADD:
			return opI32Add
		case OP_I32_SUB:
			return opI32Sub
		case OP_I32_NEG:
			return opI32Neg
		case OP_I32_MUL:
			return opI32Mul
		case OP_I32_DIV:
			return opI32Div
		case OP_I32_MOD:
			return opI32Mod
		case OP_I32_ABS:
			return opI32Abs
		case OP_I32_GT:
			return opI32Gt
		case OP_I32_GTEQ:
			return opI32Gteq
		case OP_I32_LT:
			return opI32Lt
		case OP_I32_LTEQ:
			return opI32Lteq
		case OP_I32_EQ:
			return opI32Eq
		case OP_I32_UNEQ:
			return opI32Uneq
		case OP_I32_BITAND:
			return opI32Bitand
		case OP_I32_BITOR:
			return opI32Bitor
		case OP_I32_BITXOR:
			return opI32Bitxor
		case OP_I32_BITCLEAR:
			return opI32Bitclear
		case OP_I32_BITSHL:
			return opI32Bitshl
		case OP_I32_BITSHR:
			return opI32Bitshr
		case OP_I32_MAX:
			return opI32Max
		case OP_I32_MIN:
			return opI32Min
		case OP_I32_RAND:
			return opI32Rand

		case OP_I64_STR:
			return opI64ToStr
		case OP_I64_I8:
			return opI64ToI8
		case OP_I64_I16:
			return opI64ToI16
		case OP_I64_I32:
			return opI64ToI32
		case OP_I64_UI8:
			return opI64ToUI8
		case OP_I64_UI16:
			return opI64ToUI16
		case OP_I64_UI32:
			return opI64ToUI32
		case OP_I64_UI64:
			return opI64ToUI64
		case OP_I64_F32:
			return opI64ToF32
		case OP_I64_F64:
			return opI64ToF64
		case OP_I64_PRINT:
			return opI64Print
		case OP_I64_ADD:
			return opI64Add
		case OP_I64_SUB:
			return opI64Sub
		case OP_I64_NEG:
			return opI64Neg
		case OP_I64_MUL:
			return opI64Mul
		case OP_I64_DIV:
			return opI64Div
		case OP_I64_MOD:
			return opI64Mod
		case OP_I64_ABS:
			return opI64Abs
		case OP_I64_GT:
			return opI64Gt
		case OP_I64_GTEQ:
			return opI64Gteq
		case OP_I64_LT:
			return opI64Lt
		case OP_I64_LTEQ:
			return opI64Lteq
		case OP_I64_EQ:
			return opI64Eq
		case OP_I64_UNEQ:
			return opI64Uneq
		case OP_I64_BITAND:
			return opI64Bitand
		case OP_I64_BITOR:
			return opI64Bitor
		case OP_I64_BITXOR:
			return opI64Bitxor
		case OP_I64_BITCLEAR:
			return opI64Bitclear
		case OP_I64_BITSHL:
			return opI64Bitshl
		case OP_I64_BITSHR:
			return opI64Bitshr
		case OP_I64_MAX:
			return opI64Max
		case OP_I64_MIN:
			return opI64Min
		case OP_I64_RAND:
			return opI64Rand

		case OP_UI8_STR:
			return opUI8ToStr
		case OP_UI8_I8:
			return opUI8ToI8
		case OP_UI8_I16:
			return opUI8ToI16
		case OP_UI8_I32:
			return opUI8ToI32
		case OP_UI8_I64:
			return opUI8ToI64
		case OP_UI8_UI16:
			return opUI8ToUI16
		case OP_UI8_UI32:
			return opUI8ToUI32
		case OP_UI8_UI64:
			return opUI8ToUI64
		case OP_UI8_F32:
			return opUI8ToF32
		case OP_UI8_F64:
			return opUI8ToF64
		case OP_UI8_PRINT:
			return opUI8Print
		case OP_UI8_ADD:
			return opUI8Add
		case OP_UI8_SUB:
			return opUI8Sub
		case OP_UI8_MUL:
			return opUI8Mul
		case OP_UI8_DIV:
			return opUI8Div
		case OP_UI8_MOD:
			return opUI8Mod
		case OP_UI8_GT:
			return opUI8Gt
		case OP_UI8_GTEQ:
			return opUI8Gteq
		case OP_UI8_LT:
			return opUI8Lt
		case OP_UI8_LTEQ:
			return opUI8Lteq
		case OP_UI8_EQ:
			return opUI8Eq
		case OP_UI8_UNEQ:
			return opUI8Uneq
		case OP_UI8_BITAND:
			return opUI8Bitand
		case OP_UI8_BITOR:
			return opUI8Bitor
		case OP_UI8_BITXOR:
			return opUI8Bitxor
		case OP_UI8_BITCLEAR:
			return opUI8Bitclear
		case OP_UI8_BITSHL:
			return opUI8Bitshl
		case OP_UI8_BITSHR:
			return opUI8Bitshr
		case OP_UI8_MAX:
			return opUI8Max
		case OP_UI8_MIN:
			return opUI8Min
		case OP_UI8_RAND:
			return opUI8Rand

		case OP_UI16_STR:
			return opUI16ToStr
		case OP_UI16_I8:
			return opUI16ToI8
		case OP_UI16_I16:
			return opUI16ToI16
		case OP_UI16_I32:
			return opUI16ToI32
		case OP_UI16_I64:
			return opUI16ToI64
		case OP_UI16_UI8:
			return opUI16ToUI8
		case OP_UI16_UI32:
			return opUI16ToUI32
		case OP_UI16_UI64:
			return opUI16ToUI64
		case OP_UI16_F32:
			return opUI16ToF32
		case OP_UI16_F64:
			return opUI16ToF64
		case OP_UI16_PRINT:
			return opUI16Print
		case OP_UI16_ADD:
			return opUI16Add
		case OP_UI16_SUB:
			return opUI16Sub
		case OP_UI16_MUL:
			return opUI16Mul
		case OP_UI16_DIV:
			return opUI16Div
		case OP_UI16_MOD:
			return opUI16Mod
		case OP_UI16_GT:
			return opUI16Gt
		case OP_UI16_GTEQ:
			return opUI16Gteq
		case OP_UI16_LT:
			return opUI16Lt
		case OP_UI16_LTEQ:
			return opUI16Lteq
		case OP_UI16_EQ:
			return opUI16Eq
		case OP_UI16_UNEQ:
			return opUI16Uneq
		case OP_UI16_BITAND:
			return opUI16Bitand
		case OP_UI16_BITOR:
			return opUI16Bitor
		case OP_UI16_BITXOR:
			return opUI16Bitxor
		case OP_UI16_BITCLEAR:
			return opUI16Bitclear
		case OP_UI16_BITSHL:
			return opUI16Bitshl
		case OP_UI16_BITSHR:
			return opUI16Bitshr
		case OP_UI16_MAX:
			return opUI16Max
		case OP_UI16_MIN:
			return opUI16Min
		case OP_UI16_RAND:
			return opUI16Rand

		case OP_UI32_STR:
			return opUI32ToStr
		case OP_UI32_I8:
			return opUI32ToI8
		case OP_UI32_I16:
			return opUI32ToI16
		case OP_UI32_I32:
			return opUI32ToI32
		case OP_UI32_I64:
			return opUI32ToI64
		case OP_UI32_UI8:
			return opUI32ToUI8
		case OP_UI32_UI16:
			return opUI32ToUI16
		case OP_UI32_UI64:
			return opUI32ToUI64
		case OP_UI32_F32:
			return opUI32ToF32
		case OP_UI32_F64:
			return opUI32ToF64
		case OP_UI32_PRINT:
			return opUI32Print
		case OP_UI32_ADD:
			return opUI32Add
		case OP_UI32_SUB:
			return opUI32Sub
		case OP_UI32_MUL:
			return opUI32Mul
		case OP_UI32_DIV:
			return opUI32Div
		case OP_UI32_MOD:
			return opUI32Mod
		case OP_UI32_GT:
			return opUI32Gt
		case OP_UI32_GTEQ:
			return opUI32Gteq
		case OP_UI32_LT:
			return opUI32Lt
		case OP_UI32_LTEQ:
			return opUI32Lteq
		case OP_UI32_EQ:
			return opUI32Eq
		case OP_UI32_UNEQ:
			return opUI32Uneq
		case OP_UI32_BITAND:
			return opUI32Bitand
		case OP_UI32_BITOR:
			return opUI32Bitor
		case OP_UI32_BITXOR:
			return opUI32Bitxor
		case OP_UI32_BITCLEAR:
			return opUI32Bitclear
		case OP_UI32_BITSHL:
			return opUI32Bitshl
		case OP_UI32_BITSHR:
			return opUI32Bitshr
		case OP_UI32_MAX:
			return opUI32Max
		case OP_UI32_MIN:
			return opUI32Min
		case OP_UI32_RAND:
			return opUI32Rand

		case OP_UI64_STR:
			return opUI64ToStr
		case OP_UI64_I8:
			return opUI64ToI8
		case OP_UI64_I16:
			return opUI64ToI16
		case OP_UI64_I32:
			return opUI64ToI32
		case OP_UI64_I64:
			return opUI64ToI64
		case OP_UI64_UI8:
			return opUI64ToUI8
		case OP_UI64_UI16:
			return opUI64ToUI16
		case OP_UI64_UI32:
			return opUI64ToUI32
		case OP_UI64_F32:
			return opUI64ToF32
		case OP_UI64_F64:
			return opUI64ToF64
		case OP_UI64_PRINT:
			return opUI64Print
		case OP_UI64_ADD:
			return opUI64Add
		case OP_UI64_SUB:
			return opUI64Sub
		case OP_UI64_MUL:
			return opUI64Mul
		case OP_UI64_DIV:
			return opUI64Div
		case OP_UI64_MOD:
			return opUI64Mod
		case OP_UI64_GT:
			return opUI64Gt
		case OP_UI64_GTEQ:
			return opUI64Gteq
		case OP_UI64_LT:
			return opUI64Lt
		case OP_UI64_LTEQ:
			return opUI64Lteq
		case OP_UI64_EQ:
			return opUI64Eq
		case OP_UI64_UNEQ:
			return opUI64Uneq
		case OP_UI64_BITAND:
			return opUI64Bitand
		case OP_UI64_BITOR:
			return opUI64Bitor
		case OP_UI64_BITXOR:
			return opUI64Bitxor
		case OP_UI64_BITCLEAR:
			return opUI64Bitclear
		case OP_UI64_BITSHL:
			return opUI64Bitshl
		case OP_UI64_BITSHR:
			return opUI64Bitshr
		case OP_UI64_MAX:
			return opUI64Max
		case OP_UI64_MIN:
			return opUI64Min
		case OP_UI64_RAND:
			return opUI64Rand

		case OP_F32_IS_NAN:
			return opF32Isnan
		case OP_F32_STR:
			return opF32ToStr
		case OP_F32_I8:
			return opF32ToI8
		case OP_F32_I16:
			return opF32ToI16
		case OP_F32_I32:
			return opF32ToI32
		case OP_F32_I64:
			return opF32ToI64
		case OP_F32_UI8:
			return opF32ToUI8
		case OP_F32_UI16:
			return opF32ToUI16
		case OP_F32_UI32:
			return opF32ToUI32
		case OP_F32_UI64:
			return opF32ToUI64
		case OP_F32_F64:
			return opF32ToF64
		case OP_F32_PRINT:
			return opF32Print
		case OP_F32_ADD:
			return opF32Add
		case OP_F32_SUB:
			return opF32Sub
		case OP_F32_NEG:
			return opF32Neg
		case OP_F32_MUL:
			return opF32Mul
		case OP_F32_DIV:
			return opF32Div
		case OP_F32_MOD:
			return opF32Mod
		case OP_F32_ABS:
			return opF32Abs
		case OP_F32_POW:
			return opF32Pow
		case OP_F32_GT:
			return opF32Gt
		case OP_F32_GTEQ:
			return opF32Gteq
		case OP_F32_LT:
			return opF32Lt
		case OP_F32_LTEQ:
			return opF32Lteq
		case OP_F32_EQ:
			return opF32Eq
		case OP_F32_UNEQ:
			return opF32Uneq
		case OP_F32_ACOS:
			return opF32Acos
		case OP_F32_COS:
			return opF32Cos
		case OP_F32_ASIN:
			return opF32Asin
		case OP_F32_SIN:
			return opF32Sin
		case OP_F32_SQRT:
			return opF32Sqrt
		case OP_F32_LOG:
			return opF32Log
		case OP_F32_LOG2:
			return opF32Log2
		case OP_F32_LOG10:
			return opF32Log10
		case OP_F32_MAX:
			return opF32Max
		case OP_F32_MIN:
			return opF32Min
		case OP_F32_RAND:
			return opF32Rand

		case OP_F64_IS_NAN:
			return opF64Isnan
		case OP_F64_STR:
			return opF64ToStr
		case OP_F64_I8:
			return opF64ToI8
		case OP_F64_I16:
			return opF64ToI16
		case OP_F64_I32:
			return opF64ToI32
		case OP_F64_I64:
			return opF64ToI64
		case OP_F64_UI8:
			return opF64ToUI8
		case OP_F64_UI16:
			return opF64ToUI16
		case OP_F64_UI32:
			return opF64ToUI32
		case OP_F64_UI64:
			return opF64ToUI64
		case OP_F64_F32:
			return opF64ToF32
		case OP_F64_PRINT:
			return opF64Print
		case OP_F64_ADD:
			return opF64Add
		case OP_F64_SUB:
			return opF64Sub
		case OP_F64_NEG:
			return opF64Neg
		case OP_F64_MUL:
			return opF64Mul
		case OP_F64_DIV:
			return opF64Div
		case OP_F64_MOD:
			return opF64Mod
		case OP_F64_ABS:
			return opF64Abs
		case OP_F64_POW:
			return opF64Pow
		case OP_F64_GT:
			return opF64Gt
		case OP_F64_GTEQ:
			return opF64Gteq
		case OP_F64_LT:
			return opF64Lt
		case OP_F64_LTEQ:
			return opF64Lteq
		case OP_F64_EQ:
			return opF64Eq
		case OP_F64_UNEQ:
			return opF64Uneq
		case OP_F64_ACOS:
			return opF64Acos
		case OP_F64_COS:
			return opF64Cos
		case OP_F64_ASIN:
			return opF64Asin
		case OP_F64_SIN:
			return opF64Sin
		case OP_F64_SQRT:
			return opF64Sqrt
		case OP_F64_LOG:
			return opF64Log
		case OP_F64_LOG2:
			return opF64Log2
		case OP_F64_LOG10:
			return opF64Log10
		case OP_F64_MAX:
			return opF64Max
		case OP_F64_MIN:
			return opF64Min
		case OP_F64_RAND:
			return opF64Rand

		case OP_STR_I8:
			return opStrToI8
		case OP_STR_I16:
			return opStrToI16
		case OP_STR_I32:
			return opStrToI32
		case OP_STR_I64:
			return opStrToI64
		case OP_STR_UI8:
			return opStrToUI8
		case OP_STR_UI16:
			return opStrToUI16
		case OP_STR_UI32:
			return opStrToUI32
		case OP_STR_UI64:
			return opStrToUI64
		case OP_STR_F32:
			return opStrToF32
		case OP_STR_F64:
			return opStrToF64
		case OP_STR_PRINT:
			return opStrPrint
		case OP_STR_EQ:
			return opStrEq
		case OP_STR_CONCAT:
			return opStrConcat
		case OP_STR_SUBSTR:
			return opStrSubstr
		case OP_STR_INDEX:
			return opStrIndex
		case OP_STR_TRIM_SPACE:
			return opStrTrimSpace

		case OP_APPEND:
			return opAppend
		case OP_RESIZE:
			return opResize
		case OP_INSERT:
			return opInsert
		case OP_REMOVE:
			return opRemove
		case OP_COPY:
			return opCopy

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
			return opAssertValue
		case OP_TEST:
			return opTest
		case OP_PANIC:
			return opPanic
		case OP_STRERROR:
			return opStrError

		// affordances
		case OP_AFF_PRINT:
			return opAffPrint
		case OP_AFF_QUERY:
			return opAffQuery
		case OP_AFF_ON:
			return opAffOn
		case OP_AFF_OF:
			return opAffOf
		case OP_AFF_INFORM:
			return opAffInform
		case OP_AFF_REQUEST:
			return opAffRequest
		}

		return nil
	}

	opcodeHandlerFinders = append(opcodeHandlerFinders, handleOpcode)
}
