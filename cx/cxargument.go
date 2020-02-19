package cxcore

// The CXArgument struct contains a variable, i.e. a combination of a name and a type.
//
// It is used when declaring variables and in function parameters.
//
type CXArgument struct {
	Lengths               []int // declared lengths at compile time
	DereferenceOperations []int // offset by array index, struct field, pointer
	DeclarationSpecifiers []int // used to determine finalSize
	Indexes               []*CXArgument
	Fields                []*CXArgument // strct.fld1.fld2().fld3
	Inputs                []*CXArgument // Input parameters in case `CXArgument` is of type TYPE_FUNC
	Outputs               []*CXArgument // Output parameters in case `CXArgument` is of type TYPE_FUNC
	Name                  string
	FileName              string
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
	IsInnerReference      bool // for example: &slice[0] or &struct.field
	PreviouslyDeclared    bool
	DoesEscape            bool
}

// MakeArgument ...
func MakeArgument(name string, fileName string, fileLine int) *CXArgument {
	return &CXArgument{
		Name:     name,
		FileName: fileName,
		FileLine: fileLine}
}

// MakeField ...
func MakeField(name string, typ int, fileName string, fileLine int) *CXArgument {
	return &CXArgument{
		Name:     name,
		Type:     typ,
		FileName: fileName,
		FileLine: fileLine,
	}
}

// MakeGlobal ...
func MakeGlobal(name string, typ int, fileName string, fileLine int) *CXArgument {
	size := GetArgSize(typ)
	global := &CXArgument{
		Name:     name,
		Type:     typ,
		Size:     size,
		Offset:   HeapOffset,
		FileName: fileName,
		FileLine: fileLine,
	}
	HeapOffset += size
	return global
}

// ----------------------------------------------------------------
//                             Getters

// ----------------------------------------------------------------
//                     Member handling

// AddPackage assigns CX package `pkg` to CX argument `arg`.
func (arg *CXArgument) AddPackage(pkg *CXPackage) *CXArgument {
	// pkg, err := PROGRAM.GetPackage(pkgName)
	// if err != nil {
	// 	panic(err)
	// }
	arg.Package = pkg
	return arg
}

// AddType ...
func (arg *CXArgument) AddType(typ string) *CXArgument {
	if typCode, found := TypeCodes[typ]; found {
		arg.Type = typCode
		size := GetArgSize(typCode)
		arg.Size = size
		arg.TotalSize = size
		arg.DeclarationSpecifiers = append(arg.DeclarationSpecifiers, DECL_BASIC)
	} else {
		arg.Type = TYPE_UNDEFINED
	}

	return arg
}

// AddInput adds input parameters to `arg` in case arg is of type `TYPE_FUNC`.
func (arg *CXArgument) AddInput(inp *CXArgument) *CXArgument {
	arg.Inputs = append(arg.Inputs, inp)
	if inp.Package == nil {
		inp.Package = arg.Package
	}
	return arg
}

// AddOutput adds output parameters to `arg` in case arg is of type `TYPE_FUNC`.
func (arg *CXArgument) AddOutput(out *CXArgument) *CXArgument {
	arg.Outputs = append(arg.Outputs, out)
	if out.Package == nil {
		out.Package = arg.Package
	}
	return arg
}
