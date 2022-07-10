package type_checks

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/skycoin/cx/cmd/declaration_extraction"
	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/constants"
	cxinit "github.com/skycoin/cx/cx/init"
	"github.com/skycoin/cx/cx/types"
	"github.com/skycoin/cx/cxparser/actions"
)

// Primitive Types Map
var primitiveTypesMap map[string]types.Code = map[string]types.Code{

	"bool": types.BOOL,

	"i8":  types.I8,
	"i16": types.I16,
	"i32": types.I32,
	"i64": types.I64,

	"ui8":  types.UI8,
	"ui16": types.UI16,
	"ui32": types.UI32,
	"ui64": types.UI64,

	"f32": types.F32,
	"f64": types.F64,

	"str": types.STR,
	"aff": types.AFF,
}

// ParseGlobals make the CXArguments and add them to the AST
// 1. iterates over all the global declarations
// 2. gets/makes the package
// 3.

func ParseGlobals(globals []declaration_extraction.GlobalDeclaration) {

	// Make program
	if actions.AST == nil {
		actions.AST = cxinit.MakeProgram()
	}

	// Range over global declarations and parse
	for _, global := range globals {

		// Get Package
		pkg, err := actions.AST.GetPackage(global.PackageID)

		// If package not in AST
		if err != nil || pkg == nil {

			newPkg := ast.MakePackage(global.PackageID)
			pkgIdx := actions.AST.AddPackage(newPkg)
			newPkg, err := actions.AST.GetPackageFromArray(pkgIdx)

			if err != nil {
				panic(err)
			}

			pkg = newPkg
		}

		// read bytes
		srcBytes, err := os.ReadFile(global.FileID)

		if err != nil {
			panic(err)
		}

		globalDeclaration := srcBytes[global.StartOffset : global.StartOffset+global.Length]

		tokens := strings.Fields(string(globalDeclaration))

		// type wasn't definited in declaration
		if len(tokens) < 3 {
			// error handling
			fmt.Print("error no type was declared")
		}

		// Make and add global to AST
		globalArg := ast.MakeArgument(global.GlobalVariableName, global.FileID, global.LineNumber)
		globalArg.Offset = types.InvalidPointer
		globalArg.Package = ast.CXPackageIndex(pkg.Index)

		globalArgIdx := actions.AST.AddCXArgInArray(globalArg)

		pkg = pkg.AddGlobal(actions.AST, globalArgIdx)

		// Global declaration and type setting
		reArray := regexp.MustCompile(`\[([0-9]*)\]([_a-zA-Z0-9]*)`)

		//Default declaration specifier for primitive types
		declarationSpecifier := actions.DeclarationSpecifiersBasic(primitiveTypesMap[tokens[2]])

		//Declaration specifier for arrays
		if arrayDeclarationSpecifier := reArray.FindStringSubmatch(tokens[2]); arrayDeclarationSpecifier != nil {
			declarationSpecifierBasic := actions.DeclarationSpecifiersBasic(primitiveTypesMap[arrayDeclarationSpecifier[2]])
			numberOfElements, err := strconv.Atoi(arrayDeclarationSpecifier[1])

			if err != nil {
				// error handling
			}

			declarationSpecifier = actions.DeclarationSpecifiers(declarationSpecifierBasic, []types.Pointer{types.Pointer(numberOfElements)}, constants.DECL_ARRAY)
		}

		actions.DeclareGlobalInPackage(actions.AST, pkg, actions.AST.GetCXArg(globalArgIdx), declarationSpecifier, nil, false)

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

		pkg = pkg.AddStruct(actions.AST, structCX)

		file, err := os.Open(strct.FileID)

		// if err != nil {
		// 	// error handling
		// }

		// srcBytes, err := io.ReadAll(file)

		// strctDeclaration := srcBytes[strct.StartOffset : strct.StartOffset+strct.Length]

		// if bytes.Compare(bytes.Fields(strctDeclaration)[2], []byte("struct")) == 0 {
		// 	structCX.SetType(getTypeCode(string(bytes.Fields(strctDeclaration)[2])))
		// }
		var structFields []*ast.CXArgument

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

				structFieldCXArg := ast.MakeArgument(string(tokens[0]), strct.FileID, strct.LineNumber)
				structFieldCXArg = structFieldCXArg.SetPackage(pkg)

				structField := actions.DeclarationSpecifiersBasic(primitiveTypesMap[string(tokens[1])])

				structField.Name = structFieldCXArg.Name
				structField.Package = structFieldCXArg.Package
				structField.IsLocalDeclaration = true

				structFields = append(structFields, structField)

			}
		}

		actions.DeclareStruct(actions.AST, strct.StructVariableName, structFields)

	}
}

