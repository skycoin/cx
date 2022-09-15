package type_checks

import (
	"os"
	"regexp"

	"github.com/skycoin/cx/cmd/declaration_extractor"
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

		source, err := os.ReadFile(fun.FileID)
		if err != nil {
			return err
		}

		funcDeclarationLine := source[fun.StartOffset : fun.StartOffset+fun.Length]

		reFuncMethod := regexp.MustCompile(`func\s*\(\s*(\w+)\s*.+\)`)
		funcMethod := reFuncMethod.FindSubmatch(funcDeclarationLine)

		if funcMethod != nil {
			receiver := string(funcMethod[1])
			funcName := receiver + "." + fun.FuncName
			actions.FunctionHeader(actions.AST, funcName, nil, true)
		} else {
			actions.FunctionHeader(actions.AST, fun.FuncName, nil, false)
		}

	}
	return nil
}
