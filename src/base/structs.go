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
	"printI64A", "printF32A", "printF64A", "printBool",
	
	"idStr", "idByte", "idI32", "idI64", "idF32", "idF64",
	"idByteA", "idI32A", "idI64A", "idF32A", "idF64A",
	
	"readByteA", "writeByteA", "readI32A", "writeI32A",
	"readF32A", "writeF32A", "readF64A", "writeF64A",
	
	"byteAToStr",
	"i64ToI32", "f32ToI32", "f64ToI32",
	"i32ToI64", "f32ToI64", "f64ToI64",
	"i32ToF32", "i64ToF32", "f64ToF32",
	"i32ToF64", "i64ToF64", "f32ToF64",
	
	"ltI32", "gtI32", "eqI32",
	"ltI64", "gtI64", "eqI64",

	"sleep",

	"initDef",
	"evolve",
	"goTo",
}

/*
  Context
*/

type CXProgram struct {
	Modules []*CXModule
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
	State []*CXDefinition
	ReturnAddress *CXCall
	Context *CXProgram
	Module *CXModule
}

type CXProgramStep struct {
	//Action func(*CXProgram)
	Action func()
}

/*
  Modules
*/

type CXModule struct {
	Name string
	Imports []*CXModule
	Functions []*CXFunction
	Structs []*CXStruct
	Definitions []*CXDefinition

	// Affordance inference
	//Clauses []*CXClause
	Clauses string
	Objects []string
	Query string
	//Objects []*CXArgument // Idents

	CurrentFunction *CXFunction
	CurrentStruct *CXStruct
	Context *CXProgram
}

type CXClause struct {
	//Type *CXType
	Operator *CXFunction
	Argument *CXArgument
	Object *CXArgument

	Module *CXModule
	Context *CXProgram
}

type CXDefinition struct {
	Name string
	Typ *CXType
	Value *[]byte
	Offset int
	Size int

	Module *CXModule
	Context *CXProgram
}

/*
  Structs
*/

type CXStruct struct {
	Name string
	Fields []*CXField

	Module *CXModule
	Context *CXProgram
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
	Context *CXProgram
}

type CXParameter struct {
	Name string
	Typ *CXType
}

type CXExpression struct {
	Operator *CXFunction
	Arguments []*CXArgument
	OutputNames []*CXDefinition
	Line int
	
	Function *CXFunction
	Module *CXModule
	Context *CXProgram
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
