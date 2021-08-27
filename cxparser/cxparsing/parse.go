package cxparsing

import (
	"os"

	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/constants"
	"github.com/skycoin/cx/cx/globals"
	"github.com/skycoin/cx/cxparser/actions"
	"github.com/skycoin/cx/cxparser/util/profiling"
)

func ParseProgram(fileNames []string, sourceCode []*os.File) bool {

	profile := profiling.StartCPUProfile("parse")
	defer profiling.StopCPUProfile(profile)

	defer profiling.DumpMEMProfile("parse")

	profiling.StartProfile("parse")
	defer profiling.StopProfile("parse")

	actions.AST = ast.MakeProgram()

	//corePkgsPrgrm, err := cxcore.GetCurrentCxProgram()
	var corePkgsPrgrm *ast.CXProgram = ast.PROGRAM

	if corePkgsPrgrm == nil {
		panic("CxProgram is nil")
	}
	actions.AST.Packages = corePkgsPrgrm.Packages

	// var bcPrgrm *CXProgram
	//var sPrgrm []byte
	// In case of a CX chain, we need to temporarily store the blockchain code heap elsewhere,
	// so we can then add it after the transaction code's data segment.
	//var bcHeap []byte

	// Parsing all the source code files sent as CLI arguments to CX.
	// TODO: comment what this function does
	ParseSourceCode(sourceCode, fileNames)

	//remove path variable, not used
	// setting project's working directory
	//if !options.replMode && len(sourceCode) > 0 {
	//cxgo0.PRGRM0.Path = determineWorkDir(sourceCode[0].Name())
	//}

	//globals.CxProgramPath = determineWorkDir(sourceCode[0].Name())
	//globals2.SetWorkingDir(sourceCode[0].Name())

	// Checking if a main package exists. If not, create and add it to `AST`.
	if _, err := actions.AST.GetFunction(constants.MAIN_FUNC, constants.MAIN_PKG); err != nil {
		panic("error")
	}
	initMainPkg(actions.AST)

	// Adding *init function that initializes all the global variables.
	err := AddInitFunction(actions.AST)
	if err != nil {
		return false //why return false, instead of panicing
	}

	actions.LineNo = 0

	if globals.FoundCompileErrors {
		//cleanupAndExit(cxcore.CX_COMPILATION_ERROR)
		profiling.StopCPUProfile(profile)
		exitCode := constants.CX_COMPILATION_ERROR
		os.Exit(exitCode)

	}

	return true
}

// initMainPkg adds a `main` package with an empty `main` function to `prgrm`.
func initMainPkg(prgrm *ast.CXProgram) {
	mod := ast.MakePackage(constants.MAIN_PKG)
	prgrm.AddPackage(mod)
	fn := ast.MakeFunction(constants.MAIN_FUNC, actions.CurrentFile, actions.LineNo)
	mod.AddFunction(fn)
}
