package base

const NON_ASSIGN_PREFIX = "nonAssign"
const CORE_MODULE = "core"
var BASIC_TYPES []string = []string{
	"bool", "str", "byte", "i32", "i64", "f32", "f64",
	"[]bool", "[]byte", "[]i32", "[]i64", "[]f32", "[]f64",
}
var NATIVE_FUNCTIONS = map[string]bool{
	"i32.add":true, "i32.mul":true, "i32.sub":true, "i32.div":true,
	"i64.add":true, "i64.mul":true, "i64.sub":true, "i64.div":true,
	"f32.add":true, "f32.mul":true, "f32.sub":true, "f32.div":true,
	"f64.add":true, "f64.mul":true, "f64.sub":true, "f64.div":true,
	"i32.mod":true, "i64.mod":true,
	"i32.and":true, "i32.or":true, "i32.xor":true, "i32.andNot":true,
	"i64.and":true, "i64.or":true, "i64.xor":true, "i64.andNot":true,

	"str.print":true, "byte.print":true, "i32.print":true, "i64.print":true,
	"f32.print":true, "f64.print":true, "[]byte.print":true, "[]i32.print":true,
	"[]i64.print":true, "[]f32.print":true, "[]f64.print":true, "bool.print":true,
	"[]bool.print":true,

	"str.id":true, "bool.id":true, "byte.id":true, "i32.id":true, "i64.id":true, "f32.id":true, "f64.id":true,
	"[]bool.id":true, "[]byte.id":true, "[]i32.id":true, "[]i64.id":true, "[]f32.id":true, "[]f64.id":true,

	"[]bool.make":true, "[]byte.make":true, "[]i32.make":true,
	"[]i64.make":true, "[]f32.make":true, "[]f64.make":true,

	"[]bool.read":true, "[]bool.write":true,
	"[]byte.read":true, "[]byte.write":true, "[]i32.read":true, "[]i32.write":true,
	"[]f32.read":true, "[]f32.write":true, "[]f64.read":true, "[]f64.write":true,
	"[]bool.len":true, "[]byte.len":true, "[]i32.len":true, "[]i64.len":true,
	"[]f32.len":true, "[]f64.len":true,

	"[]byte.str":true, "str.[]byte":true,
	
	"byte.i32":true, "byte.i64":true, "byte.f32":true, "byte.f64":true,
	"[]byte.[]i32":true, "[]byte.[]i64":true, "[]byte.[]f32":true, "[]byte.[]f64":true,

	"i32.byte":true, "i64.byte":true, "f32.byte":true, "f64.byte":true,
	"[]i32.[]byte":true, "[]i64.[]byte":true, "[]f32.[]byte":true, "[]f64.[]byte":true,

	"i64.i32":true, "f32.i32":true, "f64.i32":true,
	"i32.i64":true, "f32.i64":true, "f64.i64":true,
	"i32.f32":true, "i64.f32":true, "f64.f32":true,
	"i32.f64":true, "i64.f64":true, "f32.f64":true,

	"[]i64.[]i32":true, "[]f32.[]i32":true, "[]f64.[]i32":true,
	"[]i32.[]i64":true, "[]f32.[]i64":true, "[]f64.[]i64":true,
	"[]i32.[]f32":true, "[]i64.[]f32":true, "[]f64.[]f32":true,
	"[]i32.[]f64":true, "[]i64.[]f64":true, "[]f32.[]f64":true,
	
	"i32.lt":true, "i32.gt":true, "i32.eq":true, "i32.lteq":true, "i32.gteq":true,
	"i64.lt":true, "i64.gt":true, "i64.eq":true, "i64.lteq":true, "i64.gteq":true,
	"f32.lt":true, "f32.gt":true, "f32.eq":true, "f32.lteq":true, "f32.gteq":true,
	"f64.lt":true, "f64.gt":true, "f64.eq":true, "f64.lteq":true, "f64.gteq":true,
	"str.lt":true, "str.gt":true, "str.eq":true, "str.lteq":true, "str.gteq":true,
	"byte.lt":true, "byte.gt":true, "byte.eq":true, "byte.lteq":true, "byte.gteq":true,

	"i32.rand":true, "i64.rand":true,

	"and":true, "or":true, "not":true,
	"sleep":true, "halt":true, "goTo":true, "baseGoTo":true,

	"setClauses":true, "addObject":true, "setQuery":true,
	"remObject":true, "remObjects":true,

	"remExpr":true, "remArg":true, "addExpr":true, "affExpr":true,

	"serialize":true, "deserialize":true, "evolve":true,

	"initDef":true,
}
// var NATIVE_FUNCTIONS = []string{
// 	"i32.add", "i32.mul", "i32.sub", "i32.div",
// 	"i64.add", "i64.mul", "i64.sub", "i64.div",
// 	"f32.add", "f32.mul", "f32.sub", "f32.div",
// 	"f64.add", "f64.mul", "f64.sub", "f64.div",
// 	"i32.mod", "i64.mod",
// 	"i32.and", "i32.or", "i32.xor", "i32.andNot",
// 	"i64.and", "i64.or", "i64.xor", "i64.andNot",

