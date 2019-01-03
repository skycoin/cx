package base

import (
	. "github.com/satori/go.uuid"
)

type CXCall struct {
	Operator     *CXFunction
	Line         int
	FramePointer int
}

/*
  Packages
*/

type CXConstant struct {
	// native constants. only used for pre-packaged constants (e.g. math package's PI)
	// these fields are used to feed WritePrimary
	Value []byte
	Type  int
}

type CXArgument struct {
	Lengths               []int // declared lengths at compile time
	DereferenceOperations []int // offset by array index, struct field, pointer
	DeclarationSpecifiers []int // used to determine finalSize
	Indexes               []*CXArgument
	Fields                []*CXArgument // strct.fld1.fld2().fld3
	Name                  string
	FileName              string
	ElementID             UUID
	Type                  int
	Size                  int // size of underlaying basic type
	TotalSize             int // total size of an array, performance reasons
	Offset                int
	IndirectionLevels     int
	DereferenceLevels     int
	PassBy                int // pass by value or reference
	FileLine              int
	CustomType            *CXStruct
	Package               *CXPackage
	IsSlice               bool
	IsArray               bool
	IsArrayFirst          bool // and then dereference
	IsPointer             bool
	IsReference           bool
	IsDereferenceFirst    bool // and then array
	IsStruct              bool
	IsRest                bool // pkg.var <- var is rest
	IsLocalDeclaration    bool
	IsShortDeclaration    bool
	PreviouslyDeclared    bool
	DoesEscape            bool
}
