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
	Type                  int
	Size                  int // size of underlaying basic type
	TotalSize             int // total size of an array, performance reasons
	Offset                int
	IndirectionLevels     int
	DereferenceLevels     int
	PassBy                int // pass by value or reference

	FileName              string
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

//Why isnt there a ast.go

//Components of AST
//CXProgram
//CXPackage
//CXArgument
//CXExpression
//CXStruct
//CXFunction

/*
	FileName              string
- filename and line number
- can be moved to CX AST annotations (comments to be skipped or map)
	
	FileLine
*/

/*
Note: Dereference Levels, is possible unused

grep -rn "DereferenceLevels" .

./cxgo/actions/functions.go:328:			if fld.IsPointer && fld.DereferenceLevels == 0 {
./cxgo/actions/functions.go:329:				fld.DereferenceLevels++
./cxgo/actions/functions.go:333:		if arg.IsStruct && arg.IsPointer && len(arg.Fields) > 0 && arg.DereferenceLevels == 0 {
./cxgo/actions/functions.go:334:			arg.DereferenceLevels++
./cxgo/actions/functions.go:1132:					nameFld.DereferenceLevels = sym.DereferenceLevels
./cxgo/actions/functions.go:1150:						nameFld.DereferenceLevels++
./cxgo/actions/expressions.go:328:		exprOut.DereferenceLevels++
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
./cxgo/actions/functions.go:951:	sym.IndirectionLevels = arg.IndirectionLevels
./cxgo/actions/declarations.go:379:			declSpec.IndirectionLevels++
./cxgo/actions/declarations.go:383:			for c := declSpec.IndirectionLevels - 1; c > 0; c-- {
./cxgo/actions/declarations.go:384:				pointer.IndirectionLevels = c
./cxgo/actions/declarations.go:388:			declSpec.IndirectionLevels++
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
./cxgo/actions/postfix.go:60:	if !elt.IsDereferenceFirst {
./cxgo/actions/expressions.go:331:			exprOut.IsDereferenceFirst = true
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
	IsRest                bool // pkg.var <- var is rest
	IsLocalDeclaration    bool
	IsShortDeclaration    bool
	IsInnerReference      bool // for example: &slice[0] or &struct.field

*/

/*

Note: PAssBy is not used too many place
Note: Low priority for deprecation
- isnt this same as "pointer"

grep -rn "PassBy" .
./cxgo/actions/misc.go:425:			arg.PassBy = PASSBY_REFERENCE
./cxgo/actions/functions.go:666:					out.PassBy = PASSBY_VALUE
./cxgo/actions/functions.go:678:		if elt.PassBy == PASSBY_REFERENCE &&
./cxgo/actions/functions.go:712:			out.PassBy = PASSBY_VALUE
./cxgo/actions/functions.go:723:				assignElt.PassBy = PASSBY_VALUE
./cxgo/actions/functions.go:915:			expr.Inputs[0].PassBy = PASSBY_REFERENCE
./cxgo/actions/functions.go:1153:					nameFld.PassBy = fld.PassBy
./cxgo/actions/functions.go:1157:						nameFld.PassBy = PASSBY_REFERENCE
./cxgo/actions/literals.go:219:				sym.PassBy = PASSBY_REFERENCE
./cxgo/actions/expressions.go:336:		baseOut.PassBy = PASSBY_REFERENCE
./cxgo/actions/assignment.go:57:		out.PassBy = PASSBY_REFERENCE
./cxgo/actions/assignment.go:208:		to[0].Outputs[0].PassBy = from[idx].Outputs[0].PassBy
./cxgo/actions/assignment.go:234:			to[0].Outputs[0].PassBy = from[idx].Operator.Outputs[0].PassBy
./cxgo/actions/assignment.go:244:			to[0].Outputs[0].PassBy = from[idx].Operator.Outputs[0].PassBy
./cxgo/actions/declarations.go:55:			glbl.PassBy = offExpr[0].Outputs[0].PassBy
./cxgo/actions/declarations.go:69:				declaration_specifiers.PassBy = glbl.PassBy
./cxgo/actions/declarations.go:85:				declaration_specifiers.PassBy = glbl.PassBy
./cxgo/actions/declarations.go:103:			declaration_specifiers.PassBy = glbl.PassBy
./cxgo/actions/declarations.go:324:			declarationSpecifiers.PassBy = initOut.PassBy
./cxgo/actions/declarations.go:417:		arg.PassBy = PASSBY_REFERENCE
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
