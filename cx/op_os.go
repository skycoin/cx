package base

import (
	// "fmt"
	"github.com/skycoin/skycoin/src/cipher/encoder"
)

func op_os_GetWorkingDirectory(expr *CXExpression, fp int) {
	out1 := expr.Outputs[0]
	out1Offset := GetFinalOffset(fp, out1)
	
	byts := encoder.Serialize(expr.Program.Path)

	size := encoder.Serialize(int32(len(byts)))
	heapOffset := AllocateSeq(len(byts)+OBJECT_HEADER_SIZE)
	
	var header []byte = make([]byte, OBJECT_HEADER_SIZE, OBJECT_HEADER_SIZE)
	for c := 5; c < OBJECT_HEADER_SIZE; c++ {
		header[c] = size[c-5]
	}

	obj := append(header, byts...)

	WriteMemory(heapOffset, obj)

	off := encoder.SerializeAtomic(int32(heapOffset))

	WriteMemory(out1Offset, off)
}
