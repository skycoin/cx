package type_checker

import (
	"os"

	"github.com/skycoin/cx/cmd/declaration_extractor"
	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cxparser/actions"
)

// Parse Structs
// - takes in structs from cx/cmd/declaration_extractor
// - adds structs to AST
func ParseStructs(structs []declaration_extractor.StructDeclaration) error {

	for _, strct := range structs {

		pkg, err := actions.AST.GetPackage(strct.PackageID)

		if err != nil {

			newPkg := ast.MakePackage(strct.PackageID)
			pkgIdx := actions.AST.AddPackage(newPkg)
			newPkg, err = actions.AST.GetPackageFromArray(pkgIdx)

			if err != nil {
				return err
			}

			pkg = newPkg

		}

		structCX := ast.MakeStruct(strct.StructName)
		structCX.Package = ast.CXPackageIndex(pkg.Index)

		pkg = pkg.AddStruct(actions.AST, structCX)

		src, err := os.ReadFile(strct.FileID)
		if err != nil {
			return err
		}

		var structFields []*ast.CXArgument

		for _, strctFieldDec := range strct.StructFields {

			structFieldLine := src[strctFieldDec.StartOffset : strctFieldDec.StartOffset+strctFieldDec.Length]

			structFieldSpecifier, err := ParseParameterDeclaration(structFieldLine, pkg, strct.FileID, strctFieldDec.LineNumber)

			if err != nil {
				return err
			}

			structFields = append(structFields, structFieldSpecifier)
		}

		actions.DeclareStruct(actions.AST, strct.StructName, structFields)

	}
	return nil
}
