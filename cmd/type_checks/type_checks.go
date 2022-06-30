package type_checks

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/skycoin/cx/cmd/declaration_extraction"
	"github.com/skycoin/cx/cx/ast"
	cxinit "github.com/skycoin/cx/cx/init"
	"github.com/skycoin/cx/cx/types"
	"github.com/skycoin/cx/cxparser/actions"
)

// Primitive Types Map
var primitiveTypesMap map[string]types.Code = map[string]types.Code{
	"bool": types.BOOL,

	"byte": types.I8,
	"i32":  types.I32,
	"i64":  types.I64,

	"f32": types.F32,
	"f64": types.F64,

	"str": types.STR,
}

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
		srcBytes, err := os.ReadFile(global.FileID)

		if err != nil {
			// error handling
		}

		globalDeclaration := srcBytes[global.StartOffset : global.StartOffset+global.Length]

		tokens := strings.Fields(string(globalDeclaration))

		// type wasn't definited in declaration
		if len(tokens) != 3 {
			// error handling
		}

		// Make and add global to AST
		globalArg := ast.MakeArgument(global.GlobalVariableName, global.FileID, global.LineNumber)
		globalArg.Offset = types.InvalidPointer

		globalArg.Package = ast.CXPackageIndex(pkg.Index)
		globalArg.SetType(primitiveTypesMap[tokens[2]])

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

		structCX := ast.MakeStruct(strct.StructVariableName)
		structCX.Package = ast.CXPackageIndex(pkg.Index)

		file, err := os.Open(strct.FileID)

		// if err != nil {
		// 	// error handling
		// }

		// srcBytes, err := io.ReadAll(file)

		// strctDeclaration := srcBytes[strct.StartOffset : strct.StartOffset+strct.Length]

		// if bytes.Compare(bytes.Fields(strctDeclaration)[2], []byte("struct")) == 0 {
		// 	structCX.SetType(getTypeCode(string(bytes.Fields(strctDeclaration)[2])))
		// }

		scanner := bufio.NewScanner(file)

		var inBlock int
		var lineno int

		for scanner.Scan() {

			line := scanner.Bytes()
			lineno++

			if lineno < strct.LineNumber {
				continue
			}

			if lineno == strct.LineNumber && bytes.IndexAny(line, "{") != -1 {
				inBlock++
				continue
			}

			if bytes.IndexAny(line, "}") != -1 {
				inBlock--
				break
			}

			if inBlock == 1 {
				tokens := bytes.Fields(line)

				if len(tokens) > 2 {
					// syntax error
					continue
				}

				if len(tokens) == 0 {
					continue
				}

				fmt.Print(string(tokens[1]))

				typeCode := primitiveTypesMap[string(tokens[1])]

				field := ast.MakeArgument(string(tokens[0]), strct.FileID, lineno)

				structCX = structCX.AddField(actions.AST, field)

				// var typ types.Code = types.
			}
		}

		// structCX = structCX.AddField()

		// for _, field := range fields {

		// }

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

		returnArray := bytes.Split(funcDeclaration[paramParenthesisClose+1:], []byte(","))

		funcCX := ast.MakeFunction(fun.FuncVariableName, fun.FileID, fun.LineNumber)

		for _, param := range paramArray {

			if bytes.Compare(param, []byte("")) == 0 {
				continue
			}

			tokens := bytes.Fields(param)
			paramName := bytes.TrimSpace(tokens[0])
			paramArg := ast.MakeArgument(string(paramName), fun.FileID, fun.LineNumber)
			paramArg.SetType(primitiveTypesMap[string(tokens[1])])
			funcCX = funcCX.AddInput(actions.AST, paramArg)

		}

		for _, rturn := range returnArray {

			if bytes.Compare(rturn, []byte("")) == 0 {
				continue
			}

			rturn = bytes.TrimSpace(rturn)
			returnArg := ast.MakeArgument("", fun.FileID, fun.LineNumber)
			returnArg.SetType(primitiveTypesMap[string(rturn)])
			funcCX = funcCX.AddOutput(actions.AST, returnArg)

		}

		pkg.AddFunction(actions.AST, funcCX)

	}

}
