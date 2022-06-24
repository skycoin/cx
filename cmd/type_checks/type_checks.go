package type_checks

import (
	"github.com/skycoin/cx/cmd/declaration_extraction"
	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cxparser/actions"
	cxpartialparsing "github.com/skycoin/cx/cxparser/cxpartialparsing"
)

func ParseGlobals(globals []declaration_extraction.GlobalDeclaration) {

	for _, global := range globals {

		// Get Package
		pkg, err := cxpartialparsing.Program.GetPackage(global.PackageID)

		// If package not in AST
		if err != nil {

			newPkg := ast.MakePackage(global.PackageID)
			pkgIdx := cxpartialparsing.Program.AddPackage(newPkg)
			newPkg, err = cxpartialparsing.Program.GetPackageFromArray(pkgIdx)

			if err != nil {
				// error handling
			}

			pkg = newPkg
		}

		// Make and add global to AST
		globalArg := ast.MakeArgument(global.GlobalVariableName, global.FileID, global.LineNumber)
		globalArg.Offset = global.StartOffset
		globalArg.Package = ast.CXPackageIndex(pkg.Index)

		globalArgIdx := actions.AST.AddCXArgInArray(globalArg)

		pkg.AddGlobal(actions.AST, globalArgIdx)

	}
}

// func ParseEnums(enums []declaration_extraction.EnumDeclaration) {

// 	// Not sure about this one

// 	for _, enum := range enums {

// 		actions.de
// 	}
// }

func ParseStructs(structs []declaration_extraction.StructDeclaration) {

	// 1. iterate over all the structs
	// 2. add the struct name from the declaration
	// 3. search for fields
	// 4. fields to ast

	for _, strct := range structs {

		pkg, err := cxpartialparsing.Program.GetPackage(strct.PackageID)

		if err != nil {

			newPkg := ast.MakePackage(pkg.Name)
			pkgIdx := cxpartialparsing.Program.AddPackage(newPkg)
			newPkg, err = cxpartialparsing.Program.GetPackageFromArray(pkgIdx)

			if err != nil {
				// error handling
			}

			pkg = newPkg

		}

		structCX := ast.MakeStruct(strct.StructVariableName)
		structCX.Package = ast.CXPackageIndex(pkg.Index)

		cxpartialparsing.Program.AddStructInArray(structCX)

	}
}

func ParseFuncs(funcs []declaration_extraction.FuncDeclaration) {

	// 1. iterate over all the funcs
	// 2. extract inputs and outputs
	// 3. get the id and expression
	// 4. call function declaration

	for _, fun := range funcs {

		pkg, err := cxpartialparsing.Program.GetPackage(fun.PackageID)

		if err != nil {

			newPkg := ast.MakePackage(fun.PackageID)
			pkgIdx := cxpartialparsing.Program.AddPackage(newPkg)
			newPkg, err = cxpartialparsing.Program.GetPackageFromArray(pkgIdx)

			if err != nil {
				// error handling
			}

			pkg = newPkg

		}

		// srcBytes, err := os.ReadFile(fun.FileID)

		if err != nil {
			// error handling
		}

		// funcDeclaration := srcBytes[fun.StartOffset : fun.StartOffset+fun.Length]

		funcCX := ast.MakeFunction(fun.FuncVariableName, fun.FileID, fun.LineNumber)
		funcCX.Package = ast.CXPackageIndex(pkg.Index)

		cxpartialparsing.Program.AddFunctionInArray(funcCX)

	}

}
