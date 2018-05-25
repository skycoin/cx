package base

import (
	"fmt"
	// "github.com/skycoin/skycoin/src/cipher/encoder"
)


func op_test_value (expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, inp3 := expr.Inputs[0], expr.Inputs[1], expr.Inputs[2]

	byts1 := ReadMemory(stack, fp, inp1)
	byts2 := ReadMemory(stack, fp, inp2)

	var same bool
	same = true
	
	for i, byt := range byts1 {
		if byt != byts2[i] {
			same = false
		}
	}

	var message string
	message = ReadStr(stack, fp, inp3)
	
	if !same {
		if message != "" {
			fmt.Printf("%s: %d: result was not equal to the expected value\n; %s", expr.FileName, expr.FileLine, message)
		} else {
			fmt.Printf("%s: %d: result was not equal to the expected value\n", expr.FileName, expr.FileLine)
		}
	} else {
		
	}
}
