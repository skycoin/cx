package base

import (
	// "fmt"
	"github.com/skycoin/skycoin/src/cipher/encoder"
)

func op_os_GetWorkingDirectory(expr *CXExpression, mem []byte, fp int) {
	out1 := expr.Outputs[0]
	out1Offset := GetFinalOffset(mem, fp, out1, MEM_WRITE)
	
	byts := encoder.Serialize(expr.Program.Path)

	size := encoder.Serialize(int32(len(byts)))
	heapOffset := AllocateSeq(expr.Program, len(byts)+OBJECT_HEADER_SIZE)
	
	var header []byte = make([]byte, OBJECT_HEADER_SIZE, OBJECT_HEADER_SIZE)
	for c := 5; c < OBJECT_HEADER_SIZE; c++ {
		header[c] = size[c-5]
	}

	obj := append(header, byts...)

	WriteMemory(mem, heapOffset, obj)

	off := encoder.SerializeAtomic(int32(heapOffset))

	WriteMemory(mem, out1Offset, off)
}
