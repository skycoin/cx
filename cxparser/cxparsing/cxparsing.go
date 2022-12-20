package cxparsering

import (
	"bytes"
	"fmt"
	"os"

	"github.com/skycoin/cx/cmd/declaration_extractor"
	"github.com/skycoin/cx/cmd/fileloader"
	"github.com/skycoin/cx/cmd/type_checker"
	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/constants"
	"github.com/skycoin/cx/cx/globals"
	"github.com/skycoin/cx/cx/types"

	"github.com/skycoin/cx/cxparser/actions"
	"github.com/skycoin/cx/cxparser/util/profiling"
)

/*
	ParseSourceCode takes a group of files representing CX `sourceCode` and
 	parses it into CX program structures for `AST`.

	 ParseSourceCode performs the steps

	 step 1 :  preliminarystage

	 step 2 :  passone

	 step 2 : passtwo
*/
func ParseSourceCode(sourceCode []*os.File) {

	//local
	// cxpartialparsing.Program = actions.AST

	/*
		Copy the contents of the file pointers containing the CX source
		code into sourceCodeStrings
	*/
	// sourceCodeStrings := make([]string, len(sourceCode))
	// for i, source := range sourceCode {
	// 	tmp := bytes.NewBuffer(nil)
	// 	io.Copy(tmp, source)
	// 	sourceCodeStrings[i] = tmp.String()
	// }
	var sourceCodeStrings []string
	var fileNames []string

	/*
		We need to traverse the elements by hierarchy first add all the
		packages and structs at the same time then add globals, as these
		can be of a custom type (and it could be imported) the signatures
		of functions and methods are added in the cxpartialparsing.y pass
	*/
	parseErrors := 0
	if len(sourceCode) > 0 {
		var err error
		sourceCodeStrings, fileNames, err = fileloader.LoadFiles(sourceCode)
		if err != nil {
			parseErrors++
			fmt.Print(err)
		}

		Imports, Globals, _, _, Structs, Funcs, err := declaration_extractor.ExtractAllDeclarations(sourceCodeStrings, fileNames)
		if err != nil {
			parseErrors++
			fmt.Println(err)
		}

		err = type_checker.ParseAllDeclarations(Imports, Globals, Structs, Funcs)
		if err != nil {
			parseErrors++
			fmt.Println(err)
		}

	}

	// actions.AST = cxpartialparsing.Program

	if globals.FoundCompileErrors || parseErrors > 0 {
		profiling.CleanupAndExit(constants.CX_COMPILATION_ERROR)
	}

	/*
		Adding global variables `OS_ARGS` to the `os` (operating system)
		package.
	*/
	if osPkg, err := actions.AST.GetPackage(constants.OS_PKG); err == nil {
		if _, err := osPkg.GetGlobal(actions.AST, constants.OS_ARGS); err != nil {
			arg0 := ast.MakeArgument(constants.OS_ARGS, "", -1).SetType(types.UNDEFINED)
			arg0.Package = ast.CXPackageIndex(osPkg.Index)

			arg1 := ast.MakeArgument(constants.OS_ARGS, "", -1).SetType(types.STR)
			arg1 = actions.DeclarationSpecifiers(arg1, []types.Pointer{0}, constants.DECL_BASIC)
			arg1 = actions.DeclarationSpecifiers(arg1, []types.Pointer{0}, constants.DECL_SLICE)
			actions.DeclareGlobalInPackage(actions.AST, osPkg, arg0, arg1, nil, false)
		}
	}

	profiling.StartProfile("4. passtwo")

	/*
	 The pass two of parsing that generates the actual output.
	*/

	for i, source := range sourceCodeStrings {

		/*
			Because of an unkown reason, sometimes some CX programs
			throw an error related to a premature EOF (particularly in Windows).
			Adding a newline character solves this.
		*/

		source = source + "\n"

		actions.LineNo = 1

		b := bytes.NewBufferString(source)

		if len(fileNames) > 0 {
			actions.CurrentFile = fileNames[i]
		}

		profiling.StartProfile(actions.CurrentFile)

		parseErrors += Passtwo(b)

		profiling.StopProfile(actions.CurrentFile)
	}

	profiling.StopProfile("4. passtwo")

	if globals.FoundCompileErrors || parseErrors > 0 {
		profiling.CleanupAndExit(constants.CX_COMPILATION_ERROR)
	}
}
