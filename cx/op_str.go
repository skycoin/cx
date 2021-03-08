package cxcore

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/skycoin/skycoin/src/cipher/encoder"
)

func opStrToI8(expr *CXExpression, fp int) {
	outV0, err := strconv.ParseInt(ReadStr(fp, expr.Inputs[0]), 10, 8)
	if err != nil {
		panic(CX_RUNTIME_ERROR)
	}
	WriteI8(GetOffset_str(fp, expr.Outputs[0]), int8(outV0))
}

func opStrToI16(expr *CXExpression, fp int) {
	outV0, err := strconv.ParseInt(ReadStr(fp, expr.Inputs[0]), 10, 16)
	if err != nil {
		panic(CX_RUNTIME_ERROR)
	}
	WriteI16(GetOffset_str(fp, expr.Outputs[0]), int16(outV0))
}

func opStrToI32(expr *CXExpression, fp int) {
	outV0, err := strconv.ParseInt(ReadStr(fp, expr.Inputs[0]), 10, 32)
	if err != nil {
		panic(CX_RUNTIME_ERROR)
	}
	WriteI32(GetOffset_str(fp, expr.Outputs[0]), int32(outV0))
}

func opStrToI64(expr *CXExpression, fp int) {
	outV0, err := strconv.ParseInt(ReadStr(fp, expr.Inputs[0]), 10, 64)
	if err != nil {
		panic(CX_RUNTIME_ERROR)
	}
	WriteI64(GetOffset_str(fp, expr.Outputs[0]), outV0)
}

func opStrToUI8(expr *CXExpression, fp int) {
	outV0, err := strconv.ParseUint(ReadStr(fp, expr.Inputs[0]), 10, 8)
	if err != nil {
		panic(CX_RUNTIME_ERROR)
	}
	WriteUI8(GetOffset_str(fp, expr.Outputs[0]), uint8(outV0))
}

func opStrToUI16(expr *CXExpression, fp int) {
	outV0, err := strconv.ParseUint(ReadStr(fp, expr.Inputs[0]), 10, 16)
	if err != nil {
		panic(CX_RUNTIME_ERROR)
	}
	WriteUI16(GetOffset_str(fp, expr.Outputs[0]), uint16(outV0))
}

func opStrToUI32(expr *CXExpression, fp int) {
	outV0, err := strconv.ParseUint(ReadStr(fp, expr.Inputs[0]), 10, 32)
	if err != nil {
		panic(CX_RUNTIME_ERROR)
	}
	WriteUI32(GetOffset_str(fp, expr.Outputs[0]), uint32(outV0))
}

func opStrToUI64(expr *CXExpression, fp int) {
	outV0, err := strconv.ParseUint(ReadStr(fp, expr.Inputs[0]), 10, 64)
	if err != nil {
		panic(CX_RUNTIME_ERROR)
	}
	WriteUI64(GetOffset_str(fp, expr.Outputs[0]), outV0)
}

func opStrToF32(expr *CXExpression, fp int) {
	outV0, err := strconv.ParseFloat(ReadStr(fp, expr.Inputs[0]), 32)
	if err != nil {
		panic(CX_RUNTIME_ERROR)
	}
	WriteF32(GetOffset_str(fp, expr.Outputs[0]), float32(outV0))
}

func opStrToF64(expr *CXExpression, fp int) {
	outV0, err := strconv.ParseFloat(ReadStr(fp, expr.Inputs[0]), 64)
	if err != nil {
		panic(CX_RUNTIME_ERROR)
	}
	WriteF64(GetOffset_str(fp, expr.Outputs[0]), outV0)
}

func opStrPrint(expr *CXExpression, fp int) {
	inp1 := expr.Inputs[0]
	fmt.Println(ReadStr(fp, inp1))
}

func opStrEq(expr *CXExpression, fp int) {
	outB0 := ReadStr(fp, expr.Inputs[0]) == ReadStr(fp, expr.Inputs[1])
	WriteBool(GetOffset_str(fp, expr.Outputs[0]), outB0)
}

func opStrUneq(expr *CXExpression, fp int) {
	outB0 := ReadStr(fp, expr.Inputs[0]) != ReadStr(fp, expr.Inputs[1])
	WriteBool(GetOffset_str(fp, expr.Outputs[0]), outB0)
}

func opStrLt(expr *CXExpression, fp int) {
	outB0 := ReadStr(fp, expr.Inputs[0]) < ReadStr(fp, expr.Inputs[1])
	WriteBool(GetOffset_str(fp, expr.Outputs[0]), outB0)
}

func opStrLteq(expr *CXExpression, fp int) {
	outB0 := ReadStr(fp, expr.Inputs[0]) <= ReadStr(fp, expr.Inputs[1])
	WriteBool(GetOffset_str(fp, expr.Outputs[0]), outB0)
}

func opStrGt(expr *CXExpression, fp int) {
	outB0 := ReadStr(fp, expr.Inputs[0]) >= ReadStr(fp, expr.Inputs[1])
	WriteBool(GetOffset_str(fp, expr.Outputs[0]), outB0)
}

func opStrGteq(expr *CXExpression, fp int) {
	outB0 := ReadStr(fp, expr.Inputs[0]) >= ReadStr(fp, expr.Inputs[1])
	WriteBool(GetOffset_str(fp, expr.Outputs[0]), outB0)
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
	WriteI32(GetOffset_str(fp, out), int32(heapOffset))
}

func opStrConcat(expr *CXExpression, fp int) {
	WriteString(fp, ReadStr(fp, expr.Inputs[0])+ReadStr(fp, expr.Inputs[1]), expr.Outputs[0])
}

func opStrSubstr(expr *CXExpression, fp int) {
	str := ReadStr(fp, expr.Inputs[0])
	begin := ReadI32(fp, expr.Inputs[1])
	end := ReadI32(fp, expr.Inputs[2])

	WriteString(fp, str[begin:end], expr.Outputs[0])
}

func opStrIndex(expr *CXExpression, fp int) {
	str := ReadStr(fp, expr.Inputs[0])
	substr := ReadStr(fp, expr.Inputs[1])
	WriteI32(GetOffset_str(fp, expr.Outputs[0]), int32(strings.Index(str, substr)))
}

func opStrLastIndex(expr *CXExpression, fp int) {
	WriteI32(GetOffset_str(fp, expr.Outputs[0]), int32(strings.LastIndex(ReadStr(fp, expr.Inputs[0]), ReadStr(fp, expr.Inputs[1]))))
}

func opStrTrimSpace(expr *CXExpression, fp int) {
	WriteString(fp, strings.TrimSpace(ReadStr(fp, expr.Inputs[0])), expr.Outputs[0])
}
