package type_checks

import "github.com/skycoin/cx/cmd/declaration_extractor"

func ParseAllDeclarations(globals []declaration_extractor.GlobalDeclaration, structs []declaration_extractor.StructDeclaration, funcs []declaration_extractor.FuncDeclaration) {

	ParseStructs(structs)
	ParseGlobals(globals)
	ParseFuncHeaders(funcs)
}
