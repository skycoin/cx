package base

const NON_ASSIGN_PREFIX = "nonAssign"
const CORE_MODULE = "core"
var BASIC_TYPES []string = []string{
	"bool", "str", "byte", "i32", "i64", "f32", "f64",
	"[]bool", "[]byte", "[]i32", "[]i64", "[]f32", "[]f64",
}
var NATIVE_FUNCTIONS = []string{
	"i32.add", "i32.mul", "i32.sub", "i32.div",
	"i64.add", "i64.mul", "i64.sub", "i64.div",
	"f32.add", "f32.mul", "f32.sub", "f32.div",
	"f64.add", "f64.mul", "f64.sub", "f64.div",
	"i32.mod", "i64.mod",
	"i32.and", "i32.or", "i32.xor", "i32.andNot",
	"i64.and", "i64.or", "i64.xor", "i64.andNot",

	"str.print", "byte.print", "i32.print", "i64.print",
	"f32.print", "f64.print", "[]byte.print", "[]i32.print",
	"[]i64.print", "[]f32.print", "[]f64.print", "bool.print",
	"[]bool.print",

	"str.id", "bool.id", "byte.id", "i32.id", "i64.id", "f32.id", "f64.id",
	"[]bool.id", "[]byte.id", "[]i32.id", "[]i64.id", "[]f32.id", "[]f64.id",

	"[]bool.make", "[]byte.make", "[]i32.make",
	"[]i64.make", "[]f32.make", "[]f64.make",

	"[]bool.read", "[]bool.write",
	"[]byte.read", "[]byte.write", "[]i32.read", "[]i32.write",
	"[]f32.read", "[]f32.write", "[]f64.read", "[]f64.write",
	"[]bool.len", "[]byte.len", "[]i32.len", "[]i64.len",
	"[]f32.len", "[]f64.len",

	"[]byte.str", "str.[]byte",
	
	"byte.i32", "byte.i64", "byte.f32", "byte.f64",
	"[]byte.[]i32", "[]byte.[]i64", "[]byte.[]f32", "[]byte.[]f64",

	"i32.byte", "i64.byte", "f32.byte", "f64.byte",
	"[]i32.[]byte", "[]i64.[]byte", "[]f32.[]byte", "[]f64.[]byte",

	"i64.i32", "f32.i32", "f64.i32",
	"i32.i64", "f32.i64", "f64.i64",
	"i32.f32", "i64.f32", "f64.f32",
	"i32.f64", "i64.f64", "f32.f64",

	"[]i64.[]i32", "[]f32.[]i32", "[]f64.[]i32",
	"[]i32.[]i64", "[]f32.[]i64", "[]f64.[]i64",
	"[]i32.[]f32", "[]i64.[]f32", "[]f64.[]f32",
	"[]i32.[]f64", "[]i64.[]f64", "[]f32.[]f64",
	
	"i32.lt", "i32.gt", "i32.eq", "i32.lteq", "i32.gteq",
	"i64.lt", "i64.gt", "i64.eq", "i64.lteq", "i64.gteq",
	"f32.lt", "f32.gt", "f32.eq", "f32.lteq", "f32.gteq",
	"f64.lt", "f64.gt", "f64.eq", "f64.lteq", "f64.gteq",
	"str.lt", "str.gt", "str.eq", "str.lteq", "str.gteq",
	"byte.lt", "byte.gt", "byte.eq", "byte.lteq", "byte.gteq",

	"i32.rand", "i64.rand",

	"and", "or", "not",
	"sleep", "halt", "goTo",

	"setClauses", "addObject", "setQuery",
	"remObject", "remObjects",

	"remExpr", "remArg", "addExpr", "affExpr",

	"serialize", "deserialize", "evolve",

	"initDef",
}

/*
  Context
*/

type CXProgram struct {
	Modules []*CXModule
	CurrentModule *CXModule
	CallStack *CXCallStack
	Terminated bool
	// Inputs []*CXDefinition
	Outputs []*CXDefinition
	Steps []*CXCallStack
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
	Clauses string
	Objects []*CXObject
	Query string

	CurrentFunction *CXFunction
	CurrentStruct *CXStruct
	Context *CXProgram
}

type CXObject struct {
	Name string
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
	FileLine int
	Tag string
	Label string
	
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
