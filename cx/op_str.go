package cxcore

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/skycoin/skycoin/src/cipher/encoder"
)

func opStrToI8(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	outV0, err := strconv.ParseInt(ReadStr(fp, expr.Inputs[0]), 10, 8)
	if err != nil {
		panic(CX_RUNTIME_ERROR)
	}
	WriteI8(GetFinalOffset(fp, expr.Outputs[0]), int8(outV0))
}

func opStrToI16(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	outV0, err := strconv.ParseInt(ReadStr(fp, expr.Inputs[0]), 10, 16)
	if err != nil {
		panic(CX_RUNTIME_ERROR)
	}
	WriteI16(GetFinalOffset(fp, expr.Outputs[0]), int16(outV0))
}

func opStrToI32(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	outV0, err := strconv.ParseInt(ReadStr(fp, expr.Inputs[0]), 10, 32)
	if err != nil {
		panic(CX_RUNTIME_ERROR)
	}
	WriteI32(GetFinalOffset(fp, expr.Outputs[0]), int32(outV0))
}

func opStrToI64(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	outV0, err := strconv.ParseInt(ReadStr(fp, expr.Inputs[0]), 10, 64)
	if err != nil {
		panic(CX_RUNTIME_ERROR)
	}
	WriteI64(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

func opStrToUI8(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	outV0, err := strconv.ParseUint(ReadStr(fp, expr.Inputs[0]), 10, 8)
	if err != nil {
		panic(CX_RUNTIME_ERROR)
	}
	WriteUI8(GetFinalOffset(fp, expr.Outputs[0]), uint8(outV0))
}

func opStrToUI16(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	outV0, err := strconv.ParseUint(ReadStr(fp, expr.Inputs[0]), 10, 16)
	if err != nil {
		panic(CX_RUNTIME_ERROR)
	}
	WriteUI16(GetFinalOffset(fp, expr.Outputs[0]), uint16(outV0))
}

func opStrToUI32(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	outV0, err := strconv.ParseUint(ReadStr(fp, expr.Inputs[0]), 10, 32)
	if err != nil {
		panic(CX_RUNTIME_ERROR)
	}
	WriteUI32(GetFinalOffset(fp, expr.Outputs[0]), uint32(outV0))
}

func opStrToUI64(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	outV0, err := strconv.ParseUint(ReadStr(fp, expr.Inputs[0]), 10, 64)
	if err != nil {
		panic(CX_RUNTIME_ERROR)
	}
	WriteUI64(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

func opStrToF32(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	outV0, err := strconv.ParseFloat(ReadStr(fp, expr.Inputs[0]), 32)
	if err != nil {
		panic(CX_RUNTIME_ERROR)
	}
	WriteF32(GetFinalOffset(fp, expr.Outputs[0]), float32(outV0))
}

func opStrToF64(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	outV0, err := strconv.ParseFloat(ReadStr(fp, expr.Inputs[0]), 64)
	if err != nil {
		panic(CX_RUNTIME_ERROR)
	}
	WriteF64(GetFinalOffset(fp, expr.Outputs[0]), outV0)
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

	outB0 := ReadStr(fp, expr.Inputs[0]) == ReadStr(fp, expr.Inputs[1])
	WriteBool(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opStrUneq(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	outB0 := ReadStr(fp, expr.Inputs[0]) != ReadStr(fp, expr.Inputs[1])
	WriteBool(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opStrLt(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	outB0 := ReadStr(fp, expr.Inputs[0]) < ReadStr(fp, expr.Inputs[1])
	WriteBool(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opStrLteq(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	outB0 := ReadStr(fp, expr.Inputs[0]) <= ReadStr(fp, expr.Inputs[1])
	WriteBool(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opStrGt(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	outB0 := ReadStr(fp, expr.Inputs[0]) >= ReadStr(fp, expr.Inputs[1])
	WriteBool(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

func opStrGteq(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	outB0 := ReadStr(fp, expr.Inputs[0]) >= ReadStr(fp, expr.Inputs[1])
	WriteBool(GetFinalOffset(fp, expr.Outputs[0]), outB0)
}

// WriteString writes the string `str` on memory, starting at byte number `fp`.
func WriteString(fp int, str string, out *CXArgument) {

	byts := encoder.Serialize(str)
	size := len(byts) + OBJECT_HEADER_SIZE
	heapOffset := AllocateSeq(size)

	var header = make([]byte, OBJECT_HEADER_SIZE)
	WriteMemI32(header, 5, int32(size))
	obj := append(header, byts...)

	WriteMemory(heapOffset, obj)
	WriteI32(GetFinalOffset(fp, out), int32(heapOffset))
}

func opStrConcat(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	WriteString(fp, ReadStr(fp, expr.Inputs[0])+ReadStr(fp, expr.Inputs[1]), expr.Outputs[0])
}

func opStrSubstr(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	str := ReadStr(fp, expr.Inputs[0])
	begin := ReadI32(fp, expr.Inputs[1])
	end := ReadI32(fp, expr.Inputs[2])

	WriteString(fp, str[begin:end], expr.Outputs[0])
}

func opStrIndex(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	str := ReadStr(fp, expr.Inputs[0])
	substr := ReadStr(fp, expr.Inputs[1])
	WriteI32(GetFinalOffset(fp, expr.Outputs[0]), int32(strings.Index(str, substr)))
}

func opStrLastIndex(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	WriteI32(GetFinalOffset(fp, expr.Outputs[0]), int32(strings.LastIndex(ReadStr(fp, expr.Inputs[0]), ReadStr(fp, expr.Inputs[1]))))
}

func opStrTrimSpace(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	WriteString(fp, strings.TrimSpace(ReadStr(fp, expr.Inputs[0])), expr.Outputs[0])
}
