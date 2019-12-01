package main

import (
	"flag"
	"fmt"
	"os"
)

type cxCmdFlags struct {
	baseOutput          bool
	compileOutput       string
	newProject          bool
	replMode            bool
	webMode             bool
	ideMode             bool
	webPersistentMode   bool
	printHelp           bool
	printVersion        bool
	tokenizeMode        bool
	initialHeap         string
	maxHeap             string
	stackSize           string
	blockchainMode      bool
	publisherMode       bool
	peerMode            bool
	transactionMode     bool
	broadcastMode       bool
	walletMode          bool
	genAddress          bool
	port                int
	walletId            string
	walletSeed          string
	programName         string
	secKey              string
	pubKey              string
	genesisAddress      string
	genesisSignature    string
	minHeapFreeRatio    float64
	maxHeapFreeRatio    float64
	cxpath              string

	// Debug flags for the CX developers
	debugLexer          bool
}

func defaultCmdFlags() cxCmdFlags {
	return cxCmdFlags{
		baseOutput:          false,
		compileOutput:       "",
		newProject:          false,
		replMode:            false,
		webMode:             false,
		ideMode:             false,
		webPersistentMode:   false,
		printHelp:           false,
		printVersion:        false,
		blockchainMode:      false,
		transactionMode:     false,
		broadcastMode:       false,
		port:                6001,
		programName:         "cxcoin",
		walletId:            "cxcoin_cli.wlt",
		secKey:              "",
		pubKey:              "",
		genesisAddress:      "",
		genesisSignature:    "",

		debugLexer:          false,
	}
}

func registerFlags(options *cxCmdFlags) {
	args := os.Args
	if len(args) <= 1 {
		options.replMode = true
	}

	flag.BoolVar(&options.printVersion, "version", options.printVersion, "Print CX version")
	flag.BoolVar(&options.printVersion, "v", options.printVersion, "alias for -version")
	flag.BoolVar(&options.tokenizeMode, "tokenize", options.tokenizeMode, "generate a 'out.cx.txt' text file with parsed tokens")
	flag.BoolVar(&options.tokenizeMode, "t", options.tokenizeMode, "alias for -tokenize")
	flag.StringVar(&options.compileOutput, "co", options.compileOutput, "alias for -compile-output")
	flag.BoolVar(&options.newProject, "new", options.newProject, "Creates a new project located at $CXPATH/src")
	flag.BoolVar(&options.newProject, "n", options.newProject, "alias for -new")
	flag.BoolVar(&options.replMode, "repl", options.replMode, "Loads source files into memory and starts a read-eval-print loop")
	flag.BoolVar(&options.replMode, "r", options.replMode, "alias for -repl")
	flag.BoolVar(&options.webMode, "web", options.webMode, "Start CX as a web service.")
	flag.BoolVar(&options.webMode, "w", options.webMode, "alias for -web")
	flag.BoolVar(&options.ideMode, "ide", options.ideMode, "Start CX as a web service, and Leaps service start also.")
	flag.BoolVar(&options.webPersistentMode, "pw", options.webPersistentMode, "Start CX as a web service with a persistent web REPL session")
	flag.StringVar(&options.initialHeap, "heap-initial", options.initialHeap, "Set the initial heap for the CX virtual machine. The value is in bytes, but the suffixes 'G', 'M' or 'K' can be used to express gigabytes, megabytes or kilobytes, respectively. Lowercase suffixes are allowed.")
	flag.StringVar(&options.initialHeap, "hi", options.initialHeap, "alias for -initial-heap")
	flag.StringVar(&options.maxHeap, "heap-max", options.maxHeap, "Set the max heap for the CX virtual machine. The value is in bytes, but the suffixes 'G', 'M' or 'K' can be used to express gigabytes, megabytes or kilobytes, respectively. Lowercase suffixes are allowed. Note that this parameter overrides --heap-initial if --heap-max is equal to a lesser value than --heap-max's.")
	flag.StringVar(&options.maxHeap, "hm", options.maxHeap, "alias for -max-heap")
	flag.StringVar(&options.stackSize, "stack-size", options.stackSize, "Set the stack size for the CX virtual machine. The value is in bytes, but the suffixes 'G', 'M' or 'K' can be used to express gigabytes, megabytes or kilobytes, respectively. Lowercase suffixes are allowed.")
	flag.StringVar(&options.stackSize, "ss", options.stackSize, "alias for -stack-size")
	flag.Float64Var(&options.minHeapFreeRatio, "--min-heap-free", options.minHeapFreeRatio, "Minimum heap space percentage that should be free after calling the garbage collector. Value must be in the range of 0.0 and 1.0.")
	flag.Float64Var(&options.maxHeapFreeRatio, "--max-heap-free", options.maxHeapFreeRatio, "Maximum heap space percentage that should be free after calling the garbage collector. Value must be in the range of 0.0 and 1.0.")

	flag.BoolVar(&options.blockchainMode, "blockchain", options.blockchainMode, "Start a CX blockchain program")
	// flag.BoolVar(&options.blockchainMode, "bc", options.blockchainMode, "alias for -blockchain")
	flag.BoolVar(&options.publisherMode, "publisher", options.publisherMode, "Start a CX blockchain program block publisher")
	// flag.BoolVar(&options.publisherMode, "pb", options.publisherMode, "alias for -publisher")
	flag.BoolVar(&options.transactionMode, "transaction", options.transactionMode, "Test a CX blockchain transaction")
	// flag.BoolVar(&options.transactionMode, "txn", options.transactionMode, "alias for -transaction")
	flag.BoolVar(&options.broadcastMode, "broadcast", options.broadcastMode, "Broadcast a CX blockchain transaction")
	flag.BoolVar(&options.walletMode, "create-wallet", options.walletMode, "Create a wallet from a seed")
	flag.BoolVar(&options.genAddress, "generate-address", options.genAddress, "Generate a CX chain address")
	flag.BoolVar(&options.peerMode, "peer", options.peerMode, "Run a CX chain peer node")
	flag.IntVar(&options.port, "port", options.port, "Port used when running a CX chain peer node")
	flag.StringVar(&options.walletSeed, "wallet-seed", options.walletSeed, "Seed to use for a new wallet")
	flag.StringVar(&options.walletId, "wallet-id", options.walletId, "Wallet ID to use for signing transactions")
	flag.StringVar(&options.programName, "program-name", options.programName, "Name of the initial CX program on the blockchain")
	flag.StringVar(&options.secKey, "secret-key", options.secKey, "CX program blockchain security key")
	flag.StringVar(&options.pubKey, "public-key", options.pubKey, "CX program blockchain public key")
	flag.StringVar(&options.genesisAddress, "genesis-address", options.genesisAddress, "CX blockchain program genesis address")
	flag.StringVar(&options.genesisSignature, "genesis-signature", options.genesisSignature, "CX blockchain program genesis address")
	flag.StringVar(&options.cxpath, "cxpath", options.cxpath, "Used for dynamically setting the value of the environment variable CXPATH")

	// Debug flags
	flag.BoolVar(&options.debugLexer, "debug-lexer", options.debugLexer, "Debug the lexer by printing all scanner tokens")
	flag.BoolVar(&options.debugLexer, "Dl",          options.debugLexer, "alias for -debug-lexer")
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
