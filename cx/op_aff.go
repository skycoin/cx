package base

import (
	// "fmt"
	"github.com/skycoin/skycoin/src/cipher/encoder"
)

func op_aff_print(expr *CXExpression, mem []byte, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]

	_ = inp1
	_ = out1

	inp1Offset := GetFinalOffset(mem, fp, inp1, MEM_READ)

	_ = inp1Offset
}

func op_aff_query(expr *CXExpression, mem []byte, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]

	inp1Offset := GetFinalOffset(mem, fp, inp1, MEM_READ)
	inp2Offset := GetFinalOffset(mem, fp, inp2, MEM_READ)

	_ = inp2Offset

	_ = out1

	var sliceHeader []byte
	
	var len1 int32
	sliceHeader = mem[inp1Offset-SLICE_HEADER_SIZE : inp1Offset]
	encoder.DeserializeAtomic(sliceHeader[:4], &len1)

	// var obj1 []byte
	// obj1 = mem.Program.Heap.Heap[inp1Offset : int32(inp1Offset) + len1*int32(inp2.TotalSize)]
}
