package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/skycoin/cx/cx/opcodes"
	"github.com/skycoin/cx/cx/types"

	repl "github.com/skycoin/cx/cmd/cxrepl"
	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/constants"
	"github.com/skycoin/cx/cx/globals"
	cxinit "github.com/skycoin/cx/cx/init"
	"github.com/skycoin/cx/cx/util"
	"github.com/skycoin/cx/cxparser/actions"
	cxparsering "github.com/skycoin/cx/cxparser/cxparsing"
	cxparsing "github.com/skycoin/cx/cxparser/cxparsing"
	cxgo "github.com/skycoin/cx/cxparser/cxparsingcompletor"
	"github.com/skycoin/cx/cxparser/util/profiling"
)

func parseProgram(options cxCmdFlags, fileNames []string, sourceCode []*os.File) bool {

	profile := profiling.StartCPUProfile("parse")
	defer profiling.StopCPUProfile(profile)

	defer profiling.DumpMEMProfile("parse")

	profiling.StartProfile("parse")
	defer profiling.StopProfile("parse")

	if actions.AST == nil {
		actions.AST = cxinit.MakeProgram()
	}

	// Parsing all the source code files sent as CLI arguments to CX.
	cxparsing.ParseSourceCode(sourceCode, fileNames)

	// Checking if a main package exists. If not, create and add it to `AST`.
	if _, err := actions.AST.GetFunction(constants.MAIN_FUNC, constants.MAIN_PKG); err != nil {
		initMainPkg(actions.AST)
	}

	// Setting what function to start in if using the REPL.
	repl.ReplTargetFn = constants.MAIN_FUNC

	// Adding *init function that initializes all the global variables.
	err := cxparsering.AddInitFunction(actions.AST)
	if err != nil {
		panic(fmt.Sprintf("error initializing function: %v", err))
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
	modIdx := prgrm.AddPackage(mod)
	modPkg, err := prgrm.GetPackageFromArray(modIdx)
	if err != nil {
		panic(err)
	}
	fn := ast.MakeFunction(constants.MAIN_FUNC, actions.CurrentFile, actions.LineNo)
	modPkg.AddFunction(prgrm, fn)
}

// optionTokenize checks if the user wants to use CX to generate the lexer tokens
func printTokenize(options cxCmdFlags, fileNames []string) {
	var r *os.File
	var w *os.File
	var err error
	var lines []string
	var sourceFilename string

	if len(fileNames) == 0 {
		r = os.Stdin
	} else {
		sourceFilename = fileNames[0]
		if len(fileNames) > 1 {
			fmt.Fprintln(os.Stderr, "Multiple source files detected. Ignoring all except", sourceFilename)
		}
		r, err = util.CXOpenFile(sourceFilename)
		if err != nil {
			fmt.Fprintln(os.Stderr, "ProgramError reading:", sourceFilename, err)
			return
		}
		defer r.Close()

		bytesRead, err := ioutil.ReadFile(sourceFilename)
		if err != nil {
			fmt.Fprintln(os.Stderr, "ProgramError reading:", sourceFilename, err)
			return
		}

		file_content := string(bytesRead)
		lines = strings.Split(file_content, "\n")
	}

	if options.compileOutput == "" {
		w = os.Stdout
	} else {
		tokenFilename := options.compileOutput
		w, err = util.CXCreateFile(tokenFilename)
		if err != nil {
			fmt.Fprintln(os.Stderr, "ProgramError writing:", tokenFilename, err)
			return
		}
		defer w.Close()
	}

	cxgo.Tokenize(r, w, lines, sourceFilename)
}

func printProgramAST(options cxCmdFlags, cxArgs []string, sourceCode []*os.File) {
	profiling.StartProfile("run")
	defer profiling.StopProfile("run")

	if options.replMode || len(sourceCode) == 0 {
		repl.Repl()
		return
	}

	// Print CX program.
	actions.AST.PrintProgram()

	if opcodes.AssertFailed() {
		os.Exit(constants.CX_ASSERT)
	}
}

// Used for the -heap-initial, -heap-max and -stack-size flags.
// This function parses, for example, "1M" to 1048576 (the corresponding number of bytes)
// Possible suffixes are: G or g (gigabytes), M or m (megabytes), K or k (kilobytes)
func parseMemoryString(s string) types.Pointer {
	suffix := s[len(s)-1]
	_, notSuffix := strconv.ParseFloat(string(suffix), 64)

	if notSuffix == nil {
		// then we don't have a suffix
		num, err := strconv.ParseInt(s, 10, 64)

		if err != nil {
			// malformed size
			return types.InvalidPointer
		}

		return types.Cast_i64_to_ptr(num)
	} else {
		// then we have a suffix
		num, err := strconv.ParseFloat(s[:len(s)-1], 64)

		if err != nil {
			// malformed size
			return types.InvalidPointer
		}

		// The user can use suffixes to give as input gigabytes, megabytes or kilobytes.
		switch suffix {
		case 'G', 'g':
			return types.Cast_f64_to_ptr(num * 1073741824)
		case 'M', 'm':
			return types.Cast_f64_to_ptr(num * 1048576)
		case 'K', 'k':
			return types.Cast_f64_to_ptr(num * 1024)
		default:
			return types.InvalidPointer
		}
	}
}

type SourceCode struct {
	Code string //Unused?
}

// GetCXPath checks if the user has set the environment variable
// `CXPATH`. If not, CX creates a workspace at $HOME/cx, along with $HOME/cx/bin,
// $HOME/cx/pkg and $HOME/cx/src
func GetCXPath(options cxCmdFlags) {
	// Determining the filepath of the directory where the user
	// started the `cx` command.
	_, err := os.Executable()
	if err != nil {
		panic(err)
	}

	CXPATH := ""
	if os.Getenv("CXPATH") != "" {
		CXPATH = os.Getenv("CXPATH")
	}
	// `options.cxpath` overrides `os.Getenv("CXPATH")`
	if options.cxpath != "" {
		CXPATH, err = filepath.Abs(options.cxpath)
		if err != nil {
			panic(err)
		}
	}
	if os.Getenv("CXPATH") == "" && options.cxpath == "" {
		usr, err := user.Current()
		if err != nil {
			panic(err)
		}

		CXPATH = usr.HomeDir + "/cx/"
	}
	globals.BINPATH = filepath.Join(CXPATH, "bin/")
	globals.PKGPATH = filepath.Join(CXPATH, "pkg/")
	globals.SRCPATH = filepath.Join(CXPATH, "src/")
	//why would we create directories on executing every CX program?
	//directory creation should be on installation
	//CreateCxDirectories(CXPATH)
}

/*
func CreateCxDirectories(CXPATH string) {
		// Creating directories in case they do not exist.
		if _, err := cxcore.CXStatFile(CXPATH); os.IsNotExist(err) {
			cxcore.CXMkdirAll(CXPATH, 0755)
		}
		if _, err := cxcore.CXStatFile(cxcore.BINPATH); os.IsNotExist(err) {
			cxcore.CXMkdirAll(cxcore.BINPATH, 0755)
		}
		if _, err := cxcore.CXStatFile(cxcore.PKGPATH); os.IsNotExist(err) {
			cxcore.CXMkdirAll(cxcore.PKGPATH, 0755)
		}
		if _, err := cxcore.CXStatFile(cxcore.SRCPATH); os.IsNotExist(err) {
			cxcore.CXMkdirAll(cxcore.SRCPATH, 0755)
		}
	}
*/
