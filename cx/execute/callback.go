package execute

import (
	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/constants"
	"os"
)

// What calls Callback?ex
// Callback is only called from opHttpHandle, can probably be removed
// TODO: Delete and delete call from opHTTPHandle
// TODO: We probably dont need this? HTTPHandle can work in another way
func Callback(cxprogram *ast.CXProgram, fn *ast.CXFunction, inputs [][]byte) (outputs [][]byte) {
	line := cxprogram.CallStack[cxprogram.CallCounter].Line
	previousCall := cxprogram.CallCounter
	cxprogram.CallCounter++
	newCall := &cxprogram.CallStack[cxprogram.CallCounter]
	newCall.Operator = fn
	newCall.Line = 0
	newCall.FramePointer = cxprogram.StackPointer
	cxprogram.StackPointer += newCall.Operator.Size
	newFP := newCall.FramePointer

	// wiping next mem frame (removing garbage)
	for c := 0; c < fn.Size; c++ {
		cxprogram.Memory[newFP+c] = 0
	}

	for i, inp := range inputs {
		ast.WriteMemory(ast.GetFinalOffset(newFP, newCall.Operator.Inputs[i]), inp)
	}

	var nCalls = 0

	//err := cxprogram.Run(true, &nCalls, previousCall)
	err := RunCxAst(cxprogram, true, &nCalls, previousCall)
	if err != nil {
		os.Exit(constants.CX_INTERNAL_ERROR)
	}

	cxprogram.CallCounter = previousCall
	cxprogram.CallStack[cxprogram.CallCounter].Line = line

	for _, out := range fn.Outputs {
		// Making a copy of the bytes, so if we modify the bytes being held by `outputs`
		// we don't modify the program memory.
		mem := ast.ReadMemory(ast.GetFinalOffset(newFP, out), out)
		cop := make([]byte, len(mem))
		copy(cop, mem)
		outputs = append(outputs, cop)
	}
	return outputs
}
