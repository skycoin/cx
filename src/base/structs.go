package base

// used for affordances (and maybe other stuff)
var basicTypes = []string{"byte", "i32", "i64", "[]byte", "[]i32", "[]i64"}
var basicFunctions = []string{"addI32", "mulI32", "subI32", "divI32",
	"addI64", "mulI64", "subI64", "divI64",
	"readAByte", "writeAByte"}
var arrayFunctions = []string{"readAByte", "writeAByte"}

/*
  Context
*/

type cxContext struct {
	Modules map[string]*cxModule
	CurrentModule *cxModule
	CallStack *cxCallStack
	Output *cxDefinition
	Steps []*cxCallStack
	ProgramSteps []*cxProgramStep
	Heap *[]byte
}

type cxCallStack struct {
	Calls []*cxCall
}

type cxCall struct {
	Operator *cxFunction
	Line int
	State map[string]*cxDefinition
	ReturnAddress *cxCall
	Context *cxContext
	Module *cxModule
}

type cxProgramStep struct {
	Action func(*cxContext)
}

/*
  Modules
*/

type cxModule struct {
	Name string
	Imports map[string]*cxModule
	Functions map[string]*cxFunction
	Structs map[string]*cxStruct
	Definitions map[string]*cxDefinition

	CurrentFunction *cxFunction
	CurrentStruct *cxStruct
	Context *cxContext
}

type cxDefinition struct {
	Name string
	Typ *cxType
	Value *[]byte

	Module *cxModule
	Context *cxContext
}

/*
  Structs
*/

type cxStruct struct {
	Name string
	Fields []*cxField

	Module *cxModule
	Context *cxContext
}

type cxField struct {
	Name string
	Typ *cxType
}

type cxType struct {
	Name string
}

/*
  Functions
*/

type cxFunction struct {
	Name string
	Inputs []*cxParameter
	Output *cxParameter
	Expressions []*cxExpression

	CurrentExpression *cxExpression
	Module *cxModule
	Context *cxContext
}

type cxParameter struct {
	Name string
	Typ *cxType
}

type cxExpression struct {
	Operator *cxFunction
	Arguments []*cxArgument
	OutputName string
	Line int
	
	Function *cxFunction
	Module *cxModule
	Context *cxContext
}

// type cxPointer struct {
// 	//Typ *cxType // do we need to know the type?
// 	Offset int
// 	Size int
// }

type cxArgument struct {
	Typ *cxType
	Value *[]byte
	Offset int
	Size int
}

/*
  Affordances
*/

type cxAffordance struct {
	Description string
	Action func()
}
