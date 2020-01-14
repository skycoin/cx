package cxcore

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/SkycoinProject/skycoin/src/cipher/encoder"
)

func opStrToI8(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	outV0, err := strconv.ParseInt(ReadStr(fp, expr.Inputs[0]), 10, 8)
	if err != nil {
		panic(CX_RUNTIME_ERROR)
	}
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), FromI8(int8(outV0)))
}

func opStrToI16(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	outV0, err := strconv.ParseInt(ReadStr(fp, expr.Inputs[0]), 10, 16)
	if err != nil {
		panic(CX_RUNTIME_ERROR)
	}
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), FromI16(int16(outV0)))
}

func opStrToI32(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	outV0, err := strconv.ParseInt(ReadStr(fp, expr.Inputs[0]), 10, 32)
	if err != nil {
		panic(CX_RUNTIME_ERROR)
	}
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), FromI32(int32(outV0)))
}

func opStrToI64(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	outV0, err := strconv.ParseInt(ReadStr(fp, expr.Inputs[0]), 10, 64)
	if err != nil {
		panic(CX_RUNTIME_ERROR)
	}
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), FromI64(outV0))
}

func opStrToUI8(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	outV0, err := strconv.ParseUint(ReadStr(fp, expr.Inputs[0]), 10, 8)
	if err != nil {
		panic(CX_RUNTIME_ERROR)
	}
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), FromUI8(uint8(outV0)))
}

func opStrToUI16(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	outV0, err := strconv.ParseUint(ReadStr(fp, expr.Inputs[0]), 10, 16)
	if err != nil {
		panic(CX_RUNTIME_ERROR)
	}
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), FromUI16(uint16(outV0)))
}

func opStrToUI32(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	outV0, err := strconv.ParseUint(ReadStr(fp, expr.Inputs[0]), 10, 32)
	if err != nil {
		panic(CX_RUNTIME_ERROR)
	}
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), FromUI32(uint32(outV0)))
}

func opStrToUI64(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	outV0, err := strconv.ParseUint(ReadStr(fp, expr.Inputs[0]), 10, 64)
	if err != nil {
		panic(CX_RUNTIME_ERROR)
	}
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), FromUI64(outV0))
}

func opStrToF32(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	outV0, err := strconv.ParseFloat(ReadStr(fp, expr.Inputs[0]), 32)
	if err != nil {
		panic(CX_RUNTIME_ERROR)
	}
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), FromF32(float32(outV0)))
}

func opStrToF64(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	outV0, err := strconv.ParseFloat(ReadStr(fp, expr.Inputs[0]), 64)
	if err != nil {
		panic(CX_RUNTIME_ERROR)
	}
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), FromF64(outV0))
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

	outB0 := FromBool(ReadStr(fp, expr.Inputs[0]) == ReadStr(fp, expr.Inputs[1]))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opStrUneq(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	outB0 := FromBool(ReadStr(fp, expr.Inputs[0]) != ReadStr(fp, expr.Inputs[1]))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opStrLt(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	outB0 := FromBool(ReadStr(fp, expr.Inputs[0]) < ReadStr(fp, expr.Inputs[1]))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opStrLteq(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	outB0 := FromBool(ReadStr(fp, expr.Inputs[0]) <= ReadStr(fp, expr.Inputs[1]))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opStrGt(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	outB0 := FromBool(ReadStr(fp, expr.Inputs[0]) >= ReadStr(fp, expr.Inputs[1]))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opStrGteq(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	outB0 := FromBool(ReadStr(fp, expr.Inputs[0]) >= ReadStr(fp, expr.Inputs[1]))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
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
