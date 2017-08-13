package base

const NON_ASSIGN_PREFIX = "nonAssign"
const CORE_MODULE = "core"
var BASIC_TYPES []string = []string{
	"bool", "str", "byte", "i32", "i64", "f32", "f64",
	"[]bool", "[]byte", "[]i32", "[]i64", "[]f32", "[]f64",
}
var NATIVE_FUNCTIONS = []string{
	"addI32", "mulI32", "subI32", "divI32",
	"addI64", "mulI64", "subI64", "divI64",
	"addF32", "mulF32", "subF32", "divF32",
	"addF64", "mulF64", "subF64", "divF64",
	
	"printStr", "printByte", "printI32", "printI64",
	"printF32", "printF64", "printByteA", "printI32A",
	"printI64A", "printF32A", "printF64A",
	
	"idStr", "idByte", "idI32", "idI64", "idF32", "idF64",
	"idByteA", "idI32A", "idI64A", "idF32A", "idF64A",
	
	"readAByte", "writeAByte",
	
	"byteAtoStr", "i32toI64", "f32toI64", "f64toI64",
	
	"ltI32", "gtI32", "eqI32",
	"ltI64", "gtI64", "eqI64",

	"initDef",
	"evolve",
	"goTo",
}
var ARRAY_FUNCTIONS = []string{
	"readAByte", "writeAByte",
}

/*
  Context
*/

type CXContext struct {
	Modules map[string]*CXModule
	CurrentModule *CXModule
	CallStack *CXCallStack
	Outputs []*CXDefinition
	Steps []*CXCallStack
	ProgramSteps []*CXProgramStep
	Heap *[]byte
}

type CXCallStack struct {
	Calls []*CXCall
}

type CXCall struct {
	Operator *CXFunction
	Line int
	State map[string]*CXDefinition
	ReturnAddress *CXCall
	Context *CXContext
	Module *CXModule
}

type CXProgramStep struct {
	Action func(*CXContext)
}

/*
  Modules
*/

type CXModule struct {
	Name string
	Imports map[string]*CXModule
	Functions map[string]*CXFunction
	Structs map[string]*CXStruct
	Definitions map[string]*CXDefinition

	CurrentFunction *CXFunction
	CurrentStruct *CXStruct
	Context *CXContext
}

type CXDefinition struct {
	Name string
	Typ *CXType
	Value *[]byte
	Offset int
	Size int

	Module *CXModule
	Context *CXContext
}

/*
  Structs
*/

type CXStruct struct {
	Name string
	Fields []*CXField

	Module *CXModule
	Context *CXContext
}

type CXField struct {
	Name string
	Typ *CXType
}

type CXType struct {
	Name string
}

/*
  Functions
*/

type CXFunction struct {
	Name string
	Inputs []*CXParameter
	Outputs []*CXParameter
	Expressions []*CXExpression

	CurrentExpression *CXExpression
	Module *CXModule
	Context *CXContext
}

type CXParameter struct {
	Name string
	Typ *CXType
}

type CXExpression struct {
	Operator *CXFunction
	Arguments []*CXArgument
	OutputNames []string
	Line int
	
	Function *CXFunction
	Module *CXModule
	Context *CXContext
}

type CXArgument struct {
	Typ *CXType
	Value *[]byte
	Offset int
	Size int
}

/*
  Affordances
*/

type CXAffordance struct {
	Description string
	Action func()
}
