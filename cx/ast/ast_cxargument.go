package ast

import (
	"github.com/skycoin/cx/cx/constants"
	"github.com/skycoin/cx/cx/globals"
	"github.com/skycoin/cx/cx/types"
)

type CXArgumentSlice struct {
	// Lengths is used if the `CXArgument` defines an array or a
	// slice. The number of dimensions for the array/slice is
	// equal to `len(Lengths)`, while the contents of `Lengths`
	// define the sizes of each dimension. In the case of a slice,
	// `Lengths` only determines the number of dimensions and the
	// sizes are all equal to 0 (these 0s are not used for any
	// computation).
	Lengths []int
	IsSlice bool
}

// CXArgumentDebug ...
type CXArgumentDebug struct {
	FileName string
	FileLine int
}

// CXArgumentStruct ...
type CXArgumentStruct struct {
}

// CXArgumentPointer ...
type CXArgumentPointer struct {
}

type CXArgumentIndex int

// CXArgument is used to define local variables, global variables,
// literals (strings, numbers), inputs and outputs to function
// calls. All of the fields in this structure are determined at
// compile time.
type CXArgument struct {
	// Name defines the name of the `CXArgument`. Most of the
	// time, this field will be non-nil as this defines the name
	// of a variable or parameter in source code, but some
	// exceptions exist, such as in the case of literals
	// (e.g. `4`, `"Hello world!"`, `[3]i32{1, 2, 3}`.)
	Name  string
	Index int

	Package CXPackageIndex

	// Lengths is used if the `CXArgument` defines an array or a
	// slice. The number of dimensions for the array/slice is
	// equal to `len(Lengths)`, while the contents of `Lengths`
	// define the sizes of each dimension. In the case of a slice,
	// `Lengths` only determines the number of dimensions and the
	// sizes are all equal to 0 (these 0s are not used for any
	// computation).
	Lengths []types.Pointer

	// DereferenceOperations is a slice of integers where each
	// integer corresponds a `DEREF_*` constant (for example
	// `DEREF_ARRAY`, `DEREF_POINTER`.). A dereference is a
	// process where we consider the bytes at `Offset : Offset +
	// TotalSize` as an address in memory, and we use that address
	// to find the desired value (the referenced
	// value).
	DereferenceOperations []int

	// DeclarationSpecifiers is a slice of integers where each
	// integer corresponds a `DECL_*` constant (for example
	// `DECL_ARRAY`, `DECL_POINTER`.). Declarations are used to
	// create complex types such as `[5][]*Point` (an array of 5
	// slices of pointers to struct instances of type
	// `Point`).
	DeclarationSpecifiers []int

	// Indexes stores what indexes we want to access from the
	// `CXArgument`. A non-nil `Indexes` means that the
	// `CXArgument` is an index or a slice. The elements of
	// `Indexes` can be any `CXArgument` (for example, literals
	// and variables).
	Indexes []CXTypeSignatureIndex

	// Fields stores what fields are being accessed from the
	// `CXArgument` and in what order. Whenever a `DEREF_FIELD` in
	// `DereferenceOperations` is found, we consume a field from
	// `Field` to determine the new offset to the desired
	// value.
	Fields []CXArgumentIndex

	// Inputs defines the input parameters of a first-class
	// function. The `CXArgument` is of type `TYPE_FUNC` if
	// `ProgramInput` is non-nil.
	Inputs []CXArgumentIndex

	// Outputs defines the output parameters of a first-class
	// function. The `CXArgument` is of type `TYPE_FUNC` if
	// `ProgramOutput` is non-nil.
	Outputs []CXArgumentIndex

	// Type defines what's the basic or primitev type of the
	// `CXArgument`. `Type` can be equal to any of the `TYPE_*`
	// constants (e.g. `TYPE_STR`, `TYPE_I32`).
	Type types.Code

	PointerTargetType types.Code

	// Size determines the size of the basic type. For example, if
	// the `CXArgument` is of type `TYPE_STRUCT` (i.e. a
	// user-defined type or struct) and the size of the struct
	// representing the struct type is 10 bytes, then `Size == 10`.
	Size types.Pointer

	// Offset defines a relative memory offset (used in
	// conjunction with the frame pointer), in the case of local
	// variables, or it could define an absolute memory offset, in
	// the case of global variables and literals. It is used by
	// the CX virtual machine to find the bytes that represent the
	// value of the `CXArgument`.
	Offset types.Pointer

	PassBy int // pass by value or reference

	ArgDetails *CXArgumentDebug

	StructType         *CXStruct
	IsSlice            bool
	IsStruct           bool
	IsInnerReference   bool // for example: &slice[0] or &struct.field
	PreviouslyDeclared bool
}

func (arg CXArgument) IsPointer() bool {
	return arg.Type == types.POINTER
}

func (arg CXArgument) IsString() bool {
	return arg.PointerTargetType == types.STR || arg.Type == types.STR
}

