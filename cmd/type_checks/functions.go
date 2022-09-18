package type_checks

import (
	"os"
	"regexp"

	"github.com/skycoin/cx/cmd/declaration_extractor"
	"github.com/skycoin/cx/cx/ast"
	cxinit "github.com/skycoin/cx/cx/init"
	"github.com/skycoin/cx/cxparser/actions"
)

// Parse Function Headers
// - takes in funcs from cx/cmd/declaration_extractor
// - adds func headers to AST
func ParseFuncHeaders(funcs []declaration_extractor.FuncDeclaration) error {

	// Make program
	if actions.AST == nil {
		actions.AST = cxinit.MakeProgram()
	}

	for _, fun := range funcs {

		// Get Package
		pkg, err := actions.AST.GetPackage(fun.PackageID)

		// If package not in AST
		if err != nil || pkg == nil {

			newPkg := ast.MakePackage(fun.PackageID)
			pkgIdx := actions.AST.AddPackage(newPkg)
			newPkg, err := actions.AST.GetPackageFromArray(pkgIdx)

			if err != nil {
				return err
			}

			pkg = newPkg
		}

		source, err := os.ReadFile(fun.FileID)
		if err != nil {
			return err
		}

		funcDeclarationLine := source[fun.StartOffset : fun.StartOffset+fun.Length]

		reFuncMethod := regexp.MustCompile(`func\s*\(\s*(.+)\s*\)`)
		funcMethod := reFuncMethod.FindSubmatch(funcDeclarationLine)

		if funcMethod != nil {
			receiverArg, err := ParseParameterDeclaration(funcMethod[1], pkg, fun.FileID, fun.LineNumber)
			if err != nil {
				return err
			}
			actions.FunctionHeader(actions.AST, fun.FuncName, []*ast.CXArgument{receiverArg}, true)
		} else {
			actions.FunctionHeader(actions.AST, fun.FuncName, nil, false)
		}

	}
	return nil
}
