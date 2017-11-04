package base

const NON_ASSIGN_PREFIX = "nonAssign"
const CORE_MODULE = "core"
var BASIC_TYPES []string = []string{
	"bool", "str", "byte", "i32", "i64", "f32", "f64",
	"[]bool", "[]str", "[]byte", "[]i32", "[]i64", "[]f32", "[]f64",
}
var NATIVE_FUNCTIONS = map[string]bool{
	"i32.add":true, "i32.mul":true, "i32.sub":true, "i32.div":true,
	"i64.add":true, "i64.mul":true, "i64.sub":true, "i64.div":true,
	"f32.add":true, "f32.mul":true, "f32.sub":true, "f32.div":true,
	"f64.add":true, "f64.mul":true, "f64.sub":true, "f64.div":true,
	"i32.mod":true, "i64.mod":true,
	"i32.bitand":true, "i32.bitor":true, "i32.bitxor":true, "i32.bitclear":true,
	"i64.bitand":true, "i64.bitor":true, "i64.bitxor":true, "i64.bitclear":true,

	"str.print":true, "byte.print":true, "i32.print":true, "i64.print":true,
	"f32.print":true, "f64.print":true, "[]byte.print":true, "[]i32.print":true,
	"[]i64.print":true, "[]f32.print":true, "[]f64.print":true, "bool.print":true,
	"[]bool.print":true, "[]str.print":true,

	"str.id":true, "bool.id":true, "byte.id":true, "i32.id":true, "i64.id":true, "f32.id":true, "f64.id":true,
	"[]bool.id":true, "[]byte.id":true, "[]str.id":true, "[]i32.id":true, "[]i64.id":true, "[]f32.id":true, "[]f64.id":true,
	"identity":true,

	"[]bool.make":true, "[]byte.make":true, "[]str.make":true,
	"[]i32.make":true, "[]i64.make":true, "[]f32.make":true, "[]f64.make":true,
	

	"[]bool.read":true, "[]bool.write":true, "[]byte.read":true, "[]byte.write":true,
	"[]str.read":true, "[]str.write":true, "[]i32.read":true, "[]i32.write":true,
	"[]i64.read":true, "[]i64.write":true,
	"[]f32.read":true, "[]f32.write":true, "[]f64.read":true, "[]f64.write":true,
	
	"[]bool.len":true, "[]byte.len":true, "[]i32.len":true, "[]i64.len":true,
	"[]f32.len":true, "[]f64.len":true, "[]str.len":true,

	"str.concat":true, "[]byte.concat":true, "[]bool.concat":true, "[]str.concat":true,
	"[]i32.concat":true, "[]i64.concat":true, "[]f32.concat":true, "[]f64.concat":true,

	"[]byte.append":true, "[]bool.append":true, "[]str.append":true,
	"[]i32.append":true, "[]i64.append":true, "[]f32.append":true, "[]f64.append":true,
	
	"[]byte.copy":true, "[]bool.copy":true, "[]str.copy":true,
	"[]i32.copy":true, "[]i64.copy":true, "[]f32.copy":true, "[]f64.copy":true,
	
	"[]byte.str":true, "str.[]byte":true,
	
	"byte.i32":true, "byte.i64":true, "byte.f32":true, "byte.f64":true,
	"[]byte.[]i32":true, "[]byte.[]i64":true, "[]byte.[]f32":true, "[]byte.[]f64":true,

	"i32.byte":true, "i64.byte":true, "f32.byte":true, "f64.byte":true,
	"[]i32.[]byte":true, "[]i64.[]byte":true, "[]f32.[]byte":true, "[]f64.[]byte":true,

	"i64.i32":true, "f32.i32":true, "f64.i32":true,
	"i32.i64":true, "f32.i64":true, "f64.i64":true,
	"i32.f32":true, "i64.f32":true, "f64.f32":true,
	"i32.f64":true, "i64.f64":true, "f32.f64":true,

	"byte.str":true, "bool.str":true, "i32.str":true,
	"i64.str":true, "f32.str":true, "f64.str":true,

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

	"str.read":true, "i32.read":true,

	"i32.rand":true, "i64.rand":true,

	"and":true, "or":true, "not":true,
	"sleep":true, "halt":true, "goTo":true, "baseGoTo":true,



	"setClauses":true, "addObject":true, "setQuery":true,
	"remObject":true, "remObjects":true,

	"remExpr":true, "remArg":true, "addExpr":true, "affExpr":true,

	"aff.query":true, "aff.execute":true, "aff.print":true, "aff.concat":true,
	"aff.len":true,

	"serialize":true, "deserialize":true, "evolve":true,

	"initDef":true,

	"test.start":true, "test.stop":true,
	"test.error":true, "test.bool":true, "test.str":true, "test.byte":true,
	"test.i32":true, "test.i64":true, "test.f32":true, "test.f64":true,
	"test.[]bool":true, "test.[]byte":true, "test.[]str":true, "test.[]i32":true,
	"test.[]i64":true, "test.[]f32":true, "test.[]f64":true,

	"cstm.append":true, "cstm.read":true, "cstm.len":true,

	/*
          Runtime
        */

	"runtime.LockOSThread":true,
	
	/*
          OpenGL
        */
	"gl.Init":true, "gl.CreateProgram":true, "gl.LinkProgram":true,
	"gl.Clear":true, "gl.UseProgram":true,
	
	"gl.BindBuffer":true, "gl.BindVertexArray":true, "gl.EnableVertexAttribArray":true,
	"gl.VertexAttribPointer":true, "gl.DrawArrays":true, "gl.GenBuffers":true,
	"gl.BufferData":true, "gl.GenVertexArrays":true, "gl.CreateShader":true,
	
	"gl.Strs":true, "gl.Free":true, "gl.ShaderSource":true,
	"gl.CompileShader":true, "gl.GetShaderiv":true, "gl.AttachShader":true,

	"gl.MatrixMode":true,
	"gl.Rotatef":true, "gl.Translatef":true, "gl.LoadIdentity":true,
	"gl.PushMatrix":true, "gl.PopMatrix":true, "gl.EnableClientState":true,

	"gl.BindTexture":true, "gl.Color4f":true, "gl.Begin":true,
	"gl.End":true, "gl.Normal3f":true, "gl.TexCoord2f":true,
	"gl.Vertex3f":true,

	"gl.Enable":true, "gl.ClearColor":true, "gl.ClearDepth":true,
	"gl.DepthFunc":true, "gl.Lightfv":true, "gl.Frustum":true,
	"gl.Disable":true, "gl.Hint":true,

	"gl.NewTexture":true, "gl.DepthMask":true, "gl.TexEnvi":true,
	"gl.BlendFunc":true,
	
	/*
          GLFW
        */

	"glfw.Init":true, "glfw.WindowHint":true, "glfw.CreateWindow":true,
	"glfw.MakeContextCurrent":true, "glfw.ShouldClose":true,
	"glfw.PollEvents":true, "glfw.SwapBuffers":true,

	"glfw.SetKeyCallback":true,

	/*
          Operating System
        */

	"os.Create":true, "os.Open":true, "os.Close":true,
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
	Operator string
	Name string
	Typ string
	Index string
	Action func()
}
