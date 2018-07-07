package base

import (
	"fmt"
	// "github.com/skycoin/skycoin/src/cipher/encoder"
)

func op_test_value(expr *CXExpression, stack *CXStack, fp int) {
	inp1, inp2, inp3 := expr.Inputs[0], expr.Inputs[1], expr.Inputs[2]

	var byts1, byts2 []byte
	
	if inp1.Type == TYPE_STR {
		byts1 = []byte(ReadStr(stack, fp, inp1))
		byts2 = []byte(ReadStr(stack, fp, inp2))
	} else {
		byts1 = ReadMemory(stack, GetFinalOffset(stack, fp, inp1, MEM_READ), inp1)
		byts2 = ReadMemory(stack, GetFinalOffset(stack, fp, inp2, MEM_READ), inp2)
	}

	var same bool
	same = true

	// fmt.Println("byts1", byts1)
	// fmt.Println("byts2", byts2)

	for i, byt := range byts1 {
		if byt != byts2[i] {
			same = false
		}
	}

	var message string
	message = ReadStr(stack, fp, inp3)

	if !same {
		if message != "" {
			fmt.Printf("%s: %d: result was not equal to the expected value; %s\n", expr.FileName, expr.FileLine, message)
		} else {
			fmt.Printf("%s: %d: result was not equal to the expected value\n", expr.FileName, expr.FileLine)
		}
	} else {

	}
}
