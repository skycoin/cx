package type_checks

import (
	"log"

	"github.com/skycoin/cx/cmd/declaration_extraction"
	"github.com/skycoin/cx/cx/ast"
	cxinit "github.com/skycoin/cx/cx/init"
	"github.com/skycoin/cx/cxparser/actions"
)

// ParseGlobals make the CXArguments and add them to the AST
// 1. iterates over all the global declarations
// 2. gets/makes the package
// 3.

func ParseGlobals(globals []declaration_extraction.GlobalDeclaration) {

	if actions.AST == nil {
		actions.AST = cxinit.MakeProgram()
	}

	for _, global := range globals {

		newPkg := ast.MakePackage(global.PackageID)
		pkgIdx := actions.AST.AddPackage(newPkg)
		newPkg, err := actions.AST.GetPackageFromArray(pkgIdx)

		// Get Package
		pkg, err := actions.AST.GetPackage(global.PackageID)

		// If package not in AST
		if err != nil || pkg == nil {

			newPkg := ast.MakePackage(global.PackageID)
			pkgIdx := actions.AST.AddPackage(newPkg)
			newPkg, err := actions.AST.GetPackageFromArray(pkgIdx)

			if err != nil {
				log.Fatal(err)
			}

			pkg = newPkg
		}

		// read bytes
		// srcBytes, err := os.ReadFile(global.FileID)

		// if err != nil {
		// 	// error handling
		// }

		// globalDeclaration := srcBytes[global.StartOffset : global.StartOffset+global.Length]

		// tokens := strings.Fields(string(globalDeclaration))

		// // type wasn't definited in declaration
		// if len(tokens) != 3 {
		// 	// error handling
		// }

		// globalType := tokens[2]

		// Make and add global to AST
		globalArg := ast.MakeArgument(global.GlobalVariableName, global.FileID, global.LineNumber)
		// globalArg.Offset = types.InvalidPointer
		// globalArg.Type = types[]
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

// func ParseStructs(structs []declaration_extraction.StructDeclaration) {

// 	// 1. iterate over all the structs
// 	// 2. add the struct name from the declaration
// 	// 3. search for fields
// 	// 4. fields to ast

// 	for _, strct := range structs {

// 		pkg, err := cxpartialparsing.Program.GetPackage(strct.PackageID)

// 		if err != nil {

// 			newPkg := ast.MakePackage(pkg.Name)
// 			pkgIdx := cxpartialparsing.Program.AddPackage(newPkg)
// 			newPkg, err = cxpartialparsing.Program.GetPackageFromArray(pkgIdx)

// 			if err != nil {
// 				// error handling
// 			}

// 			pkg = newPkg

// 		}

// 		structCX := ast.MakeStruct(strct.StructVariableName)
// 		structCX.Package = ast.CXPackageIndex(pkg.Index)

// 		cxpartialparsing.Program.AddStructInArray(structCX)

// 	}
// }

// func ParseFuncs(funcs []declaration_extraction.FuncDeclaration) {

// 	// 1. iterate over all the funcs
// 	// 2. extract inputs and outputs
// 	// 3. get the id and expression
// 	// 4. call function declaration

// 	for _, fun := range funcs {

// 		pkg, err := cxpartialparsing.Program.GetPackage(fun.PackageID)

// 		if err != nil {

// 			newPkg := ast.MakePackage(fun.PackageID)
// 			pkgIdx := cxpartialparsing.Program.AddPackage(newPkg)
// 			newPkg, err = cxpartialparsing.Program.GetPackageFromArray(pkgIdx)

// 			if err != nil {
// 				// error handling
// 			}

// 			pkg = newPkg

// 		}

// 		// srcBytes, err := os.ReadFile(fun.FileID)

// 		if err != nil {
// 			// error handling
// 		}

// 		// funcDeclaration := srcBytes[fun.StartOffset : fun.StartOffset+fun.Length]

// 		funcCX := ast.MakeFunction(fun.FuncVariableName, fun.FileID, fun.LineNumber)
// 		funcCX.Package = ast.CXPackageIndex(pkg.Index)

// 		cxpartialparsing.Program.AddFunctionInArray(funcCX)

// 	}

// }
