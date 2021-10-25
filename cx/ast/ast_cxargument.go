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
	// Name defines the name of the `CXArgument`. Most of the
	// time, this field will be non-nil as this defines the name
	// of a variable or parameter in source code, but some
	// exceptions exist, such as in the case of literals
	// (e.g. `4`, `"Hello world!"`, `[3]i32{1, 2, 3}`.)
	Name string

	FileName string
	FileLine int
	Package  *CXPackage
}

// CXArgumentStruct ...
type CXArgumentStruct struct {
}

// CXArgumentPointer ...
type CXArgumentPointer struct {
}

// CXArgument is used to define local variables, global variables,
// literals (strings, numbers), inputs and outputs to function
// calls. All of the fields in this structure are determined at
// compile time.
type CXArgument struct {
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
	Indexes []*CXArgument
	// Fields stores what fields are being accessed from the
	// `CXArgument` and in what order. Whenever a `DEREF_FIELD` in
	// `DereferenceOperations` is found, we consume a field from
	// `Field` to determine the new offset to the desired
	// value.
	Fields []*CXArgument
	// Inputs defines the input parameters of a first-class
	// function. The `CXArgument` is of type `TYPE_FUNC` if
	// `ProgramInput` is non-nil.
	Inputs []*CXArgument
	// Outputs defines the output parameters of a first-class
	// function. The `CXArgument` is of type `TYPE_FUNC` if
	// `ProgramOutput` is non-nil.
	Outputs []*CXArgument
	// Type defines what's the basic or primitev type of the
	// `CXArgument`. `Type` can be equal to any of the `TYPE_*`
	// constants (e.g. `TYPE_STR`, `TYPE_I32`).
	Type types.Code
	// Size determines the size of the basic type. For example, if
	// the `CXArgument` is of type `TYPE_CUSTOM` (i.e. a
	// user-defined type or struct) and the size of the struct
	// representing the custom type is 10 bytes, then `Size == 10`.
	Size types.Pointer
	// TotalSize represents how many bytes are referenced by the
	// `CXArgument` in total. For example, if the `CXArgument`
	// defines an array of 5 struct instances of size 10 bytes,
	// then `TotalSize == 50`.
	TotalSize types.Pointer
	// Offset defines a relative memory offset (used in
	// conjunction with the frame pointer), in the case of local
	// variables, or it could define an absolute memory offset, in
	// the case of global variables and literals. It is used by
	// the CX virtual machine to find the bytes that represent the
	// value of the `CXArgument`.
	Offset types.Pointer
	// IndirectionLevels
	IndirectionLevels int
	DereferenceLevels int
	PassBy            int // pass by value or reference

	ArgDetails *CXArgumentDebug

	CustomType *CXStruct
	IsSlice    bool
	// IsArray                      bool
	IsPointer                    bool
	IsReference                  bool
	IsStruct                     bool
	IsInnerArg                   bool // ex. pkg.var <- var is the inner arg
	IsLocalDeclaration           bool
	IsShortAssignmentDeclaration bool // variables defined with :=
	IsInnerReference             bool // for example: &slice[0] or &struct.field
	PreviouslyDeclared           bool
	DoesEscape                   bool
}

/*
grep -rn "IsShortAssignmentDeclaration" .
IsShortAssignmentDeclaration - is this CXArgument the result of a `CASSIGN` operation (`:=`)?
./cxparser/cxgo/cxparser.y:1158:							from.Outputs[0].IsShortAssignmentDeclaration = true
./cxparser/cxgo/cxparser.y:1169:							from.Outputs[0].IsShortAssignmentDeclaration = true
./cxparser/cxgo/cxparser.go:2366:							from.Outputs[0].IsShortAssignmentDeclaration = true
./cxparser/cxgo/cxparser.go:2377:							from.Outputs[0].IsShortAssignmentDeclaration = true
./cxparser/actions/functions.go:147:		if len(expr.Outputs) > 0 && len(expr.Inputs) > 0 && expr.Outputs[0].IsShortAssignmentDeclaration && !expr.IsStructLiteral && !isParseOp(expr) {
./cxparser/actions/assignment.go:161:		sym.IsShortAssignmentDeclaration = true
./cxparser/actions/assignment.go:167:			toExpr.Outputs[0].IsShortAssignmentDeclaration = true
Binary file ./bin/cx matches
./docs/CompilerDevelopment.md:81:* IsShortAssignmentDeclaration - is this CXArgument the result of a `CASSIGN` operation (`:=`)?
./cx/serialize.go:168:	IsShortAssignmentDeclaration int32
./cx/serialize.go:337:	s.Arguments[argOff].IsShortAssignmentDeclaration = serializeBoolean(arg.IsShortAssignmentDeclaration)
./cx/serialize.go:1051:	arg.IsShortAssignmentDeclaration = dsBool(sArg.IsShortAssignmentDeclaration)
./cx/ast.go:234:	IsShortAssignmentDeclaration    bool
./cx/ast.go:1499:	IsShortAssignmentDeclaration    bool
*/

