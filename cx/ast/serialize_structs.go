package ast

type serializedCXProgramIndex struct {
	ProgramOffset     int64
	CallsOffset       int64
	PackagesOffset    int64
	StructsOffset     int64
	FunctionsOffset   int64
	ExpressionsOffset int64
	ArgumentsOffset   int64
	IntegersOffset    int64
	StringsOffset     int64
	MemoryOffset      int64
}

type serializedProgram struct {
	PackagesOffset       int64
	PackagesSize         int64
	CurrentPackageOffset int64

	InputsOffset int64
	InputsSize   int64

	OutputsOffset int64
	OutputsSize   int64

	CallStackOffset int64
	CallStackSize   int64

	CallCounter int64

	MemoryOffset int64
	MemorySize   int64

	StackPointer int64
	StackSize    int64

	DataSegmentSize     int64
	DataSegmentStartsAt int64

	HeapPointer  int64 //HeapPointer is probably related to HeapStartsAt
	HeapStartsAt int64
	HeapSize     int64

	Terminated int64

	VersionOffset int64
	VersionSize   int64
}

type serializedCall struct {
	OperatorOffset int64
	Line           int64
	FramePointer   int64
}

type serializedPackage struct {
	NameOffset            int64
	NameSize              int64
	ImportsOffset         int64
	ImportsSize           int64
	StructsOffset         int64
	StructsSize           int64
	GlobalsOffset         int64
	GlobalsSize           int64
	FunctionsOffset       int64
	FunctionsSize         int64
	CurrentFunctionOffset int64
	CurrentStructOffset   int64
}

type serializedStruct struct {
	NameOffset   int64
	NameSize     int64
	FieldsOffset int64
	FieldsSize   int64

	Size int64

	PackageOffset int64
}

type serializedFunction struct {
	NameOffset        int64
	NameSize          int64
	InputsOffset      int64
	InputsSize        int64
	OutputsOffset     int64
	OutputsSize       int64
	ExpressionsOffset int64
	ExpressionsSize   int64
	Size              int64
	Length            int64

	ListOfPointersOffset int64
	ListOfPointersSize   int64

	// We're going to determine this when procesing the expressions. Check serializedExpression type
	// IsBuiltin                        int64
	// OpCode                          int64

	CurrentExpressionOffset int64
	PackageOffset           int64
}

type serializedExpression struct {
	OperatorOffset int64
	// we add these two fields here so we don't add every native serializedFunction to the serialization
	// the CX runtime already knows about the natives properties. We just need the code if IsNative = true
	IsNative int64
	OpCode   int64

	InputsOffset  int64
	InputsSize    int64
	OutputsOffset int64
	OutputsSize   int64

	LabelOffset int64
	LabelSize   int64
	ThenLines   int64
	ElseLines   int64

	ExpressionType int64

	FunctionOffset int64
	PackageOffset  int64
}

type serializedArgument struct {
	NameOffset       int64
	NameSize         int64
	Type             int64
	StructTypeOffset int64 //WTF IS A CUSTOM TYPE!?
	Size             int64
	TotalSize        int64

	Offset int64

	IndirectionLevels           int64
	DereferenceLevels           int64
	DeclarationSpecifiersOffset int64
	DeclarationSpecifiersSize   int64

	IsSlice int64
	// IsArray      int64
	// IsArrayFirst int64
	IsPointer   int64
	IsReference int64

	IsStruct           int64
	IsRest             int64
	IsLocalDeclaration int64
	IsShortDeclaration int64
	PreviouslyDeclared int64

	PassBy     int64
	DoesEscape int64

	LengthsOffset int64
	LengthsSize   int64
	IndexesOffset int64
	IndexesSize   int64
	FieldsOffset  int64
	FieldsSize    int64
	InputsOffset  int64
	InputsSize    int64
	OutputsOffset int64
	OutputsSize   int64

	PackageOffset int64
}

type SerializedCXProgram struct {
	Index   serializedCXProgramIndex
	Program serializedProgram

	Packages     []serializedPackage
	PackagesMap  map[string]int64
	Structs      []serializedStruct
	StructsMap   map[string]int64
	Functions    []serializedFunction
	FunctionsMap map[string]int64

	Expressions []serializedExpression
	Arguments   []serializedArgument
	Calls       []serializedCall

	Strings    []byte
	StringsMap map[string]int64
	Integers   []int64

	Memory []byte
}

type SerializedDataSize struct {
	Program     int `json:"program"`
	Calls       int `json:"calls"`
	Packages    int `json:"packages"`
	Structs     int `json:"structs"`
	Functions   int `json:"functions"`
	Expressions int `json:"expressions"`
	Arguments   int `json:"arguments"`
	Integers    int `json:"integers"`
	Strings     int `json:"strings"`
	Memory      int `json:"memory"`
}
