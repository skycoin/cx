package cxparsering

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/skycoin/cx/cmd/declaration_extractor"
	"github.com/skycoin/cx/cmd/packageloader2/loader"
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
func ParseSourceCode(sourceCode []*os.File, fileNames []string, rootDir []string) {

	//local
	// cxpartialparsing.Program = actions.AST

	parseErrors := 0

	// err := loader.LoadCXProgram("main", sourceCode, rootDir, "bolt")
	// if err != nil {
	// 	fmt.Print(err)
	// 	parseErrors++
	// }

	// files, err := file_output.GetImportFiles("test", "bolt")
	// if err != nil {
	// 	fmt.Print(err)
	// 	parseErrors++
	// }

	files, err := loader.LoadCXProgramNoSave(sourceCode, rootDir)
	if err != nil {
		fmt.Print(err)
		parseErrors++
	}

	/*
		Copy the contents of the file pointers containing the CX source
		code into sourceCodeStrings
	*/
	sourceCodeStrings := make([]string, len(files))
	for i, source := range files {
		tmp := bytes.NewBuffer(nil)
		io.Copy(tmp, bytes.NewReader(source.Content))
		sourceCodeStrings[i] = tmp.String()
	}

	/*
		We need to traverse the elements by hierarchy first add all the
		packages and structs at the same time then add globals, as these
		can be of a custom type (and it could be imported) the signatures
		of functions and methods are added in the cxpartialparsing.y pass
	*/

	// if len(sourceCode) > 0 {
	// 	parseErrors = Preliminarystage(sourceCodeStrings, fileNames)
	// }

	// actions.AST = cxpartialparsing.Program

	Imports, Globals, Enums, TypeDefinitions, Structs, Funcs, err := declaration_extractor.ExtractAllDeclarations(files)
	if err != nil {
		fmt.Print(err)
		parseErrors++
	}

	if Enums != nil && TypeDefinitions != nil {

	}

	err = type_checker.ParseAllDeclarations(files, Imports, Globals, Structs, Funcs)
	if err != nil {
		fmt.Print(err)
		parseErrors++
	}

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

		if len(files) > 0 {
			actions.CurrentFile = files[i].FileName
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