/*
	FileName              string
- filename and line number
- can be moved to CX AST annotations (comments to be skipped or map)

	FileLine
*/

/*
Note: Dereference Levels, is possible unused

grep -rn "DereferenceLevels" .

./cxparser/actions/functions.go:328:			if fld.IsPointer && fld.DereferenceLevels == 0 {
./cxparser/actions/functions.go:329:				fld.DereferenceLevels++
./cxparser/actions/functions.go:333:		if arg.IsStruct && arg.IsPointer && len(arg.Fields) > 0 && arg.DereferenceLevels == 0 {
./cxparser/actions/functions.go:334:			arg.DereferenceLevels++
./cxparser/actions/functions.go:1132:					nameFld.DereferenceLevels = sym.DereferenceLevels
./cxparser/actions/functions.go:1150:						nameFld.DereferenceLevels++
./cxparser/actions/expressions.go:328:		exprOut.DereferenceLevels++
./CompilerDevelopment.md:70:* DereferenceLevels - How many dereference operations are performed to get this CXArgument?
./cx/serialize.go:149:	DereferenceLevels           int32
./cx/serialize.go:300:	s.Arguments[argOff].DereferenceLevels = int32(arg.DereferenceLevels)
./cx/serialize.go:1008:	arg.DereferenceLevels = int(sArg.DereferenceLevels)
./cx/cxargument.go:22:	DereferenceLevels     int
./cx/utilities.go:143:	if arg.DereferenceLevels > 0 {
./cx/utilities.go:144:		for c := 0; c < arg.DereferenceLevels; c++ {
*/

/*
Note: IndirectionLevels does not appear to be used at all

 grep -rn "IndirectionLevels" .
./cxparser/actions/functions.go:951:	sym.IndirectionLevels = arg.IndirectionLevels
./cxparser/actions/declarations.go:379:			declSpec.IndirectionLevels++
./cxparser/actions/declarations.go:383:			for c := declSpec.IndirectionLevels - 1; c > 0; c-- {
./cxparser/actions/declarations.go:384:				pointer.IndirectionLevels = c
./cxparser/actions/declarations.go:388:			declSpec.IndirectionLevels++
./CompilerDevelopment.md:69:* IndirectionLevels - how many discrete levels of indirection to this specific CXArgument?
Binary file ./bin/cx matches
./cx/serialize.go:148:	IndirectionLevels           int32
./cx/serialize.go:299:	s.Arguments[argOff].IndirectionLevels = int32(arg.IndirectionLevels)
./cx/serialize.go:1007:	arg.IndirectionLevels = int(sArg.IndirectionLevels)
./cx/cxargument.go:21:	IndirectionLevels     int
*/

/*
IsDereferenceFirst - is this both an array and a pointer, and if so,
is the pointer first? Mutually exclusive with IsArrayFirst.

grep -rn "IsDereferenceFirst" .
./cxparser/actions/postfix.go:60:	if !elt.IsDereferenceFirst {
./cxparser/actions/expressions.go:331:			exprOut.IsDereferenceFirst = true
./CompilerDevelopment.md:76:* IsArrayFirst - is this both a pointer and an array, and if so, is the array first? Mutually exclusive with IsDereferenceFirst
./CompilerDevelopment.md:78:* IsDereferenceFirst - is this both an array and a pointer, and if so, is the pointer first? Mutually exclusive with IsArrayFirst.
Binary file ./bin/cx matches
./cx/serialize.go:161:	IsDereferenceFirst int32
./cx/serialize.go:314:	s.Arguments[argOff].IsDereferenceFirst = serializeBoolean(arg.IsDereferenceFirst)
./cx/serialize.go:1019:	arg.IsDereferenceFirst = dsBool(sArg.IsDereferenceFirst)
./cx/cxargument.go:32:	IsDereferenceFirst    bool // and then array
./cx/cxargument.go:43:IsDereferenceFirst - is this both an array and a pointer, and if so,

*/

