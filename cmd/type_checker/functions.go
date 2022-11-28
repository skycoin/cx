package type_checker

import (
	"bytes"
	"regexp"

	"github.com/skycoin/cx/cmd/declaration_extractor"
	"github.com/skycoin/cx/cmd/packageloader2/loader"
	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cxparser/actions"
)

// Parse Function Headers
// - takes in funcs from cx/cmd/declaration_extractor
// - adds func headers to AST
func ParseFuncHeaders(files []*loader.File, funcs []declaration_extractor.FuncDeclaration) error {

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

		actions.AST.SelectPackage(fun.PackageID)

		source, err := GetSourceBytes(files, fun.FileID)
		if err != nil {
			return err
		}

		funcDeclarationLine := source[fun.StartOffset : fun.StartOffset+fun.Length]

		reFuncMethod := regexp.MustCompile(`func\s*\(\s*.+\s*\)`)
		funcMethod := reFuncMethod.Find(funcDeclarationLine)
		reParams := regexp.MustCompile(`\(([\s\w\*\[\],]*)\)`)
		params := reParams.FindAllSubmatch(funcDeclarationLine, -1)

		if funcMethod != nil {
			receiverArg, err := ParseParameterDeclaration(params[0][1], pkg, fun.FileID, fun.LineNumber)
			if err != nil {
				return err
			}

			var inputs []*ast.CXArgument
			var outputs []*ast.CXArgument

			fnName := receiverArg.StructType.Name + "." + fun.FuncName

			fn := ast.MakeFunction(fnName, actions.CurrentFile, fun.LineNumber)
			_, fnIdx := pkg.AddFunction(actions.AST, fn)
			newFn := actions.AST.GetFunctionFromArray(fnIdx)
			newFn.AddInput(actions.AST, receiverArg)

			if funcMethod[3] != nil && len(funcMethod[3]) != 0 {
				outputs, err = ParseFuncParameters(funcMethod[3], pkg, fun.FileID, fun.LineNumber)
				if err != nil {
					return err
				}
			}

			PreFunctionDeclaration(fnIdx, inputs, outputs)

		} else {

			reFuncRegular := regexp.MustCompile(`func\s*\S+\s*\(([\s\w,]*)\)(?:\s*\(([\s\w,]*)\))*`)
			funcRegular := reFuncRegular.FindSubmatch(funcDeclarationLine)

			fnIdx := actions.FunctionHeader(actions.AST, fun.FuncName, nil, false)

			var inputs []*ast.CXArgument
			var outputs []*ast.CXArgument

			if funcRegular[1] != nil && len(funcRegular[1]) != 0 {
				inputs, err = ParseFuncParameters(funcRegular[1], pkg, fun.FileID, fun.LineNumber)
				if err != nil {
					return err
				}
			}

			if funcRegular[2] != nil && len(funcRegular[2]) != 0 {
				outputs, err = ParseFuncParameters(funcRegular[2], pkg, fun.FileID, fun.LineNumber)
				if err != nil {
					return err
				}
			}

			PreFunctionDeclaration(fnIdx, inputs, outputs)
		}

	}
	return nil
}

func ParseFuncParameters(paramBytes []byte, pkg *ast.CXPackage, fileName string, lineno int) ([]*ast.CXArgument, error) {
	var parameterList []*ast.CXArgument

	tokens := bytes.Split(paramBytes, []byte(","))
	for _, token := range tokens {
		parameter, err := ParseParameterDeclaration(token, pkg, fileName, lineno)
		if err != nil {
			return nil, err
		}
		parameterList = append(parameterList, parameter)
	}

	return parameterList, nil
}

func PreFunctionDeclaration(fnIdx ast.CXFunctionIndex, inputs []*ast.CXArgument, outputs []*ast.CXArgument) {
	fn := actions.AST.GetFunctionFromArray(fnIdx)
	// adding inputs, outputs
	for _, inp := range inputs {
		fn.AddInput(actions.AST, inp)
	}
	for _, out := range outputs {
		fn.AddOutput(actions.AST, out)
	}
}
