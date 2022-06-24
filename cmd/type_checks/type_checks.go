package type_checks

import (
	"bufio"
	"go/scanner"
	"os"
	"regexp"
	"strings"

	"github.com/skycoin/cx/cmd/declaration_extraction"
	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/types"
	"github.com/skycoin/cx/cxparser/actions"
)

func ParseGlobals(globals []declaration_extraction.GlobalDeclaration) error {

	var err error

	// 1. interate over all the globals
	// 2. open file once and add to an array of src's
	// 3. if src is in array already don't open again
	// 4. extract the type and initializations in src
	// 5. add all the over information to CXProgram

	for _, global := range globals {

		// Gets the package from the AST
		pkg, err := actions.AST.GetPackage(global.PackageID)

		// If package doesn't exist in the AST
		if err != nil {

			// Create package and add to AST
			newPkg := ast.MakePackage(global.PackageID)
			pkgIdx := actions.AST.AddPackage(newPkg)
			newPkg, err = actions.AST.GetPackageFromArray(pkgIdx)

			pkg = newPkg

		}

		declarator := ast.MakeArgument(global.GlobalVariableName, global.FileID, global.LineNumber)

		srcBytes, err := os.ReadFile(global.FileID)

		if err != nil {
			continue
		}

		declaration := srcBytes[global.StartOffset : global.StartOffset+global.Length]

		tokens := strings.Fields(string(declaration))

		if tokens[3] == "" {
			continue
		}

		opTyp := ast.OpCodes[tokens[3]]

		declaration_specifiers := actions.DeclarationSpecifiersBasic(types.Code(opTyp))

		reInitializtion := regexp.MustCompile(global.GlobalVariableName + `\s+=\s+w+`)

		file, err := os.Open(global.FileID)

		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			
			line := 
		}

		actions.DeclareGlobalInPackage(actions.AST, pkg, declarator, declaration_specifiers)

	}

	return err
}

// func ParseEnums(enums []declaration_extraction.EnumDeclaration) {

// 	// Not sure about this one

// 	for _, enum := range enums {

// 		actions.de
// 	}
// }

// func ParseStructs(structs []declaration_extraction.StructDeclaration) error {

// 	// 1. iterate over all the structs
// 	// 2. add the struct name from the declaration
// 	// 3. search for fields
// 	// 4. fields to ast

// 	for _, strct := range structs {

// 		src, err := os.Open(strct.FileID)

// 		if err != nil {
// 			log.Fatal(err)
// 		}

// 		var inBlock int

// 		scanner = bufio.NewScanner(src)

// 		for scanner.Scan() {
// 			line := scanner
// 		}

// 		actions.DeclareStruct(program, strct.StructVariableName, nil)

// 	}
// }

// func ParseFuncs(funcs []declaration_extraction.FuncDeclaration) error {

// 	// 1. iterate over all the funcs
// 	// 2. extract inputs and outputs
// 	// 3. get the id and expression
// 	// 4. call function declaration
// 	var err error

// 	for _, fun := range funcs {

// 		srcBytes, err = os.ReadFile(fun.FileID)

// 		if err != nil {
// 			continue
// 		}

// 		fn := ast.MakeFunction(fun.FuncVariableName, fun.FileID, fun.LineNumber)

// 		fnIdx := actions.AST.AddFunctionInArray(fn)

// 		declaration := srcBytes[fun.StartOffset:fun.StartOffset+fun.Length]

// 		actions.FunctionDeclaration(actions.AST, fnIdx)

// 	}

// }
