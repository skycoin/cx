package main

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	cxcore "github.com/skycoin/cx/cx"
	"github.com/skycoin/cx/cxgo/actions"
	"github.com/skycoin/cx/cxgo/cxgo"
	"github.com/skycoin/cx/cxgo/cxgo0"
	"github.com/skycoin/cx/cxgo/parser"
)

const VERSION = "0.7.1"

func main() {
	//cx.CXLogFile(true)
	if os.Args != nil && len(os.Args) > 1 {
		Run(os.Args[1:])
	}
}

func Run(args []string) {
	runtime.LockOSThread()
	runtime.GOMAXPROCS(2)

	options := defaultCmdFlags()
	parseFlags(&options, args)

	// Checking if CXPATH is set, either by setting an environment variable
	// or by setting the `--cxpath` flag.
	checkCXPathSet(options)

	if checkHelp(args) {
		commandLine.PrintDefaults()
		return
	}

	// Does the user want to print the command-line help?
	if options.printHelp {
		printHelp()
		return
	}

	// Does the user want to print CX's version?
	if options.printVersion {
		printVersion()
		return
	}

	// User wants to print CX env
	if options.printEnv {
		printEnv()
		return
	}

	if options.initialHeap != "" {
		cxcore.INIT_HEAP_SIZE = parseMemoryString(options.initialHeap)
	}
	if options.maxHeap != "" {
		cxcore.MAX_HEAP_SIZE = parseMemoryString(options.maxHeap)
		if cxcore.MAX_HEAP_SIZE < cxcore.INIT_HEAP_SIZE {
			// Then MAX_HEAP_SIZE overrides INIT_HEAP_SIZE's value.
			cxcore.INIT_HEAP_SIZE = cxcore.MAX_HEAP_SIZE
		}
	}
	if options.stackSize != "" {
		cxcore.STACK_SIZE = parseMemoryString(options.stackSize)
		actions.DataOffset = cxcore.STACK_SIZE
	}
	if options.minHeapFreeRatio != float64(0) {
		cxcore.MIN_HEAP_FREE_RATIO = float32(options.minHeapFreeRatio)
	}
	if options.maxHeapFreeRatio != float64(0) {
		cxcore.MAX_HEAP_FREE_RATIO = float32(options.maxHeapFreeRatio)
	}

	// options, file pointers, filenames
	cxArgs, sourceCode, fileNames := cxcore.ParseArgsForCX(commandLine.Args(), true)

	// Propagate some options out to other packages.
	parser.DebugLexer = options.debugLexer // in package parser
	DebugProfileRate = options.debugProfile
	DebugProfile = DebugProfileRate > 0

	if run, bcHeap, sPrgrm := parseProgram(options, fileNames, sourceCode); run {
		runProgram(options, cxArgs, sourceCode, bcHeap, sPrgrm)
	}
}

// initMainPkg adds a `main` package with an empty `main` function to `prgrm`.
func initMainPkg(prgrm *cxcore.CXProgram) {
	mod := cxcore.MakePackage(cxcore.MAIN_PKG)
	prgrm.AddPackage(mod)
	fn := cxcore.MakeFunction(cxcore.MAIN_FUNC, actions.CurrentFile, actions.LineNo)
	mod.AddFunction(fn)
}

// optionTokenize checks if the user wants to use CX to generate the lexer tokens
func optionTokenize(options cxCmdFlags, fileNames []string) {
	var r *os.File
	var w *os.File
	var err error

	if len(fileNames) == 0 {
		r = os.Stdin
	} else {
		sourceFilename := fileNames[0]
		if len(fileNames) > 1 {
			fmt.Fprintln(os.Stderr, "Multiple source files detected. Ignoring all except", sourceFilename)
		}
		r, err = cxcore.CXOpenFile(sourceFilename)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading:", sourceFilename, err)
			return
		}
		defer r.Close()
	}

	if options.compileOutput == "" {
		w = os.Stdout
	} else {
		tokenFilename := options.compileOutput
		w, err = cxcore.CXCreateFile(tokenFilename)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error writing:", tokenFilename, err)
			return
		}
		defer w.Close()
	}

	parser.Tokenize(r, w)
}

