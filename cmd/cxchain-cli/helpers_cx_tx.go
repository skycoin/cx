package main

import (
	"fmt"
	"os"

	"github.com/SkycoinProject/cx-chains/src/api"
	"github.com/SkycoinProject/cx-chains/src/cipher"
	"github.com/SkycoinProject/cx-chains/src/coin"

	cxcore "github.com/SkycoinProject/cx/cx"
	"github.com/SkycoinProject/cx/cxgo/actions"
	"github.com/SkycoinProject/cx/cxgo/cxlexer"
)

// PrepareChainProg parses a program on chain, and loads additional sources onto
// the program state.
func PrepareChainProg(filenames []string, srcs []*os.File, c *api.Client, addr cipher.Address, debugLexer bool, debugProf int) (*coin.UxOut, *cxlexer.ProgBytes, error) {
	// TODO @evanlinjin: Enable profiling later.
	// _, stopProf := StartPreparationProfiling("PrepareChainProg", debugLexer, debugProf)
	// defer stopProf()

	// Obtain chain state UX.
	// TODO @evanlinjin: Implement retry logic.
	ux, err := ObtainProgramStateUxOut(c, addr)
	if err != nil {
		return nil, nil, err
	}

	// Prepare core program state for 'actions.PRGRM'.
	prog, err := cxlexer.InitProg()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to init prog: %w", err)
	}
	progB, err := cxlexer.LoadProgFromBytes(prog, ux.Body.ProgramState)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load prog state: %w", err)
	}

	// Merge chain heap.
	// If it's a CX chain transaction, we need to add the heap extracted
	// from the retrieved CX chain program state.
	if err := progB.MergeChainHeap(); err != nil {
		return nil, nil, fmt.Errorf("failed to merge chain heap: %w", err)
	}

	// Compile sources.
	if err := compileSources(prog, filenames, srcs); err != nil {
		return nil, nil, fmt.Errorf("failed to compile sources: %w", err)
	}
	return &ux, progB, nil
}

// RunChainProg runs a on-chain program state combined with a main function expression.
func RunChainProg(cxArgs []string, progB *cxlexer.ProgBytes) error {
	// TODO @evanlinjin: Enable profiling later.
	// log := log.WithField("func", "RunChainProg")
	// _, stopProf := cxprof.StartProfile(log)
	// defer stopProf()

	// Run as normal CX program (for now).
	if err := actions.PRGRM.RunCompiled(0, cxArgs); err != nil {
		return err
	}

	if cxcore.AssertFailed() {
		return fmt.Errorf("assert failed: %v", cxcore.CX_ASSERT)
	}

	return nil
}

func ObtainProgramStateUxOut(c *api.Client, addr cipher.Address) (coin.UxOut, error) {
	utxoSum, err := c.OutputsForAddresses([]string{addr.String()})
	if err != nil {
		return coin.UxOut{}, err
	}

	// Obtain ux with program state.
	uxArr, err := utxoSum.SpendableOutputs().ToUxArray()
	if err != nil {
		return coin.UxOut{}, err
	}

	var uxPS coin.UxOut
	for _, ux := range uxArr {
		if len(ux.Body.ProgramState) > 0 {
			uxPS = ux
			break
		}
	}
	if uxPS.Body.ProgramState == nil {
		return coin.UxOut{}, fmt.Errorf("failed to find output owned by '%s' that contains a program state", addr)
	}
	return uxPS, nil
}

func BroadcastMainExp(c *api.Client, genSK cipher.SecKey, ux *coin.UxOut) error {
	// Setup: extract main expression method.
	extractMainExp := func(oldProgS []byte) ([]byte, error) {
		if _, err := actions.PRGRM.SelectProgram(); err != nil {
			return nil, err
		}

		cxcore.MarkAndCompact(actions.PRGRM)

		s := cxcore.Serialize(actions.PRGRM, actions.PRGRM.BCPackageCount)
		mainExp := cxcore.ExtractBlockchainProgram(oldProgS, s)

		log.WithField("len", len(mainExp)).Info("Extracted main expression.")

		return mainExp, nil
	}

	// Setup: determine new program state.
	determineNewProgState := func(oldProgS, mainExp []byte) ([]byte, error) {
		// Running the merged program.
		if err := actions.PRGRM.RunCompiled(0, nil); err != nil {
			return nil, fmt.Errorf("failed to run prog: %w", err)
		}
		// Removing garbage from the heap. Only the global variables should be left
		// as these are independent from function calls.
		cxcore.MarkAndCompact(actions.PRGRM)

		// TODO: CX chains only work with one package at the moment (in the blockchain code). That is what that "1" is for.
		// Serializing the terminated program.
		s := cxcore.Serialize(actions.PRGRM, actions.PRGRM.BCPackageCount)
		// Extracting only the blockchain code. This is our new program state.
		newProgS := cxcore.ExtractBlockchainProgram(oldProgS, s)

		return newProgS, nil
	}

	// Setup: create cx transaction
	createTx := func(newProgS, mainExp []byte) (*coin.Transaction, error) {
		tx := new(coin.Transaction)
		tx.MainExpressions = mainExp

		if err := tx.PushInput(ux.Hash()); err != nil {
			return nil, fmt.Errorf("failed to push input: %w", err)
		}
		if err := tx.PushOutput(ux.Body.Address, ux.Body.Coins, ux.Body.Hours/2, newProgS); err != nil {
			return nil, err
		}
		if err := tx.UpdateHeader(); err != nil {
			return nil, err
		}
		if err := tx.SignInput(genSK, 0); err != nil {
			return nil, err
		}
		if err := tx.UpdateHeader(); err != nil {
			return nil, err
		}
		return tx, nil
	}

	// Run.
	mainExp, err := extractMainExp(ux.Body.ProgramState)
	if err != nil {
		return fmt.Errorf("failed to extract main exp: %w", err)
	}
	newProgS, err := determineNewProgState(ux.Body.ProgramState, mainExp)
	if err != nil {
		return fmt.Errorf("failed to determine new prog state: %w", err)
	}
	tx, err := createTx(newProgS, mainExp)
	if err != nil {
		return fmt.Errorf("failed to create tx: %w", err)
	}
	injectOut, err := c.InjectTransaction(tx)
	if err != nil {
		return fmt.Errorf("failed to inject tx: %w", err)
	}
	log.WithField("inject_out", injectOut).Info("CX transaction injected.")

	return nil
}

// func CreateCXTransaction(ux coin.UxOut, mainExp, newProgState []byte) (*coin.Transaction, error) {
// 	tx := new(coin.Transaction)
// 	tx.MainExpressions = mainExp
//
// 	if err := tx.PushInput(ux.Hash()); err != nil {
// 		return nil, fmt.Errorf("failed to push input: %w", err)
// 	}
// 	if err := tx.PushOutput(ux.Body.Address, ux.Body.Coins, ux.Body.Hours/2, newProgState); err != nil {
// 		return nil, err
// 	}
// 	if err := tx.UpdateHeader(); err != nil {
// 		return nil, err
// 	}
// 	return tx, nil
// }
