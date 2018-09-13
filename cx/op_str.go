package base

import (
	"fmt"
	"strconv"
	"github.com/skycoin/skycoin/src/cipher/encoder"
)

func op_str_str(expr *CXExpression, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	out1Offset := GetFinalOffset(fp, out1)

	switch out1.Type {
	case TYPE_BYTE:
		b, err := strconv.ParseInt(ReadStr(fp, inp1), 10, 1)
		if err != nil {
			panic("")
		}
		WriteMemory(out1Offset, encoder.Serialize(b))
	case TYPE_STR:
		WriteObject(out1Offset, []byte(ReadStr(fp, inp1)))
	case TYPE_I32:
		i, err := strconv.ParseInt(ReadStr(fp, inp1), 10, 32)
		if err != nil {
			panic("")
		}
		WriteMemory(out1Offset, encoder.SerializeAtomic(i))
	case TYPE_I64:
		l, err := strconv.ParseInt(ReadStr(fp, inp1), 10, 64)
		if err != nil {
			panic("")
		}
		WriteMemory(out1Offset, encoder.Serialize(l))
	case TYPE_F32:
		f, err := strconv.ParseFloat(ReadStr(fp, inp1), 32)
		if err != nil {
			panic("")
		}
		WriteMemory(out1Offset, encoder.Serialize(float32(f)))
	case TYPE_F64:
		d, err := strconv.ParseFloat(ReadStr(fp, inp1), 64)
		if err != nil {
			panic("")
		}
		WriteMemory(out1Offset, encoder.Serialize(d))
	}
}

func op_str_print(expr *CXExpression, fp int) {
	inp1 := expr.Inputs[0]
	fmt.Println(ReadStr(fp, inp1))
}

func op_str_eq(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadStr(fp, inp1) == ReadStr(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func op_str_concat(expr *CXExpression, fp int) {
	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	out1Offset := GetFinalOffset(fp, out1)

	inp1Str := ReadStr(fp, inp1)
	inp2Str := ReadStr(fp, inp2)

	byts := encoder.Serialize(inp1Str + inp2Str)
	size := encoder.Serialize(int32(len(byts)))
	heapOffset := AllocateSeq(len(byts) + OBJECT_HEADER_SIZE)
	
	var header []byte = make([]byte, OBJECT_HEADER_SIZE, OBJECT_HEADER_SIZE)
	for c := 5; c < OBJECT_HEADER_SIZE; c++ {
		header[c] = size[c-5]
	}

	obj := append(header, byts...)

	WriteMemory(heapOffset, obj)

	off := encoder.SerializeAtomic(int32(heapOffset + OBJECT_HEADER_SIZE))

	WriteMemory(out1Offset, off)
}
