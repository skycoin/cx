package cxcore

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/SkycoinProject/skycoin/src/cipher/encoder"
)

func opStrStr(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	out1Offset := GetFinalOffset(fp, out1)

	switch out1.Type {
	case TYPE_BYTE:
		b, err := strconv.ParseInt(ReadStr(fp, inp1), 10, 8)
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

func opStrPrint(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	inp1 := expr.Inputs[0]
	fmt.Println(ReadStr(fp, inp1))
}

func opStrEq(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	inp1, inp2, out1 := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]
	outB1 := FromBool(ReadStr(fp, inp1) == ReadStr(fp, inp2))
	WriteMemory(GetFinalOffset(fp, out1), outB1)
}

func writeString(fp int, str string, out *CXArgument) {
	if str == "" {
		return
	}
	byts := encoder.Serialize(str)
	size := encoder.Serialize(int32(len(byts)) + OBJECT_HEADER_SIZE)
	heapOffset := AllocateSeq(len(byts) + OBJECT_HEADER_SIZE)

	var header = make([]byte, OBJECT_HEADER_SIZE)
	for c := 5; c < OBJECT_HEADER_SIZE; c++ {
		header[c] = size[c-5]
	}

	obj := append(header, byts...)

	WriteMemory(heapOffset, obj)

	off := encoder.SerializeAtomic(int32(heapOffset))

	WriteMemory(GetFinalOffset(fp, out), off)
}

func opStrConcat(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	writeString(fp, ReadStr(fp, expr.Inputs[0])+ReadStr(fp, expr.Inputs[1]), expr.Outputs[0])
}

func opStrSubstr(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	str := ReadStr(fp, expr.Inputs[0])
	begin := ReadI32(fp, expr.Inputs[1])
	end := ReadI32(fp, expr.Inputs[2])

	writeString(fp, str[begin:end], expr.Outputs[0])
}

func opStrIndex(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	str := ReadStr(fp, expr.Inputs[0])
	substr := ReadStr(fp, expr.Inputs[1])
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), FromI32(int32(strings.Index(str, substr))))
}

func opStrTrimSpace(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	writeString(fp, strings.TrimSpace(ReadStr(fp, expr.Inputs[0])), expr.Outputs[0])
}
