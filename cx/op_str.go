package cxcore

import (
	"fmt"
	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/constants"
	"strconv"
	"strings"
)

func opStrToI8(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0, err := strconv.ParseInt(inputs[0].Get_str(), 10, 8)
	if err != nil {
		panic(constants.CX_RUNTIME_ERROR)
	}
    outputs[0].Set_i8(int8(outV0))
}

func opStrToI16(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0, err := strconv.ParseInt(inputs[0].Get_str(), 10, 16)
	if err != nil {
		panic(constants.CX_RUNTIME_ERROR)
	}
    outputs[0].Set_i16(int16(outV0))
}

func opStrToI32(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0, err := strconv.ParseInt(inputs[0].Get_str(), 10, 32)
	if err != nil {
		panic(constants.CX_RUNTIME_ERROR)
	}
    outputs[0].Set_i32(int32(outV0))
}

func opStrToI64(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0, err := strconv.ParseInt(inputs[0].Get_str(), 10, 64)
	if err != nil {
		panic(constants.CX_RUNTIME_ERROR)
	}
    outputs[0].Set_i64(int64(outV0))
}

func opStrToUI8(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0, err := strconv.ParseUint(inputs[0].Get_str(), 10, 8)
	if err != nil {
		panic(constants.CX_RUNTIME_ERROR)
	}
    outputs[0].Set_ui8(uint8(outV0))
}

func opStrToUI16(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0, err := strconv.ParseUint(inputs[0].Get_str(), 10, 16)
	if err != nil {
		panic(constants.CX_RUNTIME_ERROR)
	}
    outputs[0].Set_ui16(uint16(outV0))
}

func opStrToUI32(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0, err := strconv.ParseUint(inputs[0].Get_str(), 10, 32)
	if err != nil {
		panic(constants.CX_RUNTIME_ERROR)
	}
    outputs[0].Set_ui32(uint32(outV0))
}

func opStrToUI64(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0, err := strconv.ParseUint(inputs[0].Get_str(), 10, 64)
	if err != nil {
		panic(constants.CX_RUNTIME_ERROR)
	}
    outputs[0].Set_ui64(uint64(outV0))
}

func opStrToF32(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0, err := strconv.ParseFloat(inputs[0].Get_str(), 32)
	if err != nil {
		panic(constants.CX_RUNTIME_ERROR)
	}
    outputs[0].Set_f32(float32(outV0))
}

func opStrToF64(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0, err := strconv.ParseFloat(inputs[0].Get_str(), 64)
	if err != nil {
		panic(constants.CX_RUNTIME_ERROR)
	}
    outputs[0].Set_f64(float64(outV0))
}

func opStrEq(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_str() == inputs[1].Get_str()
	outputs[0].Set_bool(outV0)
}

func opStrUneq(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_str() != inputs[1].Get_str()
	outputs[0].Set_bool(outV0)
}

func opStrLt(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_str() < inputs[1].Get_str()
	outputs[0].Set_bool(outV0)
}

func opStrLteq(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_str() <= inputs[1].Get_str()
	outputs[0].Set_bool(outV0)
}

func opStrGt(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_str() >= inputs[1].Get_str()
	outputs[0].Set_bool(outV0)
}

func opStrGteq(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_str() >= inputs[1].Get_str()
	outputs[0].Set_bool(outV0)
}

func opStrPrint(inputs []ast.CXValue, outputs []ast.CXValue) {
	fmt.Println(inputs[0].Get_str())
}

func opStrConcat(inputs []ast.CXValue, outputs []ast.CXValue) {
	outputs[0].Set_str(inputs[0].Get_str() + inputs[1].Get_str())
}

func opStrSubstr(inputs []ast.CXValue, outputs []ast.CXValue) {
	str := inputs[0].Get_str()
	begin := inputs[1].Get_i32()
	end := inputs[2].Get_i32()
	outputs[0].Set_str(str[begin:end])
}

func opStrIndex(inputs []ast.CXValue, outputs []ast.CXValue) {
	str := inputs[0].Get_str()
	substr := inputs[1].Get_str()
	outputs[0].Set_i32(int32(strings.Index(str, substr)))
}

func opStrLastIndex(inputs []ast.CXValue, outputs []ast.CXValue) {
	str := inputs[0].Get_str()
	substr := inputs[1].Get_str()
	outputs[0].Set_i32(int32(strings.LastIndex(str, substr)))
}

func opStrTrimSpace(inputs []ast.CXValue, outputs []ast.CXValue) {
	outputs[0].Set_str(strings.TrimSpace(inputs[0].Get_str()))
}
