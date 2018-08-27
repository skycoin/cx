package base

import (
	// "fmt"
	"github.com/skycoin/skycoin/src/cipher/encoder"
)

func EscapeAnalysis (mem []byte, fp int, inpOffset, outOffset int, arg *CXArgument) {
	var heapOffset int
	if arg.HeapOffset > 0 {
		// then it's a reference to the symbol
		var off int32
		encoder.DeserializeAtomic(mem[fp+arg.HeapOffset:fp+arg.HeapOffset+TYPE_POINTER_SIZE], &off)

		if off > 0 {
			// non-nil, i.e. object is already allocated
			heapOffset = int(off)
		} else {
			// nil, needs to be allocated
			heapOffset = AllocateSeq(arg.Program, arg.TotalSize+OBJECT_HEADER_SIZE)
			o := GetFinalOffset(mem, fp, arg, MEM_WRITE)
			WriteMemory(mem, o, encoder.SerializeAtomic(int32(heapOffset)))
		}
	}

	byts := ReadMemory(mem, inpOffset, arg)
	// creating a header for this object
	size := encoder.SerializeAtomic(int32(len(byts)))

	var header []byte = make([]byte, OBJECT_HEADER_SIZE, OBJECT_HEADER_SIZE)
	for c := 5; c < OBJECT_HEADER_SIZE; c++ {
		header[c] = size[c-5]
	}

	obj := append(header, byts...)

	WriteMemory(mem, heapOffset, obj)

	off := encoder.SerializeAtomic(int32(heapOffset))

	WriteMemory(mem, outOffset, off)
	// WriteMemory(mem, outOffset, arg, off)
}

func op_identity(expr *CXExpression, mem []byte, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	inp1Offset := GetFinalOffset(mem, fp, inp1, MEM_READ)
	out1Offset := GetFinalOffset(mem, fp, out1, MEM_WRITE)

	if out1.DoesEscape {
		EscapeAnalysis(mem, fp, inp1Offset, out1Offset, inp1)
	} else {
		switch out1.PassBy {
		case PASSBY_VALUE:
			WriteMemory(mem, out1Offset, ReadMemory(mem, inp1Offset, inp1))
		case PASSBY_REFERENCE:
			WriteMemory(mem, out1Offset, encoder.SerializeAtomic(int32(inp1Offset)))
		}
	}
}

func op_goTo(expr *CXExpression, call *CXCall) {
	// inp1 := expr.Inputs[0]
	// call.Line = ReadI32(inp1)
}

func op_jmp(expr *CXExpression, mem []byte, fp int, call *CXCall) {
	// inp1, inp2, inp3 := expr.Inputs[0], expr.Inputs[1], expr.Inputs[2]
	inp1 := expr.Inputs[0]
	var predicate bool
	// var thenLines int32
	// var elseLines int32

	if expr.Label != "" {
		// then it's a goto
		call.Line = call.Line + expr.ThenLines
	} else {
		inp1Offset := GetFinalOffset(mem, fp, inp1, MEM_READ)

		predicateB := mem[inp1Offset : inp1Offset+inp1.Size]
		encoder.DeserializeAtomic(predicateB, &predicate)

		if predicate {
			call.Line = call.Line + expr.ThenLines
		} else {
			call.Line = call.Line + expr.ElseLines
		}
	}

}
