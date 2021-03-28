package cxparser

import (
	"bytes"
	"io"
	"os"

	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/constants"
	globals2 "github.com/skycoin/cx/cx/globals"

	"github.com/skycoin/cx/cxgo/actions"
	"github.com/skycoin/cx/cxgo/cxgo"
	"github.com/skycoin/cx/cxgo/cxgo0"
	"github.com/skycoin/cx/cxgo/util/profiling"
)

// ParseSourceCode takes a group of files representing CX `sourceCode` and
// parses it into CX program structures for `AST`.
func ParseSourceCode(sourceCode []*os.File, fileNames []string) {

	//local
	cxgo0.PRGRM0 = actions.AST

	programsourceCode := make([]string, len(sourceCode))
	for i, source := range sourceCode {
		tmp := bytes.NewBuffer(nil)
		io.Copy(tmp, source)
		programsourceCode[i] = tmp.String()

	}

	/*
		We need to traverse the elements by hierarchy first add all the
		packages and structs at the same time then add globals, as these
		can be of a custom type (and it could be imported) the signatures
		of functions and methods are added in the cxgo0.y pass
	*/

	parseErrors := 0
	if len(sourceCode) > 0 {
		parseErrors = lexer(programsourceCode, fileNames)
	}

	//package level program
	actions.AST.SetCurrentCxProgram()

	actions.AST = cxgo0.PRGRM0

	if globals2.FoundCompileErrors || parseErrors > 0 {
		profiling.CleanupAndExit(constants.CX_COMPILATION_ERROR)
	}

	// Adding global variables `OS_ARGS` to the `os` (operating system)
	// package.
	if osPkg, err := actions.AST.GetPackage(constants.OS_PKG); err == nil {
		if _, err := osPkg.GetGlobal(constants.OS_ARGS); err != nil {
			arg0 := ast.MakeArgument(constants.OS_ARGS, "", -1).AddType(constants.TypeNames[constants.TYPE_UNDEFINED])
			arg0.Package = osPkg

			arg1 := ast.MakeArgument(constants.OS_ARGS, "", -1).AddType(constants.TypeNames[constants.TYPE_STR])
			arg1 = actions.DeclarationSpecifiers(arg1, []int{0}, constants.DECL_BASIC)
			arg1 = actions.DeclarationSpecifiers(arg1, []int{0}, constants.DECL_SLICE)

			actions.DeclareGlobalInPackage(osPkg, arg0, arg1, nil, false)
		}
	}

	profiling.StartProfile("4. parse")
	// The last pass of parsing that generates the actual output.
	for i, source := range programsourceCode {
		// Because of an unkown reason, sometimes some CX programs
		// throw an error related to a premature EOF (particularly in Windows).
		// Adding a newline character solves this.
		source = source + "\n"
		actions.LineNo = 1
		b := bytes.NewBufferString(source)
		if len(fileNames) > 0 {
			actions.CurrentFile = fileNames[i]
		}
		profiling.StartProfile(actions.CurrentFile)
		parseErrors += cxgo.Parse(cxgo.NewLexer(b))
		profiling.StopProfile(actions.CurrentFile)
	}
	profiling.StopProfile("4. parse")

	if globals2.FoundCompileErrors || parseErrors > 0 {
		profiling.CleanupAndExit(constants.CX_COMPILATION_ERROR)
	}
}

/* function lexer performs a first pass for the CX cxgo.
Globals, packages and custom types are parsed added to `cxgo0.PRGRM0`.
*/
func lexer(programSourceCode, programFileNames []string) int {

	var cxpackage *ast.CXPackage
	parseErrors := 0

	re := newRegulaEexpression()

	lexstruct(programSourceCode, programFileNames, cxpackage, re)

	lexpackages(programSourceCode, programFileNames, cxpackage, re)

	lexglobal(programSourceCode, programFileNames, cxpackage, re)

	profiling.StartProfile("1. packages/structs")

	profiling.StartProfile("3. cxgo0")

	for i, source := range programSourceCode {

		profiling.StartProfile(programFileNames[i])

		source = source + "\n"

		if len(programFileNames) > 0 {
			cxgo0.CurrentFileName = programFileNames[i]
		}

		parseErrors += cxgo0.Parse(source)

		profiling.StopProfile(programFileNames[i])
	}

	profiling.StopProfile("3. cxgo0")

	return parseErrors
}
