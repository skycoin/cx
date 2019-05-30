package cxcore

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/amherag/skycoin/src/cipher/encoder"
)

func opStrToI8(expr *CXExpression, fp int) {
	outV0, err := strconv.ParseInt(ReadStr(fp, expr.Inputs[0]), 10, 8)
	if err != nil {
		panic(CX_RUNTIME_ERROR)
	}
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), FromI8(int8(outV0)))
}

func opStrToI16(expr *CXExpression, fp int) {
	outV0, err := strconv.ParseInt(ReadStr(fp, expr.Inputs[0]), 10, 16)
	if err != nil {
		panic(CX_RUNTIME_ERROR)
	}
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), FromI16(int16(outV0)))
}

func opStrToI32(expr *CXExpression, fp int) {
	outV0, err := strconv.ParseInt(ReadStr(fp, expr.Inputs[0]), 10, 32)
	if err != nil {
		panic(CX_RUNTIME_ERROR)
	}
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), FromI32(int32(outV0)))
}

func opStrToI64(expr *CXExpression, fp int) {
	outV0, err := strconv.ParseInt(ReadStr(fp, expr.Inputs[0]), 10, 64)
	if err != nil {
		panic(CX_RUNTIME_ERROR)
	}
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), FromI64(outV0))
}

func opStrToUI8(expr *CXExpression, fp int) {
	outV0, err := strconv.ParseUint(ReadStr(fp, expr.Inputs[0]), 10, 8)
	if err != nil {
		panic(CX_RUNTIME_ERROR)
	}
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), FromUI8(uint8(outV0)))
}

func opStrToUI16(expr *CXExpression, fp int) {
	outV0, err := strconv.ParseUint(ReadStr(fp, expr.Inputs[0]), 10, 16)
	if err != nil {
		panic(CX_RUNTIME_ERROR)
	}
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), FromUI16(uint16(outV0)))
}

func opStrToUI32(expr *CXExpression, fp int) {
	outV0, err := strconv.ParseUint(ReadStr(fp, expr.Inputs[0]), 10, 32)
	if err != nil {
		panic(CX_RUNTIME_ERROR)
	}
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), FromUI32(uint32(outV0)))
}

func opStrToUI64(expr *CXExpression, fp int) {
	outV0, err := strconv.ParseUint(ReadStr(fp, expr.Inputs[0]), 10, 64)
	if err != nil {
		panic(CX_RUNTIME_ERROR)
	}
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), FromUI64(outV0))
}

func opStrToF32(expr *CXExpression, fp int) {
	outV0, err := strconv.ParseFloat(ReadStr(fp, expr.Inputs[0]), 32)
	if err != nil {
		panic(CX_RUNTIME_ERROR)
	}
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), FromF32(float32(outV0)))
}

func opStrToF64(expr *CXExpression, fp int) {
	outV0, err := strconv.ParseFloat(ReadStr(fp, expr.Inputs[0]), 64)
	if err != nil {
		panic(CX_RUNTIME_ERROR)
	}
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), FromF64(outV0))
}

func opStrPrint(expr *CXExpression, fp int) {
	inp1 := expr.Inputs[0]
	fmt.Println(ReadStr(fp, inp1))
}

func opStrEq(expr *CXExpression, fp int) {
	outB0 := FromBool(ReadStr(fp, expr.Inputs[0]) == ReadStr(fp, expr.Inputs[1]))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opStrUneq(expr *CXExpression, fp int) {
	outB0 := FromBool(ReadStr(fp, expr.Inputs[0]) != ReadStr(fp, expr.Inputs[1]))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opStrLt(expr *CXExpression, fp int) {
	outB0 := FromBool(ReadStr(fp, expr.Inputs[0]) < ReadStr(fp, expr.Inputs[1]))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opStrLteq(expr *CXExpression, fp int) {
	outB0 := FromBool(ReadStr(fp, expr.Inputs[0]) <= ReadStr(fp, expr.Inputs[1]))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opStrGt(expr *CXExpression, fp int) {
	outB0 := FromBool(ReadStr(fp, expr.Inputs[0]) >= ReadStr(fp, expr.Inputs[1]))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opStrGteq(expr *CXExpression, fp int) {
	outB0 := FromBool(ReadStr(fp, expr.Inputs[0]) >= ReadStr(fp, expr.Inputs[1]))
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func writeString(fp int, str string, out *CXArgument) {

	byts := encoder.Serialize(str)
	size := encoder.Serialize(int32(len(byts)))
	heapOffset := AllocateSeq(len(byts) + OBJECT_HEADER_SIZE)

	var header = make([]byte, OBJECT_HEADER_SIZE)
	for c := 5; c < OBJECT_HEADER_SIZE; c++ {
		header[c] = size[c-5]
	}

	obj := append(header, byts...)

	WriteMemory(heapOffset, obj)

	off := encoder.SerializeAtomic(int32(heapOffset + OBJECT_HEADER_SIZE))

	WriteMemory(GetFinalOffset(fp, out), off)
}

func opStrConcat(expr *CXExpression, fp int) {
	writeString(fp, ReadStr(fp, expr.Inputs[0])+ReadStr(fp, expr.Inputs[1]), expr.Outputs[0])
}

func opStrSubstr(expr *CXExpression, fp int) {
	str := ReadStr(fp, expr.Inputs[0])
	begin := ReadI32(fp, expr.Inputs[1])
	end := ReadI32(fp, expr.Inputs[2])

	writeString(fp, str[begin:end], expr.Outputs[0])
}

func opStrIndex(expr *CXExpression, fp int) {
	str := ReadStr(fp, expr.Inputs[0])
	substr := ReadStr(fp, expr.Inputs[1])
	WriteMemory(GetFinalOffset(fp, expr.Outputs[0]), FromI32(int32(strings.Index(str, substr))))
}

func opStrTrimSpace(expr *CXExpression, fp int) {
	writeString(fp, strings.TrimSpace(ReadStr(fp, expr.Inputs[0])), expr.Outputs[0])
}
