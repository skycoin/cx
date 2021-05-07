package worker

import (
	"fmt"

	"github.com/henrylee2cn/erpc/v6"
	cxast "github.com/skycoin/cx/cx/ast"
	cxexecute "github.com/skycoin/cx/cx/execute"
)

const (
	RunProgram     = "/program_worker/run_program"
	BasePortNumber = 9090
)

type Args struct {
	Program      []byte
	Inputs       []byte
	OutputOffset int
	OutputSize   int
}

type Result struct {
	Output []byte
}
type ProgramWorker struct {
	erpc.CallCtx
}

func (pw *ProgramWorker) RunProgram(args *Args) (Result, *erpc.Status) {
	prgrmInBytes := args.Program
	prgrm := cxast.DeserializeCXProgramV2(prgrmInBytes)

	prgrm.Memory = cxast.MakeProgram().Memory
	injectMainInputs(prgrm, args.Inputs)
	err := cxexecute.RunCompiled(prgrm, 0, nil)
	if err != nil {
		return Result{}, erpc.NewStatus(1, fmt.Sprintf("%v", err))
	}

	byteOut := prgrm.Memory[args.OutputOffset : args.OutputOffset+args.OutputSize]
	res := Result{
		Output: byteOut,
	}

	return res, nil
}

// injectMainInputs injects `inps` at the beginning of `prgrm`'s memory,
// which should always represent the memory sent to the first expression contained
// in `prgrm`'s `main`'s function.
func injectMainInputs(prgrm *cxast.CXProgram, inps []byte) {
	for i := 0; i < len(inps); i++ {
		prgrm.Memory[i] = inps[i]
	}
}

func GetAvailableWorkers(numberOfAvailableWorkers int) []int {
	var workersAddr []int
	for i := 0; i < numberOfAvailableWorkers; i++ {
		workersAddr = append(workersAddr, BasePortNumber+i)
	}
	return workersAddr
}