/*
	FileName              string
- filename and line number
- can be moved to CX AST annotations (comments to be skipped or map)

	FileLine
*/

/*
All "Is" can be removed
- because there is a constants for type (int) for defining the types
- could look in definition, specifier
- but use int lookup
	IsSlice               bool
	IsStruct              bool
	IsInnerReference      bool // for example: &slice[0] or &struct.field

*/

/*

Note: PAssBy is not used too many place
Note: Low priority for deprecation
- isnt this same as "pointer"
*/

// ----------------------------------------------------------------
//                             `CXArgument` Getters

// GetAssignmentElement ...
func (arg *CXArgument) GetAssignmentElement(prgrm *CXProgram) *CXArgument {
	if len(arg.Fields) > 0 {
		return prgrm.GetCXArgFromArray(arg.Fields[len(arg.Fields)-1]).GetAssignmentElement(prgrm)
	}
	return arg

}

// GetType ...
func (arg *CXArgument) GetType(prgrm *CXProgram) types.Code {
	fieldCount := len(arg.Fields)
	if fieldCount > 0 {
		return prgrm.GetCXArgFromArray(arg.Fields[fieldCount-1]).GetType(prgrm)
	}

	if arg.Type == types.POINTER {
		return arg.PointerTargetType
	}
	return arg.Type
}

func (arg *CXArgument) GetSize() types.Pointer {
	if arg.Type != types.I8 {
		return arg.Type.Size()
	}

	return arg.Size
}

// ----------------------------------------------------------------
//                     `CXArgument` Member handling

// SetPackage sets CX package `pkg` of CX argument `arg`.
func (arg *CXArgument) SetPackage(pkg *CXPackage) *CXArgument {
	arg.Package = CXPackageIndex(pkg.Index)
	return arg
}

// SetType ...
func (arg *CXArgument) SetType(typeCode types.Code) *CXArgument {
	arg.Type = typeCode
	size := typeCode.Size()
	arg.Size = size
	arg.DeclarationSpecifiers = append(arg.DeclarationSpecifiers, constants.DECL_BASIC)
	return arg
}

// AddInput adds input parameters to `arg` in case arg is of type `TYPE_FUNC`.
func (arg *CXArgument) AddInput(prgrm *CXProgram, inp *CXArgument) *CXArgument {
	if inp.Package == -1 {
		inp.Package = arg.Package
	}

	inpIdx := prgrm.AddCXArgInArray(inp)
	arg.Inputs = append(arg.Inputs, inpIdx)

	return arg
}

// AddOutput adds output parameters to `arg` in case arg is of type `TYPE_FUNC`.
func (arg *CXArgument) AddOutput(prgrm *CXProgram, out *CXArgument) *CXArgument {
	if out.Package == -1 {
		out.Package = arg.Package
	}

	outIdx := prgrm.AddCXArgInArray(out)
	arg.Outputs = append(arg.Outputs, outIdx)

	return arg
}

// MakePointer takes an already defined `CXArgument` and turns it into a pointer.
//Only used once, deprecate
//TODO: only used by HTTP, create a better module system
func MakePointer(arg *CXArgument) *CXArgument {
	arg.DeclarationSpecifiers = append(arg.DeclarationSpecifiers, constants.DECL_POINTER)
	arg.PointerTargetType = arg.Type
	arg.Type = types.POINTER
	arg.Size = types.POINTER_SIZE

	return arg
}

// MakeStructParameter helper for creating a struct parameter. It creates a
// `CXArgument` named `argName`, that represents a structure instane of
// `strctName`, from package `pkgName`.
func MakeStructParameter(prgrm *CXProgram, pkgName, strctName, argName string) *CXArgument {
	pkg, err := prgrm.GetPackage(pkgName)
	if err != nil {
		panic(err)
	}

	strct, err := pkg.GetStruct(prgrm, strctName)
	if err != nil {
		panic(err)
	}

	arg := MakeArgument(argName, "", -1).SetType(types.STRUCT)
	// arg.DeclarationSpecifiers = append(arg.DeclarationSpecifiers, constants.DECL_STRUCT)
	arg.Size = strct.GetStructSize(prgrm)
	arg.StructType = strct

	return arg
}

// Slice Helper function for creating parameters for standard library operators.
// The current standard library only uses basic types and slices. If more options are needed, modify this function
func Slice(typeCode types.Code) *CXArgument {
	arg := Param(typeCode)
	arg.IsSlice = true
	arg.DeclarationSpecifiers = append(arg.DeclarationSpecifiers, constants.DECL_SLICE)
	return arg
}

// Func Helper function for creating function parameters for standard library operators.
// The current standard library only uses basic types and slices. If more options are needed, modify this function
// func Func(pkg *CXPackage, inputs []*CXArgument, outputs []*CXArgument) *CXArgument {
// 	arg := Param(types.FUNC)
// 	arg.Package = CXPackageIndex(pkg.Index)
// 	arg.Inputs = inputs
// 	arg.Outputs = outputs
// 	return arg
// }

