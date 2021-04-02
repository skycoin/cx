package main

import (
	repl "github.com/skycoin/cx/cmd/cxrepl"
	cxcore "github.com/skycoin/cx/cx"
	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/constants"
	"github.com/skycoin/cx/cx/execute"
	"github.com/skycoin/cx/cxgo/actions"
	"github.com/skycoin/cx/cxgo/cxgo"
	"github.com/skycoin/cx/cxgo/util/profiling"

	//"github.com/skycoin/cx/cxparser/cxgo0"
	"os"
	"runtime"
)

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
	GetCXPath(options)

	//checkHelp check command line argumenets
	//$ cx help
	if checkHelp(args) {
		commandLine.PrintDefaults()
		return
	}

	// Does the user want to print the command-line help?
	//options.printHelp works when flags are provided.
	//$ cx --vesion
	if options.printHelp {
		printHelp()
		return
	}

	// Does the user want to print CX's version?
	if options.printVersion {
		printVersion()
		return
	}

	//checkversion check command line argumenets
	//$ cx version
	if checkversion(args) {
		printVersion()
		return
	}

	// User wants to print CX env
	if options.printEnv {
		printEnv()
		return
	}

	//checkenv check command line argumenets
	//$ cx
	if checkenv(args) {
		printEnv()
		return
	}

	if options.initialHeap != "" {
		constants.INIT_HEAP_SIZE = parseMemoryString(options.initialHeap)
	}
	if options.maxHeap != "" {
		constants.MAX_HEAP_SIZE = parseMemoryString(options.maxHeap)
		if constants.MAX_HEAP_SIZE < constants.INIT_HEAP_SIZE {
			// Then MAX_HEAP_SIZE overrides INIT_HEAP_SIZE's value.
			constants.INIT_HEAP_SIZE = constants.MAX_HEAP_SIZE
		}
	}

	if options.minHeapFreeRatio != float64(0) {
		constants.MIN_HEAP_FREE_RATIO = float32(options.minHeapFreeRatio)
	}
	if options.maxHeapFreeRatio != float64(0) {
		constants.MAX_HEAP_FREE_RATIO = float32(options.maxHeapFreeRatio)
	}

	// options, file pointers, filenames
	cxArgs, sourceCode, fileNames := ast.ParseArgsForCX(commandLine.Args(), true)

	// Propagate some options out to other packages.
	cxgo.DebugLexer = options.debugLexer // in package cxgo
	profiling.DebugProfileRate = options.debugProfile
	profiling.DebugProfile = profiling.DebugProfileRate > 0

	// Load op code tables
	cxcore.LoadOpCodeTables()

	if run := parseProgram(options, fileNames, sourceCode); run {

		if checkAST(args) {
			printProgramAST(options, cxArgs, sourceCode)
			return
		}

		if options.tokenizeMode {
			printTokenize(options, fileNames)
			return
		}

		if checktokenizeMode(args) {
			printTokenize(options, fileNames)
			return
		}

		runProgram(options, cxArgs, sourceCode)
	}
}

func runProgram(options cxCmdFlags, cxArgs []string, sourceCode []*os.File) {
	profiling.StartProfile("run")
	defer profiling.StopProfile("run")

	if options.replMode || len(sourceCode) == 0 {
		actions.AST.SetCurrentCxProgram()
		repl.Repl()
		return
	}

	// Normal run of a CX program.
	//err := actions.AST.RunCompiled(0, cxArgs)
	err := execute.RunCompiled(actions.AST, 0, cxArgs)
	if err != nil {
		panic(err)
	}

	if cxcore.AssertFailed() {
		os.Exit(constants.CX_ASSERT)
	}
}
