package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/SkycoinProject/cx-chains/src/cipher"

	cxcore "github.com/SkycoinProject/cx/cx"
	"github.com/SkycoinProject/cx/cxgo/actions"
	"github.com/SkycoinProject/cx/cxgo/cxflags"
	"github.com/SkycoinProject/cx/cxgo/cxgo0"
	"github.com/SkycoinProject/cx/cxgo/cxlexer"
	"github.com/SkycoinProject/cx/cxgo/cxprof"
	"github.com/SkycoinProject/cx/cxgo/cxspec"
	"github.com/SkycoinProject/cx/cxgo/parser"
	"github.com/SkycoinProject/cx/cxutil"
)

const filePerm = 0644

type newChainFlags struct {
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

func processNewChainFlags(args []string) newChainFlags {
	// Specify default flag values.
	f := newChainFlags{
		cmd: flag.NewFlagSet(args[0], flag.ExitOnError),

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
	parseFlagSet(f.cmd, args[1:])

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

func (f *newChainFlags) postProcess() {
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

// parseProgram does...
// It returns the exit code.
func parseProgram(flags *newChainFlags, filenames []string, srcs []*os.File) int {
	log := log.WithField("func", "parseProgram")

	// Start CPU profiling.
	stopCPUProf, err := cxprof.StartCPUProfile("parseProgram", flags.debugProfile)
	if err != nil {
		log.WithError(err).Error("Failed to start CPU profiling.")
	}
	defer func() {
		if err := stopCPUProf(); err != nil {
			log.WithError(err).Error("Failed to stop CPU profiling.")
		}
	}()

	// Start log profiling.
	if flags.debugLexer {
		_, stopProf := cxprof.StartProfile(log)
		defer stopProf()
	}

	// Dump memory state.
	defer func() {
		if err := cxprof.DumpMemProfile("parseProgram"); err != nil {
			log.WithError(err).Error("Failed to dump MEM profile.")
		}
	}()

	// Prepare core program state for 'actions.PRGRM'.
	coreProgState, err := cxcore.GetProgram()
	if err != nil {
		log.WithError(err).Error("Failed to obtain prog. state of core packages.")
		return 1
	}
	prog := cxcore.MakeProgram()
	prog.Packages = coreProgState.Packages
	actions.PRGRM = prog

	// TODO @evanlinjin: We need some sort of prelude for transaction/broadcast mode.

	// Parse source code.
	if exitCode := cxlexer.ParseSourceCode(srcs, filenames); exitCode != 0 {
		log.Error("Failed to parse source code.")
		return exitCode
	}

	// Set working directory.
	if len(srcs) > 0 {
		cxgo0.PRGRM0.Path = determineWorkDir(srcs[0].Name())
	}

	// Add main function if not exist.
	ensureCXMainFunc(prog)

	// Add *init function that initializes all global variables.
	if err := ensureCXInitFunc(prog); err != nil {
		log.WithError(err).Error("Failed to setup *init CX function.")
		return 1
	}

	// Reset
	actions.LineNo = 0

	if cxcore.FoundCompileErrors {
		return cxcore.CX_COMPILATION_ERROR
	}

	return 0
}

// ensureCXMainFunc ensures that the CX program contains a main function.
func ensureCXMainFunc(prog *cxcore.CXProgram) {
	if _, err := prog.GetFunction(cxcore.MAIN_FUNC, cxcore.MAIN_PKG); err != nil {
		mainPkg := cxcore.MakePackage(cxcore.MAIN_PKG)
		prog.AddPackage(mainPkg)
		mainFn := cxcore.MakeFunction(cxcore.MAIN_FUNC, actions.CurrentFile, actions.LineNo)
		mainPkg.AddFunction(mainFn)
	}
}

// ensureCXInitFunc ensures that the CX program contains an *init function which
// initiates all global variables.
func ensureCXInitFunc(prog *cxcore.CXProgram) error {
	mainPkg, err := prog.GetPackage(cxcore.MAIN_PKG)
	if err != nil {
		return fmt.Errorf("failed to obtain main package: %w", err)
	}

	initFn := cxcore.MakeFunction(cxcore.SYS_INIT_FUNC, actions.CurrentFile, actions.LineNo)
	mainPkg.AddFunction(initFn)
	actions.FunctionDeclaration(initFn, nil, nil, actions.SysInitExprs)

	if _, err := prog.SelectFunction(cxcore.MAIN_FUNC); err != nil {
		return fmt.Errorf("failed to select main package: %w", err)
	}

	return nil
}

func determineWorkDir(filename string) (wkDir string) {
	log := log.WithField("func", "determineWorkDir")
	defer func() {
		log.WithField("work_dir", wkDir).Info()
	}()

	filename = filepath.FromSlash(filename)

	i := strings.LastIndexByte(filename, os.PathSeparator)
	if i == -1 {
		return ""
	}
	return filename[:i]
}

// initiateProgram initiates a blockchain program and returns the genesis
// program state.
func initiateProgram(cxArgs []string) ([]byte, error) {
	log := log.WithField("func", "initiateProgram")

	_, stopProf := cxprof.StartProfile(log)
	defer stopProf()

	// Initialize CX chain runtime?
	if err := actions.PRGRM.RunCompiled(0, cxArgs); err != nil {
		return nil, fmt.Errorf("failed to run compiled cx program: %w", err)
	}

	// Strip main package.
	actions.PRGRM.RemovePackage(cxcore.MAIN_PKG)

	// Remove garbage from heap.
	// Only keep global variables as these are independent from function calls.
	fmt.Println("Old heap:", actions.PRGRM.HeapPointer)
	cxcore.MarkAndCompact(actions.PRGRM)
	actions.PRGRM.HeapSize = actions.PRGRM.HeapPointer
	fmt.Println("New heap:", actions.PRGRM.HeapPointer)

	// As we have removed the 'main' pkg, blockchain pkg count is len(prog.)
	// instead of len(prog.)-1.
	actions.PRGRM.BCPackageCount = len(actions.PRGRM.Packages)

	progB := cxcore.Serialize(actions.PRGRM, actions.PRGRM.BCPackageCount)
	fmt.Println("Serialized program state:", len(progB))

	progB = cxcore.ExtractBlockchainProgram(progB, progB)
	log.WithField("size", len(progB)).Info("Obtained serialized program state.")
	return progB, nil
}

func cmdNewChain(args []string) {
	flags := processNewChainFlags(args)

	// Apply debug flags.
	parser.DebugLexer = flags.debugLexer

	// Parse for cx args for genesis program state.
	log.Info("Parsing for CX args...")
	cxRes, err := cxutil.ExtractCXArgs(flags.cmd, true)
	if err != nil {
		log.WithError(err).Fatal("Failed to extract CX args.")
	}
	cxFilenames := cxutil.ListSourceNames(cxRes.CXSources, true)
	fmt.Println("Filenames:", cxFilenames)

	// Parse and run program.
	if code := parseProgram(&flags, cxFilenames, cxRes.CXSources); code != 0 {
		os.Exit(code)
	}
	genProgState, err := initiateProgram(cxRes.CXFlags)
	if err != nil {
		errPrintf("Failed to run CX program: %v\n", err)
		os.Exit(1)
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
		errPrintf("Failed to generate chain spec: %v\n", err)
		os.Exit(1)
	}
	cSpecB, err := json.MarshalIndent(cSpec, "", "\t")
	if err != nil {
		errPrintf("Failed to encode chain spec to json: %v\n", err)
		os.Exit(1)
	}
	if err := ioutil.WriteFile(flags.chainSpecOut, cSpecB, filePerm); err != nil {
		errPrintf("Failed to write chain spec to file '%s': %v\n", flags.chainSpecOut, err)
		os.Exit(1)
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
