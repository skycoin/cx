package type_checker

import (
	"github.com/skycoin/cx/cmd/declaration_extractor"
	"github.com/skycoin/cx/cx/ast"
	cxinit "github.com/skycoin/cx/cx/init"
)

func ParseAllDeclarations(program *ast.CXProgram, imports []declaration_extractor.ImportDeclaration, globals []declaration_extractor.GlobalDeclaration, structs []declaration_extractor.StructDeclaration, funcs []declaration_extractor.FuncDeclaration) error {

	// Make program
	if program == nil {
		program = cxinit.MakeProgram()
	}

	err := ParseImports(imports)
	if err != nil {
		return err
	}

	err = ParseStructs(structs)
	if err != nil {
		return err
	}

	err = ParseGlobals(globals)
	if err != nil {
		return err
	}

	err = ParseFuncHeaders(funcs)
	if err != nil {
		return err
	}

	return nil
}
