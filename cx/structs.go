package base

/*
  Root Program
*/

// type Data []byte
// type Heap []byte
// type Stack []byte

type CXProgram struct {
	Packages       []*CXPackage
	CurrentPackage *CXPackage

	Inputs  []*CXArgument
	Outputs []*CXArgument

	CallStack   []CXCall
	CallCounter int

	// Stacks []CXStack
	// Heap   CXHeap
	// Data   Data

	Memory []byte
	HeapPointer int
	StackPointer int

	HeapStartsAt int

	Terminated bool

	Path  string
	Steps [][]CXCall
}

// type CXHeap struct {
// 	Heap        Heap
// 	HeapPointer int

// 	Program *CXProgram
// }

// type CXStack struct {
// 	Stack        Stack
// 	StackPointer int

// 	Program *CXProgram
// }

type CXCall struct {
	Operator     *CXFunction
	Line         int
	FramePointer int

	State         []*CXArgument
	ReturnAddress *CXCall

	Package *CXPackage
	Program *CXProgram
}

/*
  Packages
*/

type CXPackage struct {
	// Index int
	Name      string
	Imports   []*CXPackage
	Functions []*CXFunction
	Structs   []*CXStruct
	Globals   []*CXArgument

	CurrentFunction *CXFunction
	CurrentStruct   *CXStruct
	Program         *CXProgram
}

/*
  Structs
*/

type CXStruct struct {
	// Index int
	Name   string
	Fields []*CXArgument
	Size   int

	Package *CXPackage
	Program *CXProgram
}

/*
  Functions
*/

type CXFunction struct {
	// Index int
	Name        string
	Inputs      []*CXArgument
	Outputs     []*CXArgument
	Expressions []*CXExpression
	Size        int // automatic memory size
	Length      int // number of expressions, pre-computed for performance

	ListOfPointers []*CXArgument
	NumberOutputs  int

	IsNative bool
	OpCode   int

	CurrentExpression *CXExpression
	Package           *CXPackage
	Program           *CXProgram
}

type CXExpression struct {
	// Index int
	Operator *CXFunction
	Inputs   []*CXArgument
	Outputs  []*CXArgument
	// debugging
	Line     int
	FileLine int
	FileName string

	// used for jmp statements
	Label     string
	ThenLines int
	ElseLines int

	IsMethodCall bool
	IsStructLiteral bool
	IsArrayLiteral  bool
	IsFlattened bool // used for nested struct literals

	Function *CXFunction
	Package  *CXPackage
	Program  *CXProgram
}

type CXConstant struct {
	// native constants. only used for pre-packaged constants (e.g. math package's PI)
	// these fields are used to feed WritePrimary
	Type  int
	Value []byte
}

type CXArgument struct {
	Name        string
	Type        int
	CustomType  *CXStruct
	Size        int // size of underlaying basic type
	TotalSize   int // total size of an array, performance reasons
	PointeeSize int

	Offset      int
	// HeapOffset  int
	// OffsetOffset int // for struct fields

	IndirectionLevels     int
	DereferenceLevels     int
	Pointee               *CXArgument
	PointeeMemoryType     int
	DereferenceOperations []int // offset by array index, struct field, pointer
	DeclarationSpecifiers []int // used to determine finalSize

	IsSlice      bool
	IsArray      bool
	IsArrayFirst bool // and then dereference
	IsPointer    bool
	IsReference  bool
	// IsDereference bool
	IsDereferenceFirst bool // and then array
	IsStruct           bool
	IsField            bool
	IsRest             bool // pkg.var <- var is rest
	IsLocalDeclaration bool
	IsShortDeclaration bool

	PassBy int  // pass by value or reference
	DoesEscape bool

	// Sizes []int // used to access struct fields
	Lengths []int // declared lengths at compile time
	// NumIndexes int // how many levels we'll go deep. NumIndexes <= len(Lengths)
	Indexes []*CXArgument
	Fields  []*CXArgument // strct.fld1.fld2().fld3

	SynonymousTo string // when the symbol is just a temporary holder for another symbol

	FileLine int
	FileName string
	
	Package *CXPackage
	Program *CXProgram

	// Value *[]byte
	// Typ   string
}

/*
  Affordances
*/

type CXAffordance struct {
	Description string
	Operator    string
	Name        string
	Typ         string
	Index       string
	Action      func()
}