// 	"str.print", "byte.print", "i32.print", "i64.print",
// 	"f32.print", "f64.print", "[]byte.print", "[]i32.print",
// 	"[]i64.print", "[]f32.print", "[]f64.print", "bool.print",
// 	"[]bool.print",

// 	"str.id", "bool.id", "byte.id", "i32.id", "i64.id", "f32.id", "f64.id",
// 	"[]bool.id", "[]byte.id", "[]i32.id", "[]i64.id", "[]f32.id", "[]f64.id",

// 	"[]bool.make", "[]byte.make", "[]i32.make",
// 	"[]i64.make", "[]f32.make", "[]f64.make",

// 	"[]bool.read", "[]bool.write",
// 	"[]byte.read", "[]byte.write", "[]i32.read", "[]i32.write",
// 	"[]f32.read", "[]f32.write", "[]f64.read", "[]f64.write",
// 	"[]bool.len", "[]byte.len", "[]i32.len", "[]i64.len",
// 	"[]f32.len", "[]f64.len",

// 	"[]byte.str", "str.[]byte",
	
// 	"byte.i32", "byte.i64", "byte.f32", "byte.f64",
// 	"[]byte.[]i32", "[]byte.[]i64", "[]byte.[]f32", "[]byte.[]f64",

// 	"i32.byte", "i64.byte", "f32.byte", "f64.byte",
// 	"[]i32.[]byte", "[]i64.[]byte", "[]f32.[]byte", "[]f64.[]byte",

// 	"i64.i32", "f32.i32", "f64.i32",
// 	"i32.i64", "f32.i64", "f64.i64",
// 	"i32.f32", "i64.f32", "f64.f32",
// 	"i32.f64", "i64.f64", "f32.f64",

// 	"[]i64.[]i32", "[]f32.[]i32", "[]f64.[]i32",
// 	"[]i32.[]i64", "[]f32.[]i64", "[]f64.[]i64",
// 	"[]i32.[]f32", "[]i64.[]f32", "[]f64.[]f32",
// 	"[]i32.[]f64", "[]i64.[]f64", "[]f32.[]f64",
	
// 	"i32.lt", "i32.gt", "i32.eq", "i32.lteq", "i32.gteq",
// 	"i64.lt", "i64.gt", "i64.eq", "i64.lteq", "i64.gteq",
// 	"f32.lt", "f32.gt", "f32.eq", "f32.lteq", "f32.gteq",
// 	"f64.lt", "f64.gt", "f64.eq", "f64.lteq", "f64.gteq",
// 	"str.lt", "str.gt", "str.eq", "str.lteq", "str.gteq",
// 	"byte.lt", "byte.gt", "byte.eq", "byte.lteq", "byte.gteq",

// 	"i32.rand", "i64.rand",

// 	"and", "or", "not",
// 	"sleep", "halt", "goTo", "baseGoTo",

// 	"setClauses", "addObject", "setQuery",
// 	"remObject", "remObjects",

// 	"remExpr", "remArg", "addExpr", "affExpr",

// 	"serialize", "deserialize", "evolve",

// 	"initDef",
// }

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
	//Heap *[]byte
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
	Typ string
	Value *[]byte
	// Offset int
	// Size int

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
	Typ string
}

// type CXType struct {
// 	Name string
// }

/*
  Functions
*/

type CXFunction struct {
	Name string
	Inputs []*CXParameter
	Outputs []*CXParameter
	Expressions []*CXExpression

	// for optimization
	NumberOutputs int

	CurrentExpression *CXExpression
	Module *CXModule
	Context *CXProgram
}

type CXParameter struct {
	Name string
	Typ string
}

type CXExpression struct {
	Operator *CXFunction
	Arguments []*CXArgument
	OutputNames []*CXDefinition
	Line int
	FileLine int
	Tag string
	
	Function *CXFunction
	Module *CXModule
	Context *CXProgram
}

type CXArgument struct {
	Typ string
	Value *[]byte
	// Offset int
	// Size int
}

/*
  Affordances
*/

type CXAffordance struct {
	Description string
	Action func()
}
