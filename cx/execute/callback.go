package execute

import (
	"os"

	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/constants"
	"github.com/skycoin/cx/cx/types"
)

//TODO: Define Function Pointers and deprecate Callback

// What calls Callback?ex
// Callback is only called from opHttpHandle, can probably be removed
// TODO: Delete and delete call from opHTTPHandle
// TODO: We probably dont need this? HTTPHandle can work in another way
//TODO: Is Callback actually "CallFunction" ?
func Callback(cxprogram *ast.CXProgram, fn *ast.CXFunction, inputs [][]byte) (outputs [][]byte) {
	line := cxprogram.CallStack[cxprogram.CallCounter].Line
	previousCall := cxprogram.CallCounter
	cxprogram.CallCounter++
	newCall := &cxprogram.CallStack[cxprogram.CallCounter]
	newCall.Operator = fn
	newCall.Line = 0
	newCall.FramePointer = cxprogram.Stack.Pointer
	cxprogram.Stack.Pointer += fn.Size
	newFP := newCall.FramePointer

	// wiping next mem frame (removing garbage)
	for c := types.Pointer(0); c < fn.Size; c++ {
		cxprogram.Memory[newFP+c] = 0
	}

	for i, inp := range inputs {
		types.WriteSlice_byte(cxprogram.Memory, ast.GetFinalOffset(cxprogram, newFP, &cxprogram.CXArgs[newCall.Operator.Inputs[i]]), inp)
	}

	var maxOps = 0

	//err := cxprogram.Run(true, maxOps, previousCall)
	err := RunCxAst(cxprogram, true, maxOps, previousCall)
	if err != nil {
		os.Exit(constants.CX_INTERNAL_ERROR)
	}

	cxprogram.CallCounter = previousCall
	cxprogram.CallStack[cxprogram.CallCounter].Line = line

	for _, outIdx := range fn.Outputs {
		out := &cxprogram.CXArgs[outIdx]
		// Making a copy of the bytes, so if we modify the bytes being held by `outputs`
		// we don't modify the program memory.
		mem := types.GetSlice_byte(cxprogram.Memory, ast.GetFinalOffset(cxprogram, newFP, out), ast.GetSize(cxprogram, out))
		cop := make([]byte, len(mem))
		copy(cop, mem)
		outputs = append(outputs, cop)
	}
	return outputs
}
