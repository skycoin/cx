package type_checker

import (
	"strings"

	"github.com/skycoin/cx/cmd/declaration_extractor"
	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/packages"
	"github.com/skycoin/cx/cxparser/actions"
)

func ParseImports(imports []declaration_extractor.ImportDeclaration) error {

	// Make and add import packages to AST
	for _, imprt := range imports {
		// Get Package
		pkg, err := actions.AST.GetPackage(imprt.ImportName)
		if (err != nil || pkg == nil) && !packages.IsDefaultPackage(imprt.ImportName) {

			var imprtName string = imprt.ImportName

			if strings.Contains(imprt.ImportName, "/") {
				tokens := strings.Split(imprtName, "/")
				imprtName = tokens[len(tokens)-1]
			}

			newPkg := ast.MakePackage(imprtName)
			pkgIdx := actions.AST.AddPackage(newPkg)
			newPkg, err := actions.AST.GetPackageFromArray(pkgIdx)

			if err != nil {
				return err
			}

			pkg = newPkg
		}
	}

	// Declare import in the correct packages
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

		actions.AST.SelectPackage(imprt.PackageID)

		actions.DeclareImport(actions.AST, imprt.ImportName, imprt.FileID, imprt.LineNumber)

	}

	return nil

}
