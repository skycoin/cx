package type_checks

import (
	"github.com/skycoin/cx/cmd/declaration_extraction"
)

type GlobalType struct {
	Name string
	Type string
}

type EnumType struct {
	Name string
	Type string
}

type StructType struct {
	Name string
	Type string
}

func ParseGlobals(globals []declaration_extraction.GlobalDeclaration) ([]GlobalType, error) {

	var GlobalTypes []GlobalType

	for _, global := range globals {

	}
}

func ParseEnums(enums []declaration_extraction.EnumDeclaration) {

}

func ParseStructs(structs []declaration_extraction.StructDeclaration) {

}

func ParseFuncs(funcs []declaration_extraction.FuncDeclaration) {

}
