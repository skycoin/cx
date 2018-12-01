package base

import (
        . "github.com/satori/go.uuid"
)

/*
  Root Program
*/

type CXProgram struct {
        Packages                        []*CXPackage
	Memory                          []byte
	Inputs                          []*CXArgument
        Outputs                         []*CXArgument
	CallStack                       []CXCall
	Path                            string
        CurrentPackage                  *CXPackage
	CallCounter                     int
        HeapPointer                     int
        StackPointer                    int
        HeapStartsAt                    int
	ElementID                       UUID
        Terminated                      bool
}

type CXCall struct {
        Operator                        *CXFunction
        Line                            int
        FramePointer                    int
}

/*
  Packages
*/

type CXPackage struct {
	Imports                         []*CXPackage
        Functions                       []*CXFunction
        Structs                         []*CXStruct
        Globals                         []*CXArgument
	Name                            string
	CurrentFunction                 *CXFunction
        CurrentStruct                   *CXStruct
        ElementID                       UUID
}

/*
  Structs
*/

type CXStruct struct {
	Fields                          []*CXArgument
	Name                            string
	Size                            int
	Package                         *CXPackage
        ElementID                       UUID
}

/*
  Functions
*/

type CXFunction struct {
	ListOfPointers                  []*CXArgument
        Inputs                          []*CXArgument
        Outputs                         []*CXArgument
        Expressions                     []*CXExpression
	Name                            string
	Length                          int // number of expressions, pre-computed for performance
        Size                            int // automatic memory size
        OpCode                          int
	CurrentExpression               *CXExpression
        Package                         *CXPackage
        ElementID                       UUID
        IsNative                        bool
}

type CXExpression struct {
	Inputs                          []*CXArgument
        Outputs                         []*CXArgument
	Label                           string
	FileName                        string
	Operator                        *CXFunction
	// debugging
        FileLine                        int
	// used for jmp statements
        ThenLines                       int
        ElseLines                       int
	Function                        *CXFunction
        Package                         *CXPackage
        ElementID                       UUID
        IsMethodCall                    bool
        IsStructLiteral                 bool
        IsArrayLiteral                  bool
        IsUndType                       bool
	IsBreak                         bool
	IsContinue                      bool
}

type CXConstant struct {
        // native constants. only used for pre-packaged constants (e.g. math package's PI)
        // these fields are used to feed WritePrimary
	Value                           []byte
        Type                            int
}

type CXArgument struct {
	Lengths                         []int // declared lengths at compile time
        Indexes                         []*CXArgument
        Fields                          []*CXArgument // strct.fld1.fld2().fld3
	DereferenceOperations           []int // offset by array index, struct field, pointer
        DeclarationSpecifiers           []int // used to determine finalSize
	Name                            string
	FileName                        string
        ElementID                       UUID
        Type                            int
        Size                            int // size of underlaying basic type
        TotalSize                       int // total size of an array, performance reasons
        Offset                          int
        IndirectionLevels               int
        DereferenceLevels               int
	PassBy                          int  // pass by value or reference
	FileLine                        int
	CustomType                      *CXStruct
	Package                         *CXPackage
        IsSlice                         bool
        IsArray                         bool
        IsArrayFirst                    bool // and then dereference
        IsPointer                       bool
        IsReference                     bool
        IsDereferenceFirst              bool // and then array
        IsStruct                        bool
        IsRest                          bool // pkg.var <- var is rest
        IsLocalDeclaration              bool
	IsShortDeclaration              bool
        PreviouslyDeclared              bool
        DoesEscape                      bool
}
