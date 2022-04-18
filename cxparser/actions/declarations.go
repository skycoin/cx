package actions

import (
	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/constants"
	"github.com/skycoin/cx/cx/types"
)

// DeclarationSpecifiers is called to build a type of a variable or parameter.
//
// It is called repeatedly while the type is parsed.
//
// Returns the new type build from `declSpec` and `opTyp`.
//
// Input arguments description:
//   declSpec:     The incoming type
//   arrayLengths: The lengths of the array if `opTyp` = cxcore.DECL_ARRAY
//   opTyp:        The type of modification to `declSpec` (array of, pointer to, ...)
func DeclarationSpecifiers(declSpec *ast.CXArgument, arrayLengths []types.Pointer, opTyp int) *ast.CXArgument {
	switch opTyp {
	case constants.DECL_POINTER:
		declSpec.DeclarationSpecifiers = append(declSpec.DeclarationSpecifiers, constants.DECL_POINTER)

		declSpec.Size = types.POINTER_SIZE
		declSpec.TotalSize = types.POINTER_SIZE
		// declSpec.IndirectionLevels++

		if declSpec.Type == types.STR || declSpec.StructType != nil {
			declSpec.PointerTargetType = declSpec.Type
			declSpec.Type = types.POINTER
		}
		return declSpec
	case constants.DECL_ARRAY:
		for range arrayLengths {
			declSpec.DeclarationSpecifiers = append(declSpec.DeclarationSpecifiers, constants.DECL_ARRAY)
		}
		arg := declSpec
		// arg.IsArray = true
		arg.Lengths = arrayLengths
		arg.TotalSize = arg.Size * TotalLength(arg.Lengths)

		return arg
	case constants.DECL_SLICE:
		// for range arrayLengths {
		// 	declSpec.DeclarationSpecifiers = append(declSpec.DeclarationSpecifiers, cxcore.DECL_SLICE)
		// }

		arg := declSpec

		arg.DeclarationSpecifiers = append(arg.DeclarationSpecifiers, constants.DECL_SLICE)

		arg.IsSlice = true
		// arg.IsReference = true
		// arg.IsArray = true
		arg.PassBy = constants.PASSBY_REFERENCE

		arg.Lengths = append([]types.Pointer{0}, arg.Lengths...)
		// arg.Lengths = arrayLengths
		// arg.TotalSize = arg.Size
		// arg.Size = cxcore.TYPE_POINTER_SIZE
		arg.TotalSize = types.POINTER_SIZE

		return arg
	case constants.DECL_BASIC:
		arg := declSpec
		// arg.DeclarationSpecifiers = append(arg.DeclarationSpecifiers, cxcore.DECL_BASIC)
		arg.TotalSize = arg.Size
		return arg
	case constants.DECL_FUNC:
		// Creating this case if additional operations are needed in the
		// future.
		return declSpec
	}

	return nil
}

// DeclarationSpecifiersBasic() returns a type specifier created from one of the builtin types.
//
func DeclarationSpecifiersBasic(typeCode types.Code) *ast.CXArgument {
	arg := ast.MakeArgument("", CurrentFile, LineNo)
	arg.SetType(typeCode)
	if typeCode == types.AFF {
		// equivalent to slice of strings
		return DeclarationSpecifiers(arg, []types.Pointer{0}, constants.DECL_SLICE)
	}

	return DeclarationSpecifiers(arg, []types.Pointer{0}, constants.DECL_BASIC)
}

// DeclarationSpecifiersStruct() declares a struct
func DeclarationSpecifiersStruct(prgrm *ast.CXProgram, ident string, pkgName string,
	isExternal bool, currentFile string, lineNo int) *ast.CXArgument {
	pkg, err := prgrm.GetCurrentPackage()
	if err != nil {
		panic(err)
	}

	if isExternal {
		// custom type in an imported package
		imp, err := pkg.GetImport(prgrm, pkgName)
		if err != nil {
			panic(err)
		}

		strct, err := prgrm.GetStruct(ident, imp.Name)
		if err != nil {
			println(ast.CompilationError(currentFile, lineNo), err.Error())
			return nil
		}

		arg := ast.MakeArgument("", currentFile, lineNo)
		arg.Type = types.STRUCT
		arg.StructType = strct
		arg.Size = strct.GetStructSize(prgrm)
		arg.TotalSize = strct.GetStructSize(prgrm)
		arg.Package = ast.CXPackageIndex(pkg.Index)
		arg.DeclarationSpecifiers = append(arg.DeclarationSpecifiers, constants.DECL_STRUCT)

		return arg
	} else {
		// custom type in the current package
		strct, err := prgrm.GetStruct(ident, pkg.Name)
		if err != nil {
			println(ast.CompilationError(currentFile, lineNo), err.Error())
			return nil
		}

		arg := ast.MakeArgument("", currentFile, lineNo)
		arg.Type = types.STRUCT
		arg.DeclarationSpecifiers = append(arg.DeclarationSpecifiers, constants.DECL_STRUCT)
		arg.StructType = strct
		arg.Size = strct.GetStructSize(prgrm)
		arg.TotalSize = strct.GetStructSize(prgrm)
		arg.Package = ast.CXPackageIndex(pkg.Index)

		return arg
	}
}