// Param ...
func Param(typeCode types.Code) *CXArgument {
	arg := MakeArgument("", "", -1).SetType(typeCode)
	// arg.IsLocalDeclaration = true
	return arg
}

// MakeArgument ...
func MakeArgument(name string, fileName string, fileLine int) *CXArgument {
	return &CXArgument{
		Name:    name,
		Package: -1,
		ArgDetails: &CXArgumentDebug{
			FileName: fileName,
			FileLine: fileLine,
		},
	}

}

// MakeField ...
func MakeField(name string, typeCode types.Code, fileName string, fileLine int) *CXArgument {
	return &CXArgument{
		Name:    name,
		Package: -1,
		ArgDetails: &CXArgumentDebug{
			FileName: fileName,
			FileLine: fileLine,
		},

		Type: typeCode,
	}
}

// MakeGlobal ...
func MakeGlobal(name string, typeCode types.Code, fileName string, fileLine int) *CXArgument {
	size := typeCode.Size()
	global := &CXArgument{
		Name:    name,
		Package: -1,
		ArgDetails: &CXArgumentDebug{
			FileName: fileName,
			FileLine: fileLine,
		},
		Type:   typeCode,
		Size:   size,
		Offset: globals.HeapOffset,
	}
	globals.HeapOffset += size
	return global
}

// ------------------------------------------------------------------------------------------
//           Special functions to determine its type (atomic, array etomic, etc)

func IsTypeAtomic(arg *CXArgument) bool {
	// TODO: implement including types.IDENTIFIER
	// return (arg.Type.IsPrimitive() || arg.Type == types.IDENTIFIER) && !arg.IsSlice && len(arg.Lengths) == 0 && len(arg.Fields) == 0 && len(arg.DereferenceOperations) == 0 && (len(arg.DeclarationSpecifiers) == 0 || (len(arg.DeclarationSpecifiers) == 1 && arg.DeclarationSpecifiers[0] == constants.DECL_BASIC))

	return arg.Type.IsPrimitive() && !arg.IsSlice && len(arg.Lengths) == 0 && len(arg.Fields) == 0 && len(arg.DereferenceOperations) == 0 && (len(arg.DeclarationSpecifiers) == 0 || (len(arg.DeclarationSpecifiers) == 1 && arg.DeclarationSpecifiers[0] == constants.DECL_BASIC))
}

func IsTypePointerAtomic(arg *CXArgument) bool {
	return arg.Type.IsPrimitive() && arg.PointerTargetType == 0 && arg.StructType == nil && !arg.IsSlice && len(arg.Lengths) == 0 && len(arg.Fields) == 0 && len(arg.DereferenceOperations) == 0 && (len(arg.DeclarationSpecifiers) == 2 && arg.DeclarationSpecifiers[0] == constants.DECL_BASIC && arg.DeclarationSpecifiers[1] == constants.DECL_POINTER)
}

func IsTypeArrayAtomic(arg *CXArgument) bool {
	isThereDeclPointer := false
	for _, decl := range arg.DeclarationSpecifiers {
		if decl == constants.DECL_POINTER {
			isThereDeclPointer = true
		}
	}

	return !isThereDeclPointer && arg.Type.IsPrimitive() && !arg.IsSlice && len(arg.Lengths) > 0 && len(arg.Indexes) == 0 && len(arg.Fields) == 0 && len(arg.DereferenceOperations) == 0 && (len(arg.DeclarationSpecifiers) >= 2 && arg.DeclarationSpecifiers[0] == constants.DECL_BASIC && arg.DeclarationSpecifiers[1] == constants.DECL_ARRAY)
}

func IsTypePointerArrayAtomic(arg *CXArgument) bool {
	isThereDeclPointer := false
	for _, decl := range arg.DeclarationSpecifiers {
		if decl == constants.DECL_POINTER {
			isThereDeclPointer = true
		}
	}

	return isThereDeclPointer && arg.Type.IsPrimitive() && !arg.IsSlice && len(arg.Lengths) > 0 && len(arg.Indexes) == 0 && len(arg.Fields) == 0 && len(arg.DereferenceOperations) == 0 && (len(arg.DeclarationSpecifiers) >= 2 && arg.DeclarationSpecifiers[0] == constants.DECL_BASIC && arg.DeclarationSpecifiers[1] == constants.DECL_ARRAY)
}

func IsTypeSliceAtomic(arg *CXArgument) bool {
	isThereDeclPointer := false
	for _, decl := range arg.DeclarationSpecifiers {
		if decl == constants.DECL_POINTER {
			isThereDeclPointer = true
		}
	}

	return !isThereDeclPointer && arg.IsSlice && len(arg.Lengths) > 0 && (arg.Type.IsPrimitive() || arg.Type == types.STR)
}

func IsTypeStruct(arg *CXArgument) bool {
	return !arg.IsSlice && len(arg.Lengths) == 0 && arg.Type == types.STRUCT
}
