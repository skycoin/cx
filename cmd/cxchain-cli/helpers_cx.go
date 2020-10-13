package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/SkycoinProject/cx-chains/src/cipher"
	"github.com/sirupsen/logrus"

	cxcore "github.com/SkycoinProject/cx/cx"
	"github.com/SkycoinProject/cx/cxgo/actions"
	"github.com/SkycoinProject/cx/cxgo/cxgo0"
	"github.com/SkycoinProject/cx/cxgo/cxlexer"
	"github.com/SkycoinProject/cx/cxgo/cxprof"
)

// PrepareGenesisProg parses a genesis program (a new program).
func PrepareGenesisProg(filenames []string, srcs []*os.File, debugLexer bool, debugProf int) error {
	_, stopProf := initParseProf("PrepareGenesisProg", debugLexer, debugProf)
	defer stopProf()

	// Prepare core program state for 'actions.PRGRM'.
	prog, err := cxlexer.InitProg()
	if err != nil {
		return fmt.Errorf("failed to obtain prog. state of core packages: %w", err)
	}

	// Compile sources.
	if err := compileSources(prog, filenames, srcs); err != nil {
		return fmt.Errorf("failed to compile sources: %w", err)
	}
	return nil
}

// PrepareChainProg parses a program on chain, and loads additional sources onto
// the program state.
func PrepareChainProg(filenames []string, srcs []*os.File, nodeAddr string, addr cipher.Address, debugLexer bool, debugProf int) (*cxlexer.ProgBytes, error) {
	_, stopProf := initParseProf("PrepareChainProg", debugLexer, debugProf)
	defer stopProf()

	// Prepare core program state for 'actions.PRGRM'.
	prog, err := cxlexer.InitProg()
	if err != nil {
		return nil, fmt.Errorf("failed to init prog: %w", err)
	}
	progB, err := cxlexer.LoadProgFromChain(prog, nodeAddr, addr)
	if err != nil {
		return nil, fmt.Errorf("failed to load onto prog from chain: %w", err)
	}

	// Compile sources.
	if err := compileSources(prog, filenames, srcs); err != nil {
		return nil, fmt.Errorf("failed to compile sources: %w", err)
	}
	return progB, nil
}

// RunGenesisProg initiates a blockchain program and returns the genesis
// program state.
func RunGenesisProg(cxArgs []string) ([]byte, error) {
	log := log.WithField("func", "RunGenesisProg")

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
	cxcore.MarkAndCompact(actions.PRGRM)
	actions.PRGRM.HeapSize = actions.PRGRM.HeapPointer

	// As we have removed the 'main' pkg, blockchain pkg count is len(prog.)
	// instead of len(prog.)-1.
	actions.PRGRM.BCPackageCount = len(actions.PRGRM.Packages)

	progB := cxcore.Serialize(actions.PRGRM, actions.PRGRM.BCPackageCount)
	progB = cxcore.ExtractBlockchainProgram(progB, progB)
	log.WithField("size", len(progB)).Info("Obtained serialized program state.")
	return progB, nil
}

func RunChainProg(cxArgs []string, progB *cxlexer.ProgBytes) ([]byte, error) {
	log := log.WithField("func", "RunChainProg")

	_, stopProf := cxprof.StartProfile(log)
	defer stopProf()

	// If it's a CX chain transaction, we need to add the heap extracted
	// from the retrieved CX chain program state.
	if err := progB.MergeChainHeap(); err != nil {
		return nil, err
	}

	// Run as normal CX program (for now).
	if err := actions.PRGRM.RunCompiled(0, cxArgs); err != nil {
		return nil, err
	}

	if cxcore.AssertFailed() {
		return nil, fmt.Errorf("assert failed: %v", cxcore.CX_ASSERT)
	}

	return nil, nil
}

/*
	<<< HELPERS >>>
*/

func initParseProf(funcName string, debugLexer bool, debugProf int) (logrus.FieldLogger, func()) {
	log := log.WithField("func", funcName)

	// Start CPU profiling.
	stopCPUProf, err := cxprof.StartCPUProfile(funcName, debugProf)
	if err != nil {
		log.WithError(err).Error("Failed to start CPU profiling.")
	}

	// Start Log profiling.
	var stopProf cxprof.StopFunc
	if debugLexer {
		_, stopProf = cxprof.StartProfile(log)
	}

	return log, func() {
		// Dump memory state.
		if err := cxprof.DumpMemProfile(funcName); err != nil {
			log.WithError(err).Error("Failed to dump MEM profile.")
		}

		// Stop Log profiling.
		if stopProf != nil {
			stopProf()
		}

		// Stop CPU profiling.
		if err := stopCPUProf(); err != nil {
			log.WithError(err).Error("Failed to stop CPU profiling.")
		}
	}
}

func compileSources(prog *cxcore.CXProgram, filenames []string, srcs []*os.File) error {
	// Actually parse source code.
	if exitCode := cxlexer.ParseSourceCode(srcs, filenames); exitCode != 0 {
		return fmt.Errorf("cxlexer.ParseSourceCode returned with code %d", exitCode)
	}

	// Set working directory.
	if len(srcs) > 0 {
		cxgo0.PRGRM0.Path = determineWorkDir(srcs[0].Name())
	}

	// Add main function if not exist.
	ensureCXMainFunc(prog)

	// Add *init function that initializes all global variables.
	if err := ensureCXInitFunc(prog); err != nil {
		return fmt.Errorf("failed to setup *init func: %w", err)
	}

	// Reset.
	actions.LineNo = 0

	// Check and return.
	if cxcore.FoundCompileErrors {
		return fmt.Errorf("cxcore has compilation error code %d", cxcore.CX_COMPILATION_ERROR)
	}
	return nil
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
