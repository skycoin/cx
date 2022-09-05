package type_checks

import (
	"fmt"
	"os"
	"regexp"

	"github.com/skycoin/cx/cmd/declaration_extractor"
	"github.com/skycoin/cx/cx/ast"
	cxinit "github.com/skycoin/cx/cx/init"
	"github.com/skycoin/cx/cx/types"
	"github.com/skycoin/cx/cxparser/actions"
)

// Parse Globals
// - takes in globals from cx/cmd/declaration_extractor
// - adds globals to AST
func ParseGlobals(globals []declaration_extractor.GlobalDeclaration) error {

	// Make program
	if actions.AST == nil {
		actions.AST = cxinit.MakeProgram()
	}

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

		// Read File
		source, err := os.ReadFile(global.FileID)

		if err != nil {
			return err
		}

		// Extract Declaration from file
		reGlobalDeclaration := regexp.MustCompile(`var\s+(\w*)\s+\*{0,1}\s*(?:\[\d*\]){0,1}\s*(\w*)(?:\s*\=\s*[\s\S]+\S+){0,1}`)
		globalDeclaration := source[global.StartOffset : global.StartOffset+global.Length]
		globalTokens := reGlobalDeclaration.FindSubmatch(globalDeclaration)

		// Add Global to Pkg
		if _, err := pkg.GetGlobal(actions.AST, global.GlobalVariableName); err != nil {
			// add Global as CX Argument to CX Package
			globalArg := ast.MakeArgument(global.GlobalVariableName, "", 0)
			globalArg.Offset = types.InvalidPointer
			globalArg.Package = ast.CXPackageIndex(pkg.Index)
			globalArgIdx := actions.AST.AddCXArgInArray(globalArg)

			pkg.AddGlobal(actions.AST, globalArgIdx)
		}

		// Create Declarator
		var declarator *ast.CXArgument
		if pkg, err := actions.AST.GetCurrentPackage(); err == nil {
			declarator = ast.MakeArgument("", actions.CurrentFile, actions.LineNo)
			declarator.SetType(types.UNDEFINED)
			declarator.Name = global.GlobalVariableName
			declarator.Package = ast.CXPackageIndex(pkg.Index)
		} else {
			return err
		}

		var declaration_specifier *ast.CXArgument
		if val, ok := TypesMap[string(globalTokens[1])]; ok {
			declaration_specifier = actions.DeclarationSpecifiersBasic(val)
		} else {
			return fmt.Errorf("declaration_specifier")
		}

		actions.DeclareGlobalInPackage(actions.AST, nil, declarator, declaration_specifier, nil, false)

	}
	return nil
}
