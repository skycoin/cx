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

	OP_BYTE_BYTE
	OP_BYTE_STR
	OP_BYTE_I32
	OP_BYTE_I64
	OP_BYTE_F32
	OP_BYTE_F64

	OP_I32_BYTE
	OP_I32_STR
	OP_I32_I32
	OP_I32_I64
	OP_I32_F32
	OP_I32_F64

	OP_I64_BYTE
	OP_I64_STR
	OP_I64_I32
	OP_I64_I64
	OP_I64_F32
	OP_I64_F64

	OP_F32_BYTE
	OP_F32_STR
	OP_F32_I32
	OP_F32_I64
	OP_F32_F32
	OP_F32_F64

	OP_F64_BYTE
	OP_F64_STR
	OP_F64_I32
	OP_F64_I64
	OP_F64_F32
	OP_F64_F64

	OP_STR_BYTE
	OP_STR_STR
	OP_STR_I32
	OP_STR_I64
	OP_STR_F32
	OP_STR_F64

	END_PARSE_OPS

	OP_BOOL_PRINT

	OP_BOOL_EQUAL
	OP_BOOL_UNEQUAL
	OP_BOOL_NOT
	OP_BOOL_OR
	OP_BOOL_AND

	OP_BYTE_PRINT

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

	OP_F32_IS_NAN
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
	OP_F32_RAND
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

	OP_F64_IS_NAN
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
	OP_F64_RAND
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

	// http
	OP_HTTP_SERVE
	OP_HTTP_LISTEN_AND_SERVE
	OP_HTTP_NEW_REQUEST
	OP_HTTP_DO
	OP_HTTP_HANDLE
	OP_HTTP_CLOSE

	// dmsg
	OP_DMSG_DO
	
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

type opParam struct {
	typCode int        // The type code of the parameter.
	isSlice bool       // Is the parameter a slice.
	pkg     *CXPackage // To what package does this param belongs to.
	inputs []*CXArgument // Input parameters to a TYPE_FUNC parameter.
	outputs []*CXArgument // Output parameters to a TYPE_FUNC parameter.
}

// Helper function for creating parameters for standard library operators.
// The current standard library only uses basic types and slices. If more options are needed, modify this function
// func newOpPar(typCode int, isSlice bool) *CXArgument {
func newOpPar(paramData opParam) *CXArgument {
	arg := MakeArgument("", "", -1).AddType(TypeNames[paramData.typCode])
	arg.IsLocalDeclaration = true
	arg.Inputs = paramData.inputs
	arg.Outputs = paramData.outputs
	arg.Package = paramData.pkg
	if paramData.isSlice {
		arg.IsSlice = true
		arg.DeclarationSpecifiers = append(arg.DeclarationSpecifiers, DECL_SLICE)
	}
	return arg
}

var opParamCustomNotSlice = opParam{typCode: TYPE_CUSTOM, isSlice: false}
var opParamUndNotSlice = opParam{typCode: TYPE_UNDEFINED, isSlice: false}
var opParamFuncNotSlice = opParam{typCode: TYPE_FUNC, isSlice: false}
var opParamBoolNotSlice = opParam{typCode: TYPE_BOOL, isSlice: false}
var opParamAffNotSlice = opParam{typCode: TYPE_AFF, isSlice: false}
var opParamByteNotSlice = opParam{typCode: TYPE_BYTE, isSlice: false}
var opParamI32NotSlice = opParam{typCode: TYPE_I32, isSlice: false}
var opParamI64NotSlice = opParam{typCode: TYPE_I64, isSlice: false}
var opParamF32NotSlice = opParam{typCode: TYPE_F32, isSlice: false}
var opParamF64NotSlice = opParam{typCode: TYPE_F64, isSlice: false}
var opParamStrNotSlice = opParam{typCode: TYPE_STR, isSlice: false}

var opParamUndSlice = opParam{typCode: TYPE_UNDEFINED, isSlice: true}
var opParamF32Slice = opParam{typCode: TYPE_F32, isSlice: true}

