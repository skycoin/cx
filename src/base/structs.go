package base

type cxType struct {
	Name string
}
// used for affordances (and maybe other stuff)
var basicTypes = []string{"i32"}
var basicFunctions = []string{"addI32", "mulI32", "subI32", "divI32"}

type cxField struct {
	Name string
	Typ *cxType
}


// Are all these structures going to be c structs?

// Affordances would need to be applied before compilation (or maybe not? we could make structMeta grow automatically)
/// If it grows automatically, we just need to make a new array + the size of the affordance (this will allocate a contiguous block of memory)

// We need a byte array to store the fields (the values)
// If we want to access field "x", how do we know

// We need a structMeta[] which would be an array (or slice?)
// We could have an array, because this is going to be performed at compile time

// type cxStructC struct {
// 	fields []*cxField
// 	offset int
// }

// Do we need a structMeta for each of our structs?
// We'll use word-length offsets (for example, 4 bytes) for fields of primitive types
// What about fields of complex types?
/// For example, in an expression we have (+ Point1.x Point2.y)


// Compilation (substition, when):

// Benchmark (way to test it):

// Every relevant structure just needs to save an offset
type structMeta struct {
	//Offsets []
}

type cxStruct struct {
	Name string
	Fields []*cxField

	Module *cxModule
	Context *cxContext
}

/*
  Context
*/

type cxContext struct {
	Modules map[string]*cxModule
	CurrentModule *cxModule
	CallStack []*cxCall
	Steps [][]*cxCall
	ProgramSteps []*cxProgramStep
}

/*
  Functions
*/

type cxParameter struct {
	Name string
	Typ *cxType
}

type cxArgument struct {
	Typ *cxType
	Value *[]byte
}

type cxCall struct {
	Operator *cxFunction
	Line int
	State map[string]*cxDefinition
	ReturnAddress *cxCall
	Context *cxContext
	Module *cxModule
}

// We could somehow use the same cxCall process
// Operator could be

// The affordances option:
// The affordance could receive a context

type cxExpression struct {
	Operator *cxFunction
	Arguments []*cxArgument
	OutputName string
	Line int
	
	Function *cxFunction
	Module *cxModule
	Context *cxContext
}

type cxFunction struct {
	Name string
	Inputs []*cxParameter
	Output *cxParameter
	Expressions []*cxExpression

	CurrentExpression *cxExpression
	Module *cxModule
	Context *cxContext
}

/*
  Modules
*/

type cxDefinition struct {
	Name string
	Typ *cxType
	Value *[]byte

	Module *cxModule
	Context *cxContext
}

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

/*
  Affordances
*/

type cxAffordance struct {
	Description string
	Action func()
}


type cxProgramStep struct {
	Action func(*cxContext)
}
