package main

import (
	"fmt"
	"os"

	"github.com/SkycoinProject/cx-chains/src/cipher"

	cxcore "github.com/SkycoinProject/cx/cx"
	"github.com/SkycoinProject/cx/cxgo/actions"
	"github.com/SkycoinProject/cx/cxgo/cxlexer"
	"github.com/SkycoinProject/cx/cxgo/cxprof"
)

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

func RunChainProg(cxArgs []string, progB *cxlexer.ProgBytes) ([]byte, error) {
	log := log.WithField("func", "RunChainProg")

	// TODO @evanlinjin: Figure out how to
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
