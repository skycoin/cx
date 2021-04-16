package cxparsering

import (
	"bytes"
	"io"
	"os"

	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/constants"
	globals2 "github.com/skycoin/cx/cx/globals"

	"github.com/skycoin/cx/cxparser/actions"
	cxpartialparsing "github.com/skycoin/cx/cxparser/cxpartialparsing"
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
func ParseSourceCode(sourceCode []*os.File, fileNames []string) {

	//local
	cxpartialparsing.Program = actions.AST

	/*
		Copy the contents of the file pointers containing the CX source
		code into sourceCodeStrings
	*/
	sourceCodeStrings := make([]string, len(sourceCode))
	for i, source := range sourceCode {
		tmp := bytes.NewBuffer(nil)
		io.Copy(tmp, source)
		sourceCodeStrings[i] = tmp.String()
	}

	/*
		We need to traverse the elements by hierarchy first add all the
		packages and structs at the same time then add globals, as these
		can be of a custom type (and it could be imported) the signatures
		of functions and methods are added in the cxpartialparsing.y pass
	*/
	parseErrors := 0
	if len(sourceCode) > 0 {
		parseErrors = Preliminarystage(sourceCodeStrings, fileNames)
	}

	//package level program
	actions.AST.SetCurrentCxProgram()

	actions.AST = cxpartialparsing.Program

	if globals2.FoundCompileErrors || parseErrors > 0 {
		profiling.CleanupAndExit(constants.CX_COMPILATION_ERROR)
	}

	/*
		Adding global variables `OS_ARGS` to the `os` (operating system)
		package.
	*/
	if osPkg, err := actions.AST.GetPackage(constants.OS_PKG); err == nil {
		if _, err := osPkg.GetGlobal(constants.OS_ARGS); err != nil {
			arg0 := ast.MakeArgument(constants.OS_ARGS, "", -1).AddType(constants.TypeNames[constants.TYPE_UNDEFINED])
			arg0.ArgDetails.Package = osPkg

			arg1 := ast.MakeArgument(constants.OS_ARGS, "", -1).AddType(constants.TypeNames[constants.TYPE_STR])
			arg1 = actions.DeclarationSpecifiers(arg1, []int{0}, constants.DECL_BASIC)
			arg1 = actions.DeclarationSpecifiers(arg1, []int{0}, constants.DECL_SLICE)
			actions.DeclareGlobalInPackage(osPkg, arg0, arg1, nil, false)
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

	if globals2.FoundCompileErrors || parseErrors > 0 {
		profiling.CleanupAndExit(constants.CX_COMPILATION_ERROR)
	}
}
