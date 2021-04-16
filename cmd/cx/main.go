package main

import (
	"os"
	"runtime"

	"github.com/skycoin/cx/cx/opcodes"

	repl "github.com/skycoin/cx/cmd/cxrepl"
	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/constants"
	"github.com/skycoin/cx/cx/execute"
	"github.com/skycoin/cx/cxparser/actions"
	parsingcompletor "github.com/skycoin/cx/cxparser/cxparsingcompletor"
	"github.com/skycoin/cx/cxparser/util/profiling"
)

func main() {

	if os.Args != nil && len(os.Args) > 1 {
		Run(os.Args[1:])
	}
}

func Run(args []string) {

	runtime.LockOSThread()

	runtime.GOMAXPROCS(2)

	options := defaultCmdFlags()

	parseFlags(&options, args)

	/*
		Checking if CXPATH is set, either by setting an environment variable
		 or by setting the `--cxpath` flag.
	*/

	GetCXPath(options)

	/*
		checkHelp checks  for command line argument string "help"
		if exits print help
		$cx help
	*/

	if checkHelp(args) {
		commandLine.PrintDefaults()
		return
	}

	/*
		options.printHelp checks for flags string "help"
		if exits print help
		$cx --help
		$cx -help
	*/
	if options.printHelp {
		printHelp()
		return
	}

	/*
		options.printVersion checks for flags string "version"
		if exits print cx version
		$cx --version
		$cx -v
	*/

	if options.printVersion {
		printVersion()
		return
	}

	/*
		checkversion checks  for command line argument string "version"
		if exits print version
		$cx version
	*/
	if checkversion(args) {
		printVersion()
		return
	}

	/*
		options.printEnv checks for flags string "env"
		if exits print cx env
		$cx --env
	*/

	if options.printEnv {
		printEnv()
		return
	}

	/*
		checkenv checks  for command line argument string "env"
		if exits print env
		$cx env
	*/
	if checkenv(args) {
		printEnv()
		return
	}

	/*
		options.initialHeap checks for flags string "heap-initial"
		$cx --heap-initial
	*/
	if options.initialHeap != "" {
		constants.INIT_HEAP_SIZE = parseMemoryString(options.initialHeap)
	}

	/*
		options.maxHeap checks for flags string "maxHeap"
		$cx --maxHeap
	*/

	if options.maxHeap != "" {
		constants.MAX_HEAP_SIZE = parseMemoryString(options.maxHeap)
		if constants.MAX_HEAP_SIZE < constants.INIT_HEAP_SIZE {
			// Then MAX_HEAP_SIZE overrides INIT_HEAP_SIZE's value.
			constants.INIT_HEAP_SIZE = constants.MAX_HEAP_SIZE
		}
	}

	/*
		options.minHeapFreeRatio checks for flags string "--min-heap-free"
		$cx --min-heap-free
	*/

	if options.minHeapFreeRatio != float64(0) {
		constants.MIN_HEAP_FREE_RATIO = float32(options.minHeapFreeRatio)
	}

	/*
		options.maxHeapFreeRatio checks for flags string "--max-heap-free"
		$cx --max-heap-free
	*/
	if options.maxHeapFreeRatio != float64(0) {
		constants.MAX_HEAP_FREE_RATIO = float32(options.maxHeapFreeRatio)
	}

	// options, file pointers, filenames
	cxArgs, sourceCode, fileNames := ast.ParseArgsForCX(commandLine.Args(), true)

	// Propagate some options out to other packages.
	parsingcompletor.DebugLexer = options.debugLexer

	profiling.DebugProfileRate = options.debugProfile

	profiling.DebugProfile = profiling.DebugProfileRate > 0

	// Load op code tables
	parsingcompletor.InitCXCore()

	if run := parseProgram(options, fileNames, sourceCode); run {

		if checkAST(args) {
			printProgramAST(options, cxArgs, sourceCode)
			return
		}

		if options.tokenizeMode {
			printTokenize(options, fileNames)
			return
		}

		//if strings.Contains(args[0], "tokenize") {
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

	err := execute.RunCompiled(actions.AST, 0, cxArgs)

	if err != nil {
		panic(err)
	}

	if opcodes.AssertFailed() {
		os.Exit(constants.CX_ASSERT)
	}
}
