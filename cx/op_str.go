package cxcore

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/skycoin/skycoin/src/cipher/encoder"
)

func opStrToI8(inputs []CXValue, outputs []CXValue) {
	outV0, err := strconv.ParseInt(inputs[0].Get_str(), 10, 8)
	if err != nil {
		panic(CX_RUNTIME_ERROR)
	}
    outputs[0].Set_i8(int8(outV0))
}

func opStrToI16(inputs []CXValue, outputs []CXValue) {
	outV0, err := strconv.ParseInt(inputs[0].Get_str(), 10, 16)
	if err != nil {
		panic(CX_RUNTIME_ERROR)
	}
    outputs[0].Set_i16(int16(outV0))
}

func opStrToI32(inputs []CXValue, outputs []CXValue) {
	outV0, err := strconv.ParseInt(inputs[0].Get_str(), 10, 32)
	if err != nil {
		panic(CX_RUNTIME_ERROR)
	}
    outputs[0].Set_i32(int32(outV0))
}

func opStrToI64(inputs []CXValue, outputs []CXValue) {
	outV0, err := strconv.ParseInt(inputs[0].Get_str(), 10, 64)
	if err != nil {
		panic(CX_RUNTIME_ERROR)
	}
    outputs[0].Set_i64(int64(outV0))
}

func opStrToUI8(inputs []CXValue, outputs []CXValue) {
	outV0, err := strconv.ParseUint(inputs[0].Get_str(), 10, 8)
	if err != nil {
		panic(CX_RUNTIME_ERROR)
	}
    outputs[0].Set_ui8(uint8(outV0))
}

func opStrToUI16(inputs []CXValue, outputs []CXValue) {
	outV0, err := strconv.ParseUint(inputs[0].Get_str(), 10, 16)
	if err != nil {
		panic(CX_RUNTIME_ERROR)
	}
    outputs[0].Set_ui16(uint16(outV0))
}

func opStrToUI32(inputs []CXValue, outputs []CXValue) {
	outV0, err := strconv.ParseUint(inputs[0].Get_str(), 10, 32)
	if err != nil {
		panic(CX_RUNTIME_ERROR)
	}
    outputs[0].Set_ui32(uint32(outV0))
}

func opStrToUI64(inputs []CXValue, outputs []CXValue) {
	outV0, err := strconv.ParseUint(inputs[0].Get_str(), 10, 64)
	if err != nil {
		panic(CX_RUNTIME_ERROR)
	}
    outputs[0].Set_ui64(uint64(outV0))
}

func opStrToF32(inputs []CXValue, outputs []CXValue) {
	outV0, err := strconv.ParseFloat(inputs[0].Get_str(), 32)
	if err != nil {
		panic(CX_RUNTIME_ERROR)
	}
    outputs[0].Set_f32(float32(outV0))
}

func opStrToF64(inputs []CXValue, outputs []CXValue) {
	outV0, err := strconv.ParseFloat(inputs[0].Get_str(), 64)
	if err != nil {
		panic(CX_RUNTIME_ERROR)
	}
    outputs[0].Set_f64(float64(outV0))
}

func opStrEq(inputs []CXValue, outputs []CXValue) {
	outV0 := inputs[0].Get_str() == inputs[1].Get_str()
	outputs[0].Set_bool(outV0)
}

func opStrUneq(inputs []CXValue, outputs []CXValue) {
	outV0 := inputs[0].Get_str() != inputs[1].Get_str()
	outputs[0].Set_bool(outV0)
}

func opStrLt(inputs []CXValue, outputs []CXValue) {
	outV0 := inputs[0].Get_str() < inputs[1].Get_str()
	outputs[0].Set_bool(outV0)
}

func opStrLteq(inputs []CXValue, outputs []CXValue) {
	outV0 := inputs[0].Get_str() <= inputs[1].Get_str()
	outputs[0].Set_bool(outV0)
}

func opStrGt(inputs []CXValue, outputs []CXValue) {
	outV0 := inputs[0].Get_str() >= inputs[1].Get_str()
	outputs[0].Set_bool(outV0)
}

func opStrGteq(inputs []CXValue, outputs []CXValue) {
	outV0 := inputs[0].Get_str() >= inputs[1].Get_str()
	outputs[0].Set_bool(outV0)
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

func opStrPrint(inputs []CXValue, outputs []CXValue) {
	fmt.Println(inputs[0].Get_str())
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
