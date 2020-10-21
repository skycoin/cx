package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"
	"strings"

	"github.com/SkycoinProject/cx-chains/src/cipher"

	"github.com/SkycoinProject/cx/cxgo/cxflags"
	"github.com/SkycoinProject/cx/cxgo/cxlexer"
	"github.com/SkycoinProject/cx/cxgo/cxspec"
	"github.com/SkycoinProject/cx/cxgo/parser"
	"github.com/SkycoinProject/cx/cxutil"
)

const filePerm = 0644

type newFlags struct {
	cmd *flag.FlagSet

	replace    bool
	unifyKeys  bool
	coinName   string
	coinTicker string

	debugLexer   bool
	debugProfile int

	chainSpecOut string
	chainKeysOut string
	genKeysOut   string

	*cxflags.MemoryFlags
}

func processNewFlags(args []string) newFlags {
	// Specify default flag values.
	f := newFlags{
		cmd: flag.NewFlagSet("cxchain-cli new", flag.ExitOnError),

		replace:      false,
		unifyKeys:    false,
		coinName:     "skycoin",
		coinTicker:   "SKY",

		debugLexer:   false,
		debugProfile: 0,

		chainSpecOut: "./{coin}.chain_spec.json",
		chainKeysOut: "./{coin}.chain_keys.json",
		genKeysOut:   "./{coin}.genesis_keys.json",

		MemoryFlags: cxflags.DefaultMemoryFlags(),
	}
	f.cmd.Usage = func() {
		usage := cxutil.DefaultUsageFormat("flags", "cx source files")
		usage(f.cmd, nil)
	}

	f.cmd.BoolVar(&f.replace, "replace", f.replace, "whether to replace output file(s)")
	f.cmd.BoolVar(&f.replace, "r", f.replace, "shorthand for 'replace'")
	f.cmd.BoolVar(&f.unifyKeys, "unify", f.unifyKeys, "whether to use the same keys for genesis and chain")
	f.cmd.BoolVar(&f.unifyKeys, "u", f.unifyKeys, "shorthand for 'unify'")
	f.cmd.StringVar(&f.coinName, "coin", f.coinName, "`NAME` for cx coin")
	f.cmd.StringVar(&f.coinName, "c", f.coinName, "shorthand for 'coin'")
	f.cmd.StringVar(&f.coinTicker, "ticker", f.coinTicker, "`SYMBOL` for cx coin ticker")
	f.cmd.StringVar(&f.coinTicker, "t", f.coinTicker, "shorthand for 'ticker'")

	f.cmd.BoolVar(&f.debugLexer, "debug-lexer", f.debugLexer, "enable lexer debugging by printing all scanner tokens")
	f.cmd.IntVar(&f.debugProfile, "debug-profile", f.debugProfile, "Enable CPU+MEM profiling and set CPU profiling rate. Visualize .pprof files with 'go get github.com/google/pprof' and 'pprof -http=:8080 file.pprof'")

	f.cmd.StringVar(&f.chainSpecOut, "chain-spec-output", f.chainSpecOut, "`FILE` for chain spec output")
	f.cmd.StringVar(&f.chainKeysOut, "chain-keys-output", f.chainKeysOut, "`FILE` for chain keys output")
	f.cmd.StringVar(&f.genKeysOut, "genesis-keys-output", f.genKeysOut, "`FILE` for genesis keys output")

	f.MemoryFlags.Register(f.cmd)

	// Parse flags.
	if err := f.cmd.Parse(args); err != nil {
		os.Exit(1)
	}

	// Post process.
	f.postProcess()
	if err := f.MemoryFlags.PostProcess(); err != nil {
		log.WithError(err).Fatal()
	}

	// Set stuff.
	if f.debugLexer {
		cxlexer.SetLogger(log)
	}

	// Log stuff.
	cxflags.LogMemFlags(log)

	// Return.
	return f
}

func (f *newFlags) postProcess() {
	replaceTokens := func(s *string, coinName, coinTicker string) {
		*s = strings.ReplaceAll(*s, "{coin}", coinName)
		*s = strings.ReplaceAll(*s, "{ticker}", coinTicker)
	}

	// Replace tokens in flag values with actual values.
	replaceTokens(&f.chainSpecOut, f.coinName, f.coinTicker)
	replaceTokens(&f.chainKeysOut, f.coinName, f.coinTicker)
	replaceTokens(&f.genKeysOut, f.coinName, f.coinTicker)

	// Check replace.
	if !f.replace {
		for _, name := range []string{f.chainSpecOut, f.chainKeysOut, f.genKeysOut} {
			if _, err := os.Stat(name); !os.IsNotExist(err) {
				errPrintf("File '%s' already exists. Replace with '--replace' flag.\n", name)
				os.Exit(1)
			}
		}
	}
}

func cmdNew(args []string) {
	flags := processNewFlags(args)

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

	// Parse and run program.
	if err := PrepareGenesisProg(cxFilenames, cxRes.CXSources, flags.debugLexer, flags.debugProfile); err != nil {
		log.WithError(err).Fatal("Failed to prepare genesis CX program.")
	}
	genProgState, err := RunGenesisProg(cxRes.CXFlags)
	if err != nil {
		log.WithError(err).Fatal("Failed to run genesis CX program.")
	}

	// Generate chain keys.
	chainPK, chainSK := cipher.GenerateKeyPair()

	// Generate genesis keys.
	genPK, genSK := chainPK, chainSK
	if !flags.unifyKeys {
		genPK, genSK = cipher.GenerateKeyPair()
	}
	genAddr := cipher.AddressFromPubKey(genPK)

	// Generate and write chain spec file.
	cSpec, err := cxspec.New(flags.coinName, flags.coinTicker, chainSK, genAddr, genProgState)
	if err != nil {
		log.WithError(err).
			Fatal("Failed to generate chain spec.")
	}
	cSpecB, err := json.MarshalIndent(cSpec, "", "\t")
	if err != nil {
		log.WithError(err).
			Fatal("Failed to encode chain spec to json.")
	}
	if err := ioutil.WriteFile(flags.chainSpecOut, cSpecB, filePerm); err != nil {
		log.WithError(err).
			WithField("filename", flags.chainSpecOut).
			Fatal("Failed to write chain spec to file.")
	}

	// Write chain keys file.
	cKeys := cxspec.KeySpecFromSecKey(cxspec.ChainKey, chainSK, true, true)
	cKeysB, err := json.MarshalIndent(cKeys, "", "\t")
	if err != nil {
		errPrintf("Failed to encode chain keys to json: %v\n", err)
		os.Exit(1)
	}
	if err := ioutil.WriteFile(flags.chainKeysOut, cKeysB, filePerm); err != nil {
		errPrintf("Failed to write chain keys to file '%s': %v\n", flags.chainKeysOut, err)
		os.Exit(1)
	}

	// Write genesis keys to file.
	gKeys := cxspec.KeySpecFromSecKey(cxspec.GenesisKey, genSK, true, true)
	gKeysB, err := json.MarshalIndent(gKeys, "", "\t")
	if err != nil {
		errPrintf("Failed to encode genesis keys to json: %v\n", err)
		os.Exit(1)
	}
	if err := ioutil.WriteFile(flags.genKeysOut, gKeysB, filePerm); err != nil {
		errPrintf("Failed to write genesis keys to file '%s': %v\n", flags.genKeysOut, err)
		os.Exit(1)
	}
}
