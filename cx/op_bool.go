package cxcore

import (
	"fmt"
	"github.com/skycoin/cx/cx/ast"
)

func opBoolPrint(inputs []ast.CXValue, outputs []ast.CXValue) {
	fmt.Println(inputs[0].Get_bool())
}

func opBoolEqual(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_bool() == inputs[1].Get_bool()
	outputs[0].Set_bool(outV0)
}

func opBoolUnequal(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := inputs[0].Get_bool() != inputs[1].Get_bool()
	outputs[0].Set_bool(outV0)
}

func opBoolNot(inputs []ast.CXValue, outputs []ast.CXValue) {
	outV0 := !inputs[0].Get_bool()
	outputs[0].Set_bool(outV0)
}

func opBoolAnd(inputs []ast.CXValue, outputs []ast.CXValue) {
	inpV0 := inputs[0].Get_bool()
	inpV1 := inputs[1].Get_bool()
	outputs[0].Set_bool(inpV0 && inpV1)
}

func opBoolOr(inputs []ast.CXValue, outputs []ast.CXValue) {
	inpV0 := inputs[0].Get_bool()
	inpV1 := inputs[1].Get_bool()
	outputs[0].Set_bool(inpV0 || inpV1)
}
