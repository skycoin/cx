package type_checks

import (
	"bufio"
	"go/scanner"
	"log"
	"os"

	"github.com/skycoin/cx/cmd/declaration_extraction"
	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cxparser/actions"
)

func ParseGlobals(globals []declaration_extraction.GlobalDeclaration, program *ast.CXProgram) {

	// 1. interate over all the globals
	// 2. open file once and add to an array of src's
	// 3. if src is in array already don't open again
	// 4. extract the type and initializations in src
	// 5. add all the over information to CXProgram

	for _, global := range globals {

		pkg, err := program.GetPackage(global.PackageID)
		declarator := ast.MakeArgument(global.GlobalVariableName, global.FileID, global.LineNumber)
		declarationSpecifiers := actions.DeclarationSpecifiers()
		// initializer := ast.MakeCXLineExpression(&program, global.FileID, global.LineNumber)

		if err != nil {
			pkg = ast.MakePackage(global.PackageID)
		}

		src, err = os.ReadFile(global.FileID)

		if err != nil {
			log.Fatal(err)
		}

		actions.DeclareGlobalInPackage(program, pkg, declarator, nil, nil, false)

	}
}

func ParseEnums(enums []declaration_extraction.EnumDeclaration) {

	// Not sure about this one

	for _, enum := range enums {

		actions.de
	}
}

func ParseStructs(structs []declaration_extraction.StructDeclaration, program *ast.CXProgram) {

	// 1. iterate over all the structs
	// 2. add the struct name from the declaration
	// 3. search for fields
	// 4. fields to ast

	for _, strct := range structs {

		src, err := os.Open(strct.FileID)

		if err != nil {
			log.Fatal(err)
		}

		var inBlock int

		scanner = bufio.NewScanner(src)

		for scanner.Scan() {
			line := scanner
		}

		actions.DeclareStruct(program, strct.StructVariableName, nil)

	}
}

func ParseFuncs(funcs []declaration_extraction.FuncDeclaration, program *ast.CXProgram) {

	// 1. iterate over all the funcs
	// 2. extract inputs and outputs
	// 3. get the id and expression
	// 4. call function declaration

	for _, fun := range funcs {

		declarator := ast.MakeArgument(fun.FuncVariableName, fun.FileID, fun.LineNumber)

		src, err := os.Open(fun.FileID)

		actions.FunctionDeclaration(program)

	}

}