func init() {
	httpPkg, err := PROGRAM.GetPackage("http")
	if err != nil {
		panic(err)
	}

	AddOpCode(OP_IDENTITY, "identity",
		[]*CXArgument{newOpPar(opParamUndNotSlice)},
		[]*CXArgument{newOpPar(opParamUndNotSlice)})
	AddOpCode(OP_JMP, "jmp",
		[]*CXArgument{newOpPar(opParamBoolNotSlice)},
		[]*CXArgument{}) // newOpPar(opParamUndNotSlice) to allow 0 inputs (goto)
	AddOpCode(OP_DEBUG, "debug",
		[]*CXArgument{},
		[]*CXArgument{})

	AddOpCode(OP_SERIALIZE, "serialize",
		[]*CXArgument{newOpPar(opParamAffNotSlice)},
		[]*CXArgument{newOpPar(opParamByteNotSlice)})
	AddOpCode(OP_DESERIALIZE, "deserialize",
		[]*CXArgument{newOpPar(opParamByteNotSlice)},
		[]*CXArgument{})

	AddOpCode(OP_UND_EQUAL, "eq",
		[]*CXArgument{newOpPar(opParamUndNotSlice), newOpPar(opParamUndNotSlice)},
		[]*CXArgument{newOpPar(opParamBoolNotSlice)})
	AddOpCode(OP_UND_UNEQUAL, "uneq",
		[]*CXArgument{newOpPar(opParamUndNotSlice), newOpPar(opParamUndNotSlice)},
		[]*CXArgument{newOpPar(opParamBoolNotSlice)})
	AddOpCode(OP_UND_BITAND, "bitand",
		[]*CXArgument{newOpPar(opParamUndNotSlice), newOpPar(opParamUndNotSlice)},
		[]*CXArgument{newOpPar(opParamUndNotSlice)})
	AddOpCode(OP_UND_BITXOR, "bitxor",
		[]*CXArgument{newOpPar(opParamUndNotSlice), newOpPar(opParamUndNotSlice)},
		[]*CXArgument{newOpPar(opParamUndNotSlice)})
	AddOpCode(OP_UND_BITOR, "bitor",
		[]*CXArgument{newOpPar(opParamUndNotSlice), newOpPar(opParamUndNotSlice)},
		[]*CXArgument{newOpPar(opParamUndNotSlice)})
	AddOpCode(OP_UND_BITCLEAR, "bitclear",
		[]*CXArgument{newOpPar(opParamUndNotSlice), newOpPar(opParamUndNotSlice)},
		[]*CXArgument{newOpPar(opParamUndNotSlice)})
	AddOpCode(OP_UND_MUL, "mul",
		[]*CXArgument{newOpPar(opParamUndNotSlice), newOpPar(opParamUndNotSlice)},
		[]*CXArgument{newOpPar(opParamUndNotSlice)})
	AddOpCode(OP_UND_DIV, "div",
		[]*CXArgument{newOpPar(opParamUndNotSlice), newOpPar(opParamUndNotSlice)},
		[]*CXArgument{newOpPar(opParamUndNotSlice)})
	AddOpCode(OP_UND_MOD, "mod",
		[]*CXArgument{newOpPar(opParamUndNotSlice), newOpPar(opParamUndNotSlice)},
		[]*CXArgument{newOpPar(opParamUndNotSlice)})
	AddOpCode(OP_UND_ADD, "add",
		[]*CXArgument{newOpPar(opParamUndNotSlice), newOpPar(opParamUndNotSlice)},
		[]*CXArgument{newOpPar(opParamUndNotSlice)})
	AddOpCode(OP_UND_SUB, "sub",
		[]*CXArgument{newOpPar(opParamUndNotSlice), newOpPar(opParamUndNotSlice)},
		[]*CXArgument{newOpPar(opParamUndNotSlice)})
	AddOpCode(OP_UND_NEG, "neg",
		[]*CXArgument{newOpPar(opParamUndNotSlice)},
		[]*CXArgument{newOpPar(opParamUndNotSlice)})
	AddOpCode(OP_UND_BITSHL, "bitshl",
		[]*CXArgument{newOpPar(opParamUndNotSlice), newOpPar(opParamUndNotSlice)},
		[]*CXArgument{newOpPar(opParamUndNotSlice)})
	AddOpCode(OP_UND_BITSHR, "bitshr",
		[]*CXArgument{newOpPar(opParamUndNotSlice), newOpPar(opParamUndNotSlice)},
		[]*CXArgument{newOpPar(opParamUndNotSlice)})
	AddOpCode(OP_UND_LT, "lt",
		[]*CXArgument{newOpPar(opParamUndNotSlice), newOpPar(opParamUndNotSlice)},
		[]*CXArgument{newOpPar(opParamBoolNotSlice)})
	AddOpCode(OP_UND_GT, "gt",
		[]*CXArgument{newOpPar(opParamUndNotSlice), newOpPar(opParamUndNotSlice)},
		[]*CXArgument{newOpPar(opParamBoolNotSlice)})
	AddOpCode(OP_UND_LTEQ, "lteq",
		[]*CXArgument{newOpPar(opParamUndNotSlice), newOpPar(opParamUndNotSlice)},
		[]*CXArgument{newOpPar(opParamBoolNotSlice)})
	AddOpCode(OP_UND_GTEQ, "gteq",
		[]*CXArgument{newOpPar(opParamUndNotSlice), newOpPar(opParamUndNotSlice)},
		[]*CXArgument{newOpPar(opParamBoolNotSlice)})
	AddOpCode(OP_UND_LEN, "len",
		[]*CXArgument{newOpPar(opParamUndNotSlice)},
		[]*CXArgument{newOpPar(opParamI32NotSlice)})
	AddOpCode(OP_UND_PRINTF, "printf",
		[]*CXArgument{newOpPar(opParamUndNotSlice)},
		[]*CXArgument{})
	AddOpCode(OP_UND_SPRINTF, "sprintf",
		[]*CXArgument{newOpPar(opParamUndNotSlice)},
		[]*CXArgument{newOpPar(opParamStrNotSlice)})
	AddOpCode(OP_UND_READ, "read",
		[]*CXArgument{},
		[]*CXArgument{newOpPar(opParamStrNotSlice)})

	AddOpCode(OP_BYTE_BYTE, "byte.byte",
		[]*CXArgument{newOpPar(opParamByteNotSlice)},
		[]*CXArgument{newOpPar(opParamByteNotSlice)})
	AddOpCode(OP_BYTE_STR, "byte.str",
		[]*CXArgument{newOpPar(opParamByteNotSlice)},
		[]*CXArgument{newOpPar(opParamStrNotSlice)})
	AddOpCode(OP_BYTE_I32, "byte.i32",
		[]*CXArgument{newOpPar(opParamByteNotSlice)},
		[]*CXArgument{newOpPar(opParamI32NotSlice)})
	AddOpCode(OP_BYTE_I64, "byte.i64",
		[]*CXArgument{newOpPar(opParamByteNotSlice)},
		[]*CXArgument{newOpPar(opParamI64NotSlice)})
	AddOpCode(OP_BYTE_F32, "byte.f32",
		[]*CXArgument{newOpPar(opParamByteNotSlice)},
		[]*CXArgument{newOpPar(opParamF32NotSlice)})
	AddOpCode(OP_BYTE_F64, "byte.f64",
		[]*CXArgument{newOpPar(opParamByteNotSlice)},
		[]*CXArgument{newOpPar(opParamF64NotSlice)})

	AddOpCode(OP_BYTE_PRINT, "byte.print",
		[]*CXArgument{newOpPar(opParamByteNotSlice)},
		[]*CXArgument{})

	AddOpCode(OP_BOOL_PRINT, "bool.print",
		[]*CXArgument{newOpPar(opParamBoolNotSlice)},
		[]*CXArgument{})
	AddOpCode(OP_BOOL_EQUAL, "bool.eq",
		[]*CXArgument{newOpPar(opParamBoolNotSlice), newOpPar(opParamBoolNotSlice)},
		[]*CXArgument{newOpPar(opParamBoolNotSlice)})
	AddOpCode(OP_BOOL_UNEQUAL, "bool.uneq",
		[]*CXArgument{newOpPar(opParamBoolNotSlice), newOpPar(opParamBoolNotSlice)},
		[]*CXArgument{newOpPar(opParamBoolNotSlice)})
	AddOpCode(OP_BOOL_NOT, "bool.not",
		[]*CXArgument{newOpPar(opParamBoolNotSlice)},
		[]*CXArgument{newOpPar(opParamBoolNotSlice)})
	AddOpCode(OP_BOOL_OR, "bool.or",
		[]*CXArgument{newOpPar(opParamBoolNotSlice), newOpPar(opParamBoolNotSlice)},
		[]*CXArgument{newOpPar(opParamBoolNotSlice)})
	AddOpCode(OP_BOOL_AND, "bool.and",
		[]*CXArgument{newOpPar(opParamBoolNotSlice), newOpPar(opParamBoolNotSlice)},
		[]*CXArgument{newOpPar(opParamBoolNotSlice)})

	AddOpCode(OP_I32_BYTE, "i32.byte",
		[]*CXArgument{newOpPar(opParamI32NotSlice)},
		[]*CXArgument{newOpPar(opParamByteNotSlice)})
	AddOpCode(OP_I32_STR, "i32.str",
		[]*CXArgument{newOpPar(opParamI32NotSlice)},
		[]*CXArgument{newOpPar(opParamStrNotSlice)})
	AddOpCode(OP_I32_I32, "i32.i32",
		[]*CXArgument{newOpPar(opParamI32NotSlice)},
		[]*CXArgument{newOpPar(opParamI32NotSlice)})
	AddOpCode(OP_I32_I64, "i32.i64",
		[]*CXArgument{newOpPar(opParamI32NotSlice)},
		[]*CXArgument{newOpPar(opParamI64NotSlice)})
	AddOpCode(OP_I32_F32, "i32.f32",
		[]*CXArgument{newOpPar(opParamI32NotSlice)},
		[]*CXArgument{newOpPar(opParamF32NotSlice)})
	AddOpCode(OP_I32_F64, "i32.f64",
		[]*CXArgument{newOpPar(opParamI32NotSlice)},
		[]*CXArgument{newOpPar(opParamF64NotSlice)})

	AddOpCode(OP_I32_PRINT, "i32.print",
		[]*CXArgument{newOpPar(opParamI32NotSlice)},
		[]*CXArgument{})
	AddOpCode(OP_I32_ADD, "i32.add",
		[]*CXArgument{newOpPar(opParamI32NotSlice), newOpPar(opParamI32NotSlice)},
		[]*CXArgument{newOpPar(opParamI32NotSlice)})
	AddOpCode(OP_I32_SUB, "i32.sub",
		[]*CXArgument{newOpPar(opParamI32NotSlice), newOpPar(opParamI32NotSlice)},
		[]*CXArgument{newOpPar(opParamI32NotSlice)})
	AddOpCode(OP_I32_NEG, "i32.neg",
		[]*CXArgument{newOpPar(opParamI32NotSlice)},
		[]*CXArgument{newOpPar(opParamI32NotSlice)})
	AddOpCode(OP_I32_MUL, "i32.mul",
		[]*CXArgument{newOpPar(opParamI32NotSlice), newOpPar(opParamI32NotSlice)},
		[]*CXArgument{newOpPar(opParamI32NotSlice)})
	AddOpCode(OP_I32_DIV, "i32.div",
		[]*CXArgument{newOpPar(opParamI32NotSlice), newOpPar(opParamI32NotSlice)},
		[]*CXArgument{newOpPar(opParamI32NotSlice)})
	AddOpCode(OP_I32_ABS, "i32.abs",
		[]*CXArgument{newOpPar(opParamI32NotSlice)},
		[]*CXArgument{newOpPar(opParamI32NotSlice)})
	AddOpCode(OP_I32_POW, "i32.pow",
		[]*CXArgument{newOpPar(opParamI32NotSlice), newOpPar(opParamI32NotSlice)},
		[]*CXArgument{newOpPar(opParamI32NotSlice)})
	AddOpCode(OP_I32_GT, "i32.gt",
		[]*CXArgument{newOpPar(opParamI32NotSlice), newOpPar(opParamI32NotSlice)},
		[]*CXArgument{newOpPar(opParamBoolNotSlice)})
	AddOpCode(OP_I32_GTEQ, "i32.gteq",
		[]*CXArgument{newOpPar(opParamI32NotSlice), newOpPar(opParamI32NotSlice)},
		[]*CXArgument{newOpPar(opParamBoolNotSlice)})
	AddOpCode(OP_I32_LT, "i32.lt",
		[]*CXArgument{newOpPar(opParamI32NotSlice), newOpPar(opParamI32NotSlice)},
		[]*CXArgument{newOpPar(opParamBoolNotSlice)})
	AddOpCode(OP_I32_LTEQ, "i32.lteq",
		[]*CXArgument{newOpPar(opParamI32NotSlice), newOpPar(opParamI32NotSlice)},
		[]*CXArgument{newOpPar(opParamBoolNotSlice)})
	AddOpCode(OP_I32_EQ, "i32.eq",
		[]*CXArgument{newOpPar(opParamI32NotSlice), newOpPar(opParamI32NotSlice)},
		[]*CXArgument{newOpPar(opParamBoolNotSlice)})
	AddOpCode(OP_I32_UNEQ, "i32.uneq",
		[]*CXArgument{newOpPar(opParamI32NotSlice), newOpPar(opParamI32NotSlice)},
		[]*CXArgument{newOpPar(opParamBoolNotSlice)})
	AddOpCode(OP_I32_MOD, "i32.mod",
		[]*CXArgument{newOpPar(opParamI32NotSlice), newOpPar(opParamI32NotSlice)},
		[]*CXArgument{newOpPar(opParamI32NotSlice)})
	AddOpCode(OP_I32_RAND, "i32.rand",
		[]*CXArgument{newOpPar(opParamI32NotSlice), newOpPar(opParamI32NotSlice)},
		[]*CXArgument{newOpPar(opParamI32NotSlice)})
	AddOpCode(OP_I32_BITAND, "i32.bitand",
		[]*CXArgument{newOpPar(opParamI32NotSlice), newOpPar(opParamI32NotSlice)},
		[]*CXArgument{newOpPar(opParamI32NotSlice)})
	AddOpCode(OP_I32_BITOR, "i32.bitor",
		[]*CXArgument{newOpPar(opParamI32NotSlice), newOpPar(opParamI32NotSlice)},
		[]*CXArgument{newOpPar(opParamI32NotSlice)})
	AddOpCode(OP_I32_BITXOR, "i32.bitxor",
		[]*CXArgument{newOpPar(opParamI32NotSlice), newOpPar(opParamI32NotSlice)},
		[]*CXArgument{newOpPar(opParamI32NotSlice)})
	AddOpCode(OP_I32_BITCLEAR, "i32.bitclear",
		[]*CXArgument{newOpPar(opParamI32NotSlice), newOpPar(opParamI32NotSlice)},
		[]*CXArgument{newOpPar(opParamI32NotSlice)})
	AddOpCode(OP_I32_BITSHL, "i32.bitshl",
		[]*CXArgument{newOpPar(opParamI32NotSlice), newOpPar(opParamI32NotSlice)},
		[]*CXArgument{newOpPar(opParamI32NotSlice)})
	AddOpCode(OP_I32_BITSHR, "i32.bitshr",
		[]*CXArgument{newOpPar(opParamI32NotSlice), newOpPar(opParamI32NotSlice)},
		[]*CXArgument{newOpPar(opParamI32NotSlice)})
	AddOpCode(OP_I32_SQRT, "i32.sqrt",
		[]*CXArgument{newOpPar(opParamI32NotSlice)},
		[]*CXArgument{newOpPar(opParamI32NotSlice)})
	AddOpCode(OP_I32_LOG, "i32.log",
		[]*CXArgument{newOpPar(opParamI32NotSlice)},
		[]*CXArgument{newOpPar(opParamI32NotSlice)})
	AddOpCode(OP_I32_LOG2, "i32.log2",
		[]*CXArgument{newOpPar(opParamI32NotSlice)},
		[]*CXArgument{newOpPar(opParamI32NotSlice)})
	AddOpCode(OP_I32_LOG10, "i32.log10",
		[]*CXArgument{newOpPar(opParamI32NotSlice)},
		[]*CXArgument{newOpPar(opParamI32NotSlice)})
	AddOpCode(OP_I32_MAX, "i32.max",
		[]*CXArgument{newOpPar(opParamI32NotSlice), newOpPar(opParamI32NotSlice)},
		[]*CXArgument{newOpPar(opParamI32NotSlice)})
	AddOpCode(OP_I32_MIN, "i32.min",
		[]*CXArgument{newOpPar(opParamI32NotSlice), newOpPar(opParamI32NotSlice)},
		[]*CXArgument{newOpPar(opParamI32NotSlice)})

	AddOpCode(OP_I64_BYTE, "i64.byte",
		[]*CXArgument{newOpPar(opParamI64NotSlice)},
		[]*CXArgument{newOpPar(opParamByteNotSlice)})
	AddOpCode(OP_I64_STR, "i64.str",
		[]*CXArgument{newOpPar(opParamI64NotSlice)},
		[]*CXArgument{newOpPar(opParamStrNotSlice)})
	AddOpCode(OP_I64_I32, "i64.i32",
		[]*CXArgument{newOpPar(opParamI64NotSlice)},
		[]*CXArgument{newOpPar(opParamI32NotSlice)})
	AddOpCode(OP_I64_I64, "i64.i64",
		[]*CXArgument{newOpPar(opParamI64NotSlice)},
		[]*CXArgument{newOpPar(opParamI64NotSlice)})
	AddOpCode(OP_I64_F32, "i64.f32",
		[]*CXArgument{newOpPar(opParamI64NotSlice)},
		[]*CXArgument{newOpPar(opParamF32NotSlice)})
	AddOpCode(OP_I64_F64, "i64.f64",
		[]*CXArgument{newOpPar(opParamI64NotSlice)},
		[]*CXArgument{newOpPar(opParamF64NotSlice)})

	AddOpCode(OP_I64_PRINT, "i64.print",
		[]*CXArgument{newOpPar(opParamI64NotSlice)},
		[]*CXArgument{})
	AddOpCode(OP_I64_ADD, "i64.add",
		[]*CXArgument{newOpPar(opParamI64NotSlice), newOpPar(opParamI64NotSlice)},
		[]*CXArgument{newOpPar(opParamI64NotSlice)})
	AddOpCode(OP_I64_SUB, "i64.sub",
		[]*CXArgument{newOpPar(opParamI64NotSlice), newOpPar(opParamI64NotSlice)},
		[]*CXArgument{newOpPar(opParamI64NotSlice)})
	AddOpCode(OP_I64_NEG, "i64.neg",
		[]*CXArgument{newOpPar(opParamI64NotSlice)},
		[]*CXArgument{newOpPar(opParamI64NotSlice)})
	AddOpCode(OP_I64_MUL, "i64.mul",
		[]*CXArgument{newOpPar(opParamI64NotSlice), newOpPar(opParamI64NotSlice)},
		[]*CXArgument{newOpPar(opParamI64NotSlice)})
	AddOpCode(OP_I64_DIV, "i64.div",
		[]*CXArgument{newOpPar(opParamI64NotSlice), newOpPar(opParamI64NotSlice)},
		[]*CXArgument{newOpPar(opParamI64NotSlice)})
	AddOpCode(OP_I64_ABS, "i64.abs",
		[]*CXArgument{newOpPar(opParamI64NotSlice)},
		[]*CXArgument{newOpPar(opParamI64NotSlice)})
	AddOpCode(OP_I64_POW, "i64.pow",
		[]*CXArgument{newOpPar(opParamI64NotSlice), newOpPar(opParamI64NotSlice)},
		[]*CXArgument{newOpPar(opParamI64NotSlice)})
	AddOpCode(OP_I64_GT, "i64.gt",
		[]*CXArgument{newOpPar(opParamI64NotSlice), newOpPar(opParamI64NotSlice)},
		[]*CXArgument{newOpPar(opParamBoolNotSlice)})
	AddOpCode(OP_I64_GTEQ, "i64.gteq",
		[]*CXArgument{newOpPar(opParamI64NotSlice), newOpPar(opParamI64NotSlice)},
		[]*CXArgument{newOpPar(opParamBoolNotSlice)})
	AddOpCode(OP_I64_LT, "i64.lt",
		[]*CXArgument{newOpPar(opParamI64NotSlice), newOpPar(opParamI64NotSlice)},
		[]*CXArgument{newOpPar(opParamBoolNotSlice)})
	AddOpCode(OP_I64_LTEQ, "i64.lteq",
		[]*CXArgument{newOpPar(opParamI64NotSlice), newOpPar(opParamI64NotSlice)},
		[]*CXArgument{newOpPar(opParamBoolNotSlice)})
	AddOpCode(OP_I64_EQ, "i64.eq",
		[]*CXArgument{newOpPar(opParamI64NotSlice), newOpPar(opParamI64NotSlice)},
		[]*CXArgument{newOpPar(opParamBoolNotSlice)})
	AddOpCode(OP_I64_UNEQ, "i64.uneq",
		[]*CXArgument{newOpPar(opParamI64NotSlice), newOpPar(opParamI64NotSlice)},
		[]*CXArgument{newOpPar(opParamBoolNotSlice)})
	AddOpCode(OP_I64_MOD, "i64.mod",
		[]*CXArgument{newOpPar(opParamI64NotSlice), newOpPar(opParamI64NotSlice)},
		[]*CXArgument{newOpPar(opParamI64NotSlice)})
	AddOpCode(OP_I64_RAND, "i64.rand",
		[]*CXArgument{newOpPar(opParamI64NotSlice), newOpPar(opParamI64NotSlice)},
		[]*CXArgument{newOpPar(opParamI64NotSlice)})
	AddOpCode(OP_I64_BITAND, "i64.bitand",
		[]*CXArgument{newOpPar(opParamI64NotSlice), newOpPar(opParamI64NotSlice)},
		[]*CXArgument{newOpPar(opParamI64NotSlice)})
	AddOpCode(OP_I64_BITOR, "i64.bitor",
		[]*CXArgument{newOpPar(opParamI64NotSlice), newOpPar(opParamI64NotSlice)},
		[]*CXArgument{newOpPar(opParamI64NotSlice)})
	AddOpCode(OP_I64_BITXOR, "i64.bitxor",
		[]*CXArgument{newOpPar(opParamI64NotSlice), newOpPar(opParamI64NotSlice)},
		[]*CXArgument{newOpPar(opParamI64NotSlice)})
	AddOpCode(OP_I64_BITCLEAR, "i64.bitclear",
		[]*CXArgument{newOpPar(opParamI64NotSlice), newOpPar(opParamI64NotSlice)},
		[]*CXArgument{newOpPar(opParamI64NotSlice)})
	AddOpCode(OP_I64_BITSHL, "i64.bitshl",
		[]*CXArgument{newOpPar(opParamI64NotSlice), newOpPar(opParamI64NotSlice)},
		[]*CXArgument{newOpPar(opParamI64NotSlice)})
	AddOpCode(OP_I64_BITSHR, "i64.bitshr",
		[]*CXArgument{newOpPar(opParamI64NotSlice), newOpPar(opParamI64NotSlice)},
		[]*CXArgument{newOpPar(opParamI64NotSlice)})
	AddOpCode(OP_I64_SQRT, "i64.sqrt",
		[]*CXArgument{newOpPar(opParamI64NotSlice)},
		[]*CXArgument{newOpPar(opParamI64NotSlice)})
	AddOpCode(OP_I64_LOG, "i64.log",
		[]*CXArgument{newOpPar(opParamI64NotSlice)},
		[]*CXArgument{newOpPar(opParamI64NotSlice)})
	AddOpCode(OP_I64_LOG2, "i64.log2",
		[]*CXArgument{newOpPar(opParamI64NotSlice)},
		[]*CXArgument{newOpPar(opParamI64NotSlice)})
	AddOpCode(OP_I64_LOG10, "i64.log10",
		[]*CXArgument{newOpPar(opParamI64NotSlice)},
		[]*CXArgument{newOpPar(opParamI64NotSlice)})
	AddOpCode(OP_I64_MAX, "i64.max",
		[]*CXArgument{newOpPar(opParamI64NotSlice), newOpPar(opParamI64NotSlice)},
		[]*CXArgument{newOpPar(opParamI64NotSlice)})
	AddOpCode(OP_I64_MIN, "i64.min",
		[]*CXArgument{newOpPar(opParamI64NotSlice), newOpPar(opParamI64NotSlice)},
		[]*CXArgument{newOpPar(opParamI64NotSlice)})

	AddOpCode(OP_F32_IS_NAN, "f32.isnan",
		[]*CXArgument{newOpPar(opParamF32NotSlice)},
		[]*CXArgument{newOpPar(opParamBoolNotSlice)})
	AddOpCode(OP_F32_BYTE, "f32.byte",
		[]*CXArgument{newOpPar(opParamF32NotSlice)},
		[]*CXArgument{newOpPar(opParamByteNotSlice)})
	AddOpCode(OP_F32_STR, "f32.str",
		[]*CXArgument{newOpPar(opParamF32NotSlice)},
		[]*CXArgument{newOpPar(opParamStrNotSlice)})
	AddOpCode(OP_F32_I32, "f32.i32",
		[]*CXArgument{newOpPar(opParamF32NotSlice)},
		[]*CXArgument{newOpPar(opParamI32NotSlice)})
	AddOpCode(OP_F32_I64, "f32.i64",
		[]*CXArgument{newOpPar(opParamF32NotSlice)},
		[]*CXArgument{newOpPar(opParamI64NotSlice)})
	AddOpCode(OP_F32_F32, "f32.f32",
		[]*CXArgument{newOpPar(opParamF32NotSlice)},
		[]*CXArgument{newOpPar(opParamF32NotSlice)})
	AddOpCode(OP_F32_F64, "f32.f64",
		[]*CXArgument{newOpPar(opParamF32NotSlice)},
		[]*CXArgument{newOpPar(opParamF64NotSlice)})
	AddOpCode(OP_F32_PRINT, "f32.print",
		[]*CXArgument{newOpPar(opParamF32NotSlice)},
		[]*CXArgument{})
	AddOpCode(OP_F32_ADD, "f32.add",
		[]*CXArgument{newOpPar(opParamF32NotSlice), newOpPar(opParamF32NotSlice)},
		[]*CXArgument{newOpPar(opParamF32NotSlice)})
	AddOpCode(OP_F32_SUB, "f32.sub",
		[]*CXArgument{newOpPar(opParamF32NotSlice), newOpPar(opParamF32NotSlice)},
		[]*CXArgument{newOpPar(opParamF32NotSlice)})
	AddOpCode(OP_F32_NEG, "f32.neg",
		[]*CXArgument{newOpPar(opParamF32NotSlice)},
		[]*CXArgument{newOpPar(opParamF32NotSlice)})
	AddOpCode(OP_F32_MUL, "f32.mul",
		[]*CXArgument{newOpPar(opParamF32NotSlice), newOpPar(opParamF32NotSlice)},
		[]*CXArgument{newOpPar(opParamF32NotSlice)})
	AddOpCode(OP_F32_DIV, "f32.div",
		[]*CXArgument{newOpPar(opParamF32NotSlice), newOpPar(opParamF32NotSlice)},
		[]*CXArgument{newOpPar(opParamF32NotSlice)})
	AddOpCode(OP_F32_ABS, "f32.abs",
		[]*CXArgument{newOpPar(opParamF32NotSlice)},
		[]*CXArgument{newOpPar(opParamF32NotSlice)})
	AddOpCode(OP_F32_POW, "f32.pow",
		[]*CXArgument{newOpPar(opParamF32NotSlice), newOpPar(opParamF32NotSlice)},
		[]*CXArgument{newOpPar(opParamF32NotSlice)})
	AddOpCode(OP_F32_GT, "f32.gt",
		[]*CXArgument{newOpPar(opParamF32NotSlice), newOpPar(opParamF32NotSlice)},
		[]*CXArgument{newOpPar(opParamBoolNotSlice)})
	AddOpCode(OP_F32_GTEQ, "f32.gteq",
		[]*CXArgument{newOpPar(opParamF32NotSlice), newOpPar(opParamF32NotSlice)},
		[]*CXArgument{newOpPar(opParamBoolNotSlice)})
	AddOpCode(OP_F32_LT, "f32.lt",
		[]*CXArgument{newOpPar(opParamF32NotSlice), newOpPar(opParamF32NotSlice)},
		[]*CXArgument{newOpPar(opParamBoolNotSlice)})
	AddOpCode(OP_F32_LTEQ, "f32.lteq",
		[]*CXArgument{newOpPar(opParamF32NotSlice), newOpPar(opParamF32NotSlice)},
		[]*CXArgument{newOpPar(opParamBoolNotSlice)})
	AddOpCode(OP_F32_EQ, "f32.eq",
		[]*CXArgument{newOpPar(opParamF32NotSlice), newOpPar(opParamF32NotSlice)},
		[]*CXArgument{newOpPar(opParamBoolNotSlice)})
	AddOpCode(OP_F32_UNEQ, "f32.uneq",
		[]*CXArgument{newOpPar(opParamF32NotSlice), newOpPar(opParamF32NotSlice)},
		[]*CXArgument{newOpPar(opParamBoolNotSlice)})
	AddOpCode(OP_F32_RAND, "f32.rand",
		[]*CXArgument{},
		[]*CXArgument{newOpPar(opParamF32NotSlice)})
	AddOpCode(OP_F32_ACOS, "f32.acos",
		[]*CXArgument{newOpPar(opParamF32NotSlice)},
		[]*CXArgument{newOpPar(opParamF32NotSlice)})
	AddOpCode(OP_F32_COS, "f32.cos",
		[]*CXArgument{newOpPar(opParamF32NotSlice)},
		[]*CXArgument{newOpPar(opParamF32NotSlice)})
	AddOpCode(OP_F32_ASIN, "f32.asin",
		[]*CXArgument{newOpPar(opParamF32NotSlice)},
		[]*CXArgument{newOpPar(opParamF32NotSlice)})
	AddOpCode(OP_F32_SIN, "f32.sin",
		[]*CXArgument{newOpPar(opParamF32NotSlice)},
		[]*CXArgument{newOpPar(opParamF32NotSlice)})
	AddOpCode(OP_F32_SQRT, "f32.sqrt",
		[]*CXArgument{newOpPar(opParamF32NotSlice)},
		[]*CXArgument{newOpPar(opParamF32NotSlice)})
	AddOpCode(OP_F32_LOG, "f32.log",
		[]*CXArgument{newOpPar(opParamF32NotSlice)},
		[]*CXArgument{newOpPar(opParamF32NotSlice)})
	AddOpCode(OP_F32_LOG2, "f32.log2",
		[]*CXArgument{newOpPar(opParamF32NotSlice)},
		[]*CXArgument{newOpPar(opParamF32NotSlice)})
	AddOpCode(OP_F32_LOG10, "f32.log10",
		[]*CXArgument{newOpPar(opParamF32NotSlice)},
		[]*CXArgument{newOpPar(opParamF32NotSlice)})
	AddOpCode(OP_F32_MAX, "f32.max",
		[]*CXArgument{newOpPar(opParamF32NotSlice), newOpPar(opParamF32NotSlice)},
		[]*CXArgument{newOpPar(opParamF32NotSlice)})
	AddOpCode(OP_F32_MIN, "f32.min",
		[]*CXArgument{newOpPar(opParamF32NotSlice), newOpPar(opParamF32NotSlice)},
		[]*CXArgument{newOpPar(opParamF32NotSlice)})

	AddOpCode(OP_F64_IS_NAN, "f64.isnan",
		[]*CXArgument{newOpPar(opParamF64NotSlice)},
		[]*CXArgument{newOpPar(opParamBoolNotSlice)})
	AddOpCode(OP_F64_BYTE, "f64.byte",
		[]*CXArgument{newOpPar(opParamF64NotSlice)},
		[]*CXArgument{newOpPar(opParamByteNotSlice)})
	AddOpCode(OP_F64_STR, "f64.str",
		[]*CXArgument{newOpPar(opParamF64NotSlice)},
		[]*CXArgument{newOpPar(opParamStrNotSlice)})
	AddOpCode(OP_F64_I32, "f64.i32",
		[]*CXArgument{newOpPar(opParamF64NotSlice)},
		[]*CXArgument{newOpPar(opParamI32NotSlice)})
	AddOpCode(OP_F64_I64, "f64.i64",
		[]*CXArgument{newOpPar(opParamF64NotSlice)},
		[]*CXArgument{newOpPar(opParamI64NotSlice)})
	AddOpCode(OP_F64_F32, "f64.f32",
		[]*CXArgument{newOpPar(opParamF64NotSlice)},
		[]*CXArgument{newOpPar(opParamF32NotSlice)})
	AddOpCode(OP_F64_F64, "f64.f64",
		[]*CXArgument{newOpPar(opParamF64NotSlice)},
		[]*CXArgument{newOpPar(opParamF64NotSlice)})

	AddOpCode(OP_F64_PRINT, "f64.print",
		[]*CXArgument{newOpPar(opParamF64NotSlice)},
		[]*CXArgument{})
	AddOpCode(OP_F64_ADD, "f64.add",
		[]*CXArgument{newOpPar(opParamF64NotSlice), newOpPar(opParamF64NotSlice)},
		[]*CXArgument{newOpPar(opParamF64NotSlice)})
	AddOpCode(OP_F64_SUB, "f64.sub",
		[]*CXArgument{newOpPar(opParamF64NotSlice), newOpPar(opParamF64NotSlice)},
		[]*CXArgument{newOpPar(opParamF64NotSlice)})
	AddOpCode(OP_F64_NEG, "f64.neg",
		[]*CXArgument{newOpPar(opParamF64NotSlice)},
		[]*CXArgument{newOpPar(opParamF64NotSlice)})
	AddOpCode(OP_F64_MUL, "f64.mul",
		[]*CXArgument{newOpPar(opParamF64NotSlice), newOpPar(opParamF64NotSlice)},
		[]*CXArgument{newOpPar(opParamF64NotSlice)})
	AddOpCode(OP_F64_DIV, "f64.div",
		[]*CXArgument{newOpPar(opParamF64NotSlice), newOpPar(opParamF64NotSlice)},
		[]*CXArgument{newOpPar(opParamF64NotSlice)})
	AddOpCode(OP_F64_ABS, "f64.abs",
		[]*CXArgument{newOpPar(opParamF64NotSlice)},
		[]*CXArgument{newOpPar(opParamF64NotSlice)})
	AddOpCode(OP_F64_POW, "f64.pow",
		[]*CXArgument{newOpPar(opParamF64NotSlice), newOpPar(opParamF64NotSlice)},
		[]*CXArgument{newOpPar(opParamF64NotSlice)})
	AddOpCode(OP_F64_GT, "f64.gt",
		[]*CXArgument{newOpPar(opParamF64NotSlice), newOpPar(opParamF64NotSlice)},
		[]*CXArgument{newOpPar(opParamBoolNotSlice)})
	AddOpCode(OP_F64_GTEQ, "f64.gteq",
		[]*CXArgument{newOpPar(opParamF64NotSlice), newOpPar(opParamF64NotSlice)},
		[]*CXArgument{newOpPar(opParamBoolNotSlice)})
	AddOpCode(OP_F64_LT, "f64.lt",
		[]*CXArgument{newOpPar(opParamF64NotSlice), newOpPar(opParamF64NotSlice)},
		[]*CXArgument{newOpPar(opParamBoolNotSlice)})
	AddOpCode(OP_F64_LTEQ, "f64.lteq",
		[]*CXArgument{newOpPar(opParamF64NotSlice), newOpPar(opParamF64NotSlice)},
		[]*CXArgument{newOpPar(opParamBoolNotSlice)})
	AddOpCode(OP_F64_EQ, "f64.eq",
		[]*CXArgument{newOpPar(opParamF64NotSlice), newOpPar(opParamF64NotSlice)},
		[]*CXArgument{newOpPar(opParamBoolNotSlice)})
	AddOpCode(OP_F64_UNEQ, "f64.uneq",
		[]*CXArgument{newOpPar(opParamF64NotSlice), newOpPar(opParamF64NotSlice)},
		[]*CXArgument{newOpPar(opParamBoolNotSlice)})
	AddOpCode(OP_F64_RAND, "f64.rand",
		[]*CXArgument{},
		[]*CXArgument{newOpPar(opParamF64NotSlice)})
	AddOpCode(OP_F64_ACOS, "f64.acos",
		[]*CXArgument{newOpPar(opParamF64NotSlice)},
		[]*CXArgument{newOpPar(opParamF64NotSlice)})
	AddOpCode(OP_F64_COS, "f64.cos",
		[]*CXArgument{newOpPar(opParamF64NotSlice)},
		[]*CXArgument{newOpPar(opParamF64NotSlice)})
	AddOpCode(OP_F64_ASIN, "f64.asin",
		[]*CXArgument{newOpPar(opParamF64NotSlice)},
		[]*CXArgument{newOpPar(opParamF64NotSlice)})
	AddOpCode(OP_F64_SIN, "f64.sin",
		[]*CXArgument{newOpPar(opParamF64NotSlice)},
		[]*CXArgument{newOpPar(opParamF64NotSlice)})
	AddOpCode(OP_F64_SQRT, "f64.sqrt",
		[]*CXArgument{newOpPar(opParamF64NotSlice)},
		[]*CXArgument{newOpPar(opParamF64NotSlice)})
	AddOpCode(OP_F64_LOG, "f64.log",
		[]*CXArgument{newOpPar(opParamF64NotSlice)},
		[]*CXArgument{newOpPar(opParamF64NotSlice)})
	AddOpCode(OP_F64_LOG2, "f64.log2",
		[]*CXArgument{newOpPar(opParamF64NotSlice)},
		[]*CXArgument{newOpPar(opParamF64NotSlice)})
	AddOpCode(OP_F64_LOG10, "f64.log10",
		[]*CXArgument{newOpPar(opParamF64NotSlice)},
		[]*CXArgument{newOpPar(opParamF64NotSlice)})
	AddOpCode(OP_F64_MAX, "f64.max",
		[]*CXArgument{newOpPar(opParamF64NotSlice), newOpPar(opParamF64NotSlice)},
		[]*CXArgument{newOpPar(opParamF64NotSlice)})
	AddOpCode(OP_F64_MIN, "f64.min",
		[]*CXArgument{newOpPar(opParamF64NotSlice), newOpPar(opParamF64NotSlice)},
		[]*CXArgument{newOpPar(opParamF64NotSlice)})

	AddOpCode(OP_STR_PRINT, "str.print",
		[]*CXArgument{newOpPar(opParamStrNotSlice)},
		[]*CXArgument{})
	AddOpCode(OP_STR_CONCAT, "str.concat",
		[]*CXArgument{newOpPar(opParamStrNotSlice), newOpPar(opParamStrNotSlice)},
		[]*CXArgument{newOpPar(opParamStrNotSlice)})
	AddOpCode(OP_STR_SUBSTR, "str.substr",
		[]*CXArgument{newOpPar(opParamStrNotSlice), newOpPar(opParamI32NotSlice), newOpPar(opParamI32NotSlice)},
		[]*CXArgument{newOpPar(opParamStrNotSlice)})
	AddOpCode(OP_STR_INDEX, "str.index",
		[]*CXArgument{newOpPar(opParamStrNotSlice), newOpPar(opParamStrNotSlice)},
		[]*CXArgument{newOpPar(opParamI32NotSlice)})
	AddOpCode(OP_STR_TRIM_SPACE, "str.trimspace",
		[]*CXArgument{newOpPar(opParamStrNotSlice)},
		[]*CXArgument{newOpPar(opParamStrNotSlice)})
	AddOpCode(OP_STR_EQ, "str.eq",
		[]*CXArgument{newOpPar(opParamStrNotSlice), newOpPar(opParamStrNotSlice)},
		[]*CXArgument{newOpPar(opParamBoolNotSlice)})

	AddOpCode(OP_STR_BYTE, "str.byte",
		[]*CXArgument{newOpPar(opParamStrNotSlice)},
		[]*CXArgument{newOpPar(opParamByteNotSlice)})
	AddOpCode(OP_STR_STR, "str.str",
		[]*CXArgument{newOpPar(opParamStrNotSlice)},
		[]*CXArgument{newOpPar(opParamStrNotSlice)})
	AddOpCode(OP_STR_I32, "str.i32",
		[]*CXArgument{newOpPar(opParamStrNotSlice)},
		[]*CXArgument{newOpPar(opParamI32NotSlice)})
	AddOpCode(OP_STR_I64, "str.i64",
		[]*CXArgument{newOpPar(opParamStrNotSlice)},
		[]*CXArgument{newOpPar(opParamI64NotSlice)})
	AddOpCode(OP_STR_F32, "str.f32",
		[]*CXArgument{newOpPar(opParamStrNotSlice)},
		[]*CXArgument{newOpPar(opParamF32NotSlice)})
	AddOpCode(OP_STR_F64, "str.f64",
		[]*CXArgument{newOpPar(opParamStrNotSlice)},
		[]*CXArgument{newOpPar(opParamF64NotSlice)})

	AddOpCode(OP_APPEND, "append",
		[]*CXArgument{newOpPar(opParamUndSlice), newOpPar(opParamUndSlice)},
		[]*CXArgument{newOpPar(opParamUndSlice)})
	AddOpCode(OP_RESIZE, "resize",
		[]*CXArgument{newOpPar(opParamUndSlice), newOpPar(opParamI32NotSlice)},
		[]*CXArgument{newOpPar(opParamUndSlice)})
	AddOpCode(OP_INSERT, "insert",
		[]*CXArgument{newOpPar(opParamUndSlice), newOpPar(opParamUndSlice)},
		[]*CXArgument{newOpPar(opParamUndSlice)})
	AddOpCode(OP_REMOVE, "remove",
		[]*CXArgument{newOpPar(opParamUndSlice), newOpPar(opParamI32NotSlice)},
		[]*CXArgument{newOpPar(opParamUndSlice)})
	AddOpCode(OP_COPY, "copy",
		[]*CXArgument{newOpPar(opParamUndSlice), newOpPar(opParamUndSlice)},
		[]*CXArgument{newOpPar(opParamI32NotSlice)})

	AddOpCode(OP_ASSERT, "assert",
		[]*CXArgument{newOpPar(opParamUndNotSlice), newOpPar(opParamUndNotSlice), newOpPar(opParamStrNotSlice)},
		[]*CXArgument{newOpPar(opParamBoolNotSlice)})
	AddOpCode(OP_TEST, "test",
		[]*CXArgument{newOpPar(opParamUndNotSlice), newOpPar(opParamUndNotSlice), newOpPar(opParamStrNotSlice)},
		[]*CXArgument{})
	AddOpCode(OP_PANIC, "panic",
		[]*CXArgument{newOpPar(opParamUndNotSlice), newOpPar(opParamUndNotSlice), newOpPar(opParamStrNotSlice)},
		[]*CXArgument{})
	AddOpCode(OP_STRERROR, "strerror",
		[]*CXArgument{newOpPar(opParamI32NotSlice)},
		[]*CXArgument{newOpPar(opParamStrNotSlice)})
	// affordances
	AddOpCode(OP_AFF_PRINT, "aff.print",
		[]*CXArgument{newOpPar(opParamAffNotSlice)},
		[]*CXArgument{})
	AddOpCode(OP_AFF_QUERY, "aff.query",
		[]*CXArgument{newOpPar(opParamAffNotSlice)},
		[]*CXArgument{newOpPar(opParamAffNotSlice)})
	AddOpCode(OP_AFF_ON, "aff.on",
		[]*CXArgument{newOpPar(opParamAffNotSlice), newOpPar(opParamAffNotSlice)},
		[]*CXArgument{})
	AddOpCode(OP_AFF_OF, "aff.of",
		[]*CXArgument{newOpPar(opParamAffNotSlice), newOpPar(opParamAffNotSlice)},
		[]*CXArgument{})
	AddOpCode(OP_AFF_INFORM, "aff.inform",
		[]*CXArgument{newOpPar(opParamAffNotSlice), newOpPar(opParamI32NotSlice), newOpPar(opParamAffNotSlice)},
		[]*CXArgument{})
	AddOpCode(OP_AFF_REQUEST, "aff.request",
		[]*CXArgument{newOpPar(opParamAffNotSlice), newOpPar(opParamI32NotSlice), newOpPar(opParamAffNotSlice)},
		[]*CXArgument{})

	// http
	AddOpCode(OP_HTTP_SERVE, "http.Serve",
		[]*CXArgument{newOpPar(opParamStrNotSlice)},
		[]*CXArgument{newOpPar(opParamStrNotSlice)})
	AddOpCode(OP_HTTP_LISTEN_AND_SERVE, "http.ListenAndServe",
		[]*CXArgument{newOpPar(opParamStrNotSlice)},
		[]*CXArgument{newOpPar(opParamStrNotSlice)})
	AddOpCode(OP_HTTP_NEW_REQUEST, "http.NewRequest",
		[]*CXArgument{newOpPar(opParamStrNotSlice), newOpPar(opParamStrNotSlice), newOpPar(opParamStrNotSlice)},
		[]*CXArgument{newOpPar(opParamCustomNotSlice)})
	AddOpCode(OP_HTTP_DO, "http.Do",
		[]*CXArgument{newOpPar(opParamCustomNotSlice)},
		[]*CXArgument{newOpPar(opParamCustomNotSlice), newOpPar(opParamStrNotSlice)})
	AddOpCode(OP_DMSG_DO, "http.DmsgDo",
		[]*CXArgument{newOpPar(opParamCustomNotSlice)},
		[]*CXArgument{newOpPar(opParamStrNotSlice)})

	httpRequestType, err := httpPkg.GetStruct("Request")
	if err != nil {
		panic(err)
	}
	requestParam := MakeArgument("Request", "", -1).AddType(TypeNames[TYPE_CUSTOM])
	requestParam.DeclarationSpecifiers = append(requestParam.DeclarationSpecifiers, DECL_STRUCT)
	requestParam.DeclarationSpecifiers = append(requestParam.DeclarationSpecifiers, DECL_POINTER)
	requestParam.IsPointer = true
	requestParam.Size = TYPE_POINTER_SIZE
	requestParam.TotalSize = TYPE_POINTER_SIZE
	requestParam.CustomType = httpRequestType

	AddOpCode(OP_HTTP_HANDLE, "http.Handle",
		[]*CXArgument{newOpPar(opParamStrNotSlice), newOpPar(
			opParam{typCode: TYPE_FUNC,
				pkg: httpPkg,
				inputs: []*CXArgument{
					MakeArgument("ResponseWriter", "", -1).AddType(TypeNames[TYPE_STR]),
					requestParam,
				}})},
		[]*CXArgument{})

	AddOpCode(OP_HTTP_CLOSE, "http.Close",
		[]*CXArgument{},
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
		case OP_UND_SUB, OP_UND_NEG:
			return opSub
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

		case OP_BYTE_BYTE:
			return opByteByte
		case OP_BYTE_STR:
			return opByteByte
		case OP_BYTE_I32:
			return opByteByte
		case OP_BYTE_I64:
			return opByteByte
		case OP_BYTE_F32:
			return opByteByte
		case OP_BYTE_F64:
			return opByteByte

		case OP_BYTE_PRINT:
			return opBytePrint

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

		case OP_I32_BYTE:
			return opI32I32
		case OP_I32_STR:
			return opI32I32
		case OP_I32_I32:
			return opI32I32
		case OP_I32_I64:
			return opI32I32
		case OP_I32_F32:
			return opI32I32
		case OP_I32_F64:
			return opI32I32

		case OP_I32_PRINT:
			return opI32Print
		case OP_I32_ADD:
			return opI32Add
		case OP_I32_SUB, OP_I32_NEG:
			return opI32Sub
		case OP_I32_MUL:
			return opI32Mul
		case OP_I32_DIV:
			return opI32Div
		case OP_I32_ABS:
			return opI32Abs
		case OP_I32_POW:
			return opI32Pow
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
		case OP_I32_MOD:
			return opI32Mod
		case OP_I32_RAND:
			return opI32Rand
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
		case OP_I32_SQRT:
			return opI32Sqrt
		case OP_I32_LOG:
			return opI32Log
		case OP_I32_LOG2:
			return opI32Log2
		case OP_I32_LOG10:
			return opI32Log10

		case OP_I32_MAX:
			return opI32Max
		case OP_I32_MIN:
			return opI32Min

		case OP_I64_BYTE:
			return opI64I64
		case OP_I64_STR:
			return opI64I64
		case OP_I64_I32:
			return opI64I64
		case OP_I64_I64:
			return opI64I64
		case OP_I64_F32:
			return opI64I64
		case OP_I64_F64:
			return opI64I64

		case OP_I64_PRINT:
			return opI64Print
		case OP_I64_ADD:
			return opI64Add
		case OP_I64_SUB, OP_I64_NEG:
			return opI64Sub
		case OP_I64_MUL:
			return opI64Mul
		case OP_I64_DIV:
			return opI64Div
		case OP_I64_ABS:
			return opI64Abs
		case OP_I64_POW:
			return opI64Pow
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
		case OP_I64_MOD:
			return opI64Mod
		case OP_I64_RAND:
			return opI64Rand
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
		case OP_I64_SQRT:
			return opI64Sqrt
		case OP_I64_LOG:
			return opI64Log
		case OP_I64_LOG2:
			return opI64Log2
		case OP_I64_LOG10:
			return opI64Log10
		case OP_I64_MAX:
			return opI64Max
		case OP_I64_MIN:
			return opI64Min

		case OP_F32_IS_NAN:
			return opF32Isnan
		case OP_F32_BYTE:
			return opF32F32
		case OP_F32_STR:
			return opF32F32
		case OP_F32_I32:
			return opF32F32
		case OP_F32_I64:
			return opF32F32
		case OP_F32_F32:
			return opF32F32
		case OP_F32_F64:
			return opF32F32
		case OP_F32_PRINT:
			return opF32Print
		case OP_F32_ADD:
			return opF32Add
		case OP_F32_SUB, OP_F32_NEG:
			return opF32Sub
		case OP_F32_MUL:
			return opF32Mul
		case OP_F32_DIV:
			return opF32Div
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
		case OP_F32_RAND:
			return opF32Rand
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

		case OP_F64_IS_NAN:
			return opF64Isnan
		case OP_F64_BYTE:
			return opF64F64
		case OP_F64_STR:
			return opF64F64
		case OP_F64_I32:
			return opF64F64
		case OP_F64_I64:
			return opF64F64
		case OP_F64_F32:
			return opF64F64
		case OP_F64_F64:
			return opF64F64

		case OP_F64_PRINT:
			return opF64Print
		case OP_F64_ADD:
			return opF64Add
		case OP_F64_SUB, OP_F64_NEG:
			return opF64Sub
		case OP_F64_MUL:
			return opF64Mul
		case OP_F64_DIV:
			return opF64Div
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
		case OP_F64_RAND:
			return opF64Rand
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

		case OP_STR_BYTE:
			return opStrStr
		case OP_STR_STR:
			return opStrStr
		case OP_STR_I32:
			return opStrStr
		case OP_STR_I64:
			return opStrStr
		case OP_STR_F32:
			return opStrStr
		case OP_STR_F64:
			return opStrStr

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

		// http
		case OP_HTTP_SERVE:
			return opHTTPServe
		case OP_HTTP_LISTEN_AND_SERVE:
			return opHTTPListenAndServe
		case OP_HTTP_NEW_REQUEST:
			return opHTTPNewRequest
		case OP_HTTP_DO:
			return opHTTPDo
		case OP_HTTP_HANDLE:
			return opHTTPHandle
		case OP_HTTP_CLOSE:
			return opHTTPClose
		case OP_DMSG_DO:
			return opDMSGDo
		}

		return nil
	}

	opcodeHandlerFinders = append(opcodeHandlerFinders, handleOpcode)
}
