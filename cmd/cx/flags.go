package main

import (
	"flag"
	"fmt"
	"go/build"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type cxCmdFlags struct {
	baseOutput       bool
	compileOutput    string
	replMode         bool
	printHelp        bool
	printVersion     bool
	printEnv         bool
	printAST         bool
	tokenizeMode     bool
	initialHeap      string
	maxHeap          string
	stackSize        string
	minHeapFreeRatio float64
	maxHeapFreeRatio float64
	cxpath           string

	// Debug flags for the CX developers
	debugLexer   bool
	debugProfile int
}

func defaultCmdFlags() cxCmdFlags {
	return cxCmdFlags{
		baseOutput:    false,
		compileOutput: "",
		replMode:      false,
		printHelp:     false,
		printEnv:      false,
		printVersion:  false,
		debugLexer:    false,
		debugProfile:  0,
	}
}

var commandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

func appendDash(args []string) {

	for k, v := range args {
		switch v {
		case "version":
			args[k] = "-version"
		case "ast":
			args[k] = "-ast"
		}
	}
}

func parseFlags(options *cxCmdFlags, args []string) {
	if len(args) <= 0 {
		options.replMode = true
	}

	appendDash(args)

	commandLine.BoolVar(&options.printVersion, "version", options.printVersion, "Print CX version")
	commandLine.BoolVar(&options.printVersion, "v", options.printVersion, "alias for -version")
	commandLine.BoolVar(&options.printEnv, "env", options.printEnv, "Print CX environment information")
	commandLine.BoolVar(&options.printAST, "ast", options.printAST, "Print CX Program AST")
	commandLine.BoolVar(&options.tokenizeMode, "tokenize", options.tokenizeMode, "generate a 'out.cx.txt' text file with parsed tokens")
	commandLine.BoolVar(&options.tokenizeMode, "t", options.tokenizeMode, "alias for -tokenize")
	commandLine.StringVar(&options.compileOutput, "co", options.compileOutput, "alias for -compile-output")

	commandLine.BoolVar(&options.replMode, "repl", options.replMode, "Loads source files into memory and starts a read-eval-print loop")
	commandLine.BoolVar(&options.replMode, "r", options.replMode, "alias for -repl")
	commandLine.StringVar(&options.initialHeap, "heap-initial", options.initialHeap, "Set the initial heap for the CX virtual machine. The value is in bytes, but the suffixes 'G', 'M' or 'K' can be used to express gigabytes, megabytes or kilobytes, respectively. Lowercase suffixes are allowed.")
	commandLine.StringVar(&options.initialHeap, "hi", options.initialHeap, "alias for -initial-heap")
	commandLine.StringVar(&options.maxHeap, "heap-max", options.maxHeap, "Set the max heap for the CX virtual machine. The value is in bytes, but the suffixes 'G', 'M' or 'K' can be used to express gigabytes, megabytes or kilobytes, respectively. Lowercase suffixes are allowed. Note that this parameter overrides --heap-initial if --heap-max is equal to a lesser value than --heap-max's.")
	commandLine.StringVar(&options.maxHeap, "hm", options.maxHeap, "alias for -max-heap")
	commandLine.StringVar(&options.stackSize, "stack-size", options.stackSize, "Set the stack size for the CX virtual machine. The value is in bytes, but the suffixes 'G', 'M' or 'K' can be used to express gigabytes, megabytes or kilobytes, respectively. Lowercase suffixes are allowed.")
	commandLine.StringVar(&options.stackSize, "ss", options.stackSize, "alias for -stack-size")
	commandLine.Float64Var(&options.minHeapFreeRatio, "--min-heap-free", options.minHeapFreeRatio, "Minimum heap space percentage that should be free after calling the garbage collector. Value must be in the range of 0.0 and 1.0.")
	commandLine.Float64Var(&options.maxHeapFreeRatio, "--max-heap-free", options.maxHeapFreeRatio, "Maximum heap space percentage that should be free after calling the garbage collector. Value must be in the range of 0.0 and 1.0.")
	commandLine.StringVar(&options.cxpath, "cxpath", options.cxpath, "Used for dynamically setting the value of the environment variable CXPATH")

	// Debug flags
	commandLine.BoolVar(&options.debugLexer, "debug-lexer", options.debugLexer, "Debug the lexer by printing all scanner tokens")
	commandLine.BoolVar(&options.debugLexer, "Dl", options.debugLexer, "alias for -debug-lexer")
	commandLine.IntVar(&options.debugProfile, "debug-profile", options.debugProfile, "Enable CPU+MEM profiling and set CPU profiling rate. Visualize .pprof files with \"go get github.com/google/pprof\" and \"pprof -http=:8080 file.pprof\"")
	commandLine.IntVar(&options.debugProfile, "Dp", options.debugProfile, "alias for -debug-profile")

	commandLine.Parse(args)

}

func printHelp() {
	fmt.Printf(`Usage: cx [options] [source-files]

CX options:
-h, --help                        Prints this message.
-n, --new                         Creates a new project located at $CXPATH/src
-r, --repl                        Loads source files into memory and starts a read-eval-print loop.
-w, --web                         Start CX as a web service.

Notes:
* Option --web makes every other flag to be ignored.
`)
}

func printVersion() {
	fmt.Printf("CX version %v %v/%v\n", VERSION, runtime.GOOS, runtime.GOARCH)
}

func checkHelp(args []string) bool {
	if strings.Contains(args[0], "help") {
		return true
	}
	return false
}

func checkversion(args []string) bool {
	if strings.Contains(args[0], "version") {
		return true
	}
	return false
}

func checkenv(args []string) bool {
	if strings.Contains(args[0], "env") {
		return true
	}
	return false
}

func checkAST(args []string) bool {
	if strings.Contains(args[0], "ast") {
		return true
	}
	return false
}

func checktokenizeMode(args []string) bool {
	if strings.Contains(args[0], "tokenize") {
		return true
	}
	return false
}

func printEnv() {
	ex, _ := os.Executable()

	fmt.Println("GOROOT: ", runtime.GOROOT())
	fmt.Println("GOPATH: ", build.Default.GOPATH)
	fmt.Println("GOBIN: ", os.Getenv("GOBIN"))
	fmt.Println("GO version: ", runtime.Version())
	fmt.Println("Operating system: ", runtime.GOOS)
	fmt.Println("CX version: ", VERSION)
	fmt.Println("CX binary location: ", filepath.Dir(ex))
}