func ParseFuncHeaders(funcs []declaration_extraction.FuncDeclaration) {

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

		funcCX := ast.MakeFunction(fun.FuncVariableName, fun.FileID, fun.LineNumber)

		pkg.AddFunction(actions.AST, funcCX)

		srcBytes, err := os.ReadFile(fun.FileID)

		if err != nil {
			// error handling
		}

		reParenOpen := regexp.MustCompile(`\(`)
		reParenClose := regexp.MustCompile(`\)`)

		// Get function
		funcDeclaration := srcBytes[fun.StartOffset : fun.StartOffset+fun.Length]

		inputParenthesisOpen := reParenOpen.FindIndex(funcDeclaration)[0]
		inputParenthesisClose := reParenClose.FindIndex(funcDeclaration)[0]

		inputsArray := bytes.Split(funcDeclaration[inputParenthesisOpen+1:inputParenthesisClose], []byte(","))

		// returnParenthesisOpen := bytes.IndexAny(funcDeclaration[paramParenthesisClose:], "(")
		// returnParenthesisClose := bytes.IndexAny(funcDeclaration[returnParenthesisOpen:], ")")

		outputRemoveParenOpen := reParenOpen.ReplaceAll(funcDeclaration[inputParenthesisClose+1:], []byte(""))
		outputRemoveParenClose := reParenClose.ReplaceAll(outputRemoveParenOpen, []byte(""))

		outputsArray := bytes.Split(outputRemoveParenClose, []byte(","))

		inputArgsArray := []*ast.CXArgument{}
		outputArgsArray := []*ast.CXArgument{}

		for _, input := range inputsArray {

			if bytes.Compare(input, []byte("")) == 0 {
				continue
			}

			// Tokenize the parameter further
			tokens := bytes.Fields(input)

			// Declarator
			inputArg := ast.MakeArgument(string(tokens[0]), fun.FileID, fun.LineNumber)
			inputArg = inputArg.SetPackage(pkg)

			// Declaration Specifier and Merge with Declarator
			inputSpecifier := actions.DeclarationSpecifiersBasic(primitiveTypesMap[string(tokens[1])])
			inputSpecifier.Name = inputArg.Name
			inputSpecifier.Package = inputArg.Package

			inputArgsArray = append(inputArgsArray, inputSpecifier)

		}

		for _, output := range outputsArray {

			if bytes.Compare(output, []byte("")) == 0 {
				continue
			}

			// Tokenize
			tokens := bytes.Fields(output)

			// Declarator
			outputArg := ast.MakeArgument(string(tokens[0]), fun.FileID, fun.LineNumber)
			outputArg = outputArg.SetPackage(pkg)

			// Declaration Specifier
			outputSpecifier := actions.DeclarationSpecifiersBasic(primitiveTypesMap[string(tokens[1])])
			outputSpecifier.Name = outputArg.Name
			outputSpecifier.Package = outputArg.Package

			outputArgsArray = append(outputArgsArray, outputSpecifier)
		}

		funIdx := actions.FunctionHeader(actions.AST, fun.FuncVariableName, nil, false)
		fun := actions.AST.GetFunctionFromArray(funIdx)

		for _, inputArg := range inputArgsArray {
			fun.AddInput(actions.AST, inputArg)
		}

		for _, outputArg := range outputArgsArray {
			fun.AddOutput(actions.AST, outputArg)
		}

	}

}
