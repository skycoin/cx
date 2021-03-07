package cxcore

import (
	"fmt"
)

func opBoolPrint(expr *CXExpression, fp int) {
	fmt.Println(ReadBool(fp, expr.Inputs[0]))
}

func opBoolEqual(expr *CXExpression, fp int) {
	outV0 := ReadBool(fp, expr.Inputs[0]) == ReadBool(fp, expr.Inputs[1])
	WriteBool(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

func opBoolUnequal(expr *CXExpression, fp int) {
	outV0 := ReadBool(fp, expr.Inputs[0]) != ReadBool(fp, expr.Inputs[1])
	WriteBool(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

func opBoolNot(expr *CXExpression, fp int) {
	outV0 := !ReadBool(fp, expr.Inputs[0])
	WriteBool(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

func opBoolAnd(expr *CXExpression, fp int) {
	outV0 := ReadBool(fp, expr.Inputs[0]) && ReadBool(fp, expr.Inputs[1])
	WriteBool(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}

func opBoolOr(expr *CXExpression, fp int) {
	outV0 := ReadBool(fp, expr.Inputs[0]) || ReadBool(fp, expr.Inputs[1])
	WriteBool(GetFinalOffset(fp, expr.Outputs[0]), outV0)
}
