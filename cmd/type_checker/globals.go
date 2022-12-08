package type_checker

import (
	"regexp"

	"github.com/skycoin/cx/cmd/declaration_extractor"
	"github.com/skycoin/cx/cmd/packageloader/loader"
	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/types"
	"github.com/skycoin/cx/cxparser/actions"
)

// Parse Globals
// - takes in globals from cx/cmd/declaration_extractor
// - adds globals to AST
func ParseGlobals(files []*loader.File, globals []declaration_extractor.GlobalDeclaration) error {

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
				return err
			}

			pkg = newPkg
		}

		actions.AST.SelectPackage(global.PackageID)

		// Read File
		source, err := GetSourceBytes(files, global.FileID)
		if err != nil {
			return err
		}

		// Extract Declaration from file
		reGlobalDeclaration := regexp.MustCompile(`var\s+(\w*)\s+([\*\[\]\w\.]+)(?:\s*\=\s*[\s\S]+\S+){0,1}`)
		globalDeclaration := source[global.StartOffset : global.StartOffset+global.Length]
		globalTokens := reGlobalDeclaration.FindSubmatch(globalDeclaration)

		// add Global as CX Argument to CX Package
		globalArg := ast.MakeArgument(global.GlobalVariableName, global.FileID, global.LineNumber)
		globalArg.Offset = types.InvalidPointer
		globalArg.Package = ast.CXPackageIndex(pkg.Index)
		globalArgIdx := actions.AST.AddCXArgInArray(globalArg)

		pkg.AddGlobal(actions.AST, globalArgIdx)

		var declarationSpecifier *ast.CXArgument

		declarationSpecifier, err = ParseDeclarationSpecifier(globalTokens[2], global.FileID, global.LineNumber, declarationSpecifier)
		if err != nil {
			return err
		}

		actions.DeclareGlobalInPackage(actions.AST, nil, actions.AST.GetCXArgFromArray(globalArgIdx), declarationSpecifier, nil, false)

	}
	return nil
}
