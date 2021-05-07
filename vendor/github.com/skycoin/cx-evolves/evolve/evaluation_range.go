package evolve

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"time"

	"github.com/skycoin/cx-evolves/cxexecutes/worker"
	workerclient "github.com/skycoin/cx-evolves/cxexecutes/worker/client"
	cxast "github.com/skycoin/cx/cx/ast"
)

// perByteEvaluation for evolve with range, 1 i32 input, 1 i32 output
func perByteEvaluation_Range(ind *cxast.CXProgram, solPrototype *cxast.CXFunction, cfg *EvolveConfig) float64 {
	var points int64 = 0
	var tmp *cxast.CXProgram = cxast.PROGRAM
	cxast.PROGRAM = ind

	inpFullByteSize := 0
	for c := 0; c < len(solPrototype.Inputs); c++ {
		inpFullByteSize += solPrototype.Inputs[c].TotalSize
	}

	// We'll store the `i`th inputs on `inps`.
	inps := make([]byte, inpFullByteSize)
	for round := 0; round < cfg.NumberOfRounds; round++ {
		rand.Seed(time.Now().Unix())
		in := round
		for in == 0 {
			in = rand.Int()
		}

		inp := toByteArray(int32(in))

		// Copying the input `b`ytes.
		for b := 0; b < len(inp); b++ {
			inps[b] = inp[b]
		}

		var result worker.Result
		workerAddr := fmt.Sprintf(":%v", cfg.WorkerPortNum)
		workerclient.CallWorker(
			workerclient.CallWorkerConfig{
				Program:   ind,
				Input:     inps,
				OutputArg: solPrototype.Outputs[0],
			},
			workerAddr,
			&result,
		)

		data := int(binary.BigEndian.Uint32(result.Output))
		// if not within range, add 1 to total points
		if data > cfg.UpperRange || data < cfg.LowerRange {
			points++
		}
	}

	cxast.PROGRAM = tmp
	return float64(points)
}
