package main

import (
	"fmt"
	"os"

	cxcore "github.com/SkycoinProject/cx/cx"
	"github.com/SkycoinProject/cx/cxgo/actions"
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