/*
All "Is" can be removed
- because there is a constants for type (int) for defining the types
- could look in definition, specifier
- but use int lookup
	IsSlice               bool
	IsArray               bool
	IsArrayFirst          bool // and then dereference
	IsPointer             bool
	IsReference           bool
	IsDereferenceFirst    bool // and then array
	IsStruct              bool
	IsInnerArg                bool // pkg.var <- var is rest
	IsLocalDeclaration    bool
	IsShortAssignmentDeclaration    bool
	IsInnerReference      bool // for example: &slice[0] or &struct.field

*/

/*

Note: PAssBy is not used too many place
Note: Low priority for deprecation
- isnt this same as "pointer"

grep -rn "PassBy" .
./cxparser/actions/misc.go:425:			arg.PassBy = PASSBY_REFERENCE
./cxparser/actions/functions.go:666:					out.PassBy = PASSBY_VALUE
./cxparser/actions/functions.go:678:		if elt.PassBy == PASSBY_REFERENCE &&
./cxparser/actions/functions.go:712:			out.PassBy = PASSBY_VALUE
./cxparser/actions/functions.go:723:				assignElt.PassBy = PASSBY_VALUE
./cxparser/actions/functions.go:915:			expr.Inputs[0].PassBy = PASSBY_REFERENCE
./cxparser/actions/functions.go:1153:					nameFld.PassBy = fld.PassBy
./cxparser/actions/functions.go:1157:						nameFld.PassBy = PASSBY_REFERENCE
./cxparser/actions/literals.go:219:				sym.PassBy = PASSBY_REFERENCE
./cxparser/actions/expressions.go:336:		baseOut.PassBy = PASSBY_REFERENCE
./cxparser/actions/assignment.go:57:		out.PassBy = PASSBY_REFERENCE
./cxparser/actions/assignment.go:208:		to[0].Outputs[0].PassBy = from[idx].Outputs[0].PassBy
./cxparser/actions/assignment.go:234:			to[0].Outputs[0].PassBy = from[idx].Operator.Outputs[0].PassBy
./cxparser/actions/assignment.go:244:			to[0].Outputs[0].PassBy = from[idx].Operator.Outputs[0].PassBy
./cxparser/actions/declarations.go:55:			glbl.PassBy = offExpr[0].Outputs[0].PassBy
./cxparser/actions/declarations.go:69:				declaration_specifiers.PassBy = glbl.PassBy
./cxparser/actions/declarations.go:85:				declaration_specifiers.PassBy = glbl.PassBy
./cxparser/actions/declarations.go:103:			declaration_specifiers.PassBy = glbl.PassBy
./cxparser/actions/declarations.go:324:			declarationSpecifiers.PassBy = initOut.PassBy
./cxparser/actions/declarations.go:417:		arg.PassBy = PASSBY_REFERENCE
./CompilerDevelopment.md:71:* PassBy - an int constant representing how the variable is passed - pass by value, or pass by reference.

./cx/op_http.go:50:	headerFld.PassBy = PASSBY_REFERENCE
./cx/op_http.go:75:	transferEncodingFld.PassBy = PASSBY_REFERENCE
./cx/serialize.go:168:	PassBy     int32
./cx/serialize.go:321:	s.Arguments[argOff].PassBy = int32(arg.PassBy)
./cx/serialize.go:1009:	arg.PassBy = int(sArg.PassBy)
./cx/execute.go:366:				if inp.PassBy == PASSBY_REFERENCE {
./cx/cxargument.go:23:	PassBy                int // pass by value or reference
./cx/op_misc.go:36:		switch elt.PassBy {
./cx/utilities.go:184:	if arg.PassBy == PASSBY_REFERENCE {
*/

// ----------------------------------------------------------------
//                             `CXArgument` Getters

