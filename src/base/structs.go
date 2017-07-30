package base

// used for affordances (and maybe other stuff)
var basicTypes = []string{
	"str", "byte", "i32", "i64", "f32", "f64",
	"[]byte", "[]i32", "[]i64", "[]f64", "[]f32",
}
var basicFunctions = []string{
	"evolve",
	"addI32", "mulI32", "subI32", "divI32",
	"addI64", "mulI64", "subI64", "divI64",
	"idAI32", "idI32",
	"readAByte", "writeAByte",
}
var arrayFunctions = []string{
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

// type CXPointer struct {
// 	//Typ *CXType // do we need to know the type?
// 	Offset int
// 	Size int
// }

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
