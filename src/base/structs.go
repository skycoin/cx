package base

const MAIN_FUNC = "main"
const SYS_INIT_FUNC = "*init"
const MAIN_PKG = "main"
const NON_ASSIGN_PREFIX = "nonAssign"
const LOCAL_PREFIX = "lcl"
const CORE_MODULE = "core"
const ID_FN = "identity"
const INIT_FN = "initDef"
const SLICE_SIZE = 32
const MARK_SIZE = 1
const OBJECT_HEADER_SIZE = 9
const FORWARDING_ADDRESS_SIZE = 4
const OBJECT_SIZE = 4
const CALLSTACK_SIZE = 500000
const STACK_SIZE = 500000
const INIT_HEAP_SIZE = 500000
const NULL_HEAP_ADDRESS_OFFSET = 4
const NULL_HEAP_ADDRESS = 0

const (
	DECL_POINTER = iota // 0
	DECL_ARRAY // 1
	DECL_SLICE // 2
	DECL_STRUCT // 3
	DECL_BASIC // 4
)

const (
	MEM_READ = iota
	MEM_WRITE
)

const (
	DEREF_ARRAY = iota
	DEREF_FIELD
	DEREF_POINTER
	DEREF_DEREF
)

const TYPE_POINTER_SIZE int = 4

// types
const (
	TYPE_BOOL = iota
	TYPE_BYTE
	TYPE_STR
	TYPE_F32
	TYPE_F64
	TYPE_I8
	TYPE_I16
	TYPE_I32
	TYPE_I64
	TYPE_UI8
	TYPE_UI16
	TYPE_UI32
	TYPE_UI64

	TYPE_THRESHOLD
	
	TYPE_UNDEFINED
	TYPE_CUSTOM
	TYPE_POINTER
	TYPE_IDENTIFIER
)

var TypeCounter int
var TypeCodes map[string]int = map[string]int{
	"identifier": TYPE_IDENTIFIER,
	"bool": TYPE_BOOL,
	"byte": TYPE_BYTE,
	"str": TYPE_STR,
	"f32": TYPE_F32,
	"f64": TYPE_F64,
	"i8": TYPE_I8,
	"i16": TYPE_I16,
	"i32": TYPE_I32,
	"i64": TYPE_I64,
	"ui8": TYPE_UI8,
	"ui16": TYPE_UI16,
	"ui32": TYPE_UI32,
	"ui64": TYPE_UI64,
}

var TypeNames map[int]string = map[int]string{
	TYPE_IDENTIFIER: "ident",
	TYPE_BOOL: "bool",
	TYPE_BYTE: "byte",
	TYPE_STR: "str",
	TYPE_F32: "f32",
	TYPE_F64: "f64",
	TYPE_I8: "i8",
	TYPE_I16: "i16",
	TYPE_I32: "i32",
	TYPE_I64: "i64",
	TYPE_UI8: "ui8",
	TYPE_UI16: "ui16",
	TYPE_UI32: "ui32",
	TYPE_UI64: "ui64",
}

// memory locations
const (
	MEM_STACK = iota
	MEM_HEAP
	MEM_DATA
)

/*
  Context
*/

type Data []byte
type Heap []byte
type Stack []byte

type CXProgram struct {
	Packages []*CXPackage
	CurrentPackage *CXPackage

	Inputs []*CXArgument
	Outputs []*CXArgument
	
	CallStack []CXCall
	CallCounter int
	
	Stacks []CXStack
	Heap CXHeap
	Data Data

	Terminated bool


	// from interpreted
	Steps [][]CXCall

	
}

type CXHeap struct {
	Heap Heap
	HeapPointer int

	Program *CXProgram
}

type CXStack struct {
	Stack Stack
	StackPointer int

	Program *CXProgram
}

type CXCall struct {
	Operator *CXFunction
	Line int
	FramePointer int
}

/*
  Packages
*/

type CXPackage struct {
	Name string
	Imports []*CXPackage
	Functions []*CXFunction
	Structs []*CXStruct
	Globals []*CXArgument

	CurrentFunction *CXFunction
	CurrentStruct *CXStruct
	Program *CXProgram
}

/*
  Structs
*/

type CXStruct struct {
	Name string
	Fields []*CXArgument
	Size int

	Package *CXPackage
	Program *CXProgram
}

/*
  Functions
*/

type CXFunction struct {
	Name string
	Inputs []*CXArgument
	Outputs []*CXArgument
	Expressions []*CXExpression
	Size int // automatic memory size
	Length int // number of expressions, pre-computed for performance

	ListOfPointers []*CXArgument

	IsNative bool
	OpCode int

	CurrentExpression *CXExpression
	Package *CXPackage
	Program *CXProgram
}

type CXExpression struct {
	Operator *CXFunction
	Inputs []*CXArgument
	Outputs []*CXArgument
	// debugging
	FileLine int
	FileName string

	// used for jmp statements
	Label string
	ThenLines int
	ElseLines int

	IsStructLiteral bool
	IsArrayLiteral bool
	
	Function *CXFunction
	Package *CXPackage
	Program *CXProgram
}

type CXConstant struct {
	// native constants. only used for pre-packaged constants (e.g. math package's PI)
	// these fields are used to feed WritePrimary
	Type int
	Value []byte
}

type CXArgument struct {
	Name string
	Type int
	CustomType *CXStruct
	Size int // size of underlaying basic type
	TotalSize int // total size of an array, performance reasons
	PointeeSize int

	MemoryRead int // these will later be removed and a single memory pointer will be used
	MemoryWrite int
	Offset int
	HeapOffset int
	// OffsetOffset int // for struct fields

	IndirectionLevels int
	DereferenceLevels int
	Pointee *CXArgument
	PointeeMemoryType int
	DereferenceOperations []int // offset by array index, struct field, pointer
	DeclarationSpecifiers []int // used to determine finalSize

	IsArray bool
	IsArrayFirst bool // and then dereference
	IsPointer bool
	IsReference bool
	// IsDereference bool
	IsDereferenceFirst bool // and then array
	IsStruct bool
	IsField bool
	IsRest bool // pkg.var <- var is rest
	IsLocalDeclaration bool
	
	// Sizes []int // used to access struct fields
	Lengths []int // declared lengths at compile time
	// NumIndexes int // how many levels we'll go deep. NumIndexes <= len(Lengths)
	Indexes []*CXArgument
	Fields []*CXArgument // strct.fld1.fld2().fld3

	SynonymousTo string // when the symbol is just a temporary holder for another symbol

	Package *CXPackage
	Program *CXProgram

	// interpreted
	Value *[]byte
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

