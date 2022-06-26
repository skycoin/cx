package type_checks

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/skycoin/cx/cmd/declaration_extraction"
	"github.com/skycoin/cx/cx/ast"
	cxinit "github.com/skycoin/cx/cx/init"
	"github.com/skycoin/cx/cx/types"
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
		globalArg.Offset = types.InvalidPointer
		globalArg.SetType(0)
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

	if actions.AST == nil {
		actions.AST = cxinit.MakeProgram()
	}

	for _, strct := range structs {

		pkg, err := actions.AST.GetPackage(strct.PackageID)

		if err != nil {

			newPkg := ast.MakePackage(pkg.Name)
			pkgIdx := actions.AST.AddPackage(newPkg)
			newPkg, err = actions.AST.GetPackageFromArray(pkgIdx)

			if err != nil {
				// error handling
			}

			pkg = newPkg

		}

		// srcBytes, err := os.ReadFile(strct.FileID)

		// if err != nil {
		// 	// error handling
		// }

		structCX := ast.MakeStruct(strct.StructVariableName)
		structCX.Package = ast.CXPackageIndex(pkg.Index)
		// structCX = structCX.AddField()

		pkg.AddStruct(actions.AST, structCX)

	}
}

func ParseFuncs(funcs []declaration_extraction.FuncDeclaration) {

	// 1. iterate over all the funcs
	// 2. extract inputs and outputs
	// 3. get the id and expression
	// 4. call function declaration

	if actions.AST == nil {
		actions.AST = cxinit.MakeProgram()
	}

	for _, fun := range funcs {

		pkg, err := actions.AST.GetPackage(fun.PackageID)

		if err != nil {

			newPkg := ast.MakePackage(fun.PackageID)
			pkgIdx := actions.AST.AddPackage(newPkg)
			newPkg, err = actions.AST.GetPackageFromArray(pkgIdx)

			if err != nil {
				// error handling
			}

			pkg = newPkg

		}

		srcBytes, err := os.ReadFile(fun.FileID)

		if err != nil {
			// error handling
		}

		// Get function
		funcDeclaration := srcBytes[fun.StartOffset : fun.StartOffset+fun.Length]

		paramParenthesisOpen := bytes.IndexAny(funcDeclaration, "(")
		paramParenthesisClose := bytes.IndexAny(funcDeclaration, ")")

		paramArray := bytes.Split(funcDeclaration[paramParenthesisOpen+1:paramParenthesisClose], []byte(","))

		// returnParenthesisOpen := bytes.IndexAny(funcDeclaration[paramParenthesisClose:], "(")
		// returnParenthesisClose := bytes.IndexAny(funcDeclaration[returnParenthesisOpen:], ")")

		funcCX := ast.MakeFunction(fun.FuncVariableName, fun.FileID, fun.LineNumber)

		for _, param := range paramArray {

			if bytes.Compare(param, []byte("")) == 0 {
				continue
			}

			tokens := bytes.Fields(param)
			paramName := tokens[0]
			paramArg := ast.MakeArgument(string(paramName), fun.FileID, fun.LineNumber)
			funcCX = funcCX.AddInput(actions.AST, paramArg)

			fmt.Print(string(param))
		}

		pkg.AddFunction(actions.AST, funcCX)

	}

}
