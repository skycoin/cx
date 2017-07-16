package base

// this is how we would serialize the program

// cxExpression
// OperatorOffset
// OperatorNumber // to know how many operators are we going to read
// ArgumentsOffest
// ArgumentsNumber
// OutputNameOffset
// OutputNameNumber // how many bytes are we going to read
// Line // this is just an int
// FunctionOffset
// ModuleOffset
// ContextOffset // I don't know how to handle this. It could be the index 0 of the program

// Okay, I think this is a valid strategy
// For every singular field, we just store the offset
// For arrys, we store the offset and the size of the chunk

// from a context
// modules := getModules(cxt, offset, size)
// mainMod := getModule(modules, offset)
// mainFn := getFunction(mainMod, offset)
// mainExprs := getExpressions(mainFn, offset, size)
// firstExpr := getExpression(mainFn, offset)

// each of these functions return array of bytes
// so we are just serializing structs that tell us offset and sizes
// and we use these functions to get specific sub arrays of bytes

// to serialize, we can just iterate though all the elements and transform everything to
// byte arrays, but the important thing is *order*
// BUT, as we know that the first bytes are always going to be the context, we can deserialize
// those bytes first. These bytes will then tell us how to deserialize everything else.
// As we know the size, we also know how many of each elements are stored in the context
// so we can easily allocate an array of the correct size for this
// We can have a function which is reading all the bytes and allocating the structs with
// its arrays.
// If we wanted to keep the current structs format, we can easily transform byte arrays to strings (for names, for example) and arrays to maps



// ***
// Maybe this is important
// The calls in the call stack hold, well, the stack for each subroutine
// We need to create a heap field in the cxContext
// These heap will be holding all the values we are going to be substituting
// ***




// let's now think a little bit more about compilation
// what we need to worry about is stuff in the heap, static stuff
// stuff in the stack is dynamically allocated, so we can't do too much in here (or can we?)
// heap first:
// maybe when we are serializing a program we could replace

// result := strct.x + strct.y        <- I think the biggest problem is when we have nested types


// for example, if we have a struct with a field of type of other field, which has a type of other...
// hmmm, no, because that's the caller's job: to access the required fields of each field if that makes sense

// we would need to access a map, which is O(N log N) (it's fast, but not crazy fast; we want O(1))
// if we were using arrays, we would need to lookup for every single one, until we find it, so it would be O(n) (which is worse than a map)
// what we want to do is replace struct.x for metaStruct[2], which means


// could we group all the fields together? instead of grouping the whole structure



// as we need to *substitute*, could we just remove, for example, a cxArgument
// and replace it by another which is faster?

// currently, arguments can only hold variables (ident) and literals (i32)
// to access a struct (e.g. struct.x), this would still be an identifier
// we just substitute the value for its i32 part (this would be inline), because
// the value could change due to another function manipulating that value
// so we would need to replace it by another argument type
// for example, we could create an argument type "struct"
// internally, we would read this "struct" type, just like ident,
// and we would access
// nope, this doesn't work

// okay, back to the substituting part
// Everytime we call AddStruct we would create its corresponding structMeta

// this is for compilation
// no, we need to group the whole structure, and use word-boundaries like in C
// use paddings to fill the bytes






// currently, we can already do:
// serialization of a program
// deserialization of a program
// maybe we need to transform maps back to arrays in our structs


type metaExpression struct {
	offset int
	size int // do we need the size? we know the type because it's a metaExpression
	// and the size must be constant
}

var arguments []byte
// this slice would hold all the arguments
// at runtime, we know exactly how many arguments there are
// so we would just do a arguments = make([]arguments, 1000 * SizeOf(cxArgument))

func getArgument (args []byte, offset int) {
	// we use the offset and the known size for an argument
	// to get the desired argument (bytes)
	// I don't know if we should cast it in here
}

type metaArgument struct {
	offset int
	// an argument is always going to be, for example, 8 bytes, 4 which hold
}

// let's see, if I have Argument.Name, how do I compile this
// no, this isn't needed at compile time, because this is how we are defining our program
// a better example is, if I have strct.num, how do I compile this
// strct.num should be replaced by, for example metaStrct[313] or something like this
// well, something that we can do is modify the program structure to just "globals"
// no, we don't need to mess up with the program structure

// okay, at the first step, we're going to have just a map (when handling with interpreted programs)
// if we compile this program
//





// okay, we know that everything concerned to serializing the program is easy (let's assume)
// now, for compiling, let's see what is going to be the problem
// brandon said before that the biggest concern is to acces a structure's field
// let's get back to just serializing a program...






// okay, looking at metaArgument, it only gives us an offset


// now the more complex elements:

// An expression:
// The operator let's just leave it as a pointer (metaOperators[55]) for now
// Arguments: Okay! An argument was just a bunch of bytes
// Wait, let's return to the basics above


// with unsafe:






// Affordances would need to be applied before compilation (or maybe not? we could make structMeta grow automatically)
/// If it grows automatically, we just need to make a new array + the size of the affordance (this will allocate a contiguous block of memory)
// We could have an array, because this is going to be performed at compile time
// Every relevant structure just needs to save an offset














// used for affordances (and maybe other stuff)
var basicTypes = []string{"i32"}
var basicFunctions = []string{"addI32", "mulI32", "subI32", "divI32"}

/*
  Context
*/

// type cxContext struct {
// 	Modules map[string]*cxModule
// 	CurrentModule *cxModule
// 	CallStack []*cxCall
// 	Steps [][]*cxCall
// 	ProgramSteps []*cxProgramStep
// }

type cxContext struct {
	Modules map[string]*cxModule
	CurrentModule *cxModule
	CallStack *cxCallStack
	Steps []*cxCallStack
	ProgramSteps []*cxProgramStep
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

type cxArgument struct {
	Typ *cxType
	Value *[]byte
}

/*
  Affordances
*/

type cxAffordance struct {
	Description string
	Action func()
}
