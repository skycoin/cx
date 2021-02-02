package cxcore

import (
	"fmt"
)

// CorePackages ...
var CorePackages = []string{
	// temporary solution until we can implement these packages in pure CX I guess
	"al", "gl", "glfw", "time", "http", "os", "explorer", "aff", "gltext", "cx", "json", "regexp", "cipher",
}

// op codes
// nolint golint
const (
	OP_IDENTITY = iota + 1
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
	OP_STR_LAST_INDEX
	OP_STR_TRIM_SPACE
	OP_STR_EQ

	OP_APPEND
	OP_RESIZE
	OP_INSERT
	OP_REMOVE
	OP_COPY

	OP_ASSERT
	OP_TEST
	OP_PANIC
	OP_PANIC_IF
	OP_PANIC_IF_NOT
	OP_STRERROR

	OP_AFF_PRINT
	OP_AFF_QUERY
	OP_AFF_ON
	OP_AFF_OF
	OP_AFF_INFORM
	OP_AFF_REQUEST

	OP_HTTP_SERVE
	OP_HTTP_LISTEN_AND_SERVE
	OP_HTTP_NEW_REQUEST
	OP_HTTP_DO
	OP_HTTP_HANDLE
	OP_HTTP_CLOSE

	OP_DMSG_DO

	OP_EVOLVE_EVOLVE

	END_OF_CORE_OPS
)

// OpcodeHandler ...
type OpcodeHandler func(prgrm *CXProgram)

var (
	// OpNames ...
	OpNames = map[int]string{}

	// OpCodes ...
	OpCodes = map[string]int{}

	// Natives ...
	Natives        = map[int]*CXFunction{}
	opcodeHandlers []OpcodeHandler
)

func execNative(prgrm *CXProgram) {
	//defer RuntimeError() // High runtime cost.
	opcodeHandlers[prgrm.GetOpCode()](prgrm)
}

// RegisterPackage registers a package on the CX standard library. This does not create a `CXPackage` structure,
// it only tells the CX runtime that `pkgName` will exist by the time a CX program is run.
func RegisterPackage(pkgName string) {
	CorePackages = append(CorePackages, pkgName)
}

// GetOpCodeCount returns an op code that is available for usage on the CX standard library.
func GetOpCodeCount() int {
	return len(opcodeHandlers)
}

// Op ...
func Op(code int, name string, handler OpcodeHandler, inputs []*CXArgument, outputs []*CXArgument) {
	if code >= len(opcodeHandlers) {
		opcodeHandlers = append(opcodeHandlers, make([]OpcodeHandler, code+1)...)
	}
	if opcodeHandlers[code] != nil {
		panic(fmt.Sprintf("duplicate opcode %d : '%s' width '%s'.\n", code, name, OpNames[code]))
	}
	opcodeHandlers[code] = handler

	OpNames[code] = name
	OpCodes[name] = code
	if inputs == nil {
		inputs = []*CXArgument{}
	}
	if outputs == nil {
		outputs = []*CXArgument{}
	}
	Natives[code] = MakeNativeFunction(code, inputs, outputs)
}