// GetAssignmentElement ...
func (arg *CXArgument) GetAssignmentElement() *CXArgument {
	if len(arg.Fields) > 0 {
		return arg.Fields[len(arg.Fields)-1].GetAssignmentElement()
	}
	return arg

}

// GetType ...
func (arg *CXArgument) GetType() types.Code {
	fieldCount := len(arg.Fields)
	if fieldCount > 0 {
		return arg.Fields[fieldCount-1].GetType()
	}
	return arg.Type
}

// ----------------------------------------------------------------
//                     `CXArgument` Member handling

// AddPackage assigns CX package `pkg` to CX argument `arg`.
func (arg *CXArgument) AddPackage(pkg *CXPackage) *CXArgument {
	// pkg, err := PROGRAM.GetPackage(pkgName)
	// if err != nil {
	// 	panic(err)
	// }
	arg.ArgDetails.Package = pkg
	return arg
}

// AddType ...
func (arg *CXArgument) AddType(typeCode types.Code) *CXArgument {
	arg.Type = typeCode
	size := typeCode.Size()
	arg.Size = size
	arg.TotalSize = size
	arg.DeclarationSpecifiers = append(arg.DeclarationSpecifiers, constants.DECL_BASIC)
	return arg
}

// AddInput adds input parameters to `arg` in case arg is of type `TYPE_FUNC`.
func (arg *CXArgument) AddInput(inp *CXArgument) *CXArgument {
	arg.Inputs = append(arg.Inputs, inp)
	if inp.ArgDetails.Package == nil {
		inp.ArgDetails.Package = arg.ArgDetails.Package
	}
	return arg
}

// AddOutput adds output parameters to `arg` in case arg is of type `TYPE_FUNC`.
func (arg *CXArgument) AddOutput(out *CXArgument) *CXArgument {
	arg.Outputs = append(arg.Outputs, out)
	if out.ArgDetails.Package == nil {
		out.ArgDetails.Package = arg.ArgDetails.Package
	}
	return arg
}

// Pointer takes an already defined `CXArgument` and turns it into a pointer.
//Only used once, deprecate
//TODO: only used by HTTP, create a better module system
func Pointer(arg *CXArgument) *CXArgument {
	arg.DeclarationSpecifiers = append(arg.DeclarationSpecifiers, constants.DECL_POINTER)
	arg.IsPointer = true
	arg.Size = types.POINTER_SIZE
	arg.TotalSize = types.POINTER_SIZE

	return arg
}

// Struct helper for creating a struct parameter. It creates a
// `CXArgument` named `argName`, that represents a structure instane of
// `strctName`, from package `pkgName`.
func Struct(pkgName, strctName, argName string) *CXArgument {
	pkg, err := PROGRAM.GetPackage(pkgName)
	if err != nil {
		panic(err)
	}

	strct, err := pkg.GetStruct(strctName)
	if err != nil {
		panic(err)
	}

	arg := MakeArgument(argName, "", -1).AddType(types.CUSTOM)
	arg.DeclarationSpecifiers = append(arg.DeclarationSpecifiers, constants.DECL_STRUCT)
	arg.Size = strct.Size
	arg.TotalSize = strct.Size
	arg.CustomType = strct

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
func Func(pkg *CXPackage, inputs []*CXArgument, outputs []*CXArgument) *CXArgument {
	arg := Param(types.FUNC)
	arg.ArgDetails.Package = pkg
	arg.Inputs = inputs
	arg.Outputs = outputs
	return arg
}

// Param ...
func Param(typeCode types.Code) *CXArgument {
	arg := MakeArgument("", "", -1).AddType(typeCode)
	arg.IsLocalDeclaration = true
	return arg
}

// MakeArgument ...
func MakeArgument(name string, fileName string, fileLine int) *CXArgument {
	return &CXArgument{
		ArgDetails: &CXArgumentDebug{
			Name:     name,
			FileName: fileName,
			FileLine: fileLine,
		},
	}

}

// MakeField ...
func MakeField(name string, typeCode types.Code, fileName string, fileLine int) *CXArgument {
	return &CXArgument{
		ArgDetails: &CXArgumentDebug{
			Name:     name,
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
		ArgDetails: &CXArgumentDebug{
			Name:     name,
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
