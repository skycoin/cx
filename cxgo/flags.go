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
	baseOutput        bool
	compileOutput     string
	replMode          bool
	webMode           bool
	ideMode           bool
	webPersistentMode bool
	printHelp         bool
	printVersion      bool
	printEnv          bool
	tokenizeMode      bool
	initialHeap       string
	maxHeap           string
	stackSize         string
	blockchainMode    bool
	publisherMode     bool
	peerMode          bool
	transactionMode   bool
	broadcastMode     bool
	walletMode        bool
	genAddress        bool
	port              int
	walletId          string
	walletSeed        string
	programName       string
	secKey            string
	pubKey            string
	genesisAddress    string
	genesisSignature  string
	minHeapFreeRatio  float64
	maxHeapFreeRatio  float64
	cxpath            string

	// Debug flags for the CX developers
	debugLexer   bool
	debugProfile int
}

func defaultCmdFlags() cxCmdFlags {
	return cxCmdFlags{
		baseOutput:        false,
		compileOutput:     "",
		replMode:          false,
		webMode:           false,
		ideMode:           false,
		webPersistentMode: false,
		printHelp:         false,
		printEnv:          false,
		printVersion:      false,
		blockchainMode:    false,
		transactionMode:   false,
		broadcastMode:     false,
		port:              6001,
		programName:       "cxcoin",
		walletId:          "cxcoin_cli.wlt",
		secKey:            "",
		pubKey:            "",
		genesisAddress:    "",
		genesisSignature:  "",

		debugLexer:   false,
		debugProfile: 0,
	}
}

var commandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

func parseFlags(options *cxCmdFlags, args []string) {
	if len(args) <= 0 {
		options.replMode = true
	}

	commandLine.BoolVar(&options.printVersion, "version", options.printVersion, "Print CX version")
	commandLine.BoolVar(&options.printVersion, "v", options.printVersion, "alias for -version")
	commandLine.BoolVar(&options.printEnv, "env", options.printEnv, "Print CX environment information")
	commandLine.BoolVar(&options.tokenizeMode, "tokenize", options.tokenizeMode, "generate a 'out.cx.txt' text file with parsed tokens")
	commandLine.BoolVar(&options.tokenizeMode, "t", options.tokenizeMode, "alias for -tokenize")
	commandLine.StringVar(&options.compileOutput, "co", options.compileOutput, "alias for -compile-output")

	commandLine.BoolVar(&options.replMode, "repl", options.replMode, "Loads source files into memory and starts a read-eval-print loop")
	commandLine.BoolVar(&options.replMode, "r", options.replMode, "alias for -repl")
	commandLine.BoolVar(&options.webMode, "web", options.webMode, "Start CX as a web service.")
	commandLine.BoolVar(&options.webMode, "w", options.webMode, "alias for -web")
	commandLine.BoolVar(&options.ideMode, "ide", options.ideMode, "Start CX as a web service, and Leaps service start also.")
	commandLine.BoolVar(&options.webPersistentMode, "pw", options.webPersistentMode, "Start CX as a web service with a persistent web REPL session")
	commandLine.StringVar(&options.initialHeap, "heap-initial", options.initialHeap, "Set the initial heap for the CX virtual machine. The value is in bytes, but the suffixes 'G', 'M' or 'K' can be used to express gigabytes, megabytes or kilobytes, respectively. Lowercase suffixes are allowed.")
	commandLine.StringVar(&options.initialHeap, "hi", options.initialHeap, "alias for -initial-heap")
	commandLine.StringVar(&options.maxHeap, "heap-max", options.maxHeap, "Set the max heap for the CX virtual machine. The value is in bytes, but the suffixes 'G', 'M' or 'K' can be used to express gigabytes, megabytes or kilobytes, respectively. Lowercase suffixes are allowed. Note that this parameter overrides --heap-initial if --heap-max is equal to a lesser value than --heap-max's.")
	commandLine.StringVar(&options.maxHeap, "hm", options.maxHeap, "alias for -max-heap")
	commandLine.StringVar(&options.stackSize, "stack-size", options.stackSize, "Set the stack size for the CX virtual machine. The value is in bytes, but the suffixes 'G', 'M' or 'K' can be used to express gigabytes, megabytes or kilobytes, respectively. Lowercase suffixes are allowed.")
	commandLine.StringVar(&options.stackSize, "ss", options.stackSize, "alias for -stack-size")
	commandLine.Float64Var(&options.minHeapFreeRatio, "--min-heap-free", options.minHeapFreeRatio, "Minimum heap space percentage that should be free after calling the garbage collector. Value must be in the range of 0.0 and 1.0.")
	commandLine.Float64Var(&options.maxHeapFreeRatio, "--max-heap-free", options.maxHeapFreeRatio, "Maximum heap space percentage that should be free after calling the garbage collector. Value must be in the range of 0.0 and 1.0.")

	// commandLine.BoolVar(&options.blockchainMode, "bc", options.blockchainMode, "alias for -blockchain")
	// commandLine.BoolVar(&options.publisherMode, "pb", options.publisherMode, "alias for -publisher")
	// commandLine.BoolVar(&options.transactionMode, "txn", options.transactionMode, "alias for -transaction")
	commandLine.BoolVar(&options.broadcastMode, "broadcast", options.broadcastMode, "Broadcast a CX blockchain transaction")
	commandLine.BoolVar(&options.walletMode, "create-wallet", options.walletMode, "Create a wallet from a seed")
	commandLine.StringVar(&options.cxpath, "cxpath", options.cxpath, "Used for dynamically setting the value of the environment variable CXPATH")

	//deprecated

	//commandLine.BoolVar(&options.blockchainMode, "blockchain", options.blockchainMode, "Start a CX blockchain program")
	//commandLine.BoolVar(&options.genAddress, "generate-address", options.genAddress, "Generate a CX chain address")
	//commandLine.StringVar(&options.genesisAddress, "genesis-address", options.genesisAddress, "CX blockchain program genesis address")
	//commandLine.StringVar(&options.genesisSignature, "genesis-signature", options.genesisSignature, "CX blockchain program genesis address")
	//commandLine.BoolVar(&options.peerMode, "peer", options.peerMode, "Run a CX chain peer node")
	//commandLine.IntVar(&options.port, "port", options.port, "Port used when running a CX chain peer node")
	//commandLine.StringVar(&options.programName, "program-name", options.programName, "Name of the initial CX program on the blockchain")
	//commandLine.StringVar(&options.pubKey, "public-key", options.pubKey, "CX program blockchain public key")
	//commandLine.BoolVar(&options.publisherMode, "publisher", options.publisherMode, "Start a CX blockchain program block publisher")
	//commandLine.StringVar(&options.secKey, "secret-key", options.secKey, "CX program blockchain security key")
	//commandLine.BoolVar(&options.transactionMode, "transaction", options.transactionMode, "Test a CX blockchain transaction")
	//commandLine.StringVar(&options.walletSeed, "wallet-seed", options.walletSeed, "Seed to use for a new wallet")
	//commandLine.StringVar(&options.walletId, "wallet-id", options.walletId, "Wallet ID to use for signing transactions")

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
	fmt.Println("CX version", VERSION)
}

func checkHelp(args []string) bool {
	if strings.Contains(args[0], "help") {
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