/*
// Debug helper function used to find opcodes when they are not registered
func dumpOpCodes(opCode int) {
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

// Pointer takes an already defined `CXArgument` and turns it into a pointer.
func Pointer(arg *CXArgument) *CXArgument {
	arg.DeclarationSpecifiers = append(arg.DeclarationSpecifiers, DECL_POINTER)
	arg.IsPointer = true
	arg.Size = TYPE_POINTER_SIZE
	arg.TotalSize = TYPE_POINTER_SIZE

	return arg
}

// Struct helper for creating a struct parameter. It creates a
// `CXArgument` named `argName`, that represents a structure instane of
// `strctName`, from package `pkgName`.
func Struct(pkgName, strctName, argName string) *CXArgument {
	pkg, err := PROGRAM.GetPackage(pkgName)
	if err != nil {
		panic(err)
	}

	strct, err := pkg.GetStruct(strctName)
	if err != nil {
		panic(err)
	}

	arg := MakeArgument(argName, "", -1).AddType(TypeNames[TYPE_CUSTOM])
	arg.DeclarationSpecifiers = append(arg.DeclarationSpecifiers, DECL_STRUCT)
	arg.Size = strct.Size
	arg.TotalSize = strct.Size
	arg.CustomType = strct

	return arg
}

// Slice Helper function for creating parameters for standard library operators.
// The current standard library only uses basic types and slices. If more options are needed, modify this function
func Slice(typCode int) *CXArgument {
	arg := Param(typCode)
	arg.IsSlice = true
	arg.DeclarationSpecifiers = append(arg.DeclarationSpecifiers, DECL_SLICE)
	return arg
}

// Param ...
func Param(typCode int) *CXArgument {
	arg := MakeArgument("", "", -1).AddType(TypeNames[typCode])
	arg.IsLocalDeclaration = true
	return arg
}

// ParamData ...
type ParamData struct {
	typCode   int           // The type code of the parameter.
	paramType int           // Type of the parameter (struct, slice, etc.).
	strctName string        // Name of the struct in case we're handling a struct instance.
	pkg       *CXPackage    // To what package does this param belongs to.
	inputs    []*CXArgument // Input parameters to a TYPE_FUNC parameter.
	outputs   []*CXArgument // Output parameters to a TYPE_FUNC parameter.
}

// ParamEx Helper function for creating parameters for standard library operators.
// The current standard library only uses basic types and slices. If more options are needed, modify this function
func ParamEx(paramData ParamData) *CXArgument {
	var arg *CXArgument
	switch paramData.paramType {
	case PARAM_DEFAULT:
		arg = Param(paramData.typCode)
	case PARAM_SLICE:
		arg = Slice(paramData.typCode)
	case PARAM_STRUCT:
		arg = Struct(paramData.pkg.Name, paramData.strctName, "")
	}
	arg.Inputs = paramData.inputs
	arg.Outputs = paramData.outputs
	arg.Package = paramData.pkg
	return arg
}

// AI8 Default i8 parameter
var AI8 = Param(TYPE_I8)

// AI16 Default i16 parameter
var AI16 = Param(TYPE_I16)

// AI32 Default i32 parameter
var AI32 = Param(TYPE_I32)

// AI64 Default i64 parameter
var AI64 = Param(TYPE_I64)

// AUI8 Default ui8 parameter
var AUI8 = Param(TYPE_UI8)

// AUI16 Default ui16 parameter
var AUI16 = Param(TYPE_UI16)

// AUI32 Default ui32 parameter
var AUI32 = Param(TYPE_UI32)

// AUI64 Default ui64 parameter
var AUI64 = Param(TYPE_UI64)

// AF32 Default f32 parameter
var AF32 = Param(TYPE_F32)

// AF64 Default f64 parameter
var AF64 = Param(TYPE_F64)

// ASTR Default str parameter
var ASTR = Param(TYPE_STR)

// ABOOL Default bool parameter
var ABOOL = Param(TYPE_BOOL)

// AUND Default und parameter
var AUND = Param(TYPE_UNDEFINED)

// AAFF Default aff parameter
var AAFF = Param(TYPE_AFF)

// In Returns a slice of arguments from an argument list
func In(params ...*CXArgument) []*CXArgument {
	return params
}

// Out Returns a slice of arguments from an argument list
func Out(params ...*CXArgument) []*CXArgument {
	return params
}

func opDebug(prgrm *CXProgram) {
	prgrm.PrintStack()
}

func init() {
	httpPkg, err := PROGRAM.GetPackage("http")
	if err != nil {
		panic(err)
	}

	Op(OP_IDENTITY, "identity", opIdentity, In(AUND), Out(AUND))
	Op(OP_JMP, "jmp", opJmp, In(ABOOL), nil) // AUND to allow 0 inputs (goto)
	Op(OP_DEBUG, "debug", opDebug, nil, nil)
	Op(OP_SERIALIZE, "serialize", opSerialize, In(AAFF), Out(AUI8))
	Op(OP_DESERIALIZE, "deserialize", opDeserialize, In(AUI8), nil)
	Op(OP_UND_EQUAL, "eq", opEqual, In(AUND, AUND), Out(ABOOL))
	Op(OP_UND_UNEQUAL, "uneq", opUnequal, In(AUND, AUND), Out(ABOOL))
	Op(OP_UND_BITAND, "bitand", opBitand, In(AUND, AUND), Out(AUND))
	Op(OP_UND_BITXOR, "bitxor", opBitxor, In(AUND, AUND), Out(AUND))
	Op(OP_UND_BITOR, "bitor", opBitor, In(AUND, AUND), Out(AUND))
	Op(OP_UND_BITCLEAR, "bitclear", opBitclear, In(AUND, AUND), Out(AUND))
	Op(OP_UND_MUL, "mul", opMul, In(AUND, AUND), Out(AUND))
	Op(OP_UND_DIV, "div", opDiv, In(AUND, AUND), Out(AUND))
	Op(OP_UND_MOD, "mod", opMod, In(AUND, AUND), Out(AUND))
	Op(OP_UND_ADD, "add", opAdd, In(AUND, AUND), Out(AUND))
	Op(OP_UND_SUB, "sub", opSub, In(AUND, AUND), Out(AUND))
	Op(OP_UND_NEG, "neg", opNeg, In(AUND), Out(AUND))
	Op(OP_UND_BITSHL, "bitshl", opBitshl, In(AUND, AUND), Out(AUND))
	Op(OP_UND_BITSHR, "bitshr", opBitshr, In(AUND, AUND), Out(AUND))
	Op(OP_UND_LT, "lt", opLt, In(AUND, AUND), Out(ABOOL))
	Op(OP_UND_GT, "gt", opGt, In(AUND, AUND), Out(ABOOL))
	Op(OP_UND_LTEQ, "lteq", opLteq, In(AUND, AUND), Out(ABOOL))
	Op(OP_UND_GTEQ, "gteq", opGteq, In(AUND, AUND), Out(ABOOL))
	Op(OP_UND_LEN, "len", opLen, In(AUND), Out(AI32))
	Op(OP_UND_PRINTF, "printf", opPrintf, In(AUND), nil)
	Op(OP_UND_SPRINTF, "sprintf", opSprintf, In(AUND), Out(ASTR))
	Op(OP_UND_READ, "read", opRead, nil, Out(ASTR))

	Op(OP_BOOL_PRINT, "bool.print", opBoolPrint, In(ABOOL), nil)
	Op(OP_BOOL_EQUAL, "bool.eq", opBoolEqual, In(ABOOL, ABOOL), Out(ABOOL))
	Op(OP_BOOL_UNEQUAL, "bool.uneq", opBoolUnequal, In(ABOOL, ABOOL), Out(ABOOL))
	Op(OP_BOOL_NOT, "bool.not", opBoolNot, In(ABOOL), Out(ABOOL))
	Op(OP_BOOL_OR, "bool.or", opBoolOr, In(ABOOL, ABOOL), Out(ABOOL))
	Op(OP_BOOL_AND, "bool.and", opBoolAnd, In(ABOOL, ABOOL), Out(ABOOL))

	Op(OP_I8_STR, "i8.str", opI8ToStr, In(AI8), Out(ASTR))
	Op(OP_I8_I16, "i8.i16", opI8ToI16, In(AI8), Out(AI16))
	Op(OP_I8_I32, "i8.i32", opI8ToI32, In(AI8), Out(AI32))
	Op(OP_I8_I64, "i8.i64", opI8ToI64, In(AI8), Out(AI64))
	Op(OP_I8_UI8, "i8.ui8", opI8ToUI8, In(AI8), Out(AUI8))
	Op(OP_I8_UI16, "i8.ui16", opI8ToUI16, In(AI8), Out(AUI16))
	Op(OP_I8_UI32, "i8.ui32", opI8ToUI32, In(AI8), Out(AUI32))
	Op(OP_I8_UI64, "i8.ui64", opI8ToUI64, In(AI8), Out(AUI64))
	Op(OP_I8_F32, "i8.f32", opI8ToF32, In(AI8), Out(AF32))
	Op(OP_I8_F64, "i8.f64", opI8ToF64, In(AI8), Out(AF64))
	Op(OP_I8_PRINT, "i8.print", opI8Print, In(AI8), nil)
	Op(OP_I8_ADD, "i8.add", opI8Add, In(AI8, AI8), Out(AI8))
	Op(OP_I8_SUB, "i8.sub", opI8Sub, In(AI8, AI8), Out(AI8))
	Op(OP_I8_NEG, "i8.neg", opI8Neg, In(AI8), Out(AI8))
	Op(OP_I8_MUL, "i8.mul", opI8Mul, In(AI8, AI8), Out(AI8))
	Op(OP_I8_DIV, "i8.div", opI8Div, In(AI8, AI8), Out(AI8))
	Op(OP_I8_MOD, "i8.mod", opI8Mod, In(AI8, AI8), Out(AI8))
	Op(OP_I8_ABS, "i8.abs", opI8Abs, In(AI8), Out(AI8))
	Op(OP_I8_GT, "i8.gt", opI8Gt, In(AI8, AI8), Out(ABOOL))
	Op(OP_I8_GTEQ, "i8.gteq", opI8Gteq, In(AI8, AI8), Out(ABOOL))
	Op(OP_I8_LT, "i8.lt", opI8Lt, In(AI8, AI8), Out(ABOOL))
	Op(OP_I8_LTEQ, "i8.lteq", opI8Lteq, In(AI8, AI8), Out(ABOOL))
	Op(OP_I8_EQ, "i8.eq", opI8Eq, In(AI8, AI8), Out(ABOOL))
	Op(OP_I8_UNEQ, "i8.uneq", opI8Uneq, In(AI8, AI8), Out(ABOOL))
	Op(OP_I8_BITAND, "i8.bitand", opI8Bitand, In(AI8, AI8), Out(AI8))
	Op(OP_I8_BITOR, "i8.bitor", opI8Bitor, In(AI8, AI8), Out(AI8))
	Op(OP_I8_BITXOR, "i8.bitxor", opI8Bitxor, In(AI8, AI8), Out(AI8))
	Op(OP_I8_BITCLEAR, "i8.bitclear", opI8Bitclear, In(AI8, AI8), Out(AI8))
	Op(OP_I8_BITSHL, "i8.bitshl", opI8Bitshl, In(AI8, AI8), Out(AI8))
	Op(OP_I8_BITSHR, "i8.bitshr", opI8Bitshr, In(AI8, AI8), Out(AI8))
	Op(OP_I8_MAX, "i8.max", opI8Max, In(AI8, AI8), Out(AI8))
	Op(OP_I8_MIN, "i8.min", opI8Min, In(AI8, AI8), Out(AI8))
	Op(OP_I8_RAND, "i8.rand", opI8Rand, In(AI8, AI8), Out(AI8))

	Op(OP_I16_STR, "i16.str", opI16ToStr, In(AI16), Out(ASTR))
	Op(OP_I16_I8, "i16.i8", opI16ToI8, In(AI16), Out(AI8))
	Op(OP_I16_I32, "i16.i32", opI16ToI32, In(AI16), Out(AI32))
	Op(OP_I16_I64, "i16.i64", opI16ToI64, In(AI16), Out(AI64))
	Op(OP_I16_UI8, "i16.ui8", opI16ToUI8, In(AI16), Out(AUI8))
	Op(OP_I16_UI16, "i16.ui16", opI16ToUI16, In(AI16), Out(AUI16))
	Op(OP_I16_UI32, "i16.ui32", opI16ToUI32, In(AI16), Out(AUI32))
	Op(OP_I16_UI64, "i16.ui64", opI16ToUI64, In(AI16), Out(AUI64))
	Op(OP_I16_F32, "i16.f32", opI16ToF32, In(AI16), Out(AF32))
	Op(OP_I16_F64, "i16.f64", opI16ToF64, In(AI16), Out(AF64))
	Op(OP_I16_PRINT, "i16.print", opI16Print, In(AI16), nil)
	Op(OP_I16_ADD, "i16.add", opI16Add, In(AI16, AI16), Out(AI16))
	Op(OP_I16_SUB, "i16.sub", opI16Sub, In(AI16, AI16), Out(AI16))
	Op(OP_I16_NEG, "i16.neg", opI16Neg, In(AI16), Out(AI16))
	Op(OP_I16_MUL, "i16.mul", opI16Mul, In(AI16, AI16), Out(AI16))
	Op(OP_I16_DIV, "i16.div", opI16Div, In(AI16, AI16), Out(AI16))
	Op(OP_I16_MOD, "i16.mod", opI16Mod, In(AI16, AI16), Out(AI16))
	Op(OP_I16_ABS, "i16.abs", opI16Abs, In(AI16), Out(AI16))
	Op(OP_I16_GT, "i16.gt", opI16Gt, In(AI16, AI16), Out(ABOOL))
	Op(OP_I16_GTEQ, "i16.gteq", opI16Gteq, In(AI16, AI16), Out(ABOOL))
	Op(OP_I16_LT, "i16.lt", opI16Lt, In(AI16, AI16), Out(ABOOL))
	Op(OP_I16_LTEQ, "i16.lteq", opI16Lteq, In(AI16, AI16), Out(ABOOL))
	Op(OP_I16_EQ, "i16.eq", opI16Eq, In(AI16, AI16), Out(ABOOL))
	Op(OP_I16_UNEQ, "i16.uneq", opI16Uneq, In(AI16, AI16), Out(ABOOL))
	Op(OP_I16_BITAND, "i16.bitand", opI16Bitand, In(AI16, AI16), Out(AI16))
	Op(OP_I16_BITOR, "i16.bitor", opI16Bitor, In(AI16, AI16), Out(AI16))
	Op(OP_I16_BITXOR, "i16.bitxor", opI16Bitxor, In(AI16, AI16), Out(AI16))
	Op(OP_I16_BITCLEAR, "i16.bitclear", opI16Bitclear, In(AI16, AI16), Out(AI16))
	Op(OP_I16_BITSHL, "i16.bitshl", opI16Bitshl, In(AI16, AI16), Out(AI16))
	Op(OP_I16_BITSHR, "i16.bitshr", opI16Bitshr, In(AI16, AI16), Out(AI16))
	Op(OP_I16_MAX, "i16.max", opI16Max, In(AI16, AI16), Out(AI16))
	Op(OP_I16_MIN, "i16.min", opI16Min, In(AI16, AI16), Out(AI16))
	Op(OP_I16_RAND, "i16.rand", opI16Rand, In(AI16, AI16), Out(AI16))

	Op(OP_I32_STR, "i32.str", opI32ToStr, In(AI32), Out(ASTR))
	Op(OP_I32_I8, "i32.i8", opI32ToI8, In(AI32), Out(AI8))
	Op(OP_I32_I16, "i32.i16", opI32ToI16, In(AI32), Out(AI16))
	Op(OP_I32_I64, "i32.i64", opI32ToI64, In(AI32), Out(AI64))
	Op(OP_I32_UI8, "i32.ui8", opI32ToUI8, In(AI32), Out(AUI8))
	Op(OP_I32_UI16, "i32.ui16", opI32ToUI16, In(AI32), Out(AUI16))
	Op(OP_I32_UI32, "i32.ui32", opI32ToUI32, In(AI32), Out(AUI32))
	Op(OP_I32_UI64, "i32.ui64", opI32ToUI64, In(AI32), Out(AUI64))
	Op(OP_I32_F32, "i32.f32", opI32ToF32, In(AI32), Out(AF32))
	Op(OP_I32_F64, "i32.f64", opI32ToF64, In(AI32), Out(AF64))
	Op(OP_I32_PRINT, "i32.print", opI32Print, In(AI32), nil)
	Op(OP_I32_ADD, "i32.add", opI32Add, In(AI32, AI32), Out(AI32))
	Op(OP_I32_SUB, "i32.sub", opI32Sub, In(AI32, AI32), Out(AI32))
	Op(OP_I32_NEG, "i32.neg", opI32Neg, In(AI32), Out(AI32))
	Op(OP_I32_MUL, "i32.mul", opI32Mul, In(AI32, AI32), Out(AI32))
	Op(OP_I32_DIV, "i32.div", opI32Div, In(AI32, AI32), Out(AI32))
	Op(OP_I32_MOD, "i32.mod", opI32Mod, In(AI32, AI32), Out(AI32))
	Op(OP_I32_ABS, "i32.abs", opI32Abs, In(AI32), Out(AI32))
	Op(OP_I32_GT, "i32.gt", opI32Gt, In(AI32, AI32), Out(ABOOL))
	Op(OP_I32_GTEQ, "i32.gteq", opI32Gteq, In(AI32, AI32), Out(ABOOL))
	Op(OP_I32_LT, "i32.lt", opI32Lt, In(AI32, AI32), Out(ABOOL))
	Op(OP_I32_LTEQ, "i32.lteq", opI32Lteq, In(AI32, AI32), Out(ABOOL))
	Op(OP_I32_EQ, "i32.eq", opI32Eq, In(AI32, AI32), Out(ABOOL))
	Op(OP_I32_UNEQ, "i32.uneq", opI32Uneq, In(AI32, AI32), Out(ABOOL))
	Op(OP_I32_BITAND, "i32.bitand", opI32Bitand, In(AI32, AI32), Out(AI32))
	Op(OP_I32_BITOR, "i32.bitor", opI32Bitor, In(AI32, AI32), Out(AI32))
	Op(OP_I32_BITXOR, "i32.bitxor", opI32Bitxor, In(AI32, AI32), Out(AI32))
	Op(OP_I32_BITCLEAR, "i32.bitclear", opI32Bitclear, In(AI32, AI32), Out(AI32))
	Op(OP_I32_BITSHL, "i32.bitshl", opI32Bitshl, In(AI32, AI32), Out(AI32))
	Op(OP_I32_BITSHR, "i32.bitshr", opI32Bitshr, In(AI32, AI32), Out(AI32))
	Op(OP_I32_MAX, "i32.max", opI32Max, In(AI32, AI32), Out(AI32))
	Op(OP_I32_MIN, "i32.min", opI32Min, In(AI32, AI32), Out(AI32))
	Op(OP_I32_RAND, "i32.rand", opI32Rand, In(AI32, AI32), Out(AI32))

	Op(OP_I64_STR, "i64.str", opI64ToStr, In(AI64), Out(ASTR))
	Op(OP_I64_I8, "i64.i8", opI64ToI8, In(AI64), Out(AI8))
	Op(OP_I64_I16, "i64.i16", opI64ToI16, In(AI64), Out(AI16))
	Op(OP_I64_I32, "i64.i32", opI64ToI32, In(AI64), Out(AI32))
	Op(OP_I64_UI8, "i64.ui8", opI64ToUI8, In(AI64), Out(AUI8))
	Op(OP_I64_UI16, "i64.ui16", opI64ToUI16, In(AI64), Out(AUI16))
	Op(OP_I64_UI32, "i64.ui32", opI64ToUI32, In(AI64), Out(AUI32))
	Op(OP_I64_UI64, "i64.ui64", opI64ToUI64, In(AI64), Out(AUI64))
	Op(OP_I64_F32, "i64.f32", opI64ToF32, In(AI64), Out(AF32))
	Op(OP_I64_F64, "i64.f64", opI64ToF64, In(AI64), Out(AF64))
	Op(OP_I64_PRINT, "i64.print", opI64Print, In(AI64), nil)
	Op(OP_I64_ADD, "i64.add", opI64Add, In(AI64, AI64), Out(AI64))
	Op(OP_I64_SUB, "i64.sub", opI64Sub, In(AI64, AI64), Out(AI64))
	Op(OP_I64_NEG, "i64.neg", opI64Neg, In(AI64), Out(AI64))
	Op(OP_I64_MUL, "i64.mul", opI64Mul, In(AI64, AI64), Out(AI64))
	Op(OP_I64_DIV, "i64.div", opI64Div, In(AI64, AI64), Out(AI64))
	Op(OP_I64_MOD, "i64.mod", opI64Mod, In(AI64, AI64), Out(AI64))
	Op(OP_I64_ABS, "i64.abs", opI64Abs, In(AI64), Out(AI64))
	Op(OP_I64_GT, "i64.gt", opI64Gt, In(AI64, AI64), Out(ABOOL))
	Op(OP_I64_GTEQ, "i64.gteq", opI64Gteq, In(AI64, AI64), Out(ABOOL))
	Op(OP_I64_LT, "i64.lt", opI64Lt, In(AI64, AI64), Out(ABOOL))
	Op(OP_I64_LTEQ, "i64.lteq", opI64Lteq, In(AI64, AI64), Out(ABOOL))
	Op(OP_I64_EQ, "i64.eq", opI64Eq, In(AI64, AI64), Out(ABOOL))
	Op(OP_I64_UNEQ, "i64.uneq", opI64Uneq, In(AI64, AI64), Out(ABOOL))
	Op(OP_I64_BITAND, "i64.bitand", opI64Bitand, In(AI64, AI64), Out(AI64))
	Op(OP_I64_BITOR, "i64.bitor", opI64Bitor, In(AI64, AI64), Out(AI64))
	Op(OP_I64_BITXOR, "i64.bitxor", opI64Bitxor, In(AI64, AI64), Out(AI64))
	Op(OP_I64_BITCLEAR, "i64.bitclear", opI64Bitclear, In(AI64, AI64), Out(AI64))
	Op(OP_I64_BITSHL, "i64.bitshl", opI64Bitshl, In(AI64, AI64), Out(AI64))
	Op(OP_I64_BITSHR, "i64.bitshr", opI64Bitshr, In(AI64, AI64), Out(AI64))
	Op(OP_I64_MAX, "i64.max", opI64Max, In(AI64, AI64), Out(AI64))
	Op(OP_I64_MIN, "i64.min", opI64Min, In(AI64, AI64), Out(AI64))
	Op(OP_I64_RAND, "i64.rand", opI64Rand, In(AI64, AI64), Out(AI64))

	Op(OP_UI8_STR, "ui8.str", opUI8ToStr, In(AUI8), Out(ASTR))
	Op(OP_UI8_I8, "ui8.i8", opUI8ToI8, In(AUI8), Out(AI8))
	Op(OP_UI8_I16, "ui8.i16", opUI8ToI16, In(AUI8), Out(AI16))
	Op(OP_UI8_I32, "ui8.i32", opUI8ToI32, In(AUI8), Out(AI32))
	Op(OP_UI8_I64, "ui8.i64", opUI8ToI64, In(AUI8), Out(AI64))
	Op(OP_UI8_UI16, "ui8.ui16", opUI8ToUI16, In(AUI8), Out(AUI16))
	Op(OP_UI8_UI32, "ui8.ui32", opUI8ToUI32, In(AUI8), Out(AUI32))
	Op(OP_UI8_UI64, "ui8.ui64", opUI8ToUI64, In(AUI8), Out(AUI64))
	Op(OP_UI8_F32, "ui8.f32", opUI8ToF32, In(AUI8), Out(AF32))
	Op(OP_UI8_F64, "ui8.f64", opUI8ToF64, In(AUI8), Out(AF64))
	Op(OP_UI8_PRINT, "ui8.print", opUI8Print, In(AUI8), nil)
	Op(OP_UI8_ADD, "ui8.add", opUI8Add, In(AUI8, AUI8), Out(AUI8))
	Op(OP_UI8_SUB, "ui8.sub", opUI8Sub, In(AUI8, AUI8), Out(AUI8))
	Op(OP_UI8_MUL, "ui8.mul", opUI8Mul, In(AUI8, AUI8), Out(AUI8))
	Op(OP_UI8_DIV, "ui8.div", opUI8Div, In(AUI8, AUI8), Out(AUI8))
	Op(OP_UI8_MOD, "ui8.mod", opUI8Mod, In(AUI8, AUI8), Out(AUI8))
	Op(OP_UI8_GT, "ui8.gt", opUI8Gt, In(AUI8, AUI8), Out(ABOOL))
	Op(OP_UI8_GTEQ, "ui8.gteq", opUI8Gteq, In(AUI8, AUI8), Out(ABOOL))
	Op(OP_UI8_LT, "ui8.lt", opUI8Lt, In(AUI8, AUI8), Out(ABOOL))
	Op(OP_UI8_LTEQ, "ui8.lteq", opUI8Lteq, In(AUI8, AUI8), Out(ABOOL))
	Op(OP_UI8_EQ, "ui8.eq", opUI8Eq, In(AUI8, AUI8), Out(ABOOL))
	Op(OP_UI8_UNEQ, "ui8.uneq", opUI8Uneq, In(AUI8, AUI8), Out(ABOOL))
	Op(OP_UI8_BITAND, "ui8.bitand", opUI8Bitand, In(AUI8, AUI8), Out(AUI8))
	Op(OP_UI8_BITOR, "ui8.bitor", opUI8Bitor, In(AUI8, AUI8), Out(AUI8))
	Op(OP_UI8_BITXOR, "ui8.bitxor", opUI8Bitxor, In(AUI8, AUI8), Out(AUI8))
	Op(OP_UI8_BITCLEAR, "ui8.bitclear", opUI8Bitclear, In(AUI8, AUI8), Out(AUI8))
	Op(OP_UI8_BITSHL, "ui8.bitshl", opUI8Bitshl, In(AUI8, AUI8), Out(AUI8))
	Op(OP_UI8_BITSHR, "ui8.bitshr", opUI8Bitshr, In(AUI8, AUI8), Out(AUI8))
	Op(OP_UI8_MAX, "ui8.max", opUI8Max, In(AUI8, AUI8), Out(AUI8))
	Op(OP_UI8_MIN, "ui8.min", opUI8Min, In(AUI8, AUI8), Out(AUI8))
	Op(OP_UI8_RAND, "ui8.rand", opUI8Rand, nil, Out(AUI8))

	Op(OP_UI16_STR, "ui16.str", opUI16ToStr, In(AUI16), Out(ASTR))
	Op(OP_UI16_I8, "ui16.i8", opUI16ToI8, In(AUI16), Out(AI8))
	Op(OP_UI16_I16, "ui16.i16", opUI16ToI16, In(AUI16), Out(AI16))
	Op(OP_UI16_I32, "ui16.i32", opUI16ToI32, In(AUI16), Out(AI32))
	Op(OP_UI16_I64, "ui16.i64", opUI16ToI64, In(AUI16), Out(AI64))
	Op(OP_UI16_UI8, "ui16.ui8", opUI16ToUI8, In(AUI16), Out(AUI8))
	Op(OP_UI16_UI32, "ui16.ui32", opUI16ToUI32, In(AUI16), Out(AUI32))
	Op(OP_UI16_UI64, "ui16.ui64", opUI16ToUI64, In(AUI16), Out(AUI64))
	Op(OP_UI16_F32, "ui16.f32", opUI16ToF32, In(AUI16), Out(AF32))
	Op(OP_UI16_F64, "ui16.f64", opUI16ToF64, In(AUI16), Out(AF64))
	Op(OP_UI16_PRINT, "ui16.print", opUI16Print, In(AUI16), nil)
	Op(OP_UI16_ADD, "ui16.add", opUI16Add, In(AUI16, AUI16), Out(AUI16))
	Op(OP_UI16_SUB, "ui16.sub", opUI16Sub, In(AUI16, AUI16), Out(AUI16))
	Op(OP_UI16_MUL, "ui16.mul", opUI16Mul, In(AUI16, AUI16), Out(AUI16))
	Op(OP_UI16_DIV, "ui16.div", opUI16Div, In(AUI16, AUI16), Out(AUI16))
	Op(OP_UI16_MOD, "ui16.mod", opUI16Mod, In(AUI16, AUI16), Out(AUI16))
	Op(OP_UI16_GT, "ui16.gt", opUI16Gt, In(AUI16, AUI16), In(ABOOL))
	Op(OP_UI16_GTEQ, "ui16.gteq", opUI16Gteq, In(AUI16, AUI16), Out(ABOOL))
	Op(OP_UI16_LT, "ui16.lt", opUI16Lt, In(AUI16, AUI16), Out(ABOOL))
	Op(OP_UI16_LTEQ, "ui16.lteq", opUI16Lteq, In(AUI16, AUI16), Out(ABOOL))
	Op(OP_UI16_EQ, "ui16.eq", opUI16Eq, In(AUI16, AUI16), Out(ABOOL))
	Op(OP_UI16_UNEQ, "ui16.uneq", opUI16Uneq, In(AUI16, AUI16), Out(ABOOL))
	Op(OP_UI16_BITAND, "ui16.bitand", opUI16Bitand, In(AUI16, AUI16), Out(AUI16))
	Op(OP_UI16_BITOR, "ui16.bitor", opUI16Bitor, In(AUI16, AUI16), Out(AUI16))
	Op(OP_UI16_BITXOR, "ui16.bitxor", opUI16Bitxor, In(AUI16, AUI16), Out(AUI16))
	Op(OP_UI16_BITCLEAR, "ui16.bitclear", opUI16Bitclear, In(AUI16, AUI16), Out(AUI16))
	Op(OP_UI16_BITSHL, "ui16.bitshl", opUI16Bitshl, In(AUI16, AUI16), Out(AUI16))
	Op(OP_UI16_BITSHR, "ui16.bitshr", opUI16Bitshr, In(AUI16, AUI16), Out(AUI16))
	Op(OP_UI16_MAX, "ui16.max", opUI16Max, In(AUI16, AUI16), Out(AUI16))
	Op(OP_UI16_MIN, "ui16.min", opUI16Min, In(AUI16, AUI16), Out(AUI16))
	Op(OP_UI16_RAND, "ui16.rand", opUI16Rand, nil, Out(AUI16))

	Op(OP_UI32_STR, "ui32.str", opUI32ToStr, In(AUI32), Out(ASTR))
	Op(OP_UI32_I8, "ui32.i8", opUI32ToI8, In(AUI32), Out(AI8))
	Op(OP_UI32_I16, "ui32.i16", opUI32ToI16, In(AUI32), Out(AI16))
	Op(OP_UI32_I32, "ui32.i32", opUI32ToI32, In(AUI32), Out(AI32))
	Op(OP_UI32_I64, "ui32.i64", opUI32ToI64, In(AUI32), Out(AI64))
	Op(OP_UI32_UI8, "ui32.ui8", opUI32ToUI8, In(AUI32), Out(AUI8))
	Op(OP_UI32_UI16, "ui32.ui16", opUI32ToUI16, In(AUI32), Out(AUI16))
	Op(OP_UI32_UI64, "ui32.ui64", opUI32ToUI64, In(AUI32), Out(AUI64))
	Op(OP_UI32_F32, "ui32.f32", opUI32ToF32, In(AUI32), Out(AF32))
	Op(OP_UI32_F64, "ui32.f64", opUI32ToF64, In(AUI32), Out(AF64))
	Op(OP_UI32_PRINT, "ui32.print", opUI32Print, In(AUI32), nil)
	Op(OP_UI32_ADD, "ui32.add", opUI32Add, In(AUI32, AUI32), Out(AUI32))
	Op(OP_UI32_SUB, "ui32.sub", opUI32Sub, In(AUI32, AUI32), Out(AUI32))
	Op(OP_UI32_MUL, "ui32.mul", opUI32Mul, In(AUI32, AUI32), Out(AUI32))
	Op(OP_UI32_DIV, "ui32.div", opUI32Div, In(AUI32, AUI32), Out(AUI32))
	Op(OP_UI32_MOD, "ui32.mod", opUI32Mod, In(AUI32, AUI32), Out(AUI32))
	Op(OP_UI32_GT, "ui32.gt", opUI32Gt, In(AUI32, AUI32), Out(ABOOL))
	Op(OP_UI32_GTEQ, "ui32.gteq", opUI32Gteq, In(AUI32, AUI32), Out(ABOOL))
	Op(OP_UI32_LT, "ui32.lt", opUI32Lt, In(AUI32, AUI32), Out(ABOOL))
	Op(OP_UI32_LTEQ, "ui32.lteq", opUI32Lteq, In(AUI32, AUI32), Out(ABOOL))
	Op(OP_UI32_EQ, "ui32.eq", opUI32Eq, In(AUI32, AUI32), Out(ABOOL))
	Op(OP_UI32_UNEQ, "ui32.uneq", opUI32Uneq, In(AUI32, AUI32), Out(ABOOL))
	Op(OP_UI32_BITAND, "ui32.bitand", opUI32Bitand, In(AUI32, AUI32), Out(AUI32))
	Op(OP_UI32_BITOR, "ui32.bitor", opUI32Bitor, In(AUI32, AUI32), Out(AUI32))
	Op(OP_UI32_BITXOR, "ui32.bitxor", opUI32Bitxor, In(AUI32, AUI32), Out(AUI32))
	Op(OP_UI32_BITCLEAR, "ui32.bitclear", opUI32Bitclear, In(AUI32, AUI32), Out(AUI32))
	Op(OP_UI32_BITSHL, "ui32.bitshl", opUI32Bitshl, In(AUI32, AUI32), Out(AUI32))
	Op(OP_UI32_BITSHR, "ui32.bitshr", opUI32Bitshr, In(AUI32, AUI32), Out(AUI32))
	Op(OP_UI32_MAX, "ui32.max", opUI32Max, In(AUI32, AUI32), Out(AUI32))
	Op(OP_UI32_MIN, "ui32.min", opUI32Min, In(AUI32, AUI32), Out(AUI32))
	Op(OP_UI32_RAND, "ui32.rand", opUI32Rand, nil, Out(AUI32))

	Op(OP_UI64_STR, "ui64.str", opUI64ToStr, In(AUI64), Out(ASTR))
	Op(OP_UI64_I8, "ui64.i8", opUI64ToI8, In(AUI64), Out(AI8))
	Op(OP_UI64_I16, "ui64.i16", opUI64ToI16, In(AUI64), Out(AI16))
	Op(OP_UI64_I32, "ui64.i32", opUI64ToI32, In(AUI64), Out(AI32))
	Op(OP_UI64_I64, "ui64.i64", opUI64ToI64, In(AUI64), Out(AI64))
	Op(OP_UI64_UI8, "ui64.ui8", opUI64ToUI8, In(AUI64), Out(AUI8))
	Op(OP_UI64_UI16, "ui64.ui16", opUI64ToUI16, In(AUI64), Out(AUI16))
	Op(OP_UI64_UI32, "ui64.ui32", opUI64ToUI32, In(AUI64), Out(AUI32))
	Op(OP_UI64_F32, "ui64.f32", opUI64ToF32, In(AUI64), Out(AF32))
	Op(OP_UI64_F64, "ui64.f64", opUI64ToF64, In(AUI64), Out(AF64))
	Op(OP_UI64_PRINT, "ui64.print", opUI64Print, In(AUI64), nil)
	Op(OP_UI64_ADD, "ui64.add", opUI64Add, In(AUI64, AUI64), Out(AUI64))
	Op(OP_UI64_SUB, "ui64.sub", opUI64Sub, In(AUI64, AUI64), Out(AUI64))
	Op(OP_UI64_MUL, "ui64.mul", opUI64Mul, In(AUI64, AUI64), Out(AUI64))
	Op(OP_UI64_DIV, "ui64.div", opUI64Div, In(AUI64, AUI64), Out(AUI64))
	Op(OP_UI64_MOD, "ui64.mod", opUI64Mod, In(AUI64, AUI64), Out(AUI64))
	Op(OP_UI64_GT, "ui64.gt", opUI64Gt, In(AUI64, AUI64), Out(ABOOL))
	Op(OP_UI64_GTEQ, "ui64.gteq", opUI64Gteq, In(AUI64, AUI64), Out(ABOOL))
	Op(OP_UI64_LT, "ui64.lt", opUI64Lt, In(AUI64, AUI64), Out(ABOOL))
	Op(OP_UI64_LTEQ, "ui64.lteq", opUI64Lteq, In(AUI64, AUI64), Out(ABOOL))
	Op(OP_UI64_EQ, "ui64.eq", opUI64Eq, In(AUI64, AUI64), Out(ABOOL))
	Op(OP_UI64_UNEQ, "ui64.uneq", opUI64Uneq, In(AUI64, AUI64), Out(ABOOL))
	Op(OP_UI64_BITAND, "ui64.bitand", opUI64Bitand, In(AUI64, AUI64), Out(AUI64))
	Op(OP_UI64_BITOR, "ui64.bitor", opUI64Bitor, In(AUI64, AUI64), Out(AUI64))
	Op(OP_UI64_BITXOR, "ui64.bitxor", opUI64Bitxor, In(AUI64, AUI64), Out(AUI64))
	Op(OP_UI64_BITCLEAR, "ui64.bitclear", opUI64Bitclear, In(AUI64, AUI64), Out(AUI64))
	Op(OP_UI64_BITSHL, "ui64.bitshl", opUI64Bitshl, In(AUI64, AUI64), Out(AUI64))
	Op(OP_UI64_BITSHR, "ui64.bitshr", opUI64Bitshr, In(AUI64, AUI64), Out(AUI64))
	Op(OP_UI64_MAX, "ui64.max", opUI64Max, In(AUI64, AUI64), Out(AUI64))
	Op(OP_UI64_MIN, "ui64.min", opUI64Min, In(AUI64, AUI64), Out(AUI64))
	Op(OP_UI64_RAND, "ui64.rand", opUI64Rand, nil, Out(AUI64))

	Op(OP_F32_IS_NAN, "f32.isnan", opF32Isnan, In(AF32), Out(ABOOL))
	Op(OP_F32_STR, "f32.str", opF32ToStr, In(AF32), Out(ASTR))
	Op(OP_F32_I8, "f32.i8", opF32ToI8, In(AF32), Out(AI8))
	Op(OP_F32_I16, "f32.i16", opF32ToI16, In(AF32), Out(AI16))
	Op(OP_F32_I32, "f32.i32", opF32ToI32, In(AF32), Out(AI32))
	Op(OP_F32_I64, "f32.i64", opF32ToI64, In(AF32), Out(AI64))
	Op(OP_F32_UI8, "f32.ui8", opF32ToUI8, In(AF32), Out(AUI8))
	Op(OP_F32_UI16, "f32.ui16", opF32ToUI16, In(AF32), Out(AUI16))
	Op(OP_F32_UI32, "f32.ui32", opF32ToUI32, In(AF32), Out(AUI32))
	Op(OP_F32_UI64, "f32.ui64", opF32ToUI64, In(AF32), Out(AUI64))
	Op(OP_F32_F64, "f32.f64", opF32ToF64, In(AF32), Out(AF64))
	Op(OP_F32_PRINT, "f32.print", opF32Print, In(AF32), nil)
	Op(OP_F32_ADD, "f32.add", opF32Add, In(AF32, AF32), Out(AF32))
	Op(OP_F32_SUB, "f32.sub", opF32Sub, In(AF32, AF32), Out(AF32))
	Op(OP_F32_NEG, "f32.neg", opF32Neg, In(AF32), Out(AF32))
	Op(OP_F32_MUL, "f32.mul", opF32Mul, In(AF32, AF32), Out(AF32))
	Op(OP_F32_DIV, "f32.div", opF32Div, In(AF32, AF32), Out(AF32))
	Op(OP_F32_MOD, "f32.mod", opF32Mod, In(AF32, AF32), Out(AF32))
	Op(OP_F32_ABS, "f32.abs", opF32Abs, In(AF32), Out(AF32))
	Op(OP_F32_POW, "f32.pow", opF32Pow, In(AF32, AF32), Out(AF32))
	Op(OP_F32_GT, "f32.gt", opF32Gt, In(AF32, AF32), Out(ABOOL))
	Op(OP_F32_GTEQ, "f32.gteq", opF32Gteq, In(AF32, AF32), Out(ABOOL))
	Op(OP_F32_LT, "f32.lt", opF32Lt, In(AF32, AF32), Out(ABOOL))
	Op(OP_F32_LTEQ, "f32.lteq", opF32Lteq, In(AF32, AF32), Out(ABOOL))
	Op(OP_F32_EQ, "f32.eq", opF32Eq, In(AF32, AF32), Out(ABOOL))
	Op(OP_F32_UNEQ, "f32.uneq", opF32Uneq, In(AF32, AF32), Out(ABOOL))
	Op(OP_F32_ACOS, "f32.acos", opF32Acos, In(AF32), Out(AF32))
	Op(OP_F32_COS, "f32.cos", opF32Cos, In(AF32), Out(AF32))
	Op(OP_F32_ASIN, "f32.asin", opF32Asin, In(AF32), Out(AF32))
	Op(OP_F32_SIN, "f32.sin", opF32Sin, In(AF32), Out(AF32))
	Op(OP_F32_SQRT, "f32.sqrt", opF32Sqrt, In(AF32), Out(AF32))
	Op(OP_F32_LOG, "f32.log", opF32Log, In(AF32), Out(AF32))
	Op(OP_F32_LOG2, "f32.log2", opF32Log2, In(AF32), Out(AF32))
	Op(OP_F32_LOG10, "f32.log10", opF32Log10, In(AF32), Out(AF32))
	Op(OP_F32_MAX, "f32.max", opF32Max, In(AF32, AF32), Out(AF32))
	Op(OP_F32_MIN, "f32.min", opF32Min, In(AF32, AF32), Out(AF32))
	Op(OP_F32_RAND, "f32.rand", opF32Rand, nil, Out(AF32))

	Op(OP_F64_IS_NAN, "f64.isnan", opF64Isnan, In(AF64), Out(ABOOL))
	Op(OP_F64_STR, "f64.str", opF64ToStr, In(AF64), Out(ASTR))
	Op(OP_F64_I8, "f64.i8", opF64ToI8, In(AF64), Out(AI8))
	Op(OP_F64_I16, "f64.i16", opF64ToI16, In(AF64), Out(AI16))
	Op(OP_F64_I32, "f64.i32", opF64ToI32, In(AF64), Out(AI32))
	Op(OP_F64_I64, "f64.i64", opF64ToI64, In(AF64), Out(AI64))
	Op(OP_F64_UI8, "f64.ui8", opF64ToUI8, In(AF64), Out(AUI8))
	Op(OP_F64_UI16, "f64.ui16", opF64ToUI16, In(AF64), Out(AUI16))
	Op(OP_F64_UI32, "f64.ui32", opF64ToUI32, In(AF64), Out(AUI32))
	Op(OP_F64_UI64, "f64.ui64", opF64ToUI64, In(AF64), Out(AUI64))
	Op(OP_F64_F32, "f64.f32", opF64ToF32, In(AF64), Out(AF32))
	Op(OP_F64_PRINT, "f64.print", opF64Print, In(AF64), nil)
	Op(OP_F64_ADD, "f64.add", opF64Add, In(AF64, AF64), Out(AF64))
	Op(OP_F64_SUB, "f64.sub", opF64Sub, In(AF64, AF64), Out(AF64))
	Op(OP_F64_NEG, "f64.neg", opF64Neg, In(AF64), Out(AF64))
	Op(OP_F64_MUL, "f64.mul", opF64Mul, In(AF64, AF64), Out(AF64))
	Op(OP_F64_DIV, "f64.div", opF64Div, In(AF64, AF64), Out(AF64))
	Op(OP_F64_MOD, "f32.mod", opF64Mod, In(AF64, AF64), Out(AF64))
	Op(OP_F64_ABS, "f64.abs", opF64Abs, In(AF64), Out(AF64))
	Op(OP_F64_POW, "f64.pow", opF64Pow, In(AF64, AF64), Out(AF64))
	Op(OP_F64_GT, "f64.gt", opF64Gt, In(AF64, AF64), Out(ABOOL))
	Op(OP_F64_GTEQ, "f64.gteq", opF64Gteq, In(AF64, AF64), Out(ABOOL))
	Op(OP_F64_LT, "f64.lt", opF64Lt, In(AF64, AF64), Out(ABOOL))
	Op(OP_F64_LTEQ, "f64.lteq", opF64Lteq, In(AF64, AF64), Out(ABOOL))
	Op(OP_F64_EQ, "f64.eq", opF64Eq, In(AF64, AF64), Out(ABOOL))
	Op(OP_F64_UNEQ, "f64.uneq", opF64Uneq, In(AF64, AF64), Out(ABOOL))
	Op(OP_F64_ACOS, "f64.acos", opF64Acos, In(AF64), Out(AF64))
	Op(OP_F64_COS, "f64.cos", opF64Cos, In(AF64), Out(AF64))
	Op(OP_F64_ASIN, "f64.asin", opF64Asin, In(AF64), Out(AF64))
	Op(OP_F64_SIN, "f64.sin", opF64Sin, In(AF64), Out(AF64))
	Op(OP_F64_SQRT, "f64.sqrt", opF64Sqrt, In(AF64), Out(AF64))
	Op(OP_F64_LOG, "f64.log", opF64Log, In(AF64), Out(AF64))
	Op(OP_F64_LOG2, "f64.log2", opF64Log2, In(AF64), Out(AF64))
	Op(OP_F64_LOG10, "f64.log10", opF64Log10, In(AF64), Out(AF64))
	Op(OP_F64_MAX, "f64.max", opF64Max, In(AF64, AF64), Out(AF64))
	Op(OP_F64_MIN, "f64.min", opF64Min, In(AF64, AF64), Out(AF64))
	Op(OP_F64_RAND, "f64.rand", opF64Rand, nil, Out(AF64))

	Op(OP_STR_I8, "str.i8", opStrToI8, In(ASTR), Out(AI8))
	Op(OP_STR_I16, "str.i16", opStrToI16, In(ASTR), Out(AI16))
	Op(OP_STR_I32, "str.i32", opStrToI32, In(ASTR), Out(AI32))
	Op(OP_STR_I64, "str.i64", opStrToI64, In(ASTR), Out(AI64))
	Op(OP_STR_UI8, "str.ui8", opStrToUI8, In(ASTR), Out(AUI8))
	Op(OP_STR_UI16, "str.ui16", opStrToUI16, In(ASTR), Out(AUI16))
	Op(OP_STR_UI32, "str.ui32", opStrToUI32, In(ASTR), Out(AUI32))
	Op(OP_STR_UI64, "str.ui64", opStrToUI64, In(ASTR), Out(AUI64))
	Op(OP_STR_F32, "str.f32", opStrToF32, In(ASTR), Out(AF32))
	Op(OP_STR_F64, "str.f64", opStrToF64, In(ASTR), Out(AF64))
	Op(OP_STR_PRINT, "str.print", opStrPrint, In(ASTR), nil)
	Op(OP_STR_EQ, "str.eq", opStrEq, In(ASTR, ASTR), Out(ABOOL))
	Op(OP_STR_CONCAT, "str.concat", opStrConcat, In(ASTR, ASTR), Out(ASTR))
	Op(OP_STR_SUBSTR, "str.substr", opStrSubstr, In(ASTR, AI32, AI32), Out(ASTR))
	Op(OP_STR_INDEX, "str.index", opStrIndex, In(ASTR, ASTR), Out(AI32))
	Op(OP_STR_LAST_INDEX, "str.lastindex", opStrLastIndex, In(ASTR, ASTR), Out(AI32))
	Op(OP_STR_TRIM_SPACE, "str.trimspace", opStrTrimSpace, In(ASTR), Out(ASTR))

	Op(OP_APPEND, "append", opAppend, In(Slice(TYPE_UNDEFINED), Slice(TYPE_UNDEFINED)), Out(Slice(TYPE_UNDEFINED)))
	Op(OP_RESIZE, "resize", opResize, In(Slice(TYPE_UNDEFINED), AI32), Out(Slice(TYPE_UNDEFINED)))
	Op(OP_INSERT, "insert", opInsert, In(Slice(TYPE_UNDEFINED), Slice(TYPE_UNDEFINED)), Out(Slice(TYPE_UNDEFINED)))
	Op(OP_REMOVE, "remove", opRemove, In(Slice(TYPE_UNDEFINED), AI32), Out(Slice(TYPE_UNDEFINED)))
	Op(OP_COPY, "copy", opCopy, In(Slice(TYPE_UNDEFINED), Slice(TYPE_UNDEFINED)), Out(AI32))

	Op(OP_ASSERT, "assert", opAssertValue, In(AUND, AUND, ASTR), Out(ABOOL))
	Op(OP_TEST, "test", opTest, In(AUND, AUND, ASTR), nil)
	Op(OP_PANIC, "panic", opPanic, In(AUND, AUND, ASTR), nil)
	Op(OP_PANIC_IF, "panicIf", opPanicIf, In(ABOOL, ASTR), nil)
	Op(OP_PANIC_IF_NOT, "panicIfNot", opPanicIfNot, In(ABOOL, ASTR), nil)
	Op(OP_STRERROR, "strerror", opStrError, In(AI32), Out(ASTR))

	Op(OP_AFF_PRINT, "aff.print", opAffPrint, In(Slice(TYPE_AFF)), nil)
	Op(OP_AFF_QUERY, "aff.query", opAffQuery, In(Slice(TYPE_AFF)), Out(Slice(TYPE_AFF)))
	Op(OP_AFF_ON, "aff.on", opAffOn, In(Slice(TYPE_AFF), Slice(TYPE_AFF)), nil)
	Op(OP_AFF_OF, "aff.of", opAffOf, In(Slice(TYPE_AFF), Slice(TYPE_AFF)), nil)
	Op(OP_AFF_INFORM, "aff.inform", opAffInform, In(Slice(TYPE_AFF), AI32, Slice(TYPE_AFF)), nil)
	Op(OP_AFF_REQUEST, "aff.request", opAffRequest, In(Slice(TYPE_AFF), AI32, Slice(TYPE_AFF)), nil)

	Op(OP_HTTP_SERVE, "http.Serve", opHTTPServe, In(ASTR), Out(ASTR))
	Op(OP_HTTP_LISTEN_AND_SERVE, "http.ListenAndServe", opHTTPListenAndServe, In(ASTR), Out(ASTR))
	Op(OP_HTTP_NEW_REQUEST, "http.NewRequest", opHTTPNewRequest, In(ASTR, ASTR, ASTR), Out(ASTR))
	Op(OP_HTTP_DO, "http.Do", opHTTPDo, In(AUND), Out(AUND, ASTR))
	Op(OP_DMSG_DO, "http.DmsgDo", opDMSGDo, In(AUND), Out(ASTR))

	// Op(OP_EVOLVE_EVOLVE, "evolve.evolve", opEvolve, In(Slice(TYPE_AFF), Slice(TYPE_AFF), Slice(TYPE_F64), Slice(TYPE_F64), AI32, AI32, AI32, AF64), nil)
	// Op(OP_EVOLVE_EVOLVE, "evolve.evolve", opEvolve, In(Slice(TYPE_AFF), Slice(TYPE_AFF), Slice(TYPE_AFF), Slice(TYPE_AFF), Slice(TYPE_AFF), AI32, AI32, AI32, AF64), nil)

	Op(OP_HTTP_HANDLE, "http.Handle", opHTTPHandle,
		In(
			ASTR,
			ParamEx(ParamData{typCode: TYPE_FUNC, pkg: httpPkg, inputs: In(MakeArgument("ResponseWriter", "", -1).AddType(TypeNames[TYPE_STR]), Pointer(Struct("http", "Request", "r")))})),
		Out())

	Op(OP_HTTP_CLOSE, "http.Close", opHTTPClose, nil, nil)
}
