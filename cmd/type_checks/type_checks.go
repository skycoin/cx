package type_checks

import (
	"bufio"
	"log"
	"os"

	"github.com/skycoin/cx/cmd/declaration_extraction"
	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cxparser/actions"
	cxpartialparsing "github.com/skycoin/cx/cxparser/cxpartialparsing"
)

func ParseGlobals(globals []declaration_extraction.GlobalDeclaration) {

	for _, global := range globals {

		pkg, err := cxpartialparsing.Program.GetPackage(global.PackageID)

		if err != nil {

			newPkg := ast.MakePackage(global.PackageID)
			pkgIdx := cxpartialparsing.Program.AddPackage(newPkg)
			newPkg, err = cxpartialparsing.Program.GetPackageFromArray(pkgIdx)

			if err != nil {
				// error handling
			}

			pkg = newPkg
		}

		globalArg := ast.MakeArgument(global.GlobalVariableName, global.FileID, global.LineNumber)
		globalArg.Offset = global.StartOffset
		globalArg.Package = ast.CXPackageIndex(pkg.Index)

		globalArgIdx := actions.AST.AddCXArgInArray(globalArg)

		pkg.AddGlobal(actions.AST, globalArgIdx)

	}
}

// func ParseEnums(enums []declaration_extraction.EnumDeclaration) {

// 	// Not sure about this one

// 	for _, enum := range enums {

// 		actions.de
// 	}
// }

func ParseStructs(structs []declaration_extraction.StructDeclaration) {

	// 1. iterate over all the structs
	// 2. add the struct name from the declaration
	// 3. search for fields
	// 4. fields to ast

	for _, strct := range structs {

		src, err := os.Open(strct.FileID)

		if err != nil {
			log.Fatal(err)
		}

		var inBlock int

		scanner = bufio.NewScanner(src)

		for scanner.Scan() {
			line := scanner
		}

		actions.DeclareStruct(program, strct.StructVariableName, nil)

	}
}

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
