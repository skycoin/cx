package cxcore

import (
	"fmt"
)

func opBoolPrint(inputs []CXValue, outputs []CXValue) {
	fmt.Println(inputs[0].Get_bool())
}

func opBoolEqual(inputs []CXValue, outputs []CXValue) {
	outV0 := inputs[0].Get_bool() == inputs[1].Get_bool()
	outputs[0].Set_bool(outV0)
}

func opBoolUnequal(inputs []CXValue, outputs []CXValue) {
	outV0 := inputs[0].Get_bool() != inputs[1].Get_bool()
	outputs[0].Set_bool(outV0)
}

func opBoolNot(inputs []CXValue, outputs []CXValue) {
	outV0 := !inputs[0].Get_bool()
	outputs[0].Set_bool(outV0)
}

func opBoolAnd(inputs []CXValue, outputs []CXValue) {
	inpV0 := inputs[0].Get_bool()
	inpV1 := inputs[1].Get_bool()
	outputs[0].Set_bool(inpV0 && inpV1)
}

func opBoolOr(inputs []CXValue, outputs []CXValue) {
	inpV0 := inputs[0].Get_bool()
	inpV1 := inputs[1].Get_bool()
	outputs[0].Set_bool(inpV0 || inpV1)
}
