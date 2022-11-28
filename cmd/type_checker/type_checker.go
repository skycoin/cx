package type_checker

import (
	"github.com/skycoin/cx/cmd/declaration_extractor"
	"github.com/skycoin/cx/cmd/packageloader2/loader"
	cxinit "github.com/skycoin/cx/cx/init"
	"github.com/skycoin/cx/cxparser/actions"
)

func ParseAllDeclarations(files []*loader.File, imports []declaration_extractor.ImportDeclaration, globals []declaration_extractor.GlobalDeclaration, structs []declaration_extractor.StructDeclaration, funcs []declaration_extractor.FuncDeclaration) error {

	// Make AST if not made yet
	if actions.AST == nil {
		actions.AST = cxinit.MakeProgram()
	}

	err := ParseImports(imports)
	if err != nil {
		return err
	}

	err = ParseStructs(files, structs)
	if err != nil {
		return err
	}

	err = ParseGlobals(files, globals)
	if err != nil {
		return err
	}

	err = ParseFuncHeaders(files, funcs)
	if err != nil {
		return err
	}

	return nil
}
