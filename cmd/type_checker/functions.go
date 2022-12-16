package type_checker

import (
	"bytes"
	"io"
	"os"
	"regexp"

	"github.com/skycoin/cx/cmd/declaration_extractor"
	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cxparser/actions"
)

// Parse Function Headers
// - takes in funcs from cx/cmd/declaration_extractor
// - adds func headers to AST
func ParseFuncHeaders(funcs []declaration_extractor.FuncDeclaration) error {

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

		file, err := os.Open(fun.FileID)
		if err != nil {
			return err
		}

		tmp := bytes.NewBuffer(nil)
		io.Copy(tmp, file)
		source := tmp.Bytes()

		funcDeclarationLine := source[fun.StartOffset : fun.StartOffset+fun.Length]

		reFuncMethod := regexp.MustCompile(`func\s*\(\s*.+\s*\)`)
		funcMethod := reFuncMethod.Find(funcDeclarationLine)
		reParams := regexp.MustCompile(`\(([\s\w\*\[\],\.]*)\)`)
		params := reParams.FindAllSubmatch(funcDeclarationLine, -1)

		if funcMethod != nil {

			receiverArg, err := ParseParameterDeclaration(params[0][1], pkg, fun.FileID, fun.LineNumber)
			if err != nil {
				return err
			}

			fnName := receiverArg.StructType.Name + "." + fun.FuncName

			fn := ast.MakeFunction(fnName, actions.CurrentFile, fun.LineNumber)
			_, fnIdx := pkg.AddFunction(actions.AST, fn)
			newFn := actions.AST.GetFunctionFromArray(fnIdx)
			newFn.AddInput(actions.AST, receiverArg)

			var inputs []*ast.CXArgument
			var outputs []*ast.CXArgument

			if params[1][1] != nil && len(params[1][1]) != 0 {
				inputs, err = ParseFuncParameters(params[1][1], pkg, fun.FileID, fun.LineNumber)
				if err != nil {
					return err
				}
			}

			if len(params) == 3 {
				if params[2][1] != nil && len(params[2][1]) != 0 {
					outputs, err = ParseFuncParameters(params[2][1], pkg, fun.FileID, fun.LineNumber)
					if err != nil {
						return err
					}
				}
			}

			PreFunctionDeclaration(fnIdx, inputs, outputs)

		} else {

			fn := ast.MakeFunction(fun.FuncName, fun.FileID, fun.LineNumber)
			_, fnIdx := pkg.AddFunction(actions.AST, fn)

			var inputs []*ast.CXArgument
			var outputs []*ast.CXArgument

			if params[0][1] != nil && len(params[0][1]) != 0 {
				inputs, err = ParseFuncParameters(params[0][1], pkg, fun.FileID, fun.LineNumber)
				if err != nil {
					return err
				}
			}

			if len(params) == 2 {
				if params[1][1] != nil && len(params[1][1]) != 0 {
					outputs, err = ParseFuncParameters(params[1][1], pkg, fun.FileID, fun.LineNumber)
					if err != nil {
						return err
					}
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
