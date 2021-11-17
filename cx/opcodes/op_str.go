package opcodes

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/constants"
)

func opStrToI8(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0, err := strconv.ParseInt(inputs[0].Get_str(prgrm), 10, 8)
	if err != nil {
		panic(constants.CX_RUNTIME_ERROR)
	}
	outputs[0].Set_i8(prgrm, int8(outV0))
}

func opStrToI16(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0, err := strconv.ParseInt(inputs[0].Get_str(prgrm), 10, 16)
	if err != nil {
		panic(constants.CX_RUNTIME_ERROR)
	}
	outputs[0].Set_i16(prgrm, int16(outV0))
}

func opStrToI32(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0, err := strconv.ParseInt(inputs[0].Get_str(prgrm), 10, 32)
	if err != nil {
		panic(constants.CX_RUNTIME_ERROR)
	}
	outputs[0].Set_i32(prgrm, int32(outV0))
}

func opStrToI64(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0, err := strconv.ParseInt(inputs[0].Get_str(prgrm), 10, 64)
	if err != nil {
		panic(constants.CX_RUNTIME_ERROR)
	}
	outputs[0].Set_i64(prgrm, int64(outV0))
}

func opStrToUI8(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0, err := strconv.ParseUint(inputs[0].Get_str(prgrm), 10, 8)
	if err != nil {
		panic(constants.CX_RUNTIME_ERROR)
	}
	outputs[0].Set_ui8(prgrm, uint8(outV0))
}

func opStrToUI16(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0, err := strconv.ParseUint(inputs[0].Get_str(prgrm), 10, 16)
	if err != nil {
		panic(constants.CX_RUNTIME_ERROR)
	}
	outputs[0].Set_ui16(prgrm, uint16(outV0))
}

func opStrToUI32(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0, err := strconv.ParseUint(inputs[0].Get_str(prgrm), 10, 32)
	if err != nil {
		panic(constants.CX_RUNTIME_ERROR)
	}
	outputs[0].Set_ui32(prgrm, uint32(outV0))
}

func opStrToUI64(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0, err := strconv.ParseUint(inputs[0].Get_str(prgrm), 10, 64)
	if err != nil {
		panic(constants.CX_RUNTIME_ERROR)
	}
	outputs[0].Set_ui64(prgrm, uint64(outV0))
}

func opStrToF32(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0, err := strconv.ParseFloat(inputs[0].Get_str(prgrm), 32)
	if err != nil {
		panic(constants.CX_RUNTIME_ERROR)
	}
	outputs[0].Set_f32(prgrm, float32(outV0))
}

func opStrToF64(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0, err := strconv.ParseFloat(inputs[0].Get_str(prgrm), 64)
	if err != nil {
		panic(constants.CX_RUNTIME_ERROR)
	}
	outputs[0].Set_f64(prgrm, float64(outV0))
}

func opStrEq(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_str(prgrm) == inputs[1].Get_str(prgrm)
	outputs[0].Set_bool(prgrm, outV0)
}

func opStrUneq(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_str(prgrm) != inputs[1].Get_str(prgrm)
	outputs[0].Set_bool(prgrm, outV0)
}

func opStrLt(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_str(prgrm) < inputs[1].Get_str(prgrm)
	outputs[0].Set_bool(prgrm, outV0)
}

func opStrLteq(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_str(prgrm) <= inputs[1].Get_str(prgrm)
	outputs[0].Set_bool(prgrm, outV0)
}

func opStrGt(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_str(prgrm) >= inputs[1].Get_str(prgrm)
	outputs[0].Set_bool(prgrm, outV0)
}

func opStrGteq(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_str(prgrm) >= inputs[1].Get_str(prgrm)
	outputs[0].Set_bool(prgrm, outV0)
}

func opStrPrint(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	fmt.Println(inputs[0].Get_str(prgrm))
}

func opStrConcat(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outputs[0].Set_str(prgrm, inputs[0].Get_str(prgrm)+inputs[1].Get_str(prgrm))
}

func opStrSubstr(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	str := inputs[0].Get_str(prgrm)
	begin := inputs[1].Get_i32(prgrm)
	end := inputs[2].Get_i32(prgrm)
	outputs[0].Set_str(prgrm, str[begin:end])
}

func opStrIndex(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	str := inputs[0].Get_str(prgrm)
	substr := inputs[1].Get_str(prgrm)
	outputs[0].Set_i32(prgrm, int32(strings.Index(str, substr)))
}

func opStrLastIndex(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	str := inputs[0].Get_str(prgrm)
	substr := inputs[1].Get_str(prgrm)
	outputs[0].Set_i32(prgrm, int32(strings.LastIndex(str, substr)))
}

func opStrTrimSpace(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	outputs[0].Set_str(prgrm, strings.TrimSpace(inputs[0].Get_str(prgrm)))
}
