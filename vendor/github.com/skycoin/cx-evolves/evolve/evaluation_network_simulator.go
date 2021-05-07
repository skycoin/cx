package evolve

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/skycoin/cx-evolves/cxexecutes/worker"
	workerclient "github.com/skycoin/cx-evolves/cxexecutes/worker/client"
	cxast "github.com/skycoin/cx/cx/ast"
)

// perByteEvaluation for evolve with network sim, 1 i32 input, 1 i32 output
func perByteEvaluation_NetworkSim(ind *cxast.CXProgram, solPrototype *cxast.CXFunction, cfg *EvolveConfig) float64 {
	var score int = 0
	for rounds := 0; rounds < cfg.NumberOfRounds; rounds++ {
		// Generate random Input
		rand.Seed(time.Now().Unix())
		input := toByteArray(int32(rand.Int()))

		// Get output from transmitter
		transmitterOutput := perByteEvaluation_NetworkSim_Transmitter(ind, solPrototype, input, cfg.WorkerPortNum)

		// Input noise here

		// Get output from receiver
		receiverOutput := perByteEvaluation_NetworkSim_Receiver(ind, solPrototype, transmitterOutput, cfg.WorkerPortNum)

		// Get score by counting number of diff bits between generated input and receiverOutput
		score = score + countDifferentBits(input, receiverOutput)
	}

	return float64(score)
}

// perByteEvaluation for evolve with network sim transmitter, 1 i32 input, 1 i32 output
func perByteEvaluation_NetworkSim_Transmitter(ind *cxast.CXProgram, solPrototype *cxast.CXFunction, input []byte, workerPortNum int) []byte {
	var tmp *cxast.CXProgram = cxast.PROGRAM
	cxast.PROGRAM = ind

	inpFullByteSize := 0
	for c := 0; c < len(solPrototype.Inputs); c++ {
		inpFullByteSize += solPrototype.Inputs[c].TotalSize
	}

	// We'll store the `i`th inputs on `inps`.
	inps := make([]byte, inpFullByteSize)

	inp := input

	// Copying the input `b`ytes.
	for b := 0; b < len(inp); b++ {
		inps[b] = inp[b]
	}

	var result worker.Result
	workerAddr := fmt.Sprintf(":%v", workerPortNum)
	workerclient.CallWorker(
		workerclient.CallWorkerConfig{
			Program:   ind,
			Input:     inps,
			OutputArg: solPrototype.Outputs[0],
		},
		workerAddr,
		&result,
	)

	data := result.Output

	cxast.PROGRAM = tmp
	return data
}

// perByteEvaluation for evolve with network sim receiver, 1 i32 input, 1 i32 output
func perByteEvaluation_NetworkSim_Receiver(ind *cxast.CXProgram, solPrototype *cxast.CXFunction, input []byte, workerPortNum int) []byte {
	var tmp *cxast.CXProgram = cxast.PROGRAM
	cxast.PROGRAM = ind

	inpFullByteSize := 0
	for c := 0; c < len(solPrototype.Inputs); c++ {
		inpFullByteSize += solPrototype.Inputs[c].TotalSize
	}

	// We'll store the `i`th inputs on `inps`.
	inps := make([]byte, inpFullByteSize)

	inp := input

	// Copying the input `b`ytes.
	for b := 0; b < len(inp); b++ {
		inps[b] = inp[b]
	}

	var result worker.Result
	workerAddr := fmt.Sprintf(":%v", workerPortNum)
	workerclient.CallWorker(
		workerclient.CallWorkerConfig{
			Program:   ind,
			Input:     inps,
			OutputArg: solPrototype.Outputs[0],
		},
		workerAddr,
		&result,
	)

	data := result.Output

	cxast.PROGRAM = tmp
	return data
}

func countDifferentBits(a []byte, b []byte) int {
	var count int

	for i, val := range a {
		bitPos := 1
		for z := 0; z < 8; z++ {
			if int(val)&bitPos != int(b[i])&bitPos {
				count++
			}
			bitPos = bitPos << 1
		}
	}
	return count
}
