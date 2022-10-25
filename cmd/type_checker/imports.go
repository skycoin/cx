package type_checker

import (
	"github.com/skycoin/cx/cmd/declaration_extractor"
	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cxparser/actions"
)

func ParseImports(imports []declaration_extractor.ImportDeclaration) error {

	for _, imprt := range imports {
		// Get Package
		pkg, err := actions.AST.GetPackage(imprt.ImportName)

		// If package not in AST
		if err != nil || pkg == nil {

			newPkg := ast.MakePackage(imprt.ImportName)
			pkgIdx := actions.AST.AddPackage(newPkg)
			newPkg, err := actions.AST.GetPackageFromArray(pkgIdx)

			if err != nil {
				return err
			}

			pkg = newPkg
		}
	}

	for _, imprt := range imports {

		// Get Package
		pkg, err := actions.AST.GetPackage(imprt.PackageID)

		// If package not in AST
		if err != nil || pkg == nil {

			newPkg := ast.MakePackage(imprt.PackageID)
			pkgIdx := actions.AST.AddPackage(newPkg)
			newPkg, err := actions.AST.GetPackageFromArray(pkgIdx)

			if err != nil {
				return err
			}

			pkg = newPkg
		}

		actions.DeclareImport(actions.AST, imprt.ImportName, imprt.FileID, imprt.LineNumber)

	}

	return nil

}