func parseProgram(options cxCmdFlags, fileNames []string, sourceCode []*os.File) (bool, []byte, []byte) {
	profile := StartCPUProfile("parse")
	defer StopCPUProfile(profile)

	defer DumpMEMProfile("parse")

	StartProfile("parse")
	defer StopProfile("parse")

	actions.PRGRM = cxcore.MakeProgram()
	corePkgsPrgrm, err := cxcore.GetProgram()
	if err != nil {
		panic(err)
	}
	actions.PRGRM.Packages = corePkgsPrgrm.Packages

	// var bcPrgrm *CXProgram
	var sPrgrm []byte
	// In case of a CX chain, we need to temporarily store the blockchain code heap elsewhere,
	// so we can then add it after the transaction code's data segment.
	var bcHeap []byte

	// Parsing all the source code files sent as CLI arguments to CX.
	cxgo.ParseSourceCode(sourceCode, fileNames)

	// setting project's working directory
	if !options.replMode && len(sourceCode) > 0 {
		cxgo0.PRGRM0.Path = determineWorkDir(sourceCode[0].Name())
	}

	// Checking if a main package exists. If not, create and add it to `PRGRM`.
	if _, err := actions.PRGRM.GetFunction(cxcore.MAIN_FUNC, cxcore.MAIN_PKG); err != nil {
		initMainPkg(actions.PRGRM)
	}
	// Setting what function to start in if using the REPL.
	actions.ReplTargetFn = cxcore.MAIN_FUNC

	// Adding *init function that initializes all the global variables.
	cxgo.AddInitFunction(actions.PRGRM)

	actions.LineNo = 0

	if cxcore.FoundCompileErrors {
		//cleanupAndExit(cxcore.CX_COMPILATION_ERROR)
		StopCPUProfile(profile)
		exitCode := cxcore.CX_COMPILATION_ERROR
		os.Exit(exitCode)

	}

	return true, bcHeap, sPrgrm
}

func runProgram(options cxCmdFlags, cxArgs []string, sourceCode []*os.File, bcHeap []byte, sPrgrm []byte) {
	StartProfile("run")
	defer StopProfile("run")

	if options.replMode || len(sourceCode) == 0 {
		actions.PRGRM.SelectProgram()
		Repl()
		return
	}

	// Normal run of a CX program.
	err := actions.PRGRM.RunCompiled(0, cxArgs)
	if err != nil {
		panic(err)
	}

	if cxcore.AssertFailed() {
		os.Exit(cxcore.CX_ASSERT)
	}
}

// Used for the -heap-initial, -heap-max and -stack-size flags.
// This function parses, for example, "1M" to 1048576 (the corresponding number of bytes)
// Possible suffixes are: G or g (gigabytes), M or m (megabytes), K or k (kilobytes)
func parseMemoryString(s string) int {
	suffix := s[len(s)-1]
	_, notSuffix := strconv.ParseFloat(string(suffix), 64)

	if notSuffix == nil {
		// then we don't have a suffix
		num, err := strconv.ParseInt(s, 10, 64)

		if err != nil {
			// malformed size
			return -1
		}

		return int(num)
	} else {
		// then we have a suffix
		num, err := strconv.ParseFloat(s[:len(s)-1], 64)

		if err != nil {
			// malformed size
			return -1
		}

		// The user can use suffixes to give as input gigabytes, megabytes or kilobytes.
		switch suffix {
		case 'G', 'g':
			return int(num * 1073741824)
		case 'M', 'm':
			return int(num * 1048576)
		case 'K', 'k':
			return int(num * 1024)
		default:
			return -1
		}
	}
}

type SourceCode struct {
	Code string //Unused?
}

func determineWorkDir(filename string) string {
	filename = filepath.FromSlash(filename)

	i := strings.LastIndexByte(filename, os.PathSeparator)
	if i == -1 {
		i = 0
	}
	return filename[:i]
}

// checkCXPathSet checks if the user has set the environment variable
// `CXPATH`. If not, CX creates a workspace at $HOME/cx, along with $HOME/cx/bin,
// $HOME/cx/pkg and $HOME/cx/src
func checkCXPathSet(options cxCmdFlags) {
	// Determining the filepath of the directory where the user
	// started the `cx` command.
	_, err := os.Executable()
	if err != nil {
		panic(err)
	}
	// cxcore.COREPATH = filepath.Dir(ex) // TODO @evanlinjin: Not used.

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

	cxcore.BINPATH = filepath.Join(CXPATH, "bin/")
	cxcore.PKGPATH = filepath.Join(CXPATH, "pkg/")
	cxcore.SRCPATH = filepath.Join(CXPATH, "src/")

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
