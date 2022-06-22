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

	for _, global := range globals {

		pkg, err := program.GetPackage(global.PackageID)
		declarator := ast.MakeArgument(global.GlobalVariableName, global.FileID, global.LineNumber)
		// declarationSpecifiers := ast.MakeArgument()
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

}

func ParseStructs(structs []declaration_extraction.StructDeclaration, program *ast.CXProgram) {

	for _, strct := range structs {

		declarator := ast.MakeArgument(strct.StructVariableName, strct.FileID, strct.LineNumber)

		src, err := os.Open(strct.FileID)

		if err != nil {
			log.Fatal(err)
		}

		var inBlock int

		scanner = bufio.NewScanner(src)

		for scanner.Scan() {
			line := scanner.Scanner
		}

		actions.DeclareStruct(program, strct.StructVariableName, nil)

	}
}

func ParseFuncs(funcs []declaration_extraction.FuncDeclaration) {

	for _, fun := range funcs {

		pkg, err := program.GetPackage(fun.PackageID)
		declarator := ast.MakeArgument(fun.FuncVariableName, fun.FileID, fun.LineNumber)

		if err != nil {
			pkg = ast.MakePackage(fun.PackageID)
		}

		actions.FunctionDeclaration(&program)
	}

}
