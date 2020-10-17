package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/SkycoinProject/cx-chains/src/api"
	"github.com/SkycoinProject/cx-chains/src/cipher"
	"github.com/SkycoinProject/cx-chains/src/coin"
	"github.com/SkycoinProject/cx-chains/src/transaction"
	"github.com/SkycoinProject/cx-chains/src/util/fee"

	cxcore "github.com/SkycoinProject/cx/cx"
	"github.com/SkycoinProject/cx/cxgo/actions"
	"github.com/SkycoinProject/cx/cxgo/cxlexer"
)

// PrepareChainProg parses a program on chain, and loads additional sources onto
// the program state.
func PrepareChainProg(filenames []string, srcs []*os.File, nodeAddr string, addr cipher.Address, debugLexer bool, debugProf int) (*cxlexer.ProgBytes, error) {
	// TODO @evanlinjin: Enable profiling later.
	// _, stopProf := StartPreparationProfiling("PrepareChainProg", debugLexer, debugProf)
	// defer stopProf()

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

// RunChainProg runs a on-chain program state combined with a main function expression.
func RunChainProg(cxArgs []string, progB *cxlexer.ProgBytes) ([]byte, error) {
	// TODO @evanlinjin: Enable profiling later.
	// log := log.WithField("func", "RunChainProg")
	// _, stopProf := cxprof.StartProfile(log)
	// defer stopProf()

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

func ObtainProgramState(c *api.Client, addr cipher.Address) ([]byte, error) {
	utxoSum, err := c.OutputsForAddresses([]string{addr.String()})
	if err != nil {
		return nil, err
	}

	// Obtain ux with program state.
	uxArr, err := utxoSum.SpendableOutputs().ToUxArray()
	if err != nil {
		return nil, err
	}
	j, _ := json.MarshalIndent(uxArr[0], "", "\t")
	log.WithField("ux_len", len(uxArr)).Info(string(j))

	var progState []byte
	for _, ux := range uxArr {
		if len(ux.Body.ProgramState) > 0 {
			progState = ux.Body.ProgramState
			break
		}
	}
	if progState == nil {
		return nil, fmt.Errorf("failed to find output owned by '%s' that contains a program state", addr)
	}
	return progState, nil
}

func BroadcastMainExp(nodeAddr string, genSK cipher.SecKey, pb *cxlexer.ProgBytes) error {
	// Setting the CX runtime to run `PRGRM`.
	if _, err := actions.PRGRM.SelectProgram(); err != nil {
		return fmt.Errorf("failed to select program: %w", err)
	}
	cxcore.MarkAndCompact(actions.PRGRM)

	mainB := cxcore.ExtractBlockchainProgram(
		pb.State,
		cxcore.Serialize(actions.PRGRM, actions.PRGRM.BCPackageCount))

	genAddr, err := cipher.AddressFromSecKey(genSK)
	if err != nil {
		return err
	}

	c := api.NewClient(nodeAddr)

	utxoSum, err := c.OutputsForAddresses([]string{genAddr.String()})
	if err != nil {
		return err
	}

	// Obtain ux with program state.
	uxArr, err := utxoSum.SpendableOutputs().ToUxArray()
	if err != nil {
		return err
	}
	var progUxOut coin.UxOut // <-- UTXO with program state.
	for _, ux := range uxArr {
		if len(ux.Body.ProgramState) > 0 {
			progUxOut = ux
			break
		}
	}
	if progUxOut.Body.ProgramState == nil {
		return fmt.Errorf("failed to obtain UTXO from '%s' that contains a program state", genAddr)
	}

	tx := new(coin.Transaction)
	tx.MainExpressions = mainB

	// TODO @evanlinjin: Inject transaction here!
	injectOut, err := c.InjectTransaction(nil)
	if err != nil {
		return err
	}
	log.WithField("inject_out", injectOut).Info("CX transaction injected.")

	return nil
}

func CreateCXTransaction(genSK cipher.SecKey, uxArray coin.UxArray, headTime uint64, burnFactor uint32, progState, mainExp []byte) (*coin.Transaction, error) {
	genAddr, err := cipher.AddressFromSecKey(genSK)
	if err != nil {
		return nil, err
	}
	fmt.Println(genAddr.String())

	tx := new(coin.Transaction)
	tx.MainExpressions = mainExp

	uxBalances, err := transaction.NewUxBalances(uxArray, headTime)
	if err != nil {
		return nil, err
	}

	uxBalancesMap := make(map[cipher.SHA256]transaction.UxBalance, len(uxBalances))
	for i, b := range uxBalances {
		if _, ok := uxBalancesMap[b.Hash]; ok {
			return nil, fmt.Errorf("ux balance at index %d is a duplicate", i)
		}
		uxBalancesMap[b.Hash] = b
	}

	// sum balances
	var outputCoins uint64
	var outputHours uint64
	for _, b := range uxBalances {
		outputCoins += b.Coins
		outputHours += b.Hours
	}

	// coin hours fee
	outputHours -= fee.RequiredFee(outputHours, burnFactor)

	// Calculate total coins and minimum hours to send
	return nil, nil
}
