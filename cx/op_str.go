package base

import (
	"fmt"
)

func op_str_print(expr *CXExpression, stack *CXStack, fp int) {
	inp1 := expr.Inputs[0]
	fmt.Println(ReadStr(stack, fp, inp1))
}
