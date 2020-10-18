package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/SkycoinProject/cx-chains/src/api"
	"github.com/SkycoinProject/cx-chains/src/cipher"

	"github.com/SkycoinProject/cx/cxgo/cxflags"
	"github.com/SkycoinProject/cx/cxgo/cxspec"
	"github.com/SkycoinProject/cx/cxgo/parser"
	"github.com/SkycoinProject/cx/cxutil"
)

const (
	// ENV for chain spec filepath.
	specFileEnv          = "CXCHAIN_SPEC_FILEPATH"
	defaultChainSpecFile = "./skycoin.chain_spec.json"

	// ENV for genesis secret key.
	genSKEnv = "CXCHAIN_GENESIS_SK"
)

// parseSpecFilepathEnv parses chain spec filename from CXCHAIN_SPEC_FILEPATH env.
func parseSpecFilepathEnv() cxspec.ChainSpec {
	specFilename, ok := os.LookupEnv(specFileEnv)
	if !ok {
		specFilename = defaultChainSpecFile
	}

	log.WithField("filename", specFilename).Info("Reading chain spec file...")

	spec, err := cxspec.ReadSpecFile(specFilename)
	if err != nil {
		log.WithError(err).Fatal("Failed to start node.")
	}

	// TODO @evanlinjin: Need to fix genesis program state being atrociously massive.
	// spec.Print()

	cxspec.PopulateParamsModule(spec)

	return spec
}

// parseSecretKeyEnv parses secret key from CXCHAIN_SECRET_KEY env.
// The secret key can be null.
func parseSecretKeyEnv() cipher.SecKey {
	if skStr, ok := os.LookupEnv(genSKEnv); ok {
		sk, err := cipher.SecKeyFromHex(skStr)
		if err != nil {
			log.WithError(err).WithField("ENV", genSKEnv).Fatal("Provided genesis secret key is invalid.")
		}
		return sk
	}
	return cipher.SecKey{} // return nil secret key
}

type runFlags struct {
	cmd *flag.FlagSet

	debugLexer   bool
	debugProfile int
	*cxflags.MemoryFlags

	inject   bool   // Whether to inject transaction to cx chain.
	nodeAddr string // CX Chain node address.
}

func processRunFlags(args []string) (runFlags, cxspec.ChainSpec, cipher.SecKey) {
	spec := parseSpecFilepathEnv()
	genSK := parseSecretKeyEnv()

	// Check genesis secret key.
	if !genSK.Null() {
		genAddr, err := cipher.AddressFromSecKey(genSK)
		if err != nil {
			log.WithError(err).
				WithField(genSKEnv, genSK.Hex()).
				Fatal("Failed to extract genesis address.")
		}

		if specAddr := cipher.MustDecodeBase58Address(spec.GenesisAddr); genAddr != specAddr {
			log.WithField(genSKEnv, genSK.Hex()).
				Fatal("Provided genesis secret key does not match genesis address from chain spec.")
		}
	}

	f := runFlags{
		cmd: flag.NewFlagSet(args[0], flag.ExitOnError),

		debugLexer:   false,
		debugProfile: 0,
		MemoryFlags:  cxflags.DefaultMemoryFlags(),

		inject:   false,
		nodeAddr: fmt.Sprintf("http://127.0.0.1:%d", spec.Node.WebInterfacePort),
	}

	f.cmd.Usage = func() {
		_, _ = fmt.Fprintf(os.Stderr, "usage: %s %s [args...] [cx source files...]\n", os.Args[0], os.Args[1])
		f.cmd.PrintDefaults()
	}

	f.cmd.BoolVar(&f.debugLexer, "debug-lexer", f.debugLexer, "enable lexer debugging by printing all scanner tokens")
	f.cmd.IntVar(&f.debugProfile, "debug-profile", f.debugProfile, "enable CPU+MEM profiling and set CPU profiling rate. Visualize .pprof files with 'go get github.com/google/pprof' and 'pprof -http=:8080 file.pprof'")
	f.MemoryFlags.Register(f.cmd)

	f.cmd.BoolVar(&f.inject, "inject", f.inject, "whether to inject this as a transaction on the cx chain")
	f.cmd.BoolVar(&f.inject, "i", f.inject, "shorthand for 'inject'")

	f.cmd.StringVar(&f.nodeAddr, "node", f.nodeAddr, "HTTP API `ADDRESS` of cxchain node")
	f.cmd.StringVar(&f.nodeAddr, "n", f.nodeAddr, "shorthand for 'node'")

	// Parse flags.
	parseFlagSet(f.cmd, args[1:])

	// Log stuff.
	cxflags.LogMemFlags(log)

	// Return.
	return f, spec, genSK
}

func cmdRun(args []string) {
	flags, spec, genSK := processRunFlags(args)

	// Apply debug flags.
	parser.DebugLexer = flags.debugLexer

	// Parse for cx args for genesis program state.
	log.Info("Parsing for CX args...")
	cxRes, err := cxutil.ExtractCXArgs(flags.cmd, true)
	if err != nil {
		log.WithError(err).Fatal("Failed to extract CX args.")
	}
	cxFilenames := cxutil.ListSourceNames(cxRes.CXSources, true)
	log.WithField("filenames", cxFilenames).Info("Obtained CX sources.")

	// Prepare API Client.
	c := api.NewClient(flags.nodeAddr)

	// Prepare address.
	addr := cipher.MustDecodeBase58Address(spec.GenesisAddr)

	// Parse and run program.
	ux, progB, err := PrepareChainProg(cxFilenames, cxRes.CXSources, c, addr, flags.debugLexer, flags.debugProfile,)
	if err != nil {
		log.WithError(err).Fatal("Failed to prepare chain CX program.")
	}

	if flags.inject {
		// Run: inject.
		if err := BroadcastMainExp(c, genSK, ux); err != nil {
			log.WithError(err).Fatal("Failed to broadcast transaction.")
		}
	} else {
		// Run: without injection.
		if err := RunChainProg(cxRes.CXFlags, progB); err != nil {
			log.WithError(err).Fatal("Failed to run chain CX program.")
		}
	}
}
